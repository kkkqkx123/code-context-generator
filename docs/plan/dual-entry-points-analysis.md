# 双入口点架构分析与改进方案

## 概述

代码上下文生成器项目包含两个主要的程序入口点，分别服务于不同的使用场景和需求。本文档分析这两个入口点的关系、差异以及存在的问题，并提出改进方案。

## 入口点对比分析

### 1. 根目录 main.go（简单入口）

**路径**: `/main.go`

**特点**:
- **简单直观**: 直接运行，无需学习命令结构
- **交互式体验**: 强制文件选择和格式选择交互
- **手动参数解析**: 自定义参数解析逻辑
- **轻量级**: 代码量少，依赖简单
- **快速上手**: 适合新用户快速体验功能

**支持的命令格式**:
```bash
go run main.go [generate] -f markdown -o output.md
go run main.go -h
go run main.go  # 进入交互式模式
```

**优点**:
- 零学习成本
- 交互式体验友好
- 快速验证功能

**缺点**:
- 功能相对简单
- 配置复用不足
- 扩展性有限

### 2. cmd/cli/main.go（专业CLI）

**路径**: `/cmd/cli/main.go`

**特点**:
- **专业CLI框架**: 基于Cobra构建的完整命令行应用
- **丰富的命令结构**: 支持generate、select、config、autocomplete等子命令
- **完善的参数体系**: 支持长参数、短参数、默认值、参数验证
- **配置管理**: 支持配置文件加载和环境变量
- **高级功能**: 支持递归控制、文件过滤、内容包含等高级选项

**支持的命令格式**:
```bash
go run cmd/cli/main.go generate -f markdown -o output.md --max-depth 3
go run cmd/cli/main.go select -f json -o selected.json
go run cmd/cli/main.go config show
go run cmd/cli/main.go autocomplete --type file
```

**优点**:
- 功能完整丰富
- 配置复用充分
- 扩展性强
- 适合自动化脚本

**缺点**:
- 学习成本相对较高
- 需要理解命令结构

## 当前存在的问题

### 1. 配置复用问题

**根目录 main.go 存在的问题**:
- ❌ 未加载配置文件（config.yaml）
- ❌ 未使用环境变量中的配置
- ❌ 硬编码的排除规则（`.git`, `node_modules`等）
- ❌ 忽略规则未复用已有配置

**具体表现**:
```go
// 根目录 main.go - 硬编码配置
selectOptions := &types.SelectOptions{
    Recursive:       true,
    ShowHidden:      false,
    MaxDepth:        0,
    IncludePatterns: []string{},
    ExcludePatterns: []string{".git", "node_modules", "*.exe", "*.dll"}, // 硬编码
}
```

对比 CLI 版本的配置加载:
```go
// cmd/cli/main.go - 正确的配置加载
configManager := config.NewManager()
if configPath != "" {
    if err := configManager.Load(configPath); err != nil {
        return fmt.Errorf("加载配置文件失败: %w", err)
    }
} else {
    defaultConfigPath := "config.yaml"
    configManager.Load(defaultConfigPath) // 忽略错误，使用默认配置
}
cfg = configManager.Get()
```

### 2. 工具复用问题

- ❌ 使用了不同的文件系统遍历器
- ❌ 未复用CLI版本的高级功能（如文件内容读取、哈希计算）
- ❌ 缺少标准化处理（如换行符标准化）

### 3. 用户体验不一致

- ❌ 两个入口点的参数格式不完全兼容
- ❌ 输出格式和处理逻辑有差异
- ❌ 错误处理机制不同

## 改进方案

### 方案一：配置复用改进

修改根目录 `main.go`，使其正确加载和使用配置：

