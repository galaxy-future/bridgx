package alibaba

import (
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

const (
	_subOrderNumPerMain    = 3
	_maxNumEcsPerOperation = 100
)

//in
var _imageType = map[string]string{
	cloud.ImageGlobal:  "system",
	cloud.ImageShared:  "others",
	cloud.ImagePrivate: "self",
}

var _protocol = map[string]string{
	cloud.ProtocolIcmp: "icmp",
	cloud.ProtocolTcp:  "tcp",
	cloud.ProtocolUdp:  "udp",
	cloud.ProtocolGre:  "gre",
	cloud.ProtocolAll:  "all",
}

//out
var _chargeType = map[string]string{
	"Subscription": cloud.PrePaid,
	"PayAsYouGo":   cloud.PostPaid,
}

var _payStatus = map[string]int8{
	"Paid":      cloud.Paid,
	"Unpaid":    cloud.Unpaid,
	"Cancelled": cloud.Cancelled,
}

var _ecsStatus = map[string]string{
	"Pending":  cloud.EcsBuilding,
	"Running":  cloud.EcsRunning,
	"Starting": cloud.EcsStarting,
	"Stopping": cloud.EcsStopping,
	"Stopped":  cloud.EcsStopped,
}

var _insTypeStat = map[string]string{
	"WithStock":          cloud.InsTypeAvailable,
	"ClosedWithStock":    cloud.InsTypeLowStock,
	"WithoutStock":       cloud.InsTypeAvaSoon,
	"ClosedWithoutStock": cloud.InsTypeSellOut,
}

var _bandwidthChargeType = map[string]string{
	"PayByBandwidth": cloud.BandwidthPayByFix,
	"PayByTraffic":   cloud.BandwidthPayByTraffic,
}

var _secGrpRuleDirection = map[string]string{
	"ingress": cloud.SecGroupRuleIn,
	"egress":  cloud.SecGroupRuleOut,
}

var _osType = map[string]string{
	"linux":   cloud.OsLinux,
	"windows": cloud.OsWindows,
}

var _vpcStatus = map[string]string{
	"Pending":   cloud.VPCStatusPending,
	"Available": cloud.VPCStatusAvailable,
}

var _subnetStatus = map[string]string{
	"Pending":   cloud.SubnetPending,
	"Available": cloud.SubnetAvailable,
}
