#!/bin/bash

# JXZY 服务状态检查脚本

set -e  # 遇到错误立即退出

# 颜色定义
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

# 检查服务状态
check_service_status() {
    local service_name=$1
    local port=$2
    local pid=""
    
    if lsof -i :$port > /dev/null 2>&1; then
        pid=$(lsof -ti :$port | head -1)
        echo -e "  $service_name (端口: $port) - ${GREEN}运行中${NC} (PID: $pid)"
        return 0
    else
        echo -e "  $service_name (端口: $port) - ${RED}未运行${NC}"
        return 1
    fi
}

# 检查MySQL状态
check_mysql_status() {
    if lsof -i :13309 > /dev/null 2>&1; then
        echo -e "  MySQL (端口: 13309) - ${GREEN}运行中${NC}"
        return 0
    else
        echo -e "  MySQL (端口: 13309) - ${RED}未运行${NC}"
        return 1
    fi
}

# 主函数
main() {
    echo "=========================================="
    echo "    JXZY 服务状态检查"
    echo "=========================================="
    
    # 获取脚本所在目录
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
    
    log_info "项目根目录: $PROJECT_ROOT"
    
    # 切换到项目根目录
    cd "$PROJECT_ROOT"
    
    echo ""
    echo "服务状态:"
    echo "=========================================="
    
    # 检查各个服务状态
    local all_running=true
    
    # 检查 bs-llm 服务
    if ! check_service_status "bs-llm" 8081; then
        all_running=false
    fi
    
    # 检查 bs-rag 服务
    if ! check_service_status "bs-rag" 8082; then
        all_running=false
    fi
    
    # 检查 bll-context 服务
    if ! check_service_status "bll-context" 8080; then
        all_running=false
    fi
    
    # 检查 chat-api 服务
    if ! check_service_status "chat-api" 8888; then
        all_running=false
    fi
    
    # 检查 api-knowledge 服务
    if ! check_service_status "api-knowledge" 8889; then
        all_running=false
    fi
    
    # 检查 MySQL 服务
    if ! check_mysql_status; then
        all_running=false
    fi
    
    echo "=========================================="
    
    if [ "$all_running" = true ]; then
        log_success "所有服务运行正常！"
        echo ""
        echo "API测试:"
        echo "  # 聊天API测试:"
        echo "  curl -X POST http://localhost:8888/api/v1/chat/stream \\"
        echo "    -H \"Content-Type: application/json\" \\"
        echo "    -d '{\"user_id\": \"test\", \"message\": \"你好\", \"scene_code\": \"chat_general\"}'"
        echo ""
        echo "  # 知识库API测试:"
        echo "  curl -X POST http://localhost:8889/api/v1/knowledge/add \\"
        echo "    -H \"Content-Type: application/json\" \\"
        echo "    -d '{\"summary\": \"测试总结\", \"content\": \"测试内容\", \"user_id\": \"test\"}'"
        echo ""
        echo "  curl -X POST http://localhost:8889/api/v1/knowledge/delete \\"
        echo "    -H \"Content-Type: application/json\" \\"
        echo "    -d '{\"vector_id\": \"test_id\", \"user_id\": \"test\"}'"
    else
        log_warning "部分服务未运行"
        echo ""
        echo "启动所有服务:"
        echo "  ./scripts/restart_all.sh"
    fi
    
    echo ""
    echo "管理命令:"
    echo "  启动所有服务: ./scripts/restart_all.sh"
    echo "  停止所有服务: pkill -f \"chat-api|api-knowledge|bll-context|bs-llm|bs_rag\""
    echo "  查看日志:     tail -f logs/access.log"
    echo "  查看错误:     tail -f logs/error.log"
    echo "  查看慢查询:   tail -f logs/slow.log"
    echo "  查看统计:     tail -f logs/stat.log"
    echo ""
}

# 执行主函数
main "$@"
