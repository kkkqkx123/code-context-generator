# 代码上下文生成器 - 开发环境使用文档

## 概述

本文档为开发者提供完整的开发环境搭建指南，包括环境配置、代码结构、开发流程、测试方法和发布流程。

## 开发环境要求

### 必需工具
- **Go**: 1.24 或更高版本
- **Git**: 2.30 或更高版本
- **Make**: 可选，用于构建自动化
- **Docker**: 可选，用于容器化开发

### 推荐工具
- **IDE**: Visual Studio Code、GoLand、Vim/Neovim
- **编辑器插件**: Go扩展、语法高亮、代码格式化
- **调试工具**: Delve (dlv)
- **性能分析**: pprof、benchcmp
- **代码质量**: golangci-lint、go vet、go fmt

### 系统要求
- **操作系统**: Windows 10+/Linux/macOS
- **内存**: 4GB RAM（推荐8GB）
- **存储**: 2GB 可用空间
- **CPU**: 多核处理器

## 环境搭建

### 1. Go环境安装

#### Windows
```powershell
# 使用Scoop安装
scoop install go

# 或者从官网下载安装包
# https://golang.org/dl/
```

#### Linux
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# 或者安装最新版本
wget https://golang.org/dl/go1.24.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### macOS
```bash
# 使用Homebrew安装
brew install go

# 验证安装
go version
```

### 2. 开发工具配置

#### Git配置
```bash
# 配置Git用户信息
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"

# 配置Git别名
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit
git config --global alias.st status
```

#### Go环境变量
```bash
# 设置GOPATH和GOPROXY
export GOPATH=$HOME/go
export GOPROXY=https://goproxy.io,direct
export GO111MODULE=on
export GOSUMDB=sum.golang.org

# 添加到shell配置文件（~/.bashrc 或 ~/.zshrc）
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export GOPROXY=https://goproxy.io,direct' >> ~/.bashrc
echo 'export GO111MODULE=on' >> ~/.bashrc
```

### 3. 获取项目源码

```bash
# 克隆仓库
git clone https://github.com/yourusername/code-context-generator.git
cd code-context-generator

# 或者fork到自己的仓库后克隆
git clone https://github.com/yourusername/code-context-generator.git
cd code-context-generator

# 添加上游仓库
git remote add upstream https://github.com/original/code-context-generator.git
```

### 4. 安装依赖

```bash
# 下载项目依赖
go mod download

# 验证依赖
go mod verify

# 整理依赖
go mod tidy

# 查看依赖关系
go mod graph
```

### 5. 开发工具安装

```bash
# 安装代码质量工具
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/segmentio/golines@latest

# 安装调试工具
go install github.com/go-delve/delve/cmd/dlv@latest

# 安装测试工具
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install github.com/onsi/gomega/...@latest

# 安装文档工具
go install golang.org/x/tools/cmd/godoc@latest
```

## 项目结构

```
code-context-generator/
├── cmd/                    # 应用程序入口
│   ├── cli/               # CLI应用入口
│   │   └── main.go
│   └── tui/               # TUI应用入口
│       ├── main.go
│       └── models.go
├── internal/              # 内部包
│   ├── autocomplete/      # 自动补全功能
│   │   └── autocomplete.go
│   ├── config/           # 配置管理
│   │   └── config.go
│   ├── filesystem/       # 文件系统操作
│   │   └── filesystem.go
│   ├── formatter/        # 格式转换
│   │   └── formatter.go
│   ├── selector/         # 文件选择器
│   │   └── selector.go
│   └── utils/            # 工具函数
│       └── utils.go
├── pkg/                   # 可复用的包
│   ├── constants/        # 常量定义
│   │   └── constants.go
│   └── types/            # 类型定义
│       └── types.go
├── configs/               # 配置文件模板
├── docs/                  # 文档
├── tests/                 # 测试文件
├── scripts/               # 构建和部署脚本
├── Makefile              # 构建自动化
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

### 3. 基准测试

#### 创建基准测试
```go
// internal/formatter/formatter_bench_test.go
package formatter

import (
    "testing"
    "code-context-generator/pkg/types"
)

func BenchmarkJSONFormatter_Format(b *testing.B) {
    formatter := NewJSONFormatter(nil)
    
    // 创建测试数据
    data := createLargeTestData()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := formatter.Format(data)
        if err != nil {
            b.Fatalf("Format failed: %v", err)
        }
    }
}

