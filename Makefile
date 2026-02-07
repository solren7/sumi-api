# 变量定义，方便后续维护
BINARY_NAME=bin/server
MAIN_PATH=main.go

.PHONY: run build generate api help

# 1. 生成 SQL 代码
# 这里直接在根目录调用 sqlc，避免嵌套 make 带来的路径混乱
generate:
	@echo "Generating SQLC code..."
	sqlc generate

# 2. 运行主程序
run:
	go run $(MAIN_PATH)

# 3. 运行 API (通常与 run 类似，或者可以指定配置文件)
api:
	@echo "Starting API server..."
	go run $(MAIN_PATH) api

# 4. 编译二进制文件
build:
	@echo "Building binary..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# 帮助信息
help:
	@echo "Usage:"
	@echo "  make generate  - Run sqlc generate in /sqlc directory"
	@echo "  make run       - Run the application"
	@echo "  make api       - Alias for running the API server"
	@echo "  make build     - Build the server binary"