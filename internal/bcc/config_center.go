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

func Init(config *config.Config) error {
	clt, err := clients.NewEtcdClient(config.EtcdConfig)
	if err != nil {
		return err
	}
	configCenter = clt
	return nil
}

func GetConfig(group, dataId string) (string, error) {
	return configCenter.GetConfig(group, dataId)
}

func PublishConfig(group, dataId, content string) error {
	return configCenter.PublishConfig(group, dataId, content)
}
