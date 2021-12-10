package bcc

import (
	"github.com/galaxy-future/BridgX/config"
	"github.com/galaxy-future/BridgX/internal/clients"
)

var configCenter ConfigCenter

type ConfigCenter interface {
	GetConfig(group, dataId string) (string, error)
	PublishConfig(group, dataId, content string) error
}

func MustInit(config *config.Config) {
	clt, err := clients.NewEtcdClient(config.EtcdConfig)
	if err != nil {
		panic(err)
	}
	configCenter = clt
}

func GetConfig(group, dataId string) (string, error) {
	return configCenter.GetConfig(group, dataId)
}

func PublishConfig(group, dataId, content string) error {
	return configCenter.PublishConfig(group, dataId, content)
}
