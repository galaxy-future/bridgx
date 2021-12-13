package model

import (
	"context"
	"time"

	"github.com/galaxy-future/BridgX/internal/clients"
)

type OperationLog struct {
	Base
	Handler  string
	Params   string
	Info     string
	UserId   int64
	Response string
	UserName string `gorm:"-"`
}

func (OperationLog) TableName() string {
	return "operation_log"
}

type ExtractCondition struct {
	UserIDs    []int64
	Handlers   []string
	TimeStart  time.Time
	TimeEnd    time.Time
	PageNumber int
	PageSize   int
}

func ExtractLogs(ctx context.Context, conds ExtractCondition) (logs []OperationLog, count int64, err error) {
	query := clients.ReadDBCli.WithContext(ctx).Select(OperationLog{}.TableName())
	if len(conds.UserIDs) > 0 {
		query = query.Where("user_id IN (?)", conds.UserIDs)
	}
	if len(conds.Handlers) > 0 {
		query = query.Where("handler IN (?)", conds.Handlers)
	}
	if !conds.TimeStart.IsZero() {
		query = query.Where("create_at >= ?", conds.TimeStart)
	}
	if !conds.TimeEnd.IsZero() {
		query = query.Where("create_at < ?", conds.TimeEnd)
	}
	count, err = QueryWhere(query, conds.PageNumber, conds.PageSize, &logs, "id Desc", true)
	if err != nil {
		return nil, 0, err
	}
	return logs, count, nil
}
