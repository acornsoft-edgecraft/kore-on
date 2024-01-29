#!/bin/sh
# set -e
## Create kubernetes cert 구성

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

# Create kubernetes pki directory
function handler_crate_opt_kubernetes_pki_dir() {
  local path="/etc/kubernetes/pki /opt/kubernetes/pki"
  
  printf "%s %s %s" "mkdir -p " "$path" "&& chmod 0755 $path"
}

# Create openssl common-conf file
function handler_copy_common_openssl_conf() {
  local path="/opt/kubernetes/pki/common-openssl.conf"

  local add_dns=""
  local add_ip=""
  local service_ip=`
  awk -F'/' '{
    split($1, octets, "."); 
    last_octet = octets[4]; 
    if (last_octet == 0) last_octet++; 
    printf("%s.%s.%s.%s\n", octets[1], octets[2], octets[3], last_octet)
  }' <<< "${SERVICE_CIDR}"`

  if [[ -z "${LB_IP}" ]]; then
    add_ip+="IP.4 = ${NODE_POOL_MASTER[0]}"'\n'
  else
    add_ip+="IP.4 = ${LB_IP}"'\n'
  fi

  local api_sans_len="${#API_SANS[@]}"
  if [[ -n "${API_SANS[0]}" ]]; then
    for (( i=0; i<=${api_sans_len}-1; i++ )); do
      add_ip+="IP.$((5+$i)) = ${API_SANS[$i]}"'\n'
    done
  fi

  local result=`sed \
  -e "s/HOST_NAME/$(hostname)/" \
  -e "s/NODE_IP/${NODE_IP}/" \
  -e "s/KUBERNETES_SERVICE_IP/${service_ip}/" \
  -e "s/ADD_DNS/${add_dns}/" \
  -e "s/ADD_IP/${add_ip//\./\\\.}/" \
  <<< "$(_conf_common_openssl_conf)"`

  printf "%s %s" "echo \"${result}\" > $path" "&& chmod 0644 $path"
}

# Copy kubernetes ca certificate
function handler_copy_k8s_ca_certificate() {
  local path="/etc/kubernetes/pki"
  local copy_path="$(current_dir)/config/pki"

  printf "%s %s" "cp -f ${copy_path}/ca.* ${path}" "&& chmod 0644 ${path}/ca.crt && chmod 0600 ${path}/ca.key"
}

# Copy kubernetes ca certificate to opt directory
function handler_copy_opt_ca_certificate() {
  local path="/opt/kubernetes/pki"
  local copy_path="$(current_dir)/config/pki"

  printf "%s %s" "cp -f ${copy_path}/ca.* ${path}" "&& chmod 0644 ${path}/ca.crt && chmod 0600 ${path}/ca.key"
}

# Get registry ca certificate
function handler_get_registry_ca_certificate() {
  local path="/etc/docker/certs.d/${REGISTRY_DOMAIN}"
  local cmd="curl -kL https://${REGISTRY_DOMAIN}/api/v2.0/systeminfo/getcert"

  printf "%s %s" "${cmd} > ${path}/ca.crt" "&& chmod 0644 ${path}/ca.crt"
}
## ===== [ Handlers End ] =====

## ===== [ Tasks ] =====
# distribute certificate in worker
function distribute_cert_worker() {
  local task_name="Create kubernetes certificates"

  local res=""
  local exit_code=""
  local skip_run_index=()
  local skip_task_index=()

  ## command here
  local commands=(
    "$(handler_crate_opt_kubernetes_pki_dir)"
    "$(handler_copy_common_openssl_conf)"
    "$(handler_copy_k8s_ca_certificate)"
    "$(handler_copy_opt_ca_certificate)"
    "$(handler_get_registry_ca_certificate)"
    )
  local commands_description=(
    "Create kubernetes pki directory"
    "Create openssl common-conf file"
    "Copy kubernetes ca certificate"
    "Copy kubernetes ca certificate to opt directory"
    "Get registry ca certificate"
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

    if [[ "$skip_task" == true ]]; then
      _log_skip "${commands[i]}" "$log_path" "" "${task_name}::${commands_description[$i]}"
    else
      if [[ "$skip_run" == false  ]]; then
      res=`_run "${commands[i]}" "$log_path" "" "${task_name}::${commands_description[$i]}"`
      _response "$res"

      ## Data Handler -- START
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
  source "$(current_dir)/internal/pkg/config/common_openssl_conf.sh"
  source "$(current_dir)/config/add_node_env.rc"

  ## ===== [ Constants and Variables ] =====
  local log_path="$(current_dir)/logs/add_node.log"
  local IFS="\t\n"
  local status=""

  ## Tasks: Create kubernetes cert
  distribute_cert_worker

}

main "${@}"