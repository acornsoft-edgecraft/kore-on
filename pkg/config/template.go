package config

const Template = `
# "###"는 미사용중 제거 대상
# koreon.toml
[koreon]
cluster-name = "aaa"
#cert-validity-days = 3650 #default 36500

#closed-network = true
#local-repository = "http://192.168.88.145:8080"
#local-repository-archive-file = "/tmp/koreon/local-repo.20220224_071700.tgz"
#install-dir = "/var/lib/koreon"
#debug-mode = true

[kubernetes]
version = "1.19.10"
#container-runtime = "containerd"
#kube-proxy-mode = "ipvs"
#service-cidr ="10.96.0.0/12"
#pod-cidr="10.32.0.0/12"
#node-port-range="30000-32767"
#audit-log-enable = true
#api-sans = ["192.168.1.9"]

[kubernetes.calico]
version = "latest"
#vxlan-mode = true

[kubernetes.etcd]
ip = ["192.168.88.141"]
private-ip = ["172.33.88.141"]
encrypt-secret = true


[node-pool]
#data-dir = "/data"

[node-pool.security]
ssh-user-id = "centos"
private-key-path = "/tmp/cert/id_rsa"

[node-pool.master]
ip = ["192.168.88.141"]
private-ip = ["172.33.88.141"]
lb-ip = "192.168.88.141"
#isolated = true
#haproxy-install = true

[node-pool.node]
ip = ["192.168.88.142", "192.168.88.143"]
private-ip = ["172.33.88.142", "172.33.88.143"]


[shared-storage]
install = false
storage-ip = "192.168.88.11"
volume-dir = "/data/cluster"
volume-size = 1000


[private-registry]
install = true
registry-ip = "192.168.88.145"
data-dir = "/data/harbor"
public-cert = false
#registry-archive-file = "/tmp/koreon/harbor.20220224_072307.tgz"

[private-registry.cert-file]
ssl-certificate = ""
ssl-certificate-key = ""
`
