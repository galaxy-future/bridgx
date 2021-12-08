package eci

import (
	"encoding/json"
	"fmt"
	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/internal/gf-cluster/instance"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"strconv"
	"time"
)

//HandleCreateInstanceGroup  创建实例组
func HandleCreateInstanceGroup(c *gin.Context) {
	begin := time.Now()

	//解析请求
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse("无效的请求格式"))
		return
	}
	var cluster gf_cluster.InstanceGroupCreateRequest
	err = json.Unmarshal(data, &cluster)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse(fmt.Sprintf("无效的请求格式, err : %s", err.Error())))
		return
	}

	claims := helper.GetUserClaims(c)
	if claims == nil {
		c.JSON(400, gf_cluster.NewFailedResponse("校验身份出错"))
		return
	}
	createdUserId := claims.UserId
	createdUserName := claims.Name
	instanceGroup := gf_cluster.InstanceGroup{
		Id:            0,
		KubernetesId:  cluster.KubernetesId,
		Name:          cluster.Name,
		Image:         cluster.Image,
		Cpu:           cluster.Cpu,
		Memory:        cluster.Memory,
		Disk:          cluster.Disk,
		InstanceCount: cluster.InstanceCount,
		CreatedUser:   createdUserName,
		CreatedUserId: createdUserId,
	}
	err = instance.CreateInstanceGroup(&instanceGroup)
	if err != nil {
		c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	//统计请求时间
	defer func() {
		cost := time.Now().Sub(begin).Milliseconds()
		err := instance.AddInstanceForm(&instanceGroup, cost, createdUserId, createdUserName, gf_cluster.OptTypeExpand, instanceGroup.InstanceCount, err)
		if err != nil {
			c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
			return
		}
		if err == nil {
			c.JSON(200, gf_cluster.NewSuccessResponse())
		}
	}()

	//扩容集群
	err = instance.ExpandCustomInstanceGroup(&instanceGroup, cluster.InstanceCount)
	if err != nil {
		c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
}

//HandleBatchCreateInstanceGroup 批量新建实例组
func HandleBatchCreateInstanceGroup(c *gin.Context) {

	//解析请求
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse("无效的请求格式"))
		return
	}
	var instanceGroups []gf_cluster.InstanceGroupCreateRequest
	err = json.Unmarshal(data, &instanceGroups)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse(fmt.Sprintf("无效的请求格式, err : %s", err.Error())))
		return
	}

	claims := helper.GetUserClaims(c)
	if claims == nil {
		c.JSON(400, gf_cluster.NewFailedResponse("校验身份出错"))
		return
	}
	createdUserId := claims.UserId
	createdUserName := claims.Name
	//同步创建集群
	for _, cluster := range instanceGroups {
		begin := time.Now()
		instanceGroup := gf_cluster.InstanceGroup{
			KubernetesId:  cluster.KubernetesId,
			Name:          cluster.Name,
			Image:         cluster.Image,
			Cpu:           cluster.Cpu,
			Memory:        cluster.Memory,
			Disk:          cluster.Disk,
			InstanceCount: cluster.InstanceCount,
			CreatedUser:   createdUserName,
			CreatedUserId: createdUserId,
		}
		err = instance.CreateInstanceGroup(&instanceGroup)
		if err != nil {
			logs.Logger.Error("创建实例组失败", zap.String("groupName", instanceGroup.Name), zap.Error(err))
			continue
		}
		err = instance.ExpandCustomInstanceGroup(&instanceGroup, cluster.InstanceCount)
		if err != nil {
			logs.Logger.Error("扩容实例组失败", zap.Int64("groupId", instanceGroup.Id), zap.String("groupName", instanceGroup.Name), zap.Int("count", instanceGroup.InstanceCount), zap.Error(err))
			continue
		}
		cost := time.Now().Sub(begin).Milliseconds()
		err = instance.AddInstanceForm(&instanceGroup, cost, createdUserId, createdUserName, gf_cluster.OptTypeExpand, instanceGroup.InstanceCount, err)
		if err != nil {
			logs.Logger.Error("记录日志失败", zap.Int64("groupId", instanceGroup.Id), zap.String("groupName", instanceGroup.Name), zap.Error(err))
			continue
		}
	}
	c.JSON(200, gf_cluster.NewSuccessResponse())
}

