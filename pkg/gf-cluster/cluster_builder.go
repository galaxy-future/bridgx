package gf_cluster

type BuildMode int

const (
	//SingleMode da单节点模式
	SingleMode BuildMode = iota
	//ClusterMode HA模式
	ClusterMode
	//ClusterUnknown 未知
	ClusterUnknown
)

//ClusterBuilderParams  集群创建相关参数
type ClusterBuilderParams struct {
	PodCidr      string                `json:"pod_cidr"`
	SvcCidr      string                `json:"svc_cidr"`
	MachineList  []ClusterBuildMachine `json:"machine_list"`
	Mode         BuildMode             `json:"mode"`
	KubernetesId int64                 `json:"kubernetes_id"`
	AccessKey    string                `json:"ak"`
	AccessSecret string                `json:"sk"`
}

//ClusterBuildMachine 创建集群所用的物理机实体
type ClusterBuildMachine struct {
	IP       string            `json:"ip"`
	Hostname string            `json:"hostname"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	Labels   map[string]string `json:"labels"`
}

//String2BuildMode 字符串到创建模式转换
func String2BuildMode(mode string) BuildMode {
	switch mode {
	case KubernetesStandalone:
		return SingleMode
	case KubernetesHA:
		return ClusterMode
	}

	return ClusterUnknown
}
