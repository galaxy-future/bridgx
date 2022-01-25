package handler

import (
	"net/http"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/cmd/api/middleware/validation"
	"github.com/galaxy-future/BridgX/cmd/api/request"
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func ListOrgs(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}

	orgs, err := service.GetOrgList(ctx)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	orgThumbs := make([]response.OrgThumb, 0, len(orgs))
	for _, org := range orgs {
		if org.Id > 0 {
			orgThumb := response.OrgThumb{
				OrgId:   cast.ToString(org.Id),
				OrgName: org.OrgName,
			}
			userNum, _ := service.CountUser(ctx, org.Id)
			orgThumb.UserNum = cast.ToString(userNum)
			orgThumbs = append(orgThumbs, orgThumb)
		}
	}

	resp := &response.ListOrgsResponse{
		OrgList: orgThumbs,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func CreateOrg(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}
	req := request.CreateOrgRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}

	err := service.CreateOrg(ctx, req.OrgName, req.UserName, req.Password, user.Name)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

func GetOrgById(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}
	idParam := ctx.Param("id")
	logs.Logger.Infof("get org idParam is:%v ", idParam)
	id, err := cast.ToInt64E(idParam)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}

	org, err := service.GetOrgInfoById(ctx, id)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, org)
	return
}

func EditOrg(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TokenInvalid, nil)
		return
	}

	req := request.EditOrgRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, validation.Translate2Chinese(err), nil)
		return
	}
	if user.Name != constants.UserNameRoot {
		response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
		return
	}
	err := service.EditOrg(ctx, user.OrgId, req.OrgName)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}
