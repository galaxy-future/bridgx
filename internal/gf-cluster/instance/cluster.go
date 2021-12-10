package instance

import (
	"github.com/galaxy-future/BridgX/internal/gf-cluster/cluster"
	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
)

//CreateInstanceGroup 新建实例组
func CreateInstanceGroup(instanceGroup *gf_cluster.InstanceGroup) error {
	return model.CreateInstanceGroupFromDB(instanceGroup)
}

//DeleteInstanceGroup 删除实例组
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

//GetInstanceGroup 获取实例组
func GetInstanceGroup(instanceGroupId int64) (*gf_cluster.InstanceGroup, error) {
	return model.GetInstanceGroupFromDB(instanceGroupId)
}
