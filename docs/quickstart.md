# 🚀 快速入门指南

## 安装

### 从源码编译
```bash
git clone https://github.com/yourusername/code-context-generator.git
cd code-context-generator
go build -o code-context-generator cmd/cli/main.go
```

## 基础使用

### 扫描当前目录
```bash
./code-context-generator generate
```

### 扫描指定目录
```bash
./code-context-generator generate /path/to/your/project
```

### 指定输出格式
```bash
# Markdown格式
./code-context-generator generate -f markdown -o project.md

# XML格式
./code-context-generator generate -f xml -o project.xml
```

### 包含文件内容
```bash
./code-context-generator generate -C -o context.json
```

## 常用命令

### 文件过滤
```bash
# 排除特定文件
./code-context-generator generate -e "*.log" -e "node_modules"

# 只包含特定扩展名
./code-context-generator generate -i "*.go" -i "*.md"

# 限制文件大小（10MB）
./code-context-generator generate -s 10485760
```

### 交互式选择
```bash
# 启动文件选择器
./code-context-generator select
```

### TUI界面
```bash
./code-context-generator-tui
```

## 配置文件

### 创建默认配置
```bash
./code-context-generator config init
```

### 基础配置示例
```toml
[output]
format = "json"

[file_processing]
max_file_size = 10485760  # 10MB
exclude_patterns = ["*.log", "node_modules", ".git"]
```