func createLargeTestData() types.ContextData {
    data := types.ContextData{
        Files: make([]types.FileInfo, 1000),
        Folders: make([]types.FolderInfo, 100),
    }
    
    for i := 0; i < 1000; i++ {
        data.Files[i] = types.FileInfo{
            Path:     fmt.Sprintf("file%d.go", i),
            Size:     int64(i * 1024),
            Modified: "2024-01-01T00:00:00Z",
            Content:  fmt.Sprintf("content of file %d", i),
        }
    }
    
    data.FileCount = 1000
    data.FolderCount = 100
    data.TotalSize = 1024 * 1024 * 10 // 10MB
    
    return data
}
```

#### 运行基准测试
```bash
# 运行基准测试
go test -bench=. ./...

# 运行特定基准测试
go test -bench=BenchmarkJSONFormatter ./internal/formatter/

# 运行基准测试并生成内存分析
go test -bench=. -benchmem -memprofile=mem.prof ./...
go tool pprof mem.prof

# 运行基准测试并生成CPU分析
go test -bench=. -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

### 4. 模糊测试

#### 创建模糊测试
```go
// internal/formatter/formatter_fuzz_test.go
package formatter

import (
    "testing"
    "code-context-generator/pkg/types"
)

func FuzzJSONFormatter_Format(f *testing.F) {
    formatter := NewJSONFormatter(nil)
    
    // 添加种子语料库
    f.Add("test.go", int64(1024), "2024-01-01T00:00:00Z", "test content")
    f.Add("", int64(0), "", "")
    f.Add("../../etc/passwd", int64(9999999999), "invalid-date", string([]byte{0, 1, 2, 3, 4}))
    
    f.Fuzz(func(t *testing.T, path string, size int64, modified string, content string) {
        data := types.ContextData{
            Files: []types.FileInfo{
                {
                    Path:     path,
                    Size:     size,
                    Modified: modified,
                    Content:  content,
                },
            },
            FileCount: 1,
            TotalSize: size,
        }
        
        result, err := formatter.Format(data)
        if err != nil {
            // 预期的错误情况
            return
        }
        
        // 验证结果不为空
        if result == "" {
            t.Error("Expected non-empty result")
        }
    })
}
```

#### 运行模糊测试
```bash
# 运行模糊测试
go test -fuzz=FuzzJSONFormatter ./internal/formatter/

# 运行模糊测试指定时间
go test -fuzz=FuzzJSONFormatter -fuzztime=10s ./internal/formatter/

# 使用特定的语料库
go test -fuzz=FuzzJSONFormatter -fuzzdir=testdata/fuzz ./internal/formatter/
```

## 调试指南

### 1. 使用Delve调试器

#### 安装Delve
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

#### 启动调试会话
```bash
# 调试CLI应用
dlv debug cmd/cli/main.go

# 调试TUI应用
dlv debug cmd/tui/main.go

# 调试特定测试
dlv test ./internal/formatter/
```

#### 常用调试命令
```gdb
# 设置断点
(dlv) break main.main
(dlv) break internal/formatter/formatter.go:45

# 运行程序
(dlv) continue

# 单步执行
(dlv) next
(dlv) step

# 查看变量
(dlv) print variableName
(dlv) locals

# 查看调用栈
(dlv) stack

# 继续执行
(dlv) continue
```

### 2. 使用日志调试

#### 添加日志
```go
import (
    "log"
    "os"
)

var (
    debugLog = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile)
    infoLog  = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
    errorLog = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
)

func (f *JSONFormatter) Format(data types.ContextData) (string, error) {
    debugLog.Printf("Formatting data with %d files", data.FileCount)
    
    output, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        errorLog.Printf("JSON formatting failed: %v", err)
        return "", fmt.Errorf("JSON格式化失败: %w", err)
    }
    
    infoLog.Printf("Successfully formatted %d bytes", len(output))
    return string(output), nil
}
```

### 3. 性能分析

#### CPU分析
```go
// 在代码中添加分析
import (
    "os"
    "runtime/pprof"
)

func main() {
    // 创建CPU分析文件
    cpuProfile, err := os.Create("cpu.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer cpuProfile.Close()
    
    // 开始CPU分析
    if err := pprof.StartCPUProfile(cpuProfile); err != nil {
        log.Fatal(err)
    }
    defer pprof.StopCPUProfile()
    
    // 应用程序逻辑
    runApplication()
}
```

#### 内存分析
```go
// 内存分析
memProfile, err := os.Create("mem.prof")
if err != nil {
    log.Fatal(err)
}
defer memProfile.Close()

// 获取内存分析数据
runtime.GC()
if err := pprof.WriteHeapProfile(memProfile); err != nil {
    log.Fatal(err)
}
```

