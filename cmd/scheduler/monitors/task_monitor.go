package monitors

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/galaxy-future/BridgX/internal/bcc"
	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/pool"
	"go.etcd.io/etcd/client/v3/concurrency"
	"gorm.io/gorm"
)

type TaskMonitor struct {
}

func (m TaskMonitor) Run() {
	tasks := make([]model.Task, 0)

	err := model.QueryAll(map[string]interface{}{"status": constants.TaskStatusInit}, &tasks, "")
	if errors.Is(err, gorm.ErrRecordNotFound) || len(tasks) == 0 {
		return
	}

	for _, task := range tasks {
		//err := tryRunTask(task)
		err := bcc.SyncRun(constants.DefaultTaskMonitorInterval, constants.TaskMonitorETCDLockKeyPrefix+strconv.FormatInt(task.Id, 10), func() error {
			return scheduleTask(task)
		})
		if err != nil && err != clients.ErrReviewFailed && err != concurrency.ErrLocked {
			logs.Logger.Errorf("failed to run cluster monitor, err: %v", err)
			continue
		}
	}
}

func scheduleTask(task model.Task) error {
	//ji检查任务是否已经被执行过了
	var newTask model.Task
	err := model.Get(task.Id, &newTask)
	if err != nil {
		return err
	}
	if task.Status != newTask.Status {
		return clients.ErrReviewFailed
	}

	//执行任务
	task.Status = constants.TaskStatusRunning
	err = model.Save(&task)
	if err != nil {
		return err
	}
	switch task.TaskAction {
	case constants.TaskActionExpand:
		pool.ExpandTasksChan <- &task
	case constants.TaskActionShrink:
		pool.ShrinkTasksChan <- &task
	default:
		return fmt.Errorf("unknown task action, action : %v", task.TaskAction)
	}
	return nil
}
