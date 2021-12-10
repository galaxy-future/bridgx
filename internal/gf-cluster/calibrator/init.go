package calibrator

import (
	"runtime/debug"
	"time"

	"github.com/galaxy-future/BridgX/internal/gf-cluster/instance"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"go.uber.org/zap"
)

func Init() {
	ticker := time.NewTicker(time.Minute * 10)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Logger.Errorf("calibrate err:%v ", r)
				logs.Logger.Errorf("calibrate panic", zap.String("stack", string(debug.Stack())))
			}
		}()
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
}

func calibrate() error {
	instanceGroups, err := model.ListAllInstanceGroupFromDB()
	if err != nil {
		return err
	}
	for _, instanceGroup := range instanceGroups {
		instances, err := instance.ListCustomInstances(instanceGroup.Id)
		if err != nil {
			logs.Logger.Errorw("Trigger calibration check: failed to list instances from kubernetes.", zap.Int64("instance_group_id", instanceGroup.Id), zap.String("instance_group_name", instanceGroup.Name), zap.Error(err))
			continue
		}
		onlineCount := len(instances)
		if onlineCount == instanceGroup.InstanceCount {
			continue
		}
		if onlineCount < instanceGroup.InstanceCount {
			destCount := instanceGroup.InstanceCount
			instanceGroup.InstanceCount = onlineCount
			err := instance.ExpandCustomInstanceGroup(instanceGroup, destCount)
			if err != nil {
				logs.Logger.Errorw("Trigger calibration check: failed to expand instance", zap.Int64("instance_group_id", instanceGroup.Id), zap.String("instance_group_name", instanceGroup.Name), zap.Error(err))
				continue
			}
			logs.Logger.Infow("Trigger calibration check: success.", zap.Int64("instance_group_id", instanceGroup.Id), zap.String("instance_group_name", instanceGroup.Name), zap.String("opt_type", gf_cluster.OptTypeExpand), zap.Int("updatedInstanceCount", destCount-instanceGroup.InstanceCount))
		}
		if onlineCount > instanceGroup.InstanceCount {
			destCount := instanceGroup.InstanceCount
			instanceGroup.InstanceCount = onlineCount
			err := instance.ShrinkCustomInstanceGroup(instanceGroup, destCount)
			if err != nil {
				logs.Logger.Errorw("Trigger calibration check: failed to shrink instance", zap.Int64("instance_group_id", instanceGroup.Id), zap.String("instance_group_name", instanceGroup.Name), zap.Error(err))
				continue
			}
			logs.Logger.Infow("Trigger calibration check: success.", zap.Int64("instance_group_id", instanceGroup.Id), zap.String("instance_group_name", instanceGroup.Name), zap.String("opt_type", gf_cluster.OptTypeShrink), zap.Int("updatedInstanceCount", instanceGroup.InstanceCount-destCount))
		}
	}
	return err
}
