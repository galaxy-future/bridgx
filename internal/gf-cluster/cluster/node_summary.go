package cluster

import (
	"context"
	"fmt"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/pkg/gf-cluster"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
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
		return nil, nil
	}
	client, err := GetKubeClient(info.Id)
	if err != nil {
		return nil, err
	}

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

		cpuSize := cpuQuantity2Float(*node.Status.Capacity.Cpu())
		memorySize := storageQuantity2Float(*node.Status.Capacity.Memory())
		storageSize := storageQuantity2Float(*node.Status.Capacity.StorageEphemeral())

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
