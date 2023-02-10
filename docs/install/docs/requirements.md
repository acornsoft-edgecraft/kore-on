# **요구 사항**

Kore-on은 해당하는 최소한의 조건이 있습니다.

## **전제조건**

- control plane node와 worker node의 User명이 같아야합니다.
- control plane의 제한은 1 ~ 7대이며 홀수로 구성해야합니다.
- CLI client에서 control plane node와 worker node로 SSH 접속이 가능해야 합니다.

## **운영체제**

Kore-on은 해당하는 operating system을 테스트 및 검증하여 지원하고 있습니다.

### **Linux**

- Ubuntu 18.04, 20.04(amd64)
- CentOS/RHEL 7, 8(amd64)

## **하드웨어**

Kore-on 설치를 위해 필요한 최소한의 하드웨어 요구 사항은 다음과 같습니다.

### CLI client

- 도커를 설치하고 실행할 수 있는 사양
- 64비트 커널 및 CPU 지원
- RAM : 최소 4GB

### control plane node, worker node

- 각 노드 별 k8s를 설치하고 실행할 수 있는 사양
- RAM : 최소 2GB
- CPU : 최소 2개 CPU(코어)

## **네트워크**

중요: CLI client에서 control plane node, worker node로 SSH 접속이 가능하여야 합니다.

### **Port**

#### **CLI client, All node**

| Protocol  | Port Range  |
| ---       | ---         |
| SSH       | 22          |

#### **Control plane(s)**

| Protocol  | Direction | Port Range  | Purpose                 | Used By               |
| ---       | ---       | ---         | ---                     | ---                   |
| TCP       | Inbound   | 6443        | Kubernetes API server   | All                   |
| TCP       | Inbound   | 2379-2380   | etcd server client API  | kube-apiserver, etcd  |
| TCP       | Inbound   | 10250       | Kubelet API             | Self, Control plane   |
| TCP       | Inbound   | 10259       | kube-scheduler          | Self                  |
| TCP       | Inbound   | 10257       | kube-controller-manager | Self                  |

#### **Worker node(s)**

| Protocol  | Direction | Port Range  | Purpose             | Used By             |
| ---       | ---       | ---         | ---                 | ---                 |
| TCP       | Inbound   | 10250       | Kubelet API         | Self, Control plane |
| TCP       | Inbound   | 30000-32767 | NodePort Servicest  | All                 |
