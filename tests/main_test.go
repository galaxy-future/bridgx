package tests

import (
	"os"
	"testing"

	"github.com/galaxy-future/BridgX/config"
	"github.com/galaxy-future/BridgX/internal/cache"
	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/logs"
)

func TestMain(m *testing.M) {
	//因为是相对路径，需要把conf文件copy到tests目录下
	config.MustInit()
	logs.Init()
	clients.MustInit()
	cache.MustInit()
	exitCode := m.Run()
	os.Exit(exitCode)
}
