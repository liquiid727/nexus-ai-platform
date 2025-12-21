.PHONY: all build run clean

all: build

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

clean:
	rm -rf bin/

#Protobuf 文件存放路径
APIROOT = $(shell pwd)/pkg/api/apiserver/v1

protoc:# 编译proto文件
# 生成gRPC服务代码
	protoc --go_out=paths=source_relative:$(APIROOT) $(APIROOT)/*.proto
#
	protoc --go-grpc_out=paths=source_relative:$(APIROOT) $(APIROOT)/*.proto
