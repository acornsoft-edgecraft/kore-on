#!/bin/sh
# set -e
## Haproxy 구성

## ===== [ Handlers ] =====
# 실행 디렉토리 절대경로 찾기
function current_dir() {
  ## scripts 실행 디렉토리 절대 경로 찾기
  local scripts_path=`dirname "$(readlink -f "$0")"`
  local depth=5

  # 구분자로 문자열 분할
  IFS='/' read -ra dir <<< "$scripts_path"
  local total_len="${#dir[@]}"
  local path=""
  
  # 분할된 문자열을 출력 - $2 디렉토리 절대 경로 찾기
  for (( i=1; i<=${total_len}-$depth; i++ )); do
    path="${path}/${dir[i]}"

    if [ "$i" == "$(( ${total_len}-$depth ))" ]; then
      path="${path}"
    fi
  done
  echo "$path"
}

# Create kubelet certificate
function handler_create_k8s_certificates() {
  local cert_dir="/etc/kubernetes/pki"
  local master_cert_dir="/opt/kubernetes/pki"
  local cmd="openssl ecparam -name secp256r1 -genkey -noout -out ${cert_dir}/kubelet-server.key && \
  openssl req -new -sha256 -key ${cert_dir}/kubelet-server.key -subj '/O=system:nodes/CN=system:node:$(hostname)' | \
  openssl x509 -req -CA ${cert_dir}/ca.crt -CAkey ${cert_dir}/ca.key -CAcreateserial -out ${cert_dir}/kubelet-server.crt \
  -days ${CERT_VALIDITY_DAYS} -extensions v3_req_apiserver -extfile ${master_cert_dir}/common-openssl.conf"

  printf "%s %s" "${cmd}" "&& chmod 0644 ${cert_dir}/kubelet-server.crt && chmod 0600 ${cert_dir}/kubelet-server.key"
}

# Install Kubernetes packages (kubectl, kubelet, kubeadm)
function handler_install_k8s_packages() {
  local k8s_version_int=`sed -e "s/^v//g" <<< "$K8S_VERSION"`
  local items=(
    "kubectl-${k8s_version_int}"
    "kubelet-${k8s_version_int}"
    "kubeadm-${k8s_version_int}"
  )

  local cmd="yum install -y ${items[@]}"

  printf "%s" "${cmd}"
}

# Create kubelet directory
function handler_create_kubelet_directory() {
  local path="/var/lib/kubelet"
  local cmd="mkdir -p ${path} && chmod 0755 ${path}"

  printf "%s" "${cmd}"
}

# Create kubeadm token for joining nodes with 24h expiration
function handler_create_kubeadm_token() {
  local kubeconfig="$(current_dir)/config/acloud-client-kubeconfig"
  local cmd="kubeadm token create --kubeconfig=${kubeconfig}"

  printf "%s" "${cmd}"
}

# Copy kubeadm client config
function handler_create_kubeadm_client_config() {
  local kubeadm_token="$1"
  local kubeconfig="$(current_dir)/config/acloud-client-kubeconfig"
  local cert_dir="/etc/kubernetes/pki"
  local kubeadm_client_conf="/etc/kubernetes/kubeadm-client.conf"
  local api_lb_ip=""
  
  if [[ -z "${LB_IP}" ]]; then
    api_lb_ip="${NODE_POOL_MASTER}:6443"
  else
    api_lb_ip="${LB_IP}:${LB_PORT}"
  fi

  local result=`sed \
  -e "s/CERT_DIR/${cert_dir//\//\\\/}/g" \
  -e "s/API_LB_IP/${api_lb_ip//\./\\\.}/g" \
  -e "s/KUBEADM_TOKEN/${kubeadm_token}/g" \
  <<< "$(_conf_kubeadm_client_conf)"`

  local cmd="echo \"${result}\" > ${kubeadm_client_conf} && chmod 0644 ${kubeadm_client_conf}"

  printf "%s" "${cmd}"
}

