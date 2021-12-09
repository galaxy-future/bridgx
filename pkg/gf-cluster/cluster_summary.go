package gf_cluster

//ClusterSummary 集群概要信息
type ClusterSummary struct {
	//ClusterId 集群ID
	ClusterId int64 `json:"cluster_id"`
	//ClusterName  集群名称
	ClusterName string `json:"cluster_name"`
	//AllCpuCores 所有cpu数量
	AllCpuCores int `json:"all_cpu_cores"`
	//FreeCpuCores 剩余cpu数量
	FreeCpuCores float64 `json:"free_cpu_cores"`
	//AllMemoryGi 集群内存总数
	AllMemoryGi float64 `json:"all_memory_gi"`
	//FreeMemoryG 剩余内存
	FreeMemoryGi float64 `json:"free_memory_gi"`
	//AllDiskGi 磁盘大小
	AllDiskGi float64 `json:"all_disk_gi"`
	//FreeDiskGi 剩余磁盘
	FreeDiskGi float64 `json:"free_disk_gi"`

	//AllocatedDiskGi 最大可分配cpu
	MaxUnallocatedCpuInNode float64 `json:"max_unallocated_cpu_in_node"`
	//MaxUnallocatedMemoryInNode 最大未分配内存
	MaxUnallocatedMemoryInNode float64 `json:"max_unallocated_memory_in_node"`
	//MaxUnallocatedStorageInNode  最大没有分配存储
	MaxUnallocatedStorageInNode float64 `json:"max_unallocated_storage_in_node"`

	//PodCount pod总数
	PodCount int64 `json:"pod_count"`
	//WorkCount work节点数量
	WorkCount int `json:"work_count"`
	//MasterCount master节点数量
	MasterCount int `json:"master_count"`
	//CreatedUser
	CreatedUser string `json:"created_user"`
	//CreatedTime
	CreatedTime string `json:"created_time"`
	//Status 集群状态
	Status string `json:"status"`
	//Message 其他信息，当集群异常时，相关的额错误信息
	Message     string `json:"message"`
	InstallStep string `json:"install_step"`
}

//ListClusterSummaryResponse get api/v1/kubernetes/summary get response
type ListClusterSummaryResponse struct {
	*ResponseBase
	Pager
	Clusters []*ClusterSummary `json:"clusters"`
}

//NewListClusterSummaryResponse 新建正常返回结构
func NewListClusterSummaryResponse(clusters []*ClusterSummary, pager Pager) *ListClusterSummaryResponse {
	return &ListClusterSummaryResponse{
		ResponseBase: NewSuccessResponse(),
		Pager:        pager,
		Clusters:     clusters,
	}
}

//GetClusterSummaryResponse 获取单个集群信息
type GetClusterSummaryResponse struct {
	*ResponseBase
	Cluster *ClusterSummary `json:"cluster"`
}

//NewGetClusterSummaryResponse 新建正常返回结构
func NewGetClusterSummaryResponse(cluster *ClusterSummary) *GetClusterSummaryResponse {
	return &GetClusterSummaryResponse{
		ResponseBase: NewSuccessResponse(),
		Cluster:      cluster,
	}
}

//BridgxClusterBuildRequest 集群构建请求
type BridgxClusterBuildRequest struct {
	ClusterName       string `json:"cluster_name"`
	BridgxClusterName string `json:"bridgx_cluster_name"`
	PodCidr           string `json:"pod_cidr"`
	ServiceCidr       string `json:"service_cidr"`
	ClusterType       string `json:"cluster_type"`
}
