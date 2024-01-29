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
## ===== [ Handlers End ] =====

## ===== [ Tasks ] =====
## ===== [ Tasks End ] =====

## ===== [ Include Task ] =====
# include_sslcert
function include_distributecert_worker() {
  local task_name="Node configuration"

  local start_time=$(TZ='Asia/Seoul' date +%s.%N)
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  local end_time=""
  local duration_of_time=""

  local res=""
  local exit_code=""

  _run "sh $(current_dir)/internal/playbook/roles/node/distributecert_worker.sh" "$log_path" "" "$task_name"
  _check_res_code "$res_code"
}

# include_sslcert
function include_node_install() {
  local task_name="Node Installation"

  local start_time=$(TZ='Asia/Seoul' date +%s.%N)
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  local end_time=""
  local duration_of_time=""

  local res=""
  local exit_code=""

  _run "sh $(current_dir)/internal/playbook/roles/node/node_install.sh" "$log_path" "" "$task_name"
  _check_res_code "$res_code"
}
## ===== [ Include Task End ] =====

main() {
  ## ===== [ includes ] =====
  source "$(current_dir)/internal/pkg/utils/common.sh" "$(current_dir)/internal/pkg"
  source "$(current_dir)/internal/pkg/utils/logger.sh"
  source "$(current_dir)/config/add_node_env.rc"

  ## ===== [ Constants and Variables ] =====
  local log_path="$(current_dir)/logs/add_node.log"
  local IFS="\t\n"
  local status=""
  local kubeadm_token=""

  ## Tasks: include_distributecert_worker
  include_distributecert_worker

  ## Tasks: include_node_install
  include_node_install

}

main "${@}"