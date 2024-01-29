#!/bin/sh

## container images url
function _conf_common_openssl_conf() {
cat <<EOF 
[ req ]
distinguished_name = req_distinguished_name
[req_distinguished_name]

[ v3_ca ]
basicConstraints = critical, CA:TRUE
keyUsage = critical, digitalSignature, keyEncipherment, keyCertSign

[ v3_req_server ]
basicConstraints = CA:FALSE
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth

[ v3_req_client ]
basicConstraints = CA:FALSE
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth

[ v3_req_apiserver ]
basicConstraints = CA:FALSE
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names_cluster

[ v3_req_metricsserver ]
basicConstraints = CA:FALSE
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_metircs_server

[ alt_metircs_server ]
DNS.1 = metrics-server
DNS.2 = metrics-server.kube-system
DNS.3 = metrics-server.kube-system.svc
DNS.4 = metrics-server.kube-system.svc.cluster.local
DNS.5 = localhost
IP.1 = 127.0.0.1

[ alt_names_cluster ]
DNS.1 = kubernetes
DNS.2 = kubernetes.default
DNS.3 = kubernetes.default.svc
DNS.4 = kubernetes.default.svc.cluster.local
DNS.5 = localhost
DNS.6 = HOST_NAME
ADD_DNS
IP.1 = 127.0.0.1
IP.2 = NODE_IP
IP.3 = KUBERNETES_SERVICE_IP
ADD_IP
EOF
}