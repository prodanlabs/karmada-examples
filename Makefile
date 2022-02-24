GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
SOURCES := $(shell find ./cmd -type f  -name '*.go')

build:$(SOURCES)
		CGO_ENABLED=0 GOOS=$(GOOS) go build \
			-o karmada-custom-controller-manager \
			cmd/custom-controller-manager/custom-controller-manager.go

update-mod:
		go mod tidy && go mod vendor
.PHONY: clean
clean:
		rm -f karmada-custom-controller-manager