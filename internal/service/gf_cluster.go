package service

import (
	"context"
	"fmt"

	"github.com/galaxy-future/BridgX/cmd/api/middleware/authorization"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/types"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"github.com/gin-gonic/gin"
)

//GetBridgxUnusedCluster 获取所有没有被占用的集群
func GetBridgxUnusedCluster(ctx context.Context, user *authorization.CustomClaims, pageSize, pageNum int) ([]*gf_cluster.BridgxUnusedCluster, int, error) {
	tags := map[string]string{gf_cluster.UsageKey: gf_cluster.UnusedValue}
	clusters, total, err := GetClustersByTags(ctx, tags, pageSize, pageNum)
	if err != nil {
		return nil, 0, err
	}
	var clustersList []*gf_cluster.BridgxUnusedCluster
	for _, cluster := range clusters {
		targetCluster := &gf_cluster.BridgxUnusedCluster{
			ClusterName: cluster.ClusterName,
			CloudType:   cluster.Provider,
			Nodes:       nil,
		}
		accountKeys, err := GetAksByOrgAk(ctx, user.OrgId, "")
		if err != nil {
			return nil, 0, err
		}
		_, instances, total, err := GetInstancesByAccounts(ctx, GetInstancesCond{
			AccountKeys: accountKeys,
			Status:      []string{string(constants.Running)},
			ClusterName: cluster.ClusterName,
			PageNum:     1,
			PageSize:    10,
		})
		if err != nil {
			return nil, 0, err
		}
		for _, instance := range instances {
			targetCluster.Nodes = append(targetCluster.Nodes, instance.IpInner)
		}
		targetCluster.Total = int(total)
		clustersList = append(clustersList, targetCluster)
	}

	return clustersList, int(total), nil
}

//TagBridgxClusterUsage 将集群设置为被系统占用
func TagBridgxClusterUsage(clusterName, systemName string) error {
	tag := model.ClusterTag{
		ClusterName: clusterName,
		TagKey:      gf_cluster.UsageKey,
		TagValue:    systemName,
	}
	return EditClusterTags([]model.ClusterTag{tag})
}

//GetClusterInfo 根据集群名称获取集群名称
func GetClusterInfo(ctx context.Context, clusterName string) (*types.ClusterInfo, error) {
	cm, err := GetClusterByName(ctx, clusterName)
	if err != nil {
		return nil, err
	}
	tags, err := GetClusterTagsByClusterName(ctx, cm.ClusterName)
	if err != nil {
		return nil, err
	}
	return ConvertToClusterInfo(cm, tags)
}

//GetAllInstanceInCluster 获取集群中所有节点实例
func GetAllInstanceInCluster(ctx context.Context, user *authorization.CustomClaims, clusterName string) ([]model.Instance, error) {
	pageNum := 1
	pageSize := 50
	instances, total, err := getClusterInstances(ctx, user, clusterName, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	pageNum++

	for len(instances) < total {
		var tmpInstances []model.Instance
		tmpInstances, total, err = getClusterInstances(ctx, user, clusterName, pageNum, pageSize)
		if err != nil {
			return nil, err
		}
		instances = append(instances, tmpInstances...)
		pageNum++
	}
	return instances, nil
}

//getClusterInstances 根据peage获取实例列表
func getClusterInstances(ctx context.Context, user *authorization.CustomClaims, clusterName string, pageNum, pageSize int) ([]model.Instance, int, error) {
	accountKeys, err := GetAksByOrgAk(ctx, user.OrgId, "")
	if err != nil {
		return nil, 0, err
	}
	_, instances, total, err := GetInstancesByAccounts(ctx, GetInstancesCond{
		AccountKeys: accountKeys,
		Status:      []string{string(constants.Running)},
		ClusterName: clusterName,
		PageNum:     pageNum,
		PageSize:    pageSize,
	})
	if err != nil {
		return nil, 0, err
	}
	return instances, int(total), nil

}

//GetClusterAccount 根据集群获取Account信息
func GetClusterAccount(ctx *gin.Context, clusterName string) (*model.Account, error) {
	cluster, err := GetClusterByName(ctx, clusterName)
	if err != nil {
		return nil, err
	}
	account, err := GetAccount(cluster.Provider, cluster.AccountKey)
	if account == nil {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, err
	}
	return account, nil
}
