package service

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/types"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/cloud/alibaba"
	"github.com/galaxy-future/BridgX/pkg/cloud/huawei"
)

var clientMap sync.Map

func ExpandAndRepair(c *types.ClusterInfo, num int, taskId int64) ([]string, error) {
	tags := []cloud.Tag{{
		Key:   cloud.TaskId,
		Value: strconv.FormatInt(taskId, 10),
	},
		{
			Key:   cloud.ClusterName,
			Value: c.Name,
		}}
	expandInstanceIds := make([]string, 0)
	needExpandNum := num
	var err error
	var ids []string
	for k := 0; k < constants.Retry; k++ {
		ids, err = Expand(c, tags, needExpandNum)
		if err != nil {
			logs.Logger.Errorf("[ExpandCLuster] Expand retry error, times: %d, error: %s", k, err.Error())
		}
		expandInstanceIds = append(expandInstanceIds, ids...)
		if len(expandInstanceIds) == num {
			break
		}
		needExpandNum -= len(ids)
	}
	if len(expandInstanceIds) != num {
		logs.Logger.Infof("len(expandInstanceIds) %d, num %d", len(expandInstanceIds), num)
		_ = RepairCluster(c, taskId, expandInstanceIds)
	}
	return expandInstanceIds, err
}

func RepairCluster(c *types.ClusterInfo, taskId int64, instanceIds []string) (err error) {
	tags := []cloud.Tag{{
		Key:   cloud.TaskId,
		Value: strconv.FormatInt(taskId, 10),
	}}
	cloudInstances, err := GetInstanceByTag(c, tags)
	logs.Logger.Infof("[RepairCluster] GetInstanceByTag length %d", len(cloudInstances))
	if err != nil {
		return
	}
	cloudIds := make([]string, 0)
	for _, instance := range cloudInstances {
		cloudIds = append(cloudIds, instance.Id)
	}
	onlyCouldIds, onlyMemoryIds := cloudDiff(cloudIds, instanceIds)
	logs.Logger.Infof("[RepairCluster] taskId: %d, ClusterName: %s, Shrink InstanceIds: %v", taskId, c.Name, onlyCouldIds)
	shrink := func(attempt uint) error {
		return Shrink(c, onlyCouldIds)
	}
	err = retry.Retry(shrink, strategy.Limit(3), strategy.Backoff(backoff.BinaryExponential(10*time.Millisecond)))
	if err != nil {
		logs.Logger.Errorf("[RepairCluster] taskId: %d, ClusterName: %s, Shrink InstanceIds error: %s", taskId, c.Name, err.Error())
	}

	logs.Logger.Infof("[RepairCluster] taskId: %d, ClusterName: %s, UpdateDB InstanceIds: %v", taskId, c.Name, onlyMemoryIds)
	update := func(attempt uint) error {
		if len(onlyMemoryIds) == 0 {
			return nil
		}
		now := time.Now()
		return model.BatchUpdateByInstanceIds(onlyMemoryIds, model.Instance{
			DeleteAt: &now,
			Status:   constants.Deleted,
		})
	}
	err = retry.Retry(update, strategy.Limit(3), strategy.Backoff(backoff.BinaryExponential(10*time.Millisecond)))
	if err != nil {
		logs.Logger.Errorf("[RepairCluster] taskId: %d, ClusterName: %s, UpdateDB InstanceIds error: %s", taskId, c.Name, err.Error())
	}

	return nil
}

func cloudDiff(cloudIds, memoryIds []string) (onlyCouldIds, onlyMemoryIds []string) {
	if len(cloudIds) == 0 {
		return onlyCouldIds, memoryIds
	}
	if len(memoryIds) == 0 {
		return cloudIds, onlyMemoryIds
	}
	tmpMap := make(map[string]int, 0)
	for _, id := range cloudIds {
		tmpMap[id] += 1
	}
	for _, id := range memoryIds {
		tmpMap[id] += 2
	}
	for k, v := range tmpMap {
		if v == 1 {
			onlyCouldIds = append(onlyCouldIds, k)
		}
		if v == 2 {
			onlyMemoryIds = append(onlyMemoryIds, k)
		}
	}
	return
}

func Expand(clusterInfo *types.ClusterInfo, tags []cloud.Tag, num int) (instanceIds []string, err error) {
	batch := getBatch(num, constants.BatchMax)
	createdBatch := make(chan []string, batch)
	createdError := make(chan error, batch)
	cur := num
	provider, err := getProvider(clusterInfo.Provider, clusterInfo.AccountKey, clusterInfo.RegionId)
	if err != nil {
		return nil, err
	}
	params, err := generateParams(clusterInfo, tags)
	for ; cur > 0; cur -= constants.BatchMax {
		go func(cur int) {
			var bErr error
			defer func() {
				if bErr := recover(); bErr != nil {
					logs.Logger.Errorf("[cloud.Expand] recover error. error: %v", bErr)
					logs.Logger.Errorf("stacktrace from panic: \n" + string(debug.Stack()))
					createdError <- fmt.Errorf("panic %v", bErr)
				}
			}()
			batchInstanceIds := make([]string, 0)
			if cur < constants.BatchMax {
				batchInstanceIds, bErr = provider.BatchCreate(params, cur)
			} else {
				batchInstanceIds, bErr = provider.BatchCreate(params, constants.BatchMax)
			}
			if bErr != nil {
				logs.Logger.Errorf("[cloud.Expand] BatchCreate error. error: %s", bErr.Error())
				createdError <- bErr
				return
			}
			createdBatch <- batchInstanceIds
		}(cur)
	}
	errs := make([]error, 0)
	for i := 0; i < batch; i++ {
		select {
		case instanceIdSet := <-createdBatch:
			instanceIds = append(instanceIds, instanceIdSet...)
		case cErr := <-createdError:
			// handle errs?
			errs = append(errs, cErr)
			err = cErr
		}
	}
	return instanceIds, err
}

