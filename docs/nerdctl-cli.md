# Containerd command line used nerdctl

nerdctl은 'Docker-compatible CLI for containerd' 즉, docker client와 호환되는 명령어를 가지고 있는 도구입니다. 이름은 containerd의 뒤 4글자, 'nerd'와 'ctl' 을 합쳐 명명하였고 docker 와 거의 동일하게 사용하실 수 있습니다.

## nerdctl cli 사용법 for edgecraft

- Online Network 일때


```sh
## image pull
$ nerdctl pull ghcr.io/acornsoft-edgecraft/kore-on:latest

## run edgecraft installer container
$ nerdctl run --privileged -it ghcr.io/acornsoft-edgecraft/kore-on:latest

```