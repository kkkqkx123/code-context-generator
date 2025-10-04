// Package config 提供配置管理功能
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"code-context-generator/internal/env"
	"code-context-generator/pkg/constants"
	"code-context-generator/pkg/types"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"
)

// Manager 配置管理器接口
type Manager interface {
	Load(configPath string) error
	Get() *types.Config
	Validate() error
	Reload() error
	Save(configPath string, format string) error
	GetEnvOverrides() map[string]string
	GenerateOutput(data types.ContextData, format string) (string, error)
	GetOutputFilename(format string) string
}

// ConfigManager 配置管理器实现
type ConfigManager struct {
	config     *types.Config
	mu         sync.RWMutex
	configPath string
}

// NewManager 创建新的配置管理器
func NewManager() Manager {
	return &ConfigManager{
		config: GetDefaultConfig(),
	}
}

// Load 加载配置文件
func (cm *ConfigManager) Load(configPath string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if configPath == "" {
		configPath = constants.DefaultConfigFile
	}

	// 首先加载.env文件（如果存在）
	if err := env.LoadEnv(""); err != nil {
		// 如果.env文件加载失败，记录警告但不中断程序
		fmt.Printf("警告: 加载.env文件失败: %v\n", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 如果文件不存在，使用默认配置，不再自动创建配置文件
		cm.config = GetDefaultConfig()
		cm.configPath = configPath
		return nil
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}

	// 应用环境变量覆盖
	cm.applyEnvOverrides(config)

	// 应用格式特定的配置覆盖（基于配置文件名）
	applyFormatSpecificConfig(config, configPath)

	cm.config = config
	cm.configPath = configPath
	return nil
}

// Get 获取当前配置
func (cm *ConfigManager) Get() *types.Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// Validate 验证配置
func (cm *ConfigManager) Validate() error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if cm.config == nil {
		return fmt.Errorf("配置为空")
	}

	// 验证格式配置
	formats := []string{constants.FormatXML, constants.FormatJSON, constants.FormatTOML, constants.FormatMarkdown}
	hasEnabled := false
	for _, format := range formats {
		if cm.isFormatEnabled(format) {
			hasEnabled = true
			break
		}
	}

	if !hasEnabled {
		return fmt.Errorf("至少需要启用一种输出格式")
	}

	// 验证输出配置
	if cm.config.Output.FilenameTemplate == "" {
		return fmt.Errorf("文件名模板不能为空")
	}

	// 验证时间格式
	if _, err := time.Parse(cm.config.Output.TimestampFormat, time.Now().Format(cm.config.Output.TimestampFormat)); err != nil {
		return fmt.Errorf("时间格式无效: %w", err)
	}

	return nil
}

// Reload 重新加载配置
func (cm *ConfigManager) Reload() error {
	if cm.configPath == "" {
		return fmt.Errorf("配置文件路径未设置")
	}
	return cm.Load(cm.configPath)
}

// Save 保存配置到文件
func (cm *ConfigManager) Save(configPath string, format string) error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if cm.config == nil {
		return fmt.Errorf("配置为空")
	}

	switch strings.ToLower(format) {
	case "yaml", "yml":
		return cm.saveYAML(configPath)
	case "json":
		return cm.saveJSON(configPath)
	case "toml":
		return cm.saveTOML(configPath)
	default:
		return fmt.Errorf("不支持的格式: %s", format)
	}
}

// GetEnvOverrides 获取环境变量覆盖
func (cm *ConfigManager) GetEnvOverrides() map[string]string {
	envVars := env.GetAllEnvVars()
	overrides := make(map[string]string)

	// 将环境变量名映射到配置字段名
	mapping := map[string]string{
		env.EnvDefaultFormat:    "default_format",
		env.EnvOutputDir:        "output_dir",
		env.EnvFilenameTemplate: "filename_template",
		env.EnvTimestampFormat:  "timestamp_format",
		env.EnvMaxFileSize:      "max_file_size",
		env.EnvMaxDepth:         "max_depth",
		// env.EnvRecursive:        "recursive", // 已移除recursive参数
		env.EnvIncludeHidden:    "include_hidden",
		env.EnvFollowSymlinks:   "follow_symlinks",
		env.EnvExcludeBinary:    "exclude_binary",
		env.EnvExcludePatterns:  "exclude_patterns",
		env.EnvEncoding:         "encoding",
		env.EnvIncludeMetadata:  "include_metadata",
		// 安全扫描配置
		env.EnvSecurityEnabled:           "security_enabled",
		env.EnvSecurityFailOnCritical:    "security_fail_on_critical",
		env.EnvSecurityScanLevel:         "security_scan_level",
		env.EnvSecurityReportFormat:      "security_report_format",
		env.EnvSecurityDetectCredentials: "security_detect_credentials",
		env.EnvSecurityDetectSQLInjection: "security_detect_sql_injection",
		env.EnvSecurityDetectXSS:         "security_detect_xss",
		env.EnvSecurityDetectPathTraversal: "security_detect_path_traversal",
		env.EnvSecurityDetectQuality:     "security_detect_quality",
	}

	for envKey, fieldName := range mapping {
		if value, exists := envVars[envKey]; exists && value != "" {
			overrides[fieldName] = value
		}
	}

	return overrides
}

