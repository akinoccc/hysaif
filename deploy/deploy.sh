#!/bin/bash

# SIMS 自动化部署脚本
# 用于在Debian服务器上部署SIMS应用

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查是否为root用户
check_root() {
    if [[ $EUID -eq 0 ]]; then
        log_error "请不要使用root用户运行此脚本"
        exit 1
    fi
}

# 检查系统要求
check_system() {
    log_info "检查系统要求..."
    
    # 检查操作系统
    if ! grep -q "Debian" /etc/os-release; then
        log_error "此脚本仅支持Debian系统"
        exit 1
    fi
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    
    # 检查Docker Compose
    if ! command -v docker compose &> /dev/null; then
        log_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    
    log_info "系统检查通过"
}

# 安装依赖
install_dependencies() {
    log_info "安装系统依赖..."
    
    sudo apt-get update
    sudo apt-get install -y \
        curl \
        wget \
        git \
        unzip \
        htop \
        nginx \
        certbot \
        python3-certbot-nginx
    
    log_info "依赖安装完成"
}

# 创建项目目录
setup_directories() {
    log_info "设置项目目录..."
    
    PROJECT_DIR="/opt/sims"
    
    if [ ! -d "$PROJECT_DIR" ]; then
        sudo mkdir -p "$PROJECT_DIR"
        sudo chown $USER:$USER "$PROJECT_DIR"
    fi
    
    cd "$PROJECT_DIR"
    
    # 创建必要的子目录
    mkdir -p {logs,data,backups,ssl,nginx}
    
    log_info "项目目录设置完成: $PROJECT_DIR"
}

# 克隆或更新代码
setup_code() {
    log_info "设置应用代码..."
    
    if [ ! -d ".git" ]; then
        log_info "克隆代码仓库..."
        git clone https://github.com/akinoccc/hysaif.git .
    else
        log_info "更新代码..."
        git pull origin main
    fi
    
    log_info "代码设置完成"
}

# 设置环境变量
setup_environment() {
    log_info "设置环境变量..."
    
    if [ ! -f ".env" ]; then
        if [ -f "env.example" ]; then
            cp env.example .env
            log_warn "请编辑 .env 文件设置正确的环境变量"
            log_warn "特别注意设置安全的密码和密钥"
        else
            log_error "未找到 env.example 文件"
            exit 1
        fi
    fi
    
    # 生成随机密钥
    if grep -q "your-jwt-secret-key-here" .env; then
        JWT_SECRET=$(openssl rand -base64 64 | tr -d '\n')
        sed -i "s/your-jwt-secret-key-here/$JWT_SECRET/" .env
        log_info "已生成JWT密钥"
    fi
    
    if grep -q "your_secure_password" .env; then
        DB_PASSWORD=$(openssl rand -base64 32 | tr -d '\n')
        sed -i "s/your_secure_password/$DB_PASSWORD/" .env
        log_info "已生成数据库密码"
    fi
    
    log_info "环境变量设置完成"
}

# 设置SSL证书
setup_ssl() {
    log_info "设置SSL证书..."
    
    read -p "请输入域名 (例: example.com): " DOMAIN
    
    if [ -z "$DOMAIN" ]; then
        log_warn "跳过SSL证书设置"
        return
    fi
    
    # 使用Let's Encrypt获取证书
    sudo certbot --nginx -d "$DOMAIN" --non-interactive --agree-tos --email admin@"$DOMAIN"
    
    log_info "SSL证书设置完成"
}

# 启动服务
start_services() {
    log_info "启动服务..."
    
    # 加载环境变量
    source .env
    
    # 启动服务
    docker compose -f docker-compose.prod.yml up -d
    
    log_info "等待服务启动..."
    sleep 30
    
    # 健康检查
    if curl -f http://localhost/health > /dev/null 2>&1; then
        log_info "服务启动成功"
    else
        log_error "服务启动失败，请检查日志"
        docker compose -f docker-compose.prod.yml logs
        exit 1
    fi
}

# 设置定时任务
setup_cron() {
    log_info "设置定时任务..."
    
    # 备份数据库
    (crontab -l 2>/dev/null; echo "0 2 * * * cd /opt/sims && ./scripts/backup.sh") | crontab -
    
    # 清理Docker镜像
    (crontab -l 2>/dev/null; echo "0 3 * * 0 docker system prune -f") | crontab -
    
    # 更新SSL证书
    (crontab -l 2>/dev/null; echo "0 4 1 * * certbot renew --quiet") | crontab -
    
    log_info "定时任务设置完成"
}

# 显示部署信息
show_info() {
    log_info "部署完成！"
    echo
    echo "应用信息:"
    echo "- 前端地址: http://localhost"
    echo "- API地址: http://localhost/api"
    echo "- 项目目录: /opt/sims"
    echo
    echo "管理命令:"
    echo "- 查看日志: docker compose -f docker-compose.prod.yml logs"
    echo "- 重启服务: docker compose -f docker-compose.prod.yml restart"
    echo "- 停止服务: docker compose -f docker-compose.prod.yml down"
    echo "- 更新应用: cd /opt/sims && git pull && docker compose -f docker-compose.prod.yml up -d"
    echo
    echo "备份目录: /opt/sims/backups"
    echo "日志目录: /opt/sims/logs"
}

# 主函数
main() {
    log_info "开始SIMS自动化部署..."
    
    check_root
    check_system
    install_dependencies
    setup_directories
    setup_code
    setup_environment
    setup_ssl
    start_services
    setup_cron
    show_info
    
    log_info "部署完成！"
}

# 运行主函数
main "$@" 