# 페쇄망을 위한 local-repo, registry 백업 파일 만들기 
  
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
  - Copy sample inventory files and modify it according to your environment
```bash
$ docker run -it --name=knit --rm -v ${PWD}:/knit/work regi.k3.acornsoft.io/k3lab/knit:1.1.0 cp -R /knit/inventory/sample mycluster
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
```

  - Edit basic.yml, expert.yml
```bash
$ vi mycluster/group_vars/all/basic.yml
provider: false
cloud_provider: onpremise
cluster_name: test-cluster

# install directories
install_dir: /var/lib/knit
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
closed_network: false      <---- (주의) 압축파일 생성시에는 internet 가능한 환경임.
local_repository: http://192.168.77.194:8080
local_repository_archive_file:

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
registry_archive_file:

# option for NFS storage
storage_install: true
nfs_ip: 192.168.77.194
nfs_volume_dir: /storage

# for internal load-balancer
haproxy: true


$ vi mycluster/group_vars/all/expert.yml
...
archive_repo: true
...
```
     
  - ansible-playbook 실행
```bash
$ docker run -it --name=knit --rm -v ${PWD}:/knit/work regi.k3.acornsoft.io/k3lab/knit:1.1.0 /bin/bash

# 대상 장비 접속 여부 확인하기(필수)
$ ansible -i mycluster/inventory.ini -u root --private-key id_rsa  all -m ping

# 대상 장비 VM 파라미터 등 확인(옵션)
$ ansible -i mycluster/inventory.ini -u root --private-key id_rsa  all -m setup

# 클러스터 설치하기
$ ansible-playbook -i mycluster/inventory.ini -u root --private-key id_rsa ../scripts/prepare-repository.yml
```

  - 압축 파일 가져오기
```bash
$ scp root@192.168.77.194/tmp/*.tgz .
```  