// applyEnvOverrides 应用环境变量覆盖到配置
func (cm *ConfigManager) applyEnvOverrides(config *types.Config) {
	// 应用输出格式覆盖
	if format := env.GetDefaultFormat(); format != "" {
		config.Output.DefaultFormat = format
	}

	// 应用输出目录覆盖
	if outputDir := env.GetOutputDir(); outputDir != "" {
		config.Output.OutputDir = outputDir
	}

	// 应用安全扫描配置覆盖
	config.Security.Enabled = env.GetSecurityEnabled()
	config.Security.FailOnCritical = env.GetSecurityFailOnCritical()
	config.Security.ScanLevel = env.GetSecurityScanLevel()
	config.Security.ReportFormat = env.GetSecurityReportFormat()
	config.Security.Detectors.Credentials = env.GetSecurityDetectCredentials()
	config.Security.Detectors.SQLInjection = env.GetSecurityDetectSQLInjection()
	config.Security.Detectors.XSS = env.GetSecurityDetectXSS()
	config.Security.Detectors.PathTraversal = env.GetSecurityDetectPathTraversal()
	config.Security.Detectors.Quality = env.GetSecurityDetectQuality()

	// 应用文件名模板覆盖
	if filenameTemplate := env.GetFilenameTemplate(); filenameTemplate != "" {
		config.Output.FilenameTemplate = filenameTemplate
	}

	// 应用时间戳格式覆盖
	if timestampFormat := env.GetTimestampFormat(); timestampFormat != "" {
		config.Output.TimestampFormat = timestampFormat
	}

	// 应用最大文件大小覆盖
	if maxFileSize := env.GetMaxFileSize(); maxFileSize != "" {
		config.Filters.MaxFileSize = maxFileSize
	}

	// 应用最大深度覆盖
	config.Filters.MaxDepth = env.GetMaxDepth()

	// 应用排除模式覆盖
	if excludePatterns := env.GetExcludePatterns(); excludePatterns != "" {
		config.Filters.ExcludePatterns = strings.Split(excludePatterns, ",")
	}

	// 应用跟随符号链接覆盖
	config.Filters.FollowSymlinks = env.GetFollowSymlinks()

	// 应用排除二进制文件覆盖
	config.Filters.ExcludeBinary = env.GetExcludeBinary()
	
	// 应用编码覆盖
	if encoding := env.GetEncoding(); encoding != "" {
		config.Output.Encoding = encoding
	}
	
	// 应用元信息开关覆盖
	config.Output.IncludeMetadata = env.GetIncludeMetadata()
	
	// 注意：recursive 环境变量已被移除，使用 max_depth 控制递归行为
}

// GenerateOutput 生成输出内容
func (cm *ConfigManager) GenerateOutput(data types.ContextData, format string) (string, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	switch strings.ToLower(format) {
	case constants.FormatXML:
		return cm.generateXML(data)
	case constants.FormatJSON:
		return cm.generateJSON(data)
	case constants.FormatTOML:
		return cm.generateTOML(data)
	case constants.FormatMarkdown:
		return cm.generateMarkdown(data)
	default:
		return "", fmt.Errorf("不支持的格式: %s", format)
	}
}

