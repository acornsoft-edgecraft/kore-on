# kore-on

kore-on automates k8s installation tasks for on-premise.
It also provides upgrade cluster and scale up/down worker nodes.

## Features
- Deploys Single or Highly Available (HA) Kubernetes
- Upgrade Kubernetes cluster
- Add/Delete worker node
- Install harbor registry
- Install NFS server
- Air-Gap installation

## Documentation

Documentation is in the `/docs` directory

## Supported Linux Distributions

- **Ubuntu** 20.04
- **CentOS/RHEL** 8

## Supported Components

- Core
  - [kubernetes](https://github.com/kubernetes/kubernetes/tree/master/CHANGELOG) v1.20.6-v1.20.8, v1.21.0-v1.21.12, v1.22.1-v1.22.9, v1.23.1-v1.23.5
  - [etcd](https://github.com/etcd-io/etcd/releases) v3.4.16
  - [docker-compose](https://github.com/docker/compose/releases) v1.29.2  
  - [docker](https://www.docker.com/) v19.03.15
  - [containerd](https://containerd.io/) v1.4.3
  - [crictl](https://github.com/kubernetes-sigs/cri-tools) v1.19.0, v1.20.0, v1.21.0
  
- Network Plugin
  - [calico](https://github.com/projectcalico/calico/releases) v3.19.1, v3.20, v3.21, v3.22, v3.23
  
- Application
  - [coredns](https://github.com/coredns/coredns) v1.7.0, v1.8.0
  - [haproxy](https://hub.docker.com/_/haproxy?tab=tags&page=1&ordering=last_updated) v2.4.2  
  - [nfs-subdir-external-provisioner](https://github.com/kubernetes-sigs/nfs-subdir-external-provisioner/releases) v4.0.2  
  
- Registry
  - [harbor](https://github.com/goharbor/harbor/releases) v2.3.0
  - harbor proxy mirroring
  
## Required packages
 * docker runtime

## go-ansible 
Go-ansible is a package for running ansible-playbook or ansible commands from Golang applications.
