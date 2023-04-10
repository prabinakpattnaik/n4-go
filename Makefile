PHONY: all build_only gen precommit

ifndef MAGMA_ROOT
MAGMA_ROOT = /home/$USER/magma
endif
export MAGMA_ROOT

ifndef BUILD_OUT
BUILD_OUT := $(shell go env GOBIN)
ifdef GOOS
BUILD_OUT := $(BUILD_OUT)/$(GOOS)
endif
ifdef GOARCH
BUILD_OUT := $(BUILD_OUT)/$(GOARCH)
endif
endif

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

all: test install

build: install

download:
	go mod download

install:
	GOARCH=$(GOARCH) GOOS=$(GOOS) go build -o  $(MAGMA_ROOT)/libn4_pfcp.so -buildmode=c-shared client/client.go
	sudo mv $(MAGMA_ROOT)/libn4_pfcp.so /usr/local/lib/
	sudo mv $(MAGMA_ROOT)/libn4_pfcp.h /usr/local/include/

test:
	go test ./...

clean:
	go clean ./...

gen:
	go generate ./...

lint:
	golint -min_confidence 1. ./...

build_only:
	go build ./...

precommit: build_only test

cover:
	go test ./... -coverprofile ./cover.tmp >/dev/null
	go tool cover -func=./cover.tmp | tail -n 1
	find . -name '*.go' | xargs wc -l | tail -n 1
	rm ./cover.tmp

