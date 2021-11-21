BINARYNAME?=kaadm
VERSION?=latest
VersionPath="github.com/prodanlabs/kaadm/app/cmd"
.PHONY : clean
clean:
	rm -f ${BINARYNAME}

.PHONY : build
build:
	CGO_ENABLED=0 GOOS=linux go build -o ${BINARYNAME} -ldflags "-X '${VersionPath}.Version=${VERSION}'"  main.go

