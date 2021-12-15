package service

import (
	"context"
	"time"

	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/errs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/pkg/cmp"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm/schema"
)

type OperationLog struct {
	Operation string
	Operator  int64

	Old schema.Tabler
	New schema.Tabler
}

func RecordOperationLog(ctx context.Context, oplog OperationLog) error {
	if oplog.Operator == 0 {
		return errs.ErrOperatorIsNull
	}
	user, err := GetUserById(ctx, oplog.Operator)
	if err != nil {
		return err
	}
	if oplog.New.TableName() == "" || oplog.Old.TableName() == "" {
		return errs.ErrNewOrOldDataIsNull
	}
	res, err := cmp.Diff(oplog.Old, oplog.New)
	if err != nil {
		return err
	}
	diff, err := res.Beautiful()
	if err != nil {
		return err
	}
	diffStr, err := jsoniter.MarshalToString(diff)
	if err != nil {
		return err
	}
	now := time.Now()
	return clients.WriteDBCli.Create(model.OperationLog{
		Base: model.Base{
			CreateAt: &now,
			UpdateAt: &now,
		},
		Operation:  oplog.Operation,
		ObjectName: oplog.New.TableName(),
		Operator:   oplog.Operator,
		Diff:       diffStr,
		UserName:   user.Username,
	}).Error
}

func ExtractLogs(ctx context.Context, conds model.ExtractCondition) ([]model.OperationLog, int64, error) {
	logs, count, err := model.ExtractLogs(ctx, conds)
	if err != nil {
		return nil, 0, err
	}
	operators := make([]int64, 0, len(logs))
	for _, l := range logs {
		operators = append(operators, l.Operator)
	}
	userMap := UserMapByIDs(ctx, operators)
	for i, l := range logs {
		logs[i].UserName = userMap[l.Operator]
	}

	return logs, count, nil
}
