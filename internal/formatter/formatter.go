// Package formatter 提供多种格式的输出转换功能
package formatter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"code-context-generator/pkg/constants"
	"code-context-generator/pkg/types"
)

// Formatter 格式转换器接口
type Formatter interface {
	Format(data types.ContextData) (string, error)
	FormatFile(file types.FileInfo) (string, error)
	FormatFolder(folder types.FolderInfo) (string, error)
	GetName() string
	GetDescription() string
}

// BaseFormatter 基础格式转换器
type BaseFormatter struct {
	name        string
	description string
	config      *types.FormatConfig
}

// GetName 获取格式名称
func (f *BaseFormatter) GetName() string {
	return f.name
}

// GetDescription 获取格式描述
func (f *BaseFormatter) GetDescription() string {
	return f.description
}

// applyCustomStructure 应用自定义结构
func (f *BaseFormatter) applyCustomStructure(data types.ContextData) interface{} {
	// 根据配置应用自定义结构
	if f.config != nil && f.config.Structure != nil {
		// 创建基于实际数据的自定义结构
		result := make(map[string]interface{})
		
		// 应用结构映射
		if rootTag, ok := f.config.Structure["root"].(string); ok && rootTag != "" {
			result["XMLName"] = xml.Name{Local: rootTag}
		} else {
			result["XMLName"] = xml.Name{Local: "context"}
		}
		
		// 映射文件和文件夹数据
		if filesTag, ok := f.config.Structure["files"].(string); ok && filesTag != "" {
			result[filesTag] = map[string]interface{}{
				"file": data.Files,
			}
		} else {
			result["files"] = map[string]interface{}{
				"file": data.Files,
			}
		}
		
		if foldersTag, ok := f.config.Structure["folders"].(string); ok && foldersTag != "" {
			result[foldersTag] = map[string]interface{}{
				"folder": data.Folders,
			}
		} else {
			result["folders"] = map[string]interface{}{
				"folder": data.Folders,
			}
		}
		
		// 添加统计信息
		result["file_count"] = data.FileCount
		result["folder_count"] = data.FolderCount
		result["total_size"] = data.TotalSize
		
		return result
	}
	
	return data
}

// applyCustomFields 应用自定义字段映射
func (f *BaseFormatter) applyCustomFields(file types.FileInfo) interface{} {
	// 根据配置应用自定义字段映射
	if f.config != nil && f.config.Fields != nil {
		// 这里可以实现字段映射逻辑
		return f.config.Fields
	}
	return file
}

// JSONFormatter JSON格式转换器
type JSONFormatter struct {
	BaseFormatter
}

// NewJSONFormatter 创建JSON格式转换器
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
func (f *JSONFormatter) Format(data types.ContextData) (string, error) {
	if f.config != nil && f.config.Structure != nil {
		// 使用自定义结构
		customData := f.applyCustomStructure(data)
		output, err := json.MarshalIndent(customData, "", "  ")
		if err != nil {
			return "", fmt.Errorf("JSON格式化失败: %w", err)
		}
		return string(output), nil
	}

	// 默认结构
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON格式化失败: %w", err)
	}
	return string(output), nil
}

// FormatFile 格式化单个文件
func (f *JSONFormatter) FormatFile(file types.FileInfo) (string, error) {
	// 如果是二进制文件，不显示内容
	if file.IsBinary {
		file.Content = "[二进制文件 - 内容未显示]"
	}
	
	if f.config != nil && f.config.Fields != nil {
		// 使用自定义字段映射
		customFile := f.applyCustomFields(file)
		output, err := json.MarshalIndent(customFile, "", "  ")
		if err != nil {
			return "", fmt.Errorf("JSON文件格式化失败: %w", err)
		}
		return string(output), nil
	}

	output, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON文件格式化失败: %w", err)
	}
	return string(output), nil
}

// FormatFolder 格式化文件夹
func (f *JSONFormatter) FormatFolder(folder types.FolderInfo) (string, error) {
	output, err := json.MarshalIndent(folder, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON文件夹格式化失败: %w", err)
	}
	return string(output), nil
}

