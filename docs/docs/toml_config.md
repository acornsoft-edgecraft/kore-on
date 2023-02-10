# **toml 파일 설명**

## **[koreon]**

- debug-mode는 아직 지원하지 않습니다.
- closed-network ~ local-repository-url은 폐쇄 네트워크를 선택할 때 필요합니다.
- local-repository-url은 domain을 사용할 때 입력합니다.

| Parameter                             | Description                         | default           |
| :---                                  | :---                                | :---              |
| **`cluster-name`**                    | cluster name                        | kubernetes        |
| **`install-dir`**                     | install script dir                  | /var/lib/kore-on  |
| **`cert-validity-days`**              | cert validity                       | 36500             |
| **`debug-mode`**                      | Dry Run                             | false             |
| **`closed-network`**                  | select air-gap                      | false             |
| **`local-repository-install`**        | select local rep install            |
| **`local-repository-port`**           | input local repository port         | 8080              |
| **`local-repository-archive-file`**   | local repository archive file path  |
| **`local-repository-url`**            | local repository archive url        | registry-ip       |

## **[kubernetes]**

| Parameter               | Description                 | default           |
| ---                     | ---                         | ---               |
| **`version`**           | k8s version                 | latest            |
| **`container-runtime`** | k8s cri                     | (only)containerd  |
| **`kube-proxy-mode`**   | k8s kube-proxy Mode         | ipvs              |
| **`service-cidr`**      | k8s service network cidr    | 10.96.0.0/20      |
| **`pod-cidr`**          | k8s pod network cidr        | 10.4.0.0/24       |
| **`node-port-range`**   | k8s node port network range | 30000-32767       |
| **`audit-log-enable`**  | k8s audit log enabled       | true              |
| **`api-sans`**          | K8s apiserver SAN           | master            |

## **[kubernetes.etcd]**

- external-etcd는 아직 지원하지 않습니다.

| Parameter           | Description                       | default                       |
| :---                | :---                              | :---                          |
| **`external-etcd`** | external ETCD cluster composition | false                         |
| **`ip`**            | external ETCD node ip             | control plane node ip address |
| **`private-ip`**    | external ETCD node private ip     | control plane node ip address |

## **[kubernetes.calico]**

| Parameter         | Description               | default |
| :---              | :---                      | :---    |
| **`vxlan-mode`**  | select vxlan-mode active  | false   |

## **[node-pool]**

- data-dir는 backup, docker, log, kubelet, etcd, k8s-audit, container가 저장되는 곳입니다.

| Parameter       | Description                 | default |
| :---            | :---                        | :---    |
| **`data-dir`**  | save data of directory path | /data   |
| **`ssh-port`**  | Node ssh port               | 22      |

## **[node-pool.master]**

- 만약 private-ip가 ip와 값이 같으면 생략이 가능합니다.

| Parameter             | Description                             | default                           |
| :---                  | :---                                    | :---                              |
| **`ip`**              | control plane node ip address           |                                   |
| **`private-ip`**      | control plane nodes private ip address  |                                   |
| **`isolated`**        | control plane nodes isolated            | false                             |
| **`haproxy-install`** | used internal load-balancer             | true                              |
| **`lb-ip`**           | load-balancer ip address                | control plane[0] node ip address  |
| **`lb-port`**         | load-balancer port                      | 6443                              |

## **[node-pool.node]**

- 만약 private-ip가 ip와 값이 같으면 생략이 가능합니다.

| Parameter         | Description                     | default |
| :---              | :---                            | :---    |
| **`ip`**          | worker node ip address          |         |
| **`private-ip`**  | worker nodes private ip address |         |

## **[private-registry]**

- registry-version은 개인 레지스트리를 설치할 때 사용되는 필수 항목입니다. 주 버전만 입력하면 보조 버전이 자동으로 마지막 버전을 선택합니다.
- registry-domain은 registry가 도메인을 사용한다면 입력합니다.

| Parameter                   | Description                     | default                             |
| :---                        | :---                            | :---                                |
| **`install`**               | private registry install        | false                               |
| **`registry-version`**      | private registry version        | latest                              |
| **`registry-ip`**           | used internal load-balancer     | true                                |
| **`private-ip`**            | load-balancer ip address        | control plane[0] node ip address    |
| **`registry-domain`**       | registry domain                 | 6443                                |
| **`data-dir`**              | private registry data directory | /data/harbor                        |
| **`registry-archive-file`** | registry archive file path      | “”                                  |
| **`public-cert`**           | public cert activate            | false                               |

## **[private-registry.cert-file]**

| Parameter                 | Description               | default |
| :---                      | :---                      | :---    |
| **`ssl-certificate`**     | ssl certificate path      |         |
| **`ssl-certificate-key`** | ssl certificate key path  |         |

## **[shared-storage]**

| Parameter         | Description                 | default       |
| :---              | :---                        | :---          |
| **`install`**     | NFS Server Installation     | false         |
| **`type`**        |                             |               |
| **`storage-ip`**  | storage node ip address     |               |
| **`private-ip`**  | storage node ip address     |               |
| **`volume-dir`**  | storage node data directory | /data/storage |
| **`nfs_version`** |                             |               |

## **[prepare-airgap]**

| Parameter               | Description                 | default |
| :---                    | :---                        | :---    |
| **`k8s-version`**       | kubernetes version          | latest  |
| **`registry-version`**  | private registry version    | latest  |
| **`registry-ip`**       | private registry ip address |         |
