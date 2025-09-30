# 代码上下文生成器 - 部署文档

## 概述

本文档描述了代码上下文生成器在不同环境下的部署方案，包括开发环境、测试环境和生产环境的部署步骤。

## 部署方式

### 1. 源码部署

#### 环境准备
```bash
# 安装Go 1.24+
# Windows: 下载安装包从 https://golang.org/dl/
# Linux/macOS:
wget https://golang.org/dl/go1.24.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 验证安装
go version
```

#### 获取源码
```bash
# 克隆仓库
git clone https://github.com/yourusername/code-context-generator.git
cd code-context-generator

# 或者下载发布版本
wget https://github.com/yourusername/code-context-generator/releases/download/v1.0.0/source-code.tar.gz
tar -xzf source-code.tar.gz
cd code-context-generator-1.0.0
```

#### 构建应用
```bash
# 下载依赖
go mod download

# 构建CLI版本
go build -o code-context-generator cmd/cli/main.go

# 构建TUI版本
go build -o code-context-generator-tui cmd/tui/main.go

# 构建所有版本
make build-all  # 如果有Makefile
```

#### 安装到系统路径
```bash
# Linux/macOS
sudo cp code-context-generator /usr/local/bin/
sudo cp code-context-generator-tui /usr/local/bin/
sudo chmod +x /usr/local/bin/code-context-generator*

# Windows (以管理员身份运行PowerShell)
copy code-context-generator.exe C:\Windows\System32\
copy code-context-generator-tui.exe C:\Windows\System32\
```

### 2. 二进制部署

#### 下载预编译二进制文件
```bash
# Linux
wget https://github.com/yourusername/code-context-generator/releases/download/v1.0.0/code-context-generator-linux-amd64.tar.gz
tar -xzf code-context-generator-linux-amd64.tar.gz

# macOS
wget https://github.com/yourusername/code-context-generator/releases/download/v1.0.0/code-context-generator-darwin-amd64.tar.gz
tar -xzf code-context-generator-darwin-amd64.tar.gz

# Windows
# 下载 code-context-generator-windows-amd64.zip 并解压
```

#### 安装二进制文件
```bash
# Linux/macOS
sudo cp code-context-generator /usr/local/bin/
sudo cp code-context-generator-tui /usr/local/bin/

# Windows
# 将exe文件复制到系统PATH目录，如 C:\Windows\System32\
```

### 3. 容器化部署

#### Docker部署

##### 创建Dockerfile
```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o code-context-generator cmd/cli/main.go
RUN go build -o code-context-generator-tui cmd/tui/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/code-context-generator .
COPY --from=builder /app/code-context-generator-tui .

# 创建配置目录
RUN mkdir -p /root/.config/code-context-generator

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口（如果需要Web服务）
# EXPOSE 8080

CMD ["./code-context-generator"]
```

##### 构建镜像
```bash
# 构建镜像
docker build -t code-context-generator:latest .

# 标记版本
docker tag code-context-generator:latest code-context-generator:v1.0.0
```

##### 运行容器
```bash
# 基本运行
docker run -it --rm code-context-generator:latest --help

# 挂载当前目录进行扫描
docker run -it --rm \
  -v $(pwd):/workspace \
  -w /workspace \
  code-context-generator:latest generate . -f json

# 挂载配置文件
docker run -it --rm \
  -v $(pwd):/workspace \
  -v ~/.config/code-context-generator:/root/.config/code-context-generator \
  -w /workspace \
  code-context-generator:latest generate . --config /root/.config/code-context-generator/config.toml

# 运行TUI版本（需要TTY）
docker run -it --rm \
  --device=/dev/tty \
  -v $(pwd):/workspace \
  -w /workspace \
  code-context-generator-tui:latest
```

#### Docker Compose部署

##### docker-compose.yml
```yaml
version: '3.8'

services:
  code-context-generator:
    image: code-context-generator:latest
    container_name: code-context-generator
    volumes:
      - ./workspace:/workspace
      - ./config:/root/.config/code-context-generator
    working_dir: /workspace
    environment:
      - TZ=Asia/Shanghai
      - CODE_CONTEXT_FORMAT=json
      - CODE_CONTEXT_MAX_DEPTH=3
    stdin_open: true
    tty: true
    command: ["./code-context-generator"]
    
  code-context-generator-tui:
    image: code-context-generator-tui:latest
    container_name: code-context-generator-tui
    volumes:
      - ./workspace:/workspace
      - ./config:/root/.config/code-context-generator
    working_dir: /workspace
    environment:
      - TZ=Asia/Shanghai
    stdin_open: true
    tty: true
    command: ["./code-context-generator-tui"]
```

##### 启动服务
```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f code-context-generator

# 停止服务
docker-compose down
```

### 4. 系统服务部署

#### Linux Systemd服务

##### 创建服务文件
```ini
# /etc/systemd/system/code-context-generator.service
[Unit]
Description=Code Context Generator Service
After=network.target

[Service]
Type=simple
User=code-context
Group=code-context
WorkingDirectory=/opt/code-context-generator
ExecStart=/usr/local/bin/code-context-generator generate /var/projects --config /etc/code-context-generator/config.toml
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

# 环境变量
Environment=CODE_CONTEXT_FORMAT=json
Environment=CODE_CONTEXT_MAX_DEPTH=5

# 资源限制
LimitNOFILE=65536
MemoryLimit=1G
CPUQuota=50%

[Install]
WantedBy=multi-user.target
```

