GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOARCH=$(shell go env GOARCH)
GOOS=$(shell go env GOOS )

BASE_PAH := $(shell pwd)
BUILD_PATH = $(BASE_PAH)/build
SERVER_PATH = $(BASE_PAH)
MAIN= $(BASE_PAH)/cmd/server/main.go
APP_NAME=go-server-base

tagName := $(shell date '+%y%m%d-%H%M')
tarName = ${APP_NAME}:${tagName}.tar
imageName = ${APP_NAME}:${tagName}

xBuilderName=${APP_NAME}_xbuilder

run_dev:
	MODE=dev go run ${MAIN}

run:
	go run ${MAIN}

upx_bin:
	upx $(BUILD_PATH)/$(APP_NAME)

build_backend_on_linux:
	cd $(SERVER_PATH) \
    && GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(APP_NAME) $(MAIN)

build_backend_on_darwin:
	cd $(SERVER_PATH) \
    && GOOS=linux GOARCH=amd64 $(GOBUILD) -trimpath -ldflags '-s -w'  -o $(BUILD_PATH)/$(APP_NAME) $(MAIN)

build_all: build_backend_on_linux

build_on_local: build_backend_on_darwin upx_bin

gobuild:
	$(GOBUILD) -o $(BUILD_PATH)/$(APP_NAME) $(MAIN)

dockerbuild:
		docker buildx create --name ${xBuilderName} --driver docker-container
		docker buildx use ${xBuilderName}
		docker buildx ls
	  docker buildx build --load --platform linux/amd64 -t ${imageName} -f ./docker/Dockerfile .
		docker buildx rm ${xBuilderName}
		docker save ${imageName} -o ${tarName}

dockerremove:
		docker buildx rm ${xBuilderName}

dockersave:
		docker save ${imageName} -o ${tarName}