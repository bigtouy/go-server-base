ARG VERSION=1.23

# 使用官方Go镜像作为构建环境
FROM golang:${VERSION}-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 复制go.mod和go.sum文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目源码
COPY . .

# 编译项目生成二进制文件
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o go-server-base ./cmd/server/main.go

FROM alpine

WORKDIR /app

# 时区设置成当前时区
RUN apk add --no-cache tzdata
ENV TZ="Asia/Shanghai"

ENV MODE="prod"

# 复制二进制文件到镜像
COPY --from=builder /app/go-server-base .

EXPOSE 8080

# 设置容器启动时执行的命令
CMD ["./go-server-base"]