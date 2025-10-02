# 代码上下文生成器 - 详细设计文档

## 项目概述

代码上下文生成器是一个高性能的CLI工具，用于通过终端选择文件/文件夹，并将选中的内容打包为结构化的XML/JSON/TOML/Markdown文件，方便用户快速构建代码上下文和提示词。

## 架构设计

### 总体架构

```
code-context-generator/
├── cmd/                    # 命令入口
│   ├── cli/               # CLI命令实现

├── internal/              # 内部核心模块
│   ├── config/           # 配置管理
│   ├── filesystem/       # 文件系统操作
│   ├── formatter/        # 格式转换
│   ├── selector/         # 文件选择器

│   └── utils/            # 工具函数
├── pkg/                   # 可复用包
│   ├── types/            # 公共类型定义
│   └── constants/        # 常量定义
├── configs/               # 配置模板
├── docs/                  # 文档
└── tests/                 # 测试文件
```

### 核心模块设计

#### 1. 配置管理模块 (internal/config/)

**职责：**
- 多格式配置文件解析（YAML/JSON/TOML）
- 环境变量覆盖
- 配置验证和默认值处理
- 配置热重载

**接口定义：**
```go
type ConfigManager interface {
    Load(configPath string) error
    Get() *Config
    Validate() error
    Reload() error
    Save(configPath string, format string) error
    GetEnvOverrides() map[string]string
}
```

**现有组件重构：**
- 将现有`config_manager.go`重构为模块化设计
- 保持向后兼容性
- 增强错误处理和验证

#### 2. 文件系统模块 (internal/filesystem/)

**职责：**
- 安全的递归目录遍历
- .gitignore规则解析和应用
- 文件过滤和大小检查
- 符号链接处理
- 中文路径支持
- 并发文件读取

**接口定义：**
```go
type FileSystem interface {
    Walk(root string, options WalkOptions) (<-chan FileInfo, error)
    ReadFile(path string) (string, error)
    GetFileInfo(path string) (FileInfo, error)
    ParseGitignore(path string) ([]string, error)
    IsHidden(path string) bool
    IsSymlink(path string) bool
}

type WalkOptions struct {
    MaxDepth        int
    FollowSymlinks  bool
    ExcludePatterns []string
    IncludePatterns []string
    MaxFileSize     int64
}
```

**关键特性：**
- 使用goroutine池进行并发遍历
- 流式处理避免内存溢出
- 完善的错误处理机制
- 跨平台路径处理

#### 3. 格式转换模块 (internal/formatter/)

**职责：**
- XML/JSON/TOML/Markdown格式生成
- 模板引擎支持
- 并发格式转换
- 内存优化

**接口定义：**
```go
type Formatter interface {
    Format(data ContextData, format string) (string, error)
    GetSupportedFormats() []string
    ValidateFormat(format string) error
    SetTemplate(format string, template string) error
}

type ContextData struct {
    Files   []FileInfo
    Folders []FolderInfo
    Metadata map[string]interface{}
}
```

**格式支持：**
- XML: 结构化层次，支持自定义标签
- JSON: 标准JSON格式，支持缩进和排序
- TOML: 易读配置格式
- Markdown: 带代码高亮的文档格式

#### 4. 文件选择器模块 (internal/selector/)

**职责：**
- 交互式文件选择
- 多选支持
- 键盘导航
- 实时搜索过滤

**特性：**
- 支持方向键导航
- 空格键选择/取消选择
- 回车键进入目录
- 退格键返回上级目录
- /键进入搜索模式



#### 6. CLI模块 (cmd/cli/)

**命令结构：**
```bash
code-context-generator [command] [flags]

Commands:
  generate    生成代码上下文文件
  config      管理配置文件
  validate    验证配置文件
  version     显示版本信息

Flags:
  --format string        输出格式 (xml|json|toml|markdown)
  --output string        输出文件路径
  --config string        配置文件路径
  --exclude strings      排除模式
  --include strings      包含模式
  --max-depth int        最大遍历深度
  --follow-symlinks      跟随符号链接
  --output-dir string    输出目录
  --filename-template    文件名模板
```


