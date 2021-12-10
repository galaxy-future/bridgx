package model

import (
	"context"
	"time"

	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/constants"
	"gorm.io/gorm"
)

const (
	InstanceTypeStatusNoActivate = iota
	InstanceTypeStatusActivated
	InstanceTypeStatusExpired
)

type Instance struct {
	Id           int64 `gorm:"primary_key"`
	Status       constants.Status
	IpInner      string
	IpOuter      string
	InstanceId   string
	ClusterName  string
	TaskId       int64 //扩容任务ID
	ShrinkTaskId int64 //缩容任务ID
	ChargeType   string
	CreateAt     *time.Time
	DeleteAt     *time.Time
	RunningAt    *time.Time
	ExpireAt     *time.Time //PrePaid instance expire time
}

func (Instance) TableName() string {
	return "instance"
}

type InstanceType struct {
	Base
	Provider string
	RegionId string
	ZoneId   string
	TypeName string
	Family   string
	Core     int // 核心数量,单位 核
	Memory   int // 内存大小,单位 G
	IStatus  int `gorm':"column:i_status"` // 0 未激活 1 已激活 2 已过期
}

func (InstanceType) TableName() string {
	return "instance_type"
}

func UpdateInstanceTypeIStatus(ctx context.Context, tx *gorm.DB, status int) error {
	err := tx.WithContext(ctx).Table(InstanceType{}.TableName()).
		Where("i_status = ?", status-1).
		Update("i_status", status).Error
	return err
}

func DropExpiredInstanceType(ctx context.Context, tx *gorm.DB) error {
	err := tx.WithContext(ctx).
		Where("i_status = ?", InstanceTypeStatusExpired).
		Delete(InstanceType{}).Error
	return err
}

func BatchCreateInstance(instances []Instance) error {
	return BatchCreate(instances)
}

func UpdateByInstanceId(instance Instance) error {
	if err := clients.WriteDBCli.Where("instance_id = ?", instance.InstanceId).Updates(instance).Error; err != nil {
		logErr("UpdateByInstanceId from write db", err)
		return err
	}
	return nil
}

func BatchUpdateByInstanceIds(instanceIds []string, instance Instance) error {
	if err := clients.WriteDBCli.Where("instance_id IN (?)", instanceIds).Updates(instance).Error; err != nil {
		logErr("UpdateByInstanceId from write db", err)
		return err
	}
	return nil
}

func GetInstanceByIpInner(ipInner string) (Instance, error) {
	//fixme 内网IP可能重复
	instance := &Instance{}
	if err := clients.ReadDBCli.Where("ip_inner = ?", ipInner).First(instance).Error; err != nil {
		logErr("GetInstanceByIpInner from read db", err)
		return *instance, err
	}
	return *instance, nil
}

func GetInstancesByIPs(ipList []string, clusterName string) ([]Instance, error) {
	instances := &[]Instance{}
	if err := clients.ReadDBCli.Where("ip_inner IN (?) AND cluster_name = ?", ipList, clusterName).Find(instances).Error; err != nil {
		logErr("GetInstancesByIPs from read db", err)
		return *instances, err
	}
	return *instances, nil
}

//GetActiveInstancesByClusterName 获取当前cluster下状态不为deleted状态的所有节点
func GetActiveInstancesByClusterName(clusterName string) ([]Instance, error) {
	var instances []Instance
	if err := clients.ReadDBCli.Where("cluster_name = ? AND status != ? ", clusterName, constants.Deleted).Find(&instances).Error; err != nil {
		logErr("GetActiveInstancesByClusterName from read db", err)
		return instances, err
	}
	return instances, nil
}

//GetActiveInstancesWithCount 获取当前cluster下状态不为deleted状态的count个节点
func GetActiveInstancesWithCount(clusterName string, count int) ([]Instance, error) {
	var instances []Instance
	if err := clients.ReadDBCli.Where("cluster_name = ? AND status != ? ", clusterName, constants.Deleted).Limit(count).Find(&instances).Error; err != nil {
		logErr("GetActiveInstancesWithCount from read db", err)
		return instances, err
	}
	return instances, nil
}

//GetActiveInstancesByClusters 获取clusters下状态不为deleted状态的count个节点
func GetActiveInstancesByClusters(ctx context.Context, clusterName []string) ([]Instance, error) {
	var instances []Instance
	if err := clients.ReadDBCli.WithContext(ctx).Where("cluster_name IN (?) AND status != ? ", clusterName, constants.Deleted).Find(&instances).Error; err != nil {
		logErr("GetActiveInstancesByClusters from read db", err)
		return instances, err
	}
	return instances, nil
}