// XMLFormatter XML格式转换器
type XMLFormatter struct {
	BaseFormatter
}

// NewXMLFormatter 创建XML格式转换器
func NewXMLFormatter(config *types.FormatConfig) Formatter {
	return &XMLFormatter{
		BaseFormatter: BaseFormatter{
			name:        "XML",
			description: "Extensible Markup Language format",
			config:      config,
		},
	}
}

// Format 格式化上下文数据
func (f *XMLFormatter) Format(data types.ContextData) (string, error) {
	// 创建可序列化的结构，避免map[string]interface{}
	type SerializableContextData struct {
		XMLName     xml.Name           `xml:"context"`
		Files       []types.FileInfo   `xml:"files>file"`
		Folders     []types.FolderInfo `xml:"folders>folder"`
		FileCount   int                `xml:"file_count"`
		FolderCount int                `xml:"folder_count"`
		TotalSize   int64              `xml:"total_size"`
	}

	serializableData := SerializableContextData{
		Files:       data.Files,
		Folders:     data.Folders,
		FileCount:   data.FileCount,
		FolderCount: data.FolderCount,
		TotalSize:   data.TotalSize,
	}

	if f.config != nil && f.config.Structure != nil {
		// 使用自定义结构
		customData := f.applyCustomStructure(data)
		output, err := xml.MarshalIndent(customData, "", "  ")
		if err != nil {
			return "", fmt.Errorf("XML格式化失败: %w", err)
		}
		return xml.Header + string(output), nil
	}

	// 默认结构
	output, err := xml.MarshalIndent(serializableData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("XML格式化失败: %w", err)
	}
	return xml.Header + string(output), nil
}

// FormatFile 格式化单个文件
func (f *XMLFormatter) FormatFile(file types.FileInfo) (string, error) {
	// 如果是二进制文件，不显示内容
	if file.IsBinary {
		file.Content = "[二进制文件 - 内容未显示]"
	}
	
	output, err := xml.MarshalIndent(file, "", "  ")
	if err != nil {
		return "", fmt.Errorf("XML文件格式化失败: %w", err)
	}
	return xml.Header + string(output), nil
}

// FormatFolder 格式化文件夹
func (f *XMLFormatter) FormatFolder(folder types.FolderInfo) (string, error) {
	output, err := xml.MarshalIndent(folder, "", "  ")
	if err != nil {
		return "", fmt.Errorf("XML文件夹格式化失败: %w", err)
	}
	return xml.Header + string(output), nil
}

// TOMLFormatter TOML格式转换器
type TOMLFormatter struct {
	BaseFormatter
}

// NewTOMLFormatter 创建TOML格式转换器
func NewTOMLFormatter(config *types.FormatConfig) Formatter {
	return &TOMLFormatter{
		BaseFormatter: BaseFormatter{
			name:        "TOML",
			description: "Tom's Obvious, Minimal Language format",
			config:      config,
		},
	}
}

// Format 格式化上下文数据
func (f *TOMLFormatter) Format(data types.ContextData) (string, error) {
	var buf strings.Builder

	// 写入文件部分
	if len(data.Files) > 0 {
		buf.WriteString("[files]\n")
		for i, file := range data.Files {
			buf.WriteString("  [[files.file]]\n")
			buf.WriteString(fmt.Sprintf("    path = \"%s\"\n", file.Path))
			buf.WriteString(fmt.Sprintf("    name = \"%s\"\n", file.Name))
			buf.WriteString(fmt.Sprintf("    size = %d\n", file.Size))
			buf.WriteString(fmt.Sprintf("    content = \"%s\"\n", escapeTOMLString(file.Content)))
			if i < len(data.Files)-1 {
				buf.WriteString("\n")
			}
		}
	}

	// 写入文件夹部分
	if len(data.Folders) > 0 {
		buf.WriteString("\n[folders]\n")
		for i, folder := range data.Folders {
			buf.WriteString("  [[folders.folder]]\n")
			buf.WriteString(fmt.Sprintf("    path = \"%s\"\n", folder.Path))
			buf.WriteString(fmt.Sprintf("    name = \"%s\"\n", folder.Name))
			buf.WriteString(fmt.Sprintf("    file_count = %d\n", len(folder.Files)))
			if i < len(data.Folders)-1 {
				buf.WriteString("\n")
			}
		}
	}

	return buf.String(), nil
}

