package instance

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"
	"sync"

	"github.com/galaxy-future/BridgX/pkg/encrypt"

	"github.com/galaxy-future/BridgX/internal/gf-cluster/cluster"
	"github.com/galaxy-future/BridgX/internal/logs"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

//createInstance 创建实例
func createInstance(kubeCluster *cluster.KubernetesClient, request *gf_cluster.InstanceGroup, instanceName string) (*v1.Pod, error) {

	//获取资源
	cpu, err := strconv.ParseFloat(request.Cpu, 64)
	if err != nil {
		return nil, err
	}
	memory, err := strconv.ParseFloat(request.Memory, 64)
	if err != nil {
		return nil, err
	}
	disk, err := strconv.ParseFloat(request.Disk, 64)
	if err != nil {
		return nil, err
	}
	cpuLimit := resource.NewScaledQuantity(int64(cpu*1000), resource.Milli)
	memLimit := resource.NewScaledQuantity(int64(memory*1024), resource.Mega)
	diskLimit := resource.NewScaledQuantity(int64(disk*1024), resource.Mega)
	defaultImage := "galaxyfuture/centos-sshd:7.1"
	TerminationGracePeriodSeconds := int64(2)

	pwd, err := encrypt.AESDecrypt(encrypt.AesKeyPepper, request.SshPwd)
	if err != nil {
		return nil, err
	}
	req := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   instanceName,
			Labels: createInstanceLabels(request.Name, strconv.FormatInt(request.Id, 10)),
		},
		Spec: v1.PodSpec{
			TerminationGracePeriodSeconds: &TerminationGracePeriodSeconds,
			Containers: []v1.Container{
				{
					Name:  "instance",
					Image: defaultImage,
					Ports: []v1.ContainerPort{
						{
							Name:          "ssh",
							ContainerPort: 22,
						},
					},
					Env: []v1.EnvVar{
						{
							Name:  "PASSWORD",
							Value: pwd,
						},
					},
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{
							v1.ResourceCPU:              *cpuLimit,
							v1.ResourceMemory:           *memLimit,
							v1.ResourceEphemeralStorage: *diskLimit,
						},
						Requests: v1.ResourceList{
							v1.ResourceCPU:              *cpuLimit,
							v1.ResourceMemory:           *memLimit,
							v1.ResourceEphemeralStorage: *diskLimit,
						},
					},
				},
			},
		},
	}

	pod, err := kubeCluster.ClientSet.CoreV1().Pods("default").Create(context.Background(), req, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

//listElasticInstance 列出所有实例
func listElasticInstance(client *cluster.KubernetesClient, clusterName string, id int64) ([]*gf_cluster.Instance, error) {
	selector := metav1.LabelSelector{MatchLabels: createInstanceLabels(clusterName, strconv.FormatInt(id, 10))}
	pods, err := client.ClientSet.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{
		LabelSelector: labels.Set(selector.MatchLabels).String(),
	})

	if err == nil && errors.IsNotFound(err) {
		return make([]*gf_cluster.Instance, 0), nil
	}
	if err != nil {
		return nil, err
	}
	var instances []*gf_cluster.Instance
	for _, pod := range pods.Items {
		instances = append(instances, &gf_cluster.Instance{
			Name:   pod.Name,
			Ip:     pod.Status.PodIP,
			HostIp: pod.Status.HostIP,
		})
	}
	return instances, nil
}

func clearElasticInstance(client *cluster.KubernetesClient, instanceGroupName string, id int64) error {
	instances, err := listElasticInstance(client, instanceGroupName, id)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	var wg sync.WaitGroup
	for _, instance := range instances {
		wg.Add(1)
		go func(instance *gf_cluster.Instance) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					logs.Logger.Errorf("clearElasticInstance :%v ", r)
					logs.Logger.Errorf("clearElasticInstance panic", zap.String("stack", string(debug.Stack())))
				}
			}()
			err := client.ClientSet.CoreV1().Pods("default").Delete(context.Background(), instance.Name, metav1.DeleteOptions{})
			if err != nil {
				logs.Logger.Errorw("failed to delete pod.", zap.String("instance_group_name", instanceGroupName), zap.String("instance_name", instance.Name), zap.Error(err))
			}
		}(instance)
	}
	wg.Wait()
	return nil
}

func generateInstanceName(name string, index int) string {
	return fmt.Sprintf("%s-%d", name, index)
}

func createInstanceLabels(name string, id string) map[string]string {
	return map[string]string{
		gf_cluster.ClusterTypeKey:            gf_cluster.ClusterTypeValue,
		gf_cluster.ClusterInstanceGroupKey:   name,
		gf_cluster.ClusterInstanceGroupIdKey: id,
	}
}
