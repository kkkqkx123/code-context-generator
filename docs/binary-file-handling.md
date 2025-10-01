# 二进制文件处理文档

## 概述

本项目实现了智能的二进制文件检测和处理机制，确保在构建代码上下文时能够正确处理文本文件和二进制文件，避免二进制文件内容导致的处理问题。

## 二进制文件检测原理

### 检测方法

项目使用智能的二进制文件检测方法，结合文件扩展名和内容分析，通过`internal/utils/utils.go`中的工具函数实现：

```go
// IsTextFile 检查是否为文本文件
func IsTextFile(path string) bool {
    // 首先检查文件扩展名
    ext := strings.ToLower(filepath.Ext(path))
    textExtensions := []string{
        ".txt", ".md", ".json", ".xml", ".yaml", ".yml", ".toml",
        ".go", ".py", ".js", ".ts", ".java", ".cpp", ".c", ".h",
        ".html", ".css", ".scss", ".sass", ".sql", ".sh", ".bat",
        ".ps1", ".rb", ".php", ".rs", ".swift", ".kt", ".scala",
    }

    for _, textExt := range textExtensions {
        if ext == textExt {
            return true
        }
    }

    // 如果没有扩展名，尝试读取文件内容来判断
    if ext == "" {
        file, err := os.Open(path)
        if err != nil {
            return false // 无法打开文件，假设为二进制文件
        }
        defer file.Close()

        // 读取前512字节来判断是否为文本文件
        buffer := make([]byte, 512)
        n, err := file.Read(buffer)
        if err != nil && err != io.EOF {
            return false // 读取错误，假设为二进制文件
        }

        // 检查是否包含null字节（二进制文件的标志）
        for i := 0; i < n; i++ {
            if buffer[i] == 0 {
                return false // 包含null字节，是二进制文件
            }
        }

        // 检查是否包含可打印字符
        printableCount := 0
        for i := 0; i < n; i++ {
            b := buffer[i]
            if b >= 32 && b <= 126 { // 可打印ASCII字符
                printableCount++
            } else if b == 9 || b == 10 || b == 13 { // tab, newline, carriage return
                printableCount++
            }
        }

        // 如果大部分字符都是可打印的，认为是文本文件
        if n > 0 && float64(printableCount)/float64(n) > 0.8 {
            return true
        }
    }

    return false
}

// IsBinaryFile 检查是否为二进制文件
func IsBinaryFile(path string) bool {
    return !IsTextFile(path)
}
```

### 支持的文本文件类型

系统支持多种常见的文本文件扩展名，包括：

- **编程语言**: `.go`, `.py`, `.js`, `.ts`, `.java`, `.cpp`, `.c`, `.rb`, `.php`, `.rs`, `.swift`, `.kt`, `.scala`
- **标记语言**: `.html`, `.xml`, `.json`, `.yaml`, `.yml`, `.toml`, `.md`
- **样式文件**: `.css`, `.scss`, `.sass`
- **脚本文件**: `.sh`, `.bat`, `.ps1`, `.sql`

### 无扩展名文件处理

对于没有扩展名的文件，系统采用智能内容分析算法：

1. **内容采样**: 读取文件前512字节进行内容分析
2. **二进制检测**: 检查是否包含null字节（二进制文件的标志）
3. **字符分析**: 统计可打印字符比例（ASCII 32-126）和常见控制字符（tab、换行、回车）
4. **智能判断**: 如果超过80%的字符为可打印字符，则判定为文本文件

这种智能检测机制确保无扩展名的文本文件（如README、LICENSE、Makefile等）能够被正确识别和处理。

## 文件处理流程

### 1. 文件系统遍历阶段

在`internal/filesystem/filesystem.go`中，文件系统遍历器会在以下阶段进行二进制文件检查：

```go
func (w *FileSystemWalker) shouldIncludeFile(path string, options *types.WalkOptions) bool {
    // 检查文件大小
    if !w.FilterBySize(path, options.MaxFileSize) {
        return false
    }
    
    // 检查是否为二进制文件（如果启用了二进制文件排除）
    if options.ExcludeBinary && utils.IsBinaryFile(path) {
        return false
    }
    
    // 其他过滤逻辑...
    return true
}
```

### 2. 文件内容读取阶段

在`GetFileInfo`方法中，系统会根据文件类型决定是否读取内容：

```go
func (w *FileSystemWalker) GetFileInfo(path string) (*types.FileInfo, error) {
    // 检查是否为二进制文件
    isBinary := !utils.IsTextFile(path)
    
    var content string
    if !isBinary {
        // 只读取文本文件的内容
        fileContent, err := os.ReadFile(path)
        if err != nil {
            return nil, fmt.Errorf("读取文件内容失败: %w", err)
        }
        content = string(fileContent)
    }
    
    return &types.FileInfo{
        Path:     path,
        Name:     info.Name(),
        Size:     info.Size(),
        ModTime:  info.ModTime(),
        IsDir:    info.IsDir(),
        Content:  content,
        IsBinary: isBinary,
    }, nil
}
```

