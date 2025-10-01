# 代码上下文生成器 - 开发环境使用文档

## 开发环境要求

- **Go**: 1.24+
- **Git**: 2.30+
- **操作系统**: Windows 10+/Linux/macOS

## 环境搭建

### 1. Go环境安装

```bash
# 验证安装
go version
```

### 2. 获取项目源码

```bash
# 克隆仓库
git clone https://github.com/yourusername/code-context-generator.git
cd code-context-generator
```

### 3. 安装依赖

```bash
# 下载依赖
go mod download

# 整理依赖
go mod tidy
```

### 4. 开发工具安装

```bash
# 安装代码质量工具
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 安装调试工具
go install github.com/go-delve/delve/cmd/dlv@latest
```

## 项目结构

```
code-context-generator/
├── cmd/                    # 应用程序入口
│   ├── cli/               # CLI应用入口
│   │   ├── main.go
│   │   └── config.yaml
│   └── tui/               # TUI应用入口
│   │   ├── main.go
│   │   └── models/
├── internal/              # 内部包
│   ├── autocomplete/      # 自动补全功能
│   │   ├── autocomplete.go
│   │   └── autocomplete_test.go
│   ├── config/           # 配置管理
│   │   ├── config.go
│   │   ├── config.yaml
│   │   └── config_test.go
│   ├── env/              # 环境变量处理
│   │   └── env.go
│   ├── filesystem/       # 文件系统操作
│   │   ├── filesystem.go
│   │   └── filesystem_test.go
│   ├── formatter/        # 格式转换
│   │   ├── formatter.go
│   │   └── formatter_test.go
│   ├── selector/         # 文件选择器
│   │   ├── file_selector.go
│   │   ├── file_utils.go
│   │   ├── pattern_matcher.go
│   │   ├── selector.go
│   │   └── selector_test.go
│   └── utils/            # 工具函数
│   │   ├── color.go
│   │   ├── encoding.go
│   │   ├── file.go
│   │   ├── path.go
│   │   ├── regex.go
│   │   ├── string.go
│   │   ├── time.go
│   │   ├── utils.go
│   │   ├── utils_test.go
│   │   └── validation.go
├── pkg/                   # 可复用的包
│   ├── constants/        # 常量定义
│   │   └── constants.go
│   └── types/            # 类型定义
│       └── types.go
├── configs/               # 配置文件模板
├── docs/                  # 文档
├── examples/              # 使用示例和配置示例
├── tests/                 # 测试文件
├── go.mod                # Go模块定义
├── go.sum                # 依赖校验
├── README.md             # 项目说明
├── LICENSE               # 许可证
└── .gitignore            # Git忽略规则
```

## 开发流程

### 1. 分支管理

```bash
# 查看分支
git branch -a

# 创建功能分支
git checkout -b feature/add-new-formatter

# 创建修复分支
git checkout -b fix/memory-leak

# 创建发布分支
git checkout -b release/v1.1.0
```

### 2. 开发规范

#### 代码风格
- 遵循官方Go代码规范
- 使用gofmt格式化代码
- 使用golint检查代码质量
- 遵循项目内部的命名约定

#### 提交规范
```
类型(范围): 简短描述

详细描述（可选）

Fixes #123
```

类型包括：
- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

#### 示例提交
```bash
git commit -m "feat(formatter): add YAML format support

- Implement YAMLFormatter with custom field mapping
- Add YAML format configuration options
- Update documentation for YAML support

Fixes #45"
```

### 3. 开发步骤

#### 步骤1：创建功能分支
```bash
# 更新主分支
git checkout main
git pull upstream main

# 创建功能分支
git checkout -b feature/improve-performance
```

#### 步骤2：编写代码
```bash
# 创建新文件
touch internal/performance/optimizer.go
touch internal/performance/optimizer_test.go

# 编写代码（示例）
package performance

import (
    "runtime"
    "sync"
)

type Optimizer struct {
    workers int
    pool    *sync.Pool
}

func NewOptimizer(workers int) *Optimizer {
    return &Optimizer{
        workers: workers,
        pool: &sync.Pool{
            New: func() interface{} {
                return make([]byte, 4096)
            },
        },
    }
}
```

#### 步骤3：运行测试
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/performance/

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