│ │ > src/           │ │ │ 格式: XML                        │ │
│ │   main.go        │ │ │ 输出目录: ./output              │ │
│ │   utils.go       │ │ │ 排除: *.tmp, *.log             │ │
│ │   config/        │ │ │ 最大深度: 3                     │ │
│ │   tests/         │ │ │ 符号链接: 否                    │ │
│ │                  │ │ │                                 │ │
│ │ 空格:选择 回车:进入 │ │ │                                 │ │
│ │ /:搜索 q:退出    │ │ │                                 │ │
│ └──────────────────┘ │ └─────────────────────────────────┘ │
├──────────────────────┴────────────────────────────────────┤
│ 已选择: 3个文件, 1个目录  状态: 就绪                    │
└─────────────────────────────────────────────────────────────┘
```

## 技术选型

### 核心依赖
- **TUI框架:** bubbletea + lipgloss
- **CLI框架:** cobra
- **配置解析:** 
  - YAML: github.com/goccy/go-yaml
  - JSON: encoding/json
  - TOML: github.com/BurntSushi/toml
- **文件监控:** fsnotify (用于热重载)
- **日志:** logrus

### 开发工具
- **构建:** go build
- **测试:** go test + testify
- **代码质量:** golangci-lint
- **文档:** godoc

## 性能优化策略

### 1. 并发处理
```go
type WorkerPool struct {
    workers    int
    jobQueue   chan Job
    resultQueue chan Result
    wg         sync.WaitGroup
}
```

### 2. 内存优化
- 对象池复用
- 流式文件读取
- 及时内存清理
- 大文件分块处理

### 3. 缓存机制
- 文件信息缓存
- 配置缓存
- 模板编译缓存

## 错误处理设计

### 错误分类
```go
type ErrorType int

const (
    ErrConfig ErrorType = iota
    ErrFileSystem
    ErrFormat
    ErrValidation
    ErrPermission
)

type AppError struct {
    Type    ErrorType
    Message string
    Cause   error
    Context map[string]interface{}
}
```

### 错误处理策略
- 用户友好的错误消息
- 详细的错误上下文
- 恢复机制
- 日志记录

## 测试策略

### 单元测试
- 每个模块独立测试
- 边界条件测试
- 错误场景测试
- 性能基准测试

### 集成测试
- 端到端功能测试
- 跨平台兼容性测试
- 大文件处理测试
- 并发安全性测试

### 测试覆盖率目标
- 核心模块: >90%
- 业务逻辑: >80%
- 整体: >75%

## 部署和发布

### 构建配置
```makefile
# Makefile
BINARY_NAME=code-context-generator
VERSION=$(shell git describe --tags --always)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

build:
    go build -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -o $(BINARY_NAME) main.go

build-all:
    GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe
    GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64
    GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64
```

### 发布策略
- GitHub Releases
- 多平台二进制文件
- 安装脚本
- Docker镜像（可选）

## 监控和运维

### 指标收集
- 处理文件数量
- 处理时间
- 内存使用
- 错误率

### 日志设计
```go
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Fatal(msg string, fields ...Field)
}
```

## 安全考虑

### 输入验证
- 路径遍历防护
- 文件大小限制
- 符号链接验证
- 编码安全检查

### 权限管理
- 文件读取权限检查
- 输出目录权限验证
- 安全配置验证

## 扩展性设计

### 插件架构（未来）
```go
type Plugin interface {
    Name() string
    Version() string
    Init(config map[string]interface{}) error
    Process(data ContextData) (ContextData, error)
}
```

### 新格式支持
- 格式化器接口设计
- 模板系统扩展
- 配置架构兼容

## 兼容性保证

### 向后兼容
- 配置文件格式兼容
- API接口稳定
- 命令行参数兼容

### 平台兼容
- Windows 10+/Linux/macOS
- PowerShell/Bash/Zsh
- UTF-8编码支持

这个设计文档为项目提供了全面的架构指导，确保项目的高性能、可维护性和扩展性。