package model

import (
	"context"
	"time"

	"github.com/galaxy-future/BridgX/internal/clients"
)

type OperationLog struct {
	Base
	Operation  string `gorm:"column:operation"`   // edit/create/delete
	ObjectName string `gorm:"column:object_name"` // target model's table name
	Operator   int64  `gorm:"column:operator"`
	Diff       string `gorm:"column:diff"`

	UserName string `gorm:"-"`
}

func (OperationLog) TableName() string {
	return "operation_log"
}

type ExtractCondition struct {
	Operators  []int64
	Operations []string
	TimeStart  time.Time
	TimeEnd    time.Time
	PageNumber int
	PageSize   int
}

func ExtractLogs(ctx context.Context, conds ExtractCondition) (logs []OperationLog, count int64, err error) {
	query := clients.ReadDBCli.WithContext(ctx).Select(OperationLog{}.TableName())
	if len(conds.Operators) > 0 {
		query.Where("operator IN (?)", conds.Operators)
	}
	if len(conds.Operations) > 0 {
		query.Where("operation IN (?)", conds.Operations)
	}
	if !conds.TimeStart.IsZero() {
		query.Where("create_at >= ?", conds.TimeStart)
	}
	if !conds.TimeEnd.IsZero() {
		query.Where("create_at < ?", conds.TimeEnd)
	}
	count, err = QueryWhere(query, conds.PageNumber, conds.PageSize, &logs, "id Desc", true)
	if err != nil {
		return nil, 0, err
	}
	return logs, count, nil
}
