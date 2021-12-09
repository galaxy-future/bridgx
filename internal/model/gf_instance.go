package model

import (
	"github.com/galaxy-future/BridgX/internal/clients"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
)

func CreateInstanceGroupFromDB(instanceGroup *gf_cluster.InstanceGroup) error {

	if err := clients.WriteDBCli.Create(instanceGroup).Error; err != nil {
		logErr("CreateInstanceGroupFromDB from write db", err)
		return err
	}
	return nil
}

func DeleteInstanceGroupFromDB(instanceGroupId int64) error {
	if err := clients.WriteDBCli.Delete(&gf_cluster.InstanceGroup{}, instanceGroupId).Error; err != nil {
		logErr("DeleteInstanceGroupFromDB from write db", err)
		return err
	}

	return nil
}

func UpdateInstanceGroupFromDB(instanceGroup *gf_cluster.InstanceGroup) error {
	if err := clients.WriteDBCli.Save(instanceGroup).Error; err != nil {
		logErr("UpdateInstanceGroupFromDB from write db", err)
		return err
	}
	return nil
}

func GetInstanceGroupFromDB(instanceGroupId int64) (*gf_cluster.InstanceGroup, error) {

	var instanceGroup gf_cluster.InstanceGroup
	if err := clients.ReadDBCli.Where("id = ?", instanceGroupId).First(&instanceGroup).Error; err != nil {
		logErr("GetInstanceGroupFromDB from read db", err)
		return nil, err
	}
	return &instanceGroup, nil

}

func ListInstanceGroupFromDB(name string, pageNumber int, pageSize int) ([]*gf_cluster.InstanceGroup, int64, error) {
	clients := clients.ReadDBCli.Model(gf_cluster.InstanceGroup{})
	if name != "" {
		clients.Where("name like ?", "%"+name+"%")
	}

	var clusters []*gf_cluster.InstanceGroup
	if err := clients.Order("id desc").Offset((pageNumber - 1) * pageSize).Limit(pageSize).Find(&clusters).Error; err != nil {
		logErr("ListInstanceGroupFromDB from read db", err)
		return nil, 0, err
	}
	var total int64
	if err := clients.Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		logErr("ListInstanceGroupFromDB from read db", err)
		return nil, 0, err
	}
	return clusters, total, nil
}

func ListAllInstanceGroupFromDB() ([]*gf_cluster.InstanceGroup, error) {
	var clusters []*gf_cluster.InstanceGroup
	if err := clients.ReadDBCli.Find(&clusters).Error; err != nil {
		logErr("ListAllInstanceGroupFromDB from read db", err)
		return nil, err
	}
	return clusters, nil
}

func ListInstanceGroupInKubernetes(kubernetesId int64) ([]*gf_cluster.InstanceGroup, error) {
	var clusters []*gf_cluster.InstanceGroup
	if err := clients.ReadDBCli.Where("kubernetes_id = ?", kubernetesId).Find(&clusters).Error; err != nil {
		logErr("ListInstanceGroupFromDB from read db", err)
		return nil, err
	}

	return clusters, nil
}

func ListInstanceGroupByUser(curUserId string) ([]*gf_cluster.InstanceGroup, error) {
	var clusters []*gf_cluster.InstanceGroup
	if err := clients.ReadDBCli.Where("created_user_id = ?", curUserId).Find(&clusters).Error; err != nil {
		logErr("ListInstanceGroupFromDB from read db", err)
		return nil, err
	}
	return clusters, nil
}

func UpdateInstanceGroupInstanceCountFromDB(count int, id int64) error {
	if err := clients.WriteDBCli.Model(gf_cluster.InstanceGroup{}).Where("id", id).Update("instance_count", count).Error; err != nil {
		logErr("UpdateInstanceGroupInstanceCountFromDB from write db", err)
		return err
	}
	return nil
}

func CreateInstanceFormFromDB(instanceForm *gf_cluster.InstanceForm) error {

	if err := clients.WriteDBCli.Create(instanceForm).Error; err != nil {
		logErr("CreateInstanceFormFromDB from write db", err)
		return err
	}
	return nil
}

func ListInstanceFormFromDB(id string, pageNumber int, pageSize int) ([]*gf_cluster.InstanceForm, int64, error) {
	clients := clients.ReadDBCli.Model(gf_cluster.InstanceForm{})
	if id != "" {
		clients.Where("id = ?", id)
	}

	var forms []*gf_cluster.InstanceForm
	if err := clients.Order("id desc").Offset((pageNumber - 1) * pageSize).Limit(pageSize).Find(&forms).Error; err != nil {
		logErr("ListInstanceFormFromDB from read db", err)
		return nil, 0, err
	}
	var total int64
	if err := clients.Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		logErr("ListInstanceFormFromDB from read db", err)
		return nil, 0, err
	}
	return forms, total, nil
}
