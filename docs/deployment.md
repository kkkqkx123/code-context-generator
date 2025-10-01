# 代码上下文生成器 - 部署文档

## 部署方式

### 1. 源码部署

#### 环境准备
```bash
# 安装Go 1.24+
go version
```

#### 构建应用
```bash
# 下载依赖
go mod download

# 构建CLI版本
go build -o code-context-generator cmd/cli/main.go

# 构建TUI版本  
go build -o code-context-generator-tui cmd/tui/main.go
```

#### 安装到系统路径
```bash
# Linux/macOS
sudo cp code-context-generator /usr/local/bin/
sudo chmod +x /usr/local/bin/code-context-generator*

# Windows
copy code-context-generator.exe C:\Windows\System32\
```

### 2. 二进制部署

#### 下载预编译二进制文件
```bash
# Linux
wget https://github.com/yourusername/code-context-generator/releases/download/v1.0.0/code-context-generator-linux-amd64.tar.gz
tar -xzf code-context-generator-linux-amd64.tar.gz

# Windows: 下载zip并解压
# macOS: 下载tar.gz并解压
```

### 3. Docker部署

#### 构建镜像
```bash
docker build -t code-context-generator:latest .
```

#### 运行容器
```bash
# 基本运行
docker run -it --rm code-context-generator:latest --help

# 挂载目录扫描
docker run -it --rm -v $(pwd):/workspace -w /workspace code-context-generator:latest generate .
```

## 验证部署

```bash
# 检查版本
./code-context-generator --version

# 测试基本功能
./code-context-generator generate --help
```

### Docker Compose 部署

#### 启动服务
```bash
docker-compose up -d

# 查看日志
docker-compose logs -f code-context-generator

# 停止服务
docker-compose down
```

## 环境配置

### 开发环境

#### 环境变量配置
```bash
# Linux/macOS: ~/.bashrc 或 ~/.zshrc
export CODE_CONTEXT_ENV=development
export CODE_CONTEXT_LOG_LEVEL=debug
export CODE_CONTEXT_CONFIG_PATH=~/projects/code-context-generator/config.toml

# Windows: 系统环境变量
setx CODE_CONTEXT_ENV development
setx CODE_CONTEXT_LOG_LEVEL debug
setx CODE_CONTEXT_CONFIG_PATH "C:\projects\code-context-generator\config.toml"
```

#### 开发配置文件
```toml
# config.development.toml
[output]
format = "json"
encoding = "utf-8"

[file_processing]
include_hidden = true
max_file_size = 52428800  # 50MB
max_depth = 10
exclude_patterns = [".git", "node_modules", "*.tmp"]
include_content = true
include_hash = true

[ui]
theme = "dark"
show_progress = true
show_preview = true

#### 智能格式覆盖配置
工具支持基于配置文件名的智能格式识别功能：
- `config-json.yaml` - 自动应用 JSON 格式配置
- `config-xml.yaml` - 自动应用 XML 格式配置
- `config-toml.yaml` - 自动应用 TOML 格式配置
- `config-markdown.yaml` - 自动应用 Markdown 格式配置

例如，创建 `config-json.yaml` 文件时，工具会自动设置 `output.format = "json"` 并应用 JSON 相关的配置选项。

[performance]
max_workers = 8
buffer_size = 4096
cache_enabled = true

[logging]
level = "debug"
file_path = "logs/development.log"
```

### 测试环境

#### 测试配置
```toml
# config.test.toml
[output]
format = "xml"
encoding = "utf-8"

[file_processing]
include_hidden = false
max_file_size = 10485760  # 10MB
max_depth = 5
exclude_patterns = [".git", "node_modules", "test_*"]
include_content = false

[performance]
max_workers = 2
buffer_size = 1024

[logging]
level = "info"
file_path = "logs/test.log"
```

### 生产环境

#### 生产配置
```toml
# config.production.toml
[output]
format = "json"
encoding = "utf-8"

[file_processing]
include_hidden = false
max_file_size = 5242880  # 5MB
max_depth = 3
exclude_patterns = [
    ".git", "node_modules", "*.tmp", "*.log",
    "vendor", "cache", "temp"
]
include_content = false
include_hash = false

[performance]
max_workers = 4
buffer_size = 2048
cache_enabled = true
cache_size = 200

[logging]
level = "warn"
file_path = "/var/log/code-context-generator/production.log"
max_size = 100
max_backups = 10
max_age = 30
```

## 监控和日志

### 日志配置

#### 日志轮转配置
```bash
# Linux: /etc/logrotate.d/code-context-generator
/var/log/code-context-generator/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 code-context code-context
    postrotate
        systemctl reload code-context-generator
    endscript
}
```

