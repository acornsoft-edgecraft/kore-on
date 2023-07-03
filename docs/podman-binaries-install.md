# podman binaries and container images

> [참고] https://github.com/mgoltzsche/podman-static


## Binary installation on a host

```sh
## download binary file
$ curl -fsSL -o podman-linux-amd64.tar.gz https://github.com/mgoltzsche/podman-static/releases/latest/download/podman-linux-amd64.tar.gz
### Download a specific version:
$ VERSION=v4.5.1
$ curl -fsSL -o podman-linux-amd64.tar.gz https://github.com/mgoltzsche/podman-static/releases/download/$VERSION/podman-linux-amd64.tar.gz


### Install the binaries and configuration on your host after you've inspected the archive:
$ tar --strip-components=1 --no-overwrite-dir -C / -xzvf podman-linux-amd64.tar.gz
```

## troubleshooting

- Setup CNI networking
- [참고] https://podman.io/docs/installation#setup-cni-networking

```sh
## 오류 현상:
# contaner의 네트워크 인터페이스가 생성 되지 않느다.


## 해결 방법:
# step-1. Add configuration
$ sudo mkdir -p /etc/containers
$ sudo curl -L -o /etc/containers/registries.conf https://src.fedoraproject.org/rpms/containers-common/raw/main/f/registries.conf
$ sudo curl -L -o /etc/containers/policy.json https://src.fedoraproject.org/rpms/containers-common/raw/main/f/default-policy.json

# step-2. Network configuration
$ sudo mkdir -p /etc/cni/net.d
$ curl -qsSL https://raw.githubusercontent.com/containers/podman/main/cni/87-podman-bridge.conflist | sudo tee /etc/cni/net.d/87-podman-bridge.conflist
```