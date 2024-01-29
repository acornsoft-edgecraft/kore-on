#!/bin/sh
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

# sed 특수문자 입력
function sed_check_str() {
  echo "${LOCAL_REPO_URL//\//\\\/}"
}

## ===== [ Handlers End ] =====

## ===== [ Tasks ] =====
# Backup Repo
function backup_repo() {
  local task_name="Backup repository"
  local res=""
  local exit_code=""
  local skip_run_index=()
  local skip_task_index=()

  ## command here
  local current_date=$(date +"%Y%m%d%H%M")

  local commands=(
    "cp -R /etc/yum.repos.d /etc/yum.repos.d_back_${current_date}"
    )
  local commands_description=(
    "Backup Repository"
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

        ## 예외 처리 - START
        ## 예외 처리 - END

        _response "$res"

        ## Data Handler -- START
        ## Data Handler -- END

        ## response 결과 출력
        # 실패시 출력 후 종료
        # log 출려 후 response clear
        echo "$res_log"
        _save_task_status "$res_code" "$task_name" "$i" "$skip_index" "${commands_description[$i]}"
        _check_res_code "$res_code"
        _clear_response
      fi
    fi
  done
}

# Backup Disable Repo list
function disable_repo() {
  local task_name="Disable default repository"

  local res=""
  local exit_code=""
  local skip_run_index=()
  local skip_task_index=()

  ## command here
  local commands=(
    "subscription-manager config --rhsm.manage_repos=0"
    )
  local commands_description=(
    "disables the automatic management of repositories"
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
        ## Data Handler -- END

        ## response 결과 출력
        # 실패시 출력 후 종료
        # log 출려 후 response clear
        echo "$res_log"
        _save_task_status "$res_code" "$task_name" "$i" "$skip_index" "${commands_description[$i]}"
        _check_res_code "$res_code"
        _clear_response
      fi
    fi
  done
}

# Add local repository
function add_local_repo() {
  local task_name="Add local repository"

  local res=""
  local exit_code=""
  local skip_run_index=()
  local skip_task_index=()

  local local_repo_value=""
  
  # local_repo_value=`_conf_local_repo | sed "s/LOCAL_REPO_URL/$LOCAL_REPO_URL/g"`
  local_repo_value=`_conf_local_repo | sed "s/LOCAL_REPO_URL/$(sed_check_str)/g"`
  # command here
  local commands=(
    "echo \"$local_repo_value\" > /etc/yum.repos.d/local-repo.repo"
    "yum clean metadata"
  )
  local commands_description=(
    "Create local-repo"
    "Clean metadata"
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
        res=`_run "${commands[i]}" "$log_path" "" "${task_name}::${commands_description[i]}"`
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
  # echo "$local_repo_value" > "$(current_dir)/logs/test.log"
}

## ===== [ Tasks End ] =====

main() {
  ## ===== [ includes ] =====
  source "$(current_dir)/internal/pkg/utils/common.sh" "$(current_dir)/internal/pkg"
  source "$(current_dir)/internal/pkg/utils/config.sh"
  source "$(current_dir)/internal/pkg/utils/logger.sh"
  source "$(current_dir)/config/add_node_env.rc"
  source "$(current_dir)/config/os_release.rc"

  ## ===== [ Constants and Variables ] =====
  local log_path="$(current_dir)/logs/add_node.log"
  local IFS="\t\n"

  # ## Tasks: Backup Repo
  backup_repo
  
  ## Tasks: Disable Repo list
  disable_repo

  ## Tasks: Add Local Repo
  add_local_repo
}

main "${@}"
