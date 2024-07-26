FRONTEND_DIR = ./web
BACKEND_DIR = .

.PHONY: all build-frontend start-backend build push clean

all: build-frontend start-backend

build-frontend:
	@echo "Building frontend..."
	@cd $(FRONTEND_DIR) && yarn install --network-timeout 1000000 && DISABLE_ESLINT_PLUGIN='true' VITE_REACT_APP_VERSION=$(cat VERSION) yarn build

start-backend:
	@echo "Starting backend dev server..."
	@cd $(BACKEND_DIR) && go run main.go &


REGISTRY := your-registry.com
IMAGE_NAME := new-api
VERSION := $(shell git describe --tags --always --dirty)

FULL_IMAGE_NAME := $(REGISTRY)/$(IMAGE_NAME)
VERSION_TAG := $(FULL_IMAGE_NAME):$(VERSION)
LATEST_TAG := $(FULL_IMAGE_NAME):latest

# 构建Docker镜像
build:
	@echo "Building Docker image..."
	docker build -t $(VERSION_TAG) .
	docker tag $(VERSION_TAG) $(LATEST_TAG)

# 推送Docker镜像到仓库
push:
	@echo "Pushing Docker image to registry..."
	docker push $(VERSION_TAG)
	docker push $(LATEST_TAG)

# 清理
clean:
	@echo "Cleaning up..."
	docker rmi $(VERSION_TAG)
	docker rmi $(LATEST_TAG)