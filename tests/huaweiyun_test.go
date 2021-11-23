package tests

import (
	"testing"

	"github.com/galaxy-future/BridgX/pkg/cloud/huaweiyun"
)

func TestGetHuaweiyunClient(t *testing.T) {
	p := huaweiyun.New("", "", "cn-north-4")
	_, err := p.GetInstances(make([]string, 0))
	t.Logf("err:%v\n", err)
}
