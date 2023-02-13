# **Kore-on**

## **구성도**

![airgap_install_archtecture](./assets/airgap_install_architecture.png)

## **요구사항**

- docker v19.03.15 이상
- Ubuntu 18.04, 20.04
- CentOS/RHEL 7, 8
- SSH KEY

## **Air-gap 설치**

> 폐쇄망 설치 전 사전 준비가 있습니다.
>
> koreonctl을 통해 클러스터 구성을 진행합니다.
>
> 환경구성 예시
> >
> > - online node (prepare-airgap 용도, 인터넷환경)
> > - bastion (폐쇄망)
> > - harbor (내부 이미지 레지, 폐쇄망)
> > - air-gap[1:2] (마스터1 워커2, 폐쇄망)
