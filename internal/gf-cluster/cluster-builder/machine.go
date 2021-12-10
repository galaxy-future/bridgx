package cluster_builder

import (
	"strings"

	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
)

func Pop(slice []gf_cluster.ClusterBuildMachine) (m gf_cluster.ClusterBuildMachine, list []gf_cluster.ClusterBuildMachine) {
	m = slice[len(slice)-1]
	list = slice[0 : len(slice)-1]
	return
}

func convertHostName(hostName string) string {
	name := strings.Replace(hostName, "-", "z", 1)
	return name + "z"
}
