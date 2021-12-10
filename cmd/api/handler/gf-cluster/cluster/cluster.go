package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/gf-cluster/cluster"
	cluster_builder "github.com/galaxy-future/BridgX/internal/gf-cluster/cluster-builder"
	"github.com/galaxy-future/BridgX/internal/gf-cluster/instance"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/pkg/encrypt"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//HandleCreateCluster 创建集群
func HandleCreateCluster(c *gin.Context) {
	//读取请求体
	//TODO 封装统一方法读取请求体
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("无效的请求体"))
		return
	}
	var buildRequest gf_cluster.BridgxClusterBuildRequest
	err = json.Unmarshal(data, &buildRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(fmt.Sprintf("无效的请求体, err : %s", err.Error())))
		return
	}

	//参数校验
	//1. 集群名称不能为空
	if buildRequest.BridgxClusterName == "" {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("集群名称不能为空"))
		return
	}

	//2. 用户鉴权
	claims := helper.GetUserClaims(c)
	if claims == nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("校验身份出错"))
		return
	}

	token, err := helper.GetUserToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	//3. 获取所选 Bridgx集群信息
	clusterResponse, err := clients.GetClient().GetBriodgxClusterDetails(token, buildRequest.BridgxClusterName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(fmt.Sprintf("获取集群信息时失败,错误信息： %s", err.Error())))
		return
	}
	if clusterResponse.Code != http.StatusOK {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(fmt.Sprintf("获取集群信息时失败,错误信息： %s", clusterResponse.Msg)))
		return
	}

	//4. 获取Bridgx集群实例信息
	instanceResponse, err := clients.GetClient().GetBridgxClusterAllInstances(token, buildRequest.BridgxClusterName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(fmt.Sprintf("获取集群实例时失败,错误信息： %s", err.Error())))
		return
	}

	//5. 获取AKSK信息
	akskResponse, err := clients.GetClient().GetAKSKClusterDetails(token, buildRequest.BridgxClusterName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(fmt.Sprintf("获取集群信息认证时失败,错误信息： %s", err.Error())))
		return
	}
	if akskResponse.Code != http.StatusOK {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(fmt.Sprintf("获取集群信息认证时失败,错误信息： %s", clusterResponse.Msg)))
		return
	}
	descryptRes, err := encrypt.AESDecrypt(akskResponse.Data.AccountKey+"bridgx", akskResponse.Data.AccountSecretEncrypt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(fmt.Sprintf("解密集群信息认证时失败,错误信息： %s", err.Error())))
		return
	}

	//6. 集群搭建策略
	if buildRequest.ClusterType == gf_cluster.KubernetesHA {
		if len(instanceResponse) < gf_cluster.KubernetesHAMinMachineCount {
			c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(fmt.Sprintf("高可用集群要求最少%d台物理机", gf_cluster.KubernetesHAMinMachineCount)))
			return
		}
	}

	//7. 机器相关校验
	if len(instanceResponse) == 0 {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("集群机器数量为0"))
		return
	}
	if buildRequest.ClusterType != gf_cluster.KubernetesHA && buildRequest.ClusterType != gf_cluster.KubernetesStandalone {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(fmt.Sprintf("不支持的集群类型: %s, 目前只支持[%s,%s]", buildRequest.ClusterType, gf_cluster.KubernetesStandalone, gf_cluster.KubernetesHA)))
		return
	}

	//8. 注册集群信息
	clusterRecord := gf_cluster.KubernetesInfo{
		Id:                0,
		Name:              buildRequest.ClusterName,
		Region:            clusterResponse.Data.RegionId,
		CloudType:         clusterResponse.Data.Provider,
		Status:            gf_cluster.KubernetesStatusInitializing,
		Config:            "",
		BridgxClusterName: buildRequest.BridgxClusterName,
		Type:              buildRequest.ClusterType,
		CreatedUser:       claims.Name,
		CreatedTime:       time.Now().Unix(),
	}
	err = model.RegisterKubernetesCluster(&clusterRecord)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(fmt.Sprintf("注册集群时失败，失败信息: %s", err.Error())))
		return
	}

	//9. Bridgx 占用集群
	useResponse, err := clients.GetClient().UpdateBridgxClusterUsingTag(token, buildRequest.BridgxClusterName, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(fmt.Sprintf("请求占用集群时出错，失败信息: %s", err.Error())))
		return
	}
	if useResponse.Code != http.StatusOK {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(fmt.Sprintf("请求占用集群时出错，失败信息: %s", useResponse.Msg)))
		return
	}

	//10 搭建集群
	buildParams := gf_cluster.ClusterBuilderParams{
		PodCidr:      buildRequest.PodCidr,
		SvcCidr:      buildRequest.ServiceCidr,
		MachineList:  nil,
		Mode:         gf_cluster.String2BuildMode(buildRequest.ClusterType),
		KubernetesId: clusterRecord.Id,
		AccessKey:    akskResponse.Data.AccountKey,
		AccessSecret: descryptRes,
	}

	for _, theInstance := range instanceResponse {
		buildParams.MachineList = append(buildParams.MachineList, gf_cluster.ClusterBuildMachine{
			IP:       theInstance.IpInner,
			Hostname: theInstance.InstanceId,
			Username: "root",
			Password: clusterResponse.Data.Password,
			Labels:   map[string]string{gf_cluster.ClusterInstanceTypeKey: theInstance.InstanceType, gf_cluster.ClusterInstanceProviderLabelKey: theInstance.Provider, gf_cluster.ClusterInstanceClusterLabelKey: clusterResponse.Data.Name},
		})
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Logger.Errorf("HandleCreateCluster err:%v ", r)
				logs.Logger.Errorf("HandleCreateCluster panic", zap.String("stack", string(debug.Stack())))
			}
		}()
		cluster_builder.CreateCluster(buildParams)
	}()

	c.JSON(http.StatusOK, gf_cluster.NewSuccessResponse())

}

