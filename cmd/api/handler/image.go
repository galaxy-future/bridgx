package handler

import (
	"net/http"

	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/gin-gonic/gin"
)

func GetImageList(ctx *gin.Context) {
	account, err := GetOrgKeys(ctx)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	regionId, ok := ctx.GetQuery("region_id")
	if !ok {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	provider, ok := ctx.GetQuery("provider")
	if !ok {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	images, err := service.GetImages(ctx, service.GetImagesRequest{
		Account:  account,
		Provider: provider,
		RegionId: regionId,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, images)
}
