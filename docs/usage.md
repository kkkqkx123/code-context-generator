# 📖 使用文档

## 功能特性

### 核心功能
- **多格式输出**: 支持JSON、XML、Markdown、TOML格式
- **智能文件过滤**: 基于扩展名和内容分析的二进制文件检测
- **内容提取**: 可选择性包含文件内容
- **智能文件大小显示**: 根据文件大小自动选择B、KB、MB单位显示

### 高级特性

- **多线程处理**: 并行扫描提升性能
- **配置文件**: 支持TOML格式配置
- **路径匹配**: 支持通配符和正则表达式
- **智能去重**: 自动处理重复文件，确保每个文件只出现一次
- **精确文件夹统计**: 仅统计符合过滤条件的文件夹
- **编码控制**: 支持通过 `--encoding` 参数或 `.env` 文件设置输出文件编码格式（默认 utf-8）
- **元信息控制**: 支持通过 `.env` 配置决定输出文件是否包含元信息（默认不包含，仅保留内容、名称和相对路径）

### 环境变量配置

支持通过 `.env` 文件配置以下环境变量：

```bash
# 输出配置
CODE_CONTEXT_DEFAULT_FORMAT=json          # 默认输出格式
CODE_CONTEXT_OUTPUT_DIR=./output          # 输出目录
CODE_CONTEXT_FILENAME_TEMPLATE=context    # 文件名模板
CODE_CONTEXT_TIMESTAMP_FORMAT=2006-01-02_15-04-05  # 时间戳格式
CODE_CONTEXT_ENCODING=utf-8               # 输出文件编码格式
CODE_CONTEXT_INCLUDE_METADATA=false       # 是否包含元信息（大小、修改时间等）

# 过滤配置
CODE_CONTEXT_MAX_DEPTH=0                  # 最大扫描深度（0表示只扫描当前目录，1表示递归1层，-1表示无限制）
CODE_CONTEXT_MAX_FILE_SIZE=10MB          # 最大文件大小
CODE_CONTEXT_EXCLUDE_PATTERNS=.git,node_modules  # 排除模式
CODE_CONTEXT_INCLUDE_PATTERNS=           # 包含模式
CODE_CONTEXT_RECURSIVE=false             # 是否递归（已废弃，使用MAX_DEPTH控制）
CODE_CONTEXT_EXCLUDE_BINARY=false        # 是否排除二进制文件

# 文件处理配置
CODE_CONTEXT_INCLUDE_HIDDEN=false        # 是否包含隐藏文件
CODE_CONTEXT_INCLUDE_CONTENT=true        # 是否包含文件内容
CODE_CONTEXT_INCLUDE_HASH=false          # 是否包含文件哈希

# 安全扫描配置
CODE_CONTEXT_SECURITY_ENABLED=false      # 是否启用安全扫描（默认禁用）
CODE_CONTEXT_SECURITY_FAIL_ON_CRITICAL=false  # 发现严重问题时是否失败
CODE_CONTEXT_SECURITY_SCAN_LEVEL=standard  # 扫描级别（basic, standard, comprehensive）
CODE_CONTEXT_SECURITY_REPORT_FORMAT=text  # 报告格式（text, json, xml, html）
CODE_CONTEXT_SECURITY_DETECT_CREDENTIALS=false  # 是否检测硬编码凭证
CODE_CONTEXT_SECURITY_DETECT_SQL_INJECTION=false  # 是否检测SQL注入漏洞
CODE_CONTEXT_SECURITY_DETECT_XSS=false  # 是否检测XSS漏洞
CODE_CONTEXT_SECURITY_DETECT_PATH_TRAVERSAL=false  # 是否检测路径遍历漏洞
CODE_CONTEXT_SECURITY_DETECT_QUALITY=false  # 是否检测代码质量问题
```

## 安装

### 前置要求
- Go 1.24或更高版本
- Git（用于源码安装）

### 从源码安装
```bash
git clone https://github.com/yourusername/code-context-generator.git
cd code-context-generator
go build -o c-gen.exe cli/main.go
```

## CLI使用指南

### 基本用法
```bash
# 扫描当前目录
./c-gen generate

# 扫描指定目录
./c-gen generate /path/to/project

# 指定输出格式
./c-gen generate -f markdown -o output.md
```

### 高级用法
```bash
# 包含文件内容
./c-gen generate -C -o context.json

# 排除特定文件
./c-gen generate -e "*.log" -e "node_modules"

# 只包含特定扩展名
./c-gen generate -i "*.go" -i "*.md"

# 限制文件大小（10MB）
./c-gen generate -s 10485760

# 限制扫描深度
./c-gen generate -d 3

# 深度参数说明：
# -d 0: 只扫描当前目录，不递归子目录
# -d 1: 递归1层子目录
# -d -1: 无限递归（不限制深度）
```



### 配置管理
```bash
# 创建默认配置
./c-gen config init

# 使用自定义配置
./c-gen generate -c config.toml

# 使用智能格式覆盖配置
./c-gen generate -c config-json.yaml  # 自动使用JSON格式
./c-gen generate -c config-xml.yaml  # 自动使用XML格式
```

## 配置文件

### 智能格式覆盖
工具支持基于配置文件名的智能格式识别功能。当配置文件名包含特定格式关键词时，会自动应用对应的格式配置：

- `config-json.yaml` - 自动设置 `output.format = "json"`
- `config-xml.yaml` - 自动设置 `output.format = "xml"`
- `config-toml.yaml` - 自动设置 `output.format = "toml"`
- `config-markdown.yaml` 或 `config-md.yaml` - 自动设置 `output.format = "markdown"`

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

### generate命令(参数缺省时默认使用generate)
- `-f, --format`: 输出格式（json, xml, markdown, toml）
- `-o, --output`: 输出文件路径
- `-C, --content`: 包含文件内容
- `-H, --hash`: 包含文件哈希值
- `-e, --exclude`: 排除模式（可多次使用）
- `-i, --include`: 包含模式（可多次使用）
- `-s, --max-size`: 最大文件大小
- `-d, --max-depth`: 最大扫描深度（0表示只扫描当前目录，1表示递归1层，-1表示无限制）
- `-c, --config`: 配置文件路径
- `--encoding`: 输出文件编码格式（默认：utf-8）

### config命令
- `init`: 创建默认配置文件
- `validate`: 验证配置文件

## 实用示例

### 扫描Go项目
```bash
./c-gen generate -e "vendor" -f json -o go-project.json
```

### 扫描Python项目
```bash
./c-gen generate -e "venv" -e "__pycache__" -f markdown -o python-project.md
```

### 生成项目文档
```bash
./c-gen generate -C -H -f markdown -o documentation.md
```

## 故障排除

### 常见问题

**权限错误**: `permission denied`
```bash
chmod +x c-gen  # Linux/macOS
```

**找不到命令**: `command not found`
```bash
./c-gen  # 使用完整路径
```

**输出文件太大**: 
```bash
./c-gen generate -s 1048576  # 限制文件大小
```

### 调试模式
```bash
./c-gen generate --debug
```

### 获取帮助
```bash
./c-gen --help
./c-gen generate --help
```