func GetInstanceByTag(c *types.ClusterInfo, tags []cloud.Tag) (instances []cloud.Instance, err error) {
	provider, err := getProvider(c.Provider, c.AccountKey, c.RegionId)
	if err != nil {
		return
	}
	return provider.GetInstancesByTags(c.RegionId, tags)
}

func generateParams(clusterInfo *types.ClusterInfo, tags []cloud.Tag) (params cloud.Params, err error) {
	params.ImageId = clusterInfo.Image
	if clusterInfo.ImageConfig.Id != "" {
		params.ImageId = clusterInfo.ImageConfig.Id
	}
	params.Network = &cloud.Network{
		VpcId:                   clusterInfo.NetworkConfig.Vpc,
		SubnetId:                clusterInfo.NetworkConfig.SubnetId,
		SecurityGroup:           clusterInfo.NetworkConfig.SecurityGroup,
		InternetChargeType:      clusterInfo.NetworkConfig.InternetChargeType,
		InternetMaxBandwidthOut: clusterInfo.NetworkConfig.InternetMaxBandwidthOut,
		InternetIpType:          clusterInfo.NetworkConfig.InternetIpType,
	}
	params.InstanceType = clusterInfo.InstanceType
	params.Password = clusterInfo.Password
	params.Provider = clusterInfo.Provider
	params.Region = clusterInfo.RegionId
	params.Zone = clusterInfo.ZoneId
	params.Disks = clusterInfo.StorageConfig.Disks
	params.Tags = tags
	params.Charge = &cloud.Charge{
		ChargeType: clusterInfo.ChargeConfig.ChargeType,
		Period:     clusterInfo.ChargeConfig.Period,
		PeriodUnit: clusterInfo.ChargeConfig.PeriodUnit,
	}
	return
}

func getBatch(num, eachMax int) int {
	if num%eachMax == 0 {
		return num / eachMax
	}
	return num/eachMax + 1
}

func getProvider(provider, ak, regionId string) (cloud.Provider, error) {
	var client cloud.Provider
	key := provider + ak + regionId
	v, exist := clientMap.Load(key)
	if exist {
		return v.(cloud.Provider), nil
	}

	ctx := context.Background()
	sk, err := GetAccountSecretByAccountKey(ctx, ak)
	if err != nil {
		return nil, fmt.Errorf("found sk failed, %s", err.Error())
	}
	if sk == "" {
		return nil, errors.New("no sk found")
	}

	switch provider {
	case cloud.AlibabaCloud:
		client, err = alibaba.New(ak, sk, regionId)
	case cloud.HuaweiCloud:
		client, err = huawei.New(ak, sk, regionId)
	default:
		return nil, errors.New("invalid provider")
	}
	if err != nil {
		return nil, err
	}
	clientMap.Store(key, client)
	return client, nil
}

func Shrink(clusterInfo *types.ClusterInfo, instanceIds []string) error {
	if len(instanceIds) == 0 {
		return nil
	}
	provider, err := getProvider(clusterInfo.Provider, clusterInfo.AccountKey, clusterInfo.RegionId)
	if err != nil {
		return err
	}
	return provider.BatchDelete(instanceIds, clusterInfo.RegionId)
}

func GetInstances(clusterInfo *types.ClusterInfo, instancesIds []string) (instances []cloud.Instance, err error) {
	provider, err := getProvider(clusterInfo.Provider, clusterInfo.AccountKey, clusterInfo.RegionId)
	if err != nil {
		return
	}
	return provider.GetInstances(instancesIds)
}

func GetCloudInstancesByClusterName(clusterInfo *types.ClusterInfo) (instances []cloud.Instance, err error) {
	provider, err := getProvider(clusterInfo.Provider, clusterInfo.AccountKey, clusterInfo.RegionId)
	if err != nil {
		return
	}
	return provider.GetInstancesByCluster(clusterInfo.RegionId, clusterInfo.Name)
}

func GetInstancesByAccount(ctx context.Context, accountKey string, pageNum, pageSize int) (instances []model.Instance, total int64, err error) {
	clusterNames, err := GetEnabledClusterNamesByAccount(ctx, accountKey)
	ret := make([]model.Instance, 0)
	total, err = model.Query(map[string]interface{}{"cluster_name": clusterNames}, pageNum, pageSize, &ret, "", true)
	if err != nil {
		return ret, 0, err
	}
	return ret, total, nil
}

type GetInstancesCond struct {
	AccountKeys []string
	Status      []string
	InstanceId  string
	Ip          string
	ClusterName string
	Provider    string
	PageNum     int
	PageSize    int
}

func GetInstancesByAccounts(ctx context.Context, cond GetInstancesCond) (clusterNames []string, instances []model.Instance, total int64, err error) {
	clusterNames, err = GetEnabledClusterNamesByCond(ctx, cond.Provider, cond.ClusterName, cond.AccountKeys, true)
	if err != nil {
		return nil, nil, 0, err
	}
	instances, total, err = model.GetInstanceByCond(ctx, model.InstanceSearchCond{
		Ip:           cond.Ip,
		InstanceId:   cond.InstanceId,
		ClusterNames: clusterNames,
		Status:       cond.Status,
		PageNumber:   cond.PageNum,
		PageSize:     cond.PageSize,
	})
	return clusterNames, instances, total, err
}
