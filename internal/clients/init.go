package clients

import (
	"fmt"

	"github.com/galaxy-future/BridgX/config"
	"github.com/galaxy-future/BridgX/internal/logs"
)

func Init() {
	InitDBClients()
	InitializeClient(fmt.Sprintf("http://localhost:%v", config.GlobalConfig.ServerPort))
	if config.GlobalConfig.EtcdConfig != nil {
		_, err := NewEtcdClient(config.GlobalConfig.EtcdConfig)
		if err != nil {
			panic(err)
		}
	} else {
		logs.Logger.Warn("no ETCD Config has been defined, you should run application in standalone mod.")
	}
}
