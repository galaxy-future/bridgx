package service

import (
	"context"

	"github.com/galaxy-future/BridgX/internal/model"
)

type ExtractLogsRequest struct{}

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
