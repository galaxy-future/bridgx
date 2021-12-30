package gf_cluster

type ListClusterLogsResponse struct {
	*ResponseBase
	Logs []KubernetesInstallStep `json:"logs"`
}

type KubernetesInstallStep struct {
	Id           int64  `json:"id"`
	KubernetesId int64  `json:"kubernetes_id"`
	HostIp       string `json:"host_ip"`
	Operation    string `json:"operation"`
	Message      string `json:"message"`
}
