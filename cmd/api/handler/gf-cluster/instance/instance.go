package instance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/internal/gf-cluster/cluster"
	"github.com/galaxy-future/BridgX/internal/gf-cluster/instance"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"github.com/gin-gonic/gin"
	"github.com/wxnacy/wgo/arrays"
	"go.uber.org/zap"
)

//HandleRestartInstance  c重启实例
func HandleRestartInstance(c *gin.Context) {

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("无效的请求体"))
		return
	}
	var request gf_cluster.InstanceRestartRequest
	err = json.Unmarshal(data, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(fmt.Sprintf("无效的请求体, err : %s", err.Error())))
		return
	}
	claims := helper.GetUserClaims(c)
	if claims == nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("校验身份出错"))
		return
	}
	//重启节点
	err = instance.RestartInstance(request.InstanceGroupId, request.InstanceName)
	if err != nil {
		logs.Logger.Errorw("failed to restart instance.", zap.Int64("instance_group_id", request.InstanceGroupId), zap.String("instance_name", request.InstanceName), zap.String("operator", claims.Name), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gf_cluster.NewSuccessResponse())
}

//HandleDeleteInstance 删除节点
func HandleDeleteInstance(c *gin.Context) {
	begin := time.Now()

	//获取用户信息
	claims := helper.GetUserClaims(c)
	if claims == nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("校验身份出错"))
		return
	}
	createdUserId := claims.UserId
	createdUserName := claims.Name

	//读取请求体
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("无效的请求体"))
		return
	}
	var request gf_cluster.InstanceDeleteRequest
	err = json.Unmarshal(data, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(fmt.Sprintf("无效的请求体, err : %s", err.Error())))
		return
	}

	instanceGroup, err := instance.GetInstanceGroup(request.InstanceGroupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	err = instance.DeleteInstance(instanceGroup, request.InstanceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	//增加操作日志
	cost := time.Now().Sub(begin).Milliseconds()
	err = instance.AddInstanceForm(instanceGroup, cost, createdUserId, createdUserName, gf_cluster.OptTypeShrink, 1, err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	if err == nil {
		c.JSON(http.StatusOK, gf_cluster.NewSuccessResponse())
	}

}

//HandleListInstance 列出所欲实例
func HandleListInstance(c *gin.Context) {
	instanceGroupId, err := strconv.ParseInt(c.Param("instanceGroup"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("未指定实例组id"))
		return
	}
	items, err := instance.ListCustomInstances(instanceGroupId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(err.Error()))
	}
	c.JSON(http.StatusOK, gf_cluster.NewInstanceListResponse(items))
}

//HandleListMyInstance 列出我的实例
func HandleListMyInstance(c *gin.Context) {
	nodeIp := c.Query("node_ip")
	podIp := c.Query("pod_ip")
	claims := helper.GetUserClaims(c)
	if claims == nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("校验身份出错"))
		return
	}
	curUserId := fmt.Sprintf("%v", claims.UserId)
	groups, err := model.ListInstanceGroupByUser(curUserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	pageNumber, pageSize := helper.GetPagerParamFromQuery(c)
	var result gf_cluster.ClusterPodsSummaryArray
	kubernetesMap := make(map[int64][]string)
	for _, group := range groups {
		kubernetesMap[group.KubernetesId] = append(kubernetesMap[group.KubernetesId], group.Name)
	}
	for kubernetesId, groupNames := range kubernetesMap {
		pods, err := cluster.ListClusterPodsSummary(kubernetesId)
		if err != nil {
			logs.Logger.Errorw("failed to list pods from kubernetes cluster.", zap.Int64("kubernetes_id", kubernetesId), zap.Error(err))
			continue
		}
		for _, pod := range pods {
			if nodeIp != "" && strings.Index(pod.NodeIp, nodeIp) != 0 {
				continue
			}
			if podIp != "" && strings.Index(pod.PodIP, podIp) != 0 {
				continue
			}
			if arrays.ContainsString(groupNames, pod.GroupName) == -1 {
				continue
			}
			result = append(result, pod)
		}
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

//HandleListInstanceForm 列出所有集群
func HandleListInstanceForm(c *gin.Context) {
	id, _ := c.GetQuery("id")
	pagerNumber, pageSize := helper.GetPagerParamFromQuery(c)
	forms, total, err := model.ListInstanceFormFromDB(id, pagerNumber, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gf_cluster.NewInstanceFormListResponse(forms, gf_cluster.Pager{
		PageNumber: pagerNumber,
		PageSize:   pageSize,
		Total:      int(total),
	}))
}