#### 日志格式
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "info",
  "component": "scanner",
  "message": "Scan started",
  "context": {
    "path": "/workspace",
    "format": "json",
    "options": {
      "max_depth": 3,
      "include_hidden": false
    }
  }
}
```

### 监控指标

#### Prometheus指标
```yaml
# metrics.yml
code_context_scanner_files_scanned_total 1250
code_context_scanner_folders_scanned_total 45
code_context_scanner_errors_total 2
code_context_scanner_duration_seconds 15.3
code_context_scanner_file_size_bytes 5242880
code_context_memory_usage_bytes 67108864
code_context_cpu_usage_percent 25.5
```

#### 健康检查端点
```bash
# HTTP健康检查（如果启用Web服务）
curl -f http://localhost:8080/health || exit 1

# 进程健康检查
ps aux | grep code-context-generator | grep -v grep

# 文件健康检查
[ -f /var/run/code-context-generator.pid ] && kill -0 $(cat /var/run/code-context-generator.pid)
```

## 备份和恢复

### 配置备份
```bash
# 备份配置文件
cp /etc/code-context-generator/config.toml /backup/code-context-generator-config-$(date +%Y%m%d).toml

# 备份日志文件
tar -czf /backup/code-context-generator-logs-$(date +%Y%m%d).tar.gz /var/log/code-context-generator/
```

### 数据备份
```bash
# 备份输出文件
cp /var/code-context-generator/output/*.json /backup/output/

# 备份缓存
cp -r /var/cache/code-context-generator /backup/cache/
```

## 安全考虑

### 文件权限
```bash
# 设置适当的文件权限
chmod 755 /usr/local/bin/code-context-generator
chmod 644 /etc/code-context-generator/config.toml
chmod 750 /var/log/code-context-generator/
chown -R code-context:code-context /opt/code-context-generator/
```

### 网络安全
- 限制网络访问（如果需要网络功能）
- 使用防火墙规则
- 定期更新依赖包
- 扫描安全漏洞

### 数据安全
- 加密敏感配置文件
- 限制日志文件访问
- 定期清理临时文件
- 备份重要数据

## 性能优化

### 系统调优
```bash
# Linux系统调优
# /etc/sysctl.conf
fs.file-max = 65536
vm.swappiness = 10
net.core.rmem_max = 16777216
net.core.wmem_max = 16777216

# 应用配置
sysctl -p
```

### 应用调优
```toml
# 性能优化配置
[performance]
max_workers = 8  # 根据CPU核心数调整
buffer_size = 8192  # 增大缓冲区
batch_size = 100  # 批处理大小
cache_size = 500  # 增大缓存
```

### 配置迁移
```bash
# 检查配置兼容性
./code-context-generator config validate --config /etc/code-context-generator/config.toml

# 更新配置格式（如果需要）
./code-context-generator config migrate --from /etc/code-context-generator/config.toml.bak --to /etc/code-context-generator/config.toml
```

## 支持信息

### 获取帮助
- 项目文档: [项目文档链接]
- 问题报告: [GitHub Issues](https://github.com/yourusername/code-context-generator/issues)
- 技术支持: support@example.com

### 系统信息收集
```bash
# 收集系统信息用于支持
./scripts/collect-system-info.sh > system-info.txt

# 收集应用日志
tar -czf app-logs.tar.gz /var/log/code-context-generator/
```

## 附录

### A. 系统要求检查脚本
```bash
#!/bin/bash
# check-requirements.sh

echo "检查系统要求..."

# 检查Go版本
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo "✓ Go版本: $GO_VERSION"
else
    echo "✗ Go未安装"
fi

# 检查内存
MEMORY=$(free -m | awk 'NR==2{print $2}')
if [ $MEMORY -ge 512 ]; then
    echo "✓ 内存: ${MEMORY}MB"
else
    echo "✗ 内存不足: ${MEMORY}MB (需要至少512MB)"
fi

# 检查磁盘空间
DISK=$(df -m . | awk 'NR==2{print $4}')
if [ $DISK -ge 100 ]; then
    echo "✓ 磁盘空间: ${DISK}MB"
else
    echo "✗ 磁盘空间不足: ${DISK}MB (需要至少100MB)"
fi

echo "系统要求检查完成"
```

### B. 快速部署脚本
```bash
#!/bin/bash
# quick-deploy.sh

set -e

echo "开始快速部署..."

# 检查系统要求
./scripts/check-requirements.sh

# 下载最新版本
LATEST_VERSION=$(curl -s https://api.github.com/repos/yourusername/code-context-generator/releases/latest | grep tag_name | cut -d '"' -f 4)
wget "https://github.com/yourusername/code-context-generator/releases/download/${LATEST_VERSION}/code-context-generator-linux-amd64.tar.gz"

# 解压和安装
tar -xzf code-context-generator-linux-amd64.tar.gz
sudo cp code-context-generator /usr/local/bin/
sudo cp code-context-generator-tui /usr/local/bin/

# 创建配置目录
mkdir -p ~/.config/code-context-generator

# 初始化配置
code-context-generator config init

echo "快速部署完成！"
echo "运行 'code-context-generator --help' 开始使用"
```