#### 步骤4：代码质量检查
```bash
# 格式化代码
go fmt ./...

# 运行golint
golint ./...

# 运行go vet
go vet ./...

# 运行golangci-lint
golangci-lint run

# 检查依赖安全性
go list -json -m all | nancy sleuth
```

#### 步骤5：构建应用
```bash
# 构建CLI版本
go build -o bin/code-context-generator cmd/cli/main.go

# 构建TUI版本
go build -o bin/code-context-generator-tui cmd/tui/main.go

# 构建所有版本
make build

# 交叉编译
GOOS=windows GOARCH=amd64 go build -o bin/code-context-generator.exe cmd/cli/main.go
GOOS=linux GOARCH=amd64 go build -o bin/code-context-generator-linux cmd/cli/main.go
GOOS=darwin GOARCH=amd64 go build -o bin/code-context-generator-darwin cmd/cli/main.go
```

#### 步骤6：运行应用
```bash
# 运行CLI版本
./bin/code-context-generator --help

# 运行TUI版本
./bin/code-context-generator-tui

# 使用示例
./bin/code-context-generator generate . -f json -o output.json
```

#### 步骤7：提交代码
```bash
# 添加修改的文件
git add .

# 提交修改
git commit -m "feat(performance): add memory pool for better performance

- Implement sync.Pool for buffer reuse
- Add configurable worker pool size
- Improve memory allocation efficiency
- Add benchmarks for performance testing

Fixes #67"

# 推送到远程仓库
git push origin feature/improve-performance
```

## 测试指南

### 1. 单元测试

#### 创建测试文件
```go
// internal/formatter/formatter_test.go
package formatter

import (
    "testing"
    "code-context-generator/pkg/types"
)

func TestJSONFormatter_Format(t *testing.T) {
    formatter := NewJSONFormatter(nil)
    
    data := types.ContextData{
        Files: []types.FileInfo{
            {
                Path:     "test.go",
                Size:     1024,
                Modified: "2024-01-01T00:00:00Z",
            },
        },
        FileCount: 1,
        TotalSize: 1024,
    }
    
    result, err := formatter.Format(data)
    if err != nil {
        t.Fatalf("Format failed: %v", err)
    }
    
    if result == "" {
        t.Error("Expected non-empty result")
    }
    
    // 验证JSON格式
    if !strings.HasPrefix(result, "{") {
        t.Error("Expected JSON object")
    }
}
```

#### 运行单元测试
```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test ./internal/formatter/

# 运行测试并显示详细信息
go test -v ./...

# 运行测试并生成覆盖率报告
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 2. 集成测试

#### 创建集成测试
```go
// tests/integration_test.go
package tests

import (
    "os"
    "path/filepath"
    "testing"
    "code-context-generator/internal/filesystem"
)

func TestIntegration_FileSystemWalker(t *testing.T) {
    // 创建测试目录结构
    testDir := t.TempDir()
    
    // 创建测试文件
    files := []string{
        "file1.go",
        "file2.txt",
        "subdir/file3.json",
    }
    
    for _, file := range files {
        path := filepath.Join(testDir, file)
        os.MkdirAll(filepath.Dir(path), 0755)
        os.WriteFile(path, []byte("test content"), 0644)
    }
    
    // 测试文件系统遍历器
    walker := filesystem.NewFileSystemWalker(types.WalkOptions{
        MaxDepth: 3,
        ShowHidden: false,
    })
    
    result, err := walker.Walk(testDir, nil)
    if err != nil {
        t.Fatalf("Walk failed: %v", err)
    }
    
    // 验证结果
    if len(result.Files) != len(files) {
        t.Errorf("Expected %d files, got %d", len(files), len(result.Files))
    }
}
```


## 故障排除

### 常见问题

#### 依赖问题
```bash
# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download

# 更新依赖
go get -u ./...
go mod tidy
```

#### 构建问题
```bash
# 清理构建缓存
go clean -cache

# 强制重新构建
go build -a ./...

# 检查构建约束
//go:build linux && amd64
```

#### 测试问题
```bash
# 运行测试并显示详细输出
go test -v ./...

# 运行特定测试
go test -run TestJSONFormatter ./internal/formatter/

# 跳过某些测试
go test -short ./...
```
