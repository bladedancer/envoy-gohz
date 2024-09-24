GIT_VERSION ?= $(shell git describe --abbrev=8 --tags --always --dirty)
IMAGE_PREFIX ?= bladedancer
SERVICE_NAME=loghz

.PHONY: default
default: local.build ;

.PHONY: clean
clean:
	go clean
	rm -f bin/${SERVICE_NAME}

.PHONY: local.build
local.build: clean
	GOARCH=amd64 GOOS=linux GODEBUG=cgocheck=0 go build -o lib/libloghz.so -buildmode=c-shared pkg/main/filter.go pkg/main/config.go

dep:
	go mod tidy

vet:
	go vet

lint:
	golangci-lint run --enable-all