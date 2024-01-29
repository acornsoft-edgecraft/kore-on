#!/bin/sh
# CRI Config Data

## container images url
function _conf_haproxy_cgf() {
cat <<EOF 
global
  log 127.0.0.1 local0
  log 127.0.0.1 local1 notice
  tune.ssl.default-dh-param 2048

defaults
  log global
  mode http
  #option httplog
  option dontlognull
  timeout connect 5000ms
  timeout client 50000ms
  timeout server 50000ms

frontend healthz
  bind *:8081
  mode http
  monitor-uri /healthz

frontend api-https
   mode tcp
   bind 127.0.0.1:6443
   default_backend api-backend

backend api-backend
    mode tcp
    API_BACKENDS
EOF
}

function _conf_haproxy_yaml() {
cat <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: haproxy
  namespace: kube-system
  labels:
    addonmanager.kubernetes.io/mode: Reconcile
    k8s-app: kube-haproxy
spec:
  hostNetwork: true
  nodeSelector:
    beta.kubernetes.io/os: linux
  priorityClassName: system-node-critical
  containers:
  - name: haproxy
    image: "haproxy:latest"
    imagePullPolicy: IfNotPresent
    resources:
      requests:
        cpu: 25m
        memory: 32M
    securityContext:
      privileged: true
    livenessProbe:
      httpGet:
        path: /healthz
        port: 8081
    readinessProbe:
      httpGet:
        path: /healthz
        port: 8081
    volumeMounts:
    - mountPath: /usr/local/etc/haproxy/haproxy.cfg
      name: etc-haproxy
      readOnly: true
  volumes:
  - name: etc-haproxy
    hostPath:
      path: /etc/haproxy/haproxy.cfg
EOF
}