// GetOutputFilename 生成输出文件名
func (cm *ConfigManager) GetOutputFilename(format string) string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	template := cm.config.Output.FilenameTemplate
	if template == "" {
		template = constants.DefaultFilenameTemplate
	}

	timestamp := time.Now().Format(cm.config.Output.TimestampFormat)
	if timestamp == "" {
		timestamp = time.Now().Format(constants.DefaultTimestampFormat)
	}

	filename := strings.ReplaceAll(template, "{{.timestamp}}", timestamp)
	filename = strings.ReplaceAll(filename, "{{.extension}}", format)

	return filename
}

// saveConfig 内部保存配置（不加锁）
func (cm *ConfigManager) saveConfig(configPath string, format string) error {
	if cm.config == nil {
		return fmt.Errorf("配置为空")
	}

	switch strings.ToLower(format) {
	case "yaml", "yml":
		return cm.saveYAML(configPath)
	case "json":
		return cm.saveJSON(configPath)
	case "toml":
		return cm.saveTOML(configPath)
	default:
		return fmt.Errorf("不支持的格式: %s", format)
	}
}

// 辅助方法

func (cm *ConfigManager) isFormatEnabled(format string) bool {
	switch format {
	case constants.FormatXML:
		return cm.config.Formats.XML.Enabled
	case constants.FormatJSON:
		return cm.config.Formats.JSON.Enabled
	case constants.FormatTOML:
		return cm.config.Formats.TOML.Enabled
	case constants.FormatMarkdown:
		return cm.config.Formats.Markdown.Enabled
	default:
		return false
	}
}

func (cm *ConfigManager) saveYAML(configPath string) error {
	data, err := yaml.Marshal(cm.config)
	if err != nil {
		return fmt.Errorf("YAML序列化失败: %w", err)
	}
	return os.WriteFile(configPath, data, 0644)
}

func (cm *ConfigManager) saveJSON(configPath string) error {
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %w", err)
	}
	return os.WriteFile(configPath, data, 0644)
}

func (cm *ConfigManager) saveTOML(configPath string) error {
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(cm.config); err != nil {
		return fmt.Errorf("TOML编码失败: %w", err)
	}
	return nil
}

