package monitors

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/service"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/atomic"
)

//InstanceCountWatchJob 负责监控一个cluster是否有变更，如果有变更会schedule一个任务，保证需求可以满足
type InstanceCountWatchJob struct {
	ClusterName  string
	VersionNo    *atomic.String
	LockerClient *clients.EtcdClient
	sync.Mutex
}

func (m *InstanceCountWatchJob) Run() {
	syncKey := constants.GetClusterScheduleLockKey(m.ClusterName)
	err := m.LockerClient.SyncRun(constants.DefaultInstanceCountWatcherInterval, syncKey, func() error {
		return scheduleJob(m.ClusterName)
	})
	if err != nil && err != concurrency.ErrLocked {
		logs.Logger.Errorf("failed to watch cluster count err: %v", err)
	}
}

func (m *InstanceCountWatchJob) UniqueKey() string {
	return "instance-count-watch-" + m.ClusterName
}

func (m *InstanceCountWatchJob) GetVersionNo() string {
	return m.VersionNo.Load()
}
func (m *InstanceCountWatchJob) SetVersionNo(v string) {
	m.VersionNo.Store(v)
}

func calcWorkingCount(workingIPs string) int {
	if workingIPs == "" || workingIPs == constants.HasNoneIP {
		return 0
	}
	return len(strings.Split(workingIPs, ","))
}

func genVersionNo(e, w int) string {
	return fmt.Sprintf("%v-%v-%v", e, w, time.Now().Minute())
}

func scheduleJob(clusterName string) error {
	snapshot, err := model.GetClusterSnapshot(clusterName)
	if err != nil {
		return err
	}
	//如果存在任务，或者集群不需要调度不需要调度任务
	if len(snapshot.RunningTask) != 0 || snapshot.Cluster.ExpectCount == len(snapshot.ActiveInstances) {
		return nil
	}

	//扩容
	if snapshot.Cluster.ExpectCount > len(snapshot.ActiveInstances) {
		_, err := service.CreateExpandTask(context.Background(), snapshot.Cluster.ClusterName, snapshot.Cluster.ExpectCount-len(snapshot.ActiveInstances), "EXPECT", 0)
		if err != nil {
			logs.Logger.Errorf("CreateExpandTask err:%v", err)
			return err
		}
	}
	//缩容
	if snapshot.Cluster.ExpectCount < len(snapshot.ActiveInstances) {
		var deleteIPs []string
		for _, instance := range snapshot.ActiveInstances {
			if instance.Status == constants.Deleting {
				deleteIPs = append(deleteIPs, instance.IpInner)
			}
		}
		if len(snapshot.ActiveInstances)-len(deleteIPs) != snapshot.Cluster.ExpectCount {
			return fmt.Errorf("can not schedule shrink task because expect count != instance count - deleting instance count")
		}
		_, err := service.CreateShrinkTask(context.Background(), snapshot.Cluster.ClusterName, len(deleteIPs), strings.Join(deleteIPs, ","), "EXPECT", 0)
		if err != nil {
			logs.Logger.Errorf("CreateShrinkTask err:%v", err)
			return err
		}
	}
	return nil
}
