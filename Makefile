NAME := mirror-proxy
CGO_ENABLED = 0
GO := go
BUILD_TARGET = build
COMMIT := $(shell git rev-parse --short HEAD)
VERSION := dev-$(shell git describe --tags $(shell git rev-list --tags --max-count=1))
BUILDFLAGS =
COVERED_MAIN_SRC_FILE=./main

darwin:
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/darwin/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/darwin/$(NAME)

linux:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/linux/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/linux/$(NAME)

build-all: darwin linux

run:
	./bin/darwin/$(NAME) --cert bin/rootCA/demo.crt --key bin/rootCA/demo.key

run-linux:
	./bin/darwin/$(NAME) --cert bin/rootCA/demo.crt --key bin/rootCA/demo.key

cert:
	mkdir -p bin/rootCA
	openssl genrsa -out bin/rootCA/demo.key 1024
	openssl req -new -x509 -days 1095 -key bin/rootCA/demo.key -out bin/rootCA/demo.crt -subj "/C=CN/ST=GD/L=SZ/O=vihoo/OU=dev/CN=demo.com/emailAddress=demo@demo.com"

fmt:
	go fmt ./pkg/...

verify: tools
	go vet ./pkg/...
	go get -u golang.org/x/lint/golint
	golint -set_exit_status ./pkg/...

before-test: fmt verify

tools:
	go get -u golang.org/x/lint/golint

test: tools before-test
	go vet ./...
	go test ./... -v -coverprofile coverage.out
