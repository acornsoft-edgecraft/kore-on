#!/bin/sh
## 모든 스크립트에서 호출되어야 하는 common 스크립트.

set -e
set -o errexit

set -E
set -o errtrace

set +H
set +o histexpand

## ===== [ Constants and Variables ] =====
pkg_dir="$1"

## ===== [ includes ] =====
source "${pkg_dir}/utils/logger.sh"


## ===== [ functions ] =====
# 소수점 뺄셈 연산을 위한 함수
function subtract_decimal() {
    local result
    local result=$(awk "BEGIN {printf \"%.2f\", $1 - $2}")
    echo "$result"
}

# 소수점 나눗셈 연산을 위한 함수
function divide_decimal() {
    local result
    local result=$(awk "BEGIN {printf \"%.2f\", $1 / $2}")
    echo "$result"
}

## 실행 시간 계산
function _duration() {
  local duration=""
  local duration_sec=$(subtract_decimal "$end_time" "$start_time")
  local num2=60
  local result=$(awk -v n1="$duration_sec" -v n2="$num2" 'BEGIN { if (n1 > n2) print 1; else print 0 }')

  if [[ "$result" -eq 1 ]]; then
    duration=$(divide_decimal "$duration_sec" "60")
    duration+="분"
  else
    duration=$(echo "$duration_sec")
    duration+="초"
  fi
  echo "$duration"
}

## Response - respose data 처리
# args: $1
# $1: 실행 결과값
function _response() {
  local data="$1"
  if [[ -n `grep "\[FATAL\]" <<< "$data"` ]]; then
    res_code="1"
    res_log="$data"
  else
    # log format과 data를 분리 출력
    res_code="0"
    res_log=`grep "\[INFO\]" <<< "$data"`
    res_data=`echo "$data" | sed "/\[INFO\]/d" `
  fi
}

## Clear Response - 전역변수 clear
function _clear_response() {
  res_code=""
  res_log=""
  res_data=""
}

## Response Code 실패 처리 - 종료
# args: $1
# $1: 실행 결과 코드
function _check_res_code() {
  local code="$1"
  if [[ "$code" == "1" ]]; then
    exit 1
  fi
}

## Tasks 진행 상태 저장
# args: $1, $2, $3, $4, $5
# $1: 실행 결과 코드
# $2: 실행 Task name
# $3: 실행 Task index
# $4: skip Task index
# $5: 실행 Task name description
function _save_task_status() {
  local code="$1"
  local task_name="$2"
  local task_index="$3"
  IFS='/' read -ra skip_index <<< "$4"
  IFS='|' read -ra task_discription <<< "$5"
  # local task_discription="$5"
  local save_path="${INSTALL_DIR}"

  local cnt=1

  if [[ "$code" == 0 ]]; then
    printf "%-8s %-40s %s\n" "[RUN]" "${task_name}::${task_index}" "$task_discription" >> "${save_path}/config/task_status.index"
    for index in "${skip_index[@]}"; do
      printf "%-8s %-40s %s\n" "[SKIP]" "${task_name}::${index}" "${task_discription[$cnt]}" >> "${save_path}/config/task_status.index"
      cnt=$((cnt + 1))
    done
    # for (( i=0; i<=${cnt}-1; i++ )); do
    #   printf "%-8s %-40s %s\n" "[SKIP]" "${task_name}::${index}" "$task_discription" >> "${save_path}/config/task_status.index"
    # done
  fi
}

## Tasks 진행 상태
# args: $1, $2, $3
# $1: 실행 Task name
# $2: 실행 Task total index
# $3: 저장 위치
function _get_task_status() {
  local task_name="$1"
  local task_total_len="$2"
  local save_path="${INSTALL_DIR}/config/task_status.index"
  local skip_task=()

  if [[ -e "${save_path}" ]]; then
    for (( i=0; i<=${task_total_len}-1; i++ )); do
      if [[ -n `grep "${task_name}::$i" "${save_path}"` ]]; then
        skip_task+=("$i")
      fi
    done
  fi

  echo "${skip_task[@]}"
}

## Skip task - 진행상태별 성공한 command skip
# args: $1, $2
# $1: task index 값
# $1: 저장 위치
function _skip_task() {
  local task_index="$1"
  local save_path="${INSTALL_DIR}"

  if [[ -e "${save_path}" && -n `grep "$task_index" "${save_path}/config/task_status.index"` ]]; then
    echo true
  else
    echo false
  fi
}

