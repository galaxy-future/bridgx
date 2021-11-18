package constants

import (
	"fmt"
	"time"
)

const DefaultTaskMonitorInterval = 5
const DefaultClusterMonitorInterval = 30
const DefaultTFailedTaskInterval = 15
const DefaultInstanceCountWatcherInterval = 10
const DefaultKillExpireRunningTaskInterval = 10
const DefaultInstanceCleanerRunningInterval = 600
const DefaultQueryOrderInterval = 300
const DefaultTaskMaxRunningDuration = 20 * time.Minute

//DefaultCleanMaxRunningTTL 默认清理任务最大执行时间（秒）
const DefaultCleanMaxRunningTTL = 30

const TaskMonitorETCDLockKeyPrefix = "bridgx/task/locks/"
const ClusterMonitorETCDLockKeyPrefix = "bridgx/cluster/locks/"
const ClusterMonitorETCDReviewKeyPrefix = "bridgx/cluster/reviews/"
const ClusterInstancesCountWatcherETCDReviewKeyPrefix = "bridgx/cluster/instance-count-watcher/"

//GetClusterScheduleLockKey 对于Cluster调度任务/执行任务时 需要获取锁的key
func GetClusterScheduleLockKey(clusterName string) string {
	return fmt.Sprintf("%v/%v", ClusterInstancesCountWatcherETCDReviewKeyPrefix, clusterName)
}
