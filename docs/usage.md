# 代码上下文生成器 - 使用文档

## 概述

代码上下文生成器是一个智能的代码项目结构文档生成工具，支持CLI和TUI两种交互方式，能够扫描代码项目并生成结构化的文档输出。

## 功能特性

### 🎯 核心功能
- **多格式输出**: 支持 JSON、XML、TOML、Markdown 格式
- **智能文件选择**: 交互式文件/目录选择界面
- **自动补全**: 文件路径智能补全功能
- **配置管理**: 灵活的配置系统，支持环境变量覆盖

### 🚀 高级特性
- **并发处理**: 基于 goroutine 池的高性能文件扫描
- **大文件处理**: 流式读取，支持大文件处理
- **模式匹配**: 支持 glob 模式和正则表达式过滤
- **缓存机制**: 智能缓存提升重复扫描性能
- **跨平台**: 支持 Windows、Linux、macOS

### 🎨 用户界面
- **CLI 模式**: 功能丰富的命令行界面（基于 Cobra）
- **TUI 模式**: 现代化的终端用户界面（基于 Bubble Tea）
- **进度显示**: 实时进度条和状态信息
- **主题支持**: 可定制的界面主题

## 安装方法

### 前置要求
- Go 1.24 或更高版本
- Git（可选，用于版本控制集成）

### 从源码安装
```bash
git clone https://github.com/yourusername/code-context-generator.git
cd code-context-generator
go build -o code-context-generator cmd/cli/main.go
```

### 构建TUI版本
```bash
go build -o code-context-generator-tui cmd/tui/main.go
```

## CLI 使用指南

### 基本用法

#### 扫描当前目录
```bash
# 扫描当前目录并输出JSON格式
./code-context-generator generate
```

#### 扫描指定目录
```bash
# 扫描指定目录
./code-context-generator generate /path/to/project
```

#### 指定输出格式
```bash
# 输出为Markdown格式
./code-context-generator generate -f markdown -o project-structure.md

# 输出为XML格式
./code-context-generator generate -f xml -o project-structure.xml

# 输出为TOML格式
./code-context-generator generate -f toml -o project-structure.toml
```

### 高级用法

#### 文件过滤
```bash
# 排除特定文件/目录
./code-context-generator generate -e "*.log" -e "node_modules" -e ".git"

# 只包含特定类型的文件
./code-context-generator generate -i "*.go" -i "*.md" -i "*.json"
```

#### 扫描选项
```bash
# 包含隐藏文件，限制扫描深度
./code-context-generator generate -h -d 3

# 包含文件内容和哈希值
./code-context-generator generate -C -H

# 限制文件大小（字节）
./code-context-generator generate -s 1048576  # 1MB
```

#### 递归控制
```bash
# 禁用递归扫描（只扫描当前目录）
./code-context-generator generate --no-recursive

# 指定最大递归深度
./code-context-generator generate -d 5
```

### 交互式选择

#### 启动文件选择器
```bash
# 启动交互式文件选择器
./code-context-generator select

# 选择后输出为指定格式
./code-context-generator select -f xml -o selected-files.xml
```

### 配置管理

#### 初始化配置
```bash
# 初始化配置文件
./code-context-generator config init

# 显示当前配置
./code-context-generator config show
```

#### 自动补全

```bash
# 获取文件路径补全建议
./code-context-generator autocomplete /path/to/

# 获取目录补全建议
./code-context-generator autocomplete -t dir /path/to/
```

### 命令参数详解

#### generate 命令
```bash
./code-context-generator generate [路径] [flags]

Flags:
  -o, --output string         输出文件路径
  -f, --format string        输出格式 (json|xml|toml|markdown) (默认 "json")
  -e, --exclude strings      排除模式（可多次使用）
  -i, --include strings      包含模式（可多次使用）
  -r, --recursive            递归扫描（默认true）
  -d, --max-depth int        最大扫描深度
  -h, --hidden               包含隐藏文件
  -s, --max-size int         最大文件大小（字节）
  -C, --content              包含文件内容
  -H, --hash                 包含文件哈希值
  -v, --verbose              详细输出
  -h, --help                 帮助信息
```

#### select 命令
```bash
./code-context-generator select [flags]

Flags:
  -o, --output string         输出文件路径
  -f, --format string        输出格式 (json|xml|toml|markdown) (默认 "json")
  -m, --multi                允许多选
  -h, --help                 帮助信息
```

#### config 命令
```bash
./code-context-generator config [command]

Available Commands:
  init    初始化配置文件
  show    显示当前配置
  edit    编辑配置文件
  validate 验证配置文件
```

## TUI 使用指南

### 启动TUI界面
```bash
# 启动TUI界面
./code-context-generator-tui
```

