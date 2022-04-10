package cluster_builder

var flannel = `__EOF__
---
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: psp.flannel.unprivileged
  annotations:
    seccomp.security.alpha.kubernetes.io/allowedProfileNames: docker/default
    seccomp.security.alpha.kubernetes.io/defaultProfileName: docker/default
    apparmor.security.beta.kubernetes.io/allowedProfileNames: runtime/default
    apparmor.security.beta.kubernetes.io/defaultProfileName: runtime/default
spec:
  privileged: false
  volumes:
  - configMap
  - secret
  - emptyDir
  - hostPath
  allowedHostPaths:
  - pathPrefix: "/etc/cni/net.d"
  - pathPrefix: "/etc/kube-flannel"
  - pathPrefix: "/run/flannel"
  readOnlyRootFilesystem: false
  # Users and groups
  runAsUser:
    rule: RunAsAny
  supplementalGroups:
    rule: RunAsAny
  fsGroup:
    rule: RunAsAny
  # Privilege Escalation
  allowPrivilegeEscalation: false
  defaultAllowPrivilegeEscalation: false
  # Capabilities
  allowedCapabilities: ['NET_ADMIN', 'NET_RAW']
  defaultAddCapabilities: []
  requiredDropCapabilities: []
  # Host namespaces
  hostPID: false
  hostIPC: false
  hostNetwork: true
  hostPorts:
  - min: 0
    max: 65535
  # SELinux
  seLinux:
    # SELinux is unused in CaaSP
    rule: 'RunAsAny'
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: flannel
rules:
- apiGroups: ['extensions']
  resources: ['podsecuritypolicies']
  verbs: ['use']
  resourceNames: ['psp.flannel.unprivileged']
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - get
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - nodes/status
  verbs:
  - patch
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: flannel
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: flannel
subjects:
- kind: ServiceAccount
  name: flannel
  namespace: kube-system
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: flannel
  namespace: kube-system
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: kube-flannel-cfg
  namespace: kube-system
  labels:
    tier: node
    app: flannel
data:
  cni-conf.json: |
    {
      "name": "cbr0",
      "cniVersion": "0.3.1",
      "plugins": [
        {
          "type": "flannel",
          "delegate": {
            "hairpinMode": true,
            "isDefaultGateway": true
          }
        },
        {
          "type": "portmap",
          "capabilities": {
            "portMappings": true
          }
        }
      ]
    }
  net-conf.json: |
    {
      "Network": "{{.PodCidr}}",
      "Backend": {
        {{if eq .NetMode "vxlan"}}
        "Type": "vxlan"
        {{else}}
        "Type": "ali-vpc"
        {{end}}
      }
    }
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: kube-flannel-ds
  namespace: kube-system
  labels:
    tier: node
    app: flannel
spec:
  selector:
    matchLabels:
      app: flannel
  template:
    metadata:
      labels:
        tier: node
        app: flannel
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: kubernetes.io/os
                operator: In
                values:
                - linux
      hostNetwork: true
      priorityClassName: system-node-critical
      tolerations:
      - operator: Exists
        effect: NoSchedule
      serviceAccountName: flannel
      initContainers:
      - name: install-cni-plugin
        image: rancher/mirrored-flannelcni-flannel-cni-plugin:v1.0.0
        command:
        - cp
        args:
        - -f
        - /flannel
        - /opt/cni/bin/flannel
        volumeMounts:
        - name: cni-plugin
          mountPath: /opt/cni/bin
      - name: install-cni
        image: quay.io/coreos/flannel:v0.15.1
        command:
        - cp
        args:
        - -f
        - /etc/kube-flannel/cni-conf.json
        - /etc/cni/net.d/10-flannel.conflist
        volumeMounts:
        - name: cni
          mountPath: /etc/cni/net.d
        - name: flannel-cfg
          mountPath: /etc/kube-flannel/
      containers:
      - name: kube-flannel
        image: quay.io/coreos/flannel:v0.15.1
        command:
        - /opt/bin/flanneld
        args:
        - --ip-masq
        - --kube-subnet-mgr
        resources:
          requests:
            cpu: "100m"
            memory: "50Mi"
          limits:
            cpu: "100m"
            memory: "50Mi"
        securityContext:
          privileged: false
          capabilities:
            add: ["NET_ADMIN", "NET_RAW"]
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: ACCESS_KEY_ID
          value: {{.AccessKey}}
        - name: ACCESS_KEY_SECRET
          value: {{.AccessSecret}}
        volumeMounts:
        - name: run
          mountPath: /run/flannel
        - name: flannel-cfg
          mountPath: /etc/kube-flannel/
      volumes:
      - name: run
        hostPath:
          path: /run/flannel
      - name: cni-plugin
        hostPath:
          path: /opt/cni/bin
      - name: cni
        hostPath:
          path: /etc/cni/net.d
      - name: flannel-cfg
        configMap:
          name: kube-flannel-cfg
__EOF__
`
var initClusterCmd = `
kubeadm init \
        --apiserver-advertise-address 0.0.0.0 \
        --apiserver-bind-port 6443 \
        --cert-dir /etc/kubernetes/pki \
        --control-plane-endpoint {{.IP}} \
        --image-repository registry.cn-hangzhou.aliyuncs.com/google_containers \
        --pod-network-cidr {{.PodCidr}} \
        --service-cidr {{.SvcCidr}} \
        --service-dns-domain cluster.local \
        --upload-certs
`
var initConfig = `'__EOF__'
#!/usr/bin/env bash

#判断系统是否为Linux,内核版本是否满足要求
check_kernel() {
  kernel_name=$(uname -s)
  if [[ "$kernel_name" == "Linux" ]]; then
    printf "\033[32mStep1: Check Kernel Version\033[0m\n"
  else
    printf "\033[31mERROR: BridgeX Only Supports Linux\033[0m\n"
    exit 1
  fi

  #https://docs.docker.com/engine/install/binaries/
  req_kernel="3.10.0"
  cur_kernel=$(uname -r | awk -F- '{print $1}')
  OLD_IFS="$IFS"
  IFS="."
  ck_array=($cur_kernel)
  rk_array=($req_kernel)
  ck_array_len=${#ck_array[*]}
  if [[ "$req_kernel" == "$cur_kernel" ]]; then
    printf "\033[32mStep1.1: Kernel Version is OK\033[0m\n"
  else
    for ((i = 0; i < $ck_array_len; i++)); do
      if [[ "${ck_array[$i]}" -gt "${rk_array[$i]}" ]]; then
        printf "\033[32mStep1.1: Kernel Version is OK\033[0m\n"
        break
      elif [[ "${ck_array[$i]}" -lt "${rk_array[$i]}" ]]; then
        printf "\033[31mERROR: Recommend Kernel Version 3.10 or higher of the Linux kernel\033[0m\n"
        exit 1
      fi
    done
  fi
}

#判断Linux发行版,确定包管理工具
check_os() {
  rhl_file="/etc/redhat-release"
  if [ -f $rhl_file ]; then
    distro_linux=rhl
    printf "\033[32mStep2: Distro Linux is $distro_linux\033[0m\n"
  elif [ -f /etc/os-release ]; then
    distro_linux=$(. /etc/os-release && printf '%s' "$ID")
    printf "\033[32mStep2: distro_linux is $distro_linux\033[0m\n"
  fi
  case $distro_linux in
  rhl)
    util=yum
    ;;
  [Dd]ebian | [Uu]buntu)
    util=apt-get
    ;;
  *)
    cat 1>&2 <<'EOF'
printf "\033[31mError: Currently only rhl and deb systems are supported !!!\033[0m\n"
EOF
    exit 1
    ;;
  esac
}

#云厂商镜像地址
aliyun_mirror="mirrors.aliyun.com"
# tencent_mirror="mirrors.cloud.tencent.com"
# huawei_mirror="repo.huaweicloud.com"

#debian系列发行版初始化docker,kubernetes环境
init_k8s_deb() {
  #安装依赖,添加Docker仓库,安装Docker-ce
  printf "\033[32mStep3: Installing Docker-CE\033[0m\n"
  sudo $util update >/dev/null
  sudo $util install -y apt-transport-https ca-certificates curl gnupg2 lsb-release software-properties-common >/dev/null
  curl -fsSL https://"$aliyun_mirror"/docker-ce/linux/$distro_linux/gpg | sudo gpg --yes --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] \
  https://$aliyun_mirror/docker-ce/linux/$distro_linux $(lsb_release -cs) stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list >/dev/null
  sudo $util update >/dev/null &&\
  sudo $util -y install docker-ce docker-ce-cli containerd.io &&\
  sudo apt-mark hold docker-ce docker-ce-cli containerd.io >/dev/null &&\
  #添加kubernetes仓库,安装kubelet kubeadm kubectl
  printf "\033[32mStep4: Installing kubelet kubeadm kubectl\033[0m\n"
  curl -fsSL https://"$aliyun_mirror"/kubernetes/apt/doc/apt-key.gpg | sudo gpg --yes --dearmor -o /usr/share/keyrings/kubernetes-archive-keyring.gpg
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] \
  https://$aliyun_mirror/kubernetes/apt/ kubernetes-xenial main" | \
  sudo tee /etc/apt/sources.list.d/kubernetes.list >/dev/null
  $util update >/dev/null && $util install -y kubelet kubeadm kubectl >/dev/null && apt-mark hold kubelet kubeadm kubectl
}

#rhl系列发行版初始化docker,kubernetes环境
init_k8s_rhl() {
  #安装依赖,添加Docker仓库,安装Docker-ce
  printf "\033[32mStep3: Installing Docker-CE\033[0m\n"
  sudo $util -y -q makecache >/dev/null && sudo $util install -y -q yum-utils >/dev/null
  sudo yum-config-manager --add-repo https://"$aliyun_mirror"/docker-ce/linux/centos/docker-ce.repo >/dev/null
  sudo sed -i 's+download.docker.com+m$aliyun_mirror/docker-ce+' /etc/yum.repos.d/docker-ce.repo
  sudo $util -y -q makecache >/dev/null && sudo $util -y -q install docker-ce docker-ce-cli containerd.io >/dev/null
  sudo setenforce 0 && sudo sed -i 's/^SELINUX=enforcing$/SELINUX=disabled/' /etc/selinux/config
  #添加kubernetes仓库,安装kubelet kubeadm kubectl
  printf "\033[32mStep4: Installing kubelet kubeadm kubectl\033[0m\n"
  cat >/etc/yum.repos.d/kubernetes.repo <<EOF
[kubernetes]
name=Kubernetes
baseurl=https://$aliyun_mirror/kubernetes/yum/repos/kubernetes-el7-\$basearch/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://$aliyun_mirror/kubernetes/yum/doc/yum-key.gpg https://$aliyun_mirror/kubernetes/yum/doc/rpm-package-key.gpg
exclude=kubernetes
EOF
  sudo $util -y -q --nogpgcheck makecache >/dev/null && sudo $util install -y -q --nogpgcheck kubeadm kubelet kubectl >/dev/null
}

#初始化调用
init_k8s() {
  status=0
  [ ! -d /etc/docker ] && mkdir /etc/docker && touch /etc/docker/daemon.json
  cat >/etc/docker/daemon.json <<EOF
{
    "exec-opts": ["native.cgroupdriver=systemd"],
    "log-driver": "json-file",
    "log-opts": {
      "max-size": "300m",
      "max-file": "10"
    },
    "storage-driver": "overlay2",
    "ip": "0.0.0.0",
    "selinux-enabled": false,
    "registry-mirrors":[
        "https://kfwkfulq.mirror.aliyuncs.com",
        "https://2lqq34jg.mirror.aliyuncs.com",
        "https://pee6w651.mirror.aliyuncs.com",
        "https://docker.mirrors.ustc.edu.cn",
        "http://hub-mirror.c.163.com"
    ]
}
EOF
  sudo modprobe br_netfilter
  cat >/etc/modules-load.d/k8s.conf <<EOF
br_netfilter
EOF
  cat >/etc/sysctl.d/k8s.conf <<EOF
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
  sudo sysctl --system >/dev/null
  sudo echo 1 >/proc/sys/net/bridge/bridge-nf-call-iptables

  if [[ "$util" == "yum" ]]; then
    init_k8s_rhl
  elif [[ "$util" == "apt-get" ]]; then
    init_k8s_deb
  fi
  sudo systemctl daemon-reload &&\
  sudo systemctl enable docker --now && printf "\033[32mStep3.1: Docker Install Succeed\033[0m\n" &&\
  sudo systemctl enable kubelet --now && printf "\033[32mStep4.1: Kubelet Install Succeed\033[0m\n"
}

#阶段执行入口
main() {
  check_kernel
  check_os
  init_k8s
}

main
exit $status
__EOF__
`