##### 创建系统用户
```bash
# 创建专用用户
sudo useradd -r -s /bin/false -d /opt/code-context-generator code-context

# 创建目录
sudo mkdir -p /opt/code-context-generator
sudo mkdir -p /etc/code-context-generator
sudo mkdir -p /var/log/code-context-generator

# 设置权限
sudo chown -R code-context:code-context /opt/code-context-generator
sudo chown -R code-context:code-context /var/log/code-context-generator
```

##### 安装应用
```bash
# 复制二进制文件
sudo cp code-context-generator /usr/local/bin/
sudo cp code-context-generator-tui /usr/local/bin/

# 复制配置文件
sudo cp config.toml /etc/code-context-generator/
sudo chown code-context:code-context /etc/code-context-generator/config.toml
sudo chmod 644 /etc/code-context-generator/config.toml
```

##### 启动服务
```bash
# 重新加载systemd
sudo systemctl daemon-reload

# 启用服务
sudo systemctl enable code-context-generator

# 启动服务
sudo systemctl start code-context-generator

# 查看状态
sudo systemctl status code-context-generator

# 查看日志
sudo journalctl -u code-context-generator -f
```

#### Windows服务部署

##### 使用NSSM创建服务
```powershell
# 下载NSSM
# https://nssm.cc/download

# 安装服务
nssm install CodeContextGenerator

# 配置服务
nssm set CodeContextGenerator Application "C:\Program Files\CodeContextGenerator\code-context-generator.exe"
nssm set CodeContextGenerator AppParameters "generate C:\Projects --config C:\ProgramData\CodeContextGenerator\config.toml"
nssm set CodeContextGenerator DisplayName "Code Context Generator"
nssm set CodeContextGenerator Description "Code Context Generator Service"
nssm set CodeContextGenerator Start SERVICE_AUTO_START

# 启动服务
net start CodeContextGenerator
```

### 5. Kubernetes部署

#### 创建ConfigMap
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: code-context-generator-config
data:
  config.toml: |
    [output]
    format = "json"
    encoding = "utf-8"
    
    [file_processing]
    include_hidden = false
    max_file_size = 10485760
    max_depth = 5
    exclude_patterns = ["*.tmp", "*.log", ".git"]
```

#### 创建Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: code-context-generator
spec:
  replicas: 1
  selector:
    matchLabels:
      app: code-context-generator
  template:
    metadata:
      labels:
        app: code-context-generator
    spec:
      containers:
      - name: code-context-generator
        image: code-context-generator:latest
        command: ["./code-context-generator"]
        args: ["generate", "/workspace", "--config", "/config/config.toml"]
        volumeMounts:
        - name: config
          mountPath: /config
        - name: workspace
          mountPath: /workspace
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
      volumes:
      - name: config
        configMap:
          name: code-context-generator-config
      - name: workspace
        persistentVolumeClaim:
          claimName: workspace-pvc
```

#### 创建Service（如果需要暴露服务）
```yaml
apiVersion: v1
kind: Service
metadata:
  name: code-context-generator-service
spec:
  selector:
    app: code-context-generator
  ports:
  - protocol: TCP
    port: 8080
    targetPort: 8080
  type: ClusterIP
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

## 故障排除

### 常见问题

#### 服务启动失败
```bash
# 检查日志
journalctl -u code-context-generator -n 50

# 检查配置文件
./code-context-generator config validate --config /etc/code-context-generator/config.toml

# 检查权限
ls -la /usr/local/bin/code-context-generator
ls -la /etc/code-context-generator/
```

#### 性能问题
```bash
# 监控系统资源
top -p $(pgrep code-context-generator)
iostat -x 1
vmstat 1

# 检查内存使用
pmap -x $(pgrep code-context-generator)

# 检查文件描述符
lsof -p $(pgrep code-context-generator)
```

#### 容器问题
```bash
# 查看容器日志
docker logs code-context-generator

# 进入容器调试
docker exec -it code-context-generator /bin/sh

# 检查容器资源使用
docker stats code-context-generator
```

### 恢复程序

#### 服务恢复
```bash
# 重启服务
sudo systemctl restart code-context-generator

# 重新加载配置
sudo systemctl reload code-context-generator

# 查看服务状态
sudo systemctl status code-context-generator
```

#### 数据恢复
```bash
# 从备份恢复配置
cp /backup/code-context-generator-config-20240101.toml /etc/code-context-generator/config.toml

# 恢复权限
sudo chown code-context:code-context /etc/code-context-generator/config.toml
sudo chmod 644 /etc/code-context-generator/config.toml

# 重启服务
sudo systemctl restart code-context-generator
```

## 更新和升级

### 应用更新
```bash
# 备份当前配置
cp /etc/code-context-generator/config.toml /etc/code-context-generator/config.toml.bak

# 停止服务
sudo systemctl stop code-context-generator

# 更新应用
wget https://github.com/yourusername/code-context-generator/releases/download/v1.1.0/code-context-generator-linux-amd64.tar.gz
tar -xzf code-context-generator-linux-amd64.tar.gz
sudo cp code-context-generator /usr/local/bin/

# 重启服务
sudo systemctl start code-context-generator

# 验证更新
./code-context-generator --version
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