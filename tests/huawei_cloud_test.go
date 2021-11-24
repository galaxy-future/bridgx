package tests

import (
	"testing"

	"github.com/galaxy-future/BridgX/pkg/cloud/huawei"
)

func TestGetHuaweiCloudClient(t *testing.T) {
	p := huawei.New("", "", "cn-north-4")
	_, err := p.GetInstances(make([]string, 0))
	t.Logf("err:%v\n", err)
}
