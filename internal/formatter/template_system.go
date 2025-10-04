package formatter

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"code-context-generator/pkg/types"
)

// TemplateSystem 模板系统
type TemplateSystem struct {
	config *types.Config
}

// TemplateData 模板数据结构
type TemplateData struct {
	Project     ProjectData
	Generation  GenerationData
	Statistics  StatisticsData
	Custom      map[string]interface{}
}

// ProjectData 项目数据
type ProjectData struct {
	Name        string
	Path        string
	Description string
	Languages   []string
}

// GenerationData 生成数据
type GenerationData struct {
	Timestamp   time.Time
	Tool        string
	Version     string
	Command     string
}

// StatisticsData 统计数据
type StatisticsData struct {
	FileCount   int
	FolderCount int
	TotalSize   int64
	TotalTokens int
}

// NewTemplateSystem 创建模板系统
func NewTemplateSystem(config *types.Config) *TemplateSystem {
	return &TemplateSystem{
		config: config,
	}
}

// ProcessTemplate 处理模板
func (t *TemplateSystem) ProcessTemplate(templateStr string, data TemplateData) (string, error) {
	// 创建模板
	tmpl, err := template.New("ai_template").Funcs(t.getTemplateFuncs()).Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %w", err)
	}

	// 执行模板
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("执行模板失败: %w", err)
	}

	return buf.String(), nil
}

// getTemplateFuncs 获取模板函数
func (t *TemplateSystem) getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatSize":     formatSize,
		"formatNumber":   formatNumber,
		"formatList":     formatList,
		"formatDate":     formatDate,
		"escapeXML":      escapeXML,
		"escapeJSON":     escapeJSON,
		"truncate":       truncate,
		"wordCount":      wordCount,
		"lineCount":      lineCount,
		"join":           strings.Join,
		"split":          strings.Split,
		"replace":        strings.Replace,
		"lower":          strings.ToLower,
		"upper":          strings.ToUpper,
		"title":          strings.Title,
		"trim":           strings.TrimSpace,
	}
}

// CreateDefaultTemplateData 创建默认模板数据
func (t *TemplateSystem) CreateDefaultTemplateData(fileCount, folderCount int, totalSize int64, languages []string) TemplateData {
	return TemplateData{
		Project: ProjectData{
			Name:        t.getProjectName(),
			Path:        t.getProjectPath(),
			Description: "Code repository analysis",
			Languages:   languages,
		},
		Generation: GenerationData{
			Timestamp: time.Now(),
			Tool:      "code-context-generator",
			Version:   "1.0.0",
			Command:   t.getCommand(),
		},
		Statistics: StatisticsData{
			FileCount:   fileCount,
			FolderCount: folderCount,
			TotalSize:   totalSize,
			TotalTokens: 0, // 将在后续集成token计数
		},
		Custom: make(map[string]interface{}),
	}
}

// getProjectName 获取项目名称
func (t *TemplateSystem) getProjectName() string {
	// 可以从配置或环境变量获取
	if t.config != nil && t.config.Output.FilenameTemplate != "" {
		return t.config.Output.FilenameTemplate
	}
	return "project"
}

// getProjectPath 获取项目路径
func (t *TemplateSystem) getProjectPath() string {
	// 实现获取项目路径的逻辑
	return "."
}

// getCommand 获取生成命令
func (t *TemplateSystem) getCommand() string {
	// 实现获取命令的逻辑
	return "code-context-generator"
}

// formatSize 格式化文件大小
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatNumber 格式化数字
func formatNumber(num int) string {
	return fmt.Sprintf("%d", num)
}

// formatList 格式化列表
func formatList(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return strings.Join(items, ", ")
}

// formatDate 格式化日期
func formatDate(t time.Time) string {
	return t.Format(time.RFC3339)
}

// escapeXML 转义XML字符
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// escapeJSON 转义JSON字符
func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// truncate 截断字符串
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	if length <= 3 {
		return s[:length]
	}
	return s[:length-3] + "..."
}

// wordCount 计算单词数
func wordCount(s string) int {
	return len(strings.Fields(s))
}

// lineCount 计算行数
func lineCount(s string) int {
	return strings.Count(s, "\n") + 1
}