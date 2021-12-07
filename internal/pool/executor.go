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
	instances, err := service.ExpandCluster(clusterInfo, taskInfo.Count, task.Id)
	if len(instances) == taskInfo.Count {
		taskSuccess(task, "")
	} else {
		instanceIds := make([]string, 0)
		for _, instance := range instances {
			instanceIds = append(instanceIds, instance.Id)
		}
		_ = service.RepairCluster(clusterInfo, task.Id, instanceIds)
		if len(instances) == 0 {
			taskFailed(task, err)
		} else {
			taskPartialSuccess(task, err)
		}
	}
}

// DoExpand for test
func DoExpand(task *model.Task) {
	doExpand(task)
}

func taskPartialSuccess(task *model.Task, err error) {
	task.Status = constants.TaskStatusPartialSuccess
	if err != nil {
		task.ErrMsg = err.Error()
	}
	ft := time.Now()
	task.FinishTime = &ft
	_ = model.Save(task)
	logs.Logger.Warnf("Task PartialSuccess:%v, %v, %v", task.Id, task.TaskAction, task.TaskInfo)
	_ = model.Save(task)
}

func taskSuccess(task *model.Task, s string) {
	task.TaskResult = s
	task.Status = constants.TaskStatusSuccess
	ft := time.Now()
	task.FinishTime = &ft
	logs.Logger.Warnf("Task Success:%v, %v, %v", task.Id, task.TaskAction, task.TaskInfo)
	_ = model.Save(task)
}

func taskFailed(task *model.Task, err error) {
	task.Status = constants.TaskStatusFailed
	if err != nil {
		task.ErrMsg = err.Error()
	}
	ft := time.Now()
	task.FinishTime = &ft
	_ = model.Save(task)
	logs.Logger.Warnf("Task Failed:%v, %v, %v", task.Id, task.TaskAction, task.TaskInfo)
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
	taskSuccess(task, "")
}

func calcDeletingIPs(IPs string) int {
	if IPs == "" || IPs == constants.HasNoneIP {
		return 0
	}
	return len(strings.Split(IPs, ","))
}
