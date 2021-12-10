package instance

import (
	"context"
	"sync"
	"time"

	"github.com/galaxy-future/BridgX/internal/gf-cluster/cluster"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExpandCustomInstanceGroup 扩容实例组
func ExpandCustomInstanceGroup(instanceGroup *gf_cluster.InstanceGroup, count int) error {
	// 1 获取k8s实例列表
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
		// 2 创建k8s实例
		pod, err := createInstance(client, instanceGroup, name)
		if err != nil {
			logs.Logger.Errorw("failed to expand instance", zap.String("instance_group_name", instanceGroup.Name), zap.String("instance_name", name), zap.Error(err))
			continue
		}
		existInstances = append(existInstances, &gf_cluster.Instance{
			Name: pod.Name,
			Ip:   pod.Status.PodIP,
		})
		logs.Logger.Infow("expand instance success", zap.String("instance_group_name", instanceGroup.Name), zap.String("instance_name", name))
	}
	// 3 更新实例组实例数
	err = model.UpdateInstanceGroupInstanceCountFromDB(count, instanceGroup.Id)
	if err != nil {
		return err
	}
	return nil
}

// ShrinkCustomInstanceGroup 缩容实例组
func ShrinkCustomInstanceGroup(instanceGroup *gf_cluster.InstanceGroup, count int) error {
	// 1 缩容数量
	shrinkCount := instanceGroup.InstanceCount - count
	// 2 获取k8s实例列表
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
		// 3 执行缩容
		wg.Add(1)
		go func(instance *gf_cluster.Instance) {
			defer wg.Done()
			err := client.ClientSet.CoreV1().Pods("default").Delete(context.Background(), instance.Name, v1.DeleteOptions{})
			if err != nil {
				logs.Logger.Errorw("failed to shrink instance.", zap.String("instance_group_name", instanceGroup.Name), zap.String("instance_name", instance.Name), zap.Error(err))
				return
			}
			logs.Logger.Infow("shrink instance success", zap.String("instance_group_name", instanceGroup.Name), zap.String("instance_name", instance.Name))
		}(instance)
	}
	wg.Wait()
	// 4 更新实例组实例数
	err = model.UpdateInstanceGroupInstanceCountFromDB(count, instanceGroup.Id)
	if err != nil {
		return err
	}
	return nil
}

//ListCustomInstances 列出所有eci
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

//RestartInstance 重启实例
func RestartInstance(instanceGroupId int64, name string) error {
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
			logs.Logger.Errorw("server run failed ", zap.Error(err))
		}
	}()
	return nil
}

//DeleteInstance 删除实例
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
