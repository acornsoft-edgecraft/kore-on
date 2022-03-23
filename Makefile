# Directory
ROOTDIR=${PWD}
TARGETDIR=~/gows/bin

REGI_SVR=regi.k3.acornsoft.io

GIT_COMMIT = `git rev-parse HEAD`
VERSION = v1.1.3
BUILD_DATE = `date +'%Y-%m-%dT%H:%M:%S'`
BUILD_OPTIONS = -ldflags "-X main.Version=$(VERSION) -X main.CommitId=$(GIT_COMMIT) -X main.BuildDate=$(BUILD_DATE)"
GOARCH=amd64

all: darwin-ctl docker

getfiles:
	@echo "Get necessary files ..."
	curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.19.0/crictl-v1.19.0-linux-amd64.tar.gz -o ./Dockerfile/scripts/files/crictl-v1.19.0-linux-amd64.tar.gz
	curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.20.0/crictl-v1.20.0-linux-amd64.tar.gz -o ./Dockerfile/scripts/files/crictl-v1.20.0-linux-amd64.tar.gz
	curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.21.0/crictl-v1.21.0-linux-amd64.tar.gz -o ./Dockerfile/scripts/files/crictl-v1.21.0-linux-amd64.tar.gz
	curl -L https://github.com/docker/compose/releases/download/1.29.2/docker-compose-Linux-x86_64 -o ./Dockerfile/scripts/files/docker-compose-v1.29.2-Linux-x86_64
	curl -L https://github.com/etcd-io/etcd/releases/download/v3.4.16/etcd-v3.4.16-linux-amd64.tar.gz -o ./Dockerfile/scripts/files/etcd-v3.4.16-linux-amd64.tar.gz
	curl -L https://github.com/goharbor/harbor/releases/download/v2.3.0/harbor-offline-installer-v2.3.0.tgz -o ./Dockerfile/scripts/files/harbor-offline-installer-v2.3.0.tgz

darwin:
	@echo "Make darwin binary ..."
	cd ${ROOTDIR}/koreon && GOOS=darwin GOARCH=${GOARCH} go build ${BUILD_OPTIONS} -o ${TARGETDIR}/darwin/koreon_darwin_${VERSION}
	ln -s ${TARGETDIR}/darwin/koreon_darwin_${VERSION} ${TARGETDIR}/darwin/koreon

linux:
	@echo "Make linux binary ..."
	cd ${ROOTDIR}/koreon && GOOS=linux GOARCH=${GOARCH} go build ${BUILD_OPTIONS} -o ${TARGETDIR}/linux/koreon_linux_${VERSION}
	ln -s ${TARGETDIR}/linux/koreon_linux_${VERSION} ~/gows/bin/linux/koreon

docker:
	@echo "Make docker image ..."
	cd ${ROOTDIR}/Dockerfile && docker build -t ${REGI_SVR}/k3lab/koreon:${VERSION} .

pushimage:
	echo "Push koreon image to ${REGI_SVR} ..."
	docker push ${REGI_SVR}/k3lab/koreon:${VERSION}

clean:
	rm -f ${TARGETDIR}/darwin/koreon_darwin_${VERSION}
	rm -f ${TARGETDIR}/darwin/koreon
	rm -f ${TARGETDIR}/linux/koreon_linux_${VERSION}
	rm -f ${TARGETDIR}/linux/koreon
	rm -f ${TARGETDIR}/windows/koreon_windows_${VERSION}.exe
	rm -f ${TARGETDIR}/windows/koreon.exe

linux-ctl:
	@echo "Make linux binary ..."

	GOOS=linux GOARCH=${GOARCH} go build ${BUILD_OPTIONS} -o ${TARGETDIR}/koreonctl_linux_${VERSION}
	ln -s ${TARGETDIR}/koreonctl_linux_${VERSION} ~/gows/bin/linux/koreonctl

darwin-ctl:
	@echo "Make darwin binary ..."
	rm -f ${TARGETDIR}/koreonctl_darwin_${VERSION}
	rm -f ~/gows/bin/darwin/koreonctl
	GOOS=darwin GOARCH=${GOARCH} go build ${BUILD_OPTIONS} -o ${TARGETDIR}/koreonctl_darwin_${VERSION}
	ln -s ${TARGETDIR}/koreonctl_darwin_${VERSION} ~/gows/bin/darwin/koreonctl

template:
	@echo "make template.go file ..."
	sh -c "./make_template_go.sh"