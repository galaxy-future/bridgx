package helper

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/galaxy-future/BridgX/internal/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

func ConvertToInstanceThumbList(ctx context.Context, instances []model.Instance, clusters []model.Cluster) []response.InstanceThumb {
	if len(instances) == 0 {
		return nil
	}
	clusterMap := genClusterMap(clusters)
	ret := make([]response.InstanceThumb, 0)
	for _, instance := range instances {
		startupTime := 0
		if instance.RunningAt != nil {
			startupTime = int(instance.RunningAt.Sub(*instance.CreateAt).Seconds())
		}
		instanceTypeDesc := getInstanceTypeDesc(instance.ClusterName, clusterMap)
		r := response.InstanceThumb{
			InstanceId:         instance.InstanceId,
			IpInner:            instance.IpInner,
			IpOuter:            instance.IpOuter,
			Provider:           getProvider(instance.ClusterName, clusterMap),
			ClusterName:        instance.ClusterName,
			InstanceType:       instanceTypeDesc,
			LoginName:          getLoginName(instance.ClusterName, clusterMap),
			LoginPassword:      getLoginPassword(instance.ClusterName, clusterMap),
			CreateAt:           instance.CreateAt.String(),
			Status:             getStringStatus(instance.Status),
			StartupTime:        startupTime,
			ChargeType:         instance.ChargeType,
			ComputingPowerType: getComputingPowerType(instanceTypeDesc),
		}
		ret = append(ret, r)
	}
	return ret
}

func getComputingPowerType(instanceTypeDesc string) string {
	if strings.Contains(instanceTypeDesc, constants.IsGpu) {
		return constants.GPU
	} else {
		return constants.CPU
	}
}

func getProvider(clusterName string, m map[string]model.Cluster) string {
	cluster, ok := m[clusterName]
	if ok {
		return cluster.Provider
	}
	return ""
}

func getLoginName(clusterName string, m map[string]model.Cluster) string {
	return "root"
}

func getLoginPassword(clusterName string, m map[string]model.Cluster) string {
	cluster, ok := m[clusterName]
	if ok {
		return cluster.Password
	}
	return ""
}

func getInstanceType(clusterName string, m map[string]model.Cluster) string {
	cluster, ok := m[clusterName]
	if ok {
		return cluster.InstanceType
	}
	return ""
}
func getInstanceTypeDesc(clusterName string, m map[string]model.Cluster) string {
	cluster, ok := m[clusterName]
	if ok {
		instanceType := service.GetInstanceTypeByName(cluster.InstanceType)
		return instanceType.GetDesc()
	}
	return ""
}

func genClusterMap(clusters []model.Cluster) map[string]model.Cluster {
	m := make(map[string]model.Cluster)
	if len(clusters) == 0 {
		return m
	}
	for _, cluster := range clusters {
		m[cluster.ClusterName] = cluster
	}
	return m
}

func ConvertToInstanceUsageList(ctx context.Context, instances []model.Instance) []response.InstanceUsage {
	if len(instances) == 0 {
		return nil
	}
	cluster, _ := service.GetClusterByName(ctx, instances[0].ClusterName)
	instanceType := ""
	if cluster != nil {
		instanceType = cluster.InstanceType
	}
	ret := make([]response.InstanceUsage, 0)
	for _, instance := range instances {
		shutdownAt := "-"
		startupTime := time.Now().Sub(*instance.CreateAt).Seconds()
		if instance.Status == constants.Deleted {
			shutdownAt = instance.DeleteAt.String()
			startupTime = instance.DeleteAt.Sub(*instance.CreateAt).Seconds()
		}
		r := response.InstanceUsage{
			Id:           cast.ToString(instance.Id),
			ClusterName:  instance.ClusterName,
			InstanceId:   instance.InstanceId,
			StartupAt:    instance.CreateAt.String(),
			ShutdownAt:   shutdownAt,
			StartupTime:  int(startupTime),
			InstanceType: instanceType,
		}
		ret = append(ret, r)
	}
	return ret
}

func ConvertToInstanceDetail(ctx context.Context, instance *model.Instance) (*response.InstanceDetail, error) {
	cluster, err := service.GetClusterByName(ctx, instance.ClusterName)
	if err != nil || cluster == nil {
		return nil, err
	}
	ret := response.InstanceDetail{
		InstanceId:    instance.InstanceId,
		Provider:      cluster.Provider,
		RegionId:      cluster.RegionId,
		ImageId:       cluster.Image,
		InstanceType:  cluster.InstanceType,
		IpInner:       instance.IpInner,
		IpOuter:       instance.IpOuter,
		CreateAt:      instance.CreateAt.String(),
		StorageConfig: parseStorageConfig(cluster.StorageConfig),
		NetworkConfig: parseNetworkConfig(cluster.NetworkConfig),
	}
	return &ret, nil
}

func FilterByComputingPowerType(ctx *gin.Context, zones service.ListInstanceTypeResponse) []service.InstanceTypeByZone {
	computingPowerType := ctx.Query("computing_power_type")
	instanceTypes := zones.InstanceTypes
	if computingPowerType == "" {
		return instanceTypes
	}

	ret := make([]service.InstanceTypeByZone, 0)
	if computingPowerType == constants.GPU {
		for i, instanceType := range instanceTypes {
			if strings.Contains(instanceType.InstanceTypeFamily, constants.IsGpu) {
				ret = append(ret, instanceTypes[i])
			}
		}
		return ret
	}

	if computingPowerType == constants.CPU {
		for i, instanceType := range instanceTypes {
			if !strings.Contains(instanceType.InstanceTypeFamily, constants.IsGpu) {
				ret = append(ret, instanceTypes[i])
			}
		}
		return ret
	}
	return instanceTypes
}

func parseNetworkConfig(config string) *response.NetworkConfig {
	nc := types.NetworkConfig{}
	err := jsoniter.UnmarshalFromString(config, &nc)
	if err != nil {
		return nil
	}
	resp := &response.NetworkConfig{
		VpcName:           nc.Vpc,
		SubnetIdName:      nc.SubnetId,
		SecurityGroupName: nc.SecurityGroup,
	}
	return resp
}

func parseStorageConfig(config string) *response.StorageConfig {
	sc := types.StorageConfig{}
	err := jsoniter.UnmarshalFromString(config, &sc)
	if err != nil || &sc == nil || sc.Disks == nil {
		return nil
	}
	resp := &response.StorageConfig{
		SystemDiskType: sc.Disks.SystemDisk.Category,
		SystemDiskSize: sc.Disks.SystemDisk.Size,
	}
	dataDisks := make([]response.DataDisk, 0)
	for _, disk := range sc.Disks.DataDisk {
		dd := response.DataDisk{
			DataDiskType: disk.Category,
			DataDiskSize: disk.Size,
		}
		dataDisks = append(dataDisks, dd)
	}
	resp.DataDisks = dataDisks
	resp.DataDiskNum = len(dataDisks)
	return resp
}

func getStringStatus(s constants.Status) string {
	switch s {
	case constants.Undefined:
		return "Undefined"
	case constants.Pending:
		return "Pending"
	case constants.Timeout:
		return "Timeout"
	case constants.Starting:
		return "Starting"
	case constants.Running:
		return "Running"
	case constants.Deleted:
		return "Deleted"
	case constants.Deleting:
		return "Deleting"
	}
	return ""
}
