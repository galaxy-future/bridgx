package main

import (
	"fmt"

	"github.com/galaxy-future/BridgX/cmd/api/middleware"
	"github.com/galaxy-future/BridgX/cmd/api/routers"
	"github.com/galaxy-future/BridgX/config"
	"github.com/galaxy-future/BridgX/internal/bcc"
	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/service"
)

func main() {
	config.Init()
	logs.Init()
	clients.Init()
	if err := bcc.Init(config.GlobalConfig); err != nil {
		panic(err)
	}
	service.Init(100)
	middleware.Init()
	r := routers.Init()
	err := r.Run(fmt.Sprintf(":%d", config.GlobalConfig.ServerPort))
	if err != nil {
		logs.Logger.Fatal(err.Error())
	}
}
