package formatter

import (
	"fmt"
	"strings"

	"code-context-generator/internal/formatter/encoding"
	"code-context-generator/pkg/types"
)

// MarkdownFormatter Markdown格式转换器
type MarkdownFormatter struct {
	BaseFormatter
	config *types.Config
}

// NewMarkdownFormatter 创建Markdown格式转换器
func NewMarkdownFormatter(config *types.Config) Formatter {
	return &MarkdownFormatter{
		BaseFormatter: BaseFormatter{
			name:        "Markdown",
			description: "Markdown format",
			config:      nil,
		},
		config: config,
	}
}

// Format 格式化上下文数据
func (f *MarkdownFormatter) Format(data types.ContextData) (string, error) {
	// 检查是否包含元信息
	includeMetadata := false // 默认不包含元信息
	if f.config != nil {
		includeMetadata = f.config.Output.IncludeMetadata
	}

	if f.config != nil && f.config.Formats.Markdown.Structure != nil {
		// 使用自定义结构
		customData := f.applyCustomStructure(data)
		// 将自定义结构转换为Markdown格式
		return f.convertToMarkdown(customData)
	}

	var result strings.Builder

	// 标题
	result.WriteString("# 代码上下文\n\n")

	// 统计信息
	result.WriteString("## 统计信息\n\n")
	result.WriteString(fmt.Sprintf("- **文件数量**: %d\n", data.FileCount))
	result.WriteString(fmt.Sprintf("- **文件夹数量**: %d\n", data.FolderCount))
	result.WriteString(fmt.Sprintf("- **总大小**: %d 字节\n", data.TotalSize))
	result.WriteString("\n")

	// 元信息（如果包含）
	if includeMetadata && len(data.Metadata) > 0 {
		result.WriteString("## 元信息\n\n")
		for key, value := range data.Metadata {
			result.WriteString(fmt.Sprintf("- **%s**: %v\n", key, value))
		}
		result.WriteString("\n")
	}

	// 文件夹列表
	if len(data.Folders) > 0 {
		result.WriteString("## 文件夹\n\n")
		for _, folder := range data.Folders {
			result.WriteString(fmt.Sprintf("- **%s** (%s)\n", folder.Name, folder.Path))
			result.WriteString(fmt.Sprintf("  - 大小: %d 字节\n", folder.Size))
			result.WriteString(fmt.Sprintf("  - 文件数量: %d\n", folder.Count))
			result.WriteString(fmt.Sprintf("  - 修改时间: %s\n", folder.ModTime.Format("2006-01-02 15:04:05")))
			if folder.IsHidden {
				result.WriteString("  - 隐藏: 是\n")
			}
			result.WriteString("\n")
		}
	}

	// 文件列表
	if len(data.Files) > 0 {
		result.WriteString("## 文件\n\n")
		for _, file := range data.Files {
			result.WriteString(fmt.Sprintf("### %s\n\n", file.Name))
			result.WriteString(fmt.Sprintf("- **路径**: %s\n", file.Path))
			result.WriteString(fmt.Sprintf("- **大小**: %d 字节\n", file.Size))
			result.WriteString(fmt.Sprintf("- **修改时间**: %s\n", file.ModTime.Format("2006-01-02 15:04:05")))
			if file.IsBinary {
				result.WriteString("- **类型**: 二进制文件\n")
			} else {
				result.WriteString("- **类型**: 文本文件\n")
			}
			if file.IsHidden {
				result.WriteString("- **隐藏**: 是\n")
			}
			
			// 文件内容
			if !file.IsBinary {
				result.WriteString("\n#### 内容\n\n")
				result.WriteString("```\n")
				// 限制内容长度以避免Markdown文件过大
				content := file.Content
				if len(content) > 1000 {
					content = content[:1000] + "\n... (内容已截断)"
				}
				result.WriteString(content)
				result.WriteString("\n```\n")
			}
			result.WriteString("\n")
		}
	}

	resultStr := result.String()

	// 检查配置中是否指定了编码格式
	if f.config != nil && f.config.Formats.Markdown.Encoding != "" && f.config.Formats.Markdown.Encoding != "utf-8" {
		// 转换编码
		encodedOutput, err := encoding.ConvertEncoding(resultStr, f.config.Formats.Markdown.Encoding)
		if err != nil {
			return "", fmt.Errorf("编码转换失败: %w", err)
		}
		return encodedOutput, nil
	}

	return resultStr, nil
}

