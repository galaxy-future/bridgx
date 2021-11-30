package cloud

import (
	"time"
)

const (
	Pending     = "Pending"
	TaskId      = "TaskId"
	ClusterName = "ClusterName"
)

type Params struct {
	Provider     string
	InstanceType string
	ImageId      string
	Network      *Network
	Zone         string
	Region       string
	Disks        *Disks
	Charge       *Charge
	Password     string
	Tags         []Tag
}

type Tag struct {
	Key   string
	Value string
}

type Instance struct {
	Id       string   `json:"id"`
	CostWay  string   `json:"cost_way"`
	Provider string   `json:"provider"`
	IpInner  string   `json:"ip_inner"`
	IpOuter  string   `json:"ip_outer"`
	Network  *Network `json:"network"`
	ImageId  string   `json:"image_id"`
	Status   string   `json:"status"`
}

type Network struct {
	VpcId                   string `json:"vpc_id"`
	SubnetId                string `json:"subnet_id"`
	SecurityGroup           string `json:"security_group"`
	InternetChargeType      string `json:"internet_charge_type"`
	InternetMaxBandwidthOut int    `json:"internet_max_bandwidth_out"`
}

type Charge struct {
	Period     int    `json:"period"`
	PeriodUnit string `json:"period_unit"`
	ChargeType string `json:"charge_type"`
}

type Disks struct {
	SystemDisk DiskConf   `json:"system_disk"`
	DataDisk   []DiskConf `json:"data_disk"`
}

type DiskConf struct {
	Category         string `json:"category"`
	Size             int    `json:"size"`
	PerformanceLevel string `json:"performance_level"`
}

type CreateVpcRequest struct {
	RegionId  string
	VpcName   string
	CidrBlock string
}

type CreateVpcResponse struct {
	//RouterId       *string `json:"VRouterId,omitempty" xml:"VRouterId,omitempty"`
	//RouteTableId    *string `json:"RouteTableId,omitempty" xml:"RouteTableId,omitempty"`
	VpcId     string
	RequestId string
	//ResourceGroupId *string `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`

}

type GetVpcRequest struct {
	VpcId    string
	RegionId string
	VpcName  string
}

type GetVpcResponse struct {
	Vpc VPC
}

type VPC struct {
	VpcId     string
	VpcName   string
	CidrBlock string
	SwitchIds []string
	RegionId  string
	Status    string
	CreateAt  string
}

type CreateSwitchRequest struct {
	RegionId    string
	ZoneId      string
	CidrBlock   string
	VSwitchName string
	VpcId       string
}
type CreateSwitchResponse struct {
	SwitchId  string
	RequestId string
}

type CreateSecurityGroupRequest struct {
	RegionId          string
	SecurityGroupName string
	VpcId             string
	SecurityGroupType string
}

type CreateSecurityGroupResponse struct {
	SecurityGroupId string
	RequestId       string
}

type AddSecurityGroupRuleRequest struct {
	RegionId        string
	VpcId           string
	SecurityGroupId string
	IpProtocol      string
	PortRange       string
	GroupId         string
	CidrIp          string
	PrefixListId    string
}

type DescribeSecurityGroupsRequest struct {
	VpcId    string
	RegionId string
}

type DescribeSecurityGroupsResponse struct {
	Groups []SecurityGroup
}

type SecurityGroup struct {
	SecurityGroupId   string
	SecurityGroupType string
	SecurityGroupName string
	CreateAt          string
	VpcId             string
	RegionId          string
}
type AddSecurityGroupRuleResponse struct {
}

type GetSwitchRequest struct {
	SwitchId string
}

type Switch struct {
	VpcId                   string
	SwitchId                string
	Name                    string
	IsDefault               int
	AvailableIpAddressCount int
	VStatus                 string
	CreateAt                string
	ZoneId                  string
	CidrBlock               string
}

type GetSwitchResponse struct {
	Switch Switch
}

type GetRegionsResponse struct {
	Regions []Region
}

type Region struct {
	RegionId  string
	LocalName string
}

type GetZonesRequest struct {
	RegionId string
}

type GetZonesResponse struct {
	Zones []Zone
}

type Zone struct {
	ZoneId    string
	LocalName string
}

type InstanceType struct {
	Status         string
	StatusCategory string
	Value          string
}

type InstanceInfo struct {
	Core        int
	Memory      int
	Family      string
	InsTypeName string
}

type DescribeAvailableResourceRequest struct {
	RegionId string
	ZoneId   string
}

type DescribeAvailableResourceResponse struct {
	InstanceTypes map[string][]InstanceType
}

type AvailableZone struct {
	ZoneId string
	Status string
}

type DescribeInstanceTypesRequest struct {
	TypeName []string
}
type DescribeInstanceTypesResponse struct {
	Infos []InstanceInfo
}
type DescribeImagesRequest struct {
	RegionId string
}

type DescribeImagesResponse struct {
	Images []Image
}

type Image struct {
	OsType  string
	OsName  string
	ImageId string
}

type DescribeVpcsRequest struct {
	RegionId string
}

type DescribeVpcsResponse struct {
	Vpcs []VPC
}

type DescribeSwitchesRequest struct {
	VpcId string
}

type DescribeSwitchesResponse struct {
	Switches []Switch
}

type SecurityGroupRule struct {
	VpcId           string
	SecurityGroupId string
	PortRange       string
	Protocol        string
	Direction       string
	GroupId         string
	CidrIp          string
	PrefixListId    string
	CreateAt        string
}

type DescribeGroupRulesRequest struct {
	RegionId        string
	SecurityGroupId string
}

type DescribeGroupRulesResponse struct {
	Rules []SecurityGroupRule
}

type GetOrdersRequest struct {
	StartTime time.Time
	EndTime   time.Time
	PageNum   int
	PageSize  int
}

type GetOrdersResponse struct {
	Orders []Order
}

type Order struct {
	OrderId        string
	OrderTime      time.Time
	Product        string
	Quantity       int32
	UsageStartTime time.Time
	UsageEndTime   time.Time
	RegionId       string
	ChargeType     string
	PayStatus      int8
	Currency       string
	Cost           float32
	Extend         map[string]interface{}
}
