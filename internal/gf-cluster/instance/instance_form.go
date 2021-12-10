package instance

import (
	"time"

	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
)

func AddInstanceForm(instanceGroup *gf_cluster.InstanceGroup, cost int64, createdUserId int64, createdUserName string, opt string, updatedInstanceCount int, err error) error {
	executeStatus := gf_cluster.InstanceInit
	if err == nil {
		executeStatus = gf_cluster.InstanceNormal
	}
	if err != nil {
		executeStatus = gf_cluster.InstanceError
	}
	kubernetes, err := model.GetKubernetesCluster(instanceGroup.KubernetesId)
	if err != nil {
		return err
	}
	instanceForms := gf_cluster.InstanceForm{
		Id:                   0,
		ExecuteStatus:        executeStatus,
		InstanceGroup:        instanceGroup.Name,
		Cpu:                  instanceGroup.Cpu,
		Memory:               instanceGroup.Memory,
		Disk:                 instanceGroup.Disk,
		OptType:              opt,
		UpdatedInstanceCount: updatedInstanceCount,
		HostTime:             cost,
		CreatedUserId:        createdUserId,
		CreatedUserName:      createdUserName,
		CreatedTime:          time.Now().Unix(),
		ClusterName:          kubernetes.Name,
	}

	err = model.CreateInstanceFormFromDB(&instanceForms)
	if err != nil {
		return err
	}
	return err
}