// FormatFile 格式化单个文件
func (f *TOMLFormatter) FormatFile(file types.FileInfo) (string, error) {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("path = \"%s\"\n", file.Path))
	buf.WriteString(fmt.Sprintf("name = \"%s\"\n", file.Name))
	buf.WriteString(fmt.Sprintf("size = %d\n", file.Size))
	
	// 如果是二进制文件，不显示内容
	if file.IsBinary {
		buf.WriteString("content = \"[二进制文件 - 内容未显示]\"\n")
	} else {
		buf.WriteString(fmt.Sprintf("content = \"%s\"\n", escapeTOMLString(file.Content)))
	}
	
	buf.WriteString(fmt.Sprintf("mod_time = \"%s\"\n", file.ModTime.Format(time.RFC3339)))

	return buf.String(), nil
}

// FormatFolder 格式化文件夹
func (f *TOMLFormatter) FormatFolder(folder types.FolderInfo) (string, error) {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("path = \"%s\"\n", folder.Path))
	buf.WriteString(fmt.Sprintf("name = \"%s\"\n", folder.Name))
	buf.WriteString(fmt.Sprintf("file_count = %d\n", len(folder.Files)))
	buf.WriteString(fmt.Sprintf("mod_time = \"%s\"\n", folder.ModTime.Format(time.RFC3339)))

	return buf.String(), nil
}

// MarkdownFormatter Markdown格式转换器
type MarkdownFormatter struct {
	BaseFormatter
}

// NewMarkdownFormatter 创建Markdown格式转换器
func NewMarkdownFormatter(config *types.FormatConfig) Formatter {
	return &MarkdownFormatter{
		BaseFormatter: BaseFormatter{
			name:        "Markdown",
			description: "Markdown format with code blocks",
			config:      config,
		},
	}
}

// Format 格式化上下文数据
func (f *MarkdownFormatter) Format(data types.ContextData) (string, error) {
	var sb strings.Builder

	// 添加标题
	sb.WriteString("# 代码上下文\n\n")
	sb.WriteString(fmt.Sprintf("生成时间: %s\n\n", time.Now().Format(time.RFC3339)))

	// 添加文件部分
	if len(data.Files) > 0 {
		sb.WriteString("## 文件\n\n")
		for _, file := range data.Files {
			sb.WriteString(fmt.Sprintf("### %s\n\n", file.Name))
			sb.WriteString(fmt.Sprintf("- **路径**: `%s`\n", file.Path))
			sb.WriteString(fmt.Sprintf("- **大小**: %d 字节\n", file.Size))
			sb.WriteString(fmt.Sprintf("- **修改时间**: %s\n\n", file.ModTime.Format(time.RFC3339)))

			// 添加代码块（只针对文本文件）
			if !file.IsBinary {
				sb.WriteString("```")
				if ext := filepath.Ext(file.Path); ext != "" {
					sb.WriteString(strings.TrimPrefix(ext, "."))
				}
				sb.WriteString("\n")
				sb.WriteString(file.Content)
				sb.WriteString("\n```\n\n")
			} else {
				sb.WriteString("**[二进制文件 - 内容未显示]**\n\n")
			}
		}
	}

	// 添加文件夹部分
	if len(data.Folders) > 0 {
		sb.WriteString("## 文件夹\n\n")
		for _, folder := range data.Folders {
			sb.WriteString(fmt.Sprintf("### %s\n\n", folder.Name))
			sb.WriteString(fmt.Sprintf("- **路径**: `%s`\n", folder.Path))
			sb.WriteString(fmt.Sprintf("- **文件数**: %d\n", len(folder.Files)))
			sb.WriteString(fmt.Sprintf("- **文件数**: %d\n\n", len(folder.Files)))

			// 添加文件夹中的文件
			if len(folder.Files) > 0 {
				sb.WriteString("#### 文件列表\n\n")
				for _, file := range folder.Files {
					sb.WriteString(fmt.Sprintf("- `%s` (%d 字节)\n", file.Name, file.Size))
				}
				sb.WriteString("\n")
			}
		}
	}

	return sb.String(), nil
}

