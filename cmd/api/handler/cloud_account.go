package handler

import (
	"net/http"
	"strings"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/cmd/api/request"
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func CreateCloudAccount(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}
	req := request.CreateCloudAccountRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	validErr := service.CheckAccountValid(req.AccountKey, req.AccountSecret)
	if validErr != nil {
		response.MkResponse(ctx, http.StatusBadRequest, validErr.Error(), nil)
		return
	}
	err := service.CreateCloudAccount(ctx, req.AccountName, req.Provider, req.AccountKey, req.AccountSecret, user.OrgId, user.Name)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

func ListCloudAccounts(ctx *gin.Context) {
	provider, _ := ctx.GetQuery("provider")
	accountName, _ := ctx.GetQuery("account_name")
	account, _ := ctx.GetQuery("account")
	pageNum, pageSize := getPager(ctx)
	accounts, total, err := service.GetAccounts(provider, accountName, account, pageNum, pageSize)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	pager := response.Pager{
		PageNumber: pageNum,
		PageSize:   pageSize,
		Total:      int(total),
	}
	resp := &response.ListCloudAccountResponse{
		CloudAccountList: helper.ConvertToCloudAccountList(accounts),
		Pager:            pager,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func EditCloudAccount(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}
	req := request.EditCloudAccountRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	err := service.EditCloudAccount(ctx, cast.ToInt64(req.AccountId), req.AccountName, req.Provider, user.Name)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

func DeleteCloudAccount(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	idParam := ctx.Param("ids")
	input := strings.Split(idParam, ",")
	ids := make([]int64, 0)
	for _, v := range input {
		ids = append(ids, cast.ToInt64(v))
	}
	err := service.DeleteCloudAccount(ctx, ids, user.OrgId)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

//GetAccountInfo get accountKey and encrypt accountSecret by clusterName
func GetAccountInfo(ctx *gin.Context) {
	clusterName, ok := ctx.GetQuery("cluster_name")
	if !ok || clusterName == "" {
		response.MkResponse(ctx, http.StatusBadRequest, "missing param cluster_name", nil)
		return
	}
	cluster, err := service.GetClusterByName(ctx, clusterName)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	account, err := service.GetAccount(cluster.Provider, cluster.AccountKey)
	if account == nil {
		response.MkResponse(ctx, http.StatusBadRequest, "account not found", nil)
		return
	}
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp, err := helper.ConvertToEncryptAccountInfo(account)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}
