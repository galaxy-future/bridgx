package clients

import (
	"github.com/galaxy-future/BridgX/config"
	"github.com/galaxy-future/BridgX/internal/logs"
)

func MustInit() {
	InitDBClients()
	if config.GlobalConfig.EtcdConfig != nil {
		_, err := NewEtcdClient(config.GlobalConfig.EtcdConfig)
		if err != nil {
			panic(err)
		}
	} else {
		logs.Logger.Warn("no ETCD Config has been defined, you should run application in standalone mod.")
	}
}
