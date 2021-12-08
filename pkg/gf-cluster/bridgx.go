package gf_cluster

const (
	//UsageKey 集群的占用的Tag
	UsageKey = "usage"
	//UnusedValue 未被占用的TagValue
	UnusedValue = "unused"
	//GalaxyfutureCloudUsage 星云集群占用Value
	GalaxyfutureCloudUsage = "galaxyfuture-cloud"
)

//ListBridgxClusterByTagRequest 获取集群请求的Request
type ListBridgxClusterByTagRequest struct {
	Tags       map[string]string `json:"tags"`
	PageNumber int               `json:"page_number"`
	PageSize   int               `json:"page_size"`
}

//BridgxClusterName 集群名称信息
type BridgxClusterName struct {
	ClusterId   string            `json:"cluster_id"`
	ClusterName string            `json:"cluster_name"`
	Provider    string            `json:"provider"`
	Tags        map[string]string `json:"tags"`
	CreateAt    string            `json:"create_at"`
	CreateBy    string            `json:"create_by"`
}

// ListBirdgxClusterData 列出集群中返回体数据
type ListBirdgxClusterData struct {
	ClusterList []*BridgxClusterName `json:"cluster_list"`
	Pager       Pager                `json:"pager"`
}

//ListBirdgxClusterByClusterResponse Bridgx返回
type ListBirdgxClusterByClusterResponse struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data *ListBirdgxClusterData `json:"data"`
}

//EditBridgxClusterTagRequst 编辑tag请求
type EditBridgxClusterTagRequst struct {
	Tags        map[string]string `json:"tags"`
	ClusterName string            `json:"cluster_name"`
}

//EditBridgxClusterTagResponse 编辑请求返回
type EditBridgxClusterTagResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

//GetBridgxClusterInstanceRequest 获得Bridgx信息的Request
type GetBridgxClusterInstanceRequest struct {
	ClusterName string `json:"cluster_name"`
	Status      string `json:"status"`
}

//BridgxInstance Bridgx集群实例
type BridgxInstance struct {
	InstanceId   string `json:"instance_id"`
	IpInner      string `json:"ip_inner"`
	IpOuter      string `json:"ip_outer"`
	Provider     string `json:"provider"`
	ClusterName  string `json:"cluster_name"`
	InstanceType string `json:"instance_type"`
	CreateAt     string `json:"create_at"`
	Status       string `json:"status"`
}

//GetBridgxClusterInstanceData 获取集群data
type GetBridgxClusterInstanceData struct {
	InstanceList []*BridgxInstance `json:"instance_list"`
	Pager        Pager             `json:"pager"`
}

//GetBridgxClusterInstanceResponse 获得集群实例返回
type GetBridgxClusterInstanceResponse struct {
	Code int                           `json:"code"`
	Msg  string                        `json:"msg"`
	Data *GetBridgxClusterInstanceData `json:"data"`
}

//BridgxClusterDetails 集群相信信息
type BridgxClusterDetails struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	RegionId     string `json:"region_id"`
	ZoneId       string `json:"zone_id"`
	InstanceType string `json:"instance_type"`
	ChargeType   string `json:"charge_type"`
	Image        string `json:"image"`
	Provider     string `json:"provider"`
	Password     string `json:"password"`
}

//BridgxClusterDetailsResponse 集群相信信息response
type BridgxClusterDetailsResponse struct {
	Code int                   `json:"code"`
	Msg  string                `json:"msg"`
	Data *BridgxClusterDetails `json:"data"`
}

//AKSKData 集群AKSK数据
type AKSKData struct {
	AccountName          string `json:"account_name"`
	AccountKey           string `json:"account_key"`
	AccountSecretEncrypt string `json:"account_secret_encrypt"`
	Provider             string `json:"provider"`
}

//GetAKSKResponse 获得AKSKresponse
type GetAKSKResponse struct {
	Code int       `json:"code"`
	Data *AKSKData `json:"data"`
	Msg  string    `json:"msg"`
}

//BridgxUnusedCluster we未被占用的集群
type BridgxUnusedCluster struct {
	ClusterName string   `json:"cluster_name"`
	CloudType   string   `json:"cloud_type"`
	Nodes       []string `json:"nodes"`
}

//ListUnusedBridgxClusterResponse 未被占用集群list相应体
type ListUnusedBridgxClusterResponse struct {
	*ResponseBase
	Pager
	Clusters []*BridgxUnusedCluster `json:"clusters"`
}

//NewListUnusedBridgxClusterResponse 新建未被占用的集群的相应体
func NewListUnusedBridgxClusterResponse(clusters []*BridgxUnusedCluster, pager Pager) *ListUnusedBridgxClusterResponse {
	return &ListUnusedBridgxClusterResponse{
		ResponseBase: NewSuccessResponse(),
		Pager:        pager,
		Clusters:     clusters,
	}
}
