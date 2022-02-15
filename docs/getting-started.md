- [Deploy kubernetes cluster using kore-on container](#deploy-kubernetes-cluster-using-kore-on-container)
  - [Create VM with vagrant and virtualbox](#create-vm-with-vagrant-and-virtualbox)
  - [페쇄망을 위한 local-repo, registry 백업 파일 만들기](#페쇄망을-위한-local-repo-registry-백업-파일-만들기)
  - [Deploy k8s cluster to VM](#deploy-k8s-cluster-to-vm)

# Deploy kubernetes cluster using kore-on container
## Create VM with vagrant and virtualbox

 [vagrant를 이용한 vm 생성](vagrant-virtualbox.md)
 
## 페쇄망을 위한 local-repo, registry 백업 파일 만들기 

 [local-repo, harbor 백업파일 만들](prepare-airgap.md)
  
## Deploy k8s cluster to VM
1. Make clean working directory.
```bash
$ mkdir /tmp/koreon; cd /tmp/koreon
```

2. Generate ssh key and copy it to target servers
```bash
$ ssh-keygen -t rsa -N '' -f ./id_rsa
$ ssh-copy-id ./id_rsa root@192.168.77.190
$ ssh-copy-id ./id_rsa root@192.168.77.191
$ ssh-copy-id ./id_rsa root@192.168.77.192
$ ssh-copy-id ./id_rsa root@192.168.77.193
$ ssh-copy-id ./id_rsa root@192.168.77.194
```

3. Run container
  - Inventory 및 클러스터 설치 관련 변수 설정
mycluster
├── group_vars
│   └── all
│       ├── basic.yml
│       └── expert.yml
└── inventory.ini

  - Copy sample inventory files and modify it according to your environment
```bash
$ docker run -it --name=koreon --rm -v ${PWD}:/koreon/work regi.k3.acornsoft.io/k3lab/koreon:1.1.1 cp -R /koreon/inventory/sample mycluster
$ cp local-repo.20210726_120901.tgz mycluster
$ cp harbor.20210726_122024.tgz mycluster
```

  - Edit inventory file and configuration files
```bash
$ vi mycluster/inventory.ini
# Inventory sample
[all]
master-01   ansible_ssh_host=192.168.77.190  ip=192.168.77.190
master-02   ansible_ssh_host=192.168.77.191  ip=192.168.77.191
master-03   ansible_ssh_host=192.168.77.192  ip=192.168.77.192
etcd-01     ansible_ssh_host=192.168.77.190  ip=192.168.77.190
etcd-02     ansible_ssh_host=192.168.77.191  ip=192.168.77.191
etcd-03     ansible_ssh_host=192.168.77.192  ip=192.168.77.192
node-01     ansible_ssh_host=192.168.77.193  ip=192.168.77.193
storage-01  ansible_ssh_host=192.168.77.194  ip=192.168.77.194
registry-01 ansible_ssh_host=192.168.77.194  ip=192.168.77.194

[etcd]
etcd-01
etcd-02
etcd-03

[masters]
master-01
master-02
master-03

[sslhost]
master-01

[node]
node-01

[storage]
storage-01

[registry]
registry-01

[cluster:children]
masters
node


$ vi mycluster/group_vars/all/basic.yml
provider: false
cloud_provider: onpremise
cluster_name: test-cluster

# install directories
install_dir: /var/lib/koreon
data_root_dir: /data

# kubernetes options
k8s_version: 1.21.1
cluster_id: test-cluster
api_lb_ip: https://192.168.77.190:6443
lb_ip: 192.168.77.190
lb_port: 6443
pod_ip_range: 10.0.0.0/16
service_ip_range: 172.20.0.0/16

# for air gap installation
closed_network: true
local_repository: http://192.168.77.194:8080
local_repository_archieve_file: local-repo.20210726_120901.tgz

# option for master isolation
master_isolated: false
audit_log_enable: true
cert_validity_days: 36500

# container runtime [containerd | docker]
container_runtime: containerd

# kube-proxy mode [iptables | ipvs]
kube_proxy_mode: ipvs

# option for harbor registry
registry_install: true
registry_data_dir: /data/harbor
registry: 192.168.77.194
registry_domain: 192.168.77.194
registry_public_cert: false
registry_archieve_file: harbor.20210726_122024.tgz

# option for NFS storage
storage_install: true
nfs_ip: 192.168.77.194
nfs_volume_dir: /storage

# for internal load-balancer
haproxy: true
```

  - cluster install
<p align="center">
  <a href="https://asciinema.org/a/MkTRg4ke0RFo6T9lBOwv79Z9G">
  <img src="https://asciinema.org/a/MkTRg4ke0RFo6T9lBOwv79Z9G.png" width="885"></image>
  </a>
</p>
        
<p align="center">
  <a href="https://asciinema.org/a/XnnbMUFCiSdCJUMPc64wAoA8U">
  <img src="https://asciinema.org/a/XnnbMUFCiSdCJUMPc64wAoA8U.png" width="885"></image>
  </a>
</p>
   
      
```bash
$ docker run -it --name=koreon --rm -v ${PWD}:/koreon/work regi.k3.acornsoft.io/k3lab/koreon:1.1.1 /bin/bash

# 대상 장비 접속 여부 확인하기(필수)
$ ansible -i mycluster/inventory.ini -u root --private-key id_rsa  all -m ping

# 대상 장비 VM 파라미터 등 확인(옵션)
$ ansible -i mycluster/inventory.ini -u root --private-key id_rsa  all -m setup

# 클러스터 설치하기
$ ansible-playbook -i mycluster/inventory.ini -u root --private-key id_rsa ../scripts/cluster.yml

# 클러스터 설치 후 확인
$ kubectl --kubeconfig=mycluster/acloud-client-kubeconfig get nodes -o wide
NAME             STATUS   ROLES                  AGE    VERSION   INTERNAL-IP      EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION     CONTAINER-RUNTIME
ubuntu2004-190   Ready    control-plane,master   104m   v1.21.1   192.168.77.190   <none>        Ubuntu 20.04.2 LTS   5.4.0-74-generic   containerd://1.4.3
ubuntu2004-191   Ready    control-plane,master   104m   v1.21.1   192.168.77.191   <none>        Ubuntu 20.04.2 LTS   5.4.0-74-generic   containerd://1.4.3
ubuntu2004-192   Ready    control-plane,master   104m   v1.21.1   192.168.77.192   <none>        Ubuntu 20.04.2 LTS   5.4.0-74-generic   containerd://1.4.3
ubuntu2004-193   Ready    node                   102m   v1.21.1   192.168.77.193   <none>        Ubuntu 20.04.2 LTS   5.4.0-74-generic   containerd://1.4.3


$ kubectl --kubeconfig=mycluster/acloud-client-kubeconfig get po -A
NAMESPACE     NAME                                       READY   STATUS    RESTARTS   AGE
kube-system   calico-kube-controllers-78d6f96c7b-qdxgh   1/1     Running   0          88m
kube-system   calico-node-fzshs                          1/1     Running   0          88m
kube-system   calico-node-nxblb                          1/1     Running   0          88m
kube-system   calico-node-rksqs                          1/1     Running   0          88m
kube-system   calico-node-sljvm                          1/1     Running   0          87m
kube-system   coredns-589f7445-6f5xt                     1/1     Running   0          89m
kube-system   coredns-589f7445-skk89                     1/1     Running   0          89m
kube-system   haproxy-ubuntu2004-193                     1/1     Running   0          87m
kube-system   kube-apiserver-ubuntu2004-190              1/1     Running   0          89m
kube-system   kube-apiserver-ubuntu2004-191              1/1     Running   0          89m
kube-system   kube-apiserver-ubuntu2004-192              1/1     Running   0          89m
kube-system   kube-controller-manager-ubuntu2004-190     1/1     Running   0          89m
kube-system   kube-controller-manager-ubuntu2004-191     1/1     Running   0          89m
kube-system   kube-controller-manager-ubuntu2004-192     1/1     Running   0          89m
kube-system   kube-proxy-hrztd                           1/1     Running   0          88m
kube-system   kube-proxy-m2gxp                           1/1     Running   0          88m
kube-system   kube-proxy-mgcfd                           1/1     Running   0          88m
kube-system   kube-proxy-rcv8s                           1/1     Running   0          87m
kube-system   kube-scheduler-ubuntu2004-190              1/1     Running   0          89m
kube-system   kube-scheduler-ubuntu2004-191              1/1     Running   0          89m
kube-system   kube-scheduler-ubuntu2004-192              1/1     Running   0          89m
kube-system   metrics-server-65cbb6d659-mdzvr            1/1     Running   0          88m
kube-system   metrics-server-65cbb6d659-n6996            1/1     Running   0          88m
kube-system   nfs-pod-provisioner-7b7944c46d-tgbkx       1/1     Running   0          87m
```

  - 클러스터 Worker node 추가하기
```bash
$ vi mycluster/inventory.ini
# Inventory sample
[all]
master-01   ansible_ssh_host=192.168.77.190  ip=192.168.77.190
master-02   ansible_ssh_host=192.168.77.191  ip=192.168.77.191
master-03   ansible_ssh_host=192.168.77.192  ip=192.168.77.192
etcd-01     ansible_ssh_host=192.168.77.190  ip=192.168.77.190
etcd-02     ansible_ssh_host=192.168.77.191  ip=192.168.77.191
etcd-03     ansible_ssh_host=192.168.77.192  ip=192.168.77.192
node-01     ansible_ssh_host=192.168.77.193  ip=192.168.77.193
node-02     ansible_ssh_host=192.168.77.195  ip=192.168.77.195
storage-01  ansible_ssh_host=192.168.77.194  ip=192.168.77.194
registry-01 ansible_ssh_host=192.168.77.194  ip=192.168.77.194

[etcd]
etcd-01
etcd-02
etcd-03

[masters]
master-01
master-02
master-03

[sslhost]
master-01

[node]
node-01
node-02

[storage]
storage-01

[registry]
registry-01

[cluster:children]
masters
node

$ ansible-playbook -i mycluster/inventory.ini -u root --private-key id_rsa ../scripts/add-node.yml
```  

  - 클러스터 업그레이드 하기
```bash
$ vi mycluster/group_vars/all/basic.yml
...
k8s_version: 1.21.2
...

$ ansible-playbook -i mycluster/inventory.ini -u root --private-key id_rsa ../scripts/upgrade.yml
```  

  - 클러스터 Worker node 삭제하기
```bash
$ ansible-playbook -i mycluster/inventory.ini -u root --private-key id_rsa -e remove_node_name=ubuntu2004-195 -e target=192.168.77.195 remove-node.yml
```

  - 클러스터 삭제하기
```bash
# optional [reset-cluster/reset-registry/reset-storage]
$ ansible-playbook -i mycluster/inventory.ini -u root --private-key id_rsa --tags reset-cluster ../scripts/reset.yml
```