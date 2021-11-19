package monitors

import (
	"errors"
	"github.com/galaxy-future/BridgX/internal/clients"
	"time"

	"github.com/galaxy-future/BridgX/cmd/scheduler/crond"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/model"
	"go.uber.org/atomic"
	"gorm.io/gorm"
)

var lastUpdateTime = time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC)

//ClusterMonitor 负责发现新建立的cluster，并启动一个定时任务，已监控该集群是否有变更
type ClusterMonitor struct {
	LockerClient *clients.EtcdClient
}

func (m ClusterMonitor) Run() {
	clusters := make([]model.Cluster, 0)
	err := model.QueryAll(map[string]interface{}{}, &clusters, "")
	if errors.Is(err, gorm.ErrRecordNotFound) || len(clusters) == 0 {
		return
	}
	var tmp = lastUpdateTime
	for _, cluster := range clusters {
		if cluster.UpdateAt.After(lastUpdateTime) {

			if cluster.Status != constants.ClusterStatusEnable {
				m.removeClusterMonitorJobs(&cluster)
				continue
			}

			m.addClusterMonitorJobs(&cluster)
			if cluster.UpdateAt.After(tmp) {
				tmp = *cluster.UpdateAt
			}
		}
	}
	lastUpdateTime = tmp
}

func (m ClusterMonitor) addClusterMonitorJobs(cluster *model.Cluster) {
	//instanceCountJob := &InstanceCountWatchJob{
	//	ClusterName: cluster.ClusterName,
	//	VersionNo:   atomic.NewString(""),
	//}
	//crond.AddFixedIntervalSecondsXJob(constants.DefaultInstanceCountWatcherInterval, instanceCountJob)

	cleanerJob := &InstanceCleaner{
		clusterName: cluster.ClusterName,
		VersionNo:   atomic.NewString(""),
	}
	crond.AddFixedIntervalSecondsXJob(constants.DefaultInstanceCleanerRunningInterval, cleanerJob)
}

func (m ClusterMonitor) removeClusterMonitorJobs(cluster *model.Cluster) {
	//instanceCountJob := &InstanceCountWatchJob{
	//	ClusterName: cluster.ClusterName,
	//	VersionNo:   atomic.NewString(""),
	//}
	//crond.RemoveXJob(instanceCountJob.UniqueKey())

	cleanerJob := &InstanceCleaner{
		clusterName:  cluster.ClusterName,
		VersionNo:    atomic.NewString(""),
		LockerClient: m.LockerClient,
	}
	crond.RemoveXJob(cleanerJob.UniqueKey())
}
