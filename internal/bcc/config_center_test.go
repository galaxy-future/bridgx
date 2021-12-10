package bcc

import (
	"testing"
	"time"

	"github.com/galaxy-future/BridgX/config"
)

func TestGetConfig(t *testing.T) {
	conf := config.Config{
		EtcdConfig: &config.EtcdConfig{
			Endpoints:   []string{"127.0.0.1:2379"},
			DailTimeout: time.Second * 5,
		},
	}
	MustInit(&conf)
	testGroup := "group1"
	testDataId := "data1"
	testContent := "content1"

	err := PublishConfig(testGroup, testDataId, testContent)
	if err != nil {
		t.Errorf("failed to run PublishConfig with empty etcd config, err : %v", err)
	}
	gotContent, err := GetConfig(testGroup, testDataId)
	if err != nil {
		t.Errorf("failed to run GetConfig with empty etcd config, err : %v", err)
	}
	if gotContent != testContent {
		t.Errorf("un equal got :%s want:%s", gotContent, testContent)
	}
}