func HandleListInstanceGroup(c *gin.Context) {
	name, _ := c.GetQuery("name")
	pagerNumber, pageSize := helper.GetPagerParamFromQuery(c)
	instanceGroups, total, err := model.ListInstanceGroupFromDB(name, pagerNumber, pageSize)
	if err != nil {
		c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	pager := gf_cluster.Pager{
		PageNumber: pagerNumber,
		PageSize:   pageSize,
		Total:      int(total),
	}

	c.JSON(200, gf_cluster.NewListInstanceGroupResponse(instanceGroups, pager))
}

func HandleDeleteInstanceGroup(c *gin.Context) {
	begin := time.Now()
	clusterId, err := strconv.ParseInt(c.Param("instanceGroup"), 10, 64)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse("未指定实例组id"))
		return
	}
	instanceGroup, err := instance.GetInstanceGroup(clusterId)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	err = instance.DeleteInstanceGroup(instanceGroup)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	claims := helper.GetUserClaims(c)
	if claims == nil {
		c.JSON(400, gf_cluster.NewFailedResponse("校验身份出错"))
		return
	}
	createdUserId := claims.UserId
	createdUserName := claims.Name
	defer func() {
		cost := time.Now().Sub(begin).Milliseconds()
		err = instance.AddInstanceForm(instanceGroup, cost, createdUserId, createdUserName, gf_cluster.OptTypeShrink, instanceGroup.InstanceCount, err)
		if err != nil {
			c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
			return
		}
		if err == nil {
			c.JSON(200, gf_cluster.NewSuccessResponse())
		}
	}()
}

//HandleBatchDeleteInstanceGroup 批量删除集群
func HandleBatchDeleteInstanceGroup(c *gin.Context) {
	//解析请求体
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse("无效的请求格式"))
		return
	}
	var request gf_cluster.InstanceGroupBatchDeleteRequest
	err = json.Unmarshal(data, &request)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse(fmt.Sprintf("无效的请求格式, err : %s", err.Error())))
		return
	}

	claims := helper.GetUserClaims(c)
	if claims == nil {
		c.JSON(400, gf_cluster.NewFailedResponse("校验身份出错"))
		return
	}
	createdUserId := claims.UserId
	createdUserName := claims.Name
	for _, clusterId := range request.Ids {
		begin := time.Now()
		instanceGroup, err := instance.GetInstanceGroup(clusterId)
		if err != nil {
			logs.Logger.Error("获取实例组失败", zap.Int64("groupId", clusterId), zap.String("groupName", instanceGroup.Name), zap.Error(err))
			continue
		}
		err = instance.DeleteInstanceGroup(instanceGroup)
		if err != nil {
			logs.Logger.Error("删除实例组失败", zap.Int64("groupId", clusterId), zap.String("groupName", instanceGroup.Name), zap.Error(err))
			continue
		}
		cost := time.Now().Sub(begin).Milliseconds()
		err = instance.AddInstanceForm(instanceGroup, cost, createdUserId, createdUserName, gf_cluster.OptTypeShrink, instanceGroup.InstanceCount, err)
		if err != nil {
			logs.Logger.Error("记录操作日志失败", zap.Int64("groupId", clusterId), zap.String("groupName", instanceGroup.Name), zap.Error(err))
		}
	}
	c.JSON(200, gf_cluster.NewSuccessResponse())
}

//HandleGetInstanceGroup 获取实例组信息
func HandleGetInstanceGroup(c *gin.Context) {
	clusterId, err := strconv.ParseInt(c.Param("instanceGroup"), 10, 64)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse("未指定实例组id"))
		return
	}
	cluster, err := instance.GetInstanceGroup(clusterId)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	c.JSON(200, gf_cluster.NewGetInstanceGroupResponse(cluster))
}

//HandleUpdateInstanceGroup 更新实例组信息
func HandleUpdateInstanceGroup(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse("无效的请求格式"))
		return
	}
	var cluster gf_cluster.InstanceGroupUpdateRequest
	err = json.Unmarshal(data, &cluster)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse(fmt.Sprintf("无效的请求格式, err: %s", err.Error())))
		return
	}

	instanceGroup := gf_cluster.InstanceGroup{
		Id:            cluster.Id,
		KubernetesId:  cluster.KubernetesId,
		Name:          cluster.Name,
		Image:         cluster.Image,
		Cpu:           cluster.Cpu,
		Memory:        cluster.Memory,
		Disk:          cluster.Disk,
		InstanceCount: cluster.InstanceCount,
	}
	err = model.UpdateInstanceGroupFromDB(&instanceGroup)
	if err != nil {
		c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	c.JSON(200, gf_cluster.NewSuccessResponse())

}

