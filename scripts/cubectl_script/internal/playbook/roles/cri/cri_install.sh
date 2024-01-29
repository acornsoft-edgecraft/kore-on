#!/bin/sh
# set -e
## Container runtime interface 구성

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

# Create containerd directory
function handler_crate_containerd_dir() {
  local path="/etc/containerd/certs.d/${REGISTRY_DOMAIN}"
  printf "%s %s %s" "mkdir -p " "$path" "&& chmod 0755 $path"
}

# Create containerd directory
function handler_crate_containerd_dirs() {
  local path=`_conf_cri_dirs | tr "\n" " "`
  printf "%s %s %s" "mkdir -p " "$path" "&& chmod 0755 $path"
}

# Create containerd directory
function handler_copy_containerd_config_file() {
  local path="/etc/containerd/config.toml"
  local dir="${DATA_ROOT_DIR//\//\\\/}"
  local image_pause="${IMAGE_PAUSE_VERSION//\//\\\/}"

  local result=`sed \
  -e "s/^root = .*/root = \"${dir}\/containerd\"/" \
  -e "s/sandbox_image = .*/sandbox_image = \"${REGISTRY_DOMAIN}\/${image_pause}\"/" \
  <<< "$(_conf_cri_config_toml)"`

  echo "${result}" > "$(current_dir)/tmp/config.toml"

  printf "%s %s" "cp -f $(current_dir)/tmp/config.toml $path" "&& chmod 0644 $path"
}

# Create containerd data directory
function handler_crate_containerd_data_dir() {
  local path="${DATA_ROOT_DIR}/containerd"
  printf "%s %s %s" "mkdir -p " "$path" "&& chmod 0755 $path"
}

# Create containerd hosts.toml
function handler_create_containerd_hosts_toml() {
  local path="/etc/containerd/certs.d/"
  local file_name="hosts.toml"
  local delimiter=" && "
  local commands="$(_conf_script)"
  readarray -t list <<< "$(_conf_cri_dirs | sed "1i${REGISTRY_DOMAIN}" | sed "s/${path//\//\\\/}"//)"

  local cnt="${#list[@]}"

  for (( i=0; i<=${cnt}-1; i++)); do
    local server_url="https://${list[$i]}"
    local host_url="https://${list[$i]}"
    local ca_path="/etc/docker/certs.d/${REGISTRY_DOMAIN}/ca.crt"

    if [[ "$i" != 0 ]]; then
      host_url="${REGISTRY_DOMAIN}/v2/${list[$i]}"
    fi
    
    local data=`sed \
    -e "s/SERVER_URL/${server_url//\//\\\/}/" \
    -e "s/HOST_URL/${host_url//\//\\\/}/" \
    -e "s/CA_PATH/${ca_path//\//\\\/}/" \
    <<< "$(_conf_containerd_hosts_toml)"`

    if [[ "$i" == 0 ]]; then
      data=`echo "${data}" | sed -e "/capabilities/d" -e "/override_path/d"`
    fi
    echo "${data}" > "$(current_dir)/tmp/${list[$i]}_hosts.toml"

    commands+=$'\n'
    commands+="cp -f $(current_dir)/tmp/${list[i]}_hosts.toml ${path}${list[i]}/hosts.toml && chmod 0644 ${path}${list[i]}/hosts.toml"
  done


  echo "$commands" > "$(current_dir)/tmp/containerd_commands.sh"

  printf "%s" "sh $(current_dir)/tmp/containerd_commands.sh"
}

# Configure crictl.yaml
function handler_containerd_crictl_yaml() {
  local path="/etc/crictl.yaml"
  printf "%s %s" "echo \"$(_conf_containerd_crictl_yaml)\" > $path" "&& chmod 0644 $path"
}

## ===== [ Handlers End ] =====

## ===== [ Tasks ] =====
# Install containerd
function install_containerd() {
  local task_name="Install containerd"

  local res=""
  local exit_code=""
  local skip_run_index=()
  local skip_task_index=()

  ## command here
  local commands=(
    "yum install -y --nogpgcheck ${CONTAINERD_IO}"
    "$(handler_crate_containerd_dir)"
    "$(handler_crate_containerd_dirs)"
    "$(handler_crate_containerd_data_dir)"
    "$(handler_create_containerd_hosts_toml)"
    "$(handler_copy_containerd_config_file)"
    "systemctl daemon-reload && systemctl enable --now containerd.service && systemctl restart containerd.service"
    "$(handler_containerd_crictl_yaml)"
    )
  local commands_description=(
    "Install containerd"
    "Create containerd directory"
    "Create containerd directory for container images"
    "Create containerd data directory"
    "Add containerd config for mirrors"
    "Copy containerd config file"
    "Enable containerd service"
    "Configure crictl.yaml"
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
  source "$(current_dir)/internal/pkg/utils/config.sh"
  source "$(current_dir)/internal/pkg/config/cri_conf.sh"
  source "$(current_dir)/config/add_node_env.rc"

  ## ===== [ Constants and Variables ] =====
  local log_path="$(current_dir)/logs/add_node.log"
  local IFS="\t\n"
  local status=""

  ## Tasks: Install containerd
  install_containerd

}

main "${@}"