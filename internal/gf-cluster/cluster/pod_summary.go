package cluster

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

//ListClusterPodsSummary 获得集群下所有Pod详情
func ListClusterPodsSummary(clusterId int64) (gf_cluster.ClusterPodsSummaryArray, error) {
	cluster, err := model.GetKubernetesCluster(clusterId)
	if err != nil {
		return nil, err
	}
	pods, err := getClusterPodInfo(cluster)
	if err != nil {
		return nil, err
	}
	return pods, nil

}

//getClusterPodInfo 获取集群pod信息，当前使用clientset直接查询kubernetes集群，有性能压力
//TODO 使用client-to watcher/informer机制，缓存对象，防止频繁获取信息对与k8s集群的压力
func getClusterPodInfo(info *gf_cluster.KubernetesInfo) ([]*gf_cluster.PodSummary, error) {

	if info.Status != gf_cluster.KubernetesStatusRunning {
		return make([]*gf_cluster.PodSummary, 0), nil
	}
	client, err := GetKubeClient(info.Id)
	if err != nil {
		return nil, err
	}
	//TODO remove hadrd code
	ctxPodQuery, cancelPodQuery := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelPodQuery()
	selector := metav1.LabelSelector{MatchLabels: createInstanceLabels()}
	pods, err := client.ClientSet.CoreV1().Pods("default").List(ctxPodQuery, metav1.ListOptions{
		LabelSelector: labels.Set(selector.MatchLabels).String(),
	},
	)
	if err != nil {
		return nil, fmt.Errorf("从集群获取pod信息时出错，出错原因为:%s", err.Error())
	}

	var podInfos []*gf_cluster.PodSummary
	for _, pod := range pods.Items {

		cpuResource := resource.NewScaledQuantity(0, resource.Kilo)
		memoryResource := resource.NewScaledQuantity(0, resource.Kilo)
		storageResource := resource.NewScaledQuantity(0, resource.Kilo)
		if pod.Status.Phase == v1.PodRunning {
			for _, container := range pod.Spec.Containers {
				cpuResource.Add(*container.Resources.Limits.Cpu())
				memoryResource.Add(*container.Resources.Limits.Memory())
				storageResource.Add(*container.Resources.Limits.StorageEphemeral())
			}
		}

		groupName, exist := pod.Labels[gf_cluster.ClusterInstanceGroupKey]
		if !exist {
			groupName = "unknown"
		}

		groupIdStr, exist := pod.Labels[gf_cluster.ClusterInstanceGroupIdKey]
		if !exist {
			groupIdStr = "0"
		}
		groupId, err := strconv.ParseInt(groupIdStr, 10, 64)
		if err != nil {
			logs.Logger.Error("获取GroupId失败", zap.String("value", cpuResource.String()), zap.Error(err))
		}
		cpuUsed := cpuQuantity2Float(*cpuResource)
		memoryUsed := storageQuantity2Float(*memoryResource)
		storageUSed := storageQuantity2Float(*storageResource)

		status := pod.Status.Phase

		startTime := ""
		if pod.Status.StartTime != nil {
			startTime = formatHumanReadableDuration(time.Now().Sub(pod.Status.StartTime.Time))
		}

		podInfos = append(podInfos, &gf_cluster.PodSummary{
			NodeName:          pod.Spec.NodeName,
			NodeIp:            pod.Status.HostIP,
			PodName:           pod.GetName(),
			PodIP:             pod.Status.PodIP,
			AllocatedCpuCores: cpuUsed,
			AllocatedMemoryGi: memoryUsed,
			AllocatedDiskGi:   storageUSed,
			GroupName:         groupName,
			RunningTime:       startTime,
			Status:            string(status),
			GroupId:           groupId,
		})
	}

	return podInfos, nil

}

func formatHumanReadableDuration(duration time.Duration) string {
	duration = duration.Round(time.Minute)
	d := duration / (time.Hour * 24)
	duration -= d * time.Hour * 24
	h := duration / time.Hour
	duration -= h * time.Hour
	m := duration / time.Minute
	if d > 0 {
		return fmt.Sprintf("%04d天%02d小时%02d分", d, h, m)
	}
	if h > 0 {
		return fmt.Sprintf("%02d小时%02d分", h, m)
	}
	return fmt.Sprintf("%02d分", m)
}
