# Containerd install on ubuntu 20.04

> [참고] https://fwany0708.tistory.com/13

Containerd는 kubernetes, docker에서 사용되는 컨테이너 런타임 중 하나임
오픈소스이며 단독적으로 사용할 때는 사용하기 힘들다. 그래서 docker, kubernetes같은 시스템을 같이 사용해주는 것이 좋다. 
root계정으로 설치하는게 마음 편함.

```sh
$ VERSION=1.6.18

$ wget https://github.com/containerd/containerd/releases/download/v${VERSION}/cri-containerd-cni-${VERSION}-linux-amd64.tar.gz

## / 디렉토리에 압축을 해제 하면 cri/cni 설치 된다.
$ tar --no-overwrite-dir -C / -xzf cri-containerd-cni-${VERSION}-linux-amd64.tar.gz
etc/
etc/systemd/
etc/systemd/system/
etc/systemd/system/containerd.service
etc/crictl.yaml
etc/cni/
etc/cni/net.d/
etc/cni/net.d/10-containerd-net.conflist
usr/
usr/local/
usr/local/bin/
usr/local/bin/containerd-shim-runc-v2
usr/local/bin/ctd-decoder
usr/local/bin/containerd-stress
usr/local/bin/ctr
usr/local/bin/critest
usr/local/bin/containerd
usr/local/bin/containerd-shim
usr/local/bin/crictl
usr/local/bin/containerd-shim-runc-v1
usr/local/sbin/
usr/local/sbin/runc
opt/
opt/containerd/
opt/containerd/cluster/
opt/containerd/cluster/gce/
opt/containerd/cluster/gce/cloud-init/
opt/containerd/cluster/gce/cloud-init/master.yaml
opt/containerd/cluster/gce/cloud-init/node.yaml
opt/containerd/cluster/gce/env
opt/containerd/cluster/gce/configure.sh
opt/containerd/cluster/gce/cni.template
opt/containerd/cluster/version
opt/cni/
opt/cni/bin/
opt/cni/bin/sbr
opt/cni/bin/vlan
opt/cni/bin/host-device
opt/cni/bin/tuning
opt/cni/bin/bandwidth
opt/cni/bin/firewall
opt/cni/bin/ipvlan
opt/cni/bin/macvlan
opt/cni/bin/bridge

$ cat <<EOF | sudo tee /etc/modules-load.d/containerd.conf
overlay
br_netfilter
EOF

$ sudo modprobe overlay
$ sudo modprobe br_netfilter

## 재부팅하지 않고 sysctl 파라미터 적용
$ sudo sysctl --system

## containerd 구성
$ sudo mkdir /etc/containerd
$ containerd config default | sudo tee /etc/containerd/config.toml

## SystemdCgroup을 활성화
$ sudo sed -i 's/SystemdCgroup \= false/SystemdCgroup \= true/g' /etc/containerd/config.toml

## Containerd 서비스를 시작하고

$ sudo systemctl daemon-reload
$ sudo systemctl systemctl enable --now containerd

## nerdctl cli 설치
$ wget https://github.com/containerd/nerdctl/releases/download/v1.4.0/nerdctl-1.4.0-linux-amd64.tar.gz
$ tar -C /usr/bin -zxvf ./nerdctl-1.4.0-linux-amd64.tar.gz nerdctl
```