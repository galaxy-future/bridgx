package handler

import (
	"net/http"
	"strings"

	"github.com/galaxy-future/BridgX/cmd/api/helper"
	"github.com/galaxy-future/BridgX/cmd/api/response"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func GetTaskCount(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	accountKey := ctx.Query("account")
	accountKeys, err := service.GetAksByOrgAk(ctx, user.GetOrgIdForTest(), accountKey)
	cnt, err := service.GetTaskCount(ctx, accountKeys)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp := &response.TaskCountResponse{
		TaskNum: cnt,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func GetTaskDescribe(ctx *gin.Context) {
	taskId, ok := ctx.GetQuery("task_id")
	if !ok {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	task, err := service.GetTask(ctx, taskId)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	instances, err := service.GetInstancesByTaskId(ctx, taskId, task.TaskAction)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp := helper.ConvertToTaskDetail(instances, task)
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func GetTaskDescribeAll(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.PermissionDenied, nil)
		return
	}
	taskId := ctx.Query("task_id")
	taskName := ctx.Query("task_name")
	clusterName := ctx.Query("cluster_name")
	taskStatus := strings.ToUpper(ctx.Query("task_status"))
	pn, ps := getPager(ctx)
	accountKeys, _ := service.GetAksByOrgId(user.OrgId)
	tasks, total, err := service.GetTaskListByCond(ctx, accountKeys, model.TaskSearchCond{
		TaskId:      taskId,
		TaskName:    taskName,
		ClusterName: clusterName,
		TaskStatus:  taskStatus,
		PageNumber:  pn,
		PageSize:    ps,
	})
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	pager := response.Pager{
		PageNumber: pn,
		PageSize:   ps,
		Total:      int(total),
	}
	list, err := helper.ConvertToTaskDetailList(ctx, tasks)
	resp := response.TaskDetailListResponse{
		TaskList: list,
		Pager:    pager,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func GetTaskList(ctx *gin.Context) {
	user := helper.GetUserClaims(ctx)
	if user == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}

	accountKey, ok := ctx.GetQuery("account")
	accountKeys := make([]string, 0)
	if ok && accountKey != "" {
		accountKeys = append(accountKeys, accountKey)
	} else {
		accountKeys, _ = service.GetAksByOrgId(user.OrgId)
	}
	pn, ps := getPager(ctx)

	tasks, total, err := service.GetTaskListByAks(ctx, accountKeys, pn, ps)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	pager := response.Pager{
		PageNumber: pn,
		PageSize:   ps,
		Total:      int(total),
	}
	resp := &response.TaskListResponse{
		TaskList: helper.ConvertToTaskThumbList(tasks),
		Pager:    pager,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}

func GetTaskInstances(ctx *gin.Context) {
	taskId, ok := ctx.GetQuery("task_id")
	if !ok {
		response.MkResponse(ctx, http.StatusBadRequest, response.ParamInvalid, nil)
		return
	}
	instanceStatus := strings.ToUpper(ctx.Query("instance_status"))
	pn, ps := getPager(ctx)
	task, err := service.GetTask(ctx, taskId)
	if err != nil || task == nil {
		response.MkResponse(ctx, http.StatusBadRequest, response.TaskNotFound, nil)
		return
	}
	condition := service.InstancesSearchCond{
		TaskId:     cast.ToInt64(taskId),
		TaskAction: task.TaskAction,
		Status:     instanceStatus,
		PageNumber: pn,
		PageSize:   ps,
	}
	instances, total, err := service.GetInstancesByCond(ctx, condition)
	if err != nil {
		response.MkResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	pager := response.Pager{
		PageNumber: pn,
		PageSize:   ps,
		Total:      int(total),
	}
	cluster, err := service.GetClusterByName(ctx, task.TaskFilter)
	if err != nil {
		response.MkResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp := &response.TaskInstancesResponse{
		InstanceList: helper.ConvertToInstanceThumbList(ctx, instances, []model.Cluster{*cluster}),
		Pager:        pager,
	}
	response.MkResponse(ctx, http.StatusOK, response.Success, resp)
	return
}
