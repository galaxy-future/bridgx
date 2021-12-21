package constants

const (
	DeletingIPs          = "deleting_ips"
	Instances            = "instances"
	WorkingIPs           = "working_ips"
	ExpectInstanceNumber = "expect_instance_number"

	HasNoneIP       = "-"
	HasNoneInstance = "-"

	Interval = 15
	Retry    = 3

	BatchMax = 100

	DefaultUsername           = "root"
	DefaultClusterUsageKey    = "usage"
	DefaultClusterUsageUnused = "unused"

	ClusterStatusEnable  = "ENABLE"
	ClusterStatusDisable = "DISABLE"

	ClusterTypeStandard = "standard"
	ClusterTypeCustom   = "custom"
)

const (
	ErrClusterNotExist           = "集群: [%s] 不存在"
	ErrPrePaidShrinkNotSupported = "不支持对包年包月的集群机器进行缩容操作"
)

const (
	GPU                     = "GPU"
	CPU                     = "CPU"
	IsAlibabaCloudGpuType   = "gn"
	IsHuaweiCloudGpuType    = "G"
	IsHuaweiCloudGpuTypeTwo = "P"
)
