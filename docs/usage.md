# 📖 使用文档

## 功能特性

### 核心功能
- **多格式输出**: 支持JSON、XML、Markdown、TOML格式
- **智能文件过滤**: 基于扩展名和内容分析的二进制文件检测
- **内容提取**: 可选择性包含文件内容
- **交互式选择**: TUI界面支持文件选择

### 高级特性
- **多线程处理**: 并行扫描提升性能
- **配置文件**: 支持TOML格式配置
- **路径匹配**: 支持通配符和正则表达式

## 安装

### 前置要求
- Go 1.24或更高版本
- Git（用于源码安装）

### 从源码安装
```bash
git clone https://github.com/yourusername/code-context-generator.git
cd code-context-generator
go build -o code-context-generator cmd/cli/main.go
go build -o code-context-generator-tui cmd/tui/main.go
```

## CLI使用指南

### 基本用法
```bash
# 扫描当前目录
./code-context-generator generate

# 扫描指定目录
./code-context-generator generate /path/to/project

# 指定输出格式
./code-context-generator generate -f markdown -o output.md
```

### 高级用法
```bash
# 包含文件内容
./code-context-generator generate -C -o context.json

# 排除特定文件
./code-context-generator generate -e "*.log" -e "node_modules"

# 只包含特定扩展名
./code-context-generator generate -i "*.go" -i "*.md"

# 限制文件大小（10MB）
./code-context-generator generate -s 10485760

# 限制扫描深度
./code-context-generator generate -d 3
```

### 交互式选择
```bash
# 启动文件选择器
./code-context-generator select

# 多选模式
./code-context-generator select -m -f json -o selected.json
```

### 配置管理
```bash
# 创建默认配置
./code-context-generator config init

# 使用自定义配置
./code-context-generator generate -c config.toml
```

## 配置文件

### 基础配置
```toml
[output]
format = "json"
encoding = "utf-8"

[file_processing]
max_file_size = 10485760  # 10MB
exclude_patterns = ["*.log", "node_modules", ".git"]
```

### 高级配置
```toml
[output]
format = "json"
pretty = true

[file_processing]
include_content = true
include_hash = true
max_file_size = 52428800  # 50MB
max_depth = 5
workers = 4
exclude_patterns = [
    "*.exe", "*.dll", "*.so", "*.dylib",
    "*.pyc", "*.pyo", "*.pyd",
    "node_modules", ".git", ".svn", ".hg",
    "__pycache__", "*.egg-info", "dist", "build"
]

[formats.json]
indent = "  "
sort_keys = true

[formats.markdown]
include_toc = true
```

## 命令参数详解

### generate命令
- `-f, --format`: 输出格式（json, xml, markdown, toml）
- `-o, --output`: 输出文件路径
- `-C, --content`: 包含文件内容
- `-H, --hash`: 包含文件哈希值
- `-e, --exclude`: 排除模式（可多次使用）
- `-i, --include`: 包含模式（可多次使用）
- `-s, --max-size`: 最大文件大小
- `-d, --max-depth`: 最大扫描深度
- `-c, --config`: 配置文件路径

### select命令
- `-m, --multi`: 多选模式
- `-f, --format`: 输出格式
- `-o, --output`: 输出文件路径

### config命令
- `init`: 创建默认配置文件
- `validate`: 验证配置文件

## 实用示例

### 扫描Go项目
```bash
./code-context-generator generate -e "vendor" -f json -o go-project.json
```

### 扫描Python项目
```bash
./code-context-generator generate -e "venv" -e "__pycache__" -f markdown -o python-project.md
```

### 生成项目文档
```bash
./code-context-generator generate -C -H -f markdown -o documentation.md
```

## 故障排除

### 常见问题

**权限错误**: `permission denied`
```bash
chmod +x code-context-generator  # Linux/macOS
```

**找不到命令**: `command not found`
```bash
./code-context-generator  # 使用完整路径
```

**输出文件太大**: 
```bash
./code-context-generator generate -s 1048576  # 限制文件大小
```

### 调试模式
```bash
./code-context-generator generate --debug
```

### 获取帮助
```bash
./code-context-generator --help
./code-context-generator generate --help
```