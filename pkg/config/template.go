package config

const Template = `
[koreon]
## Required
## - local-repository-install: local repository installation activate. (Required when selecting the closed network.)
##                             It is installed on the registry host.
## - local-repository-port: Port number used as local repository. (Required when selecting the closed network.)
##                          If you use the default value, you can skip it. (default: 8080)
## - local-repository-archive-file: local repository packages archive file path (Required when selecting the closed network.)
## - local-repository-url: local repository service url (Required when selecting the closed network.)
##                         If you are installing a private repository, you can skip it. (default: registry-ip)
## Optional
## - cluster-name: use cluster name in config context (default: "kubernetes")
## - install-dir: installation scripts(harbor, shell scripts) save directory (default: "/var/lib/kore-on")
## - cert-validity-days: SSL validity days(default: 36500)
## - debug-mode: Check mode is just a simulation, and if you have steps that use conditionals that depend on the results of prior commands,
##				 it may be less useful for you. (default: false)
## - closed-network: Enable Air Gap (default: false)
#cluster-name = "test-cluster"
#install-dir = "/var/lib/kore-on"
#cert-validity-days = 36500
#debug-mode = true
#closed-network = true
#local-repository-install = true
#local-repository-port = 8080
#local-repository-archive-file = ""
#local-repository-url = "http://x.x.x.x"

[kubernetes]
## Required
## - 
## Optional
## - version: Kubernetes version (default: "latest")
##            If you input only the major version, the minor version automatically selects the last version.
## - container-runtime: use k8s cri (only containerd)
## - kube-proxy-mode: use k8s proxy mode [iptables | ipvs] (default: "ipvs")
## - service-cidr: k8s service network cidr (default: "10.96.0.0/20")
## - pod-cidr: k8s pod network cidr (default: "10.4.0.0/24")
## - node-port-range: k8s node port network range (default: "30000-32767")
## - audit-log-enable: k8s audit log enabled (default: true)
## - api-sans: Add k8s apiserver SAN [--apiserver-cert-extra-sans same as setting] (default: master[0] ip address)
#version = "v1.23.12"
#container-runtime = "containerd"
#kube-proxy-mode = "ipvs"
#service-cidr = "172.20.0.0/24"
#pod-cidr = "10.10.0.0/24"
#node-port-range = "30000-32767"
#audit-log-enable = true
#api-sans = ["x.x.x.x"]

[kubernetes.etcd]
## Required
## - ip: k8s control plane nodes ip address. (Required when selecting the external-etcd="true")
## - private-ip: K8s control plane nodes private ip address. (Required when selecting the external-etcd="true")
##               If you use the same IP address, you can skip it.
## Optional
## - external-etcd: used external etcd than input the ip and private-ip address (default: false)
##                  not used than skip ip address. it is used control plane nodes as automatic.
#external-etcd = true
#ip = ["x.x.x.x"]
#private-ip = ["x.x.x.x"]


[kubernetes.calico]
## Required
## - 
## Optional
## - vxlan-mode: calico VXLAN mode activate (default: false)
#vxlan-mode = true

[node-pool]
## Required
## - 
## Optional
## - data-dir: data(backup, docker, log, kubelet, etcd, k8s-audit, containerd) root dir (default: "/data") 
## - ssh-port: Node ssh port (default: 22) 
#data-dir = "/data"
#ssh-port = 22

[node-pool.master]
## Required
## - ip: k8s control plane nodes ip address.
## - private-ip: K8s control plane nodes private ip address.
##               If you use the same IP address, you can skip it.
## Optional
## - lb-ip: loadbalancer ip address (default: master[0] node ip address)
## - isolated: K8s control plane nodes isolated (default: false)
## - haproxy-install: used internal load-balancer (default: true)
## - lb-ip: Enter the IP address when using a load balancer (default: master[0] ip address)
## - lb-port: Enter the port when using a load balancer (default: "6443")
#ip = ["x.x.x.x","x.x.x.x","x.x.x.x"]
#private-ip = ["x.x.x.x","x.x.x.x","x.x.x.x"]
#isolated = true
#haproxy-install = true
#lb-ip = "x.x.x.x"
#lb-port = "6443"

[node-pool.node]
## Required
## - ip: k8s work nodes ip address.
## - private-ip: K8s work nodes private ip address.
##               If you use the same IP address, you can skip it.
## Optional
#ip = ["x.x.x.x", "x.x.x.x"]
#private-ip = ["x.x.x.x", "x.x.x.x"]

[private-registry]
## Required
## - registry-ip: Public IP address of the private registry node.
##                This is a required entry used when installing a private registry.
## - private-ip: Private IP address of the private registry node. 
##               This is a required entry used when installing a private registry.
##               If you use the same IP address, you can skip it.			  
## Optional
## - install: private registry install (default: false)
## - registry-version: Private registry version. (default: latest)
##                     This is a required entry used when installing a private registry.
##                     If you input only the major version, the minor version automatically selects the last version.
## - registry-domain: If using a domain, please enter the domain. (default: registry-ip)
## - data-dir: private registry data directory (default: "/data/harbor")
## - registry-archive-file: registry archive file path (default: "")
## - public-cert: public cert activate (default: false)
#install = true
#registry-version = "v2.6.0"
#registry-ip = "x.x.x.x"
#private-ip = "x.x.x.x"
#registry-domain = "x.x.x.x"
#data-dir = "/data/harbor"
#registry-archive-file = ""
#public-cert = false

[private-registry.cert-file]
## Required
## - ca-cert: The certificate path used when using public-cert.
##            This is a required field used when using a public certificate.
## Optional
#ca-cert = ""

[shared-storage]
## Required
## - storage-ip: Storage node ip address.
##               This is a required field used when installing the nfs server.
##               (this is using it to generate an inventory and generate an extra vars)
## - private-ip: Storage node ip address.
##               This is a required field used when installing the nfs server.
##               If you use the same IP address, you can skip it.
##               (this is using it to generate an inventory)
## - volume-dir: Storage node data directory. (default: /data/storage)
##               This is a required field used when installing the nfs server.
##               (this is using it to generate an extra vars)
## Optional
## - install: NFS Server Installation (default: false)
#install = true
#type = nfs
#storage-ip = "x.x.x.x"
#private-ip = "x.x.x.x"
#volume-dir = "/data/storage"
#nfs_version = "4.1"

[prepare-airgap]
## Required
## - k8s-version: Kubernetes version (default: "latest")
##                If you input only the major version, the minor version automatically selects the last version.
## - registry-version: Private registry version. (default: latest)
##                     This is a required entry used when installing a private registry.
##                     If you input only the major version, the minor version automatically selects the last version.
## - registry-ip: Private registry node ip address.
##                This is a required field used when the pre-preparation stage of the air gap network.
## Optional
## - 
#k8s-version = "v1.21"
#registry-version = "v2.6"
#registry-ip = "x.x.x.x"
`
