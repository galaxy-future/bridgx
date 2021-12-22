package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/types"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

type Cluster struct {
	Base
	ClusterName  string //uniq_key
	ClusterType  string
	ClusterDesc  string
	ExpectCount  int
	Status       string //ENABLE, DISABLE
	RegionId     string
	ZoneId       string
	InstanceType string

	Image    string
	Provider string
	Password string

	//Advanced Config
	ImageConfig   string
	NetworkConfig string
	StorageConfig string
	ChargeConfig  string
	AccountKey    string

	CreateBy      string
	UpdateBy      string
	DeleteUniqKey int64
	DeletedAt     gorm.DeletedAt
}

func (c *Cluster) GetChargeType() string {
	conf, err := c.UnmarshalChargeConfig()
	if err != nil {
		return ""
	}
	return conf.ChargeType
}

func (c *Cluster) UnmarshalChargeConfig() (*types.ChargeConfig, error) {
	chargeConfig := types.ChargeConfig{}
	if c.ChargeConfig != "" {
		err := jsoniter.UnmarshalFromString(c.ChargeConfig, &chargeConfig)
		if err != nil {
			return nil, err
		}
	}
	return &chargeConfig, nil
}

type ClusterSnapshot struct {
	Cluster         Cluster
	ActiveInstances []Instance
	RunningTask     []Task
}

func (Cluster) TableName() string {
	return "cluster"
}

func CreateClusterWithTagsAndInstances(ctx context.Context, cluster *Cluster, tags []*ClusterTag, instances []Instance) error {
	return clients.WriteDBCli.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(cluster).Error; err != nil {
			return err
		}
		if len(tags) > 0 {
			if err := tx.Create(&tags).Error; err != nil {
				return err
			}
		}
		if len(instances) > 0 {
			if err := tx.Create(&instances).Error; err != nil {
				return err
			}
		}
		return nil
	})
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

type ClusterSearchCond struct {
	AccountKeys []string
	ClusterName string
	ClusterType string
	Provider    string
	Usage       string
	PageNum     int
	PageSize    int
}

func ListClustersByCond(ctx context.Context, cond ClusterSearchCond) ([]Cluster, int, error) {
	res := make([]Cluster, 0)
	sql := clients.ReadDBCli.Debug().WithContext(ctx).Where(map[string]interface{}{})
	if cond.ClusterType != constants.ClusterTypeCustom && len(cond.AccountKeys) > 0 {
		sql.Where("cluster.account_key IN (?)", cond.AccountKeys)
	}
	if cond.Provider != "" {
		sql.Where("cluster.provider = ?", cond.Provider)
	} else {
		sql.Where("cluster.provider != ?", cloud.PrivateCloud)
	}
	if cond.ClusterName != "" {
		sql.Where("cluster.cluster_name LIKE ?", fmt.Sprintf("%%%v%%", cond.ClusterName))
	}
	if cond.ClusterType != "" {
		sql.Where("cluster.cluster_type = ?", strings.ToLower(cond.ClusterType))
	}
	joins := ""
	if cond.Usage != "" {
		joins = "LEFT JOIN cluster_tag ct on ct.cluster_name = cluster.cluster_name"
		sql.Where("ct.tag_key = ? and ct.tag_value LIKE ?", constants.DefaultClusterUsageKey, fmt.Sprintf("%%%v%%", cond.Usage))
	}

	err := sql.Distinct().Order("id DESC").Offset((cond.PageNum - 1) * cond.PageSize).Limit(cond.PageSize).Joins(joins).Find(&res).Error
	if err != nil {
		return res, 0, err
	}
	var cnt int64
	err = sql.Distinct("cluster.cluster_name").Offset(-1).Limit(-1).Count(&cnt).Error
	if err != nil {
		return res, 0, err
	}
	return res, int(cnt), err
}
