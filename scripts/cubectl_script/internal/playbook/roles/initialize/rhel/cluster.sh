#!/bin/sh
set -e
## initialize 구성 sub script - cluster

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

# Load kernel modules for Calico
function handler_modprobe_calico() {
  printf "%s %s" "modprobe" "$(_conf_cluster_calico_kernel)" | tr "\n" " " 
}

# Persist kernel modules for Calico
function handler_calico_persist() {
  local path="/etc/modules"
  printf "%s %s %s %s" "echo" "\"$(_conf_cluster_calico_kernel)\"" " > $path" "&& chmod 0644 $path"
}

# Persist br_netfilter module
function handler_br_netfilter_persist() {
  local path="/etc/modules-load.d/br_netfilter.conf"
  printf "%s %s %s %s" "echo" "br_netfilter" " > $path" "&& chmod 0644 $path"
}

# Modprobe Kernel Module for IPVS
function handler_modprobe_ipvs() {
  printf "%s %s" "modprobe" "$(_conf_modpwobe_ipvs)" | tr "\n" " " 
}

# Persist ip_vs modules
function handler_ipvs_persist() {
  local path="/etc/modules-load.d/ipvs.conf"
  printf "%s %s %s %s" "echo" "\"$(_conf_modpwobe_ipvs)\"" " > $path" "&& chmod 0644 $path"
}

# Forwarding IPv4 and letting iptables see bridged traffic
function handler_crate_sysctl_k8s() {
  local path="/etc/sysctl.d/k8s.conf"
  printf "%s %s %s %s" "echo" "\"$(_conf_sysctl_k8s)\"" " > $path" "&& chmod 0644 $path && sysctl --system"
}

# Create docker cert directory
function handler_crate_docker_cert_dir() {
  local path="/etc/docker/certs.d/${REGISTRY_DOMAIN}"
  printf "%s %s %s" "mkdir -p " "$path" "&& chmod 0755 $path"
}

# Persist ip_vs modules
function handler_sysctl_dind() {
  local path="/etc/sysctl.d/99-dind.conf"
  printf "%s %s %s %s" "echo" "\"$(_conf_sysctl_dind)\"" " > $path" "&& chmod 0755 $path && sysctl --system"
}

## ===== [ Handlers End ] =====

## ===== [ Tasks ] =====
# Kernel configuration
function initialize_cluster() {
  local task_name="Initialize Cluster"

  local res=""
  local exit_code=""
  local service_list=""
  local skip_run_index=()
  local skip_task_index=()

  ## command here
  local commands=(
    "swapoff -a && sed -i '/swap/s/^/#/' /etc/fstab"
    "systemctl list-unit-files --type=service"
    "systemctl disable --now firewalld.service"
    "yum install -y NetworkManager"
    "mkdir -p /etc/NetworkManager/conf.d && chmod 0755 /etc/NetworkManager/conf.d"
    "echo \"$(_conf_networkmanager_calico)\" > /etc/NetworkManager/conf.d/calico.conf && chmod 0644 /etc/NetworkManager/conf.d/calico.conf"
    "echo \"$(_conf_networkmanager)\" > /etc/NetworkManager/conf.d/90-dns-none.conf && chmod 0644 /etc/NetworkManager/conf.d/90-dns-none.conf"
    "systemctl reload NetworkManager.service"
    "$(handler_modprobe_calico)"
    "$(handler_calico_persist)"
    "modprobe br_netfilter"
    "$(handler_br_netfilter_persist)"
    "$(handler_modprobe_ipvs)"
    "$(handler_ipvs_persist)"
    "yum install -y ipvsadm ipset"
    "$(handler_crate_sysctl_k8s)"
    "$(handler_crate_docker_cert_dir)"
    "systemctl enable --now systemd-resolved.service"
    "yum install -y nfs-utils"
    "$(handler_sysctl_dind)"
    )
  local commands_description=(
    "Swapoff && Remove swapfile from /etc/fstab"
    "Get system services"
    "Disable a service and stop for firewalld.service"
    "Install NetworkManager"
    "Ensure NetworkManager conf.d dir"
    "Prevent NetworkManager from managing Calico interfaces"
    "Disable NetworkManager DNS processing"
    "Reload NetworkManager.service"
    "Load kernel modules for Calico"
    "Persist kernel modules for Calico"
    "Enable br_netfilter module"
    "Persist br_netfilter module"
    "Modprobe Kernel Module for IPVS"
    "Persist ip_vs modules"
    "Install package for ipvs"
    "Initialize | Forwarding IPv4 and letting iptables see bridged traffic"
    "Create docker cert directory"
    "Enable service systemd-resolved"
    "Install nfs-utils"
    "Run the Docker daemon as a non-root user"
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
    local skip_description=""

    if [[ "$skip_task" == true ]]; then
      _log_skip "${commands[i]}" "$log_path" "" "${task_name}::${commands_description[$i]}"
    else
      if [[ "$skip_run" == false  ]]; then
        res=`_run "${commands[i]}" "$log_path" "" "${task_name}::${commands_description[$i]}"`

        ## 예외 처리 - START
        ## 예외 처리 - END

        _response "$res"

        ## Data Handler -- START
        skip_description="${commands_description[$i]}|"

        # ServiceList 저장 후 결과 처리
        if [[ "$i" == 1 && "$res_code" == 0 ]]; then
          service_list="$res_data"
        fi

        if [[ "$i" == 1  ]]; then
          if [[ -z `grep "firewalld.service" <<< "$service_list"` ]];then
            skip_run_index+=(2)
            skip_index+="2/"
            skip_description+="${commands_description[2]}|"
          fi
          if [[ -n `grep "NetworkManager.service" <<< "$service_list"` ]]; then
            skip_run_index+=(3)
            skip_index+="3/"
            skip_description+="${commands_description[3]}|"
          fi
          if [[ -n `grep "systemd-resolved.service" <<< "$service_list" | grep "enabled"` ]]; then
            skip_run_index+=(17)
            skip_index+="17/"
            skip_description+="${commands_description[18]}|"
          fi
        fi

        # KUBE-PROXY-MODE == 'ipvs'
        if [[ "$i" == 11 && "$KUBE_PROXY_MODE" != "ipvs" ]]; then
          skip_run_index+=(12 13 14)
          skip_index+="12/13/14/"
          skip_description="${commands_description[12]}|"
          skip_description+="${commands_description[13]}|"
          skip_description+="${commands_description[14]}|"
        fi

        # sysctl-dind
        if [[ "$i" == 19 && "0" -eq `awk -v n1="$VERSION_ID" -v n2="8" 'BEGIN { if (n1 > n2) print 1; else print 0 }'` ]]; then
          skip_run_index+=(20)
          skip_index+="20/"
          skip_description+="${commands_description[20]}|"
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
  local status=""
 
  ## Tasks: initialize_cluster
  initialize_cluster

  ## 스크립트 종료 상태 변경
  if [[ "$status" == "0" ]]; then
    exit 0
  fi
}

main "${@}"