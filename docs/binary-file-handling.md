# 二进制文件处理文档

## 概述

本项目实现了智能的二进制文件检测和处理机制，确保在构建代码上下文时能够正确处理文本文件和二进制文件，避免二进制文件内容导致的处理问题。

## 二进制文件检测原理

### 检测方法

项目使用基于文件扩展名的检测方法，通过`internal/utils/utils.go`中的工具函数实现：

```go
// IsTextFile 判断文件是否为文本文件
func IsTextFile(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    
    // 预定义的文本文件扩展名列表
    textExts := map[string]bool{
        ".txt": true, ".md": true, ".go": true, ".py": true, ".js": true,
        ".html": true, ".css": true, ".json": true, ".xml": true, ".yaml": true,
        ".yml": true, ".toml": true, ".ini": true, ".cfg": true, ".conf": true,
        ".sh": true, ".bat": true, ".cmd": true, ".ps1": true, ".bash": true,
        ".zsh": true, ".fish": true, ".c": true, ".cpp": true, ".h": true,
        ".hpp": true, ".java": true, ".cs": true, ".php": true, ".rb": true,
        ".rs": true, ".swift": true, ".kt": true, ".scala": true, ".r": true,
        ".m": true, ".mm": true, ".pl": true, ".lua": true, ".vim": true,
        ".el": true, ".lisp": true, ".sql": true, ".vimrc": true, ".gitignore": true,
        ".dockerignore": true, ".editorconfig": true, ".env": true, ".properties": true,
        ".gradle": true, ".cmake": true, ".make": true, ".mk": true, ".dockerfile": true,
        ".jenkinsfile": true, ".travis.yml": true, ".gitlab-ci.yml": true, ".github": true,
    }
    
    // 检查文件扩展名是否在文本文件列表中
    if isText, exists := textExts[ext]; exists && isText {
        return true
    }
    
    // 没有扩展名的文件默认为文本文件
    if ext == "" {
        return true
    }
    
    return false
}

// IsBinaryFile 判断文件是否为二进制文件
func IsBinaryFile(filename string) bool {
    return !IsTextFile(filename)
}
```

### 支持的文本文件类型

系统支持超过50种常见的文本文件扩展名，包括：

- **编程语言**: `.go`, `.py`, `.js`, `.java`, `.c`, `.cpp`, `.cs`, `.php`, `.rb`, `.rs`等
- **标记语言**: `.html`, `.xml`, `.json`, `.yaml`, `.toml`, `.md`等
- **配置文件**: `.ini`, `.cfg`, `.conf`, `.properties`, `.env`等
- **脚本文件**: `.sh`, `.bat`, `.cmd`, `.ps1`, `.bash`等
- **构建文件**: `.gradle`, `.cmake`, `.make`, `.mk`, `dockerfile`等
- **版本控制**: `.gitignore`, `.dockerignore`, `.editorconfig`等
- **CI/CD配置**: `.travis.yml`, `.gitlab-ci.yml`, `.jenkinsfile`等

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