## Skip run - 명령어 실행 skip
# args: $1
# $1: index 배열 값
function _skip_run() {
  local data=("$@")
  local total_len="${#data[@]}"
  local key=""
  local skip_run_index=()

  # key, value 분리
  for (( i=0; i<=${total_len}-1; i++ )); do
    key="${data[0]}"
    if [[ "$i" > "0"  ]]; then
      skip_run_index+=("${data[$i]}")
    fi
  done

  # return
  for i in "${skip_run_index[@]}"; do
    if [[ "$i" == "$key" ]]; then
      echo true
      return
    fi
  done
  echo false
}

## Skip index - 배열값 대신 구분자 사용
# args: $1
# $1: index 배열 값
function _skip_index() {
  local index="$@"
  local skip_index=""
  local total_len="${#index[@]}"
  local delimiter="/"

  for (( i=0; i<=${total_len}-1; i++ )); do
    skip_index+="${index[$i]}"
    if [[ $i != ${total_len}-1 ]]; then
      skip_index+="$delimiter"
    fi
  done
  echo "$skip_index"
}

## command 실행 - root 권한 사용
# args: $1, $2, $3, $4
# $1: 실핼 명형어
# $2: 로그 디렉토리
# $3: play_book role name
# $4: task name
# $5: sub command name
function _run() {
  local command="$1"
  local log_path="$2"
  local role_name="$3"
  local task_name="$4"
  local sub_cmd_name="$5"

  local start_time=$(TZ='Asia/Seoul' date +%s.%N)
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  local end_time=""
  local exit_code=""
  local command_output=""
  local duration_of_time=""
  local log_level=""
  local log=""

  # sudo -H -S -n -u root /bin/sh -c "$command"
  if command_output=`sudo -H -S -n -u root /bin/sh -c "$command" 2>&1`; then
    end_time=$(TZ='Asia/Seoul' date +%s.%N)
    duration_of_time="$(_duration "$end_time" "$start_time")"
    exit_code="$?"
    log=`_log_info "$command" "$log_path" "$command_output" "$exit_code" "$duration_of_time" "$role_name" "$task_name" "$sub_cmd_name"`
    echo "$log" | tee -a "$log_path"
    echo "$command_output" | tee -a "$log_path"
    res_code="0"
  else
    end_time=$(TZ='Asia/Seoul' date +%s.%N)
    duration_of_time="$(_duration "$end_time" "$start_time")"
    exit_code="$?"
    log=`_log_fatal "$command" "$log_path" "$command_output" "$exit_code" "$duration_of_time" "$role_name" "$task_name" "$sub_cmd_name"`
    # log=`_log_fatal "$command" "$log_path" "$command_output" "$exit_code" "$duration_of_time" "$role_name" "$task_name "$sub_cmd_name"`
    echo "$log" | tee -a "$log_path"
    echo "[command]'$command'" | tee -a "$log_path"
    echo "$command_output" | tee -a "$log_path"
    res_code="1"
    # exit 1
  fi
}

## Clear tmp files
# args: $1, $2, $3
# $1: 로그 디렉토리
# $2: play_book role name
# $3: tmp 디텍토리
function _clear_tmp() {
  local path="$1"
  local log_path="$2"
  local command="rm -f ${path}/tmp/*"
  if command_output=`sudo -H -S -n -u root /bin/sh -c "$command" 2>&1`; then
    echo "${path}/tmp directory clear succesed!" | tee -a "$log_path" > /dev/null
  else
    echo "${path}/tmp directory clear failed!" | tee -a "$log_path" > /dev/null
  fi
}

function _get_first_ip() {
  local cidr=$1

  # CIDR 표기법에서 네트워크 주소 및 서브넷 마스크 얻기
  local network_address=$(echo "$cidr" | cut -d'/' -f1)
  local subnet_mask=$(echo "$cidr" | cut -d'/' -f2)

  # 서브넷 마스크를 이용하여 네트워크 주소를 계산
  IFS=. read -r -a network_octets <<< "$network_address"
  IFS=. read -r -a subnet_mask_octets <<< "$(printf "1%.0s" $(seq 1 "$subnet_mask"))"

  # 네트워크 주소를 서브넷 마스크로 AND 연산하여 첫 번째 IP 주소 계산
  local first_ip_octets=()
  for i in {0..3}; do
      first_ip_octets+=("$((network_octets[i] & subnet_mask_octets[i]))")
  done

  # 첫 번째 IP 주소를 출력
  echo "${first_ip_octets[*]}" | tr ' ' '.'
}


## Check 특수문자 in sed
# args: $1
# $1: 문자열
function _check_sed_str() {
  local str="$1"

  sed -e "s/\([!@#$%^&*/]\)/\1\\\2/g" <<< "${str}"
}