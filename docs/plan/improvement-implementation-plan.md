# 配置复用改进实施计划

## 目标
解决根目录 `main.go` 未使用配置文件、忽略规则和环境变量的问题，使其与 `cmd/cli/main.go` 保持一致的行为和输出质量。

## 修改方案

### 1. 修改根目录 main.go - 配置加载

**修改前**（当前代码片段）：
```go
func main() {
    // 解析命令行参数
    var format, outputFile string
    var showHelp bool
    
    args := os.Args[1:]
    for i := 0; i < len(args); i++ {
        // ... 参数解析逻辑
    }
    
    if showHelp {
        // ... 帮助信息
    }
    
    // 硬编码的配置
    selectOptions := &types.SelectOptions{
        Recursive:       true,
        ShowHidden:      false,
        MaxDepth:        0,
        IncludePatterns: []string{},
        ExcludePatterns: []string{".git", "node_modules", "*.exe", "*.dll"}, // 硬编码排除规则
    }
    // ...
}
```

**修改后**（建议实现）：
```go
func main() {
    // 首先加载.env文件（如果存在）
    if err := env.LoadEnv(""); err != nil {
        log.Printf("警告: 加载.env文件失败: %v", err)
    }

    // 解析命令行参数
    var format, outputFile string
    var showHelp bool
    
    args := os.Args[1:]
    for i := 0; i < len(args); i++ {
        // ... 参数解析逻辑
    }
    
    if showHelp {
        // ... 帮助信息
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
    // ...
}
```

### 2. 修改根目录 main.go - 文件系统遍历

**修改前**：
```go
// 使用简单的文件遍历
walker := filesystem.NewFileSystemWalker(types.WalkOptions{})
result, err := walker.Walk(currentDir, selectOptions)
```

**修改后**：
```go
// 使用与CLI版本一致的文件系统遍历器
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

### 3. 修改根目录 main.go - 输出标准化

**修改前**：
```go
// 格式化并直接输出
outputData, err := formatter.Format(contextData)
if err != nil {
    log.Fatalf("格式化输出失败: %v", err)
}

// 保存到文件
if err := os.WriteFile(finalOutputFile, []byte(outputData), 0644); err != nil {
    log.Fatalf("写入输出文件失败: %v", err)
}
```

**修改后**：
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

### 4. 修改根目录 main.go - 错误处理增强

**修改前**：
```go
// 简单的错误处理
if err != nil {
    log.Fatalf("文件选择失败: %v", err)
}
```

**修改后**：
```go
// 增强的错误处理和日志记录
if err != nil {
    if verbose {
        log.Printf("详细错误信息: %+v", err)
    }
    log.Fatalf("文件选择失败: %v", err)
}
```

## 实施步骤

### 步骤 1: 添加必要的导入

在根目录 `main.go` 的导入部分添加：
```go
import (
    // ... 现有导入
    "github.com/yourusername/code-context-generator/internal/config"
    "github.com/yourusername/code-context-generator/internal/env"
    "github.com/yourusername/code-context-generator/internal/utils"
)
```

### 步骤 2: 修改主函数结构

1. 在 `main()` 函数开始处添加配置加载
2. 修改 `selectOptions` 的创建逻辑
3. 更新文件系统遍历调用
4. 添加输出标准化处理

### 步骤 3: 测试验证

1. **配置加载测试**：
   ```bash
   # 创建测试配置文件
   echo "filters:" > config.yaml
   echo "  exclude_patterns:" >> config.yaml
   echo "    - \"*.log\"" >> config.yaml
   echo "    - \"*.tmp\"" >> config.yaml
   
   # 测试配置是否生效
   go run main.go -f markdown -o test.md
   ```

2. **环境变量测试**：
   ```bash
   # 设置环境变量
   $env:CCG_MAX_DEPTH="2"
   
   # 测试环境变量是否生效
   go run main.go -f markdown -o test.md
   ```

3. **输出一致性测试**：
   ```bash
   # 使用相同参数测试两个入口点
   go run main.go -f markdown -o simple.md
   go run cmd/cli/main.go generate -f markdown -o cli.md
   
   # 比较输出文件
   fc simple.md cli.md
   ```

## 预期效果

### 改进前的问题
- ❌ 根目录版本忽略 `.git` 等目录，但使用硬编码规则
- ❌ 不支持配置文件中的自定义忽略规则
- ❌ 不支持环境变量配置
- ❌ 输出格式可能与CLI版本不一致

### 改进后的效果
- ✅ 正确加载和使用配置文件
- ✅ 支持环境变量配置
- ✅ 使用与CLI版本一致的忽略规则和处理逻辑
- ✅ 输出文件格式和质量保持一致
- ✅ 保持简单易用的特点

## 注意事项

1. **向后兼容性**：确保修改后不影响现有用户的使用习惯
2. **错误处理**：配置加载失败时应优雅降级，使用合理默认值
3. **性能考虑**：配置加载不应显著影响启动时间
4. **日志级别**：提供适当的日志输出，既不过于冗长也不过于简略

## 后续优化

1. **参数统一**：逐步支持更多CLI版本的参数选项
2. **功能增强**：在保持简洁的前提下添加实用功能
3. **文档完善**：更新相关文档，说明配置使用方法
4. **测试覆盖**：为配置加载和使用添加单元测试

通过实施这个改进计划，可以使根目录 `main.go` 在保持简单易用的同时，充分利用项目的配置体系和工具链，提供与专业CLI版本一致的高质量输出。