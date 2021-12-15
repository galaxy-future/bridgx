package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/galaxy-future/BridgX/config"
	"github.com/galaxy-future/BridgX/internal/bcc"
	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/types"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	jsoniter "github.com/json-iterator/go"
)

func CreateCluster(ctx context.Context, cluster *model.Cluster, username string, userId int64) error {
	cluster.Status = constants.ClusterStatusEnable
	now := time.Now()
	cluster.CreateAt = &now
	cluster.UpdateAt = &now
	cluster.CreateBy = username
	cluster.UpdateBy = username
	err := model.Create(cluster)
	if err != nil {
		return err
	}
	err = RecordOperationLog(ctx, OperationLog{
		Operation: "create cluster",
		Operator:  userId,
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
	for _, cluster := range clusters {
		c := cluster
		c.DeleteUniqKey = c.Id
		err = model.Save(&c)
		if err != nil {
			return err
		}
	}
	return model.Delete(clusters)
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

func ListClusters(ctx context.Context, accountKeys []string, clusterName, provider string, pageNum, pageSize int) ([]model.Cluster, int, error) {
	return model.ListClustersByCond(ctx, accountKeys, clusterName, provider, pageNum, pageSize)
}

func GetEnabledClusterNamesByAccount(ctx context.Context, accountKey string) ([]string, error) {
	res := make([]string, 0)
	err := clients.ReadDBCli.WithContext(ctx).Model(&model.Cluster{}).Select("cluster_name").Where("account_key = ? AND status IN (?)", accountKey, constants.ClusterStatusEnable).Find(&res).Error
	if err != nil {
		return res, err
	}
	return res, nil

}

func GetEnabledClusterNamesByCond(ctx context.Context, ak, clusterName string, aks []string, strict bool) ([]string, error) {
	res := make([]string, 0)
	query := clients.ReadDBCli.WithContext(ctx).
		Model(&model.Cluster{}).
		Select("cluster_name").
		Where("status = ?", constants.ClusterStatusEnable)
	if ak == "" && len(aks) > 0 {
		query = query.Where("account_key IN (?)", aks)
	}
	if ak != "" {
		query = query.Where("account_key = ?", ak)
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

func GetEnabledClusterNamesByAccounts(ctx context.Context, accountKeys []string) ([]string, error) {
	res := make([]string, 0)
	err := clients.ReadDBCli.WithContext(ctx).Model(&model.Cluster{}).Select("cluster_name").Where("account_key in (?) AND status = ?", accountKeys, constants.ClusterStatusEnable).Find(&res).Error
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

	err := jsoniter.UnmarshalFromString(m.NetworkConfig, networkConfig)
	if err != nil {
		return nil, err
	}
	err = jsoniter.UnmarshalFromString(m.StorageConfig, storageConfig)
	if err != nil {
		return nil, err
	}
	err = jsoniter.UnmarshalFromString(m.ChargeConfig, chargeConfig)
	if err != nil {
		return nil, err
	}
	var mt = make(map[string]string, 0)
	for _, clusterTag := range tags {
		mt[clusterTag.TagKey] = clusterTag.TagValue
	}
	instanceType := GetInstanceTypeByName(m.InstanceType)
	clusterInfo := &types.ClusterInfo{
		Id:             m.Id,
		Name:           m.ClusterName,
		Desc:           m.ClusterDesc,
		RegionId:       m.RegionId,
		ZoneId:         m.ZoneId,
		InstanceType:   m.InstanceType,
		Image:          m.Image,
		Provider:       m.Provider,
		Username:       constants.DefaultUsername,
		Password:       m.Password,
		AccountKey:     m.AccountKey,
		ImageConfig:    imageConfig,
		NetworkConfig:  networkConfig,
		StorageConfig:  storageConfig,
		ChargeConfig:   chargeConfig,
		Tags:           mt,
		InstanceCore:   instanceType.Core,
		InstanceMemory: instanceType.Memory,
	}
	return clusterInfo, nil
}

func ExpandCluster(c *types.ClusterInfo, num int, taskId int64) (instanceIds []cloud.Instance, err error) {
	//调用云厂商接口进行扩容
	expandInstanceIds, err := ExpandAndRepair(c, num, taskId)
	if len(expandInstanceIds) == 0 && err != nil {
		return nil, err
	}

	//将扩容的Instance信息保存到DB
	err = saveExpandInstancesToDB(c, expandInstanceIds, taskId)
	if err != nil {
		logs.Logger.Errorf("[ExpandCluster] Expand error. cluster name: %s, error: %v", c.Name, err)
		return nil, err
	}

	//查询扩容的Instance的IP并保存
	expandIPs, expandInstances, err := queryAndSaveExpandIPs(c, err, expandInstanceIds)
	if err != nil {
		logs.Logger.Errorf("[ExpandCluster] queryAndSaveExpandIPs error. cluster name: %s, error: %v", c.Name, err)
		return expandInstances, err
	}
	//发布扩容信息到配置中心
	_ = publishExpandConfig(c.Name, expandInstanceIds, expandIPs)
	return expandInstances, nil
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

func queryAndSaveExpandIPs(c *types.ClusterInfo, err error, expandInstanceIds []string) ([]string, []cloud.Instance, error) {
	expandIps := make([]string, 0)
	expandInstances := make([]cloud.Instance, 0)
	// TODO scheduler
	for k := 0; k < constants.Interval; k++ {
		expandInstances, err = GetInstances(c, expandInstanceIds)
		logs.Logger.Infof("[queryAndSaveExpandIPs] expandInstances: %v, err: %v", expandInstances, err)
		if err == nil && len(expandInstances) == len(expandInstanceIds) && judgeInstancesIsReady(expandInstances, c.NetworkConfig) {
			break
		}
		time.Sleep(constants.Delay * time.Second)
	}
	if err != nil {
		logs.Logger.Errorf("[ExpandCluster] GetInstances error. cluster name: %s, error: %s", c.Name, err.Error())
	}
	for _, instance := range expandInstances {
		if instance.IpInner != "" {
			expandIps = append(expandIps, instance.IpInner)
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
		} else {
			logs.Logger.Errorf("[syncDbAndConfig] InstanceId:%v, GOT NO IP", instance.Id)
		}
		if err != nil {
			logs.Logger.Errorf("[syncDbAndConfig] UpdateByInstanceId Error IP:%v, instanceId:%v", instance.IpInner, instance.Id)
		}
	}
	return expandIps, expandInstances, err
}

func saveExpandInstancesToDB(c *types.ClusterInfo, expandInstanceIds []string, taskId int64) error {
	instances := make([]model.Instance, 0)
	now := time.Now()
	for _, instanceId := range expandInstanceIds {
		instances = append(instances, model.Instance{
			TaskId:      taskId,
			InstanceId:  instanceId,
			Status:      constants.Pending,
			ClusterName: c.Name,
			CreateAt:    &now,
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

func judgeInstancesIsReady(instances []cloud.Instance, chargeConfig *types.NetworkConfig) bool {
	var internetChargeType string
	if chargeConfig != nil {
		internetChargeType = chargeConfig.InternetChargeType
	}
	needCheckIpOuter := internetChargeType == cloud.InternetChargeTypePayByTraffic || internetChargeType == cloud.InternetChargeTypePayByBandwidth
	for _, instance := range instances {
		if instance.Status == cloud.EcsBuilding || instance.IpInner == "" {
			return false
		}
		if needCheckIpOuter && instance.IpOuter == "" {
			return false
		}
	}
	return true
}
