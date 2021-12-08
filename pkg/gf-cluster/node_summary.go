package gf_cluster

//ClusterNodeSummary 节点综述
type ClusterNodeSummary struct {
	//Status节点状态
	Status string `json:"status"`
	//IpAddress ip信息
	IpAddress string `json:"ip_address"`
	//HostName 机器名称
	HostName string `json:"host_name"`
	//ClusterName 集群名称
	ClusterName string `json:"cluster_name"`
	//AllCpuCores 所有cpu数量
	AllCpuCores int `json:"all_cpu_cores"`
	//AllocatedCpuCores 已经分配cpu数量
	FreeCpuCores float64 `json:"free_cpu_cores"`
	//AllMemoryGi 集群内存总数
	AllMemoryGi float64 `json:"all_memory_gi"`
	//AllocatedMemoryGi 已经分配内存
	FreeMemoryGi float64 `json:"free_memory_gi"`
	//AllDiskGi 磁盘大小
	AllDiskGi float64 `json:"all_disk_gi"`
	//AllocatedDiskGi 已经分配磁盘大小
	FreeDiskGi float64 `json:"free_disk_gi"`
	//MachineType 机器型号
	MachineType string `json:"machine_type"`
	//CloudProvider 为所属云厂商
	CloudProvider string `json:"cloud_provider"`
	//PodCount 实例中pod数量
	PodCount int `json:"pod_count"`
	//master/node
	Role string `json:"role"`
	//Message 获取节点失败信息，正常返回为success，异常时返回出错信息
	Message string `json:"message"`
}

type ClusterNodeSummaryArray []*ClusterNodeSummary

func (array ClusterNodeSummaryArray) Len() int {
	return len(array)
}
func (array ClusterNodeSummaryArray) Less(i, j int) bool {
	return array[i].ClusterName < array[j].ClusterName && array[i].IpAddress < array[i].IpAddress
}
func (array ClusterNodeSummaryArray) Swap(i, j int) {
	array[i], array[j] = array[j], array[i]
}

//ListClusterNodesResponse get api/v1/kubernetes/nodes get response
type ListClusterNodesResponse struct {
	*ResponseBase
	Pager
	Nodes ClusterNodeSummaryArray `json:"nodes"`
}

func NewListClusterNodesResponse(nodes []*ClusterNodeSummary, page Pager) *ListClusterNodesResponse {
	return &ListClusterNodesResponse{
		ResponseBase: NewSuccessResponse(),
		Pager:        page,
		Nodes:        nodes,
	}
}