### 3. 格式化输出阶段

各种格式化器会根据`IsBinary`字段处理文件内容：

#### JSON格式化器
```go
func (f *JSONFormatter) FormatFile(file types.FileInfo) (string, error) {
    // 如果是二进制文件，不显示内容
    if file.IsBinary {
        file.Content = "[二进制文件 - 内容未显示]"
    }
    
    // 格式化逻辑...
}
```

#### Markdown格式化器
```go
func (f *MarkdownFormatter) FormatFile(file types.FileInfo) (string, error) {
    // 添加代码块（只针对文本文件）
    if !file.IsBinary {
        sb.WriteString("```")
        // 添加代码高亮等...
        sb.WriteString(file.Content)
        sb.WriteString("\n```\n")
    } else {
        sb.WriteString("**[二进制文件 - 内容未显示]**\n")
    }
    // ...
}
```

## 配置选项

### 命令行选项

在CLI工具中，可以通过以下选项控制二进制文件处理：

```bash
# 排除二进制文件（默认行为）
code-context-generator generate --exclude-binary path/to/directory

# 包含二进制文件（不推荐）
code-context-generator generate --exclude-binary=false path/to/directory
```

### 配置文件选项

在配置文件中，可以通过以下设置控制二进制文件处理：

```toml
[filters]
max_file_size = "10MB"
exclude_patterns = ["*.log", "*.tmp"]
exclude_binary = true  # 排除二进制文件
max_depth = 5
```

## 类型定义

### FileInfo结构体

在`pkg/types/types.go`中，文件信息结构体包含二进制文件标识：

```go
type FileInfo struct {
    Name     string    `yaml:"name" json:"name" toml:"name"`
    Path     string    `yaml:"path" json:"path" toml:"path"`
    Content  string    `yaml:"content" json:"content" toml:"content"`
    Size     int64     `yaml:"size" json:"size" toml:"size"`
    ModTime  time.Time `yaml:"mod_time" json:"mod_time" toml:"mod_time"`
    IsDir    bool      `yaml:"is_dir" json:"is_dir" toml:"is_dir"`
    IsHidden bool      `yaml:"is_hidden" json:"is_hidden" toml:"is_hidden"`
    IsBinary bool      `yaml:"is_binary" json:"is_binary" toml:"is_binary"`
}
```

### WalkOptions结构体

文件遍历选项包含二进制文件过滤设置：

```go
type WalkOptions struct {
    MaxDepth        int
    MaxFileSize     int64
    ExcludePatterns []string
    IncludePatterns []string
    FollowSymlinks  bool
    ShowHidden      bool
    ExcludeBinary   bool  // 是否排除二进制文件
}
```

### FiltersConfig结构体

过滤器配置包含二进制文件排除选项：

```go
type FiltersConfig struct {
    MaxFileSize     string   `yaml:"max_file_size" json:"max_file_size" toml:"max_file_size"`
    ExcludePatterns []string `yaml:"exclude_patterns" json:"exclude_patterns" toml:"exclude_patterns"`
    IncludePatterns []string `yaml:"include_patterns" json:"include_patterns" toml:"include_patterns"`
    MaxDepth        int      `yaml:"max_depth" json:"max_depth" toml:"max_depth"`
    FollowSymlinks  bool     `yaml:"follow_symlinks" json:"follow_symlinks" toml:"follow_symlinks"`
    ExcludeBinary   bool     `yaml:"exclude_binary" json:"exclude_binary" toml:"exclude_binary"`
}
```

## 优势

1. **安全性**: 避免二进制文件内容导致的编码问题或程序崩溃
2. **性能**: 不读取二进制文件内容，提高处理速度
3. **清晰度**: 在输出中明确标识二进制文件
4. **灵活性**: 可通过配置控制是否排除二进制文件
5. **扩展性**: 易于添加新的文本文件类型支持

## 使用建议

1. **默认配置**: 建议保持`exclude_binary = true`的默认设置
2. **自定义类型**: 如需支持新的文件类型，可修改`IsTextFile`函数
3. **性能优化**: 对于大型项目，二进制文件排除可以显著提高扫描速度
4. **内容验证**: 对于重要文件，建议先验证文件类型再进行处理

## 相关文件

- `internal/utils/utils.go` - 二进制文件检测函数
- `internal/filesystem/filesystem.go` - 文件系统遍历和过滤
- `internal/formatter/formatter.go` - 格式化输出处理
- `pkg/types/types.go` - 类型定义
- `internal/config/config.go` - 配置文件处理
- `cmd/cli/main.go` - CLI命令行选项