#!/bin/sh
set -e
### Kubernetes 클러스터 노드 추가

## ===== [ Handlers ] =====
# 실행 디렉토리 절대경로 찾기
function current_dir() {
  ## scripts 실행 디렉토리 절대 경로 찾기
  local scripts_path=`dirname "$(readlink -f "$0")"`
  local depth=3

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
## ===== [ Handlers End ] =====

## ===== [ Tasks ] =====
# OS release
function os_release_export() {
  local role_name="GET OS Release"

  local res=""
  local exit_code=""
  local values=""

  res=`_run "cat /etc/*release" "$log_path" "$role_name" ""`
  if [[ -z `grep "\[FATAL\]" <<< "$res"` ]]; then
    grep -i "=" <<< "$res" > $(current_dir)/config/os_release.rc
    source "$(current_dir)/config/os_release.rc"
  fi
}

# Init
function init() {
  local role_name="INIT"

  local start_time=$(TZ='Asia/Seoul' date +%s.%N)
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  local end_time=""
  local duration_of_time=""

  local res=""
  local exit_code=""
  
  _run "sh $(current_dir)/internal/playbook/roles/init/$ID/init.sh" "$log_path" "$role_name" ""
  _clear_tmp "$(current_dir)" "$log_path"
  _check_res_code "$res_code"
}

# Bootstrap OS
function bootstrap_os() {
  local role_name="Bootstrap OS"

  local start_time=$(TZ='Asia/Seoul' date +%s.%N)
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  local end_time=""
  local duration_of_time=""

  local res=""
  local exit_code=""
    
  _run "sh $(current_dir)/internal/playbook/roles/bootstrap_os/$ID/add_local_repo.sh" "$log_path" "$role_name" ""
  _clear_tmp "$(current_dir)" "$log_path"
  _check_res_code "$res_code"
}

## Node Initialize
function initialize() {
  local role_name="Node Initialize"

  local start_time=$(TZ='Asia/Seoul' date +%s.%N)
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  local end_time=""
  local duration_of_time=""

  local res=""
  local exit_code=""

  _run "sh $(current_dir)/internal/playbook/roles/initialize/$ID/initialize.sh" "$log_path" "$role_name" ""
  _clear_tmp "$(current_dir)" "$log_path"
  _check_res_code "$res_code"
}

## CRI 구성
function cri_configuration() {
  local role_name="CRI Configuration"

  local start_time=$(TZ='Asia/Seoul' date +%s.%N)
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  local end_time=""
  local duration_of_time=""

  local res=""
  local exit_code=""

  _run "sh $(current_dir)/internal/playbook/roles/cri/cri_install.sh" "$log_path" "$role_name" ""
  _clear_tmp "$(current_dir)" "$log_path"
  _check_res_code "$res_code"
}

## Haproxy 구성
function haproxy_configuration() {
  local role_name="Haproxy Configuration"

  local start_time=$(TZ='Asia/Seoul' date +%s.%N)
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  local end_time=""
  local duration_of_time=""

  local res=""
  local exit_code=""

  _run "sh $(current_dir)/internal/playbook/roles/haproxy/haproxy_install.sh" "$log_path" "$role_name" ""
  _clear_tmp "$(current_dir)" "$log_path"
  _check_res_code "$res_code"
}

## Node 구성
function node_configuration() {
  local role_name="Node Configuration"

  local start_time=$(TZ='Asia/Seoul' date +%s.%N)
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  local end_time=""
  local duration_of_time=""

  local res=""
  local exit_code=""

  _run "sh $(current_dir)/internal/playbook/roles/node/node_run_tasks.sh" "$log_path" "$role_name" ""
  _clear_tmp "$(current_dir)" "$log_path"
  _check_res_code "$res_code"
}
## ===== [ Tasks End ] =====

main() {
  ## ===== [ includes ] =====
  source "$(current_dir)/internal/pkg/utils/common.sh" "$(current_dir)/internal/pkg"
  source "$(current_dir)/internal/pkg/utils/config.sh"
  source "$(current_dir)/config/add_node_env.rc"

  ## ===== [ Constants and Variables ] =====
  # local start_time=$(TZ='Asia/Seoul' date +%s.%N)
  # local end_time=""
  # local duration_of_time=""
  # res_code: 1 -> 싪패, 0 -> 성공
  local log_path="$(current_dir)/logs/add_node.log"
  local res_log=""
  local res_data=""
  local res_code=""

  ## Role: os_release_export - First run required
  os_release_export

  ## Role: init
  init

  ## Role: bootstrap_os
  bootstrap_os

  ## Role: bootstrap_os
  initialize

  ## Role: cri_configuration
  cri_configuration

  ## Role: haproxy_configuration
  haproxy_configuration

  ## Role: node_configuration
  node_configuration

}

main "${@}"
