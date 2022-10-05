# Cubescripts

It contains terraform and ansible scripts for installing k8s cluster and addons using cube command 

## Features
- Automates the provisioning of Kubernetes clusters in AWS, Azure and On-promise
- Deploys Single or Highly Available (HA) Kubernetes
- Upgrade Kubernetes cluster
- Add/Delete worker node
- Encrypt sensitive data using ansible-vault

## Supported Components
* Core
  - kubernetes v1.9.10
  - etcd v3.2.18
  - docker v17.03.2.ce
* Network Plugin
  - weave v2.1.3
* Core Addons
  - coredns v1.0.1
  - nginx-ingress-controller v0.14.0
  - defaultbackend v0.4  
  - nfs-client-provisioner v2.0.1  
* Monitoring
  - kubernetes-dashboard v1.10.0
  - prometheus v2.3.2
  - alertmanager v0.15.0
  - grafana v5.0.4  
  - node-exporter v0.16.0
  - kube-state-metrics v1.3.1
  - cocktail-monitoring v1.0.0.B000002
  - disk-usage-exporter v1.0.B000067
  
## Documentation

Documentation is in the `/docs` directory, 

## Changelog from 2.5.4 to 2.5.5
* Restrict supporting k8s version
	* 1.9.1 ~ 1.9.11
	* 1.10.1 ~ 1.10.5
* Add default values for registry and cert values for build server in all.yaml
* Has been Deployed to cengjiyun (master branch)
* Acloud hosting (develop branch)

## Changelog from 2.0.3 to 2.5.0
* Replace deploying method for k8s control panel from yaml to kubeadm
* Add encrypt sensitive data using ansible-vault
* Support kubernetes upgrade as follow
	* 1.8.13, 1.8.14
	* 1.9.1 ~ 1.9.8
	* 1.10.1 ~ 1.10.3
	* Not allow downgrade
* Status query Etcd, k8s control panel, monitoring and cocktail(if installed) pods
* Disable yum auto-upgrade option except security
* Dedicated minikube