#!/bin/sh
set -e
## local repository를 구성

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

# Create log directory
function handler_crate_install_dir() {
  local path="${INSTALL_DIR}/config"
  printf "%s %s %s" "mkdir -p " "$path" "&& chmod 0755 $path"
}

# Create log directory
function handler_crate_log_dir() {
  local path="${DATA_ROOT_DIR}/logs"
  printf "%s %s %s" "mkdir -p " "$path" "&& chmod 0755 $path"
}
## ===== [ Handlers End ] =====

## ===== [ Tasks ] =====
# Backup Repo
function init_conf() {
  local task_name="Init Configrations"

  local res=""
  local exit_code=""
  local skip_run_index=()
  local skip_task_index=()

  ## command here
  local commands=(
    "$(handler_crate_install_dir)"
    "$(handler_crate_log_dir)"
    )
  local commands_description=(
    "Create $INSTALL_DIR"
    "Create $DATA_ROOT_DIR"
  )

  local total_len="${#commands[@]}"

  # get task status - 성공한 task는 skip 하기 위함
  skip_task_index=(`_get_task_status "$task_name" "$total_len"`)

  for (( i=0; i<=${total_len}-1; i++ )); do
    ## skip 설정
    # skip_run: 명령어 실행 skip
    # skip_task: 진행상태별 성공한 command skip
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
        if [[ ! -e "${INSTALL_DIR}/config/task_status.index" ]]; then
          cat "$(current_dir)/config/task_status_index.template" > "${INSTALL_DIR}/config/task_status.index"
        fi
        ## Data Handler -- END

        ## response 결과 출력
        # log 출력 후 response clear
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
  source "$(current_dir)/config/add_node_env.rc"

  ## ===== [ Constants and Variables ] =====
  local log_path="$(current_dir)/logs/add_node.log"
  local IFS="\t\n"

  # ## Tasks: Init Configrations
  init_conf
}

main "${@}"
