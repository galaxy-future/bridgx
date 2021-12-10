package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/galaxy-future/BridgX/cmd/scheduler/crond"
	"github.com/galaxy-future/BridgX/config"
	"github.com/galaxy-future/BridgX/internal/bcc"
	"github.com/galaxy-future/BridgX/internal/cache"
	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/logs"
)

func main() {
	config.MustInit()
	logs.Init()
	clients.MustInit()
	crond.Init()
	bcc.MustInit(config.GlobalConfig)
	cache.MustInit()

	err := Init()
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		Run()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	fmt.Println("terminal with signal, waiting for running tasks stop... ")

	//停止所有task，并等待所有任务执行完成
	Stop()
	wg.Wait()
	fmt.Println("all running tasks stopped, goodbye! ")

}
