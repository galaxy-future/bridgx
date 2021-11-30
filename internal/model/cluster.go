package model

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/constants"
)

type Cluster struct {
	Base
	ClusterName  string //uniq_key
	ClusterDesc  string
	ExpectCount  int
	Status       string //ENABLE, DISABLE
	RegionId     string
	ZoneId       string
	InstanceType string

	ChargeType string
	Image      string
	Provider   string
	Password   string

	//Advanced Config
	NetworkConfig string
	StorageConfig string
	ChargeConfig  string
	AccountKey    string

	CreateBy      string
	UpdateBy      string
	DeleteUniqKey int64
	DeletedAt     gorm.DeletedAt
}

type ClusterSnapshot struct {
	Cluster         Cluster
	ActiveInstances []Instance
	RunningTask     []Task
}

func (Cluster) TableName() string {
	return "cluster"
}

// GetByClusterName find first record that match given conditions
func GetByClusterName(clusterName string) (*Cluster, error) {
	var out Cluster
	if err := clients.ReadDBCli.Where("cluster_name = ?", clusterName).First(&out).Error; err != nil {
		logErr("GetByClusterName from read db", err)
		return nil, err
	}
	return &out, nil
}

// GetByClusterNames find first record that match given conditions
func GetByClusterNames(clusterNames []string) ([]Cluster, error) {
	out := make([]Cluster, 0)
	if err := clients.ReadDBCli.Where("cluster_name IN (?)", clusterNames).Find(&out).Error; err != nil {
		logErr("GetByClusterName from read db", err)
		return nil, err
	}
	return out, nil
}

//GetClusterById find cluster with given cluster id
func GetClusterById(id int64) (*Cluster, error) {
	var cluster Cluster
	if err := clients.ReadDBCli.Where("id = ?", id).First(&cluster).Error; err != nil {
		logErr("GetByClusterName from read db", err)
		return nil, err
	}
	return &cluster, nil
}

//GetOneRegionByAccKey find one region_id with given accountKey
func GetOneRegionByAccKey(accountKey string) (*Cluster, error) {
	var cluster Cluster
	if err := clients.ReadDBCli.Where("account_key = ?", accountKey).First(&cluster).Error; err != nil {
		logErr("GetOneRegionByAccKey from read db", err)
		return nil, err
	}
	return &cluster, nil
}

//GetUpdatedCluster 获取任务更新时间大于指定时间的所有cluster实例
func GetUpdatedCluster(currentTime time.Time) ([]Cluster, error) {
	var clusters []Cluster
	if err := clients.ReadDBCli.Where("update_at >=  ", currentTime).First(&clusters).Error; err != nil {
		logErr("GetUpdatedCluster from read db", err)
		return nil, err
	}
	return clusters, nil
}

//GetClusterSnapshot 获取集群现状快照
func GetClusterSnapshot(clusterName string) (*ClusterSnapshot, error) {
	cluster, err := GetByClusterName(clusterName)
	if err != nil {
		logErr("GetClusterById from read db", err)
		return nil, err
	}
	instances, err := GetActiveInstancesByClusterName(cluster.ClusterName)
	if err != nil {
		logErr("GetActiveInstancesByClusterName from read db", err)
		return nil, err
	}
	tasks, err := GetTaskByStatus(cluster.ClusterName, []string{constants.TaskStatusInit, constants.TaskStatusRunning})
	if err != nil {
		logErr("GetActiveTaskByClusterName from read db", err)
		return nil, err
	}
	return &ClusterSnapshot{
		Cluster:         *cluster,
		ActiveInstances: instances,
		RunningTask:     tasks,
	}, nil
}

func ListClustersByCond(ctx context.Context, accountKeys []string, clusterName, provider string, pageNum, pageSize int) ([]Cluster, int, error) {
	res := make([]Cluster, 0)
	sql := clients.ReadDBCli.WithContext(ctx).Where(map[string]interface{}{})
	if len(accountKeys) > 0 {
		sql.Where("account_key IN (?)", accountKeys)
	}
	if provider != "" {
		sql.Where("provider = ?", provider)
	}
	if clusterName != "" {
		sql.Where("cluster_name LIKE ?", fmt.Sprintf("%%%v%%", clusterName))
	}
	err := sql.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&res).Error
	if err != nil {
		return res, 0, err
	}
	var cnt int64
	err = sql.Offset(-1).Limit(-1).Count(&cnt).Error
	if err != nil {
		return res, 0, err
	}
	return res, int(cnt), err
}
