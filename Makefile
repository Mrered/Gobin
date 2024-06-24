# 项目名称
PROJECT_NAME := Gobin

# 二进制文件输出目录
BIN_DIR := bin

# Go 源代码目录
SRC_DIR := ./cmd

# 自动发现所有二进制文件名
BINARIES := $(shell find $(SRC_DIR) -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)

# 默认目标
.PHONY: all
all: build

# 构建所有二进制文件
.PHONY: build
build: $(BINARIES)

$(BINARIES):
	@echo "Building $@..."
	@go build -o $(BIN_DIR)/$@ $(SRC_DIR)/$@/main.go

# 清理生成的文件
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)

# 更新依赖
.PHONY: tidy
tidy:
	@echo "Tidying up dependencies..."
	@go mod tidy

# 显示帮助信息
.PHONY: help
help:
	@echo "Makefile for $(PROJECT_NAME)"
	@echo
	@echo "Usage:"
	@echo "  make [target]"
	@echo
	@echo "Targets:"
	@echo "  all       - 构建所有二进制文件 (默认目标)"
	@echo "  build     - 构建所有二进制文件"
	@echo "  clean     - 清理生成的文件"
	@echo "  tidy      - 更新依赖"
	@echo "  help      - 显示帮助信息"