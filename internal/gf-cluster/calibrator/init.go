package calibrator

import (
	"github.com/galaxy-future/BridgX/internal/gf-cluster/instance"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"go.uber.org/zap"
	"time"
)

func Init() error {
	ticker := time.NewTicker(time.Minute * 60)
	go func() {
		for {
			<-ticker.C
			err := calibrate()
			if err != nil {
				logs.Logger.Error("failed to calibrate instances info", zap.Error(err))
				continue
			}
		}
	}()
	logs.Logger.Info("bridgx-kubernetes calibrator start success.")
	return nil
}

func calibrate() error {
	instanceGroups, err := model.ListAllInstanceGroupFromDB()
	if err != nil {
		return err
	}
	for _, instanceGroup := range instanceGroups {
		instances, err := instance.ListCustomInstances(instanceGroup.Id)
		if err != nil {
			logs.Logger.Error("failed to list instances from kubernetes.", zap.Int64("instance_group_id", instanceGroup.Id), zap.String("instance_group_name", instanceGroup.Name), zap.Error(err))
			continue
		}
		onlineCount := len(instances)
		if onlineCount == instanceGroup.InstanceCount {
			continue
		}
		if onlineCount < instanceGroup.InstanceCount {
			begin := time.Now()
			err := instance.ExpandCustomInstanceGroup(instanceGroup, instanceGroup.InstanceCount)
			if err != nil {
				logs.Logger.Error("failed to expand instance", zap.Int64("instance_group_id", instanceGroup.Id), zap.String("instance_group_name", instanceGroup.Name), zap.Error(err))
			}
			cost := time.Now().Sub(begin).Milliseconds()
			err = instance.AddInstanceForm(instanceGroup, cost, gf_cluster.ExpandAndShrinkDefaultUserId, gf_cluster.ExpandAndShrinkDefaultUser, gf_cluster.OptTypeExpand, instanceGroup.InstanceCount-onlineCount, err)
			if err != nil {
				logs.Logger.Error("failed to add instance form while calibrate",
					zap.Int64("instance_group_id", instanceGroup.Id), zap.String("instance_group_name", instanceGroup.Name),
					zap.String("opt_type", gf_cluster.OptTypeExpand), zap.Error(err))
				continue
			}
		}
		if onlineCount > instanceGroup.InstanceCount {
			begin := time.Now()
			// replace shrink count
			instanceGroup.InstanceCount = onlineCount
			err := instance.ShrinkCustomInstanceGroup(instanceGroup, instanceGroup.InstanceCount)
			if err != nil {
				logs.Logger.Error("failed to shrink instance", zap.Int64("instance_group_id", instanceGroup.Id), zap.String("instance_group_name", instanceGroup.Name), zap.Error(err))
			}
			cost := time.Now().Sub(begin).Milliseconds()
			err = instance.AddInstanceForm(instanceGroup, cost, gf_cluster.ExpandAndShrinkDefaultUserId, gf_cluster.ExpandAndShrinkDefaultUser, gf_cluster.OptTypeShrink, instanceGroup.InstanceCount-onlineCount, err)
			if err != nil {
				logs.Logger.Error("failed to add instance form while calibrate", zap.Int64("instance_group_id", instanceGroup.Id),
					zap.String("instance_group_name", instanceGroup.Name),
					zap.String("opt_type", gf_cluster.OptTypeShrink), zap.Error(err))
				continue
			}
		}
	}
	return err
}