```go
func main() {
    // 首先加载.env文件（如果存在）
    if err := env.LoadEnv(""); err != nil {
        log.Printf("警告: 加载.env文件失败: %v", err)
    }

    // 创建配置管理器并加载配置
    configManager := config.NewManager()
    
    // 尝试加载配置文件
    defaultConfigPath := "config.yaml"
    if err := configManager.Load(defaultConfigPath); err != nil {
        log.Printf("警告: 加载配置文件失败，使用默认配置: %v", err)
    }
    
    cfg := configManager.Get()
    
    // 使用配置中的设置
    selectOptions := &types.SelectOptions{
        Recursive:       cfg.FileProcessing.Recursive,
        ShowHidden:      cfg.FileProcessing.IncludeHidden,
        MaxDepth:        cfg.Filters.MaxDepth,
        IncludePatterns: cfg.Filters.IncludePatterns,
        ExcludePatterns: cfg.Filters.ExcludePatterns, // 使用配置文件中的排除规则
    }
    
    // ... 其余代码
}
```

### 方案二：工具复用改进

统一使用CLI版本的文件系统遍历器和处理逻辑：

```go
// 替换原有的walker使用方式
walker := filesystem.NewFileSystemWalker(types.WalkOptions{})

walkOptions := &types.WalkOptions{
    MaxDepth:        cfg.Filters.MaxDepth,
    MaxFileSize:     cfg.Filters.MaxFileSize,
    ExcludePatterns: cfg.Filters.ExcludePatterns,
    IncludePatterns: cfg.Filters.IncludePatterns,
    FollowSymlinks:  cfg.Filters.FollowSymlinks,
    ShowHidden:      cfg.FileProcessing.IncludeHidden,
    ExcludeBinary:   cfg.Filters.ExcludeBinary,
}

result, err := walker.Walk(currentDir, walkOptions)
```

### 方案三：参数标准化

使根目录版本的参数与CLI版本保持一致：

```go
// 支持更多CLI版本的参数
var (
    format, outputFile string
    showHelp, verbose bool
    maxDepth int
    excludeBinary bool
)

// 解析更多参数
for i := 0; i < len(args); i++ {
    arg := args[i]
    switch arg {
    case "-f", "--format":
        if i+1 < len(args) {
            format = args[i+1]
            i++
        }
    case "-o", "--output":
        if i+1 < len(args) {
            outputFile = args[i+1]
            i++
        }
    case "-v", "--verbose":
        verbose = true
    case "--max-depth":
        if i+1 < len(args) {
            maxDepth, _ = strconv.Atoi(args[i+1])
            i++
        }
    case "--exclude-binary":
        excludeBinary = true
    case "-h", "--help":
        showHelp = true
    case "generate":
        // 忽略generate参数
    }
}
```

### 方案四：输出标准化

使用CLI版本的输出标准化处理：

```go
// 格式化输出后，使用标准化处理
outputData, err := formatter.Format(contextData)
if err != nil {
    log.Fatalf("格式化输出失败: %v", err)
}

// 标准化换行符为当前操作系统格式
normalizedData := utils.NormalizeLineEndings(outputData)

// 保存到文件
if err := os.WriteFile(finalOutputFile, []byte(normalizedData), 0644); err != nil {
    log.Fatalf("写入输出文件失败: %v", err)
}
```

## 推荐实施策略

### 阶段一：配置复用（高优先级）
1. 修改根目录 `main.go` 加载配置文件
2. 使用配置中的忽略规则和扫描参数
3. 保持向后兼容性

### 阶段二：工具统一（中优先级）
1. 统一文件系统遍历器使用
2. 标准化输出处理
3. 错误处理机制统一

### 阶段三：功能增强（低优先级）
1. 添加更多CLI版本的参数支持
2. 增强交互式体验
3. 完善帮助文档

## 实施建议

1. **保持简单入口的简洁性**：即使改进配置复用，也要保持根目录版本的简单易用特点
2. **渐进式改进**：分阶段实施，避免一次性大幅修改导致稳定性问题
3. **充分测试**：每个阶段都要充分测试，确保不影响现有功能
4. **文档同步**：更新相关文档，说明两个入口点的差异和使用场景

## 使用建议

### 何时使用根目录 main.go
- 快速体验功能
- 简单的交互式使用
- 不需要复杂配置的场景
- 新用户初次使用

### 何时使用 cmd/cli/main.go
- 需要自动化脚本
- 复杂的文件过滤和处理需求
- 需要配置文件管理
- 专业用户和高级功能需求

通过合理的架构设计和渐进式改进，可以让两个入口点各自发挥优势，为不同需求的用户提供最佳的使用体验。