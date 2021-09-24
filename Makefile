# Directory
ROOTDIR=${PWD}
TARGETDIR=~/gows/bin

REGI_SVR=regi.k3.acornsoft.io

GIT_COMMIT = `git rev-parse HEAD`
VERSION = 1.1.0
BUILD_DATE = `date +'%Y-%m-%dT%H:%M:%S'`
BUILD_OPTIONS = -ldflags "-X main.Version=$(VERSION) -X main.CommitId=$(GIT_COMMIT) -X main.BuildDate=$(BUILD_DATE)"
GOARCH=amd64

all: docker

darwin:
	@echo "Make darwin binary ..."
	cd ${ROOTDIR}/knit && GOOS=darwin GOARCH=${GOARCH} go build ${BUILD_OPTIONS} -o ${TARGETDIR}/darwin/knit_darwin_${VERSION}
	ln -s ${TARGETDIR}/darwin/knit_darwin_${VERSION} ${TARGETDIR}/darwin/knit

linux:
	@echo "Make linux binary ..."
	cd ${ROOTDIR}/knit && GOOS=linux GOARCH=${GOARCH} go build ${BUILD_OPTIONS} -o ${TARGETDIR}/linux/knit_linux_${VERSION}
	ln -s ${TARGETDIR}/linux/knit_linux_${VERSION} ~/gows/bin/linux/knit

docker:
	@echo "Make docker image ..."
	cd ${ROOTDIR}/Dockerfile && docker build -t ${REGI_SVR}/k3lab/knit:${VERSION} .

pushimage:
	echo "Push knit image to ${REGI_SVR} ..."
	docker push ${REGI_SVR}/k3lab/knit:${VERSION}

clean:
	rm -f ${TARGETDIR}/darwin/knit_darwin_${VERSION}
	rm -f ${TARGETDIR}/darwin/knit
	rm -f ${TARGETDIR}/linux/knit_linux_${VERSION}
	rm -f ${TARGETDIR}/linux/knit
	rm -f ${TARGETDIR}/windows/knit_windows_${VERSION}.exe
	rm -f ${TARGETDIR}/windows/knit.exe