package handler

import "github.com/galaxy-future/BridgX/cmd/api/helper"

func Init() {
	helper.RegisterHandlerLogReader(ExpandCluster, new(ExpandClusterLogReader))
	helper.RegisterHandlerLogReader(ModifyAdminPassword, new(ModifyAdminPasswordLogReader))
}
