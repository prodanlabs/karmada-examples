BINARYNAME?=kaadm
VERSION?=latest
KarmadaVersion?=latest
GitCommitID?=`git rev-parse HEAD`
VersionPath?="github.com/prodanlabs/kaadm/app/cmd"
.PHONY : clean
clean:
	rm -f ${BINARYNAME}

.PHONY : build
build:
	CGO_ENABLED=0 GOOS=linux go build -o ${BINARYNAME} -ldflags "-X '${VersionPath}.Version=${VERSION}' -X '${VersionPath}.GitCommitID=${GitCommitID}' -X '${VersionPath}.KarmadaVersion=${KarmadaVersion}'"  main.go

.PHONY : arm64
arm64:
	GOARCH=arm64  GOOS="linux"  GOARM=7  CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ AR=aarch64-linux-gnu-ar \
	go build -o ${BINARYNAME} -ldflags "-X '${VersionPath}.Version=${VERSION}' -X '${VersionPath}.GitCommitID=${GitCommitID}' -X '${VersionPath}.KarmadaVersion=${KarmadaVersion}'"  main.go