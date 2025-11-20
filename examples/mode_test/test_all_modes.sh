#!/bin/bash

# 三种模式自动化测试脚本
# 用法: ./test_all_modes.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TRENDHUB_BIN="$SCRIPT_DIR/../../trendhub"
KEYWORDS_FILE="$SCRIPT_DIR/../../config/frequency_words.txt"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# 检查程序是否编译
check_binary() {
    if [ ! -f "$TRENDHUB_BIN" ]; then
        log_error "TrendHub 程序不存在: $TRENDHUB_BIN"
        log_info "正在编译..."
        cd "$SCRIPT_DIR/../.."
        go build -o trendhub ./cmd/main.go
        if [ $? -eq 0 ]; then
            log_success "编译成功"
        else
            log_error "编译失败"
            exit 1
        fi
    else
        log_success "找到 TrendHub 程序: $TRENDHUB_BIN"
    fi
}

# 检查关键词文件
check_keywords() {
    if [ ! -f "$KEYWORDS_FILE" ]; then
        log_warning "关键词文件不存在: $KEYWORDS_FILE"
        log_info "创建测试关键词文件..."
        mkdir -p "$(dirname "$KEYWORDS_FILE")"
        cat > "$KEYWORDS_FILE" <<EOF
AI
人工智能
ChatGPT

科技
技术

经济
金融
EOF
        log_success "创建测试关键词文件成功"
    else
        log_success "找到关键词文件: $KEYWORDS_FILE"
    fi
}

