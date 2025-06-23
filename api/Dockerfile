# 构建阶段
FROM golang:1.24-bullseye AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apt-get update && apt-get install -y git gcc libc6-dev ca-certificates tzdata && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 构建应用 - 优化构建参数
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o sims-api .

# 运行阶段
FROM debian:bullseye-slim

# 安装必要的包
RUN apt-get update && apt-get install -y \
    ca-certificates \
    tzdata \
    wget \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# 创建非root用户
RUN groupadd -g 1001 sims && \
    useradd -r -u 1001 -g sims sims

# 设置工作目录
WORKDIR /app

# 从构建阶段复制文件
COPY --from=builder /app/sims-api .
COPY --from=builder /app/rbac_model.conf .

# 注意：config.json 将通过部署脚本或挂载卷提供，不在镜像中包含

# 创建数据目录并设置权限
RUN mkdir -p /app/data && chown -R sims:sims /app

# 切换到非root用户
USER sims

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./sims-api"]