package tests

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/galaxy-future/BridgX/internal/service"
	"github.com/galaxy-future/BridgX/internal/types"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/cloud/alibaba"
	"github.com/galaxy-future/BridgX/pkg/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func TestCreateCluster(t *testing.T) {
	err := service.CreateCluster4Test("TEST_CLUSTER")
	assert.Nil(t, err, "should be nil", err)
}

func TestCreateClusterByApi(t *testing.T) {
	cluster := types.ClusterInfo{
		Name: fmt.Sprintf("gf.metrics.pi"),
		//Name:         fmt.Sprintf("gf.metrics.pi.cluster-%v", time.Now().Unix()),
		RegionId: "cn-qingdao",
		//ZoneId:       "cn-beijing-h",
		InstanceType: "ecs.g6.large",
		Image:        "centos_8_4_uefi_x64_20G_alibase_20210611.vhd",
		Provider:     alibaba.CloudName,
		Password:     "ASDqwe123",
		AccountKey:   "LTAI5tAwAMpXAQ78pePcRb6t",
		NetworkConfig: &types.NetworkConfig{
			Vpc:           "vpc-2zelmmlfd5c5duibc2xb2",
			SubnetId:      "vsw-2zennaxawzq6sa2fdj8l5",
			SecurityGroup: "sg-2zefbt9tw0yo1r7vc3ac",
		},
		StorageConfig: &types.StorageConfig{
			Disks: &cloud.Disks{
				SystemDisk: cloud.DiskConf{Size: 40, Category: "cloud_efficiency"},
				DataDisk: []cloud.DiskConf{{
					Size:     100,
					Category: "cloud_efficiency",
				}},
			},
		},
		ChargeConfig: &types.ChargeConfig{ChargeType: "PostPaid"},
	}
	b, _ := jsoniter.MarshalToString(cluster)
	t.Logf(b)
	ret, _ := utils.HttpPostJsonDataT("http://0.0.0.0:9090/api/v1/cluster/create", []byte(b), 3)
	t.Logf("Response:%v", string(ret))
}

func TestCreateClusterTagsByApi(t *testing.T) {
	name := time.Now().Unix()
	req := fmt.Sprintf(`{"name":"Cluster-%v","desc":"k","region_id":"cn-bj","zone_id":"cn-bj-h","instance_type":"2c4g","charge_type":"by_month","network_config":{"vpc":"vpc-ikw1swp1"},"storage_config":{"mountPoint":"/opt/data","nas":""},"tags":{"dc":"lf","env":"prod"}}`, name)
	ret, err := utils.HttpPostJsonDataT("http://0.0.0.0:9090/api/v1/cluster/create", []byte(req), 3)
	t.Logf("Response:%v", string(ret))
	assert.Nil(t, err, "err not nil")
	tagReq := fmt.Sprintf(`{"cluster_name": "Cluster-%v", "tags": {"k1": "v1", "k2": "v2"}}`, name)
	ret2, err := utils.HttpPostJsonDataT("http://0.0.0.0:9090/api/v1/cluster/add_tags", []byte(tagReq), 3)
	t.Logf("Response:%v", string(ret2))
	assert.Nil(t, err, "err not nil")
}

func TestCreateClusterErr(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = service.CreateCluster4Test("TEST_CLUSTER")
		r := 100 + rand.Int31n(50)
		time.Sleep(time.Duration(r) * time.Millisecond)
	}
}

func TestExpandClusterUseMockCluster(t *testing.T) {
	cluster := types.ClusterInfo{
		Name:         fmt.Sprintf("cluster-%v", time.Now()),
		RegionId:     "cn-beijing",
		ZoneId:       "cn-beijing-h",
		InstanceType: "ecs.s6-c1m1.small",
		Image:        "centos_7_9_x64_20G_alibase_20210623.vhd",
		Provider:     alibaba.CloudName,
		Password:     "ASDqwe123",
		AccountKey:   "LTAI5tAwAMpXAQ78pePcRb6t",
		NetworkConfig: &types.NetworkConfig{
			Vpc:           "vpc-2zelmmlfd5c5duibc2xb2",
			SubnetId:      "vsw-2zennaxawzq6sa2fdj8l5",
			SecurityGroup: "sg-2zefbt9tw0yo1r7vc3ac",
		},
		StorageConfig: &types.StorageConfig{
			Disks: &cloud.Disks{
				SystemDisk: cloud.DiskConf{Size: 40, Category: "cloud_efficiency"},
				DataDisk: []cloud.DiskConf{{
					Size:     100,
					Category: "cloud_efficiency",
				}},
			},
		},
		ChargeConfig: &types.ChargeConfig{ChargeType: "PostPaid"},
	}
	instanceIds, err := service.Expand(&cluster, nil, 2)
	t.Logf("instanceIds: %v", strings.Join(instanceIds, ","))
	t.Log("err: ", err)
}

func TestGetInstance(t *testing.T) {
	cluster := types.ClusterInfo{
		RegionId:   "cn-beijing",
		Provider:   alibaba.CloudName,
		AccountKey: "LTAI5tAwAMpXAQ78pePcRb6t",
	}
	res, err := service.GetInstances(&cluster, []string{"i-2ze5ysm1hx7o9q3mz218", "i-2ze5ysm1hx7o9q3mz219"})
	t.Logf("infos: %v", res)
	t.Log("err: ", err)
}

func TestShrink(t *testing.T) {
	cluster := types.ClusterInfo{
		RegionId:   "cn-beijing",
		Provider:   alibaba.CloudName,
		AccountKey: "LTAI5tAwAMpXAQ78pePcRb6t",
	}
	err := service.Shrink(&cluster, []string{"i-2ze5ysm1hx7o9q3mz218", "i-2ze5ysm1hx7o9q3mz219"})
	t.Log("err: ", err)
}

func TestCreateExpandTask(t *testing.T) {
	req := fmt.Sprintf(`{"cluster_name":"gf.bridgx.online", "count": 1}`)
	ret, err := utils.HttpPostJsonDataT("http://10.192.219.2:9090/api/v1/cluster/expand", []byte(req), 3)
	t.Logf("Response:%v", string(ret))
	assert.Nil(t, err, "err not nil")
}

func TestCreateShrinkTask(t *testing.T) {
	req := fmt.Sprintf(`{"cluster_name":"gf.bridgx.online", "ips":[], "count": 2}`)
	ret, err := utils.HttpPostJsonDataT("http://10.192.219.2:9090/api/v1/cluster/shrink", []byte(req), 3)
	t.Logf("Response:%v", string(ret))
	assert.Nil(t, err, "err not nil")
}

func TestGetClusterCount(t *testing.T) {
	cnt, err := service.GetClusterCount(context.Background(), []string{"LTAI5tAwAMpXAQ78pePcRb6t"})
	t.Logf("get account cnt:%v", cnt)
	assert.Nil(t, err)
	assert.NotZero(t, cnt)
	cnt, err = service.GetClusterCount(context.Background(), []string{"account_not_exist"})
	assert.Nil(t, err)
	assert.Zero(t, cnt)
}
