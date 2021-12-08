package gf_cluster

type PodSummary struct {
	NodeName string `json:"node_name"`
	NodeIp   string `json:"node_ip"`

	PodName string `json:"pod_name"`
	PodIP   string `json:"pod_ip"`

	//AllocatedCpuCores 已经分配cpu数量
	AllocatedCpuCores float64 `json:"allocated_cpu_cores"`
	//AllocatedMemoryGi 已经分配内存
	AllocatedMemoryGi float64 `json:"allocated_memory_gi"`
	//AllocatedDiskGi 已经分配磁盘大小
	AllocatedDiskGi float64 `json:"allocated_disk_gi"`
	//GroupName 组名
	GroupName string `json:"group_name"`
	//RunningTime 运行时间
	RunningTime string `json:"running_time"`
	//Status 状态
	Status  string `json:"status"`
	GroupId int64  `json:"group_id"`
}
type ClusterPodsSummaryArray []*PodSummary

func (array ClusterPodsSummaryArray) Len() int {
	return len(array)
}
func (array ClusterPodsSummaryArray) Less(i, j int) bool {
	return array[i].GroupName < array[j].GroupName && array[i].PodName < array[i].PodName
}
func (array ClusterPodsSummaryArray) Swap(i, j int) {
	array[i], array[j] = array[j], array[i]
}

//ListClusterPodDetail 节点详细信息
type ListClusterPodDetail struct {
	ClusterPodsSummaryArray `json:"pods_information"`
}

//ListClusterPodsDetailResponse 返回集群pod信息
type ListClusterPodsDetailResponse struct {
	*ResponseBase
	Pager
	Pods ClusterPodsSummaryArray `json:"pods"`
}

func NewListClusterPodsDetailResponse(array ClusterPodsSummaryArray, pager Pager) *ListClusterPodsDetailResponse {
	return &ListClusterPodsDetailResponse{
		ResponseBase: NewSuccessResponse(),
		Pager:        pager,
		Pods:         array,
	}
}

type CreateClusterRequest struct {
	ClusterName     string         `json:"cluster_name"`
	BridgxClusterId int64          `json:"bridgx_cluster_id"`
	Type            KubernetesType `json:"type"`
}
