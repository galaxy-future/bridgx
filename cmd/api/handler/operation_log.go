package handler

import (
	"net/http"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/cmd/api/middleware/validation"
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/galaxy-future/BridgX/internal/types"
	"github.com/galaxy-future/BridgX/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ExtractLogRequest struct {
	UserIDs    []int64  `json:"user_ids" form:"user_ids"`
	Handlers   []string `json:"handlers" form:"handlers"`
	TimeStart  string   `json:"time_start" form:"time_start"`
	TimeEnd    string   `json:"time_end" form:"time_end"`
	PageNumber int      `json:"page_number" form:"page_number"`
	PageSize   int      `json:"page_size" form:"page_size"`
}
type ExtractLogsResponse struct {
	Logs  []Log
	Pager types.Pager
}
type Log struct {
	ID              int64  `json:"id"`
	UserID          int64  `json:"user_id"`
	UserName        string `json:"user_name"`
	Operation       string `json:"operation"`
	OperationDetail string `json:"operation_detail"`
	ExecTime        string `json:"exec_time"`
}

func ExtractLog(ctx *gin.Context) {
	// TODO: Org??
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusForbidden, response.PermissionDenied, nil)
		return
	}
	req := ExtractLogRequest{}
	err := ctx.BindQuery(&req)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, validation.Translate2Chinese(err), nil)
		return
	}

	logs, total, err := service.ExtractLogs(ctx, model.ExtractCondition{
		UserIDs:    req.UserIDs,
		Handlers:   req.Handlers,
		TimeStart:  utils.ParseTime(req.TimeStart),
		TimeEnd:    utils.ParseTime(req.TimeEnd),
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
	})

	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
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
		reader := helper.GetLogReader(log.Handler)
		res.Logs = append(res.Logs, Log{
			ID:              log.Id,
			UserID:          log.UserId,
			UserName:        log.UserName,
			Operation:       reader.GetOperation(log.Handler),
			OperationDetail: reader.GetOperationDetail(log.Info),
			ExecTime:        utils.FormatTime(*log.CreateAt),
		})
	}
	return res
}