//HandleExpandInstanceGroup 扩容实例组
func HandleExpandInstanceGroup(c *gin.Context) {
	//读取请求体
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse("无效的请求体"))
		return
	}
	var request gf_cluster.InstanceGroupExpandRequest
	err = json.Unmarshal(data, &request)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse(fmt.Sprintf("无效的请求体,err:%s", err.Error())))
		return
	}

	//获取实例组
	instanceGroup, err := instance.GetInstanceGroup(request.InstanceGroupId)
	if err != nil {
		c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	if request.Count <= 0 {
		c.JSON(400, gf_cluster.NewFailedResponse("实实例数量应该大于等于1"))
		return
	}
	if request.Count <= instanceGroup.InstanceCount {
		c.JSON(400, gf_cluster.NewFailedResponse("目标实例数应该大于当前实例数"))
		return
	}
	err = instance.ExpandCustomInstanceGroup(instanceGroup, request.Count)
	if err != nil {
		c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
		return
	}

	c.JSON(200, gf_cluster.NewSuccessResponse())
}

//HandleShrinkInstanceGroup 缩容实例组
func HandleShrinkInstanceGroup(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse("无效的请求体"))
		return
	}
	var request gf_cluster.InstanceGroupShrinkRequest
	err = json.Unmarshal(data, &request)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse(fmt.Sprintf("无效的请求体,err:%s", err.Error())))
		return
	}
	instanceGroup, err := instance.GetInstanceGroup(request.InstanceGroupId)
	if err != nil {
		c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	if request.Count <= 0 {
		c.JSON(400, gf_cluster.NewFailedResponse("实例数量应该大于等于1"))
		return
	}
	if request.Count >= instanceGroup.InstanceCount {
		c.JSON(400, gf_cluster.NewFailedResponse("目标实例数应该小于当前实例数"))
		return
	}

	err = instance.ShrinkCustomInstanceGroup(instanceGroup, request.Count)
	if err != nil {
		c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	c.JSON(200, gf_cluster.NewSuccessResponse())
}

//HandleExpandOrShrinkInstanceGroup 扩缩容接口
func HandleExpandOrShrinkInstanceGroup(c *gin.Context) {
	begin := time.Now()
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse("无效的请求体"))
		return
	}
	var request gf_cluster.InstanceGroupExpandOrShrinkRequest
	err = json.Unmarshal(data, &request)
	if err != nil {
		c.JSON(400, gf_cluster.NewFailedResponse(fmt.Sprintf("无效的请求体,err:%s", err.Error())))
		return
	}
	instanceGroup, err := instance.GetInstanceGroup(request.InstanceGroupId)
	if err != nil {
		c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
		return
	}
	if request.Count < 0 {
		c.JSON(400, gf_cluster.NewFailedResponse("实例数量应该大于等于0"))
		return
	}
	if request.Count == instanceGroup.InstanceCount {
		c.JSON(400, gf_cluster.NewFailedResponse("目标实例数不应该等于当前实例数"))
		return
	}
	var optType string
	var updatedInstanceCount int
	if request.Count > instanceGroup.InstanceCount {
		optType = gf_cluster.OptTypeExpand
		err = instance.ExpandCustomInstanceGroup(instanceGroup, request.Count)
		if err != nil {
			c.JSON(400, gf_cluster.NewFailedResponse(err.Error()))
			return
		}
		updatedInstanceCount = request.Count - instanceGroup.InstanceCount
	}
	if request.Count < instanceGroup.InstanceCount {
		optType = gf_cluster.OptTypeShrink
		err = instance.ShrinkCustomInstanceGroup(instanceGroup, request.Count)
		if err != nil {
			c.JSON(400, gf_cluster.NewFailedResponse(err.Error()))
			return
		}
		updatedInstanceCount = instanceGroup.InstanceCount - request.Count
	}
	claims := helper.GetUserClaims(c)
	if claims == nil {
		c.JSON(400, gf_cluster.NewFailedResponse("校验身份出错"))
		return
	}

	createdUserId := claims.UserId
	createdUserName := claims.Name

	defer func() {
		cost := time.Now().Sub(begin).Milliseconds()
		err = instance.AddInstanceForm(instanceGroup, cost, createdUserId, createdUserName, optType, updatedInstanceCount, err)
		if err != nil {
			c.JSON(500, gf_cluster.NewFailedResponse(err.Error()))
			return
		}
		if err == nil {
			c.JSON(200, gf_cluster.NewSuccessResponse())
		}
	}()
}