func (cm *ConfigManager) generateXML(data types.ContextData) (string, error) {
	// 获取XML配置
	xmlConfig := cm.config.Formats.XML

	var sb strings.Builder

	// 添加XML声明
	if xmlConfig.Formatting.Declaration {
		encoding := xmlConfig.Formatting.Encoding
		if encoding == "" {
			encoding = "UTF-8"
		}
		sb.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="%s"?>`, encoding))
		sb.WriteString("\n")
	}

	// 生成根元素
	rootTag := xmlConfig.RootTag
	if rootTag == "" {
		rootTag = "context"
	}

	sb.WriteString(fmt.Sprintf("<%s>\n", rootTag))

	// 生成元数据
	if data.Metadata != nil {
		sb.WriteString("  <metadata>\n")
		for key, value := range data.Metadata {
			sb.WriteString(fmt.Sprintf("    <%s>%v</%s>\n", key, value, key))
		}
		sb.WriteString("  </metadata>\n")
	}

	// 生成文件部分
	if len(data.Files) > 0 {
		filesTag := xmlConfig.FilesTag
		if filesTag == "" {
			filesTag = "files"
		}
		sb.WriteString(fmt.Sprintf("  <%s>\n", filesTag))

		fileTag := xmlConfig.FileTag
		if fileTag == "" {
			fileTag = "file"
		}

		for _, file := range data.Files {
			sb.WriteString(fmt.Sprintf("    <%s>\n", fileTag))

			// 获取字段映射
			pathField := xmlConfig.Fields["path"]
			if pathField == "" {
				pathField = "path"
			}
			sb.WriteString(fmt.Sprintf("      <%s>%s</%s>\n", pathField, escapeXML(file.Path), pathField))

			if file.Content != "" {
				contentField := xmlConfig.Fields["content"]
				if contentField == "" {
					contentField = "content"
				}
				sb.WriteString(fmt.Sprintf("      <%s><![CDATA[%s]]></%s>\n", contentField, file.Content, contentField))
			}

			sb.WriteString(fmt.Sprintf("    </%s>\n", fileTag))
		}
		sb.WriteString(fmt.Sprintf("  </%s>\n", filesTag))
	}

	// 生成文件夹部分
	if len(data.Folders) > 0 {
		folderTag := xmlConfig.FolderTag
		if folderTag == "" {
			folderTag = "folder"
		}

		for _, folder := range data.Folders {
			sb.WriteString(fmt.Sprintf("  <%s>\n", folderTag))

			pathField := xmlConfig.Fields["path"]
			if pathField == "" {
				pathField = "path"
			}
			sb.WriteString(fmt.Sprintf("    <%s>%s</%s>\n", pathField, escapeXML(folder.Path), pathField))

			if len(folder.Files) > 0 {
				filesTag := xmlConfig.FilesTag
				if filesTag == "" {
					filesTag = "files"
				}
				sb.WriteString(fmt.Sprintf("    <%s>\n", filesTag))

				fileTag := xmlConfig.FileTag
				if fileTag == "" {
					fileTag = "file"
				}

				for _, file := range folder.Files {
					sb.WriteString(fmt.Sprintf("      <%s>\n", fileTag))

					filenameField := xmlConfig.Fields["filename"]
					if filenameField == "" {
						filenameField = "filename"
					}
					sb.WriteString(fmt.Sprintf("        <%s>%s</%s>\n", filenameField, escapeXML(file.Name), filenameField))

					if file.Content != "" {
						contentField := xmlConfig.Fields["content"]
						if contentField == "" {
							contentField = "content"
						}
						sb.WriteString(fmt.Sprintf("        <%s><![CDATA[%s]]></%s>\n", contentField, file.Content, contentField))
					}

					sb.WriteString(fmt.Sprintf("      </%s>\n", fileTag))
				}
				sb.WriteString(fmt.Sprintf("    </%s>\n", filesTag))
			}

			sb.WriteString(fmt.Sprintf("  </%s>\n", folderTag))
		}
	}

	sb.WriteString(fmt.Sprintf("</%s>", rootTag))

	return sb.String(), nil
}

// escapeXML 转义XML特殊字符
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

func (cm *ConfigManager) generateJSON(data types.ContextData) (string, error) {
	// 实现JSON生成逻辑
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON生成失败: %w", err)
	}
	return string(output), nil
}

func (cm *ConfigManager) generateTOML(data types.ContextData) (string, error) {
	// 实现TOML生成逻辑
	var buf strings.Builder
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return "", fmt.Errorf("TOML生成失败: %w", err)
	}
	return buf.String(), nil
}

func (cm *ConfigManager) generateMarkdown(data types.ContextData) (string, error) {
	// 实现Markdown生成逻辑
	var sb strings.Builder

	// 添加文件部分
	for _, file := range data.Files {
		sb.WriteString(fmt.Sprintf("## 文件: %s\n\n", file.Path))
		sb.WriteString("```\n")
		sb.WriteString(file.Content)
		sb.WriteString("\n```\n\n")
	}

	// 添加文件夹部分
	for _, folder := range data.Folders {
		sb.WriteString(fmt.Sprintf("### 文件夹: %s\n\n", folder.Path))
		for _, file := range folder.Files {
			sb.WriteString(fmt.Sprintf("#### 文件: %s\n\n", file.Name))
			sb.WriteString("```\n")
			sb.WriteString(file.Content)
			sb.WriteString("\n```\n\n")
		}
	}

	return sb.String(), nil
}

// LoadConfig 从文件加载配置（辅助函数）
func LoadConfig(configPath string) (*types.Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(configPath))
	var config types.Config

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("YAML解析失败: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("JSON解析失败: %w", err)
		}
	case ".toml":
		if _, err := toml.Decode(string(data), &config); err != nil {
			return nil, fmt.Errorf("TOML解析失败: %w", err)
		}
	case ".xml":
		// 对于XML配置文件，我们需要特殊处理，因为XML结构通常与Go结构体不匹配
		// 这里我们使用一个简化的XML解析，或者返回错误提示用户使用支持的格式
		return nil, fmt.Errorf("XML配置文件格式暂不支持，请使用YAML、JSON或TOML格式")
	default:
		return nil, fmt.Errorf("不支持的配置文件格式: %s", ext)
	}

	// 根据配置文件名应用相应的格式配置
	applyFormatSpecificConfig(&config, configPath)

	return &config, nil
}

