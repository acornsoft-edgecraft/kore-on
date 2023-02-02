# changelog

## Test

- serchTest

# koreonctl

- init
  - 이 명령은 시스템 정보 및 설치 정보를 설정할 수 있는 샘플 파일을 다운로드합니다.
- prepare-airgap
  - 에어갭 네트워크를 위한 쿠버네티스 클러스터 및 레지스트리 준비
- prepare-airgap download-archive
  - 보관 파일을 localhost로 다운로드
- create
- -p [key] -u [user]
  - 이 명령은 Kubernetes 클러스터 및 레지스트리를 설치합니다.
- destroy
  - 이 명령은 [Kubernetes 클러스터 / 레지스트리 / 저장소]를 삭제할 수 있습니다. 하위 명령을 사용하지 않으면 모두 삭제됩니다.
  - cluster
  - registry
  - storage
  - prepare-airgap

mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# 설치 순서 온라인

```bash
korectl init

```