### TUI界面功能

#### 主界面
- **路径输入**: 可视化路径输入框
- **格式选择**: 下拉选择输出格式
- **选项配置**: 复选框配置扫描选项
- **快速操作**: 常用功能快捷键

#### 文件选择器
- **目录树展示**: 层级化目录结构
- **多选支持**: Ctrl+Space 多选文件
- **键盘导航**: 方向键导航，Enter确认
- **搜索模式**: / 键进入搜索模式
- **实时过滤**: 动态过滤文件列表

#### 配置编辑器
- **实时预览**: 配置更改实时生效
- **格式验证**: 输入验证和错误提示
- **模板支持**: 自定义输出模板
- **主题切换**: 多种界面主题

#### 进度显示
- **实时进度条**: 扫描进度可视化
- **状态信息**: 当前操作状态
- **速度显示**: 处理速度统计
- **剩余时间**: 预估完成时间

### TUI快捷键

#### 全局快捷键
```
Ctrl+C: 退出程序
Ctrl+R: 重新扫描
Ctrl+S: 保存配置
Ctrl+H: 显示帮助
Tab: 切换面板
```

#### 文件选择器快捷键
```
↑/↓: 上下移动
←/→: 展开/收起目录
Space: 选择/取消选择
Ctrl+A: 全选
Ctrl+N: 取消全选
/: 进入搜索模式
Esc: 退出搜索模式
```

#### 配置编辑器快捷键
```
Ctrl+Z: 撤销
Ctrl+Y: 重做
Ctrl+S: 保存配置
Ctrl+L: 加载配置
```

## 配置文件详解

### 配置文件格式

支持三种格式：TOML、YAML、JSON，默认使用 TOML 格式。

#### TOML 配置示例
```toml
[output]
format = "json"
encoding = "utf-8"
file_path = ""

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
include_patterns = []
include_content = false
include_hash = false

[ui]
theme = "default"
show_progress = true
show_size = true
show_date = true
show_preview = true

[performance]
max_workers = 4
buffer_size = 1024
cache_enabled = true
cache_size = 100

[logging]
level = "info"
file_path = ""
max_size = 10
max_backups = 3
max_age = 7
```

#### YAML 配置示例
```yaml
output:
  format: json
  encoding: utf-8
  file_path: ""

file_processing:
  include_hidden: false
  max_file_size: 10485760
  max_depth: 0
  exclude_patterns:
    - "*.exe"
    - "*.dll"
    - "node_modules"
    - ".git"
  include_patterns: []
  include_content: false
  include_hash: false

ui:
  theme: default
  show_progress: true
  show_size: true
  show_date: true
  show_preview: true

performance:
  max_workers: 4
  buffer_size: 1024
  cache_enabled: true
  cache_size: 100

logging:
  level: info
  file_path: ""
  max_size: 10
  max_backups: 3
  max_age: 7
```

### 配置项说明

#### 输出配置 (output)
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| format | string | "json" | 输出格式：json、xml、toml、markdown |
| encoding | string | "utf-8" | 文件编码：utf-8、gbk、gb2312 |
| file_path | string | "" | 输出文件路径，空字符串表示标准输出 |

#### 文件处理配置 (file_processing)
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| include_hidden | bool | false | 是否包含隐藏文件 |
| max_file_size | int | 10485760 | 最大文件大小（字节） |
| max_depth | int | 0 | 最大扫描深度，0表示无限制 |
| exclude_patterns | []string | [] | 排除模式列表 |
| include_patterns | []string | [] | 包含模式列表 |
| include_content | bool | false | 是否包含文件内容 |
| include_hash | bool | false | 是否包含文件哈希值 |

#### 界面配置 (ui)
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| theme | string | "default" | 界面主题：default、dark、light |
| show_progress | bool | true | 是否显示进度条 |
| show_size | bool | true | 是否显示文件大小 |
| show_date | bool | true | 是否显示修改日期 |
| show_preview | bool | true | 是否显示预览 |

#### 性能配置 (performance)
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| max_workers | int | 4 | 最大工作线程数 |
| buffer_size | int | 1024 | 缓冲区大小 |
| cache_enabled | bool | true | 是否启用缓存 |
| cache_size | int | 100 | 缓存大小 |

#### 日志配置 (logging)
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| level | string | "info" | 日志级别：debug、info、warn、error |
| file_path | string | "" | 日志文件路径，空字符串表示控制台输出 |
| max_size | int | 10 | 日志文件最大大小（MB） |
| max_backups | int | 3 | 最大备份文件数 |
| max_age | int | 7 | 日志文件最大保存天数 |

## 使用示例

