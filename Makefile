.PHONY: help dev build clean

help:
	@echo "使用方式: make + [usage]"
	@echo "可用命令: [help, dev, build, clean]"

dev:
	@go fmt ./
	@go run ./main.go

build:
	@go fmt ./
	@go vet ./
	@go build main.go

clean:
	@rm -rf main
	
