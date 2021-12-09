package cluster

import (
	"fmt"
	"time"

	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/resource"
)

//ListClustersSummary 列出所有集群summary
func ListClustersSummary(id string, name string, pageNumber int, pageSize int) ([]*gf_cluster.ClusterSummary, int, error) {
	clusters, total, err := model.ListKubernetesClusters(id, name, pageNumber, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("查询集群信息失败，错误信息: %v", err.Error())
	}

	var summaries []*gf_cluster.ClusterSummary
	for _, cluster := range clusters {

		summary, err := queryClusterInfo(cluster)
		if err != nil {
			summary.Message = err.Error()
		} else {
			summary.Message = gf_cluster.StatusSuccess
		}
		summaries = append(summaries, summary)

	}

	return summaries, total, nil

}

//GetClustersSummary 获得集群概述
func GetClustersSummary(clusterId int64) (*gf_cluster.ClusterSummary, error) {
	cluster, err := model.GetKubernetesCluster(clusterId)
	if err != nil {
		return nil, err
	}

	summary, err := queryClusterInfo(cluster)
	if err != nil {
		summary.Message = err.Error()
	} else {
		summary.Message = gf_cluster.StatusSuccess
	}

	return summary, nil
}

func queryClusterInfo(cluster *gf_cluster.KubernetesInfo) (*gf_cluster.ClusterSummary, error) {

	summary := &gf_cluster.ClusterSummary{
		ClusterId:                   cluster.Id,
		ClusterName:                 cluster.Name,
		AllCpuCores:                 0,
		FreeCpuCores:                0,
		AllMemoryGi:                 0,
		FreeMemoryGi:                0,
		AllDiskGi:                   0,
		FreeDiskGi:                  0,
		MaxUnallocatedCpuInNode:     0,
		MaxUnallocatedMemoryInNode:  0,
		MaxUnallocatedStorageInNode: 0,
		PodCount:                    0,
		WorkCount:                   0,
		MasterCount:                 0,
		CreatedUser:                 cluster.CreatedUser,
		CreatedTime:                 time.Unix(cluster.CreatedTime, 0).Format("2006-01-02 15:04:05"),
		Status:                      cluster.Status,
		InstallStep:                 cluster.InstallStep,
	}

	nodes, err := getClusterNodeInfo(cluster)
	if err != nil {
		return summary, err
	}

	pods, err := getClusterPodInfo(cluster)
	if err != nil {
		return summary, err
	}

	calcNodeResourceUsage(nodes, pods)
	maxFreeCpuInNode := 0.0
	maxMemoryInNode := 0.0
	maxStorageInNode := 0.0
	workerCount := 0
	masterCount := 0
	for _, node := range nodes {
		summary.AllCpuCores += node.AllCpuCores
		summary.AllMemoryGi += node.AllMemoryGi
		summary.AllDiskGi += node.AllDiskGi
		summary.FreeCpuCores += node.FreeCpuCores
		summary.FreeMemoryGi += node.FreeMemoryGi
		summary.FreeDiskGi += node.FreeDiskGi
		if node.FreeCpuCores > maxFreeCpuInNode {
			maxFreeCpuInNode = node.FreeCpuCores
		}
		if node.FreeMemoryGi > maxMemoryInNode {
			maxMemoryInNode = node.FreeMemoryGi
		}
		if node.FreeDiskGi > maxStorageInNode {
			maxStorageInNode = node.FreeDiskGi
		}
		if node.Role == gf_cluster.KubernetesRoleWorker {
			workerCount++
		}
		if node.Role == gf_cluster.KubernetesRoleMaster {
			masterCount++
		}
	}

	summary.MaxUnallocatedCpuInNode = maxFreeCpuInNode
	summary.MaxUnallocatedMemoryInNode = maxMemoryInNode
	summary.MaxUnallocatedStorageInNode = maxStorageInNode
	summary.PodCount = int64(len(pods))
	summary.WorkCount = workerCount
	summary.MasterCount = masterCount

	return summary, nil

}

func calcNodeResourceUsage(nodes []*gf_cluster.ClusterNodeSummary, pods []*gf_cluster.PodSummary) {
	nodeResourceMap := make(map[string]*gf_cluster.ClusterNodeSummary)
	for _, node := range nodes {
		nodeResourceMap[node.HostName] = node
	}

	for _, pod := range pods {
		node, exist := nodeResourceMap[pod.NodeName]
		if !exist {
			logs.Logger.Error("can not find pod host", zap.String("nodeName", pod.NodeName))
			continue
		}
		node.FreeCpuCores -= pod.AllocatedCpuCores
		node.FreeMemoryGi -= pod.AllocatedMemoryGi
		node.FreeDiskGi -= pod.AllocatedDiskGi
	}

}

func createInstanceLabels() map[string]string {
	return map[string]string{
		gf_cluster.ClusterTypeKey: gf_cluster.ClusterTypeValue,
	}
}

func cpuQuantity2Float(cpu resource.Quantity) float64 {
	memoryContent := cpu.ScaledValue(resource.Milli)
	return float64(memoryContent) / 1000.0
}

func storageQuantity2Float(mem resource.Quantity) float64 {
	memoryContent := mem.ScaledValue(resource.Mega)
	return float64(memoryContent) / 1024.0
}
