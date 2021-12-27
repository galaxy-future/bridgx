package cloud

const (
	PrivateCloud = "PrivateCloud" //自建机房 私有云

	AlibabaCloud = "AlibabaCloud"
	HuaweiCloud  = "HuaweiCloud"
	TencentCloud = "TencentCloud"
)

const (
	TaskId      = "TaskId"
	ClusterName = "ClusterName"
)

const (
	PrePaid  = "PrePaid"
	PostPaid = "PostPaid"
)

const (
	Paid = iota + 1
	Unpaid
	Cancelled
)

const (
	EcsBuilding = "Pending"
	EcsRunning  = "Running"
	EcsStarting = "Starting"
	EcsStopping = "Stopping"
	EcsStopped  = "Stopped"
	EcsAbnormal = "Abnormal"
	EcsDeleted  = "Deleted"
)

const (
	BandwidthPayByTraffic = "PayByTraffic"
	BandwidthPayByFix     = "PayByBandwidth"
)

const (
	InsTypeAvailable = "Available"
	InsTypeAvaSoon   = "AvailableSoon"
	InsTypeLowStock  = "LowStock"
	InsTypeSellOut   = "Sellout"
)

const (
	ImageGlobal  = "global"
	ImagePrivate = "private"
	ImageShared  = "shared"
)

const (
	SecGroupRuleIn  = "ingress"
	SecGroupRuleOut = "egress"
)

const (
	SecGroupAllow = "allow"
	SecGroupDeny  = "deny"
)

const (
	IpV4 = "IPv4"
	IpV6 = "IPv6"
)

const (
	ProtocolIcmp   = "icmp"
	ProtocolIcmpV6 = "icmpV6"
	ProtocolTcp    = "tcp"
	ProtocolUdp    = "udp"
	ProtocolGre    = "gre"
	ProtocolAll    = "all"
)

const (
	OsLinux   = "linux"
	OsWindows = "windows"
	OsOther   = "other"
)

const (
	VPCStatusPending   = "Pending"
	VPCStatusAvailable = "Available"
	VPCStatusAbnormal  = "Abnormal"
)

const (
	SubnetPending   = "Pending"
	SubnetAvailable = "Available"
	SubnetAbnormal  = "Abnormal"
)
