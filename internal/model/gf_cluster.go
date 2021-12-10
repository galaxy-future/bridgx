package model

import (
	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"gorm.io/gorm"
)

//RegisterKubernetesCluster  注册集群
func RegisterKubernetesCluster(kubernetes *gf_cluster.KubernetesInfo) error {
	if err := clients.WriteDBCli.Create(kubernetes).Error; err != nil {
		logErr("CreateCluster from WriteDBCli db", err)
		return err
	}
	return nil
}

//DeleteKubernetesCluster 删除集群记录
func DeleteKubernetesCluster(kubernetesId int64) error {
	err := clients.WriteDBCli.Delete(&gf_cluster.KubernetesInfo{}, kubernetesId).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		logErr("DeleteClusterById from WriteDBCli db", err)
		return err
	}
	return nil
}

//UpdateKubernetesCluster 更新集群信息
func UpdateKubernetesCluster(kubernetes *gf_cluster.KubernetesInfo) error {
	if err := clients.WriteDBCli.Save(kubernetes).Error; err != nil {
		logErr("UpdateCluster from WriteDBCli db", err)
		return err
	}
	return nil
}

//GetKubernetesCluster 获取集群
func GetKubernetesCluster(kubernetesId int64) (*gf_cluster.KubernetesInfo, error) {

	var cluster gf_cluster.KubernetesInfo
	if err := clients.ReadDBCli.Where("id = ?", kubernetesId).First(&cluster).Error; err != nil {
		logErr("GetClusterById from read db", err)
		return nil, err
	}
	return &cluster, nil
}

//ListKubernetesClusters 列出所有集群
func ListKubernetesClusters() ([]*gf_cluster.KubernetesInfo, error) {
	var clusters []*gf_cluster.KubernetesInfo
	if err := clients.ReadDBCli.Find(&clusters).Error; err != nil {
		logErr("GetClusterById from read db", err)
		return nil, err
	}
	return clusters, nil
}

//ListRunningKubernetesClusters 列出所有正在运行的集群
func ListRunningKubernetesClusters() ([]*gf_cluster.KubernetesInfo, error) {
	var clusters []*gf_cluster.KubernetesInfo
	if err := clients.ReadDBCli.Where("status = ?", "running").Find(&clusters).Error; err != nil {
		logErr("GetClusterById from read db", err)
		return nil, err
	}
	return clusters, nil
}
