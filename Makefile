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
	GOPROXY=https://mirrors.aliyun.com/goproxy/ CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/linux/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/linux/$(NAME)

win:
	go get github.com/inconshreveable/mousetrap
	go get github.com/mattn/go-isatty
	GOPROXY=https://mirrors.aliyun.com/goproxy/ CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=386 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/windows/$(NAME).exe $(MAIN_SRC_FILE)

build-all: darwin linux win

clean:
	rm -rfv bin

run: darwin cert
	./bin/darwin/$(NAME) --cert bin/rootCA/demo.crt --key bin/rootCA/demo.key --config config/.mirror-proxy.yaml

run-no-dev: darwin
	./bin/darwin/$(NAME) --config config/.mirror-proxy.yaml

run-linux: linux cert
	./bin/darwin/$(NAME) --cert bin/rootCA/demo.crt --key bin/rootCA/demo.key

run-win: win cert
	./bin/windows/$(NAME) --cert bin/rootCA/demo.crt --key bin/rootCA/demo.key

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

before-test: fmt

tools:
	go get -u golang.org/x/lint/golint

test: tools before-test
	go vet ./...
	go test ./... -v -coverprofile coverage.out

image: linux
	docker build -t jenkinszh/mirror-proxy:dev-$(COMMIT) .

image-github: linux
	docker build -t docker.pkg.github.com/jenkins-zh/mirror-proxy/mirror-proxy:dev-$(COMMIT) .

push-image: image
	docker push jenkinszh/mirror-proxy:dev-$(COMMIT)

push-github-image: image-github
	docker push docker.pkg.github.com/jenkins-zh/mirror-proxy/mirror-proxy:dev-$(COMMIT)

run-image: image
	docker run -p 7070:7070 --rm jenkinszh/mirror-proxy:dev

front-image:
	docker build -t jenkinszh/mirror-proxy-front:dev-$(COMMIT) front

front-github-image:
	docker build -t docker.pkg.github.com/jenkins-zh/mirror-proxy/mirror-proxy-front:dev-$(COMMIT) front

front-image-push: front-image
	docker push jenkinszh/mirror-proxy-front:dev-$(COMMIT)

front-github-image-push: front-github-image
	docker push docker.pkg.github.com/jenkins-zh/mirror-proxy/mirror-proxy-front:dev-$(COMMIT)
