package service

import (
	"context"
	"errors"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/config"
	"github.com/galaxy-future/BridgX/internal/bcc"
	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/types"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/utils"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateClusterWithTagsAndInstances(ctx context.Context, cluster *model.Cluster, tags []*model.ClusterTag, instances []model.Instance, username string, uid int64) error {
	now := time.Now()
	cluster.CreateAt = &now
	cluster.UpdateAt = &now
	cluster.CreateBy = username
	cluster.UpdateBy = username
	cluster.Status = constants.ClusterStatusEnable
	if len(tags) > 0 {
		for _, tag := range tags {
			tag.CreateAt = &now
			tag.UpdateAt = &now
		}
	}
	err := model.CreateClusterWithTagsAndInstances(ctx, cluster, tags, instances)
	if err != nil {
		return err
	}
	err = RecordOperationLog(ctx, OperationLog{
		Operation: OperationCreate,
		Operator:  uid,
		Old:       nil,
		New:       cluster,
	})
	if err != nil {
		logs.Logger.Errorf("RecordOperationLog failed.Err:[%s]", err.Error())
	}
	return nil
}

func EditCluster(cluster *model.Cluster, username string) error {
	clusterInDB, err := model.GetByClusterName(cluster.ClusterName)
	if err != nil {
		return err
	}
	if clusterInDB == nil {
		return errors.New("editing cluster not exist")
	}
	now := time.Now()
	cluster.Id = clusterInDB.Id
	cluster.Status = clusterInDB.Status
	cluster.CreateAt = clusterInDB.CreateAt
	cluster.CreateBy = clusterInDB.CreateBy
	cluster.UpdateAt = &now
	cluster.UpdateBy = username
	return model.Save(cluster)
}

func DeleteClusters(ctx context.Context, ids []int64, orgId int64) error {
	clusters := make([]model.Cluster, 0)
	if len(ids) == 0 {
		return nil
	}
	err := model.Gets(ids, &clusters)
	if err != nil {
		return err
	}
	if len(clusters) == 0 {
		return nil
	}
	return clients.WriteDBCli.Transaction(func(tx *gorm.DB) error {
		for _, cluster := range clusters {
			err = tx.Delete(model.ClusterTag{}, model.ClusterTag{ClusterName: cluster.ClusterName}).Error
			if err != nil {
				return err
			}
			c := cluster
			c.DeleteUniqKey = c.Id
			err = tx.Save(&c).Error
			if err != nil {
				return err
			}
		}
		err = tx.Delete(clusters).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func CreateCluster4Test(clusterName string) error {
	cluster := &model.Cluster{ClusterName: clusterName}
	cluster.Status = constants.ClusterStatusEnable
	return model.Create(cluster)
}

func GetClusterById(ctx context.Context, Id int64) (*model.Cluster, error) {
	cluster := &model.Cluster{}
	err := model.Get(Id, cluster)
	if err != nil {
		return nil, err
	}
	return cluster, err
}

func GetClusterByName(ctx context.Context, name string) (*model.Cluster, error) {
	cluster, err := model.GetByClusterName(name)
	if err != nil {
		return nil, err
	}
	return cluster, err
}

func GetClustersByNames(ctx context.Context, names []string) ([]model.Cluster, error) {
	cluster, err := model.GetByClusterNames(names)
	if err != nil {
		return nil, err
	}
	return cluster, err
}

func GetClusterTagsByClusterName(ctx context.Context, name string) ([]model.ClusterTag, error) {
	clusterTags := make([]model.ClusterTag, 0)
	err := model.QueryAll(map[string]interface{}{"cluster_name": name}, &clusterTags, "")
	if err != nil {
		return nil, err
	}
	return clusterTags, err
}

func GetClusterCount(ctx context.Context, accountKeys []string) (count int64, err error) {
	err = clients.ReadDBCli.WithContext(ctx).Model(&model.Cluster{}).
		Where("account_key in (?) and status = ?", accountKeys, constants.ClusterStatusEnable).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func ListClusters(ctx context.Context, cond model.ClusterSearchCond) ([]model.Cluster, int, error) {
	return model.ListClustersByCond(ctx, cond)
}

func GetEnabledClusterNamesByAccount(ctx context.Context, accountKey string) ([]string, error) {
	res := make([]string, 0)
	err := clients.ReadDBCli.WithContext(ctx).Model(&model.Cluster{}).Select("cluster_name").Where("account_key = ? AND status IN (?)", accountKey, constants.ClusterStatusEnable).Find(&res).Error
	if err != nil {
		return res, err
	}
	return res, nil

}

func GetEnabledClusterNamesByCond(ctx context.Context, provider, clusterName string, aks []string, strict bool) ([]string, error) {
	res := make([]string, 0)
	query := clients.ReadDBCli.WithContext(ctx).
		Model(&model.Cluster{}).
		Select("cluster_name").
		Where("status = ?", constants.ClusterStatusEnable)
	if len(aks) > 0 && provider != cloud.PrivateCloud {
		query = query.Where("account_key IN (?)", aks)
	}
	if provider != "" {
		query = query.Where("provider = ?", provider)
	}
	if clusterName != "" {
		if strict {
			query = query.Where("cluster_name = ?", clusterName)
		} else {
			query = query.Where("cluster_name like ?", "%"+clusterName+"%")
		}
	}
	err := query.Find(&res).Error
	if err != nil {
		return res, err
	}
	return res, nil
}

func GetStandardClusterNamesByAccounts(ctx context.Context, accountKeys []string) ([]string, error) {
	res := make([]string, 0)
	err := clients.ReadDBCli.WithContext(ctx).Model(&model.Cluster{}).Select("cluster_name").Where("account_key in (?) AND status = ? AND cluster_type = ?", accountKeys, constants.ClusterStatusEnable, constants.ClusterTypeStandard).Find(&res).Error
	if err != nil {
		return res, err
	}
	return res, nil

}

//ConvertToClusterInfo 将cluster，和tags转换为一个Cloud clusterInfo
func ConvertToClusterInfo(m *model.Cluster, tags []model.ClusterTag) (*types.ClusterInfo, error) {
	imageConfig := &types.ImageConfig{}
	networkConfig := &types.NetworkConfig{}
	storageConfig := &types.StorageConfig{}
	chargeConfig := &types.ChargeConfig{}
	if m.ImageConfig != "" {
		err := jsoniter.UnmarshalFromString(m.ImageConfig, imageConfig)
		if err != nil {
			return nil, err
		}
	}
	if m.NetworkConfig != "" {
		err := jsoniter.UnmarshalFromString(m.NetworkConfig, networkConfig)
		if err != nil {
			return nil, err
		}
	}
	if m.StorageConfig != "" {
		err := jsoniter.UnmarshalFromString(m.StorageConfig, storageConfig)
		if err != nil {
			return nil, err
		}
	}
	if m.ChargeConfig != "" {
		err := jsoniter.UnmarshalFromString(m.ChargeConfig, chargeConfig)
		if err != nil {
			return nil, err
		}
	}
	extendConfig := &types.ExtendConfig{}
	if m.ExtendConfig != "" {
		err := jsoniter.UnmarshalFromString(m.ExtendConfig, extendConfig)
		if err != nil {
			return nil, err
		}
	}
	var mt = make(map[string]string, 0)
	for _, clusterTag := range tags {
		mt[clusterTag.TagKey] = clusterTag.TagValue
	}

	clusterInfo := &types.ClusterInfo{
		Id:            m.Id,
		Name:          m.ClusterName,
		Desc:          m.ClusterDesc,
		RegionId:      m.RegionId,
		ZoneId:        m.ZoneId,
		ClusterType:   m.ClusterType,
		InstanceType:  m.InstanceType,
		Image:         m.Image,
		Provider:      m.Provider,
		Username:      constants.DefaultUsername,
		Password:      m.Password,
		AccountKey:    m.AccountKey,
		ImageConfig:   imageConfig,
		NetworkConfig: networkConfig,
		StorageConfig: storageConfig,
		ChargeConfig:  chargeConfig,
		ExtendConfig:  extendConfig,
		Tags:          mt,
	}
	return clusterInfo, nil
}

func ExpandCluster(c *types.ClusterInfo, num int, taskId int64) ([]string, []string, error) {
	//调用云厂商接口进行扩容
	expandInstanceIds, expandErr := ExpandInDeed(c, num, taskId)
	if len(expandInstanceIds) == 0 && expandErr != nil {
		return nil, nil, expandErr
	}

	//将扩容的Instance信息保存到DB
	err := saveExpandInstancesToDB(c, expandInstanceIds, taskId)
	if err != nil {
		logs.Logger.Errorf("[ExpandCluster] saveExpandInstancesToDB error. cluster name: %s, error: %v", c.Name, err)
		return nil, expandInstanceIds, err
	}

	//查询扩容的Instance的IP并保存
	expandIPs, availableIds, err := queryAndSaveExpandIPs(c, taskId, len(expandInstanceIds))
	if err != nil {
		logs.Logger.Errorf("[ExpandCluster] queryAndSaveExpandIPs error. cluster name: %s, error: %v", c.Name, err)
		return availableIds, expandInstanceIds, err
	}
	//发布扩容信息到配置中心
	_ = publishExpandConfig(c.Name, expandInstanceIds, expandIPs)
	return availableIds, expandInstanceIds, expandErr
}

func ShrinkClusterBySpecificIps(c *types.ClusterInfo, deletingIPs string, count int, taskId int64) (err error) {
	toBeDeletedIds, notExistIds := getMappingInstanceIdList(c.Name, deletingIPs)
	if len(toBeDeletedIds) == 0 {
		logs.Logger.Warnf("%v has no deletingIPs %v", c.Name, deletingIPs)
		return nil
	}
	if len(toBeDeletedIds)+len(notExistIds) != count {
		logs.Logger.Warnf("%v toBeDeleted:%v + alreadyDeleted:%v not match expect_shrink_count:%v", c.Name, toBeDeletedIds, notExistIds, count)
		return errors.New("need delete instance count NOT MATCH expect delete count")
	}
	logs.Logger.Infof("cluster:%v, DELETING ip list:%v, instances list:%v", c.Name, deletingIPs, toBeDeletedIds)
	err = Shrink(c, toBeDeletedIds)
	if err != nil {
		logs.Logger.Errorf("[ShrinkCluster] Shrink instance error. cluster name: %s, error: %s", c.Name, err.Error())
		return
	}
	now := time.Now()
	err = model.BatchUpdateByInstanceIds(toBeDeletedIds, model.Instance{
		ShrinkTaskId: taskId,
		Status:       constants.Deleted,
		DeleteAt:     &now,
	})
	if err != nil {
		logs.Logger.Errorf("[ShrinkClusterBySpecificIps] update db error. cluster name: %s, error: %s", c.Name, err.Error())
		return
	}
	_ = publishShrinkConfig(c.Name)
	return err
}

func ShrinkCluster(c *types.ClusterInfo, num int, taskId int64) (err error) {
	logs.Logger.Infof("Shrink %v, with count:%v", c.Name, num)
	instances, err := model.GetActiveInstancesWithCount(c.Name, num)
	if err != nil {
		logs.Logger.Errorf("[ShrinkCluster] Get instanceIdStr error. cluster name: %s, error: %s", c.Name, err.Error())
		return err
	}
	toBeDeletedInstanceIds := make([]string, 0)
	for _, instance := range instances {
		toBeDeletedInstanceIds = append(toBeDeletedInstanceIds, instance.InstanceId)
	}
	err = Shrink(c, toBeDeletedInstanceIds)
	if err != nil {
		logs.Logger.Errorf("[ShrinkCluster] Shrink instance error. cluster name: %s, error: %s", c.Name, err.Error())
		return
	}
	now := time.Now()
	err = model.BatchUpdateByInstanceIds(toBeDeletedInstanceIds, model.Instance{
		ShrinkTaskId: taskId,
		Status:       constants.Deleted,
		DeleteAt:     &now,
	})
	if err != nil {
		logs.Logger.Errorf("[ShrinkCluster] Shrink instance update db error. cluster name: %s, error: %s", c.Name, err.Error())
		return
	}
	_ = publishShrinkConfig(c.Name)
	return err
}

func CreateShrinkAllTask(ctx context.Context, clusterName, taskName string, uid int64) (int64, error) {
	count, err := model.CountActiveInstancesByClusterName(ctx, []string{clusterName})
	if err != nil {
		return 0, err
	}
	return CreateShrinkTask(ctx, clusterName, int(count), "", taskName, uid)
}

//CleanClusterUnusedInstances 清除由于系统异常导致的云厂商中残留的机器
func CleanClusterUnusedInstances(clusterInfo *types.ClusterInfo) (int, error) {
	instancesInBridgx, err := model.GetActiveInstancesByClusterName(clusterInfo.Name)
	if err != nil {
		return 0, err
	}
	instanceInCloud, err := GetCloudInstancesByClusterName(clusterInfo)
	if err != nil {
		return 0, err
	}
	instanceIds := calcUnusedInstancesId(instanceInCloud, instancesInBridgx)
	if len(instanceIds) > 0 {
		err := Shrink(clusterInfo, instanceIds)
		if err != nil {
			return 0, err
		}
	}
	return len(instanceIds), nil
}

func calcUnusedInstancesId(cloudInstances []cloud.Instance, bridgeXInstances []model.Instance) []string {
	var unusedInstanceIds []string
	bridgxInstanceExists := make(map[string]struct{})
	for _, bridgxInstance := range bridgeXInstances {
		bridgxInstanceExists[bridgxInstance.InstanceId] = struct{}{}
	}

	for _, cloudInstance := range cloudInstances {
		if _, exists := bridgxInstanceExists[cloudInstance.Id]; !exists {
			unusedInstanceIds = append(unusedInstanceIds, cloudInstance.Id)
		}
	}
	return unusedInstanceIds
}

func getMappingInstanceIdList(clusterName, deletingIPs string) (toBeDeletedIds, notExistIds []string) {
	activeInstances, err := model.GetActiveInstancesByClusterName(clusterName)
	if err != nil || len(activeInstances) == 0 {
		return nil, nil
	}
	m := make(map[string]string, 0)
	for _, instance := range activeInstances {
		m[instance.IpInner] = instance.InstanceId
	}
	for _, ip := range strings.Split(deletingIPs, ",") {
		if insId, ok := m[ip]; ok {
			toBeDeletedIds = append(toBeDeletedIds, insId)
		} else {
			notExistIds = append(notExistIds, insId)
		}
	}
	logs.Logger.Infof("%v real delete working IPs:%v", clusterName, toBeDeletedIds)
	return
}

func getDelayFactor(n int) int {
	if n < 2 {
		return 5
	} else if n < 5 {
		return 8
	} else if n < 8 {
		return 13
	} else {
		return 20
	}
}

func queryAndSaveExpandIPs(c *types.ClusterInfo, taskId int64, idNum int) ([]string, []string, error) {
	var err error
	needPublicIp := false
	if c.NetworkConfig != nil && c.NetworkConfig.InternetMaxBandwidthOut > 0 {
		needPublicIp = true
	}
	var instances []cloud.Instance
	var insNum int
	tags := []cloud.Tag{{
		Key:   cloud.TaskId,
		Value: strconv.FormatInt(taskId, 10),
	}}
	// TODO scheduler
	for k := 0; k < constants.Interval; k++ {
		instances, err = GetInstanceByTag(c, tags)
		insNum = len(instances)
		logs.Logger.Infof("[queryAndSaveExpandIPs] insNum: %d, idNum: %d, err: %v", insNum, idNum, err)
		if err == nil && insNum == idNum && judgeInstancesIsReady(instances, needPublicIp) {
			logs.Logger.Infof("[queryAndSaveExpandIPs] is ready, %d", insNum)
			break
		}
		time.Sleep(time.Duration(getDelayFactor(k)) * time.Second)
	}
	if err != nil {
		return nil, nil, err
	}

	expandIps := make([]string, 0, insNum)
	expandIds := make([]string, 0, insNum)
	for _, instance := range instances {
		if !IsInstanceReady(instance, needPublicIp) {
			logs.Logger.Errorf("[syncDbAndConfig] InstanceId:%v is not ready", instance.Id)
			continue
		}
		update := func(attempt uint) error {
			now := time.Now()
			var expireAt *time.Time
			if instance.CostWay == cloud.InstanceChargeTypePrePaid {
				expireAt = instance.ExpireAt
			}
			return model.UpdateByInstanceId(model.Instance{
				InstanceId:  instance.Id,
				IpInner:     instance.IpInner,
				IpOuter:     instance.IpOuter,
				ClusterName: c.Name,
				Status:      constants.Running,
				RunningAt:   &now,
				ExpireAt:    expireAt,
			})
		}
		err = retry.Retry(update, strategy.Limit(3), strategy.Backoff(backoff.Fibonacci(10*time.Millisecond)))
		if err != nil {
			logs.Logger.Errorf("[syncDbAndConfig] UpdateByInstanceId Error IP:%v, instanceId:%v", instance.IpInner, instance.Id)
			continue
		}

		expandIps = append(expandIps, instance.IpInner)
		expandIds = append(expandIds, instance.Id)
	}
	return expandIps, expandIds, nil
}

func saveExpandInstancesToDB(c *types.ClusterInfo, expandInstanceIds []string, taskId int64) error {
	instances := make([]model.Instance, 0)
	now := time.Now()
	for _, instanceId := range expandInstanceIds {
		instances = append(instances, model.Instance{
			Base: model.Base{
				CreateAt: &now,
			},
			TaskId:      taskId,
			InstanceId:  instanceId,
			Status:      constants.Pending,
			ClusterName: c.Name,
			ChargeType:  c.ChargeConfig.ChargeType,
		})
	}
	return model.BatchCreateInstance(instances)
}

func publishExpandConfig(clusterName string, expandInstanceIds []string, expandIPs []string) error {
	if !config.GlobalConfig.NeedPublishConfig {
		logs.Logger.Infof("expand cluster:%v no need publish config", clusterName)
		return nil
	}
	//将扩容的实例信息发布到配置中心
	err := publishExpandInstanceConfig(clusterName, expandInstanceIds)
	if err != nil {
		logs.Logger.Errorf("[ExpandCluster] Publish instance_id_list error. cluster name: %s, error: %v", clusterName, err)
		return err
	}

	//将扩容的IP信息发布到配置中心
	err = publishExpandIPConfig(clusterName, expandIPs)
	if err != nil {
		logs.Logger.Errorf("[ExpandCluster] publishExpandIPConfig error. cluster name: %s, error: %v", clusterName, err)
		return err
	}
	return nil
}

func publishExpandIPConfig(clusterName string, expandIPs []string) error {
	existingIPs, _ := bcc.GetConfig(clusterName, constants.WorkingIPs)
	totalIps := expandIPs
	if existingIPs != "" && existingIPs != constants.HasNoneIP {
		totalIps = append(totalIps, strings.Split(existingIPs, ",")...)
	}
	err := bcc.PublishConfig(clusterName, constants.WorkingIPs, strings.Join(totalIps, ","))
	return err
}

func publishExpandInstanceConfig(clusterName string, expandInstanceIds []string) error {
	instanceIdsStr, _ := bcc.GetConfig(clusterName, constants.Instances)
	totalInstanceIds := expandInstanceIds
	if instanceIdsStr != "" && instanceIdsStr != constants.HasNoneInstance {
		totalInstanceIds = append(totalInstanceIds, strings.Split(instanceIdsStr, ",")...)
	}
	return bcc.PublishConfig(clusterName, constants.Instances, strings.Join(totalInstanceIds, ","))
}

func publishShrinkConfig(clusterName string) error {
	if !config.GlobalConfig.NeedPublishConfig {
		logs.Logger.Infof("shrink cluster:%v no need publish config", clusterName)
		return nil
	}
	instances, err := model.GetActiveInstancesByClusterName(clusterName)
	if err != nil || len(instances) == 0 {
		return err
	}
	restInstanceIds := make([]string, 0)
	restInstanceIPs := make([]string, 0)

	restInstancesStr := constants.HasNoneInstance
	restIps := constants.HasNoneIP

	for _, instance := range instances {
		restInstanceIds = append(restInstanceIds, instance.InstanceId)
		restInstanceIPs = append(restInstanceIPs, instance.IpInner)
	}
	if len(restInstanceIds) > 0 {
		restInstancesStr = strings.Join(restInstanceIds, ",")
	}
	if len(restInstanceIPs) > 0 {
		restIps = strings.Join(restInstanceIPs, ",")
	}

	err = bcc.PublishConfig(clusterName, constants.Instances, restInstancesStr)
	if err != nil {
		logs.Logger.Errorf("[ExpandCluster] Publish instance_id_list error. cluster name: %s, error: %s", clusterName, err.Error())
	}
	err = bcc.PublishConfig(clusterName, constants.WorkingIPs, restIps)
	if err != nil {
		logs.Logger.Errorf("[ExpandCluster] Publish ip_list error. cluster name: %s, error: %s", clusterName, err.Error())
	}
	return err
}

func IsInstanceReady(instance cloud.Instance, needPublicIp bool) bool {
	if instance.Status != cloud.EcsRunning || instance.IpInner == "" {
		return false
	}
	if needPublicIp && instance.IpOuter == "" {
		return false
	}
	return true
}

func judgeInstancesIsReady(instances []cloud.Instance, needPublicIp bool) bool {
	for _, instance := range instances {
		if !IsInstanceReady(instance, needPublicIp) {
			return false
		}
	}
	return true
}

// CheckInstanceConnectable 检测机器连通性
func CheckInstanceConnectable(instances []model.CustomClusterInstance) response.CheckInstanceConnectableResponse {
	resMachines := make([]*model.ConnectableResult, 0)
	ch := make(chan *model.ConnectableResult, len(instances))
	var wg sync.WaitGroup
	for _, req := range instances {
		wg.Add(1)
		go func(req model.CustomClusterInstance) {
			defer func() {
				if r := recover(); r != nil {
					logs.Logger.Errorf("CheckInstanceConnectable err:%v ", r)
					logs.Logger.Errorw("CheckInstanceConnectable panic", zap.String("stack", string(debug.Stack())))
				}
				wg.Done()
			}()
			isPass := utils.SshCheck(req.InstanceIp, req.LoginName, req.LoginPassword)
			res := &model.ConnectableResult{
				InstanceIp: req.InstanceIp,
				IsPass:     isPass,
			}
			ch <- res
		}(req)
	}
	wg.Wait()
	close(ch)
	isAllPass := true
	for i := 0; i < len(instances); i++ {
		machine, ok := <-ch
		if !ok {
			continue
		}
		if !machine.IsPass {
			isAllPass = false
		}
		resMachines = append(resMachines, machine)
	}
	res := response.CheckInstanceConnectableResponse{
		IsAllPass:    isAllPass,
		InstanceList: resMachines,
	}
	return res
}
