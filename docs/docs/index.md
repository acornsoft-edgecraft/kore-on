# Kore-On

## Requirements

- docker
- Ubuntu 20.04
- SSH KEY

## Online Install

1. `apt-get update`

2. `apt-get upgrade -y`

3. docker install
   - `apt-get install -y docker.io`
   - Error
     1. Got permission denied while trying to connect to the Docker daemon socket
        - `sudo usermod -a -G docker $USER`
        - reconnect
        - `id`
        - find docker

4. `curl -LO https://github.com/acornsoft-edgecraft/kore-on/releases/download/**[last version]**/koreonctl-linux-amd64`
   - ex) `curl -LO https://github.com/acornsoft-edgecraft/kore-on/releases/download/v1.3.0/koreonctl-linux-amd64`

5. `chmod +x koreonctl-linux-amd64`

6. `cp koreonctl-linux-amd64 /usr/bin/korectl`

7. `korectl init`

8. edit koreon.toml
    - example online toml file link

9. `korectl create -p [SSH KEY] -u [USERNAME]`

## Verification

Going to controllplane node and do it the root.
Use that anything kubernetes CLI.

If want you general user?

```bash
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```
