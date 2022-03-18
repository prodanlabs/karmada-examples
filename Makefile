GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
SOURCES := $(shell find ./ -type f  -name '*.go')

karmada-custom-controller-manager:$(SOURCES)
		CGO_ENABLED=0 GOOS=$(GOOS) go build \
			-o karmada-custom-controller-manager \
			cmd/custom-controller-manager/custom-controller-manager.go

karmada-custom-webhook:$(SOURCES)
		CGO_ENABLED=0 GOOS=$(GOOS) go build \
			-o karmada-custom-webhook \
			cmd/custom-webhook/custom-webhook.go
update:
		go mod tidy && go mod vendor

tls:
		bash manifests/webhook/create-tls.sh

.PHONY: clean
clean:
		rm -rf karmada-custom-controller-manager karmada-custom-webhook manifests/webhook/test-certs