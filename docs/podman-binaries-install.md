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
