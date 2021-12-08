package instance

import (
	"context"
	"github.com/galaxy-future/BridgX/internal/gf-cluster/cluster"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"sync"
	"time"
)

func ExpandCustomInstanceGroup(instanceGroup *gf_cluster.InstanceGroup, count int) error {
	client, err := cluster.GetKubeClient(instanceGroup.KubernetesId)
	if err != nil {
		return err
	}
	existInstances, err := listElasticInstance(client, instanceGroup.Name, instanceGroup.Id)
	if err != nil {
		return err
	}

	existsMap := make(map[string]struct{})
	for _, instance := range existInstances {
		existsMap[instance.Name] = struct{}{}
	}
	for i := 0; i < count; i++ {
		if len(existInstances) >= count {
			break
		}
		name := generateInstanceName(instanceGroup.Name, i)
		if _, exist := existsMap[name]; exist {
			continue
		}
		pod, err := createInstance(client, instanceGroup, name)
		if err != nil {
			logs.Logger.Error("failed to expand instance", zap.String("instance_group_name", instanceGroup.Name), zap.String("instance_name", name), zap.Error(err))
			continue
		}
		existInstances = append(existInstances, &gf_cluster.Instance{
			Name: pod.Name,
			Ip:   pod.Status.PodIP,
		})
		logs.Logger.Info("expand instance success", zap.String("instance_group_name", instanceGroup.Name), zap.String("instance_name", name))
	}
	err = model.UpdateInstanceGroupInstanceCountFromDB(count, instanceGroup.Id)
	if err != nil {
		return err
	}
	return nil
}

func ShrinkCustomInstanceGroup(instanceGroup *gf_cluster.InstanceGroup, count int) error {
	shrinkCount := instanceGroup.InstanceCount - count
	client, err := cluster.GetKubeClient(instanceGroup.KubernetesId)
	if err != nil {
		return err
	}
	existInstances, err := listElasticInstance(client, instanceGroup.Name, instanceGroup.Id)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, instance := range existInstances {
		if shrinkCount <= 0 {
			break
		}
		shrinkCount--
		wg.Add(1)
		go func(instance *gf_cluster.Instance) {
			defer wg.Done()
			err := client.ClientSet.CoreV1().Pods("default").Delete(context.Background(), instance.Name, v1.DeleteOptions{})
			if err != nil {
				logs.Logger.Error("failed to shrink instance.", zap.String("instance_group_name", instanceGroup.Name), zap.String("instance_name", instance.Name), zap.Error(err))
				return
			}
			logs.Logger.Info("shrink instance success", zap.String("instance_group_name", instanceGroup.Name), zap.String("instance_name", instance.Name))
		}(instance)
	}
	wg.Wait()
	err = model.UpdateInstanceGroupInstanceCountFromDB(count, instanceGroup.Id)
	if err != nil {
		return err
	}
	return nil
}

func ClearCustomECIClusterInstances(instanceGroupId int64) error {

	instanceGroup, err := GetInstanceGroup(instanceGroupId)
	if err != nil {
		return err
	}

	client, err := cluster.GetKubeClient(instanceGroup.KubernetesId)
	if err != nil {
		return err
	}

	return clearElasticInstance(client, instanceGroup.Name, instanceGroupId)
}

//ListCustomInstances 列出所有eci
//TODO fe分页
func ListCustomInstances(instanceGroupId int64) ([]*gf_cluster.Instance, error) {
	instanceGroup, err := GetInstanceGroup(instanceGroupId)
	if err != nil {
		return nil, err
	}

	client, err := cluster.GetKubeClient(instanceGroup.KubernetesId)
	if err != nil {
		return nil, err
	}
	return listElasticInstance(client, instanceGroup.Name, instanceGroupId)
}

func RestartInstanceGroup(instanceGroupId int64, name string) error {
	instanceGroup, err := GetInstanceGroup(instanceGroupId)
	if err != nil {
		return err
	}
	client, err := cluster.GetKubeClient(instanceGroup.KubernetesId)
	if err != nil {
		return err
	}
	err = client.ClientSet.CoreV1().Pods("default").Delete(context.Background(), name, v1.DeleteOptions{})
	if err != nil {
		return err
	}
	go func() {
		time.Sleep(time.Duration(2) * time.Second)
		_, err = createInstance(client, instanceGroup, name)
		if err != nil {
			logs.Logger.Error("server run failed ", zap.Error(err))
		}
	}()
	return nil
}

func DeleteInstance(instanceGroup *gf_cluster.InstanceGroup, name string) error {
	client, err := cluster.GetKubeClient(instanceGroup.KubernetesId)
	if err != nil {
		return err
	}
	err = client.ClientSet.CoreV1().Pods("default").Delete(context.Background(), name, v1.DeleteOptions{})
	if err != nil {
		return err
	}
	err = model.UpdateInstanceGroupInstanceCountFromDB(instanceGroup.InstanceCount-1, instanceGroup.Id)
	if err != nil {
		return err
	}
	return nil
}
