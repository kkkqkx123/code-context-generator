# 代码上下文生成器 (Code Context Generator)

一个智能的代码项目结构文档生成工具，支持CLI和TUI两种交互方式，能够扫描代码项目并生成结构化的文档输出。

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

## 安装

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

## 快速开始

### CLI 使用

#### 基本用法
```bash
# 扫描当前目录并输出JSON格式
./code-context-generator generate

# 扫描指定目录
./code-context-generator generate /path/to/project

# 输出为Markdown格式
./code-context-generator generate -f markdown -o project-structure.md
```

#### 高级用法
```bash
# 排除特定文件/目录
./code-context-generator generate -e "*.log" -e "node_modules" -e ".git"

# 包含隐藏文件，限制扫描深度
./code-context-generator generate -h -d 3

# 包含文件内容和哈希值
./code-context-generator generate -C -H

# 限制文件大小
./code-context-generator generate -s 1048576  # 1MB
```

#### 交互式选择
```bash
# 启动交互式文件选择器
./code-context-generator select

# 选择后输出为指定格式
./code-context-generator select -f xml -o selected-files.xml
```

#### 配置管理
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

### TUI 使用

```bash
# 启动TUI界面
./code-context-generator-tui
```

TUI界面提供：
- 可视化路径输入
- 交互式文件选择
- 实时配置编辑
- 进度显示
- 结果预览

## 📚 文档

我们提供了完整的文档体系，帮助你快速上手和深入了解本工具：

### 🎯 新用户
- [**快速入门指南**](docs/quickstart.md) - 5分钟快速上手 🚀
- [**使用文档**](docs/usage.md) - 完整的使用指南 📖
- [**配置详解**](docs/usage.md#配置文件详解) - 配置项详细说明 ⚙️

### 🚀 部署和运维
- [**部署文档**](docs/deployment.md) - 多种部署方式指南 📦
- [**系统服务**](docs/deployment.md#系统服务部署) - 配置为系统服务 🔧
- [**容器化部署**](docs/deployment.md#容器化部署) - Docker/Kubernetes部署 🐳

### 💻 开发贡献
- [**开发环境文档**](docs/development.md) - 开发环境搭建指南 🛠️
- [**开发流程**](docs/development.md#开发流程) - 贡献代码流程 📋
- [**API文档**](docs/development.md#api文档) - 代码API文档 📊

### 📖 文档导航
- [**文档中心**](docs/README.md) - 所有文档的索引和导航 📑
- [**常见问题**](docs/usage.md#常见问题) - 常见问题解答 ❓
- [**故障排除**](docs/usage.md#故障排除) - 问题排查指南 🔍

## 配置

配置文件支持 TOML、YAML、JSON 格式，默认配置文件示例：

```toml
[output]
format = "json"
encoding = "utf-8"
file_path = ""
pretty = true

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
buffer_size = 8192
cache_enabled = true
cache_size = 100

[logging]
level = "info"
file_path = ""
max_size = 10
max_backups = 3
max_age = 7

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

## 架构设计

### 模块结构
```
code-context-generator/
├── cmd/
│   ├── cli/          # CLI应用程序入口
│   └── tui/          # TUI应用程序入口
├── internal/
│   ├── config/       # 配置管理
│   ├── filesystem/   # 文件系统操作
│   ├── formatter/    # 格式转换
│   ├── selector/     # 文件选择器
│   ├── autocomplete/ # 自动补全
│   └── utils/        # 工具函数
├── pkg/
│   ├── types/        # 类型定义
│   └── constants/    # 常量定义
├── configs/          # 配置文件
├── docs/            # 文档
└── tests/           # 测试文件
```

### 核心组件

#### 1. 配置管理器 (Config Manager)
- 支持多格式配置文件（TOML、YAML、JSON）
- 环境变量覆盖
- 配置验证和默认值
- 热重载支持

#### 2. 文件系统遍历器 (File System Walker)
- 并发文件扫描
- 灵活的过滤机制
- 大文件流式处理
- 进度报告

#### 3. 格式化器 (Formatter)
- 多格式输出支持
- 自定义字段映射
- 模板系统
- 代码高亮

#### 4. 选择器 (Selector)
- 交互式文件选择
- 多选/单选模式
- 搜索和过滤
- 预览功能

#### 5. 自动补全器 (Autocompleter)
- 智能路径补全
- 上下文感知
- 缓存机制
- 模糊匹配

### 技术栈

#### 核心依赖
- **CLI框架**: [Cobra](https://github.com/spf13/cobra) - 现代化的CLI应用框架
- **TUI框架**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) - 优雅的终端用户界面
- **样式库**: [Lipgloss](https://github.com/charmbracelet/lipgloss) - 终端样式和布局

#### 配置和序列化
- **配置解析**: [TOML](https://github.com/BurntSushi/toml), [YAML](https://github.com/goccy/go-yaml)
- **JSON处理**: 标准库 `encoding/json`
- **XML处理**: 标准库 `encoding/xml`

#### 文件处理
- **文件监控**: [fsnotify](https://github.com/fsnotify/fsnotify) - 文件系统事件监控
- **路径处理**: 标准库 `path/filepath`
- **并发控制**: 标准库 `sync`, `context`

#### 日志和错误处理
- **日志库**: [logrus](https://github.com/sirupsen/logrus) - 结构化日志
- **错误处理**: 自定义错误类型和包装

## 性能优化

### 并发处理
- 使用 goroutine 池控制并发数量
- 工作队列模式处理文件扫描
- 上下文取消支持

### 内存管理
- 对象池复用减少GC压力
- 流式处理避免大内存占用
- 智能缓存策略

### I/O优化
- 批量文件操作
- 异步I/O模式
- 预读取和延迟写入

## 错误处理

### 错误类型
- **文件系统错误**: 权限、不存在、磁盘空间
- **配置错误**: 格式、验证、不兼容
- **网络错误**: 远程文件访问
- **内存错误**: 大文件处理

### 错误处理策略
- 优雅降级
- 重试机制
- 详细错误信息
- 恢复建议

## 测试

### 单元测试
```bash
go test ./internal/... -v
```

### 集成测试
```bash
go test ./tests/... -v
```

### 性能测试
```bash
go test -bench=. ./internal/filesystem
```

## 贡献

### 开发环境设置
```bash
git clone https://github.com/yourusername/code-context-generator.git
cd code-context-generator
go mod download
```

### 代码规范
- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加充分的注释和文档
- 编写单元测试

### 提交规范
- 使用清晰的提交信息
- 关联相关Issue
- 添加适当的标签

## 路线图

### 近期计划 (v1.1)
- [ ] 远程文件系统支持（FTP、SFTP）
- [ ] 插件系统
- [ ] 主题自定义
- [ ] 多语言支持

### 中期计划 (v1.2)
- [ ] Web界面
- [ ] API服务
- [ ] 数据库集成
- [ ] 云存储支持

### 长期计划 (v2.0)
- [ ] AI智能分析
- [ ] 代码质量检测
- [ ] 依赖关系图
- [ ] 实时同步

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 致谢

- [Cobra](https://github.com/spf13/cobra) - CLI框架
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI框架
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - 样式库
- 所有贡献者和支持者

## 联系方式

- **Issue**: [GitHub Issues](https://github.com/yourusername/code-context-generator/issues)
- **邮件**: your.email@example.com
- **文档**: [Wiki](https://github.com/yourusername/code-context-generator/wiki)

---

⭐ 如果这个项目对你有帮助，请给我们一个星标！