package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/cmd/api/middleware/validation"
	"github.com/galaxy-future/BridgX/cmd/api/request"
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/galaxy-future/BridgX/internal/types"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

func GetClusterById(ctx *gin.Context) {
	idParam := ctx.Param("id")
	logs.Logger.Infof("idParam is:%v ", idParam)
	id, err := cast.ToInt64E(idParam)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	cm, err := service.GetClusterById(ctx, id)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	tags, err := service.GetClusterTagsByClusterName(ctx, cm.ClusterName)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp, err := service.ConvertToClusterInfo(cm, tags)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func GetInstanceStat(ctx *gin.Context) {
	clusterName, ok := ctx.GetQuery("cluster_name")
	if !ok || clusterName == "" {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	cluster, err := service.GetClusterByName(ctx, clusterName)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	instanceType := service.GetInstanceTypeByName(cluster.InstanceType)
	instanceCount, err := service.GetInstanceCount(ctx, nil, clusterName)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, helper.ConvertToInstanceStat(instanceType, instanceCount))
	return
}

func GetClusterCount(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	accountKey := ctx.Query("account")
	accountKeys, err := service.GetAksByOrgAkProvider(ctx, user.GetOrgIdForTest(), accountKey, "")
	cnt, err := service.GetClusterCount(ctx, accountKeys)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp := &response.ClusterCountResponse{
		ClusterNum: cnt,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func ListClustersByTags(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}
	req := request.ListClusterByTagsRequest{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	pn, ps := buildPager(req.PageNumber, req.PageSize)
	clusters, total, err := service.GetClustersByTags(ctx, req.Tags, ps, pn)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	pager := response.Pager{
		PageNumber: pn,
		PageSize:   ps,
		Total:      int(total),
	}
	tags, err := service.GetClusterTagsByClusters(ctx, clusters)
	resp := &response.ListClustersWithTagResponse{
		ClusterList: helper.ConvertToClusterThumbListWithTag(clusters, tags),
		Pager:       pager,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func buildPager(pageNumber int, pageSize int) (int, int) {
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 || pageSize > constants.DefaultPageSize {
		pageSize = constants.DefaultPageSize
	}
	return pageNumber, pageSize
}

func ListClusters(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}
	var accountKeys []string
	accountKey, ok := ctx.GetQuery("account")
	if !ok || accountKey == "" {
		var err error
		accountKeys, err = service.GetAksByOrgId(user.OrgId)
		if err != nil {
			response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	} else {
		accountKeys = append(accountKeys, accountKey)
	}
	clusterName, _ := ctx.GetQuery("cluster_name")
	provider, _ := ctx.GetQuery("provider")

	pn, ps := getPager(ctx)

	clusters, total, err := service.ListClusters(ctx, accountKeys, clusterName, provider, pn, ps)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	instanceCountMap := service.GetInstanceCountByCluster(ctx, clusters)
	pager := response.Pager{
		PageNumber: pn,
		PageSize:   ps,
		Total:      total,
	}
	tags, err := service.GetClusterTagsByClusters(ctx, clusters)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp := &response.ListClustersResponse{
		ClusterList: helper.ConvertToClusterThumbList(clusters, instanceCountMap, tags),
		Pager:       pager,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func GetClusterByName(ctx *gin.Context) {
	name := ctx.Param("name")
	logs.Logger.Infof("name is:%v ", name)
	cm, err := service.GetClusterByName(ctx, name)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	tags, err := service.GetClusterTagsByClusterName(ctx, cm.ClusterName)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp, err := service.ConvertToClusterInfo(cm, tags)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func CreateCluster(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}
	clusterInput := types.ClusterInfo{}
	err := ctx.ShouldBindJSON(&clusterInput)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, validation.Translate2Chinese(err), nil)
		return
	}
	m, err := convertToClusterModel(&clusterInput)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	tags := make([]model.ClusterTag, 0)
	tags = append(tags, model.ClusterTag{
		ClusterName: clusterInput.Name,
		TagKey:      constants.DefaultClusterUsageKey,
		TagValue:    constants.DefaultClusterUsageUnused,
	})
	for k, v := range clusterInput.Tags {
		if k == "" || v == "" {
			response.MkResponse(ctx, http.StatusBadRequest, "empty key/value for tags", nil)
			return
		}
		tag := model.ClusterTag{
			ClusterName: clusterInput.Name,
			TagKey:      k,
			TagValue:    v,
		}
		tags = append(tags, tag)
	}
	err = service.CreateCluster(m, user.Name)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if len(tags) > 0 {
		err = service.CreateClusterTags(tags)
		if err != nil {
			response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, m.ClusterName)
	return
}

func EditCluster(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}
	clusterInput := types.ClusterInfo{}
	err := ctx.ShouldBindJSON(&clusterInput)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, validation.Translate2Chinese(err), err)
		return
	}
	m, err := convertToClusterModel(&clusterInput)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	err = service.EditCluster(m, user.Name)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, m.ClusterName)
	return
}

func convertToClusterModel(clusterInput *types.ClusterInfo) (*model.Cluster, error) {
	if clusterInput.NetworkConfig == nil {
		return nil, errors.New("missing network config")
	}
	if clusterInput.StorageConfig == nil {
		return nil, errors.New("missing storage config")
	}
	if clusterInput.ChargeConfig == nil {
		return nil, errors.New("missing charge config")
	}
	nc, _ := jsoniter.MarshalToString(clusterInput.NetworkConfig)
	sc, _ := jsoniter.MarshalToString(clusterInput.StorageConfig)
	cc, _ := jsoniter.MarshalToString(clusterInput.ChargeConfig)
	m := model.Cluster{
		ClusterName:  clusterInput.Name,
		ClusterDesc:  clusterInput.Desc,
		RegionId:     clusterInput.RegionId,
		ZoneId:       clusterInput.ZoneId,
		InstanceType: clusterInput.InstanceType,
		Image:        clusterInput.Image,
		Password:     clusterInput.Password,
		Provider:     clusterInput.Provider,
		AccountKey:   clusterInput.AccountKey,

		NetworkConfig: nc,
		StorageConfig: sc,
		ChargeConfig:  cc,
	}
	return &m, nil
}

func AddClusterTags(ctx *gin.Context) {
	req := request.TagRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	tags := make([]model.ClusterTag, 0)
	for k, v := range req.Tags {
		if k == "" || v == "" {
			response.MkResponse(ctx, http.StatusBadRequest, "empty key/value for tags", nil)
			return
		}
		tag := model.ClusterTag{
			ClusterName: req.ClusterName,
			TagKey:      k,
			TagValue:    v,
		}
		tags = append(tags, tag)
	}
	err = service.CreateClusterTags(tags)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

func GetClusterTags(ctx *gin.Context) {
	req := request.GetTagsRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	pn, ps := buildPager(req.PageNumber, req.PageSize)
	ret, total, err := service.GetClusterTags(ctx, req.ClusterName, req.TagKey, pn, ps)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	pager := response.Pager{PageNumber: pn, PageSize: ps, Total: int(total)}
	resp := response.ClusterTagsResponse{ClusterTags: helper.ConvertToClusterTags(ret), Pager: pager}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func EditClusterTags(ctx *gin.Context) {
	req := request.TagRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	tags := make([]model.ClusterTag, 0)
	for k, v := range req.Tags {
		if k == "" || v == "" {
			response.MkResponse(ctx, http.StatusBadRequest, "empty key/value for tags", nil)
			return
		}
		tag := model.ClusterTag{
			ClusterName: req.ClusterName,
			TagKey:      k,
			TagValue:    v,
		}
		tags = append(tags, tag)
	}
	err = service.EditClusterTags(tags)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

func DeleteClusterTags(ctx *gin.Context) {
	req := request.TagRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	tags := make([]model.ClusterTag, 0)
	for k, v := range req.Tags {
		if k == "" {
			continue
		}
		tag := model.ClusterTag{
			ClusterName: req.ClusterName,
			TagKey:      k,
			TagValue:    v,
		}
		tags = append(tags, tag)
	}
	err = service.DeleteClusterTags(tags)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

func DeleteClusters(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	idParam := ctx.Param("ids")
	input := strings.Split(idParam, ",")
	ids := make([]int64, 0)
	for _, v := range input {
		ids = append(ids, cast.ToInt64(v))
	}
	err := service.DeleteClusters(ctx, ids, user.OrgId)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

func ExpandCluster(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
		return
	}
	req := request.ExpandClusterRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, validation.Translate2Chinese(err), nil)
		return
	}
	taskId, err := service.CreateExpandTask(ctx, req.ClusterName, req.Count, req.TaskName, user.UserId)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, taskId)
	return
}

func ShrinkCluster(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
		return
	}
	req := request.ShrinkClusterRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, validation.Translate2Chinese(err), nil)
		return
	}
	taskId, err := service.CreateShrinkTask(ctx, req.ClusterName, req.Count, strings.Join(req.IPs, ","), req.TaskName, user.UserId)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, taskId)
	return
}

func ShrinkAllInstances(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
		return
	}
	req := request.ShrinkAllInstancesRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, validation.Translate2Chinese(err), nil)
		return
	}
	taskId, err := service.CreateShrinkAllTask(ctx, req.ClusterName, req.TaskName, user.UserId)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, taskId)
	return
}
