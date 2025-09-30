#!/bin/bash

# 代码上下文生成器 - 使用示例脚本
# 该脚本展示了如何使用不同的配置文件和参数

echo "🚀 代码上下文生成器 - 使用示例"
echo "=================================="

# 检查是否已安装工具
if ! command -v ./code-context-generator &> /dev/null; then
    echo "❌ 请先编译代码上下文生成器:"
    echo "   go build -o code-context-generator cmd/cli/main.go"
    exit 1
fi

# 创建输出目录
mkdir -p output/logs

echo ""
echo "1️⃣ 基础扫描 - 使用默认配置"
echo "------------------------"
./code-context-generator generate \
    -o output/basic-scan.json

echo ""
echo "2️⃣ 使用基础配置文件"
echo "------------------"
./code-context-generator generate \
    -c examples/basic-config.toml \
    -o output/basic-config-output.json

echo ""
echo "3️⃣ 生成项目文档（包含文件内容）"
echo "--------------------------------"
./code-context-generator generate \
    -c examples/project-documentation.toml \
    -o output/project-documentation.md

echo ""
echo "4️⃣ 高性能扫描（适合大项目）"
echo "----------------------------"
./code-context-generator generate \
    -c examples/performance-optimized.toml \
    -o output/fast-scan.json

echo ""
echo "5️⃣ 自定义参数扫描"
echo "----------------"
./code-context-generator generate \
    -f xml \
    -e "node_modules" -e ".git" -e "*.log" \
    -s 1048576 \
    -d 3 \
    -o output/custom-scan.xml

echo ""
echo "6️⃣ 交互式文件选择"
echo "----------------"
echo "📝 这将启动交互式选择器..."
./code-context-generator select \
    -m \
    -f markdown \
    -o output/selected-files.md

echo ""
echo "7️⃣ 生成配置文件"
echo "--------------"
./code-context-generator config init \
    -o output/my-config.toml

echo ""
echo "8️⃣ 验证配置文件"
echo "--------------"
./code-context-generator config validate \
    -c examples/basic-config.toml

echo ""
echo "9️⃣ 显示当前配置"
echo "--------------"
./code-context-generator config show

echo ""
echo "🔟 性能测试"
echo "----------"
echo "📊 扫描性能测试..."
time ./code-context-generator generate \
    -c examples/performance-optimized.toml \
    -o output/performance-test.json

echo ""
echo "📋 批处理示例 - 扫描多个目录"
echo "============================="

# 定义要扫描的目录数组
PROJECT_DIRS=(
    "src"
    "internal"
    "pkg"
    "cmd"
)

for dir in "${PROJECT_DIRS[@]}"; do
    if [ -d "$dir" ]; then
        echo "📁 扫描目录: $dir"
        ./code-context-generator generate \
            "$dir" \
            -f json \
            -e "*.log" -e "*.tmp" \
            -o "output/scan_${dir//\//_}.json"
    fi
done

echo ""
echo "🔄 定时任务示例"
echo "============="

# 创建定时任务脚本
cat > output/scheduled-scan.sh << 'EOF'
#!/bin/bash
# 定时扫描脚本
DATE=$(date +%Y%m%d_%H%M%S)
PROJECT_NAME=$(basename "$PWD")

./code-context-generator generate \
    -c examples/project-documentation.toml \
    -o "backup/${PROJECT_NAME}_documentation_${DATE}.md"

echo "✅ 备份完成: backup/${PROJECT_NAME}_documentation_${DATE}.md"
EOF

chmod +x output/scheduled-scan.sh

echo "✅ 定时任务脚本已创建: output/scheduled-scan.sh"
echo "   可以添加到crontab中实现定时备份"
echo "   示例: 0 2 * * * /path/to/scheduled-scan.sh"

echo ""
echo "📊 结果统计"
echo "=========="
echo "生成的文件:"
ls -la output/ | grep -E "\.(json|xml|md|toml)$"

echo ""
echo "文件大小统计:"
du -h output/* | sort -hr

echo ""
echo "✨ 示例完成！"
echo "============"
echo "📁 输出文件保存在: output/"
echo "📝 日志文件保存在: output/logs/"
echo ""
echo "💡 提示:"
echo "   - 使用 '-c' 参数指定配置文件"
echo "   - 使用 '-f' 参数指定输出格式"
echo "   - 使用 '-e' 参数排除文件/目录"
echo "   - 使用 '-s' 参数限制文件大小"
echo "   - 使用 '-d' 参数限制扫描深度"
echo "   - 使用 '--debug' 参数启用调试模式"
echo ""
echo "📚 更多帮助:"
echo "   ./code-context-generator --help"
echo "   ./code-context-generator generate --help"
echo "   cat docs/quickstart.md"