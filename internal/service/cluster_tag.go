package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/model"
)

func GetClusterTags(ctx context.Context, clusterName, tagKey string, pageNum, pageSize int) ([]model.ClusterTag, int64, error) {
	return model.GetClusterTags(ctx, clusterName, tagKey, pageNum, pageSize)
}

func CreateClusterTags(tags []model.ClusterTag) error {
	return model.Create(&tags)
}

func EditClusterTags(tags []model.ClusterTag) error {
	for _, tag := range tags {
		tagDB, err := model.GetBySpecifyClusterTag(tag.ClusterName, tag.TagKey)
		if strings.ToLower(tag.TagKey) == constants.DefaultClusterUsageKey {
			if strings.ToLower(tagDB.TagValue) != constants.DefaultClusterUsageUnused && strings.ToLower(tag.TagValue) != constants.DefaultClusterUsageUnused {
				return fmt.Errorf("tag[usage] unsupport state transition %v->%v", tagDB.TagValue, tag.TagValue)
			}
		}
		if err != nil {
			return err
		}
		tagDB.TagValue = tag.TagValue
		tagDB.ValueDesc = tag.ValueDesc
		err = model.Save(tagDB)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteClusterTags(tags []model.ClusterTag) error {
	return model.Delete(&tags)
}

func GetClustersByTags(ctx context.Context, tags map[string]string, pageSize, pageNum int) ([]model.Cluster, int64, error) {
	clusters := make([]model.Cluster, 0)
	clusterTags, total, err := model.GetClusterNamesByTags(ctx, tags, pageSize, pageNum)
	if err != nil || total == 0 || len(clusterTags) == 0 {
		return clusters, 0, err
	}
	clusterNames := make([]string, 0)
	for _, tag := range clusterTags {
		clusterNames = append(clusterNames, tag.ClusterName)
	}
	clusters, err = GetClustersByNames(ctx, clusterNames)
	if err != nil {
		return nil, 0, err
	}
	return clusters, total, err
}

func GetClusterTagsByClusters(ctx context.Context, clusters []model.Cluster) (map[string]map[string]string, error) {
	res := make(map[string]map[string]string, 0)
	if len(clusters) == 0 {
		return res, nil
	}
	clusterNames := make([]string, 0)
	for _, cluster := range clusters {
		clusterNames = append(clusterNames, cluster.ClusterName)
	}
	tags, err := model.GetClusterTagsByClusterNames(ctx, clusterNames)
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		clusterTags, ok := res[tag.ClusterName]
		if ok {
			clusterTags[tag.TagKey] = tag.TagValue
		} else {
			clusterTags = make(map[string]string, 0)
			clusterTags[tag.TagKey] = tag.TagValue
		}
		if tag.ValueDesc != "" {
			clusterTags[tag.TagKey] = tag.ValueDesc
		}
		res[tag.ClusterName] = clusterTags
	}
	return res, nil
}
