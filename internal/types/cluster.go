package types

import (
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

type ClusterInfo struct {
	//Base Config
	Id           int64  `json:"id"`
	Name         string `json:"name" binding:"required,max=20"`
	Desc         string `json:"desc"`
	RegionId     string `json:"region_id" binding:"required"`
	ZoneId       string `json:"zone_id" binding:"required"`
	InstanceType string `json:"instance_type" binding:"required"`
	Image        string `json:"image" binding:"required"`
	Provider     string `json:"provider" binding:"required,mustIn=cloud"`
	Username     string `json:"username"`
	Password     string `json:"password" binding:"required,min=8,max=30,charTypeGT3"`
	AccountKey   string `json:"account_key" binding:"required"` //阿里云ak

	//Advanced Config
	NetworkConfig *NetworkConfig `json:"network_config"`
	StorageConfig *StorageConfig `json:"storage_config"`
	ChargeConfig  *ChargeConfig  `json:"charge_config"`

	//Custom Config
	Tags map[string]string `json:"tags"`

	InstanceCore   int `json:"instance_core"`   // 核心数量,单位 核
	InstanceMemory int `json:"instance_memory"` // 内存大小,单位 G
}

type NetworkConfig struct {
	Vpc                     string `json:"vpc" binding:"required"`
	SubnetId                string `json:"subnet_id" binding:"required"`
	SecurityGroup           string `json:"security_group" binding:"required"`
	InternetChargeType      string `json:"internet_charge_type"`
	InternetMaxBandwidthOut int    `json:"internet_max_bandwidth_out"`
}

type StorageConfig struct {
	MountPoint string       `json:"mount_point"`
	NAS        string       `json:"nas"`
	Disks      *cloud.Disks `json:"disks"`
}

type ChargeConfig struct {
	ChargeType string `json:"charge_type"`
	Period     int    `json:"period"`
	PeriodUnit string `json:"period_unit"`
}

type OrgKeys struct {
	OrgId int64     `json:"org_id"`
	Info  []KeyInfo `json:"info"`
}

type KeyInfo struct {
	AK       string
	SK       string
	Provider string
}

type Pager struct {
	PageNumber int
	PageSize   int
	Total      int
}