#### 分析工具
```bash
# 查看CPU分析
go tool pprof cpu.prof
(pprof) top
(pprof) list functionName
(pprof) web

# 查看内存分析
go tool pprof mem.prof
(pprof) top
(pprof) list functionName
(pprof) web
```

## 代码质量

### 1. 代码格式化

```bash
# 格式化所有代码
go fmt ./...

# 使用goimports（自动管理导入）
goimports -w .

# 使用golines（格式化长行）
golines -w .
```

### 2. 静态分析

```bash
# 运行go vet
go vet ./...

# 运行golangci-lint
golangci-lint run

# 运行staticcheck
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

### 3. 安全检查

```bash
# 检查依赖安全性
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
gosec ./...

# 检查依赖漏洞
go install github.com/sonatypecommunity/nancy@latest
go list -json -m all | nancy sleuth
```

## 文档生成

### 1. Go文档

```bash
# 生成文档
go doc

# 查看特定包的文档
go doc code-context-generator/internal/formatter

# 查看特定函数的文档
go doc JSONFormatter.Format

# 启动文档服务器
godoc -http=:6060
# 访问 http://localhost:6060
```

### 2. 代码注释规范

```go
// JSONFormatter JSON格式转换器
type JSONFormatter struct {
    BaseFormatter
}

// NewJSONFormatter 创建JSON格式转换器
// 
// 参数:
//   - config: 格式配置，可为nil
//
// 返回:
//   - Formatter: JSON格式转换器实例
func NewJSONFormatter(config *types.FormatConfig) Formatter {
    return &JSONFormatter{
        BaseFormatter: BaseFormatter{
            name:        "JSON",
            description: "JavaScript Object Notation format",
            config:      config,
        },
    }
}

// Format 格式化上下文数据
// 
// 该方法将ContextData格式化为JSON字符串，支持自定义结构和字段映射。
// 如果配置中指定了自定义结构，将使用自定义结构进行格式化。
//
// 参数:
//   - data: 要格式化的上下文数据
//
// 返回:
//   - string: 格式化的JSON字符串
//   - error: 格式化过程中的错误
func (f *JSONFormatter) Format(data types.ContextData) (string, error) {
    // 实现代码
}
```

## 发布流程

### 1. 版本管理

#### 语义化版本
- **主版本号(MAJOR)**: 不兼容的API修改
- **次版本号(MINOR)**: 向下兼容的功能性新增
- **修订号(PATCH)**: 向下兼容的问题修正

#### 创建版本标签
```bash
# 更新版本号（在代码中）
# 通常在main.go或version.go中

# 提交版本更新
git add .
git commit -m "chore(version): bump version to v1.1.0"

# 创建标签
git tag -a v1.1.0 -m "Release version 1.1.0"

# 推送标签
git push origin v1.1.0
```

### 2. 构建发布版本

#### 创建构建脚本
```bash
#!/bin/bash
# scripts/build-release.sh

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

echo "Building version $VERSION..."

# 清理之前的构建
rm -rf dist/
mkdir -p dist/

# 构建不同平台的二进制文件
platforms=("linux/amd64" "darwin/amd64" "windows/amd64")

for platform in "${platforms[@]}"; do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    
    output_name="code-context-generator-$GOOS-$GOARCH"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "-X main.version=$VERSION -X main.buildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" \
        -o "dist/$output_name" \
        cmd/cli/main.go
    
    # 构建TUI版本
    tui_output_name="code-context-generator-tui-$GOOS-$GOARCH"
    if [ $GOOS = "windows" ]; then
        tui_output_name+='.exe'
    fi
    
    env GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "-X main.version=$VERSION -X main.buildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" \
        -o "dist/$tui_output_name" \
        cmd/tui/main.go
done

# 创建压缩包
cd dist/
for file in *; do
    if [[ $file == *.exe ]]; then
        zip "${file%.exe}.zip" "$file"
    else
        tar -czf "${file}.tar.gz" "$file"
    fi
done

echo "Build complete!"
```

#### 执行构建
```bash
# 赋予执行权限
chmod +x scripts/build-release.sh

