getfiles:
	@echo "Get necessary files ..."
	curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.19.0/crictl-v1.19.0-linux-amd64.tar.gz -o ./internal/playbooks/cubescripts/files/crictl-v1.19.0-linux-amd64.tar.gz
	curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.20.0/crictl-v1.20.0-linux-amd64.tar.gz -o ./internal/playbooks/cubescripts/files/crictl-v1.20.0-linux-amd64.tar.gz
	curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.21.0/crictl-v1.21.0-linux-amd64.tar.gz -o ./internal/playbooks/cubescripts/files/crictl-v1.21.0-linux-amd64.tar.gz
	curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.22.0/crictl-v1.22.0-linux-amd64.tar.gz -o ./internal/playbooks/cubescripts/files/crictl-v1.22.0-linux-amd64.tar.gz
	curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.23.0/crictl-v1.23.0-linux-amd64.tar.gz -o ./internal/playbooks/cubescripts/files/crictl-v1.23.0-linux-amd64.tar.gz
	curl -L https://github.com/docker/compose/releases/download/1.29.2/docker-compose-Linux-x86_64 -o ./internal/playbooks/cubescripts/files/docker-compose-v1.29.2-Linux-x86_64
	curl -L https://github.com/etcd-io/etcd/releases/download/v3.4.16/etcd-v3.4.16-linux-amd64.tar.gz -o ./internal/playbooks/cubescripts/files/etcd-v3.4.16-linux-amd64.tar.gz
	curl -L https://github.com/goharbor/harbor/releases/download/v2.3.0/harbor-offline-installer-v2.3.0.tgz -o ./internal/playbooks/cubescripts/files/harbor-offline-installer-v2.3.0.tgz

# @REM Windows 64bit binary compile
build_windows_64:
	GOOS=windows & GOARCH=amd64 & go build -o koreon_windows_64_ctl.exe main.go

# @REM Windows 32bit binary compile
build_windows_32:
	GOOS=windows & GOARCH=386 & go build -o koreon_windows_32_ctl.exe main.go

# @REM Linux 64bit binary compile
build_linux_64:
	GOOS=linux & GOARCH=amd64 & go build -o koreon_linux_64_ctl main.go

# @REM Linux 32bit binary compile
build_linux_32:
	GOOS=linux & GOARCH=386 & go build -o koreon_linux_32_ctl main.go

# @REM MacOS 64bit binary compile
build_macos:
	GOOS=darwin & GOARCH=amd64 & go build -o koreon_macos_ctl main.go
