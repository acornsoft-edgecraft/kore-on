# Kore-On

## 구성도

![online_install_archtecture](./assets/online_install_architecture.jpeg)

## 요구사항

- docker v19.03.15 이상
- Ubuntu 20.04
- CentOS/RHEL 8
- SSH KEY

## 온라인 설치
> koreonctl을 통해 클러스터 구성을 진행합니다.
> root 유저에서 실행을 권장하며 설치 될 타켓 노드 외 별도의 Client에서 설치 koreonctl을 실행해야 합니다.
> 아래 샘플에서는 설치 될 타켓 노드를 master 노드 1대 와 worker 노드 2대의 설정으로 진행합니다.
 
<!-- to be! 이부분은 설치 Client OS별로 나누어야 할 듯  -->
1. apt 업데이트

        apt-get update
        apt-get upgrade -y
       

2. docker install
        koreonctl의 실행을 위해 실행 할 Client에 Docker 설치가 필요합니다. 다음 링크를 참조하여 Client OS에 맞는 도커를 설치합니다. 이미 설치 되었다면 다음 단계로 진행합니다.
        https://docs.docker.com/engine/install/
           
3. 설치 CLI Tool 인 koreonctl 을 다운로드 합니다

        curl -LO https://github.com/acornsoft-edgecraft/kore-on/releases/download/[last version]/koreonctl-linux-amd64

        ex) curl -LO https://github.com/acornsoft-edgecraft/kore-on/releases/download/v1.3.0/koreonctl-linux-amd64

4. 다운로드 한 설치파일에 실행 권한을 부여 합니다

        chmod +x koreonctl-linux-amd64

5. 실행 파일명 변경 및 위치 이동

        cp koreonctl-linux-amd64 /usr/bin/koreonctl

6. 설치 설정파일 koreon.toml 을 기본값으로 생성 합니다 

        koreonctl init

7. koreon.toml 파일을 클러스터 구성에 맞게 수정 합니다

        ```toml
        [koreon]
        # 하단의 해당하는 부분만 변경
        cluster-name = "testing-cluster"

        [node-pool.node]
        # 하단의 해당하는 부분만 변경
        ip = ["x.x.x.x","x.x.x.x","x.x.x.x"]

        [node-pool.master]
        # 하단의 해당하는 부분만 변경
        ip = ["x.x.x.x"]
        ```

8. 클러스터 설치 시작
        
        [SSH KEY PATH] - 설치 될 클러스터의 SSH 접근 key값을 설정합니다 
        [USERNAME] - 설치 될 클러스터 노드의 SSH 노드의 접속 user를 설정합니다 
        
        korectl create -p [SSH KEY PATH] -u [USERNAME]

## 검증

> 관리자 계정에서 kubetnetes CLI를 실행하여야한다.

controllplane node에서 관리자 계정이 아닌 일반 사용자가 Kubernetes CLI를 사용하기를 원하면 아래 명령어를 사용해 주세요

    mkdir -p $HOME/.kube
    sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
    sudo chown $(id -u):$(id -g) $HOME/.kube/config
