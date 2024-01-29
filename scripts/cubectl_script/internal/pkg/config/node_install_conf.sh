#!/bin/sh

# kubeadm client config
function _conf_kubeadm_client_conf() {
cat <<EOF
apiVersion: kubeadm.k8s.io/v1beta2
kind: JoinConfiguration
caCertPath: CERT_DIR/ca.crt
discovery:
  bootstrapToken:
    apiServerEndpoint: API_LB_IP
    token: KUBEADM_TOKEN
    unsafeSkipCAVerification: true
EOF
}

# kubelet extra config file
function _conf_kubelet_extra_conf() {
cat <<EOF
KUBELET_EXTRA_ARGS="--root-dir=DATA_ROOT_DIR/kubelet \\
LOG_DIR
--v=2 \\
--runtime-request-timeout=15m \\
RESOLV_CONF
--node-ip=NODE_IP \\
--node-labels=cubectl.acornsoft.io/clusterid=CLUSTER_NAME,cubectl.acornsoft.io/ansible_ssh_host=NODE_IP"
EOF
}

function _conf_kubelet_extra_conf_log_dir() {
cat <<EOF
--log-dir=DATA_ROOT_DIR/log \\
--logtostderr=false \\
EOF
}

function _conf_kubelet_extra_conf_resolv_ubuntu() {
cat <<EOF
--resolv-conf=/run/systemd/resolve/resolv.conf \\
EOF
}

function _conf_kubelet_extra_conf_resolv() {
cat <<EOF
--resolv-conf=/etc/resolv.conf \\
EOF
}