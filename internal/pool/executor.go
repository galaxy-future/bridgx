package pool

import (
	"context"
	"strings"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/galaxy-future/BridgX/pkg/utils"
	jsoniter "github.com/json-iterator/go"
)

func doExpand(task *model.Task) {
	logs.Logger.Infof("Executing Task:%v, %v [%v], task info:%v", task.Id, task.TaskAction, task.TaskFilter, task.TaskInfo)
	taskInfo := &model.ExpandTaskInfo{}
	err := jsoniter.UnmarshalFromString(task.TaskInfo, taskInfo)
	if err != nil {
		taskFailed(task, err)
		return
	}
	taskInfo.TaskExecHost = utils.PrivateIPv4()
	task.TaskInfo, _ = jsoniter.MarshalToString(taskInfo)
	cluster, err := model.GetByClusterName(taskInfo.ClusterName)
	if err != nil {
		taskFailed(task, err)
		return
	}
	tags, _ := service.GetClusterTagsByClusterName(context.Background(), taskInfo.ClusterName)
	clusterInfo, err := service.ConvertToClusterInfo(cluster, tags)
	if err != nil {
		taskFailed(task, err)
		return
	}
	availableIds, allIds, err := service.ExpandCluster(clusterInfo, taskInfo.Count, task.Id)
	if len(availableIds) == taskInfo.Count {
		taskSuccess(task, taskInfo.Count)
	} else {
		successNum := service.RepairCluster(clusterInfo, task.Id, availableIds, allIds)
		if successNum == 0 {
			taskFailed(task, err)
		} else {
			taskPartialSuccess(task, successNum, err)
		}
	}
}

// DoExpand for test
func DoExpand(task *model.Task) {
	doExpand(task)
}

func saveTaskResult(task *model.Task, taskResult *model.TaskResult, stat string, err error) {
	task.Status = stat
	if err != nil {
		task.ErrMsg = err.Error()
	}
	task.TaskResult, err = jsoniter.MarshalToString(taskResult)
	if err != nil {
		logs.Logger.Warnf("saveTaskResult taskResult to string failed, %v", err)
	}
	ft := time.Now()
	task.FinishTime = &ft
	_ = model.Save(task)
	logs.Logger.Warnf("Task %v:%v, %v, %v", stat, task.Id, task.TaskAction, task.TaskInfo)
}

func taskPartialSuccess(task *model.Task, successNum int, err error) {
	saveTaskResult(task, &model.TaskResult{SuccessNum: successNum}, constants.TaskStatusPartialSuccess, err)
}

func taskSuccess(task *model.Task, successNum int) {
	saveTaskResult(task, &model.TaskResult{SuccessNum: successNum}, constants.TaskStatusSuccess, nil)
}

func taskFailed(task *model.Task, err error) {
	saveTaskResult(task, &model.TaskResult{}, constants.TaskStatusFailed, err)
}

func doShrink(task *model.Task) {
	logs.Logger.Infof("Executing Task:%v, %v [%v], task info:%v", task.Id, task.TaskAction, task.TaskFilter, task.TaskInfo)
	taskInfo := &model.ShrinkTaskInfo{}
	err := jsoniter.UnmarshalFromString(task.TaskInfo, taskInfo)
	if err != nil {
		taskFailed(task, err)
		return
	}
	taskInfo.TaskExecHost = utils.PrivateIPv4()
	task.TaskInfo, _ = jsoniter.MarshalToString(taskInfo)
	cluster, err := model.GetByClusterName(taskInfo.ClusterName)
	if err != nil {
		taskFailed(task, err)
		return
	}
	tags, _ := service.GetClusterTagsByClusterName(context.Background(), taskInfo.ClusterName)
	clusterInfo, err := service.ConvertToClusterInfo(cluster, tags)
	if err != nil {
		taskFailed(task, err)
		return
	}
	deletingIPs := calcDeletingIPs(taskInfo.IPs)
	shrink := func(attempt uint) error {
		logs.Logger.Infof("shrink cluster:%v with retry times:%v", clusterInfo.Name, attempt)
		if deletingIPs > 0 {
			return service.ShrinkClusterBySpecificIps(clusterInfo, taskInfo.IPs, taskInfo.Count, task.Id)
		} else {
			return service.ShrinkCluster(clusterInfo, taskInfo.Count, task.Id)
		}
	}
	err = retry.Retry(shrink, strategy.Limit(3), strategy.Backoff(backoff.BinaryExponential(time.Second)))
	if err != nil {
		taskFailed(task, err)
		return
	}
	taskSuccess(task, taskInfo.Count)
}

func calcDeletingIPs(IPs string) int {
	if IPs == "" || IPs == constants.HasNoneIP {
		return 0
	}
	return len(strings.Split(IPs, ","))
}
