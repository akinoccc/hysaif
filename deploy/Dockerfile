# 运行时镜像：直接使用 CI 产物，不在镜像内执行构建
FROM node:22 AS runtime

# 安装 nginx 和 supervisor
RUN apt-get update && apt-get install -y nginx supervisor && rm -rf /var/lib/apt/lists/*

# 工作目录
WORKDIR /app

# 复制前端与后端的已构建产物（由 CI 预先放置在构建上下文中）
# 期望上下文包含：
# - frontend  (admin 前端静态资源)
# - hysaif-api (Go 编译后的二进制)
COPY frontend ./frontend
COPY api ./api

# 复制配置与启动脚本
COPY deploy/nginx.conf /etc/nginx/nginx.conf
COPY deploy/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY deploy/start.sh /app/start.sh
RUN chmod +x /app/start.sh

# 暴露端口
EXPOSE 3000 9000

# 启动命令
CMD ["/app/start.sh"]