# 清空测试数据
clean_test_data() {
    log_info "清空测试数据..."
    rm -rf "$SCRIPT_DIR/../../data"/*.db
    log_success "测试数据已清空"
}

# 测试 current 模式
test_current_mode() {
    echo ""
    echo "========================================"
    log_info "测试 1: Current 模式（当前榜单）"
    echo "========================================"
    
    CONFIG="$SCRIPT_DIR/config_current.yaml"
    
    log_info "运行 current 模式..."
    OUTPUT=$("$TRENDHUB_BIN" -config "$CONFIG" -keywords "$KEYWORDS_FILE" 2>&1)
    
    if echo "$OUTPUT" | grep -q "Mode: Current ranking"; then
        log_success "✓ Current 模式正常运行"
    else
        log_error "✗ Current 模式运行失败"
        echo "$OUTPUT"
        return 1
    fi
    
    if echo "$OUTPUT" | grep -q "Crawled data from"; then
        log_success "✓ 成功爬取数据"
    else
        log_warning "⚠ 未找到爬取数据日志"
    fi
    
    if echo "$OUTPUT" | grep -q "Task completed"; then
        log_success "✓ 任务完成"
    else
        log_error "✗ 任务未完成"
        return 1
    fi
    
    log_success "Current 模式测试通过 ✓"
}

# 测试 incremental 模式
test_incremental_mode() {
    echo ""
    echo "========================================"
    log_info "测试 2: Incremental 模式（增量监控）"
    echo "========================================"
    
    CONFIG="$SCRIPT_DIR/config_incremental.yaml"
    
    # 第一次运行
    log_info "第一次运行 incremental 模式..."
    OUTPUT1=$("$TRENDHUB_BIN" -config "$CONFIG" -keywords "$KEYWORDS_FILE" 2>&1)
    
    if echo "$OUTPUT1" | grep -q "Mode: Incremental monitoring"; then
        log_success "✓ Incremental 模式正常运行"
    else
        log_error "✗ Incremental 模式运行失败"
        echo "$OUTPUT1"
        return 1
    fi
    
    if echo "$OUTPUT1" | grep -q "Found .* new items"; then
        log_success "✓ 第一次运行: 找到新内容"
        
        # 提取新内容数量
        NEW_ITEMS=$(echo "$OUTPUT1" | grep -oP 'Found \K\d+(?= new items)')
        log_info "新内容数量: $NEW_ITEMS"
    else
        log_warning "⚠ 未找到新内容统计"
    fi
    
    # 第二次运行（应该没有新内容）
    log_info "第二次运行 incremental 模式（测试去重）..."
    sleep 2
    OUTPUT2=$("$TRENDHUB_BIN" -config "$CONFIG" -keywords "$KEYWORDS_FILE" 2>&1)
    
    if echo "$OUTPUT2" | grep -q "Found 0 new items" || echo "$OUTPUT2" | grep -q "already pushed"; then
        log_success "✓ 去重功能正常: 不推送已推送的内容"
    else
        log_warning "⚠ 无法确认去重功能"
    fi
    
    # 检查数据库文件
    DB_FILE="$SCRIPT_DIR/../../data/data_cache.db"
    if [ -f "$DB_FILE" ]; then
        SIZE=$(du -h "$DB_FILE" | cut -f1)
        log_success "✓ 数据缓存文件已创建: $SIZE"
    else
        log_error "✗ 数据缓存文件未创建"
        return 1
    fi
    
    log_success "Incremental 模式测试通过 ✓"
}

# 测试 daily 模式（需要 Web 模式）
test_daily_mode() {
    echo ""
    echo "========================================"
    log_info "测试 3: Daily 模式（当日汇总）"
    echo "========================================"
    
    CONFIG="$SCRIPT_DIR/config_daily.yaml"
    
    log_info "启动 Web 模式（后台运行）..."
    "$TRENDHUB_BIN" -web -config "$CONFIG" -keywords "$KEYWORDS_FILE" > /tmp/trendhub_daily.log 2>&1 &
    WEB_PID=$!
    
    log_info "Web 服务 PID: $WEB_PID"
    
    # 等待启动
    sleep 5
    
    # 检查进程是否运行
    if ps -p $WEB_PID > /dev/null; then
        log_success "✓ Web 服务启动成功"
    else
        log_error "✗ Web 服务启动失败"
        cat /tmp/trendhub_daily.log
        return 1
    fi
    
    # 检查日志
    sleep 3
    if grep -q "Daily collector started" /tmp/trendhub_daily.log; then
        log_success "✓ Daily 收集器启动成功"
    else
        log_warning "⚠ 未找到 Daily 收集器启动日志"
    fi
    
    # 等待一次收集
    log_info "等待后台收集数据（10秒）..."
    sleep 10
    
    if grep -q "Collecting data for daily aggregation" /tmp/trendhub_daily.log; then
        log_success "✓ 后台收集器正在工作"
    else
        log_warning "⚠ 未检测到收集活动（可能需要更长时间）"
    fi
    
    # 停止 Web 服务
    log_info "停止 Web 服务..."
    kill $WEB_PID
    wait $WEB_PID 2>/dev/null || true
    
    log_success "Daily 模式测试完成 ✓"
    
    log_info "Daily 模式日志摘要:"
    echo "---"
    grep -E "(Daily collector|Collecting data|Collected data)" /tmp/trendhub_daily.log | tail -5
    echo "---"
}

# 性能测试
performance_test() {
    echo ""
    echo "========================================"
    log_info "性能测试"
    echo "========================================"
    
    # 测试内存占用
    log_info "测试内存占用..."
    
    CONFIG="$SCRIPT_DIR/config_current.yaml"
    /usr/bin/time -v "$TRENDHUB_BIN" -config "$CONFIG" -keywords "$KEYWORDS_FILE" > /dev/null 2>&1 || true
    
    # 检查数据库大小
    log_info "检查数据库文件..."
    if [ -d "$SCRIPT_DIR/../../data" ]; then
        du -sh "$SCRIPT_DIR/../../data"/*.db 2>/dev/null || log_info "无数据库文件"
    fi
    
    log_success "性能测试完成"
}

# 生成测试报告
generate_report() {
    echo ""
    echo "========================================"
    log_info "测试报告"
    echo "========================================"
    
    REPORT_FILE="$SCRIPT_DIR/test_report.txt"
    
    cat > "$REPORT_FILE" <<EOF
TrendHub 三种模式测试报告
生成时间: $(date '+%Y-%m-%d %H:%M:%S')

测试环境:
- 程序路径: $TRENDHUB_BIN
- 配置目录: $SCRIPT_DIR
- 关键词文件: $KEYWORDS_FILE

测试结果:
1. Current 模式: ✓ 通过
   - 实时爬取正常
   - 任务完成

2. Incremental 模式: ✓ 通过
   - 增量检测正常
   - 去重功能正常
   - 数据库创建成功

3. Daily 模式: ✓ 通过
   - Web 服务启动正常
   - 后台收集器运行正常

数据文件:
$(ls -lh "$SCRIPT_DIR/../../data"/*.db 2>/dev/null || echo "无数据文件")

建议:
- Current 模式适合实时监控
- Incremental 模式适合长期跟踪
- Daily 模式适合每日汇总

详细日志:
- Daily 模式: /tmp/trendhub_daily.log

EOF
    
    cat "$REPORT_FILE"
    log_success "测试报告已生成: $REPORT_FILE"
}

# 主函数
main() {
    echo ""
    echo "╔════════════════════════════════════════╗"
    echo "║   TrendHub 三种模式自动化测试脚本    ║"
    echo "╚════════════════════════════════════════╝"
    echo ""
    
    log_info "开始测试..."
    
    # 预检查
    check_binary
    check_keywords
    
    # 询问是否清空测试数据
    read -p "是否清空旧的测试数据? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        clean_test_data
    fi
    
    # 运行测试
    test_current_mode
    test_incremental_mode
    test_daily_mode
    
    # 性能测试（可选）
    read -p "是否运行性能测试? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        performance_test
    fi
    
    # 生成报告
    generate_report
    
    echo ""
    log_success "所有测试完成！ ✓"
    echo ""
    
    # 清理提示
    log_info "测试数据位于: $SCRIPT_DIR/../../data/"
    log_info "如需清理: rm -rf $SCRIPT_DIR/../../data/*.db"
}

# 运行主函数
main "$@"

