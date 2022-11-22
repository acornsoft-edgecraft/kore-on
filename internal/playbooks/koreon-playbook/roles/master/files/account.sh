#!/bin/bash

KUBE_CERT="/etc/kubernetes/admin.conf"

cat > openssl.cnf <<EOF
[ req ]
distinguished_name = req_distinguished_name
[req_distinguished_name]

[ v3_req_client ]
basicConstraints = CA:FALSE
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth
EOF

openssl genrsa -out acloud-client.key 2048
openssl req -new -key acloud-client.key -days 3650 -out acloud-client.csr -subj "/CN=acloud-client" -config ./openssl.cnf

cat > acloud-client-csr.yml <<EOF
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: acloud-client
spec:
  groups:
  - system:authenticated
  request: $(cat acloud-client.csr | base64 | tr -d '\n')
  usages:
  - digital signature
  - key encipherment
  - client auth
EOF

kubectl --kubeconfig=${KUBE_CERT} delete csr acloud-client
kubectl --kubeconfig=${KUBE_CERT} apply -f acloud-client-csr.yml

while true
do
	ret=$(kubectl --kubeconfig=${KUBE_CERT} get csr acloud-client | grep -v NAME | awk '{print $4}')

	if [ "$ret" == "Pending" ]; then
	    echo "Found acloud-client csr with Pending status"
		break
	else
		echo "acloud-client csr not found, sleep 1 sec ..."
		sleep 1
	fi
done

echo "Approve acloud-client csr"
kubectl --kubeconfig=${KUBE_CERT} certificate approve acloud-client

while true
do
	ret=$(kubectl --kubeconfig=${KUBE_CERT} get csr acloud-client | grep -v NAME | awk '{print $4}')

	if [ "$ret" == "Approved,Issued" ]; then
		echo "Found acloud-client csr with Approved,Issued status"
		break
	else
		echo "acloud-client did not approved, sleep 1 sec ..."
		sleep 1
	fi
done

kubectl --kubeconfig=${KUBE_CERT} get csr acloud-client -o jsonpath='{.status.certificate}' | base64 --decode > acloud-client.crt
echo "Make acloud-client.crt is done"

cat > acloud-client-crb.yml <<EOF
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: acloud-binding
subjects:
- kind: User
  name: acloud-client
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
EOF

kubectl --kubeconfig=${KUBE_CERT} apply -f acloud-client-crb.yml
echo "Create clusterrolebinding fo acloud-client is done"

cat > acloud-client-kubeconfig <<EOF
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: $(cat ${KUBE_CERT} | grep certificate-authority-data | awk '{print $2}')
    server: $(cat ${KUBE_CERT} | grep server | awk '{print $2}')
  name: acloud-client
contexts:
- context:
    cluster: acloud-client
    user: acloud-client
  name: acloud-client
current-context: acloud-client
kind: Config
preferences: {}
users:
- name: acloud-client
  user:
    client-certificate-data: $(cat acloud-client.crt | base64 | tr -d '\n')
    client-key-data: $(cat acloud-client.key | base64 | tr -d '\n')
EOF

echo "Create acloud-client-kubeconfig is done"