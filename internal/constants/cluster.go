package constants

const (
	DeletingIPs          = "deleting_ips"
	Instances            = "instances"
	WorkingIPs           = "working_ips"
	ExpectInstanceNumber = "expect_instance_number"

	HasNoneIP       = "-"
	HasNoneInstance = "-"

	Interval = 60
	Delay    = 5
	Retry    = 3

	BatchMax = 100

	DefaultUsername           = "root"
	DefaultClusterUsageKey    = "usage"
	DefaultClusterUsageUnused = "unused"
)

const (
	ClusterStatusEnable  = "ENABLE"
	ClusterStatusDisable = "DISABLE"
)

const (
	ErrClusterNotExist           = "集群: [%s] 不存在"
	ErrPrePaidShrinkNotSupported = "不支持对包年包月的集群机器进行缩容操作"
)