//GetDeletedInstancesByTime 获取clusters下状态为deleted状态的节点
func GetDeletedInstancesByTime(ctx context.Context, clusterName []string, createBefore, deleteAfter time.Time) ([]Instance, error) {
	//取在createBefore之前&&在deleteAfter之后的实例
	var instances []Instance
	if err := clients.ReadDBCli.WithContext(ctx).Where("cluster_name IN (?) AND status = ? AND create_at < ? AND delete_at >= ?", clusterName, constants.Deleted, createBefore, deleteAfter).Find(&instances).Error; err != nil {
		logErr("GetDeletedInstancesByTime from read db", err)
		return instances, err
	}
	return instances, nil
}

//GetUsageInstancesBySpecifyDay 获取clusters下 指定时间段内存活过或仍然存活的节点
func GetUsageInstancesBySpecifyDay(ctx context.Context, clusterName []string, createBefore, deleteAfter time.Time, pageNum, pageSize int) ([]Instance, int64, error) {
	//取在createBefore之前&&在deleteAfter之后的实例
	var instances []Instance
	whereSql := clients.ReadDBCli.Debug().WithContext(ctx).Where("cluster_name IN (?) AND create_at < ? AND (delete_at >= ? OR delete_at IS NULL)", clusterName, createBefore, deleteAfter)
	if err := whereSql.Order("id DESC").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&instances).Error; err != nil {
		logErr("GetUsageInstancesBySpecifyDay from read db", err)
		return instances, 0, err
	}
	var total int64
	if err := whereSql.Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return instances, 0, err
	}
	return instances, total, nil
}

//CountActiveInstancesByClusterName 获取clusters下状态不为deleted状态节点数量
func CountActiveInstancesByClusterName(ctx context.Context, clusterNames []string) (int64, error) {
	if len(clusterNames) == 0 {
		return 0, nil
	}
	var ret int64
	if err := clients.ReadDBCli.WithContext(ctx).Model(&Instance{}).Where("cluster_name IN (?) AND status != ? ", clusterNames, constants.Deleted).Count(&ret).Error; err != nil {
		logErr("CountActiveInstancesByClusterName from read db", err)
		return 0, err
	}
	return ret, nil
}

//GetInstanceByInstanceId 获取Instance
func GetInstanceByInstanceId(instanceId string) (*Instance, error) {
	var ret Instance
	if err := clients.ReadDBCli.Model(&ret).Where("instance_id IN (?) ", instanceId).Find(&ret).Error; err != nil {
		logErr("CountActiveInstancesByClusterName from read db", err)
		return nil, err
	}
	return &ret, nil
}

type InstanceTypeCondition struct {
	Provider string
	RegionId string
	ZoneId   string
	Core     int
	Memory   int
}

func ScanInstanceType(ctx context.Context) (ins []InstanceType, err error) {
	err = clients.ReadDBCli.WithContext(ctx).Table(InstanceType{}.TableName()).
		Where("i_status = ?", InstanceTypeStatusActivated).
		Find(&ins).Error
	return ins, err
}

func GetInstanceTypeByName(ctx context.Context, instanceTypeName string) (*InstanceType, error) {
	ins := InstanceType{}
	err := clients.ReadDBCli.WithContext(ctx).Table(InstanceType{}.TableName()).
		Where("type_name = ?", instanceTypeName).
		First(&ins).Error
	return &ins, err
}

type InstanceSearchCond struct {
	Ip           string
	InstanceId   string
	ClusterNames []string
	Status       []string
	PageNumber   int
	PageSize     int
}

func GetInstanceByCond(ctx context.Context, cond InstanceSearchCond) (ret []Instance, total int64, err error) {
	query := clients.ReadDBCli.WithContext(ctx).Table(Instance{}.TableName()).
		Where("cluster_name in (?)", cond.ClusterNames)
	if cond.Ip != "" {
		query = query.Where("ip_inner = ? OR ip_outer = ?", cond.Ip, cond.Ip)
	}
	if cond.InstanceId != "" {
		query = query.Where("instance_id = ?", cond.InstanceId)
	}
	if len(cond.Status) != 0 && cond.Status[0] != "" {
		query = query.Where("status in (?)", cond.Status)
	}
	total, err = QueryWhere(query, cond.PageNumber, cond.PageSize, &ret, "id", true)
	if err != nil {
		return nil, 0, err
	}
	return ret, total, nil
}
