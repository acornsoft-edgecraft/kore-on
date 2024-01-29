#!/bin/sh
set -e
## initialize 구성

## ===== [ Handlers ] =====
# 실행 디렉토리 절대경로 찾기
function current_dir() {
  ## scripts 실행 디렉토리 절대 경로 찾기
  local scripts_path=`dirname "$(readlink -f "$0")"`
  local depth=6

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
# yum clean all & delete dnf cache
function yum_clean_all() {
  local task_name="yum clean all"

  local res=""
  local exit_code=""
  local skip_run_index=()
  local skip_task_index=()

  ## command here
  local commands=(
    "yum clean all"
    "rm -rf /var/cache/yum"
    )
  local commands_description=(
    "yum clean all"
    "delete (yum/dnf) cache"
  )

  local num2=8
  local result=$(awk -v n1="$VERSION_ID" -v n2="$num2" 'BEGIN { if (n1 > n2) print 1; else print 0 }')
  if [[ "$result" -eq 1 ]]; then
    commands[1]=`sed 's/\/var\/cache\/yum/\/var\/cache\/dnf/g' <<< "${commands[1]}"`
  fi

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
      # echo "$res"
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

# Disable selinux
function disable_selinux() {
  local task_name="Disable selinux"

  local res=""
  local exit_code=""
  local selinux_status=""
  local skip_run_index=()
  local skip_task_index=()

  ## command here
  local commands=(
    "getenforce"
    "setenforce 0"
    "sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config"
    )
  local commands_description=(
    "SELinux 상태 확인"
    "disable selinux"
    "permissive 모드로 SELinux 설정"
  )
  local total_len="${#commands[@]}"

  # get task status - 성공한 task는 skip 하기 위함
  skip_task_index=(`_get_task_status "$task_name" "$total_len"`)

  for (( i=0; i<=${total_len}-1; i++ )); do
    ## skip 설정
    # skip_run: 명령어 실행 skip
    # skip_task: 진행 상태별 성공한 command skip
    # skip_run_index: 실행중 skip할 명령어 설정
    # skip_index: 진행 상태별 성공한 command save skip task 설정
    local skip_run=`_skip_run "$i" "${skip_run_index[@]}"`
    local skip_task=`_skip_task "${task_name}::$i"`
    local skip_index=""
    local skip_description=""

    if [[ "$skip_task" == true ]]; then
      _log_skip "${commands[i]}" "$log_path" "" "${task_name}::${commands_description[$i]}"
    else
      if [[ "$skip_run" == false  ]]; then
        res=`_run "${commands[i]}" "$log_path" "" "${task_name}::${commands_description[$i]}"`
        _response "$res"
        
        ## Data Handler -- START
        skip_description="${commands_description[$i]}|"

        if [[ "$i" == 0 && `grep "Disabled" <<< "$res"` ]]; then
          skip_run_index+=(1 2)
          skip_index+="1/2/"
            skip_description+="${commands_description[1]}|"
            skip_description+="${commands_description[2]}|"
        fi

        ## Data Handler -- END

        ## response 결과 출력
        # 실패시 결과값 출력 후 종료
        # 실행 상태 저장
        # log 출력 
        # response clear
        echo "$res_log"
        _save_task_status "$res_code" "$task_name" "$i" "$skip_index" "$skip_description"
        _check_res_code "$res_code"
        _clear_response
      fi
    fi
  done
}
 
## ===== [ Tasks End ] =====

## ===== [ Include Task ] =====
# cluster.sh
function include_cluster() {
  local task_name="Kernel configuration"

  local start_time=$(TZ='Asia/Seoul' date +%s.%N)
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  local end_time=""
  local duration_of_time=""

  local res=""
  local exit_code=""

  _run "sh $(current_dir)/internal/playbook/roles/initialize/rhel/cluster.sh" "$log_path" "" "$task_name"
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

  ## Tasks: yum_clean_all
  yum_clean_all

  ## Tasks: Disable selinux
  disable_selinux

  ## Include Tasks: cluster.sh
  include_cluster

  # ## 스크립트 종료 상태 변경
  # if [[ "$status" == "0" ]]; then
  #   exit 0
  # fi
}

main "${@}"