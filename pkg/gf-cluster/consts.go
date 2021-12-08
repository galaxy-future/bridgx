package gf_cluster

const (
	DefaultPageSize = 10
)

//Kubernetes step
const (
	KubernetesStepInitializeCluster = "Initialize the cluster"
	KubernetesStepInstallMaster     = "Install master: "
	KubernetesStepInstallFlannel    = "Install flannel"
	KubernetesStepInstallNode       = "Install Node: "
	KubernetesStepDone              = "Done"
)

//KubernetesStatus
const (
	KubernetesStatusInitializing = "initialize"
	KubernetesStatusFailed       = "failed"
	KubernetesStatusRunning      = "running"
)

//Kubernetes Labels
const (
	KubernetesRoleKey    = "node-role.kubernetes.io/master"
	KubernetesRoleMaster = "master"
	KubernetesRoleWorker = "worker"

	ClusterTypeKey   = "galaxy-future.org/app-type"
	ClusterTypeValue = "gf-cluster"

	ClusterInstanceGroupKey   = "galaxy-future.org/group"
	ClusterInstanceGroupIdKey = "galaxy-future.org/group-id"

	ClusterInstanceTypeKey          = "galaxy-future.org/machine-type"
	ClusterInstanceProviderLabelKey = "galaxy-future.org/machine-provider"
	ClusterInstanceClusterLabelKey  = "galaxy-future.org/bridgx-group"
)

const (
	KubernetesStandalone        string = "standalone"
	KubernetesHA                string = "HA"
	KubernetesHAMinMachineCount        = 4
)

const (
	HeaderTokenName = "Trans-UserToken"
)

const (
	InstanceInit   = "INIT"
	InstanceNormal = "NORMAL"
	InstanceError  = "ERROR"
)

const (
	OptTypeExpand = "EXPAND"
	OptTypeShrink = "SHRINK"
)

const (
	ExpandAndShrinkDefaultUser   = "bridgx"
	ExpandAndShrinkDefaultUserId = 0
)
