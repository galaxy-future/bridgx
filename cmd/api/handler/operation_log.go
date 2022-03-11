package handler

import (
	"net/http"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/galaxy-future/BridgX/internal/types"
	"github.com/galaxy-future/BridgX/pkg/utils"
	"github.com/gin-gonic/gin"
)

func CreateLog(ctx *gin.Context) {
	req := service.OperationLog{}
	err := ctx.BindJSON(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = service.RecordOperationLog(ctx, req)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, nil)
	return
}

type ExtractLogsResponse struct {
	Logs  []Log
	Pager types.Pager
}
type Log struct {
	ID       int64  `json:"id"`
	Operator string `json:"operator"`
	Action   string `json:"action"`
	Detail   string `json:"detail"`
	ExecTime string `json:"exec_time"`
}

func ExtractLog(ctx *gin.Context) {
	// TODO: Org??
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusForbidden, response.PermissionDenied, nil)
		return
	}
	req := service.ExtractCondition{}
	err := ctx.BindQuery(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	logs, total, err := service.ExtractLogs(ctx, req)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, modelLog2Res(logs, int(total), req.PageNumber, req.PageSize))
	return
}

func modelLog2Res(logs []model.OperationLog, total, page, size int) ExtractLogsResponse {
	res := ExtractLogsResponse{
		Pager: types.Pager{
			PageNumber: page,
			PageSize:   size,
			Total:      total,
		},
	}
	for _, log := range logs {
		var execTime string
		if log.CreateAt != nil {
			execTime = utils.FormatTime(*log.CreateAt)
		}
		res.Logs = append(res.Logs, Log{
			ID:       log.Id,
			Operator: log.Operator,
			Action:   log.Action,
			Detail:   log.Detail,
			ExecTime: execTime,
		})
	}
	return res
}
