package model

import (
	"github.com/galaxy-future/BridgX/internal/clients"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
)

func ListClusterInstallStep(clusterId int64) (result []gf_cluster.KubernetesInstallStep, err error) {
	err = clients.ReadDBCli.Model(gf_cluster.KubernetesInstallStep{}).Where("cluster_id = ?", clusterId).
		Find(&result).Error
	return
}
