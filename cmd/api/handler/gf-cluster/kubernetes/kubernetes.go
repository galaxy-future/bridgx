package kubernetes

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"github.com/gin-gonic/gin"
)

func HandleRegisterKubernetes(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"message": "read request body failed"})
		return
	}
	var theCluster gf_cluster.KubernetesInfo
	err = json.Unmarshal(data, &theCluster)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	err = model.RegisterKubernetesCluster(&theCluster)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	c.JSON(200, gf_cluster.NewSuccessResponse())
}
func HandleListKubernetes(c *gin.Context) {
	kubernetes, err := model.ListRunningKubernetesClusters()
	if err != nil {
		c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	c.JSON(200, gf_cluster.NewKubernetesInfoListResponse(kubernetes))
}

func HandleGetKubernetes(c *gin.Context) {
	clusterId, err := strconv.ParseInt(c.Param("cluster"), 10, 64)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse("未提供ClusterId"))
		return
	}
	kubernetes, err := model.GetKubernetesCluster(clusterId)
	if err != nil {
		c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	if kubernetes == nil {
		c.JSON(400, gf_cluster.NewFailedResponse("没有找到相关记录"))
		return
	}
	c.JSON(200, gf_cluster.NewKubernetesInfoGetResponse(kubernetes))
}

func HandleUpdateKubernetes(c *gin.Context) {

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse("无效的请求信息"))
		return
	}
	var cluster gf_cluster.KubernetesInfo
	err = json.Unmarshal(data, &cluster)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	err = model.UpdateKubernetesCluster(&cluster)
	if err != nil {
		c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	c.JSON(200, gf_cluster.NewSuccessResponse())
}