# Copy kubelet extra config file
function handler_create_kubelet_extra_config() {
  local path="/etc/sysconfig/kubelet"
  readarray -t temp <<< "$(_conf_kubelet_extra_conf)"

  local temp_cnt="${#temp[@]}"

  local temp_path="$(current_dir)/tmp/kubelet_extra_conf"
  touch "${temp_path}"

  if [[ "1" -eq `awk -v n1="${K8S_VERSION//[^v]//}" -v n2="1.26" 'BEGIN { if (n1 < n2) print 1; else print 0 }'` ]]; then
    for (( i=0; i<=${temp_cnt}-1; i++ )); do
      if [[ "$i" == 1 ]]; then
        sed -e "s/DATA_ROOT_DIR/${DATA_ROOT_DIR//\//\\\/}/g" <<< "$(_conf_kubelet_extra_conf_log_dir)" >> "${temp_path}"
      elif [[ "$i" == 4 ]]; then
        echo "$(_conf_kubelet_extra_conf_resolv)" >> "${temp_path}"
      else
        sed -e "s/DATA_ROOT_DIR/${DATA_ROOT_DIR//\//\\\/}/g" \
        -e "s/NODE_IP/${NODE_IP//\./\\\.}/g" \
        -e "s/CLUSTER_NAME/${CLUSTER_NAME}/g" <<< "${temp[$i]}" >> "${temp_path}"
      fi
    done
  else
    for (( i=0; i<=${temp_cnt}-1; i++ )); do
      if [[ "$i" == 4 ]]; then
        echo "$(_conf_kubelet_extra_conf_resolv)" >> "${temp_path}"
      else
        sed -e "s/DATA_ROOT_DIR/${DATA_ROOT_DIR//\//\\\/}/g" \
        -e '/LOG_DIR/d' \
        -e "s/NODE_IP/${NODE_IP//\./\\\.}/g" \
        -e "s/CLUSTER_NAME/${CLUSTER_NAME}/g" <<< "${temp[$i]}" >> "${temp_path}"
      fi
    done
  fi

  local cmd="cat \"${temp_path}\" > ${path} && chmod 0644 ${path}"

  printf "%s" "${cmd}"
}

# Join to cluster
function handler_join_to_cluster() {
  local path="/etc/kubernetes/kubeadm-client.conf"

  local cmd="kubeadm join --config ${path} --ignore-preflight-errors=all"

  printf "%s" "${cmd}"
}

# Wait for kubelet bootstrap to create config
function handler_save_acloud_kueconfig() {
  local path="/etc/kubernetes/acloud/acloud-client-kubeconfig"
  local kubeconfig="$(current_dir)/config/acloud-client-kubeconfig"

  local cmd="mkdir -p ${path} && cp ${kubeconfig} ${path} && chmod 0644 ${path}"

  printf "%s" "${cmd}"
}

# Update server field in kubelet kubeconfig (haproxy)
function handler_update_kubelet_conf() {
  local path="/etc/kubernetes/kubelet.conf"
  
  local cmd="sed -i 's#server:.*#server: https://localhost:6443#g' ${path}"

  printf "%s" "${cmd}"
}

# kubectl label node
function handler_kubectl_label_node() {
  local kubeconfig="$(current_dir)/config/acloud-client-kubeconfig"

  local cmd="kubectl --kubeconfig=${kubeconfig} label node $(hostname) node-role.kubernetes.io/node='' --overwrite"

  printf "%s" "${cmd}"
}

# kubectl label node
function handler_kubectl_status_node() {
  local kubeconfig="$(current_dir)/config/acloud-client-kubeconfig"

  local cmd="kubectl --kubeconfig=${kubeconfig} get nodes -o wide"

  printf "%s" "${cmd}"
}
## ===== [ Handlers End ] =====

