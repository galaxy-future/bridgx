package tests

import (
	"context"
	"testing"
	"time"

	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/model"
	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestCreateInstance(t *testing.T) {
	instances := make([]model.Instance, 0)
	now := time.Now()
	instances = append(instances, model.Instance{
		InstanceId:  "test1",
		Status:      constants.Pending,
		ClusterName: "test_cluster1",
		CreateAt:    &now,
		DeleteAt:    &now,
	})
	instances = append(instances, model.Instance{
		InstanceId:  "test2",
		Status:      constants.Pending,
		ClusterName: "test_cluster2",
		CreateAt:    &now,
		DeleteAt:    &now,
	})
	err := model.BatchCreateInstance(instances)
	t.Log(err)
}

func TestUpdateInstance(t *testing.T) {
	instance1 := model.Instance{
		InstanceId:  "test1",
		Status:      constants.Running,
		ClusterName: "test_cluster1",
		IpInner:     "10.0.0.1",
	}
	err := model.UpdateByInstanceId(instance1)
	t.Log(err)

	instance2 := model.Instance{
		InstanceId: "test2",
		Status:     constants.Deleted,
	}
	err = model.UpdateByInstanceId(instance2)
	t.Log(err)
}

func TestGetInstanceByIp(t *testing.T) {
	instance, err := model.GetInstanceByIpInner("10.0.0.1")
	t.Log(err)
	t.Log(instance)
}

func TestGetInstanceByIps(t *testing.T) {
	ips := []string{"10.0.0.1", "10.0.0.2"}
	instances, err := model.GetInstancesByIPs(ips, "")
	t.Log(err)
	t.Log(instances)
}

func TestBatchUpdate(t *testing.T) {
	instanceIds := []string{"test1", "test2"}
	err := model.BatchUpdateByInstanceIds(instanceIds, model.Instance{Status: constants.Deleted})
	t.Log(err)
}

func TestSyncInstanceExpireTime(t *testing.T) {
	err := service.SyncInstanceExpireTime(context.Background(), "gf.cloud_4")
	assert.Nil(t, err)
}
