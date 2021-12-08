package gf_cluster

type KubernetesType string

//KubernetesInfo 代表一个Kubernetes集群
type KubernetesInfo struct {
	Id                int64  `json:"id"`
	Name              string `json:"name"`
	Region            string `json:"region"`
	CloudType         string `json:"cloud_type"`
	Status            string `json:"status"`
	Config            string `json:"config"` //kube config
	InstallStep       string `json:"install_step"`
	Message           string `json:"message"`
	BridgxClusterName string `json:"bridgx_cluster_name"`
	Type              string `json:"type"`
	CreatedUser       string `json:"created_user"`
	CreatedTime       int64  `json:"created_time"`
}

type KubernetesInfoListResponse struct {
	*ResponseBase
	Clusters []*KubernetesInfo `json:"clusters"`
}

func NewKubernetesInfoListResponse(clusters []*KubernetesInfo) *KubernetesInfoListResponse {
	return &KubernetesInfoListResponse{
		ResponseBase: NewSuccessResponse(),
		Clusters:     clusters,
	}
}

type KubernetesInfoGetResponse struct {
	*ResponseBase
	Cluster *KubernetesInfo
}

func NewKubernetesInfoGetResponse(cluster *KubernetesInfo) *KubernetesInfoGetResponse {
	return &KubernetesInfoGetResponse{
		ResponseBase: NewSuccessResponse(),
		Cluster:      cluster,
	}
}
