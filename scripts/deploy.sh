#!/bin/bash

# JXZY AI应用服务部署脚本
# 支持开发和生产环境的一键部署

set -e

# 配置参数
PROJECT_NAME="jxzy"
DOCKER_COMPOSE_FILE="docker-compose.yml"
ENV_FILE=".env"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "检查部署依赖..."
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    
    # 检查Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    
    # 检查Docker服务状态
    if ! docker info &> /dev/null; then
        log_error "Docker服务未启动，请先启动Docker"
        exit 1
    fi
    
    log_success "依赖检查完成"
}

# 创建环境配置文件
create_env_file() {
    log_info "创建环境配置文件..."
    
    if [ ! -f "$ENV_FILE" ]; then
        cat > "$ENV_FILE" << EOF
# JXZY AI应用服务环境配置

# 数据库配置
MYSQL_ROOT_PASSWORD=jxzy123456
MYSQL_DATABASE=jxzy
MYSQL_USER=jxzy
MYSQL_PASSWORD=jxzy123456

# Redis配置
REDIS_PASSWORD=

# 豆包API配置 (请填入实际的API密钥)
DOUBAO_API_KEY=your-doubao-api-key
DOUBAO_BASE_URL=https://ark.cn-beijing.volces.com/api/v3
DOUBAO_MODEL=ep-20241201-xxxx

# 阿里云配置 (请填入实际的密钥)
ALIBABA_CLOUD_ACCESS_KEY_ID=your-access-key-id
ALIBABA_CLOUD_ACCESS_KEY_SECRET=your-access-key-secret
ALIBABA_CLOUD_REGION=cn-beijing

# 向量数据库配置
DASHVECTOR_ENDPOINT=https://dashvector.cn-beijing.aliyuncs.com
DASHVECTOR_API_KEY=your-dashvector-api-key

# JWT配置
JWT_SECRET=jxzy-gateway-secret-key-change-in-production

# 监控配置
GRAFANA_ADMIN_PASSWORD=admin123
EOF
        log_warning "已创建环境配置文件 $ENV_FILE，请填入实际的API密钥"
        log_warning "编辑文件: vi $ENV_FILE"
        return 1
    fi
    
    log_success "环境配置文件已存在"
    return 0
}

# 创建必要的目录
create_directories() {
    log_info "创建必要的目录..."
    
    mkdir -p logs/{apis,blls,bs}
    mkdir -p bs/rag-rpc/storage/documents
    mkdir -p monitor/grafana/provisioning
    
    log_success "目录创建完成"
}

# 构建镜像
build_images() {
    log_info "构建Docker镜像..."
    
    # 并行构建所有服务镜像
    docker-compose build --parallel
    
    log_success "Docker镜像构建完成"
}

# 启动基础服务
start_infrastructure() {
    log_info "启动基础服务..."
    
    # 启动基础服务：Redis, Etcd, MySQL
    docker-compose up -d redis etcd mysql
    
    # 等待MySQL启动完成
    log_info "等待MySQL启动完成..."
    sleep 30
    
    # 检查MySQL连接
    until docker-compose exec mysql mysqladmin ping -h"localhost" --silent; do
        log_info "等待MySQL连接..."
        sleep 5
    done
    
    log_success "基础服务启动完成"
}

# 启动RPC服务
start_rpc_services() {
    log_info "启动RPC服务..."
    
    # 按依赖顺序启动RPC服务
    docker-compose up -d llm-rpc
    sleep 10
    
    docker-compose up -d rag-rpc
    sleep 10
    
    docker-compose up -d context-rpc prompt-rpc
    sleep 10
    
    log_success "RPC服务启动完成"
}

# 启动Gateway服务
start_gateway() {
    log_info "启动Gateway服务..."
    
    docker-compose up -d gateway-api
    
    # 等待Gateway启动
    sleep 15
    
    # 健康检查
    if curl -f http://localhost:8888/health &> /dev/null; then
        log_success "Gateway服务启动成功"
    else
        log_warning "Gateway服务可能未完全启动，请检查日志"
    fi
}

# 启动监控服务
start_monitoring() {
    log_info "启动监控服务..."
    
    docker-compose up -d prometheus grafana
    
    log_success "监控服务启动完成"
    log_info "Prometheus: http://localhost:9090"
    log_info "Grafana: http://localhost:3000 (admin/admin123)"
}

# 显示服务状态
show_status() {
    log_info "服务状态："
    docker-compose ps
    
    echo ""
    log_info "服务地址："
    echo "Gateway API: http://localhost:8888"
    echo "Prometheus: http://localhost:9090"
    echo "Grafana: http://localhost:3000"
    echo ""
    log_info "健康检查："
    echo "curl http://localhost:8888/health"
}

# 停止所有服务
stop_services() {
    log_info "停止所有服务..."
    docker-compose down
    log_success "所有服务已停止"
}

# 清理资源
cleanup() {
    log_info "清理Docker资源..."
    
    # 停止并删除容器
    docker-compose down -v
    
    # 删除未使用的镜像
    docker image prune -f
    
    # 删除未使用的卷
    docker volume prune -f
    
    log_success "资源清理完成"
}

# 查看日志
view_logs() {
    local service=$1
    if [ -z "$service" ]; then
        log_info "查看所有服务日志..."
        docker-compose logs -f
    else
        log_info "查看 $service 服务日志..."
        docker-compose logs -f "$service"
    fi
}

# 重启服务
restart_service() {
    local service=$1
    if [ -z "$service" ]; then
        log_error "请指定要重启的服务名称"
        exit 1
    fi
    
    log_info "重启 $service 服务..."
    docker-compose restart "$service"
    log_success "$service 服务重启完成"
}

# 显示帮助信息
show_help() {
    echo "JXZY AI应用服务部署脚本"
    echo ""
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  deploy      - 完整部署（默认）"
    echo "  start       - 启动所有服务"
    echo "  stop        - 停止所有服务"
    echo "  restart     - 重启所有服务"
    echo "  status      - 显示服务状态"
    echo "  logs [服务] - 查看日志"
    echo "  build       - 重新构建镜像"
    echo "  cleanup     - 清理资源"
    echo "  help        - 显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 deploy"
    echo "  $0 logs gateway-api"
    echo "  $0 restart llm-rpc"
}

# 主函数
main() {
    local command=${1:-deploy}
    
    case $command in
        "deploy")
            check_dependencies
            if ! create_env_file; then
                log_error "请先配置环境变量文件 $ENV_FILE"
                exit 1
            fi
            create_directories
            build_images
            start_infrastructure
            start_rpc_services
            start_gateway
            start_monitoring
            show_status
            ;;
        "start")
            check_dependencies
            docker-compose up -d
            show_status
            ;;
        "stop")
            stop_services
            ;;
        "restart")
            if [ -n "$2" ]; then
                restart_service "$2"
            else
                docker-compose restart
                show_status
            fi
            ;;
        "status")
            show_status
            ;;
        "logs")
            view_logs "$2"
            ;;
        "build")
            build_images
            ;;
        "cleanup")
            cleanup
            ;;
        "help"|"--help"|"-h")
            show_help
            ;;
        *)
            log_error "未知命令: $command"
            show_help
            exit 1
            ;;
    esac
}

# 脚本入口
main "$@"
