package handler

import "github.com/galaxy-future/BridgX/cmd/api/helper"

func Init() {
	helper.RegisterHandlerLogReader(CreateCluster, new(CreateClusterLogReader))
	helper.RegisterHandlerLogReader(ModifyAdminPassword, new(ModifyAdminPasswordLogReader))
}
