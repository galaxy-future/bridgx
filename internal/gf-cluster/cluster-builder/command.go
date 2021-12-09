package cluster_builder

import (
	"bytes"
	"fmt"

	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"golang.org/x/crypto/ssh"
)

func CreateCluster(params gf_cluster.ClusterBuilderParams) {
	updateInstallStep(params.KubernetesId, gf_cluster.KubernetesStepInitializeCluster)
	machineList := params.MachineList
	master, machineList := Pop(machineList)

	updateStatus(params.KubernetesId, gf_cluster.KubernetesStatusInitializing)

	initOutput, err := initCluster(master, params.PodCidr, params.SvcCidr)
	recordStep(params.KubernetesId, master.IP, gf_cluster.KubernetesStepInitializeCluster, err)
	if err != nil {
		failed(params.KubernetesId, err.Error())
		return
	}

	masterCmd, nodeCmd := parseInitResult(initOutput)

	//获取kube config
	config, err := initKubeConfig(master)
	if err != nil {
		failed(params.KubernetesId, err.Error())
		return
	}

	if err = recordConfig(params.KubernetesId, config); err != nil {
		failed(params.KubernetesId, "record config err:"+err.Error())
		return
	}

	//安装flannel
	updateInstallStep(params.KubernetesId, gf_cluster.KubernetesStepInstallFlannel)
	err = initFlannel(master, FlannelData{
		AccessKey:    params.AccessKey,
		AccessSecret: params.AccessSecret,
		PodCidr:      params.PodCidr,
	})
	recordStep(params.KubernetesId, master.IP, gf_cluster.KubernetesStepInstallFlannel, err)
	if err != nil {
		failed(params.KubernetesId, "flannel init err:"+err.Error())
		return
	}

	taintMaster(master, master.Hostname)
	//安装master节点
	if params.Mode == gf_cluster.ClusterMode {
		for i := 0; i < 2; i++ {
			var masterNode gf_cluster.ClusterBuildMachine
			masterNode, machineList = Pop(machineList)
			updateInstallStep(params.KubernetesId, gf_cluster.KubernetesStepInstallMaster)
			resetMachine(masterNode)
			_, err = Run(masterNode, masterCmd)
			recordStep(params.KubernetesId, masterNode.IP, gf_cluster.KubernetesStepInstallMaster+masterNode.Hostname, err)
			if err != nil {
				failed(params.KubernetesId, "add master err:"+err.Error())
				return
			}
			taintMaster(master, masterNode.Hostname)
		}
	}

	//安装node节点
	length := len(machineList)
	for i := 0; i < length; i++ {
		var node gf_cluster.ClusterBuildMachine
		node, machineList = Pop(machineList)
		updateInstallStep(params.KubernetesId, gf_cluster.KubernetesStepInstallNode)
		resetMachine(node)
		_, err = Run(node, nodeCmd)
		recordStep(params.KubernetesId, node.IP, gf_cluster.KubernetesStepInstallNode+node.Hostname, err)
		if err != nil {
			failed(params.KubernetesId, "add node err:"+err.Error())
			return
		}
	}

	//给节点打标签
	labelCluster(master, params.MachineList)

	updateInstallStep(params.KubernetesId, gf_cluster.KubernetesStepDone)
	recordStep(params.KubernetesId, "", gf_cluster.KubernetesStepDone, nil)
	updateStatus(params.KubernetesId, gf_cluster.KubernetesStatusRunning)
}

func Run(machine gf_cluster.ClusterBuildMachine, cmd string) (string, error) {
	fmt.Println(cmd)
	client, err := ssh.Dial("tcp", machine.IP+":22", &ssh.ClientConfig{
		User:            machine.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(machine.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return "", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err = session.Run(cmd); err != nil {
		return "", err
	}

	fmt.Println(b.String())
	return b.String(), nil
}