### 示例1：基础项目扫描
```bash
# 扫描Go项目并生成JSON格式的项目结构
./code-context-generator generate ~/projects/my-go-app -f json -o my-app-structure.json
```

### 示例2：前端项目文档生成
```bash
# 扫描React项目，排除node_modules和构建文件
./code-context-generator generate ~/projects/react-app \
  -e "node_modules" -e "build" -e "dist" -e "*.log" \
  -f markdown -o react-app-docs.md
```

### 示例3：代码审查准备
```bash
# 扫描代码并包含内容，用于代码审查
./code-context-generator generate ~/projects/code-review \
  -C -H -f xml -o code-review-context.xml \
  -i "*.go" -i "*.js" -i "*.py" -s 5242880  # 5MB限制
```

### 示例4：配置文件使用
```bash
# 使用自定义配置文件
./code-context-generator generate ~/projects/my-app --config my-config.toml

# 使用环境变量覆盖配置
export CODE_CONTEXT_FORMAT=xml
export CODE_CONTEXT_MAX_SIZE=10485760
./code-context-generator generate ~/projects/my-app
```

### 示例5：交互式选择
```bash
# 使用TUI选择特定文件
./code-context-generator-tui

# 使用CLI选择器
./code-context-generator select -m -f json -o selected-files.json
```

## 常见问题

### Q: 如何处理大文件？
A: 使用 `-s` 参数限制文件大小，例如 `-s 10485760` 限制为10MB。对于超大文件，建议使用流式处理模式。

### Q: 如何排除特定目录？
A: 使用 `-e` 参数指定排除模式，支持glob模式：`-e "node_modules" -e ".git" -e "*.log"`

### Q: 如何包含隐藏文件？
A: 使用 `-h` 或 `--hidden` 参数包含隐藏文件。

### Q: 如何自定义输出格式？
A: 通过配置文件中的模板系统自定义输出格式，支持字段映射和结构自定义。

### Q: 性能优化建议？
A: 
1. 合理设置 `max_workers` 参数
2. 启用缓存机制
3. 使用适当的缓冲区大小
4. 限制扫描深度和文件大小
5. 排除不必要的目录

### Q: 如何处理中文路径？
A: 工具原生支持UTF-8编码，对于中文路径和文件名无需特殊配置。

### Q: 如何调试问题？
A: 
1. 使用 `-v` 参数启用详细输出
2. 设置日志级别为 `debug`
3. 检查配置文件语法
4. 验证文件权限

## 最佳实践

### 1. 项目文档化
```bash
# 为项目创建完整的文档结构
./code-context-generator generate . \
  -f markdown \
  -o PROJECT_STRUCTURE.md \
  -e "node_modules" -e ".git" -e "vendor" \
  -C -H
```

### 2. 代码审查准备
```bash
# 生成包含内容的代码上下文
./code-context-generator generate src/ \
  -f xml \
  -o code-review.xml \
  -i "*.go" -i "*.md" \
  -C -H -s 5242880
```

### 3. 持续集成
```bash
# 在CI中生成项目结构报告
./code-context-generator generate . \
  -f json \
  -o project-report.json \
  --no-recursive \
  -e "*.tmp" -e "*.log"
```

### 4. 配置文件模板
创建项目专用的配置文件模板，包含常用的排除模式和格式设置。

## 故障排除

### 常见错误

#### 权限错误
```
Error: 扫描失败: 权限被拒绝
```
解决方案：检查文件和目录的读取权限，使用管理员权限运行或修改文件权限。

#### 内存不足
```
Error: 内存不足
```
解决方案：减小 `max_file_size` 和 `buffer_size`，降低 `max_workers` 数量。

#### 配置文件错误
```
Error: 配置文件解析失败
```
解决方案：验证配置文件语法，检查字段名称和类型。

### 性能问题

#### 扫描速度慢
- 启用缓存机制
- 增加工作线程数
- 排除大文件和不必要的目录
- 使用适当的缓冲区大小

#### 内存使用过高
- 减小缓冲区大小
- 限制文件大小
- 降低并发线程数
- 及时清理缓存

### 获取帮助

#### 查看帮助信息
```bash
./code-context-generator --help
./code-context-generator generate --help
./code-context-generator select --help
```

#### 查看版本信息
```bash
./code-context-generator --version
```

#### 获取详细输出
```bash
./code-context-generator generate -v [其他参数]
```

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 支持CLI和TUI界面
- 支持JSON、XML、TOML、Markdown格式
- 基础文件过滤功能
- 配置管理系统

## 联系方式

如有问题或建议，请通过以下方式联系：
- 项目Issues: [GitHub Issues](https://github.com/yourusername/code-context-generator/issues)
- 邮箱: your.email@example.com