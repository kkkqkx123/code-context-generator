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
	// 检查是否启用AI优化
	if f.config != nil && f.config.Output.AIOptimized {
		return f.formatAIOptimized(data)
	}

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

// formatAIOptimized AI优化格式化
func (f *MarkdownFormatter) formatAIOptimized(data types.ContextData) (string, error) {
	var result strings.Builder

	// 获取语言列表
	// 根据文件扩展名获取语言列表
	languageMap := make(map[string]struct{})
	for _, file := range data.Files {
		language := f.detectLanguage(file.Name)
		languageMap[language] = struct{}{}
	}
	languages := make([]string, 0, len(languageMap))
	for lang := range languageMap {
		languages = append(languages, lang)
	}

	// 创建AI摘要生成器
	summaryGenerator := NewAISummaryGenerator(f.config)
	aiSummary := summaryGenerator.GenerateSummary(data.FileCount, data.TotalSize, languages)

	// 创建模板系统
	templateSystem := NewTemplateSystem(f.config)
	
	// 创建模板数据（用于未来的模板扩展）
	templateData := templateSystem.CreateDefaultTemplateData(data.FileCount, data.FolderCount, data.TotalSize, languages)
	_ = templateData // 当前未使用，但为模板系统预留

	// 生成AI摘要部分
	result.WriteString("# AI优化代码上下文分析\n\n")
	result.WriteString("## 项目摘要\n\n")
	result.WriteString(fmt.Sprintf("%v", aiSummary))
	result.WriteString("\n\n")

	// 生成目录结构
	result.WriteString("## 项目结构\n\n")
	directoryStructure := f.generateDirectoryStructure(data)
	result.WriteString(directoryStructure)
	result.WriteString("\n\n")

	// 生成文件内容
	result.WriteString("## 代码内容\n\n")
	for _, file := range data.Files {
		fileContent := f.formatFileAIOptimized(file)
		result.WriteString(fileContent)
		result.WriteString("\n\n")
	}

	// 生成自定义指令
	if f.config != nil && f.config.Output.AIInstructions.Enabled {
		instructionLoader := NewInstructionLoader(f.config)
		customInstructions, err := instructionLoader.LoadInstructions()
		if err == nil && customInstructions != "" {
			result.WriteString("## AI分析指令\n\n")
			result.WriteString(customInstructions)
			result.WriteString("\n\n")
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

// formatFileAIOptimized AI优化文件格式化
func (f *MarkdownFormatter) formatFileAIOptimized(file types.FileInfo) string {
	var result strings.Builder

	// 文件标题和信息
	result.WriteString(fmt.Sprintf("### %s\n\n", file.Name))
	result.WriteString(fmt.Sprintf("- **路径**: `%s`\n", file.Path))
	result.WriteString(fmt.Sprintf("- **大小**: %d 字节\n", file.Size))
	result.WriteString(fmt.Sprintf("- **修改时间**: %s\n", file.ModTime.Format("2006-01-02 15:04:05")))

	// 检测语言类型
	language := f.detectLanguage(file.Name)
	result.WriteString(fmt.Sprintf("- **语言**: %s\n", language))

	// 估算token数量
	tokenCount := f.estimateTokens(file.Content)
	result.WriteString(fmt.Sprintf("- **Token数量**: %d\n", tokenCount))

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
		result.WriteString("#### 代码内容\n\n")
		result.WriteString(fmt.Sprintf("```%s\n", language))
		result.WriteString(file.Content)
		result.WriteString("\n```\n")
	} else {
		result.WriteString("#### 文件内容\n\n")
		result.WriteString("[二进制文件 - 内容未显示]\n")
	}

	return result.String()
}

// generateDirectoryStructure 生成目录结构
func (f *MarkdownFormatter) generateDirectoryStructure(data types.ContextData) string {
	var result strings.Builder
	result.WriteString("```\n")

	// 简化的目录结构表示
	for _, folder := range data.Folders {
		result.WriteString(fmt.Sprintf("%s/\n", folder.Name))
	}
	for _, file := range data.Files {
		result.WriteString(fmt.Sprintf("  %s\n", file.Name))
	}

	result.WriteString("```\n")
	return result.String()
}

// getFileTypes 获取文件类型统计
func (f *MarkdownFormatter) getFileTypes(data types.ContextData) map[string]int {
	fileTypes := make(map[string]int)
	for _, file := range data.Files {
		ext := f.getFileExtension(file.Name)
		if ext != "" {
			fileTypes[ext]++
		}
	}
	return fileTypes
}

// getFileExtension 获取文件扩展名
func (f *MarkdownFormatter) getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}

// detectLanguage 检测编程语言
func (f *MarkdownFormatter) detectLanguage(filename string) string {
	ext := f.getFileExtension(filename)
	switch ext {
	case "go":
		return "go"
	case "py":
		return "python"
	case "js":
		return "javascript"
	case "ts":
		return "typescript"
	case "java":
		return "java"
	case "cpp", "cc", "cxx":
		return "cpp"
	case "c":
		return "c"
	case "cs":
		return "csharp"
	case "php":
		return "php"
	case "rb":
		return "ruby"
	case "rs":
		return "rust"
	case "swift":
		return "swift"
	case "kt":
		return "kotlin"
	case "scala":
		return "scala"
	case "r":
		return "r"
	case "m":
		return "matlab"
	case "pl":
		return "perl"
	case "sh", "bash":
		return "bash"
	case "ps1":
		return "powershell"
	case "sql":
		return "sql"
	case "html":
		return "html"
	case "css":
		return "css"
	case "xml":
		return "xml"
	case "json":
		return "json"
	case "yaml", "yml":
		return "yaml"
	case "md":
		return "markdown"
	default:
		return "text"
	}
}

// getLanguages 获取语言列表
func (f *MarkdownFormatter) getLanguages(data types.ContextData) []string {
	languageMap := make(map[string]bool)
	for _, file := range data.Files {
		lang := f.detectLanguage(file.Name)
		if lang != "text" && lang != "" {
			languageMap[lang] = true
		}
	}
	
	languages := make([]string, 0, len(languageMap))
	for lang := range languageMap {
		languages = append(languages, lang)
	}
	return languages
}

// estimateTokens 估算token数量（简化版本）
func (f *MarkdownFormatter) estimateTokens(content string) int {
	// 简化的token估算：假设每个token大约4个字符
	return len(content) / 4
}
