package kubernetes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"github.com/gin-gonic/gin"
)

//HandleRegisterKubernetes 注册集群，用于支持已有k8s集群注册
//后期用于其他集群直接录入
func HandleRegisterKubernetes(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "read request body failed"})
		return
	}
	var theCluster gf_cluster.KubernetesInfo
	err = json.Unmarshal(data, &theCluster)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	err = model.RegisterKubernetesCluster(&theCluster)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gf_cluster.NewSuccessResponse())
}

//HandleListKubernetes 列出所有集群
func HandleListKubernetes(c *gin.Context) {
	kubernetes, err := model.ListRunningKubernetesClusters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gf_cluster.NewKubernetesInfoListResponse(kubernetes))
}

//HandleGetKubernetes 获取指定集群详细信息
func HandleGetKubernetes(c *gin.Context) {
	clusterId, err := strconv.ParseInt(c.Param("cluster"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("未提供ClusterId"))
		return
	}
	kubernetes, err := model.GetKubernetesCluster(clusterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	if kubernetes == nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("没有找到相关记录"))
		return
	}
	c.JSON(http.StatusOK, gf_cluster.NewKubernetesInfoGetResponse(kubernetes))
}

//HandleUpdateKubernetes 更新集群信息
func HandleUpdateKubernetes(c *gin.Context) {

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse("无效的请求信息"))
		return
	}
	var cluster gf_cluster.KubernetesInfo
	err = json.Unmarshal(data, &cluster)
	if err != nil {
		c.JSON(http.StatusBadRequest, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	err = model.UpdateKubernetesCluster(&cluster)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gf_cluster.NewSuccessResponse())
}
