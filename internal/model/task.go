package model

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/spf13/cast"
)

type Task struct {
	Base
	TaskName      string     `json:"task_name"`
	Status        string     `json:"status"`      //INIT, RUNNING, SUCCESS, FAILED
	TaskAction    string     `json:"task_action"` //expand, shrink
	TaskFilter    string     `json:"task_filter"` //任务过滤，业务标识（如集群名等）
	TaskInfo      string     `json:"task_info"`   //不同任务需要的不同的参数
	ErrMsg        string     `json:"err_msg"`
	TaskResult    string     `json:"task_result"`
	SupportCancel bool       `json:"support_cancel"`
	FinishTime    *time.Time `json:"finish_time"`
}

func (Task) TableName() string {
	return "task"
}

type TaskInfo interface {
	GetBeforeAndExpectInstanceCount() (beforeCount int, expectCount int)
}

type ExpandTaskInfo struct {
	ClusterName    string `json:"cluster_name"`
	Count          int    `json:"count"`
	TaskExecHost   string `json:"task_exec_host"`
	TaskSubmitHost string `json:"task_submit_host"`
	UserId         int64  `json:"user_id"`
	BeforeCount    int    `json:"before_count"`
}

func (e *ExpandTaskInfo) GetBeforeAndExpectInstanceCount() (beforeCount int, expectCount int) {
	return e.BeforeCount, e.BeforeCount + e.Count
}

type ExpandTaskRes struct {
	InstanceIdList []string `json:"instance_id_list"`
}

type ShrinkTaskInfo struct {
	ClusterName    string `json:"cluster_name"`
	Count          int    `json:"count"`
	IPs            string `json:"ips"`
	TaskExecHost   string `json:"task_exec_host"`
	TaskSubmitHost string `json:"task_submit_host"`
	UserId         int64  `json:"user_id"`
	BeforeCount    int    `json:"before_count"`
}

func (e *ShrinkTaskInfo) GetBeforeAndExpectInstanceCount() (beforeCount int, expectCount int) {
	return e.BeforeCount, e.BeforeCount - e.Count
}

func CountByTaskStatus(taskFilter string, statuses []string) (int64, error) {
	var cnt int64
	if err := clients.ReadDBCli.Model(&Task{}).Where("task_filter = ? AND status IN (?)", taskFilter, statuses).Count(&cnt).Error; err != nil {
		logErr("CountByTaskStatus from read db", err)
		return 0, err
	}
	return cnt, nil
}

//GetTaskByStatus 获取当前cluster下状态zai在列表中的所有实例
func GetTaskByStatus(clusterName string, statuses []string) ([]Task, error) {
	var tasks []Task
	if err := clients.ReadDBCli.Where("task_filter = ? AND status in (?) ", clusterName, statuses).Find(&tasks).Error; err != nil {
		logErr("GetActiveTaskByClusterName from read db", err)
		return tasks, err
	}
	return tasks, nil
}

//GetExpireRunningTask 获取执行状态为Running并且最后更新时间（应该为执行时间）大于指定时间的所有的task
func GetExpireRunningTask(duration time.Duration) ([]Task, error) {
	var tasks []Task
	if err := clients.ReadDBCli.Where("status = ? AND update_at < ?  ", constants.TaskStatusRunning, time.Now().Add(-duration)).Find(&tasks).Error; err != nil {
		logErr("GetExpireRunningTask from read db", err)
		return tasks, err
	}
	return tasks, nil
}

func UpdateTaskStatus(taskIds []int64, status string) error {
	now := time.Now()
	err := Updates(Task{}, taskIds, map[string]interface{}{"status": status, "update_at": &now})
	if err != nil {
		logErr("GetExpireRunningTask from read db", err)
		return fmt.Errorf("can not update task stauts : %w", err)
	}
	return nil
}

func GetTaskCount(ctx context.Context, clusterNames []string) (int64, error) {
	var cnt int64
	if err := clients.ReadDBCli.WithContext(ctx).Model(&Task{}).Where("task_filter IN (?) ", clusterNames).Count(&cnt).Error; err != nil {
		logErr("GetTaskCount from read db", err)
		return 0, err
	}
	return cnt, nil
}

type TaskSearchCond struct {
	TaskId      string
	TaskName    string
	ClusterName string
	TaskStatus  string
	PageNumber  int
	PageSize    int
}

func SearchTask(ctx context.Context, clusterNames []string, cond TaskSearchCond) (ret []Task, total int64, err error) {
	query := clients.ReadDBCli.WithContext(ctx).Table(Task{}.TaskName).Where("task_filter IN (?) ", clusterNames)
	if cond.TaskId != "" {
		query = query.Where("id = ? ", cast.ToInt64(cond.TaskId))
	}
	if cond.TaskName != "" {
		query = query.Where("task_name LIKE ?", "%"+cond.TaskName+"%")
	}
	if cond.TaskStatus != "" {
		query = query.Where("status = ?", cond.TaskStatus)
	}
	total, err = QueryWhere(query, cond.PageNumber, cond.PageSize, &ret, "id DESC", true)
	if err != nil {
		return nil, 0, err
	}
	return ret, total, nil
}
