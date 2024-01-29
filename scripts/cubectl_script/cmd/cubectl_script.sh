#!/bin/sh
set -e

## ===== [ Handlers ] =====
# 실행 디렉토리 절대경로 찾기
function current_dir() {
  ## scripts 실행 디렉토리 절대 경로 찾기
  local scripts_path=`dirname "$(readlink -f "$0")"`
  local depth=2

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

## ===== [ Sub Commands Start] =====
# Add Node
function run_command() {
  local sub_command="$1"

  ## sub command list
  # command_index name과 실행 쉘스크립트 name은 같아야 한다. (예: test => tesh.sh)
  local command_index=(
    "add_node"
    "upgrade_cluster"
  )
  local commands_description=(
    "Add node in k8s cluster"
    "Kubernetes version upgrade for k8s cluster"
  )

  local total_len="${#command_index[@]}"

  if [[ -z "${sub_command}" ]]; then
    echo "Required input the sub-command"
    return
  elif [[ "list" == "$sub_command" ]]; then
    printf "%-30s %-40s\n" "NAME" "DESCRIPTION" 
    printf "%-30s %-40s\n" "--------------------" "----------------------------------------" 
    for (( i=0; i<=${total_len}-1; i++)); do
      printf "%-30s %-40s\n" "${command_index[$i]}" "${commands_description[$i]}" 
    done
    return
  elif [[ -z `grep "${sub_command}" <<< "${command_index[@]}"` ]]; then
    echo "Not found!! sub-command: ${sub_command}"
    return
  fi

  # _run "sh $(current_dir)/internal/playbook/${sub_command}.sh" "$log_path" "" "" "${sub_command}"
  # _check_res_code "$res_code"

  sh $(current_dir)/internal/playbook/${sub_command}.sh
  _check_res_code "$res_code"
}
## ===== [ Sub Commands End] =====

main() {
  ## ===== [ includes ] =====
  source "$(current_dir)/internal/pkg/utils/common.sh" "$(current_dir)/internal/pkg"

  ## ===== [ Constants and Variables ] =====
  # local start_time=$(TZ='Asia/Seoul' date +%s.%N)
  # local end_time=""
  # local duration_of_time=""
  # res_code: 1 -> 싪패, 0 -> 성공
  local log_path="$(current_dir)/logs/add_node.log"
  local res_log=""
  local res_data=""
  local res_code=""

  ## ===== [ run sub command ] =====
  ## Sub Command 실행
  # args: $1
  # $1: sub command parameter
  run_command $1
}

main "${@}"