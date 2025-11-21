.PHONY: build run clean web dev install help test lint build-all build-linux build-windows build-darwin build-linux-amd64 build-linux-arm64 build-windows-amd64 build-windows-arm64 build-darwin-amd64 build-darwin-arm64

# 变量定义
BINARY_NAME=trendhub
BUILD_DIR=build
VERSION_FILE=version
VERSION=$(shell cat $(VERSION_FILE) 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION=$(shell go version | awk '{print $$3}')
LDFLAGS=-X 'github.com/gotoailab/trendhub/internal/version.Version=$(VERSION)' \
        -X 'github.com/gotoailab/trendhub/internal/version.BuildTime=$(BUILD_TIME)' \
        -X 'github.com/gotoailab/trendhub/internal/version.GitCommit=$(GIT_COMMIT)' \
        -X 'github.com/gotoailab/trendhub/internal/version.GoVersion=$(GO_VERSION)' \
        -s -w

# 默认配置路径
CONFIG_PATH?=config/config.yaml
KEYWORDS_PATH?=config/frequency_words.txt
WEB_ADDR?=:8080
PUSH_DB_PATH?=data/push_records.db
CACHE_DB_PATH?=data/data_cache.db

# 帮助信息
help:
	@echo "TrendHub Makefile 命令："
	@echo ""
	@echo "构建相关："
	@echo "  make build          - 构建可执行文件（当前平台，带版本信息）"
	@echo "  make build-fast     - 快速构建（不带版本信息）"
	@echo "  make build-all      - 构建所有平台和架构"
	@echo "  make install        - 构建并安装到系统 PATH"
	@echo ""
	@echo "跨平台构建："
	@echo "  make build-linux    - 构建 Linux 版本（amd64 + arm64）"
	@echo "  make build-windows  - 构建 Windows 版本（amd64 + arm64）"
	@echo "  make build-darwin   - 构建 macOS 版本（amd64 + arm64）"
	@echo "  make build-linux-amd64    - 构建 Linux amd64"
	@echo "  make build-linux-arm64    - 构建 Linux arm64"
	@echo "  make build-windows-amd64  - 构建 Windows amd64"
	@echo "  make build-windows-arm64 - 构建 Windows arm64"
	@echo "  make build-darwin-amd64   - 构建 macOS amd64 (Intel)"
	@echo "  make build-darwin-arm64   - 构建 macOS arm64 (Apple Silicon)"
	@echo ""
	@echo "运行相关："
	@echo "  make run            - 构建并运行（命令行模式）"
	@echo "  make web            - 构建并运行 Web 模式"
	@echo "  make dev            - 开发模式运行（不构建，直接 go run）"
	@echo ""
	@echo "运行参数示例："
	@echo "  make run CONFIG_PATH=config/config.yaml KEYWORDS_PATH=config/frequency_words.txt"
	@echo "  make web WEB_ADDR=:8080 CONFIG_PATH=config/config.yaml"
	@echo ""
	@echo "其他："
	@echo "  make clean          - 清理所有构建文件"
	@echo "  make clean-cross    - 清理跨平台构建文件（保留当前平台）"
	@echo "  make test           - 运行测试"
	@echo "  make lint           - 代码检查"
	@echo "  make version        - 显示版本信息"
	@echo "  make help           - 显示此帮助信息"

# 构建（带版本信息）
build:
	@echo "构建 $(BINARY_NAME) v$(VERSION)..."
	@echo "构建时间: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go 版本: $(GO_VERSION)"
	@mkdir -p $(BUILD_DIR)
	@go mod tidy
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) cmd/main.go
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"
	@$(BUILD_DIR)/$(BINARY_NAME) -version 2>/dev/null || echo "提示: 程序可能不支持 -version 参数"

# 快速构建（不带版本信息）
build-fast:
	@echo "快速构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go mod tidy
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/main.go
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 安装到系统
install: build
	@echo "安装 $(BINARY_NAME) 到系统..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "安装完成: /usr/local/bin/$(BINARY_NAME)"

