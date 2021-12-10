package cluster

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//ListClusterNodeSummary 获得集群下所有节点详情
func ListClusterNodeSummary(clusterId int64) (gf_cluster.ClusterNodeSummaryArray, error) {
	cluster, err := model.GetKubernetesCluster(clusterId)
	if err != nil {
		return nil, err
	}
	nodes, err := getClusterNodeInfo(cluster)
	if err != nil {
		return nil, err
	}

	pods, err := getClusterPodInfo(cluster)
	if err != nil {
		return nil, err
	}

	calcNodeResourceUsage(nodes, pods)

	return nodes, nil

}

func getClusterNodeInfo(info *gf_cluster.KubernetesInfo) ([]*gf_cluster.ClusterNodeSummary, error) {

	if info.Status != gf_cluster.KubernetesStatusRunning {
		return make([]*gf_cluster.ClusterNodeSummary, 0), nil
	}
	client, err := GetKubeClient(info.Id)
	if err != nil {
		return nil, err
	}

	//TODO remove hadrd code
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	nodes, err := client.ClientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("从集群中查询节点信息失败，失败原因：%s", err.Error())
	}

	var nodesSummary []*gf_cluster.ClusterNodeSummary
	for _, node := range nodes.Items {
		var ipAddress, hostName string
		for _, address := range node.Status.Addresses {
			if address.Type == v1.NodeHostName {
				hostName = address.Address
			}
			if address.Type == v1.NodeInternalIP {
				ipAddress = address.Address
			}
		}

		nodeCpu := *node.Status.Capacity.Cpu()
		nodeMemory := *node.Status.Capacity.Memory()
		nodeStorage := *node.Status.Capacity.StorageEphemeral()

		//
		//kubeUsed ,exist := systemUsed[ipAddress]
		//if exist {
		//	nodeCpu.Sub(*kubeUsed.Cpu)
		//	nodeMemory.Sub(*kubeUsed.Memory)
		//	nodeStorage.Sub(*kubeUsed.Storage)
		//}

		cpuSize := cpuQuantity2Float(nodeCpu)
		memorySize := storageQuantity2Float(nodeMemory)
		storageSize := storageQuantity2Float(nodeStorage)

		role, exists := node.Labels[gf_cluster.KubernetesRoleKey]
		if !exists {
			role = gf_cluster.KubernetesRoleWorker
		} else {
			role = gf_cluster.KubernetesRoleMaster
		}

		nodeStatus := "notReady"
		for _, condition := range node.Status.Conditions {
			if condition.Type == v1.NodeReady && condition.Status == v1.ConditionTrue {
				nodeStatus = "Ready"
			}
		}

		nodesSummary = append(nodesSummary, &gf_cluster.ClusterNodeSummary{
			Status:        nodeStatus,
			IpAddress:     ipAddress,
			HostName:      hostName,
			ClusterName:   info.BridgxClusterName,
			AllCpuCores:   int(cpuSize),
			AllMemoryGi:   memorySize,
			AllDiskGi:     storageSize,
			FreeCpuCores:  cpuSize,
			FreeMemoryGi:  memorySize,
			FreeDiskGi:    storageSize,
			MachineType:   "",
			CloudProvider: info.CloudType,
			Role:          role,
		})
	}

	calcNodePodCounts(nodesSummary, info)

	return nodesSummary, nil

}

func calcNodePodCounts(nodesSummary []*gf_cluster.ClusterNodeSummary, info *gf_cluster.KubernetesInfo) {
	pods, err := getClusterPodInfo(info)
	if err != nil {
		for _, nodeSummary := range nodesSummary {
			if nodeSummary.Message != "" {
				nodeSummary.Message = "获取节点pod实例时失败"
			}
		}
	}
	summaries := make(map[string]*gf_cluster.ClusterNodeSummary)
	for _, nodeSummary := range nodesSummary {
		summaries[nodeSummary.IpAddress] = nodeSummary
	}

	for _, pod := range pods {
		nodeSummary, exist := summaries[pod.NodeIp]
		if !exist {
			logs.Logger.Error("没有查询到指定pod所属节点")
			continue
		}
		nodeSummary.PodCount++
	}
}

type PodResources struct {
	Cpu     *resource.Quantity `json:"cpu"`
	Memory  *resource.Quantity `json:"memory"`
	Storage *resource.Quantity `json:"storage"`
}

//ListKubeSystemInstance 列出kuber-system所占用资源
func ListKubeSystemInstance(info *gf_cluster.KubernetesInfo) (map[string]*PodResources, error) {
	if info.Status != gf_cluster.KubernetesStatusRunning {
		return nil, nil
	}
	client, err := GetKubeClient(info.Id)
	if err != nil {
		return nil, err
	}
	pods, err := client.ClientSet.CoreV1().Pods("kube-system").List(context.Background(), metav1.ListOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	allocatedResource := make(map[string]*PodResources)
	for _, pod := range pods.Items {
		if pod.Status.Phase == v1.PodRunning {
			item, exist := allocatedResource[pod.Status.HostIP]
			if !exist {
				item = &PodResources{
					Cpu:     resource.NewScaledQuantity(0, resource.Milli),
					Memory:  resource.NewScaledQuantity(0, resource.Mega),
					Storage: resource.NewScaledQuantity(0, resource.Giga),
				}
				allocatedResource[pod.Status.HostIP] = item
			}
			for _, container := range pod.Spec.Containers {
				item.Cpu.Add(*container.Resources.Limits.Cpu())
				item.Memory.Add(*container.Resources.Limits.Memory())
				item.Storage.Add(*container.Resources.Limits.StorageEphemeral())
			}
		}
	}

	return allocatedResource, nil
}
