#!/bin/sh
## local repository를 구성

## 실행 디렉토리 절대경로 찾기
function current_dir() {
  ## scripts 실행 디렉토리 절대 경로 찾기
  local scripts_path=`dirname "$(readlink -f "$0")"`
  local depth=4

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

## ===== [ Tasks ] =====
# Backup Repo
function backup_repo() {
  local items=("/etc/apt/sources.list.d" "/etc/apt/sources.list")
  local current_date=$(date +"%Y%m%d%H%M")
  local total_len="${#items[@]}"

  for (( i=0; i<=${total_len}-1; i++ )); do
    ## command here
    # _run "cp -R ${items[i]}" "$logs_dir/add_node.log"
    echo "${items[i]}"
  done
}

# Backup Disable Repo list
function disable_repo() {
  local items=("/etc/apt/sources.list.d" "/etc/apt/sources.list")
  local total_len="${#items[@]}"

  # Remove apt repository
  ## command here
  # _run "rm -rf ${items[0]}"
  echo "${items[0]}"

  # Create apt repository
  ## command here
  # _run "mkdir -p ${items[0]}"
  echo "${items[1]}"

  # Replace apt repository
  ## command here
  # _run 'sed -i "s/^deb/#deb/g" ${items[1]}'
  echo "${items[1]}"
}
## ===== [ Tasks End ] =====

main() {
  ## ===== [ includes ] =====
  source "$(current_dir)"/pkg/utils/common.sh "$(current_dir)/internal/pkg"
  # source "$pkg_dir"/utils/logger.sh


  ## Tasks: Backup Repo
  backup_repo
  
  ## Tasks: Disable Repo list
  disable_repo

  ## Tasks: Add Local Repo

}

main "${@}"
