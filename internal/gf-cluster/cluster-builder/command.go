package cluster_builder

import (
	"bytes"
	"fmt"
	"sync"

	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"golang.org/x/crypto/ssh"
)

func CreateCluster(params gf_cluster.ClusterBuilderParams) {
	var wg sync.WaitGroup
	updateInstallStep(params.KubernetesId, gf_cluster.KubernetesStepInitializeCluster)
	machineList := params.MachineList
	master, machineList := Pop(machineList)
	defer taintMaster(master, master.Hostname)

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
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = initFlannel(master, FlannelData{
			AccessKey:    params.AccessKey,
			AccessSecret: params.AccessSecret,
			PodCidr:      params.PodCidr,
		})
		recordStep(params.KubernetesId, master.IP, gf_cluster.KubernetesStepInstallFlannel, err)
		if err != nil {
			failed(params.KubernetesId, "flannel init err:"+err.Error())
		}
	}()

	//安装master节点
	if params.Mode == gf_cluster.ClusterMode {
		updateInstallStep(params.KubernetesId, gf_cluster.KubernetesStepInstallMaster)
		for i := 0; i < 2; i++ {
			var masterNode gf_cluster.ClusterBuildMachine
			masterNode, machineList = Pop(machineList)
			wg.Add(1)
			go func(masterMachine gf_cluster.ClusterBuildMachine) {
				defer wg.Done()
				resetMachine(masterMachine)
				_, err = sshRun(masterMachine, masterCmd)
				recordStep(params.KubernetesId, masterMachine.IP, gf_cluster.KubernetesStepInstallMaster+masterMachine.Hostname, err)
				if err != nil {
					failed(params.KubernetesId, "add master err:"+err.Error())
				}
				taintMaster(master, masterMachine.Hostname)
			}(masterNode)
		}
	}

	//安装node节点
	length := len(machineList)
	updateInstallStep(params.KubernetesId, gf_cluster.KubernetesStepInstallNode)
	for i := 0; i < length; i++ {
		var node gf_cluster.ClusterBuildMachine
		node, machineList = Pop(machineList)
		wg.Add(1)
		go func(nodeMachine gf_cluster.ClusterBuildMachine) {
			defer wg.Done()
			resetMachine(nodeMachine)
			_, err = sshRun(nodeMachine, nodeCmd)
			recordStep(params.KubernetesId, nodeMachine.IP, gf_cluster.KubernetesStepInstallNode+nodeMachine.Hostname, err)
			if err != nil {
				failed(params.KubernetesId, "add node err:"+err.Error())
			}
		}(node)
	}

	//给节点打标签
	labelCluster(master, params.MachineList)

	wg.Wait()
	updateInstallStep(params.KubernetesId, gf_cluster.KubernetesStepDone)
	recordStep(params.KubernetesId, "", gf_cluster.KubernetesStepDone, nil)
	updateStatus(params.KubernetesId, gf_cluster.KubernetesStatusRunning)
}

func AddMachine(master gf_cluster.ClusterBuildMachine, machine gf_cluster.ClusterBuildMachine) error {
	cmd, err := getJoinCommand(master)
	if err != nil {
		return err
	}

	resetMachine(machine)

	_, err = sshRun(machine, cmd)
	if err != nil {
		return err
	}

	return nil
}

func sshRun(machine gf_cluster.ClusterBuildMachine, cmd string) (string, error) {
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
