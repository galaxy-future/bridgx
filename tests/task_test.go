package tests

import (
	"testing"
	"time"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/pool"
	"github.com/galaxy-future/BridgX/pkg/id_generator"
	"github.com/galaxy-future/BridgX/pkg/utils"
	jsoniter "github.com/json-iterator/go"
)

func TestCountByTaskStatus(t *testing.T) {
	cnt, _ := model.CountByTaskStatus("gf.metrics.pi.cluster-1634627868", []string{"SUCCESS"})
	t.Logf("task count:%v", cnt)
}

func TestTaskExpand(t *testing.T) {
	info := &model.ExpandTaskInfo{
		ClusterName:    "gf.bridgx.online",
		Count:          200,
		TaskSubmitHost: utils.PrivateIPv4(),
	}
	s, _ := jsoniter.MarshalToString(info)
	taskId := id_generator.GetNextId()
	task := &model.Task{
		TaskName:      "yulong_test",
		TaskAction:    constants.TaskActionExpand,
		Status:        constants.TaskStatusInit,
		TaskFilter:    "gf.bridgx.online",
		TaskInfo:      s,
		SupportCancel: false,
	}
	now := time.Now()
	task.Id = int64(taskId)
	task.CreateAt = &now
	task.UpdateAt = &now
	pool.DoExpand(task)
}

func Test_getTaskInfoCountDiff(t *testing.T) {
	success := 1
	type args struct {
		task *model.Task
	}
	tests := []struct {
		name         string
		args         args
		wantBefore   int
		wantExpect   int
		wantAfter    int
		wantUserName string
	}{
		{
			name: "expand instance count diff",
			args: args{
				task: &model.Task{
					TaskAction: constants.TaskActionExpand,
					TaskInfo:   "{\"count\":5,\"before_count\":10,\"user_id\":1}",
				},
			},
			wantBefore:   10,
			wantExpect:   15,
			wantAfter:    11,
			wantUserName: "root",
		},
		{
			name: "shrink instance count diff",
			args: args{
				task: &model.Task{
					TaskAction: constants.TaskActionShrink,
					TaskInfo:   "{\"count\":6,\"before_count\":10,\"user_id\":1}",
				},
			},
			wantBefore:   10,
			wantExpect:   4,
			wantAfter:    9,
			wantUserName: "root",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := helper.ExtractTaskInfo(tt.args.task)
			if info.GetBeforeInstanceCount() != tt.wantBefore {
				t.Errorf("extractTaskInfo() gotBefore = %v, want %v", info.GetBeforeInstanceCount(), tt.wantBefore)
			}
			if info.GetExpectInstanceCount() != tt.wantExpect {
				t.Errorf("extractTaskInfo() gotExpect = %v, want %v", info.GetExpectInstanceCount(), tt.wantExpect)
			}
			if info.GetAfterInstanceCount(success) != tt.wantAfter {
				t.Errorf("extractTaskInfo() gotAfter = %v, want %v", info.GetAfterInstanceCount(1), tt.wantAfter)
			}
			if info.GetCreateUsername() != tt.wantUserName {
				t.Errorf("extractTaskInfo() GetCreateUsername = %v, want %v", info.GetCreateUsername(), tt.wantUserName)
			}
		})
	}
}
