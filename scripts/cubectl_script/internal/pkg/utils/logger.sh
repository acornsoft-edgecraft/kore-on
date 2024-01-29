#!/bin/sh
## 모든 스크립트에서 호출되어야 하는 logger 스크립트.

set -e
set -o errexit

set -E
set -o errtrace

set +H
set +o histexpand

## log 출력/저장
# args: $1, $2, $3, $4, $5, $6, $7
# $1: 실핼 명형어
# $2: 로그 디렉토리
# $3: 샐행 결과 message
# $4: 샐행 결과 code
# $5: 소요 시간
# $6: Playbook role name
# $7: task name
function _log_info() {
  local log_command="$1"
  local log_path="$2"
  local log_message="$3"
  local log_exit_code="$4"
  local duration_of_time=`printf "%10s" "$5"`
  # local log_duration_of_time="$5"
  local log_role_name="$6"
  local log_task_name="$7"
  local log_sub_cmd_name="$8"

  local log_level="INFO"
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  
  if [[ -z "$log_exit_code" ]]; then
    log_exit_code="0"
  fi
  if [[ -n "$log_role_name" ]]; then
    echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [role]:'$log_role_name'"
    # echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [role]:'$log_role_name'"  >> "$log_path"
  fi
  if [[ -n "$log_task_name" ]]; then
    echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_task_name'"
    # echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_task_name'" >> "$log_path"
    # echo "$log_message" >> "$log_path"
  fi
  if [[ -n "$log_sub_cmd_name" ]]; then
    echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [sub-command]:'$log_sub_cmd_name'"
  fi
  # echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_title'" >> "$log_path"
  # echo "$log_message" >> "$log_path"
  # echo "" >> "$log_path"
}

function _log_warning() {
  local log_command="$1"
  local log_path="$2"
  local log_message="$3"
  local log_exit_code="$4"
  local duration_of_time=`printf "%10s" "$5"`
  # local log_duration_of_time="$5"
  local log_role_name="$6"
  local log_task_name="$7"

  local log_level="WARNING"
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")

  if [[ -n "$log_role_name" ]]; then
    echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [role]:'$log_role_name'"
    # echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [role]:'$log_role_name'"  >> "$log_path"
  fi
  if [[ -n "$log_task_name" ]]; then
    echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_task_name'"
    # echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_task_name'" >> "$log_path"
    # echo "$log_message" >> "$log_path"
  fi
  # echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_title'"
  # echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_title'" >> "$log_path"
  # echo "$log_message" >> "$log_path"
  # echo "" >> "$log_path"
}

function _log_error() {
  local log_command="$1"
  local log_path="$2"
  local log_message="$3"
  local log_exit_code="$4"
  local duration_of_time=`printf "%10s" "$5"`
  # local log_duration_of_time="$5"
  local log_role_name="$6"
  local log_task_name="$7"

  local log_level="ERROR"
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")

  if [[ -n "$log_role_name" ]]; then
    echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [role]:'$log_role_name'"
    # echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [role]:'$log_role_name'"  >> "$log_path"
  fi
  if [[ -n "$log_task_name" ]]; then
    echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_task_name'"
    # echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_task_name'" >> "$log_path"
    # echo "$log_message" >> "$log_path"
  fi
  # echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_title'"
  # echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_title'" >> "$log_path"
  # echo "$log_message" >> "$log_path"
  # echo "" >> "$log_path"
}

function _log_fatal() {
  local log_command="$1"
  local log_path="$2"
  local log_message="$3"
  local log_exit_code="$4"
  local duration_of_time=`printf "%10s" "$5"`
  # local log_duration_of_time="$5"
  local log_role_name="$6"
  local log_task_name="$7"
  local log_sub_cmd_name="$8"

  local log_level="FATAL"
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  
  if [[ -z "$log_exit_code" ]]; then
    log_exit_code="1"
  fi
  if [[ -n "$log_role_name" ]]; then
    echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [role]:'$log_role_name'"
    # echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [role]:'$log_role_name'"  >> "$log_path"
  fi
  if [[ -n "$log_task_name" ]]; then
    echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_task_name'"
    # echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_task_name'" >> "$log_path"
    # echo "$log_message" >> "$log_path"
  fi
  if [[ -n "$log_sub_cmd_name" ]]; then
    echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [sub-command]:'$log_sub_cmd_name'"
  fi
  # echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_title'" >> "$log_path"
  # echo "$log_message" >> "$log_path"
  # echo "" >> "$log_path"
}

function _log_skip() {
  local log_command="$1"
  local log_path="$2"
  local log_role_name="$3"
  local log_task_name="$4"

  local start_time=$(TZ='Asia/Seoul' date +%s.%N)
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  local end_time=""
  local log_exit_code="0"
  local command_output=""
  local duration_of_time=`printf "%9s" "$5"`
  local log_level="SKIP"
  local log=""

  local log_level="SKIP"
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")

  if [[ -n "$log_role_name" ]]; then
    echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [role]:'$log_role_name'"
  fi
  if [[ -n "$log_task_name" ]]; then
    echo "[$timestamp] [Duration of time]:'$duration_of_time' [$log_level] [code]:'$log_exit_code' [task]:'$log_task_name'"
  fi

}