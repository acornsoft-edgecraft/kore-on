package conf

const Template = `#koreon.toml
#koreon.toml
[koreon]
cluster-name = "ml-central2"
#cert-validity-days = 3650

#closed-network = true
#local-repository = "http://192.168.1.251:8080"
#local-repository-archive-file = "/Users/okpiri/git/kore-on/work/local-repo.20220224_071700.tgz"
#install-dir = "/var/lib/knit"
#debug-mode = true

[kubernetes]
version = "1.19.10"
#container-runtime = "containerd"
#kube-proxy-mode = "ipvs"
#vxlan-mode = true
#service-cidr ="10.96.0.0/12"
#pod-cidr="10.32.0.0/12"
#node-port-range="30000-32767"
#audit-log-enable = true
#api-sans = ["192.168.1.9"]


[kubernetes.etcd]
ip = ["192.168.88.141"]
private-ip = ["172.33.88.141"]
encrypt-secret = true


[node-pool]
#data-dir = "/data"

[node-pool.security]
ssh-user-id = "centos"
private-key-path = "/Users/okpiri/cert/hostacloud/id_rsa"

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
volume-dir = "/data/nvme/mlops/161"
volume-size = 1000


[private-registry]
install = true
registry-ip = "192.168.88.145"
data-dir = "/data/harbor"
public-cert = false
registry-archive-file = "/Users/okpiri/git/kore-on/work/harbor.20220224_072307.tgz"

[private-registry.cert-file]
ssl-certificate = ""
ssl-certificate-key = ""
`
