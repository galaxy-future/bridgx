package main

import (
	"github.com/galaxy-future/BridgX/cmd/scheduler/crond"
	"github.com/galaxy-future/BridgX/cmd/scheduler/monitors"
	"github.com/galaxy-future/BridgX/cmd/scheduler/types"
	"github.com/galaxy-future/BridgX/internal/constants"
)

var schedulers = []*types.Scheduler{
	{
		//扫库，查看是否有待执行的Task，分配Task到WorkerPool
		Interval: constants.DefaultTaskMonitorInterval,
		Monitor:  &monitors.TaskMonitor{},
	},
	// 自动监控当前实例数量与预期实例数量是否相等并执行扩缩容，待启用
	{
		Interval: constants.DefaultClusterMonitorInterval,
		Monitor:  &monitors.ClusterMonitor{},
	},
	{
		Interval: constants.DefaultKillExpireRunningTaskInterval,
		Monitor:  &monitors.TaskKiller{},
	},
	{
		Interval: constants.DefaultQueryOrderInterval,
		Monitor:  &monitors.QueryOrderJobs{},
	},
}

func Run() {
	for _, s := range schedulers {
		crond.AddFixedIntervalSecondsJob(s.Interval, s.Monitor)
	}
	crond.Run()
}

func Stop() {
	crond.Stop()
}
