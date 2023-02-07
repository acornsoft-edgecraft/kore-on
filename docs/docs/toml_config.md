# **toml 파일 설명**

## **[koreon]**

| Parameter | Description |
| :--- | :--- |
| **`cluster-name`** | **클러스터명** |
| **`install-dir`** | **설치 스크립트 디렉토리** |
| **`cert-validity-days`** | **인증서 유효기간 (default: 365)** |
| **`debug-mode`** | **Dry Run (not yet supported.)** |
| **`closed-network`** | **폐쇄망 선택** |
| **`local-repository-install`** | **로컬 Repo 설치 선택** |
| **`local-repository-port`** | **로컬 Repo 서비스 포트** |
| **`local-repository-archive-file`** | **로컬 Repo 패키지 아카이브 파일명** |
| **`local-repository-url`** | **로컬 Repo Url (default: registry-ip)** |

## **[kubernetes]**

| Parameter | Description | default |
| --- | --- | --- |
| **`version`** | **k8s 버전** | **latest** |
| **`container-runtime`** | **k8s cri** | **(only)containerd** |
| **`kube-proxy-mode`** | **k8s kube-proxy Mode** | **ipvs** |
| **`service-cidr`** | **k8s service network cidr** | **"10.96.0.0/20”** |
| **`pod-cidr`** | **k8s pod network cidr** | **"10.4.0.0/24”** |
| **`node-port-range`** | **k8s node port network range** | **"30000-32767”** |
| **`audit-log-enable`** | **k8s audit log enabled** | **true** |
| **`api-sans`** | **K8s apiserver에 SAN** | **master** |

## **[kubernetes.etcd]**

