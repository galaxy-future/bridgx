package cluster

import (
	"fmt"

	"github.com/galaxy-future/BridgX/internal/model"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	clients map[int64]*KubernetesClient
)

type KubernetesClient struct {
	ClientSet *kubernetes.Clientset
}

func CreateKubernetesClusterClient(data []byte) (*KubernetesClient, error) {
	config, err := clientcmd.NewClientConfigFromBytes(data)
	if err != nil {
		return nil, err
	}
	clientConfig, err := config.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		panic(err.Error())
	}

	return &KubernetesClient{ClientSet: clientSet}, nil
}

func GetKubeClient(kubernetesClusterId int64) (*KubernetesClient, error) {
	client := clients[kubernetesClusterId]
	if client == nil {
		kubernetesCluster, err := model.GetKubernetesCluster(kubernetesClusterId)
		if err != nil {
			return nil, err
		}
		if kubernetesCluster.Status != gf_cluster.KubernetesStatusRunning {
			return nil, fmt.Errorf("当前集群不可用, 当前状态为:%s", kubernetesCluster.Status)
		}
		client, err = CreateKubernetesClusterClient([]byte(kubernetesCluster.Config))
		if err != nil {
			return nil, fmt.Errorf("创建集群连接失败,错误原因: %s", err.Error())
		}
		clients[kubernetesClusterId] = client
	}
	return client, nil
}

func init() {
	clients = make(map[int64]*KubernetesClient)
}