//HandleListClusterSummary 列出所有集群
func HandleListClusterSummary(c *gin.Context) {
	id, _ := c.GetQuery("id")
	name, _ := c.GetQuery("name")
	pageNumber, pageSize := helper.GetPagerParamFromQuery(c)
	clusters, total, err := cluster.ListClustersSummary(id, name, pageNumber, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gf_cluster.NewListClusterSummaryResponse(clusters, gf_cluster.Pager{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      total,
	}))
}

//HandleGetClusterSummary 获得集群概述信息
func HandleGetClusterSummary(c *gin.Context) {
	clusterId, err := strconv.ParseInt(c.Param("clusterId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("should assign theCluster id"))
		return
	}
	theCluster, err := cluster.GetClustersSummary(clusterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gf_cluster.NewGetClusterSummaryResponse(theCluster))
}

//HandleDeleteKubernetes 删除集群
func HandleDeleteKubernetes(c *gin.Context) {
	clusterId, err := strconv.ParseInt(c.Param("clusterId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("未指定集群Id"))
		return
	}
	//获得用户token，与bridgx交互
	token, err := helper.GetUserToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	//获取当前集群信息
	theCluster, err := model.GetKubernetesCluster(clusterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	// 删除group
	groups, err := model.ListInstanceGroupInKubernetes(clusterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	for _, group := range groups {
		err := instance.DeleteInstanceGroup(group)
		if err != nil {
			logs.Logger.Errorw("failed to delete theCluster", zap.Error(err))
		}
	}

	//释放Bridgx集群
	useResponse, err := clients.GetClient().UpdateBridgxClusterUsingTag(token, theCluster.BridgxClusterName, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(fmt.Sprintf("释放集群时出错，失败信息: %s", err.Error())))
		return
	}
	if useResponse.Code != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(fmt.Sprintf("释放集群时出错，失败信息: %s", useResponse.Msg)))
		return
	}

	//删除集群
	err = model.DeleteKubernetesCluster(clusterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gf_cluster.NewSuccessResponse())
}

//HandleListNodesSummary 获取集群节点概要信息
func HandleListNodesSummary(c *gin.Context) {
	nodeIp := c.Query("node_ip")
	clusterName := c.Query("cluster_name")
	role := c.Query("role")

	clusterId, err := strconv.ParseInt(c.Param("clusterId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("should assign cluster id"))
		return
	}
	pageNumber, pageSize := helper.GetPagerParamFromQuery(c)

	nodes, err := cluster.ListClusterNodeSummary(clusterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	var result gf_cluster.ClusterNodeSummaryArray

	for _, node := range nodes {
		if nodeIp != "" && strings.Index(node.IpAddress, nodeIp) != 0 {
			continue
		}
		if clusterName != "" && strings.Index(node.ClusterName, clusterName) != 0 {
			continue
		}
		if role != "" && strings.Index(node.Role, role) != 0 {
			continue
		}
		result = append(result, node)
	}

	sort.Sort(result)

	start := (pageNumber - 1) * pageSize
	if start >= len(result) {
		c.JSON(http.StatusOK, gf_cluster.NewListClusterNodesResponse(nil, gf_cluster.Pager{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			Total:      len(result),
		}))
		return
	}

	end := pageNumber * pageSize
	if end >= len(result) {
		end = len(result)
	}
	c.JSON(http.StatusOK, gf_cluster.NewListClusterNodesResponse(result[start:end], gf_cluster.Pager{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      len(result),
	}))

}

//HandleListClusterPodsSummary 获取集群pod概述信息
func HandleListClusterPodsSummary(c *gin.Context) {
	nodeIp := c.Query("node_ip")
	podIp := c.Query("pod_ip")

	clusterId, err := strconv.ParseInt(c.Param("clusterId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("should assign cluster id"))
		return
	}
	pageNumber, pageSize := helper.GetPagerParamFromQuery(c)

	pods, err := cluster.ListClusterPodsSummary(clusterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	var result gf_cluster.ClusterPodsSummaryArray

	for _, pod := range pods {
		if nodeIp != "" && strings.Index(pod.NodeIp, nodeIp) != 0 {
			continue
		}
		if podIp != "" && strings.Index(pod.PodIP, podIp) != 0 {
			continue
		}
		result = append(result, pod)
	}

	sort.Sort(result)

	start := (pageNumber - 1) * pageSize
	if start >= len(result) {
		c.JSON(http.StatusOK, gf_cluster.NewListClusterPodsDetailResponse(nil, gf_cluster.Pager{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			Total:      len(result),
		}))
		return
	}

	end := pageNumber * pageSize
	if end >= len(result) {
		end = len(result)
	}
	c.JSON(http.StatusOK, gf_cluster.NewListClusterPodsDetailResponse(result[start:end], gf_cluster.Pager{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      len(result),
	}))

}
