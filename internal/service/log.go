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
	userIds := make([]int64, 0, len(logs))
	for _, l := range logs {
		userIds = append(userIds, l.UserId)
	}
	userMap := UserMapByIDs(ctx, userIds)
	for i, l := range logs {
		logs[i].UserName = userMap[l.UserId]
	}

	return logs, count, nil
}
