SVC=mem_prometheus_exporter

VERSION := $(shell git rev-parse HEAD)
BUILDDATE := $(shell date -u +"%d-%m-%Y_%H:%M:%S_UTC")

LDFLAGS=-ldflags '-s -w -extldflags "-static"'

.PHONY: default
default: bin

.PHONY: test
test:
	GO111MODULE=on go test ./... -count=1

.PHONY: bin
bin: test
	GO111MODULE=on CGO_ENABLED=0 GOARCH=amd64 go build ${LDFLAGS} -o bin/${SVC}-amd64-linux

.PHONY: install
install: test
	go install ./...

.PHONY: clean
	rm -rf bin
