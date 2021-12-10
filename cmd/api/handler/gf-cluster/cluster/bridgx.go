package cluster

import (
	"net/http"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/internal/clients"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"github.com/gin-gonic/gin"
)

//HandleListUnusedBridgxCluster lie列出没有被使用的集群
//TODO 使用Brdigx方法
func HandleListUnusedBridgxCluster(c *gin.Context) {
	claims := helper.GetUserClaims(c)
	if claims == nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("校验身份出错"))
		return
	}
	token, err := helper.GetUserToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	pageNumber, pageSize := helper.GetPagerParamFromQuery(c)

	//从bridgx中获取未被使用的集群
	//TODO 使用原生接口
	response, err := clients.GetClient().GetUnusedCluster(token, pageNumber, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	if response.Code != http.StatusOK {
		c.JSON(response.Code, gf_cluster.NewFailedResponse(response.Msg))
		return
	}

	var clustersList []*gf_cluster.BridgxUnusedCluster
	for _, cluster := range response.Data.ClusterList {
		targetCluster := &gf_cluster.BridgxUnusedCluster{
			ClusterName: cluster.ClusterName,
			CloudType:   cluster.Provider,
			Nodes:       nil,
		}

		//获取集群实例列表
		//TODO 使用原生接口
		instanceResponse, err := clients.GetClient().GetBridgxClusterInstances(token, cluster.ClusterName, pageNumber, pageSize)
		if err != nil {
			c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(err.Error()))
			return
		}
		if instanceResponse.Code != http.StatusOK {
			c.JSON(instanceResponse.Code, gf_cluster.NewFailedResponse(instanceResponse.Msg))
			return
		}

		for _, instance := range instanceResponse.Data.InstanceList {
			targetCluster.Nodes = append(targetCluster.Nodes, instance.IpInner)
		}

		clustersList = append(clustersList, targetCluster)
	}

	c.JSON(http.StatusOK, gf_cluster.NewListUnusedBridgxClusterResponse(clustersList, gf_cluster.Pager{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      response.Data.Pager.Total,
	}))
}
