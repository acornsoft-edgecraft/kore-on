# 단일 노드 추가 스크립트

## 작업 계획서

- 1. 마스터 노드에서 필요한 파일을 가져온다. (ca.crt, ca.key, acloud-client-kubeconfig)
- 2. 작업노드에서 설치 파일과 필수 파일을 준비 한다. (directory tree 확인)
- 3. 작업노드에서 노드추가 실행 쉘스크립트를 실행 한다.
- 4. 작업노드에서 로그를 확인 한다.
- 5. 작업노드에서 노드 추가 완료후 클러스터 확인 한다.

실행 명령어:
```sh
## 1. 마스터 노드에서 필수 파일을 가져온다.
cat /etc/kubernetes/pki/ca.crt
cat /etc/kubernetes/pki/ca.key
cat /etc/kubernetes/acloud/acloud-client-kubeconfig

## 2-1. 작업노드에서 설치 파일과 필수 파일을 준비 한다.
# Directory tree 설명을 참고해서 준비된 필수 파일과 실행 파일들이 해당 디렉토리에 있는지 확인 한다.
# ~/cubectl_script/config/pki/ca.crt
# ~/cubectl_script/config/pki/ca.key
# ~/cubectl_script/config/acloud-client-kubeconfig

## 2-2. 작업노드에서 설치에 필요한 값을 작성 한다.
# 아래 설명되어 있는 환경변수 값 정의 항목을 참고해서 작성.
vi ~/cubectl_script/config/add_node_env.rc

## 3. 작업노드에서 설치 파일이 있는 디렉토리에서 명령어 실행
# sub-command parameter는 필수 항목
# ~/cubectl_script/cmd/cubectl_script.sh list 명령으로 sub-command list를 확인할 수 있다.
sh ~/cubectl_script/cmd/cubectl_script.sh add_node

## 4. 작업노드에서 상세 로그 확인
tail -f ~/cubectl_script/logs/add_node.log

## 5. 작업 노드에서 acloud-client-kubeconfig를 사용해서 클러스터 확인
# dnsutil을 daemonset으로 설치되어 있으면 노드 추가할때 자동배포 된다.
# 추가되 노드에 생성된 dnsutils로 nslookup 확인
# 참고: dnsutils배포는 마스터노드에서 /etc/kubernetes/addon/test/dnsutils.yaml 파일로 배포할 수 있다.
export KUBECONFIG=~/shell_scripts/config/acloud-client-kubeconfig
kubectl get nodes -o wide
kubectl exec -i -t dnsutils-xxxx -- nslookup kubernetes.default
```

## 필수사항

인증서를 생성 하기위해 ca.crt ca.key 파일이 있어야 합니다.
k8s cluster에 조인 할때 토큰을 생성을 위해 kubeconfig가 필요 합니다.

- 마스터 노드에서 ca.crt ca.key 값을 가져와서 지정된 경로에 저장해야 합니다.
- 마스터 노드에서 kubeconfig 값을 가져와서 지정된 경로에 저장해야 합니다.
- 파일명과 저장위치는 고정값으로 변경할 수 없습니다.

```sh
## 마스터 노드에서 ca.crt 와 ca.key 값을 받아와서 해당 위치에 저장한다.
# 인증서 생성시 사용
# 저장 위치: <쉘스크립트 위치>/config/pki/ca.crt
# 저장 위치: <쉘스크립트 위치>/config/pki/ca.key
cat /etc/kubernets/pki/ca.crt > [저장위치]
cat /etc/kubernets/pki/ca.crt > [저장위치]

## 마스터 노드에서 kubeconfig 값을 받아와서 해당 위치에 저장한다.
# node join 할때 token 생성에 사용
# 저장 위치: <쉘스크립트 위치>/config/acloud-client-kubeconfig
cat /etc/kubernetes/acloud/acloud-client-kubeconfig > [저장위치]
```

## 환경변수 값 정의

> 환경변수 값은 설치시 필수로 필요한 요소들 입니다.
> 설치정보는 "INSTALL_DIR" 값의 경로에 저장 됩니다. 실행 단계별로 저장되어 만약 실패시 재설시 할때 완료단계를 skip 할수 있습니다.
> "LOCAL_REPO_URL" 값은 OS direcory명을 포함해야 합니다.

