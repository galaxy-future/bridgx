package handler

import (
	"net/http"

	"github.com/galaxy-future/BridgX/cmd/api/middleware/validation"
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/gin-gonic/gin"
)

type GetImageListRequest struct {
	RegionID string `json:"region_id" binding:"required" form:"region_id"`
	Provider string `json:"provider" binding:"required,mustIn=cloud" form:"provider"`
}

func GetImageList(ctx *gin.Context) {
	account, err := GetOrgKeys(ctx)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	req := GetImageListRequest{}
	err = ctx.Bind(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, validation.Translate2Chinese(err), nil)
		return
	}
	images, err := service.GetImages(ctx, service.GetImagesRequest{
		Account:  account,
		Provider: req.Provider,
		RegionId: req.RegionID,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, images)
}
