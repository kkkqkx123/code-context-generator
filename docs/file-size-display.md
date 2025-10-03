# 📊 智能文件大小显示

本文档介绍 code-context-generator 的智能文件大小显示功能。

## 🎯 功能概述

工具会根据文件大小自动选择合适的单位显示，提供更佳的用户体验：

- **≤1KB**: 显示为 `X B` (字节)
- **>1KB 且 ≤1MB**: 显示为 `X.X KB` (千字节)  
- **>1MB**: 显示为 `X.X MB` (兆字节)

## 📋 显示示例

```bash
# 1字节文件
💾 总大小: 1 B

# 1KB文件
💾 总大小: 1.0 KB

# 100KB文件
💾 总大小: 100.0 KB

# 1MB文件
💾 总大小: 1.0 MB

# 大文件
💾 总大小: 15.7 MB
```

## ⚙️ 技术实现

### 核心函数

使用 `internal/utils.FormatFileSize()` 函数进行单位转换：

```go
func FormatFileSize(bytes int64) string {
    if bytes <= 1024 {
        return fmt.Sprintf("%d B", bytes)
    } else if bytes <= 1024*1024 {
        return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
    } else {
        return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
    }
}
```

### 使用位置

该功能在以下位置使用：
- `main.go`: 命令行工具输出
- `cli/main.go`: CLI应用程序输出

### 修改历史

- **之前**: 固定显示为 `💾 总大小: 0.00 MB`
- **现在**: 智能单位选择 `💾 总大小: X.X 单位`

## 🧪 测试验证

创建不同大小的测试文件验证功能：

```bash
# 创建测试文件
echo "A" > test_1b.txt                    # 1字节
echo "A" * 1024 > test_1kb.txt          # 1KB  
echo "A" * 102400 > test_100kb.txt      # 100KB
echo "A" * 1048576 > test_1mb.txt       # 1MB

# 测试显示效果
./c-gen generate test_1b.txt -f markdown
./c-gen generate test_1kb.txt -f markdown
./c-gen generate test_100kb.txt -f markdown
./c-gen generate test_1mb.txt -f markdown
```

## 🔧 相关配置

虽然文件大小显示是自动的，但可以通过以下配置影响扫描的文件大小：

```toml
[file_processing]
max_file_size = 10485760  # 最大文件大小限制 (10MB)
```

## 💡 注意事项

1. **精度**: 显示精确到小数点后1位
2. **单位**: 严格遵循1024进制 (1KB = 1024B, 1MB = 1024KB)
3. **边界**: 1KB边界使用≤1024B规则，1MB边界使用≤1024KB规则
4. **性能**: 单位转换在显示时进行，不影响扫描性能

## 🚀 扩展可能

未来可考虑支持：
- GB、TB单位显示
- 可配置的小数精度
- 用户自定义单位阈值
- 国际化单位显示（KiB、MiB等）