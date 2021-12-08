package cluster_builder

import (
	"fmt"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"testing"
)

func TestParseInitResult(t *testing.T) {
	master, node := parseInitResult(testResult)
	fmt.Println(master)
	fmt.Println(node)
}

func TestCreateCluster(t *testing.T) {
	list := []gf_cluster.ClusterBuildMachine{
		{
			IP:       "192.168.16.84",
			Hostname: "i-2zehhkzw8hsbloo1htvb",
			Username: "root",
			Password: "mzQ7chN2",
			Labels: map[string]string{
				"just": "test",
			},
		},
		{
			IP:       "192.168.16.83",
			Hostname: "i-2zehhkzw8hsbloo1htva",
			Username: "root",
			Password: "mzQ7chN2",
			Labels: map[string]string{
				"just": "test",
			},
		},
	}

	params := gf_cluster.ClusterBuilderParams{
		AccessKey:    "LTAI5t9h55p5b23qHiJ1vxTx",
		AccessSecret: "W6rcn1boYTCJlYMj2mM1MeBNCOHMEg",
		PodCidr:      "10.10.0.0/16",
		SvcCidr:      "10.20.0.0/16",
		MachineList:  list,
		Mode:         gf_cluster.SingleMode,
		KubernetesId: 1,
	}
	CreateCluster(params)
}

func TestPop(t *testing.T) {
	m := gf_cluster.ClusterBuildMachine{}
	list := []gf_cluster.ClusterBuildMachine{
		m, m, m, m, m,
	}

	one, list := Pop(list)
	t.Log("one", one, list, len(list))

	two, list := Pop(list)
	t.Log("two", two, list, len(list))

	three, list := Pop(list)
	t.Log("three", three, list, len(list))

	four, list := Pop(list)
	t.Log("four", four, list, len(list))

	five, list := Pop(list)
	t.Log("five", five, list, len(list))
}