# 执行构建
./scripts/build-release.sh v1.1.0
```

### 3. 创建发布

#### GitHub Release
```bash
# 创建发布（使用GitHub CLI）
gh release create v1.1.0 \
    --title "Release v1.1.0" \
    --notes "## What's New\n\n- Performance improvements\n- Bug fixes\n- New features" \
    dist/*.tar.gz \
    dist/*.zip
```

#### 发布说明模板
```markdown
# Release v1.1.0

## 🚀 新功能
- 添加YAML格式支持
- 改进文件选择器界面
- 增加性能优化选项

## 🐛 问题修复
- 修复内存泄漏问题
- 修复大文件处理问题
- 修复Windows路径问题

## 📈 性能改进
- 提升扫描速度30%
- 减少内存使用20%
- 优化并发处理

## 📝 文档更新
- 更新使用文档
- 添加新的示例
- 改进API文档

## 🔧 其他
- 更新依赖包
- 改进测试覆盖
- 代码重构

## 📥 下载
- [Linux AMD64](link-to-linux-binary)
- [macOS AMD64](link-to-macos-binary)
- [Windows AMD64](link-to-windows-binary)

## 🙏 致谢
感谢所有贡献者的支持！
```

## 持续集成

### 1. GitHub Actions配置

#### .github/workflows/ci.yml
```yaml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.24, 1.23]
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Download dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Generate coverage report
      run: go tool cover -html=coverage.out -o coverage.html
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
    
    - name: Run golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
    
    - name: Run gosec security scanner
      uses: securecodewarrior/github-action-gosec@master
      with:
        args: '-no-fail -fmt sarif -out results.sarif ./...'
    
    - name: Upload SARIF file
      uses: github/codeql-action/upload-sarif@v2
      with:
        sarif_file: results.sarif

  build:
    needs: test
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: 1.24
    
    - name: Build CLI
      run: go build -v -o code-context-generator cmd/cli/main.go
    
    - name: Build TUI
      run: go build -v -o code-context-generator-tui cmd/tui/main.go
    
    - name: Test build artifacts
      run: |
        ./code-context-generator --help
        ./code-context-generator-tui --help || true
    
    - name: Upload build artifacts
      uses: actions/upload-artifact@v3
      with:
        name: binaries
        path: |
          code-context-generator
          code-context-generator-tui

  release:
    needs: build
    runs-on: ubuntu-latest
    if: startsWith(github.ref, 'refs/tags/')
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: 1.24
    
    - name: Run GoReleaser
      uses: goreleaser/goreleaser-action@v4
      with:
        version: latest
        args: release --rm-dist
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### 2. 代码质量检查

#### golangci-lint配置
```yaml
# .golangci.yml
run:
  timeout: 5m
  issues-exit-code: 1
  tests: true

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true
  
  govet:
    check-shadowing: true
    enable-all: true
  
  gocyclo:
    min-complexity: 15
  
  maligned:
    suggest-new: true
  
  dupl:
    threshold: 100
  
  goconst:
    min-len: 3
    min-occurrences: 3

linters:
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - gochecknoinits
    - goconst
    - gocritic
    - gocyclo
    - gofmt
    - goimports
    - golint
    - gomnd
    - goprintffuncname
    - gosec
    - gosimple
    - govet
    - ineffassign
    - interfacer
    - lll
    - misspell
    - nakedret
    - rowserrcheck
    - scopelint
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace
  
  disable:
    - maligned  # 已被govet取代

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - gocyclo
        - errcheck
        - dupl
        - gosec
        - lll
```

## 最佳实践

### 1. 代码组织
- 保持包的小而专注
- 使用清晰的命名约定
- 遵循单一职责原则
- 编写可测试的代码

### 2. 错误处理
- 总是检查错误
- 包装错误以添加上下文
- 使用自定义错误类型
- 提供有用的错误信息

### 3. 性能优化
- 使用基准测试识别性能瓶颈
- 避免过早优化
- 使用性能分析工具
- 考虑内存分配

### 4. 文档编写
- 为所有导出的类型和函数编写文档
- 使用示例代码
- 保持文档更新
- 使用清晰的示例

### 5. 测试策略
- 编写单元测试覆盖核心逻辑
- 使用表格驱动测试
- 测试错误情况
- 保持测试简单和快速

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

## 获取帮助

### 资源
- [Go官方文档](https://golang.org/doc/)
- [Go语言规范](https://golang.org/ref/spec)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go代码审查评论](https://github.com/golang/go/wiki/CodeReviewComments)

### 社区
- [Go Forum](https://forum.golangbridge.org/)
- [Reddit r/golang](https://www.reddit.com/r/golang/)
- [Stack Overflow](https://stackoverflow.com/questions/tagged/go)
- [Go Slack](https://gophers.slack.com/)

### 项目支持
- 项目Issues: [GitHub Issues](https://github.com/yourusername/code-context-generator/issues)
- 开发文档: [开发文档链接]
- 邮件列表: dev@example.com