// FormatFile 格式化单个文件
func (f *MarkdownFormatter) FormatFile(file types.FileInfo) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## %s\n\n", file.Name))
	sb.WriteString(fmt.Sprintf("- **路径**: `%s`\n", file.Path))
	sb.WriteString(fmt.Sprintf("- **大小**: %d 字节\n", file.Size))
	sb.WriteString(fmt.Sprintf("- **修改时间**: %s\n\n", file.ModTime.Format(time.RFC3339)))

	// 添加代码块（只针对文本文件）
	if !file.IsBinary {
		sb.WriteString("```")
		if ext := filepath.Ext(file.Path); ext != "" {
			sb.WriteString(strings.TrimPrefix(ext, "."))
		}
		sb.WriteString("\n")
		sb.WriteString(file.Content)
		sb.WriteString("\n```\n")
	} else {
		sb.WriteString("**[二进制文件 - 内容未显示]**\n")
	}

	return sb.String(), nil
}

// FormatFolder 格式化文件夹
func (f *MarkdownFormatter) FormatFolder(folder types.FolderInfo) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## %s\n\n", folder.Name))
	sb.WriteString(fmt.Sprintf("- **路径**: `%s`\n", folder.Path))
	sb.WriteString(fmt.Sprintf("- **文件数**: %d\n", len(folder.Files)))
	sb.WriteString(fmt.Sprintf("- **文件数**: %d\n\n", len(folder.Files)))

	// 添加文件列表
	if len(folder.Files) > 0 {
		sb.WriteString("### 文件列表\n\n")
		for _, file := range folder.Files {
			sb.WriteString(fmt.Sprintf("- `%s` (%d 字节)\n", file.Name, file.Size))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// FormatterFactory 格式转换器工厂
type FormatterFactory struct {
	formatters map[string]Formatter
}

// NewFormatterFactory 创建格式转换器工厂
func NewFormatterFactory() *FormatterFactory {
	return &FormatterFactory{
		formatters: make(map[string]Formatter),
	}
}

// Register 注册格式转换器
func (ff *FormatterFactory) Register(format string, formatter Formatter) {
	ff.formatters[strings.ToLower(format)] = formatter
}

// Get 获取格式转换器
func (ff *FormatterFactory) Get(format string) (Formatter, error) {
	formatter, exists := ff.formatters[strings.ToLower(format)]
	if !exists {
		return nil, fmt.Errorf("不支持的格式: %s", format)
	}
	return formatter, nil
}

// GetSupportedFormats 获取支持的格式列表
func (ff *FormatterFactory) GetSupportedFormats() []string {
	formats := make([]string, 0, len(ff.formatters))
	for format := range ff.formatters {
		formats = append(formats, format)
	}
	return formats
}

// NewFormatter 创建格式转换器
func NewFormatter(format string) (Formatter, error) {
	factory := CreateDefaultFactory(nil)
	return factory.Get(format)
}

// CreateDefaultFactory 创建默认的格式转换器工厂
func CreateDefaultFactory(configs map[string]*types.FormatConfig) *FormatterFactory {
	factory := NewFormatterFactory()

	// 注册所有支持的格式
	factory.Register(constants.FormatJSON, NewJSONFormatter(configs[constants.FormatJSON]))
	factory.Register(constants.FormatXML, NewXMLFormatter(configs[constants.FormatXML]))
	factory.Register(constants.FormatTOML, NewTOMLFormatter(configs[constants.FormatTOML]))
	factory.Register(constants.FormatMarkdown, NewMarkdownFormatter(configs[constants.FormatMarkdown]))

	return factory
}

// 辅助方法

func escapeTOMLString(s string) string {
	// 简单的TOML字符串转义
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}
