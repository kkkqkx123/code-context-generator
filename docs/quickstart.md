# 🚀 快速入门指南

## 安装

### 从源码编译
```bash
git clone https://github.com/yourusername/code-context-generator.git
cd code-context-generator
go build -o c-gen cmd/cli/main.go
```

## 基础使用

### 扫描当前目录
```bash
./c-gen generate
```

### 扫描指定目录
```bash
./c-gen generate /path/to/your/project
```

### 指定输出格式
```bash
# Markdown格式
./c-gen generate -f markdown -o project.md

# XML格式
./c-gen generate -f xml -o project.xml
```

### 智能格式覆盖
```bash
# 使用config-json.yaml自动应用JSON格式
./c-gen generate -c config-json.yaml

# 使用config-xml.yaml自动应用XML格式  
./c-gen generate -c config-xml.yaml
```

### 包含文件内容
```bash
./c-gen generate -C -o context.json
```

## 常用命令

### 文件过滤
```bash
# 排除特定文件
./c-gen generate -e "*.log" -e "node_modules"

# 只包含特定扩展名
./c-gen generate -i "*.go" -i "*.md"

# 限制文件大小（10MB）
./c-gen generate -s 10485760
```



## 配置文件

### 创建默认配置
```bash
./c-gen config init
```

### 基础配置示例
```toml
[output]
format = "json"

[file_processing]
max_file_size = 10485760  # 10MB
exclude_patterns = ["*.log", "node_modules", ".git"]
```

### 智能格式覆盖配置
工具支持基于配置文件名的智能格式识别：
- `config-json.yaml` - 自动应用 JSON 格式配置
- `config-xml.yaml` - 自动应用 XML 格式配置
- `config-toml.yaml` - 自动应用 TOML 格式配置
- `config-markdown.yaml` - 自动应用 Markdown 格式配置

例如，创建 `config-json.yaml` 文件时，工具会自动设置 `output.format = "json"` 并应用 JSON 相关的配置选项。