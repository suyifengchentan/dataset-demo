# 使用 golang 官方镜像作为构建环境
FROM golang:1.24-alpine AS builder

# 安装必要的构建工具
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目源代码
COPY . .

# 构建应用程序，禁用 CGO 以确保二进制文件的可移植性
RUN CGO_ENABLED=0 GOOS=linux go build -o database-demo .

# 使用轻量级的 alpine 镜像作为运行环境
FROM alpine:latest

# 安装必要的运行时库（如 ca-certificates，如果需要访问 HTTPS）
RUN apk add --no-cache ca-certificates

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/database-demo .

# 暴露应用程序端口
EXPOSE 8080

# 启动应用程序
CMD ["./database-demo"]
