package instance

import (
	"github.com/galaxy-future/BridgX/internal/gf-cluster/cluster"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/pkg/gf-cluster"
)

func CreateInstanceGroup(instanceGroup *gf_cluster.InstanceGroup) error {
	return model.CreateInstanceGroupFromDB(instanceGroup)
}

func DeleteInstanceGroup(instanceGroup *gf_cluster.InstanceGroup) error {

	client, err := cluster.GetKubeClient(instanceGroup.KubernetesId)
	if err != nil {
		return err
	}
	err = model.DeleteInstanceGroupFromDB(instanceGroup.Id)
	if err != nil {
		return err
	}
	err = clearElasticInstance(client, instanceGroup.Name, instanceGroup.Id)
	if err != nil {
		return err
	}

	return nil
}

func GetInstanceGroup(instanceGroupId int64) (*gf_cluster.InstanceGroup, error) {
	return model.GetInstanceGroupFromDB(instanceGroupId)
}
