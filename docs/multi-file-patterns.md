# 多文件处理和模式文件使用指南

本文档详细介绍如何使用多文件处理功能和模式文件功能。

## 多文件处理

使用 `-m, --multiple-files` 参数可以指定多个文件进行处理。

### 基本用法

```bash
# 指定多个文件
code-context-generator generate -m file1.go -m file2.go -m file3.go

# 指定输出格式和文件名
code-context-generator generate -m main.go -m utils.go -f markdown -o output.md
```

### 高级用法

```bash
# 结合模式文件过滤多个文件
code-context-generator generate -m src/main.go -m src/utils.go -p patterns.txt

# 结合其他参数使用
code-context-generator generate -m config.json -m readme.md -C -H -f json
```

### 注意事项

1. 当指定多个文件时，程序会忽略目录扫描，只处理指定的文件
2. 输出文件名默认基于第一个文件的名称生成
3. 支持绝对路径和相对路径
4. 可以与 `-p, --pattern-file` 参数结合使用进行进一步过滤

## 模式文件使用

使用 `-p, --pattern-file` 参数可以从文件加载过滤模式。

### 模式文件格式

模式文件支持以下格式：

```
# 注释行以 # 开头
*.go          # 匹配所有 .go 文件
*.json        # 匹配所有 .json 文件

# 目录匹配
src/*.js      # 匹配 src 目录下的 .js 文件
test/**/*_test.go  # 匹配 test 目录下所有子目录的测试文件

# 通配符匹配
test_*        # 匹配以 test_ 开头的文件
*_test.*      # 匹配以 _test 结尾的文件名

# 路径分隔符兼容（支持 Windows 和 Linux 格式）
test_files\data.*    # Windows 格式
test_files/config.*  # Linux 格式
```

### 使用示例

```bash
# 使用模式文件过滤
code-context-generator generate -p patterns.txt

# 结合目录扫描使用
code-context-generator generate /path/to/project -p patterns.txt

# 结合多文件处理使用
code-context-generator generate -m file1.go -m file2.go -p patterns.txt
```

### 模式文件示例

创建 `patterns.txt` 文件：

```
# Go 项目模式文件
*.go          # 包含所有 Go 文件
*.mod         # 包含 go.mod 文件
*.sum         # 包含 go.sum 文件

# 排除测试文件（如果需要）
# *_test.go    # 注释掉以排除测试文件

# 配置文件
*.json        # JSON 配置文件
*.yaml        # YAML 配置文件
*.yml         # YML 配置文件
*.toml        # TOML 配置文件

# 文档文件
*.md          # Markdown 文档
*.txt         # 文本文件
readme*       # README 文件

# 构建文件
Makefile      # Makefile
Dockerfile    # Docker 文件
*.dockerfile  # Docker 文件

# 排除目录
# .git/        # Git 目录
# node_modules/ # Node.js 依赖
# vendor/      # Go vendor 目录
```

### 路径格式兼容性

程序支持 Windows 和 Linux 两种路径格式：

```bash
# Windows 格式
code-context-generator generate -p patterns_windows.txt

# Linux 格式  
code-context-generator generate -p patterns_linux.txt

# 混合格式（在同一模式文件中）
code-context-generator generate -p patterns_mixed.txt
```

模式文件中的路径分隔符会被自动处理，确保在不同操作系统上都能正确工作。

## 组合使用示例

### 场景 1：处理特定源文件

```bash
# 只处理主要的 Go 源文件
code-context-generator generate \
  -m main.go \
  -m cmd/server.go \
  -m internal/config/config.go \
  -m internal/handlers/*.go \
  -f markdown \
  -o project-core.md
```

### 场景 2：使用模式文件进行精确过滤

创建 `core-patterns.txt`：
```
# 核心代码文件
*.go
!*_test.go     # 排除测试文件
!*/mock/*      # 排除 mock 目录

# 配置文件
*.json
*.yaml
config.toml

# 重要文档
*.md
LICENSE
```

然后使用：
```bash
code-context-generator generate /path/to/project -p core-patterns.txt -o core-docs.md
```

### 场景 3：多文件 + 模式文件

```bash
# 处理多个指定文件，并用模式文件进一步过滤
code-context-generator generate \
  -m src/main.go \
  -m src/utils.go \
  -m src/config.json \
  -p include-patterns.txt \
  -f json \
  -o selected-files.json
```

## 常见问题

### Q: 多文件处理和目录扫描有什么区别？
A: 多文件处理只处理指定的文件，不会扫描目录；目录扫描会递归遍历整个目录结构。

### Q: 模式文件和命令行模式参数有什么区别？
A: 模式文件支持更复杂的模式定义，可以包含注释，支持路径格式兼容，适合保存常用的过滤规则。

### Q: 可以同时使用多个模式文件吗？
A: 目前不支持同时使用多个模式文件，但可以在一个模式文件中定义所有需要的模式。

### Q: 模式文件支持哪些通配符？
A: 支持标准的 glob 通配符：`*` 匹配任意字符，`?` 匹配单个字符，`[abc]` 匹配字符集等。

### Q: 如何调试模式文件是否工作正常？
A: 使用 `-v, --verbose` 参数可以查看详细的处理过程，确认哪些文件被包含或排除。