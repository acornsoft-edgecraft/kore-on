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

# Create kubernetes manifests / haproxy config directory
function handler_crate_haproxy_config_dir() {
  local path="/etc/kubernetes/manifests /etc/haproxy"
  printf "%s %s %s" "mkdir -p " "$path" "&& chmod 0755 $path"
}

# Copy kubernetes manifests haproxy.yaml
function handler_copy_haproxy_yaml() {
  local path="/etc/kubernetes/manifests/haproxy.yaml"

  printf "%s %s" "echo \"$(_conf_haproxy_yaml)\" > $path" "&& chmod 0644 $path"
}

# Copy haproxy.cfg
function handler_copy_haproxy_cfg() {
  local path="/etc/haproxy/haproxy.cfg"
  local mater_pool=("${NODE_POOL_MASTER[@]}")
  local data=""

  local total_len="${#mater_pool[@]}"

  for (( i=0; i<=${total_len}-1; i++ )); do
    data+="server  api$i  ${mater_pool[$i]}:6443  check"'\n    '
  done

  local result=`sed \
  -e "s/API_BACKENDS/${data}/g" \
  <<< "$(_conf_haproxy_cgf)"`

  printf "%s %s" "echo \"${result}\" > $path" "&& chmod 0644 $path"
}
## ===== [ Handlers End ] =====

## ===== [ Tasks ] =====
# Install containerd
function install_haproxy() {
  local task_name="Install haproxy"

  local res=""
  local exit_code=""
  local skip_run_index=()
  local skip_task_index=()

  ## command here
  local commands=(
    "$(handler_crate_haproxy_config_dir)"
    "$(handler_copy_haproxy_yaml)"
    "$(handler_copy_haproxy_cfg)"
    )
  local commands_description=(
    "Create haproxy directory (kubernetes manifests/haproxy config)"
    "Copy kubernetes manifests haproxy.yaml"
    "Copy haproxy.cfg"
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
  source "$(current_dir)/internal/pkg/config/haproxy_conf.sh"
  source "$(current_dir)/config/add_node_env.rc"

  ## ===== [ Constants and Variables ] =====
  local log_path="$(current_dir)/logs/add_node.log"
  local IFS="\t\n"
  local status=""

  ## Tasks: Install containerd
  install_haproxy

}

main "${@}"