# 运行（命令行模式）
run: build
	@echo "运行 $(BINARY_NAME) (命令行模式)..."
	@$(BUILD_DIR)/$(BINARY_NAME) \
		-config $(CONFIG_PATH) \
		-keywords $(KEYWORDS_PATH) \
		-pushdb $(PUSH_DB_PATH) \
		-cachedb $(CACHE_DB_PATH)

# 运行 Web 模式
web: build
	@echo "运行 $(BINARY_NAME) (Web 模式)..."
	@echo "访问地址: http://localhost$(WEB_ADDR)"
	@$(BUILD_DIR)/$(BINARY_NAME) \
		-web \
		-addr $(WEB_ADDR) \
		-config $(CONFIG_PATH) \
		-keywords $(KEYWORDS_PATH) \
		-pushdb $(PUSH_DB_PATH) \
		-cachedb $(CACHE_DB_PATH)

# 开发模式（直接运行，不构建）
dev:
	@echo "开发模式运行..."
	@go run cmd/main.go \
		-web \
		-addr $(WEB_ADDR) \
		-config $(CONFIG_PATH) \
		-keywords $(KEYWORDS_PATH) \
		-pushdb $(PUSH_DB_PATH) \
		-cachedb $(CACHE_DB_PATH)

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@echo "清理完成"

# 清理跨平台构建（保留当前平台构建）
clean-cross:
	@echo "清理跨平台构建文件..."
	@find $(BUILD_DIR) -type d -name "*-*" -exec rm -rf {} + 2>/dev/null || true
	@echo "跨平台构建文件清理完成"

# 运行测试
test:
	@echo "运行测试..."
	@go test -v ./...

# 代码检查
lint:
	@echo "代码检查..."
	@go vet ./...
	@echo "提示: 可以使用 golangci-lint 进行更详细的检查"

# 创建必要的目录
init-dirs:
	@echo "创建必要的目录..."
	@mkdir -p build data config
	@echo "目录创建完成"

# 显示版本信息
version: build
	@echo "显示版本信息..."
	@$(BUILD_DIR)/$(BINARY_NAME) -version

# 跨平台构建函数
define build-platform
	@echo "构建 $(1)/$(2) 版本..."
	@mkdir -p $(BUILD_DIR)/$(1)-$(2)
	@GOOS=$(1) GOARCH=$(2) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)_$(1)_$(2)_$(3) cmd/main.go
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)_$(1)_$(2)_$(3)"
endef

# Linux amd64
build-linux-amd64:
	$(call build-platform,linux,amd64,)

# Linux arm64
build-linux-arm64:
	$(call build-platform,linux,arm64,)

# Linux (amd64 + arm64)
build-linux: build-linux-amd64 build-linux-arm64
	@echo "Linux 版本构建完成"

# Windows amd64
build-windows-amd64:
	$(call build-platform,windows,amd64,.exe)

# Windows arm64
build-windows-arm64:
	$(call build-platform,windows,arm64,.exe)

# Windows (amd64 + arm64)
build-windows: build-windows-amd64 build-windows-arm64
	@echo "Windows 版本构建完成"

# macOS amd64 (Intel)
build-darwin-amd64:
	$(call build-platform,darwin,amd64,)

# macOS arm64 (Apple Silicon)
build-darwin-arm64:
	$(call build-platform,darwin,arm64,)

# macOS (amd64 + arm64)
build-darwin: build-darwin-amd64 build-darwin-arm64
	@echo "macOS 版本构建完成"

# 构建所有平台
build-all: build-linux build-windows build-darwin
	@echo ""
	@echo "=========================================="
	@echo "所有平台构建完成！"
	@echo "构建文件位于: $(BUILD_DIR)/"
	@echo "=========================================="
	@find $(BUILD_DIR) -name "$(BINARY_NAME)*" -type f | sort
