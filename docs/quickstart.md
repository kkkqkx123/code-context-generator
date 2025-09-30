# 🚀 快速入门指南

本指南将帮助你在5分钟内快速上手使用代码上下文生成器。

## 📋 目录
- [安装](#安装)
- [基础使用](#基础使用)
- [常用命令](#常用命令)
- [配置文件](#配置文件)
- [故障排除](#故障排除)

## 安装

### 方式1：使用预编译二进制文件（推荐）
```bash
# Windows
wget https://github.com/yourusername/code-context-generator/releases/download/v1.0.0/code-context-generator-windows-amd64.exe

# Linux
wget https://github.com/yourusername/code-context-generator/releases/download/v1.0.0/code-context-generator-linux-amd64

# macOS
wget https://github.com/yourusername/code-context-generator/releases/download/v1.0.0/code-context-generator-darwin-amd64
```

### 方式2：从源码编译
```bash
# 克隆仓库
git clone https://github.com/yourusername/code-context-generator.git
cd code-context-generator

# 编译
go build -o code-context-generator cmd/cli/main.go

# 编译TUI版本
go build -o code-context-generator-tui cmd/tui/main.go
```

### 方式3：使用Go安装
```bash
go install github.com/yourusername/code-context-generator/cmd/cli@latest
go install github.com/yourusername/code-context-generator/cmd/tui@latest
```

## 基础使用

### 1. 扫描当前目录（最简单用法）
```bash
./code-context-generator generate
```
输出示例：
```json
{
  "files": [
    {
      "path": "README.md",
      "size": 1024,
      "modified": "2024-01-01T10:00:00Z"
    },
    {
      "path": "main.go",
      "size": 2048,
      "modified": "2024-01-01T09:30:00Z"
    }
  ],
  "total_files": 2,
  "total_size": 3072
}
```

### 2. 扫描指定目录
```bash
./code-context-generator generate /path/to/your/project
```

### 3. 指定输出格式
```bash
# 输出为Markdown格式
./code-context-generator generate -f markdown -o project.md

# 输出为XML格式
./code-context-generator generate -f xml -o project.xml

# 输出为TOML格式
./code-context-generator generate -f toml -o project.toml
```

### 4. 包含文件内容
```bash
# 包含文件内容
./code-context-generator generate -C -o context.json

# 同时包含内容和哈希值
./code-context-generator generate -C -H -o context.json
```

## 常用命令

### 文件过滤
```bash
# 排除特定文件/目录
./code-context-generator generate -e "*.log" -e "node_modules" -e ".git"

# 只包含特定扩展名
./code-context-generator generate -i "*.go" -i "*.md"

# 限制文件大小（10MB）
./code-context-generator generate -s 10485760

# 限制扫描深度（2层）
./code-context-generator generate -d 2
```

### 交互式选择
```bash
# 启动交互式文件选择器
./code-context-generator select

# 多选模式
./code-context-generator select -m -f json -o selected.json
```

### TUI界面
```bash
# 启动TUI界面
./code-context-generator-tui
```

### 使用配置文件
```bash
# 使用自定义配置文件
./code-context-generator generate -c myconfig.toml

# 生成默认配置文件
./code-context-generator config init
```

## 配置文件

### 创建默认配置
```bash
./code-context-generator config init
```

### 基础配置示例（config.toml）
```toml
[output]
format = "json"
encoding = "utf-8"

[file_processing]
include_hidden = false
max_file_size = 10485760  # 10MB
max_depth = 0  # 无限制
exclude_patterns = [
    "*.exe", "*.dll", "*.so", "*.dylib",
    "*.pyc", "*.pyo", "*.pyd",
    "node_modules", ".git", ".svn", ".hg",
    "__pycache__", "*.egg-info", "dist", "build"
]

[ui]
theme = "default"
show_progress = true
```

### 高级配置示例
```toml
[output]
format = "json"
encoding = "utf-8"
pretty = true
output_dir = "./output"
filename_template = "context_{{.Timestamp}}.{{.Format}}"

[file_processing]
include_hidden = false
include_content = true
include_hash = true
max_file_size = 52428800  # 50MB
max_depth = 5
buffer_size = 8192
workers = 4
exclude_patterns = [
    "*.exe", "*.dll", "*.so", "*.dylib",
    "*.pyc", "*.pyo", "*.pyd",
    "node_modules", ".git", ".svn", ".hg",
    "__pycache__", "*.egg-info", "dist", "build",
    "*.log", "*.tmp", "*.temp", "*.cache"
]

[formats.json]
enabled = true
indent = "  "
sort_keys = true

[formats.xml]
enabled = true
indent = "  "
use_cdata = false

[formats.toml]
enabled = true
indent = "  "

[formats.markdown]
enabled = true
template = "default"
include_toc = true

[ui]
theme = "default"
show_progress = true
auto_complete = true
confirm_deletion = true
```

## 实用示例

### 示例1：扫描Go项目
```bash
# 扫描Go项目，排除vendor目录
./code-context-generator generate -e "vendor" -e "*.exe" -f json -o go-project.json
```

### 示例2：扫描Python项目
```bash
# 扫描Python项目，排除虚拟环境和缓存
./code-context-generator generate \
  -e "venv" -e "__pycache__" -e "*.pyc" \
  -e "*.egg-info" -e "dist" -e "build" \
  -f markdown -o python-project.md
```

### 示例3：扫描Web项目
```bash
# 扫描前端项目，排除node_modules和构建产物
./code-context-generator generate \
  -e "node_modules" -e "dist" -e "build" \
  -e "*.min.js" -e "*.min.css" \
  -f xml -o web-project.xml
```

### 示例4：生成项目文档
```bash
# 生成包含内容的完整项目文档
./code-context-generator generate \
  -C -H -f markdown \
  -e "node_modules" -e ".git" -e "*.log" \
  -o project-documentation.md
```

### 示例5：快速备份文件列表
```bash
# 生成文件列表用于备份
./code-context-generator generate -f json -o backup-list.json
```

## 故障排除

### 常见问题

#### Q: 权限错误
**问题**: `permission denied`
**解决**: 
```bash
# Linux/macOS
chmod +x code-context-generator

# Windows
# 确保文件没有被系统阻止
```

#### Q: 找不到命令
**问题**: `command not found`
**解决**: 
```bash
# 添加到PATH或指定完整路径
./code-context-generator

# 或移动到系统目录
sudo mv code-context-generator /usr/local/bin/
```

#### Q: 输出文件太大
**问题**: 生成的文件太大
**解决**: 
```bash
# 限制文件大小
./code-context-generator generate -s 1048576  # 1MB

# 排除大文件
./code-context-generator generate -e "*.mp4" -e "*.zip"

# 限制扫描深度
./code-context-generator generate -d 3
```

#### Q: 扫描速度太慢
**问题**: 扫描大型项目很慢
**解决**: 
```bash
# 增加工作线程数
./code-context-generator generate --workers 8

# 排除不必要的目录
./code-context-generator generate -e "node_modules" -e ".git" -e "vendor"

# 使用缓存（如果支持）
./code-context-generator generate --cache
```

#### Q: 格式错误
**问题**: 输出格式不正确
**解决**: 
```bash
# 检查配置文件格式
./code-context-generator config validate

# 使用默认配置
./code-context-generator generate -c default.toml
```

### 调试模式
```bash
# 启用调试模式
./code-context-generator generate --debug

# 查看详细日志
./code-context-generator generate -v -v  # 最高详细级别
```

### 获取帮助
```bash
# 查看帮助
./code-context-generator --help

# 查看子命令帮助
./code-context-generator generate --help
./code-context-generator select --help

# 查看版本信息
./code-context-generator --version
```

## 🎯 下一步

完成快速入门后，你可以：

1. **深入学习** - 阅读[完整使用文档](usage.md)
2. **部署应用** - 查看[部署文档](deployment.md)
3. **参与开发** - 阅读[开发环境文档](development.md)
4. **高级配置** - 探索更多配置选项
5. **性能优化** - 学习如何优化扫描性能

## 📞 获取帮助

如果在使用过程中遇到问题：

1. 查看[使用文档](usage.md)中的详细说明
2. 查看[故障排除](#故障排除)部分
3. 提交Issue到项目仓库
4. 参与社区讨论

---

*🎉 恭喜！现在你已经掌握了代码上下文生成器的基本使用方法。开始探索更多高级功能吧！*