package model

import (
	"context"
	"fmt"

	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/constants"
)

// ClusterTag
//使用Tags来描述Cluster的用途，属性等
// 比如使用
//	{
//		"group": "serviceName"
//  	"imageId": "xxx:1.2"
//  }
type ClusterTag struct {
	Base
	ClusterName string
	TagKey      string
	TagValue    string
}

func (ClusterTag) TableName() string {
	return "cluster_tag"
}

func GetTagsByClusterName(clusterName string) ([]ClusterTag, error) {
	clusterTags := make([]ClusterTag, 0)
	err := QueryAll(map[string]interface{}{"cluster_name": clusterName}, &clusterTags, "")
	if err != nil {
		logErr("GetTagsByClusterName from read db", err)
		return nil, err
	}
	return clusterTags, nil
}

func GetClusterTags(ctx context.Context, clusterName, tagKey string, pageNum, pageSize int) ([]ClusterTag, int64, error) {
	clusterTags := make([]ClusterTag, 0)
	cond := map[string]interface{}{"cluster_name": clusterName}
	if tagKey != "" {
		cond["tag_key"] = tagKey
	}
	total, err := Query(cond, pageNum, pageSize, &clusterTags, "", true)
	if err != nil {
		logErr("GetTagsByClusterName from read db", err)
		return nil, 0, err
	}
	return clusterTags, total, nil
}

func GetClusterTagsByClusterNames(ctx context.Context, clusterNames []string) ([]ClusterTag, error) {
	clusterTags := make([]ClusterTag, 0)
	err := QueryAll(map[string]interface{}{"cluster_name": clusterNames}, &clusterTags, "")
	if err != nil {
		logErr("GetClusterTagsByClusterNames from read db", err)
		return nil, err
	}
	return clusterTags, nil
}

func GetBySpecifyClusterTag(clusterName string, tagKey string) (*ClusterTag, error) {
	var ret ClusterTag
	err := clients.ReadDBCli.Where("cluster_name = ? AND tag_key = ?", clusterName, tagKey).First(&ret).Error
	if err != nil {
		logErr("GetBySpecifyClusterTag from read db", err)
		return nil, err
	}
	return &ret, nil
}

//GetClusterNamesByTags list cluster names(contains deleted cluster) by tags:
// e.g.
// {
//		"k1": "v1",
//		"k2": "v2",
//		"k3": "",
// }
// ===>
// should search distinct(cluster_name) from database using condition:
// (tag_key = 'k1' AND tag_value = 'v1') OR tag_key = 'k3' OR (tag_key = 'k2' AND tag_value = 'v2'))
//	group by cluster_name
//
func GetClusterNamesByTags(ctx context.Context, tags map[string]string, pageSize, pageNum int) ([]ClusterTag, int64, error) {
	where := clients.ReadDBCli.WithContext(ctx)
	validCondCount := 0
	for tagKey, tagValue := range tags {
		if tagKey == "" {
			continue
		}
		validCondCount++
		orCondition := clients.ReadDBCli.WithContext(ctx)
		if validCondCount > 1 {
			orCondition = orCondition.Where("tag_key = ?", tagKey)
			if tagValue != "" {
				orCondition = orCondition.Where("tag_value = ?", tagValue)
			}
			where = where.Or(orCondition)
			continue
		}
		orCondition = orCondition.Where("tag_key = ?", tagKey)
		if tagValue != "" {
			orCondition = orCondition.Where("tag_value = ?", tagValue)
		}
		where = where.Or(orCondition)
	}
	if validCondCount == 0 {
		return nil, 0, fmt.Errorf("tags invalid")
	}
	ret := make([]ClusterTag, 0)
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 || pageSize > constants.DefaultPageSize {
		pageSize = constants.DefaultPageSize
	}
	offset := (pageNum - 1) * pageSize
	query := clients.ReadDBCli.Select("cluster_name").Where(where).
		Group("cluster_name").Having("count(cluster_name) = ?", validCondCount).
		Offset(offset).Limit(pageSize).Find(&ret)

	if err := query.Error; err != nil {
		logErr("query data from read db", err)
		return nil, 0, err
	}
	var count int64
	if err := query.Offset(-1).Limit(-1).Distinct().Count(&count).Error; err != nil {
		logErr("query data from read db", err)
		return nil, 0, err
	}
	return ret, count, nil
}
