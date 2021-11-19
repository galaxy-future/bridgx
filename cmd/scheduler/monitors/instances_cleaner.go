package monitors

import (
	"fmt"

	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/service"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/atomic"
)

//InstanceCleaner 负责定时扫描云厂商由于系统异常创建的实例，并释放这部分实例
type InstanceCleaner struct {
	clusterName  string
	VersionNo    *atomic.String
	LockerClient *clients.EtcdClient
}

func (cleaner *InstanceCleaner) Run() {
	recycleCount := 0
	err := cleaner.LockerClient.SyncRun(constants.DefaultCleanMaxRunningTTL, constants.GetClusterScheduleLockKey(cleaner.clusterName), func() error {
		cluster, err := model.GetByClusterName(cleaner.clusterName)
		if err != nil {
			return err
		}
		//查看是否有正在执行的任务，如果有任务执行，则不进行任何清理
		tasks, err := model.GetTaskByStatus(cleaner.clusterName, []string{constants.TaskStatusInit, constants.TaskStatusRunning})
		if err != nil {
			return err
		}
		if len(tasks) != 0 {
			return clients.ErrReviewFailed
		}

		tags, err := model.GetTagsByClusterName(cluster.ClusterName)
		if err != nil {
			return err
		}

		info, err := service.ConvertToClusterInfo(cluster, tags)
		if err != nil {
			return fmt.Errorf("failed to convert cluster to cluster info , %w", err)
		}
		recycleCount, err = service.CleanClusterUnusedInstances(info)
		return err
	})
	if err != concurrency.ErrLocked {
		logs.Logger.Errorf("failed to clean cluster unused instance err:%v", err)
		return
	}
}

func (cleaner *InstanceCleaner) UniqueKey() string {
	return "cleaner-" + cleaner.clusterName
}

func (cleaner *InstanceCleaner) GetVersionNo() string {
	return cleaner.VersionNo.Load()
}
func (cleaner *InstanceCleaner) SetVersionNo(v string) {
	cleaner.VersionNo.Store(v)
}
