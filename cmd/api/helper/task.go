package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/service"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

func ConvertToTaskDetail(instances []model.Instance, task *model.Task) *response.TaskDetailResponse {
	if len(instances) == 0 {
		return defaultTaskDetailByType(task)
	}
	ret := &response.TaskDetailResponse{}
	ret.TaskName = task.TaskName
	ret.TaskAction = task.TaskAction
	ret.TaskStatus = task.Status
	ret.ClusterName = task.TaskFilter
	ret.TaskResult = task.TaskResult
	ret.FailReason = task.ErrMsg
	ret.CreateAt = task.CreateAt.String()
	ret.TaskId = cast.ToString(task.Id)
	var running, suspending, success, fail, total int
	for _, instance := range instances {
		total++
		switch instance.Status {
		case constants.Undefined:
			suspending++
		case constants.Pending:
			running++
		case constants.Timeout:
			fail++
		case constants.Starting:
			running++
		case constants.Running:
			success++
		case constants.Deleted:
			success++
		case constants.Deleting:
			success++
		}
	}
	successRate := fmt.Sprintf("%0.2f", float64(success)/float64(total))
	ret.FailNum = fail
	ret.SuspendNum = suspending
	ret.RunNum = running
	ret.SuccessNum = success
	ret.TotalNum = total
	ret.SuccessRate = successRate
	endTime := time.Now()
	if task.FinishTime != nil {
		endTime = *task.FinishTime
	}
	ret.ExecuteTime = int(endTime.Sub(*task.CreateAt).Seconds())
	taskInfo := ExtractTaskInfo(task)
	ret.BeforeInstanceCount = taskInfo.GetBeforeInstanceCount()
	ret.AfterInstanceCount = taskInfo.GetAfterInstanceCount(success)
	ret.ExpectInstanceCount = taskInfo.GetExpectInstanceCount()
	ret.CreateBy = taskInfo.GetCreateUsername()
	return ret
}

func ExtractTaskInfo(task *model.Task) model.TaskInfo {
	var info model.TaskInfo
	switch task.TaskAction {
	case constants.TaskActionExpand:
		info = &model.ExpandTaskInfo{}
		_ = jsoniter.UnmarshalFromString(task.TaskInfo, info)
	case constants.TaskActionShrink:
		info = &model.ShrinkTaskInfo{}
		_ = jsoniter.UnmarshalFromString(task.TaskInfo, info)
	}
	return info
}

func defaultTaskDetailByType(task *model.Task) *response.TaskDetailResponse {
	if task == nil {
		return nil
	}
	endTime := time.Now()
	if task.FinishTime != nil {
		endTime = *task.FinishTime
	}
	resp := &response.TaskDetailResponse{
		TaskName:    task.TaskName,
		ClusterName: task.TaskFilter,
		TaskStatus:  task.Status,
		TaskResult:  task.TaskResult,
		TaskAction:  task.TaskAction,
		FailReason:  task.ErrMsg,
		TaskId:      cast.ToString(task.Id),
		CreateAt:    task.CreateAt.String(),
		ExecuteTime: int(endTime.Sub(*task.CreateAt).Seconds()),
	}
	if task.TaskAction == constants.TaskActionExpand {
		resp.SuccessRate = "0.00"
		taskInfo := model.ExpandTaskInfo{}
		_ = jsoniter.UnmarshalFromString(task.TaskInfo, &taskInfo)
		resp.TotalNum = taskInfo.Count
		if task.Status == constants.TaskStatusFailed {
			resp.FailNum = taskInfo.Count
		}
		resp.CreateBy = taskInfo.GetCreateUsername()
	}
	if task.TaskAction == constants.TaskActionShrink {
		taskInfo := model.ShrinkTaskInfo{}
		_ = jsoniter.UnmarshalFromString(task.TaskInfo, &taskInfo)
		resp.CreateBy = taskInfo.GetCreateUsername()
		if task.Status == constants.TaskStatusSuccess {
			resp.SuccessRate = "1.00"
			resp.SuccessNum = taskInfo.Count
			resp.TotalNum = taskInfo.Count
		} else {
			resp.SuccessRate = "0.00"
			if task.Status == constants.TaskStatusFailed {
				resp.FailNum = taskInfo.Count
			}
			resp.TotalNum = taskInfo.Count
		}
	}
	return resp
}

func ConvertToTaskDetailList(ctx context.Context, tasks []model.Task) ([]*response.TaskDetailResponse, error) {
	detailList := make([]*response.TaskDetailResponse, 0)
	if len(tasks) == 0 {
		return detailList, nil
	}
	for _, task := range tasks {
		t := task
		instances, err := service.GetInstancesByTaskId(ctx, cast.ToString(task.Id), task.TaskAction)
		if err != nil {
			return nil, err
		}
		r := ConvertToTaskDetail(instances, &t)
		if r != nil {
			detailList = append(detailList, r)
		}
	}

	return detailList, nil
}
