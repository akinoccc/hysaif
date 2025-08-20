#!/bin/sh

# 启动脚本 - 用于Docker容器中启动所有服务
# 支持通过单一 Dockerfile 的运行时镜像构建方式

set -e

echo "=== 容器启动脚本开始执行 ==="

# 创建必要的日志目录
echo "创建日志目录..."
mkdir -p /var/log/supervisor
mkdir -p /var/log/nginx

# 检查Go API可执行文件
if [ ! -f "/app/api/hysaif-api" ]; then
    echo "错误: Go API可执行文件不存在: /app/api/hysaif-api"
    exit 1
fi

# 检查前端文件
echo "检查前端文件..."

# 检查Admin后台
if [ -d "/app/frontend" ]; then
    echo "检测到前端构建产物"
else
    echo "警告: 前端构建产物不存在，前端路由可能无法正常工作"
fi

[supervisord]
nodaemon=true
user=root
logfile=/var/log/supervisor/supervisord.log
pidfile=/var/run/supervisord.pid

[program:nginx]
command=nginx -g "daemon off;"
autostart=true
autorestart=true
stderr_logfile=/var/log/supervisor/nginx.err.log
stdout_logfile=/var/log/supervisor/nginx.out.log

[program:hysaif-api]
command=/app/api/hysaif-api -c /config/config.json
directory=/app/api
autostart=true
autorestart=true
environment=GIN_MODE="release"
stderr_logfile=/var/log/supervisor/api.err.log
stdout_logfile=/var/log/supervisor/api.out.log
EOF
fi

# 检查nginx配置文件
echo "验证nginx配置..."
if ! nginx -t; then
    echo "错误: nginx配置文件验证失败"
    exit 1
fi

# 设置Go API可执行权限
chmod +x /app/api/hysaif-api

echo "=== 启动supervisor管理所有服务 ==="
echo "服务列表:"
echo "  - nginx (反向代理)"
echo "  - hysaif-api (后端API服务)"
echo "  - hysaif-web (前端服务)"

# 启动supervisor
exec /usr/bin/supervisord -c /etc/supervisor/conf.d/supervisord.conf