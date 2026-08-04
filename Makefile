.PHONY: all build run gotool clean swagger help

BINARY := bluebell

all: gotool build

ifeq ($(OS),Windows_NT)
build:
	set "CGO_ENABLED=0" && go build -o "$(BINARY).exe" .

run: build
	@"$(BINARY).exe"

clean:
	@if exist "$(BINARY)" del /Q "$(BINARY)"
	@if exist "$(BINARY).exe" del /Q "$(BINARY).exe"
else
build:
	CGO_ENABLED=0 go build -o "$(BINARY)" .

run: build
	@./"$(BINARY)"

clean:
	@rm -f "$(BINARY)" "$(BINARY).exe"
endif

gotool:
	@go fmt ./...
	@go vet ./...

swagger:
	@swag init

help:
	@echo "make - 格式化、检查代码并构建当前平台的二进制文件"
	@echo "make build - 构建当前平台的二进制文件"
	@echo "make run - 构建并运行应用程序"
	@echo "make gotool - 格式化代码并执行静态检查"
	@echo "make swagger - 重新生成 Swagger 接口文档"
	@echo "make clean - 删除二进制文件"
	@echo "make help - 显示此帮助信息"