// applyFormatSpecificConfig 根据配置文件名应用相应的格式配置
func applyFormatSpecificConfig(config *types.Config, configPath string) {
	// 将配置文件名转换为小写以便匹配
	configName := strings.ToLower(filepath.Base(configPath))

	// 根据配置文件名中是否包含特定格式名称来应用相应的格式配置
	if strings.Contains(configName, "xml") {
		// 如果配置文件名包含xml，应用XML格式配置
		if config.Formats.XML.Enabled {
			config.Output.DefaultFormat = "xml"
		}
	} else if strings.Contains(configName, "json") {
		// 如果配置文件名包含json，应用JSON格式配置
		if config.Formats.JSON.Enabled {
			config.Output.DefaultFormat = "json"
		}
	} else if strings.Contains(configName, "toml") {
		// 如果配置文件名包含toml，应用TOML格式配置
		if config.Formats.TOML.Enabled {
			config.Output.DefaultFormat = "toml"
		}
	} else if strings.Contains(configName, "markdown") || strings.Contains(configName, "md") {
		// 如果配置文件名包含markdown或md，应用Markdown格式配置
		if config.Formats.Markdown.Enabled {
			config.Output.DefaultFormat = "markdown"
		}
	}
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *types.Config {
	return &types.Config{
		FileProcessing: types.FileProcessingConfig{
			IncludeHidden:  false,
			IncludeContent: true, // 默认包含文件内容
			IncludeHash:    false,
		},
		Security: types.SecurityConfig{
			Enabled:     false, // 默认禁用，仅由gitignore规则控制
			ScanLevel:   "standard",
			FailOnCritical: false,
			ReportFormat: "text",
			Detectors: types.DetectorConfig{
				Credentials:    false, // 默认禁用
				SQLInjection: false, // 默认禁用
				XSS:          false, // 默认禁用
				PathTraversal: false, // 默认禁用
				Quality:      false, // 默认禁用
			},
		},
		Formats: types.FormatsConfig{
			XML: types.XMLFormatConfig{
				FormatConfig: types.FormatConfig{
					Enabled: true,
					Fields: map[string]string{
						"path":     "path",
						"content":  "content",
						"filename": "filename",
					},
				},
				RootTag:   "context",
				FileTag:   "file",
				FilesTag:  "files",
				FolderTag: "folder",
				Formatting: types.XMLFormattingConfig{
					Indent:      "  ",
					Declaration: true,
					Encoding:    "UTF-8",
				},
			},
			JSON: types.FormatConfig{
				Enabled: true,
				Structure: map[string]interface{}{
					"file":   "file",
					"folder": "folder",
				},
				Fields: map[string]string{
					"path":     "path",
					"content":  "content",
					"filename": "filename",
				},
			},
			TOML: types.FormatConfig{
				Enabled: true,
				Structure: map[string]interface{}{
					"file_section":   "file",
					"folder_section": "folder",
				},
				Fields: map[string]string{
					"path":     "path",
					"content":  "content",
					"filename": "filename",
				},
			},
			Markdown: types.FormatConfig{
				Enabled: true,
				Structure: map[string]interface{}{
					"file_header":   "##",
					"folder_header": "###",
					"code_block":    "```",
				},
				Formatting: map[string]interface{}{
					"separator":     "\n\n",
					"add_toc":       false,
					"code_language": true,
				},
			},
		},
		Fields: types.FieldsConfig{
			CustomNames: map[string]string{
				"filepath":    "path",
				"filecontent": "content",
				"filename":    "name",
			},
			Filter: struct {
				Include []string `yaml:"include"`
				Exclude []string `yaml:"exclude"`
			}{
				Include: []string{},
				Exclude: []string{},
			},
			Processing: struct {
				MaxLength      int  `yaml:"max_length"`
				AddLineNumbers bool `yaml:"add_line_numbers"`
				TrimWhitespace bool `yaml:"trim_whitespace"`
				CodeHighlight  bool `yaml:"code_highlight"`
			}{
				MaxLength:      0,
				AddLineNumbers: false,
				TrimWhitespace: true,
				CodeHighlight:  false,
			},
		},
		Filters: types.FiltersConfig{
			MaxFileSize:     "10MB",
			ExcludePatterns: []string{".git", "node_modules", "*.exe", "*.dll"},
			IncludePatterns: []string{},
			MaxDepth:        0,
			FollowSymlinks:  false,
			ExcludeBinary:   false,
		},
		Output: types.OutputConfig{
			Format:           "json",
			OutputDir:        "./output",
			Encoding:         "utf-8",
			DefaultFormat:    "xml",
			FilenameTemplate: "context_{{.timestamp}}.{{.extension}}",
			IncludeMetadata:  false,
		},
	}
}
