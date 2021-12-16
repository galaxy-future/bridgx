package handler

import (
	"github.com/galaxy-future/BridgX/internal/constants"
	"net/http"
	"strings"
	"time"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/cmd/api/middleware/validation"
	"github.com/galaxy-future/BridgX/cmd/api/request"
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/gin-gonic/gin"
)

func GetInstanceCount(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
		return
	}
	clusterName := ctx.Query("cluster_name")
	accountKey := ctx.Query("account")
	accountKeys, err := service.GetAksByOrgAkProvider(ctx, user.OrgId, accountKey, "")
	cnt, err := service.GetInstanceCount(ctx, accountKeys, clusterName)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp := &response.InstanceCountResponse{
		InstanceNum: cnt,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func GetInstance(ctx *gin.Context) {
	instanceId, ok := ctx.GetQuery("instance_id")
	if !ok {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return

	}
	instance, err := service.GetInstance(ctx, instanceId)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	detail, err := helper.ConvertToInstanceDetail(ctx, instance)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp := &response.InstanceResponse{
		InstanceDetail: *detail,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func GetInstanceList(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
		return
	}
	accountKey := ctx.Query("account")
	instanceId := ctx.Query("instance_id")
	ip := ctx.Query("ip")
	provider := ctx.Query("provider")
	clusterName := ctx.Query("cluster_name")
	statusStr := ctx.Query("status")
	var status []string
	if statusStr != "" {
		status = strings.Split(statusStr, ",")
		for i, st := range status {
			status[i] = strings.ToUpper(st)
		}
	}
	accountKeys, err := service.GetAksByOrgAkProvider(ctx, user.OrgId, accountKey, provider)
	pn, ps := getPager(ctx)
	clusterNames, instances, total, err := service.GetInstancesByAccounts(ctx, accountKeys, status, pn, ps, instanceId, ip, clusterName)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	clusters, err := service.GetClustersByNames(ctx, clusterNames)
	pager := response.Pager{
		PageNumber: pn,
		PageSize:   ps,
		Total:      int(total),
	}
	resp := &response.InstanceListResponse{
		InstanceList: helper.ConvertToInstanceThumbList(ctx, instances, clusters),
		Pager:        pager,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func GetInstanceUsageTotal(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
		return
	}
	clusterName := ctx.Query("cluster_name")
	dateStr := ctx.Query("date")
	date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	usageTotal, err := service.GetInstanceUsageTotal(ctx, clusterName, date, user.OrgId)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, gin.H{
		"usage_total": usageTotal,
	})
	return
}

func GetInstanceUsageStatistics(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
		return
	}
	clusterName := ctx.Query("cluster_name")
	dateStr := ctx.Query("date")
	date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	pn, ps := getPager(ctx)
	instances, total, err := service.GetInstanceUsageStatistics(ctx, clusterName, date, user.OrgId, pn, ps)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	pager := response.Pager{
		PageNumber: pn,
		PageSize:   ps,
		Total:      int(total),
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, &response.InstanceUsageResponse{
		InstanceList: helper.ConvertToInstanceUsageList(ctx, instances),
		Pager:        pager,
	})
	return
}

func SyncInstanceExpireTime(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
		return
	}
	req := request.SyncInstanceExpireTimeRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, validation.Translate2Chinese(err), nil)
		return
	}
	err = service.SyncInstanceExpireTime(ctx, req.ClusterName)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

func ListRegions(ctx *gin.Context) {
	account, err := GetOrgKeys(ctx)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	provider, ok := ctx.GetQuery("provider")
	if !ok || provider == "" {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	logs.Logger.Infof("provider:[%s], %v", provider, account)
	regions, err := service.GetRegions(ctx, service.GetRegionsRequest{
		Provider: provider,
		Account:  account,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, regions)
}

func ListZones(ctx *gin.Context) {
	account, err := GetOrgKeys(ctx)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	provider := ctx.Query("provider")
	regionId := ctx.Query("region_id")
	logs.Logger.Infof("provider:[%s] regionId[:%s]", provider, regionId)
	zones, err := service.GetZones(ctx, service.GetZonesRequest{
		Provider: provider,
		RegionId: regionId,
		Account:  account,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, zones)
}

func ListInstanceType(ctx *gin.Context) {
	account, err := GetOrgKeys(ctx)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	provider := ctx.Query("provider")
	regionId := ctx.Query("region_id")
	zoneId := ctx.Query("zone_id")
	if !checkListInstanceTypeParams(provider, regionId, zoneId) {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	logs.Logger.Infof("provider:[%s] regionId:[%s] zoneId:[%s]", provider, regionId, zoneId)
	zones, err := service.ListInstanceType(ctx, service.ListInstanceTypeRequest{
		Provider: provider,
		RegionId: regionId,
		ZoneId:   zoneId,
		Account:  account,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	filterByComputingPowerType(ctx, zones)
}

func filterByComputingPowerType(ctx *gin.Context, zones service.ListInstanceTypeResponse) {
	computingPowerType := ctx.Query("computing_power_type")
	instanceTypes := zones.InstanceTypes
	if computingPowerType == "" {
		response.MkResponse(ctx, http.StatusOK, response.Success, instanceTypes)
	}

	ret := make([]service.InstanceTypeByZone, len(instanceTypes))
	if computingPowerType == constants.GPU {
		for _, instanceType := range instanceTypes {
			if strings.Contains(instanceType.InstanceTypeFamily, constants.IsGpu) {
				ret = append(ret, instanceType)
			}
		}
		response.MkResponse(ctx, http.StatusOK, response.Success, ret)
	}

	if computingPowerType == constants.CPU {
		for _, instanceType := range instanceTypes {
			if !strings.Contains(instanceType.InstanceTypeFamily, constants.IsGpu) {
				ret = append(ret, instanceType)
			}
		}
		response.MkResponse(ctx, http.StatusOK, response.Success, ret)
	}
}

func checkListInstanceTypeParams(provider, regionId, zoneId string) bool {
	// TODO: 后续得查一下缓存而不是值判断是否为空
	return provider != "" && regionId != "" && zoneId != ""
}
