#!/bin/sh
# Config Data

## local repository 설정값
function _conf_local_repo() {
cat <<EOF 
[local-repo]
async = 1
baseurl = LOCAL_REPO_URL
enabled = 1
gpgcheck = 0
name = Local Repo configure
EOF
}

## kernel modules for Calico
function _conf_cluster_calico_kernel () {
cat <<EOF
ip_set
ip_tables
ip6_tables
ipt_REJECT
ipt_set
xt_addrtype
xt_comment
xt_conntrack
xt_ipvs
xt_mark
xt_multiport
xt_sctp
xt_set
ipip
nf_conntrack_netlink
ipt_rpfilter
EOF
}

## Disable NetworkManager DNS processing when RHEL 8
function _conf_networkmanager () {
cat <<EOF
[main]
dns=none
EOF
}

## Prevent NetworkManager from managing Calico interfaces
function _conf_networkmanager_calico () {
cat <<EOF
[keyfile]
unmanaged-devices=interface-name:cali*;interface-name:tunl*;interface-name:vxlan.calico
EOF
}

## Modprobe Kernel Module for IPVS
function _conf_modpwobe_ipvs() {
cat <<EOF
ip_vs
ip_vs_rr
ip_vs_wrr
ip_vs_sh
nf_conntrack
EOF
}

## Forwarding IPv4 and letting iptables see bridged traffic
function _conf_sysctl_k8s() {
cat <<EOF
net.bridge.bridge-nf-call-iptables=1
net.bridge.bridge-nf-call-ip6tables=1
net.bridge.bridge-nf-call-arptables=1
net.ipv4.ip_forward=1
EOF
}

## Run the Docker daemon as a non-root user (Rootless mode)
function _conf_sysctl_dind() {
cat <<EOF
user.max_user_namespaces=28633
EOF
}


## Run Shell script
function _conf_script() {
cat <<EOF
#!/bin/sh
set -e

## Run commands

EOF
}
