COMMAND_NAME := go-exif-extract
OUTPUT_DIR := ${GOPATH}/bin
OUTPUT_FILE := ${OUTPUT_DIR}/${COMMAND_NAME}
GOOS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH := $(subst x86_64,amd64,$(shell uname -m))
GO_FILES := $(shell find . -type f -not -path './vendor/*' -name '*.go')

.DEFAULT_GOAL: build
.PHONY: build clean vendor test

build: clean vendor
	@mkdir -p ${OUTPUT_DIR}
	env GO111MODULE=auto GOFLAGS=-mod=mod GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -o "$(COMMAND_NAME)"

clean:
	@rm -f ${OUTPUT_FILE}

vendor:
	go mod tidy && go mod vendor

test: vendor
	@go fmt ./...
	go test ./...