## ===== [ Tasks ] =====
# distribute certificate in worker
function node_install() {
  local task_name="Node installation"

  local res=""
  local exit_code=""
  local skip_run_index=()
  local skip_task_index=()
  local kubeadm_token=""

  ## command here
  local commands=(
    "$(handler_create_k8s_certificates)"
    "$(handler_install_k8s_packages)"
    "$(handler_create_kubelet_directory)"
    "$(handler_create_kubeadm_token)"
    "$(handler_create_kubeadm_client_config "${kubeadm_token}")"
    "$(handler_create_kubelet_extra_config)"
    "systemctl daemon-reload"
    "systemctl enable --now kubelet.service"
    "$(handler_join_to_cluster)"
    # "$(handler_save_acloud_kueconfig)"
    "$(handler_update_kubelet_conf)"
    "$(handler_kubectl_label_node)"
    "$(handler_kubectl_status_node)"
    )
  local commands_description=(
    "Create kubelet certificate"
    "Install Kubernetes packages (kubectl, kubelet, kubeadm)"
    "Create kubelet directory"
    "Create kubeadm token for joining nodes with 24h expiration"
    "Copy kubeadm client config"
    "Copy kubelet extra config file"
    "systemctl daemon-reload"
    "Start and enable kubelet on worker node"
    "Join to cluster"
    # "Save acloud_kueconfig"
    "Update server field in kubelet kubeconfig (haproxy)"
    "kubectl label node"
    "kubectl status node"
  )
  
  local total_len="${#commands[@]}"

  # get task status - 성공한 task는 skip 하기 위함
  skip_task_index=(`_get_task_status "$task_name" "$total_len"`)

  for (( i=0; i<=${total_len}-1; i++ )); do
    ## skip 설정
    # skip_run: 명령어 실행 skip
    # skip_task: 진행상태별 성공한 command skip
    # skip_run_index: 실행중 skip할 명령어 설정
    # skip_index: 진행 상태별 성공한 command save skip task 설정
    local skip_run=`_skip_run "$i" "${skip_run_index[@]}"`
    local skip_task=`_skip_task "${task_name}::$i"`
    local skip_index=""

    if [[ "$i" == 4 ]]; then
      commands[4]="$(handler_create_kubeadm_client_config "${kubeadm_token}")"
    fi

    if [[ "$skip_task" == true ]]; then
      _log_skip "${commands[i]}" "$log_path" "" "${task_name}::${commands_description[$i]}"
    else
      if [[ "$skip_run" == false  ]]; then
      res=`_run "${commands[i]}" "$log_path" "" "${task_name}::${commands_description[$i]}"`
      _response "$res"

      ## Data Handler -- START
      if [[ "$i" == 3 ]]; then
        kubeadm_token="${res_data}"
      fi
      if [[ "$i" == "$((${total_len}-1))" ]]; then
        res_log+='\n'
        res_log+="${res_data}"
      fi
      ## Data Handler -- END

      ## response 결과 출력
      # 실패시 결과값 출력 후 종료
      # 실행 상태 저장
      # log 출력 
      # response clear
      echo "$res_log"
      _save_task_status "$res_code" "$task_name" "$i" "$skip_index" "${commands_description[$i]}"
      _check_res_code "$res_code"
      _clear_response
      fi
    fi
  done
}
## ===== [ Tasks End ] =====

main() {
  ## ===== [ includes ] =====
  source "$(current_dir)/internal/pkg/utils/common.sh" "$(current_dir)/internal/pkg"
  source "$(current_dir)/internal/pkg/utils/logger.sh"
  source "$(current_dir)/internal/pkg/config/node_install_conf.sh"
  source "$(current_dir)/config/add_node_env.rc"

  ## ===== [ Constants and Variables ] =====
  local log_path="$(current_dir)/logs/add_node.log"
  local IFS="\t\n"
  local status=""

  ## Tasks: node_install
  node_install

}

main "${@}"