// FormatFile 格式化单个文件
func (f *MarkdownFormatter) FormatFile(file types.FileInfo) (string, error) {
	var result strings.Builder

	// 文件标题
	result.WriteString(fmt.Sprintf("# %s\n\n", file.Name))

	// 文件信息
	result.WriteString("## 文件信息\n\n")
	result.WriteString(fmt.Sprintf("- **路径**: %s\n", file.Path))
	result.WriteString(fmt.Sprintf("- **大小**: %d 字节\n", file.Size))
	result.WriteString(fmt.Sprintf("- **修改时间**: %s\n", file.ModTime.Format("2006-01-02 15:04:05")))
	if file.IsBinary {
		result.WriteString("- **类型**: 二进制文件\n")
	} else {
		result.WriteString("- **类型**: 文本文件\n")
	}
	if file.IsHidden {
		result.WriteString("- **隐藏**: 是\n")
	}
	result.WriteString("\n")

	// 文件内容
	if !file.IsBinary {
		result.WriteString("## 内容\n\n")
		result.WriteString("```\n")
		result.WriteString(file.Content)
		result.WriteString("\n```\n")
	} else {
		result.WriteString("## 内容\n\n")
		result.WriteString("[二进制文件 - 内容未显示]\n")
	}

	resultStr := result.String()

	// 检查配置中是否指定了编码格式
	if f.config != nil && f.config.Formats.Markdown.Encoding != "" && f.config.Formats.Markdown.Encoding != "utf-8" {
		// 转换编码
		encodedOutput, err := encoding.ConvertEncoding(resultStr, f.config.Formats.Markdown.Encoding)
		if err != nil {
			return "", fmt.Errorf("编码转换失败: %w", err)
		}
		return encodedOutput, nil
	}

	return resultStr, nil
}

// FormatFolder 格式化文件夹
func (f *MarkdownFormatter) FormatFolder(folder types.FolderInfo) (string, error) {
	var result strings.Builder

	// 文件夹标题
	result.WriteString(fmt.Sprintf("# %s\n\n", folder.Name))

	// 文件夹信息
	result.WriteString("## 文件夹信息\n\n")
	result.WriteString(fmt.Sprintf("- **路径**: %s\n", folder.Path))
	result.WriteString(fmt.Sprintf("- **大小**: %d 字节\n", folder.Size))
	result.WriteString(fmt.Sprintf("- **修改时间**: %s\n", folder.ModTime.Format("2006-01-02 15:04:05")))
	result.WriteString(fmt.Sprintf("- **文件数量**: %d\n", folder.Count))
	if folder.IsHidden {
		result.WriteString("- **隐藏**: 是\n")
	}
	result.WriteString("\n")

	// 子文件夹
	if len(folder.Folders) > 0 {
		result.WriteString("## 子文件夹\n\n")
		for _, subFolder := range folder.Folders {
			result.WriteString(fmt.Sprintf("### %s\n\n", subFolder.Name))
			result.WriteString(fmt.Sprintf("- **路径**: %s\n", subFolder.Path))
			result.WriteString(fmt.Sprintf("- **大小**: %d 字节\n", subFolder.Size))
			result.WriteString(fmt.Sprintf("- **文件数量**: %d\n", subFolder.Count))
			result.WriteString(fmt.Sprintf("- **修改时间**: %s\n", subFolder.ModTime.Format("2006-01-02 15:04:05")))
			if subFolder.IsHidden {
				result.WriteString("- **隐藏**: 是\n")
			}
			result.WriteString("\n")
		}
	}

	// 文件
	if len(folder.Files) > 0 {
		result.WriteString("## 文件\n\n")
		for _, file := range folder.Files {
			result.WriteString(fmt.Sprintf("### %s\n\n", file.Name))
			result.WriteString(fmt.Sprintf("- **路径**: %s\n", file.Path))
			result.WriteString(fmt.Sprintf("- **大小**: %d 字节\n", file.Size))
			result.WriteString(fmt.Sprintf("- **修改时间**: %s\n", file.ModTime.Format("2006-01-02 15:04:05")))
			if file.IsBinary {
				result.WriteString("- **类型**: 二进制文件\n")
			} else {
				result.WriteString("- **类型**: 文本文件\n")
			}
			if file.IsHidden {
				result.WriteString("- **隐藏**: 是\n")
			}
			
			// 文件内容
			if !file.IsBinary {
				result.WriteString("\n#### 内容\n\n")
				result.WriteString("```\n")
				// 限制内容长度以避免Markdown文件过大
				content := file.Content
				if len(content) > 1000 {
					content = content[:1000] + "\n... (内容已截断)"
				}
				result.WriteString(content)
				result.WriteString("\n```\n")
			}
			result.WriteString("\n")
		}
	}

	resultStr := result.String()

	// 检查配置中是否指定了编码格式
	if f.config != nil && f.config.Formats.Markdown.Encoding != "" && f.config.Formats.Markdown.Encoding != "utf-8" {
		// 转换编码
		encodedOutput, err := encoding.ConvertEncoding(resultStr, f.config.Formats.Markdown.Encoding)
		if err != nil {
			return "", fmt.Errorf("编码转换失败: %w", err)
		}
		return encodedOutput, nil
	}

	return resultStr, nil
}

// convertToMarkdown 将自定义结构转换为Markdown格式
func (f *MarkdownFormatter) convertToMarkdown(data interface{}) (string, error) {
	// 这里可以实现更复杂的转换逻辑
	// 目前简单地将数据转换为JSON字符串，然后包装在代码块中
	var result strings.Builder
	result.WriteString("# 自定义结构数据\n\n")
	result.WriteString("```json\n")
	// 这里应该实现更复杂的转换逻辑
	result.WriteString("// 自定义结构转换待实现\n")
	result.WriteString("```\n")
	return result.String(), nil
}