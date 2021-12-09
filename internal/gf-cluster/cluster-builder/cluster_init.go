package cluster_builder

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
)

type InitClusterData struct {
	IP      string
	PodCidr string
	SvcCidr string
}

type FlannelData struct {
	PodCidr      string
	AccessKey    string
	AccessSecret string
}

func initClusterTmpl(data InitClusterData) (string, error) {
	tmpl, err := template.New("init_cluster").Parse(initClusterCmd)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func initFlannel(machine gf_cluster.ClusterBuildMachine, data FlannelData) error {
	tmpl, err := template.New("flannel").Parse(flannel)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return err
	}

	_, err = Run(machine, "tee flannel.yaml <<"+buf.String())
	if err != nil {
		return err
	}

	command := `kubectl apply -f flannel.yaml`
	_, err = Run(machine, command)
	return err
}

func initCluster(machine gf_cluster.ClusterBuildMachine, podCidr, svcCidr string) (string, error) {
	resetMachine(machine)

	data := InitClusterData{
		IP:      machine.IP,
		PodCidr: podCidr,
		SvcCidr: svcCidr,
	}

	initCommand, err := initClusterTmpl(data)
	if err != nil {
		return "", err
	}

	result, err := Run(machine, initCommand)
	if err != nil {
		return "", err
	}

	return result, nil
}

func initKubeConfig(machine gf_cluster.ClusterBuildMachine) (string, error) {
	result, err := Run(machine, "mkdir -p $HOME/.kube && sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && cat $HOME/.kube/config")
	if err != nil {
		return "", err
	}

	if strings.Contains(result, "kubernetes-admin") && strings.Contains(result, "client-certificate-data") {
		return result, nil
	} else {
		return "", errors.New("config format is wrong")
	}
}

func parseInitResult(result string) (master string, node string) {
	reg, _ := regexp.Compile("kubeadm join.*\n.*\n.*")
	res := reg.FindAllString(result, 2)
	master = strings.Trim(res[0], " ")
	node = strings.Trim(res[1], " ")
	return
}

func labelCluster(master gf_cluster.ClusterBuildMachine, list []gf_cluster.ClusterBuildMachine) {
	for _, machine := range list {
		for label, values := range machine.Labels {
			cmd := fmt.Sprintf("kubectl label nodes %s %s=%s", convertHostName(machine.Hostname), label, values)
			_, _ = Run(master, cmd)
		}
	}
}

func resetMachine(machine gf_cluster.ClusterBuildMachine) {
	initMachine(machine)
	_, _ = Run(machine, "echo y | kubeadm reset")
	_, _ = Run(machine, "rm -rf .kube & rm flannel.yaml")
	resetFlannel(machine)
}

func taintMaster(master gf_cluster.ClusterBuildMachine, node string) {
	cmd := fmt.Sprintf("kubectl taint nodes %s node-role.kubernetes.io/master:NoSchedule-",
		convertHostName(node))
	_, _ = Run(master, cmd)
}

func initMachine(machine gf_cluster.ClusterBuildMachine) {
	result, err := Run(machine, "ls -lah")
	if err != nil {
		return
	}

	if strings.Contains(result, "init.sh") {
		return
	}

	_, _ = Run(machine, "tee init.sh <<"+initConfig)
	_, _ = Run(machine, "sh init.sh")
}

func resetFlannel(machine gf_cluster.ClusterBuildMachine) {
	cmd := "ifconfig cni0 down && ip link delete cni0 && ifconfig flannel.1 down && ip link delete flannel.1 && rm -rf /var/lib/cni/ && rm -f /etc/cni/net.d/* && systemctl stop kubelet"
	_, _ = Run(machine, cmd)
}
