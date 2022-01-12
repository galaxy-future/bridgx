package tencent

import (
	"errors"

	"github.com/galaxy-future/BridgX/pkg/cloud"
)

const (
	_maxNumEcsPerOperation = 100
	_offset                = 0
	_pageSize              = 100
)

const (
	_vpcEndpoint       = "vpc.tencentcloudapi.com"
	_cvmEndpoint       = "cvm.tencentcloudapi.com"
	_apiEndpoint       = "api.tencentcloudapi.com"
	_subnetFilterVpcId = "vpc-id"
)

var (
	_errResponseIsNil = errors.New("response is nil")
	_errIsNotOne      = errors.New("totalCount isn't one")
)

//in
var _inEcsChargeType = map[string]string{
	cloud.InstanceChargeTypePrePaid:  "PREPAID",
	cloud.InstanceChargeTypePostPaid: "POSTPAID_BY_HOUR",
}

var _imageType = map[string]string{
	cloud.ImageGlobal:  "PUBLIC_IMAGE",
	cloud.ImageShared:  "SHARED_IMAGE",
	cloud.ImagePrivate: "PRIVATE_IMAGE",
}

var _bandwidthChargeMode = map[string]string{
	cloud.BandwidthPayByTraffic: "TRAFFIC_POSTPAID_BY_HOUR",
	cloud.BandwidthPayByFix:     "BANDWIDTH_PREPAID",
}

var _protocol = map[string]string{
	cloud.ProtocolIcmp:   "ICMP",
	cloud.ProtocolIcmpV6: "ICMPV6",
	cloud.ProtocolTcp:    "TCP",
	cloud.ProtocolUdp:    "UDP",
	cloud.ProtocolGre:    "GRE",
	cloud.ProtocolAll:    "ALL",
}

//out
var _ecsChargeType = map[string]string{
	"POSTPAID_BY_HOUR": cloud.InstanceChargeTypePostPaid,
	"PREPAID":          cloud.InstanceChargeTypePrePaid,
}

var _ecsStatus = map[string]string{
	"PENDING":       cloud.EcsBuilding,
	"STARTING":      cloud.EcsStarting,
	"REBOOTING":     cloud.EcsStarting,
	"RUNNING":       cloud.EcsRunning,
	"STOPPING":      cloud.EcsStopping,
	"STOPPED":       cloud.EcsStopped,
	"LAUNCH_FAILED": cloud.EcsAbnormal,
	"SHUTDOWN":      cloud.EcsDeleted,
	"TERMINATING":   cloud.EcsDeleted,
}

var _insTypeChargeType = map[string]string{
	"PREPAID":          cloud.InsTypeChargeTypePrePaid,
	"POSTPAID_BY_HOUR": cloud.InsTypeChargeTypePostPaid,
}

var _insTypeStat = map[string]string{
	"SELL":     cloud.InsTypeAvailable,
	"SOLD_OUT": cloud.InsTypeSellOut,
}

var _bandwidthChargeType = map[string]string{
	"BANDWIDTH_PREPAID":        cloud.BandwidthPayByFix,
	"TRAFFIC_POSTPAID_BY_HOUR": cloud.BandwidthPayByTraffic,
}

var _outProtocol = map[string]string{
	"icmp":   cloud.ProtocolIcmp,
	"icmpv6": cloud.ProtocolIcmpV6,
	"tcp":    cloud.ProtocolTcp,
	"udp":    cloud.ProtocolUdp,
	"gre":    cloud.ProtocolGre,
	"ALL":    cloud.ProtocolAll,
}