- "K8S_VERSION": k8s cluster version과  동일 해야 한다. (kubectl, kubelet, kubeadm 패키지 버전으로 사용된다)
- "CLUSTER_NAME": k8s cluster name과 동일해야 한다. (default:"kubernetes", kubelet extra args에서 사용)
- "REGISTRY_DOMAIN": 사설 레지스트리 ip address
- "LOCAL_REPO_URL": local repository의 URL (OS direcory명을 포함해야 한다)
- "INSTALL_DIR": 설치정보 저장 디렉토리 (default: /var/lib/cubectl)
- "DATA_ROOT_DIR": DATA 저장 디렠토리 (default: /data)
- "KUBE_PROXY_MODE": K8s kube-proxy 모드 설정 (default: "ipvs")
- "SERVICE_CIDR": K8s 서비스 ip 대역 설정 (default: "10.96.0.0/20")
- "CERT_VALIDITY_DAYS": 인증성 기간 설정 (default: "36500")
- "API_SANS": K8s apiserver 인증서에 ip 추가 (apiserver.crt에 추가됨)
- "CONTAINERD_IO": 설치 할 containerd package명 (default: containerd.io-1.6.18-3.1.el8)
- "IMAGE_PAUSE_VERSION": containerd 서비스가 사용하는 이미지명 (default: registry.k8s.io/pause:3.9)
- "NODE_POOL_MASTER": 마스터 ip (haproxy.cfg에 구성 된다)
- "LB_IP": 로드밸런서 IP 주소 (인증서 생성시 추가된다.)
- "LB_PORT": 로드밸런서 PORT 번호 (kubeconfig에 추가 된다.)
- "NODE_IP": 노드 IP 주소 (인증서 생성에 추가 된다)

```sh
## 아래 내용은 예제입니다:
## DEFAULT
K8S_VERSION="v1.26.7"
CLUSTER_NAME="kubernetes"
REGISTRY_DOMAIN="192.168.88.219"
LOCAL_REPO_URL="http://192.168.88.219:8080/rhel8"
INSTALL_DIR="/var/lib/cubectl"
DATA_ROOT_DIR="/data"
KUBE_PROXY_MODE="ipvs"
SERVICE_CIDR="10.96.0.0/20"
CERT_VALIDITY_DAYS="36500"
API_SANS=(
  ""
)

## CRI
CONTAINERD_IO="containerd.io-1.6.18-3.1.el8"
IMAGE_PAUSE_VERSION="registry.k8s.io/pause:3.9"

## MASTER
NODE_POOL_MASTER=(
  "10.10.30.4"
  "10.10.30.175"
  "10.10.30.218"
)
LB_IP=""
LB_PORT=""

## NODE
NODE_IP="10.10.30.213"
```

## 디렉토리 트리 설명

- 디렉토리 구조
  
```sh
cubectl_script
├── cmd                               실행 명령어 디렉토리:
│   └── cubectl_script.sh               - 실행 명령어
├── config                            실행 설정 디렉토리:
│   ├── acloud-client-kubeconfig        - 마스터 노드의 acloud-client-kubeconfig 파일을 여기에 저장
│   ├── add_node_env.rc                 - 작업 노드 설정 파일
│   ├── os_release.rc                   - 작업 노드 OS 정보(자동 입력됨 - 변경 X)
│   ├── pki                             인증서 파일 디렉토리:
│   │   ├── ca.crt                        - 마스터 노드의 ca.crt 파일을 여기에 저장
│   │   └── ca.key                        - 마스터 노드의 ca.key 파일을 여기에 저장
│   └── task_status_index.template      - 설치 진행 정보 저장 파일
├── docs                              문서 디렉토리
│   └── README.md                       - README 파일
├── internal                          내부 소스 디렉토리(전체 진행 소스)
└── logs                              로그 디렉토리
│   └── add_node.log                    - 노드 추가 로그 파일
└── tmp                               임시 파일 저장 디렉토리(설치 진행에 필요한 임시파일이 저장되고 TASK가 완료되면 내용 삭제 된다.)

```