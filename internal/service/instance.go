package service

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/galaxy-future/BridgX/internal/constants"

	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/types"
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

var zoneInsTypeCache = map[string]map[string][]InstanceTypeByZone{} // key: provider key: zoneID

func GetInstanceCount(ctx context.Context, accountKeys []string, clusterName string) (int64, error) {
	clusterNames, err := GetEnabledClusterNamesByCond(ctx, "", clusterName, accountKeys, true)
	if err != nil {
		return 0, err
	}
	ret, err := model.CountActiveInstancesByClusterName(ctx, clusterNames)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func GetInstanceCountByCluster(ctx context.Context, clusters []model.Cluster) map[string]int64 {
	retMap := make(map[string]int64, 0)
	for _, cluster := range clusters {
		ret, err := model.CountActiveInstancesByClusterName(ctx, []string{cluster.ClusterName})
		if err != nil {
			ret = 0
		}
		retMap[cluster.ClusterName] = ret
	}
	return retMap
}

func GetInstanceTypeByName(ctx context.Context, instanceTypeName string) (*model.InstanceType, error) {
	return model.GetInstanceTypeByName(ctx, instanceTypeName)
}

func GetInstancesByTaskId(ctx context.Context, taskId string, taskAction string) ([]model.Instance, error) {
	ret := make([]model.Instance, 0)
	m := make(map[string]interface{}, 0)
	if taskAction == constants.TaskActionExpand {
		m["task_id"] = taskId
	} else if taskAction == constants.TaskActionShrink {
		m["shrink_task_id"] = taskId
	} else {
		return nil, errors.New("not support task action")
	}
	err := model.QueryAll(m, &ret, "")
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func GetInstanceUsageTotal(ctx context.Context, clusterName string, specifyDay time.Time, orgId int64) (int64, error) {
	clusterNames := make([]string, 0)
	if clusterName == "" {
		accounts, err := GetAksByOrgId(orgId)
		if err != nil {
			return 0, err
		}
		if len(accounts) == 0 {
			return 0, nil
		}
		clusterNames, err = GetEnabledClusterNamesByAccounts(ctx, accounts)
		if err != nil {
			return 0, err
		}
		if len(clusterNames) == 0 {
			return 0, nil
		}
	} else {
		clusterNames = append(clusterNames, clusterName)
	}
	notDeletedInstances, err := model.GetActiveInstancesByClusters(ctx, clusterNames)
	if err != nil {
		return 0, err
	}
	var totalUsage int64
	var specDayStart, specDayEnd = specifyDay, specifyDay.Add(24 * time.Hour)
	var timeEnd = specDayEnd
	if timeEnd.After(time.Now()) {
		timeEnd = time.Now()
	}
	deletedInstances, err := model.GetDeletedInstancesByTime(ctx, clusterNames, specDayEnd, specDayStart)
	if err != nil {
		return 0, nil
	}
	for _, instance := range notDeletedInstances {
		if instance.CreateAt.After(specDayEnd) {
			continue
		}
		if instance.CreateAt.Before(specDayStart) {
			totalUsage += int64(timeEnd.Sub(specDayStart).Seconds())
		} else {
			totalUsage += int64(timeEnd.Sub(*instance.CreateAt).Seconds())
		}
	}
	for _, instance := range deletedInstances {
		start := instance.CreateAt
		end := *instance.DeleteAt
		if start.Before(specDayStart) {
			start = &specDayStart
		}
		if end.After(specDayEnd) {
			end = specDayEnd
		}
		totalUsage += int64(end.Sub(*start).Seconds())
	}
	return totalUsage, nil
}

func GetInstanceUsageStatistics(ctx context.Context, clusterName string, specifyDay time.Time, orgId int64, pageNum, pageSize int) ([]model.Instance, int64, error) {
	clusterNames := make([]string, 0)
	if clusterName == "" {
		accounts, err := GetAksByOrgId(orgId)
		if err != nil {
			return nil, 0, err
		}
		if len(accounts) == 0 {
			return nil, 0, nil
		}
		clusterNames, err = GetEnabledClusterNamesByAccounts(ctx, accounts)
		if err != nil {
			return nil, 0, err
		}
		if len(clusterNames) == 0 {
			return nil, 0, nil
		}
	} else {
		clusterNames = append(clusterNames, clusterName)
	}
	var specDayStart, specDayEnd = specifyDay, specifyDay.Add(24 * time.Hour)

	instances, total, err := model.GetUsageInstancesBySpecifyDay(ctx, clusterNames, specDayEnd, specDayStart, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return instances, total, nil
}

type InstancesSearchCond struct {
	TaskId     int64
	TaskAction string
	Status     string
	PageNumber int
	PageSize   int
}

func GetInstancesByCond(ctx context.Context, cond InstancesSearchCond) (ret []model.Instance, total int64, err error) {
	queryMap := map[string]interface{}{}
	if cond.TaskAction == constants.TaskActionExpand {
		queryMap["task_id"] = cond.TaskId
	}
	if cond.TaskAction == constants.TaskActionShrink {
		queryMap["shrink_task_id"] = cond.TaskId
	}
	if cond.Status != "" {
		queryMap["status"] = cond.Status
	}
	total, err = model.Query(queryMap, cond.PageNumber, cond.PageSize, &ret, "id", true)
	if err != nil {
		return ret, 0, err
	}
	return ret, total, nil
}

func GetInstance(ctx context.Context, instanceId string) (*model.Instance, error) {
	ret, err := model.GetInstanceByInstanceId(instanceId)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func SyncInstanceTypes(ctx context.Context, provider string) error {
	accounts, err := GetDefaultAccount(provider)
	if err != nil {
		return err
	}
	regions, err := GetRegions(ctx, GetRegionsRequest{
		Provider: provider,
		Account:  accounts,
	})
	if err != nil {
		return err
	}
	ak := getFirstAk(accounts, provider)
	instanceInfoMap := make(map[string]*cloud.InstanceInfo)
	insInfoReq := make([]string, 0, 10)
	instanceTypes := make([]model.InstanceType, 0, 1000)
	for _, region := range regions {
		p, err := getProvider(provider, ak, region.RegionId)
		if err != nil {
			logs.Logger.Errorf("region[%s] getProvider failed,err: %v", region.RegionId, err)
			continue
		}
		res, err := p.DescribeAvailableResource(cloud.DescribeAvailableResourceRequest{
			RegionId: region.RegionId,
		})
		if err != nil {
			logs.Logger.Errorf("region[%s] DescribeAvailableResource failed,err: %v", region.RegionId, err)
		}
		for zone, ins := range res.InstanceTypes {
			for i, in := range ins {
				if _, ok := instanceInfoMap[in.Value]; !ok {
					instanceInfoMap[in.Value] = new(cloud.InstanceInfo)
					insInfoReq = append(insInfoReq, in.Value)
				}
				if len(insInfoReq) == 10 || len(ins)-1 == i && len(insInfoReq) > 0 {
					res, err := p.DescribeInstanceTypes(cloud.DescribeInstanceTypesRequest{TypeName: insInfoReq})
					if err != nil {
						logs.Logger.Errorf("region[%s] DescribeInstanceTypes failed,err: %v req: %v", region.RegionId, err, insInfoReq)
					}
					for _, info := range res.Infos {
						instanceInfoMap[info.InsTypeName].Family = info.Family
						instanceInfoMap[info.InsTypeName].Memory = info.Memory
						instanceInfoMap[info.InsTypeName].Core = info.Core
					}
					insInfoReq = insInfoReq[0:0]
				}
				instanceTypes = append(instanceTypes, model.InstanceType{
					Provider: provider,
					RegionId: region.RegionId,
					ZoneId:   zone,
					TypeName: in.Value,
				})
			}
		}

	}
	inss := make([]model.InstanceType, 0, 100)
	for i, insType := range instanceTypes {
		insInfo := instanceInfoMap[insType.TypeName]
		insType.Family = insInfo.Family
		insType.Memory = insInfo.Memory
		insType.Core = insInfo.Core
		now := time.Now()
		insType.CreateAt = &now
		insType.UpdateAt = &now
		if len(inss) == 100 || len(instanceTypes)-1 == i {
			err := BatchCreateInstanceType(ctx, inss)
			if err != nil {
				logs.Logger.Errorf("inss[%v] BatchCreateInstanceType failed,err: %v", inss, err)
			}
			inss = inss[0:0]
		}
		inss = append(inss, insType)
	}
	return exchangeStatus(ctx)
}

type ListInstanceTypeRequest struct {
	Provider string
	RegionId string
	ZoneId   string
	Account  *types.OrgKeys
}

type ListInstanceTypeResponse struct {
	InstanceTypes []InstanceTypeByZone `json:"instance_types"`
}

type InstanceTypeByZone struct {
	InstanceTypeFamily string `json:"instance_type_family"`
	InstanceType       string `json:"instance_type"`
	Core               int    `json:"core"`
	Memory             int    `json:"memory"`
}

func ListInstanceType(ctx context.Context, req ListInstanceTypeRequest) (ListInstanceTypeResponse, error) {
	if len(zoneInsTypeCache) == 0 {
		RefreshCache()
	}
	zoneMap, ok := zoneInsTypeCache[req.Provider]
	if !ok {
		return ListInstanceTypeResponse{}, nil
	}
	res, ok := zoneMap[req.ZoneId]
	if !ok {
		return ListInstanceTypeResponse{}, nil
	}
	return ListInstanceTypeResponse{InstanceTypes: res}, nil
}

func BatchCreateInstanceType(ctx context.Context, inss []model.InstanceType) error {
	return model.BatchCreate(inss)
}

func exchangeStatus(ctx context.Context) error {
	tx := clients.WriteDBCli.Begin()
	err := model.UpdateInstanceTypeIStatus(ctx, tx, model.InstanceTypeStatusExpired)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = model.UpdateInstanceTypeIStatus(ctx, tx, model.InstanceTypeStatusActivated)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = model.DropExpiredInstanceType(ctx, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	err = RefreshCache()
	if err != nil {
		logs.Logger.Infof("RefreshCache error:%v", err)
	}
	return nil
}

func RefreshCache() error {
	ctx := context.Background()
	ins, err := model.ScanInstanceType(ctx)
	if err != nil {
		logs.Logger.Error("RefreshCache Error err:%v", err)
		return err
	}
	if len(ins) == 0 {
		// TODO: SELECT `provider`,`access_key` FROM ACCOUNT GROUP BY `provider`.
		err = SyncInstanceTypes(ctx, cloud.AlibabaCloud)
		if err != nil {
			logs.Logger.Error("SyncInstanceTypes Error err:%v", err)
			return err
		}
		ins, err = model.ScanInstanceType(ctx)
		if err != nil {
			logs.Logger.Error("ScanInstanceType Error err:%v", err)
			return err
		}
	}
	for _, in := range ins {
		provider := in.Provider
		providerMap, ok := zoneInsTypeCache[provider]
		if !ok {
			zoneInsTypeCache[provider] = make(map[string][]InstanceTypeByZone)
			providerMap = zoneInsTypeCache[provider]
		}

		zoneId := in.ZoneId
		_, ok = providerMap[zoneId]
		if !ok {
			providerMap[zoneId] = make([]InstanceTypeByZone, 0, 400)
		}
		providerMap[zoneId] = append(providerMap[zoneId], InstanceTypeByZone{
			InstanceTypeFamily: in.Family,
			InstanceType:       in.TypeName,
			Core:               in.Core,
			Memory:             in.Memory,
		})
	}
	for provider, zoneMap := range zoneInsTypeCache {
		for zone, typeList := range zoneMap {
			sort.Slice(typeList, func(i, j int) bool {
				typeI := typeList[i]
				typeJ := typeList[j]
				if typeI.Core < typeJ.Core {
					return true
				}
				if typeI.Core == typeJ.Core && typeI.Memory <= typeJ.Memory {
					return true
				}
				return false
			})
			zoneMap[zone] = typeList
		}
		zoneInsTypeCache[provider] = zoneMap
	}
	return nil
}
