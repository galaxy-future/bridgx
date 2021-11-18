package tests

import (
	"testing"
	"time"

	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/pool"
	"github.com/galaxy-future/BridgX/pkg/id_generator"
	"github.com/galaxy-future/BridgX/pkg/utils"
	jsoniter "github.com/json-iterator/go"

	"github.com/galaxy-future/BridgX/internal/model"
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
