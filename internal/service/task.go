package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/id_generator"
	"github.com/galaxy-future/BridgX/pkg/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

func CreateExpandTask(ctx context.Context, clusterName string, count int, taskName string, uid int64) (int64, error) {
	if hasUnfinishedTask(clusterName) {
		return 0, errors.New(fmt.Sprintf("Cluster:%v has unfinished task", clusterName))
	}
	cluster, err := model.GetByClusterName(clusterName)
	if err != nil {
		return 0, err
	}
	if cluster == nil {
		return 0, fmt.Errorf(constants.ErrClusterNotExist, clusterName)
	}
	currentCount, err := model.CountActiveInstancesByClusterName(ctx, []string{clusterName})
	if err != nil {
		return 0, err
	}
	info := &model.ExpandTaskInfo{
		ClusterName:    clusterName,
		Count:          count,
		TaskSubmitHost: utils.PrivateIPv4(),
		UserId:         uid,
		BeforeCount:    int(currentCount),
	}
	s, _ := jsoniter.MarshalToString(info)
	logs.Logger.Infof("cluster:%v expand task info:%v", clusterName, s)
	taskId := id_generator.GetNextId()
	task := &model.Task{
		TaskName:      taskName,
		TaskAction:    constants.TaskActionExpand,
		Status:        constants.TaskStatusInit,
		TaskFilter:    clusterName,
		TaskInfo:      s,
		SupportCancel: false,
	}
	now := time.Now()
	task.Id = int64(taskId)
	task.CreateAt = &now
	task.UpdateAt = &now
	err = model.Create(task)
	if err != nil {
		return 0, err
	}
	return task.Id, nil
}
func CreateShrinkTask(ctx context.Context, clusterName string, count int, ips string, taskName string, uid int64) (int64, error) {
	if hasUnfinishedTask(clusterName) {
		return 0, errors.New(fmt.Sprintf("Cluster:%v has unfinished task", clusterName))
	}
	cluster, err := model.GetByClusterName(clusterName)
	if err != nil {
		return 0, err
	}
	if cluster == nil {
		return 0, fmt.Errorf(constants.ErrClusterNotExist, clusterName)
	}
	if chargeType := cluster.GetChargeType(); chargeType == cloud.InstanceChargeTypePrePaid {
		return 0, errors.New(constants.ErrPrePaidShrinkNotSupported)
	}
	currentCount, err := model.CountActiveInstancesByClusterName(ctx, []string{clusterName})
	if err != nil {
		return 0, err
	}
	info := &model.ShrinkTaskInfo{
		ClusterName:    clusterName,
		Count:          count,
		IPs:            ips,
		TaskSubmitHost: utils.PrivateIPv4(),
		UserId:         uid,
		BeforeCount:    int(currentCount),
	}
	s, _ := jsoniter.MarshalToString(info)
	logs.Logger.Infof("cluster:%v shrink task info:%v", clusterName, s)
	taskId := id_generator.GetNextId()
	task := &model.Task{
		TaskName:      taskName,
		TaskAction:    constants.TaskActionShrink,
		Status:        constants.TaskStatusInit,
		TaskFilter:    clusterName,
		TaskInfo:      s,
		SupportCancel: false,
	}
	now := time.Now()
	task.Id = int64(taskId)
	task.CreateAt = &now
	task.UpdateAt = &now
	err = model.Create(task)
	if err != nil {
		return 0, err
	}
	return task.Id, nil
}

func hasUnfinishedTask(clusterName string) bool {
	cnt, err := model.CountByTaskStatus(clusterName, []string{constants.TaskStatusInit, constants.TaskStatusRunning})
	if err != nil {
		return false
	}
	return cnt != 0
}

func GetTaskCount(ctx context.Context, accountKeys []string) (int64, error) {
	clusterNames, err := GetEnabledClusterNamesByAccounts(ctx, accountKeys)
	if err != nil {
		return 0, err
	}
	ret, err := model.GetTaskCount(ctx, clusterNames)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func GetTaskListByAk(ctx context.Context, accountKey string, pageNum, pageSize int) ([]model.Task, int64, error) {
	clusterNames, err := GetEnabledClusterNamesByAccount(ctx, accountKey)
	if err != nil {
		return nil, 0, err
	}
	ret := make([]model.Task, 0)
	total, err := model.Query(map[string]interface{}{"task_filter": clusterNames}, pageNum, pageSize, &ret, "id desc", true)
	if err != nil {
		return ret, 0, err
	}
	return ret, total, nil
}

func GetTaskListByAks(ctx context.Context, accountKeys []string, pageNum, pageSize int) ([]model.Task, int64, error) {
	clusterNames, err := GetEnabledClusterNamesByAccounts(ctx, accountKeys)
	if err != nil {
		return nil, 0, err
	}
	ret := make([]model.Task, 0)
	total, err := model.Query(map[string]interface{}{"task_filter": clusterNames}, pageNum, pageSize, &ret, "id desc", true)
	if err != nil {
		return ret, 0, err
	}
	return ret, total, nil
}

func GetTaskListByCond(ctx context.Context, accountKey []string, cond model.TaskSearchCond) ([]model.Task, int64, error) {
	clusterNames, err := GetEnabledClusterNamesByCond(ctx, "", cond.ClusterName, accountKey, false)
	if err != nil {
		return nil, 0, err
	}
	ret := make([]model.Task, 0)
	ret, total, err := model.SearchTask(ctx, filterClusterNames(clusterNames, cond.ClusterName), cond)
	if err != nil {
		return ret, 0, err
	}
	return ret, total, nil
}

func filterClusterNames(clusterNames []string, searchName string) []string {
	ret := make([]string, 0, len(clusterNames))
	for _, name := range clusterNames {
		if strings.Contains(name, searchName) {
			ret = append(ret, name)
		}
	}
	return ret
}

func GetTask(ctx context.Context, taskId string) (*model.Task, error) {
	ret := &model.Task{}
	err := model.Get(cast.ToInt64(taskId), ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
