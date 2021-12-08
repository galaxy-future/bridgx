package cluster_builder

import (
	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/logs"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"go.uber.org/zap"
)

func updateStatus(id int64, status string) error {
	kubernetes := gf_cluster.KubernetesInfo{
		Id:     id,
		Status: status,
	}

	return update(kubernetes)
}

func updateInstallStep(id int64, step string) error {
	kubernetes := gf_cluster.KubernetesInfo{
		Id:          id,
		InstallStep: step,
	}

	return update(kubernetes)
}

func recordConfig(id int64, config string) error {
	kubernetes := gf_cluster.KubernetesInfo{
		Id:     id,
		Config: config,
	}

	return update(kubernetes)
}

func failed(id int64, message string) error {
	logs.Logger.Errorf("cluster create failed.", zap.Int64("id", id), zap.String("message", message))
	kubernetes := gf_cluster.KubernetesInfo{
		Id:      id,
		Status:  gf_cluster.KubernetesStatusFailed,
		Message: message,
	}

	return update(kubernetes)
}

func update(kubernetes gf_cluster.KubernetesInfo) error {
	connection := clients.WriteDBCli
	tx := connection.Model(&gf_cluster.KubernetesInfo{}).Where("id = ?", kubernetes.Id)
	kubernetes.Id = 0
	if err := tx.Updates(kubernetes).Error; err != nil {
		logs.Logger.Errorf("UpdateCluster from WriteDBCli db", err)
		return err
	}

	return nil
}
