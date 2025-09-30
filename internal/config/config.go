// Package config 提供配置管理功能
package config

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"
	"code-context-generator/pkg/constants"
	"code-context-generator/pkg/types"
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

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 如果文件不存在，创建默认配置
		return cm.Save(configPath, "yaml")
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}

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
	overrides := make(map[string]string)
	
	// 获取格式相关的环境变量
	if format := os.Getenv(constants.EnvPrefix + "DEFAULT_FORMAT"); format != "" {
		overrides["default_format"] = format
	}
	
	// 获取输出相关的环境变量
	if outputDir := os.Getenv(constants.EnvPrefix + "OUTPUT_DIR"); outputDir != "" {
		overrides["output_dir"] = outputDir
	}
	
	// 获取过滤相关的环境变量
	if maxDepth := os.Getenv(constants.EnvPrefix + "MAX_DEPTH"); maxDepth != "" {
		overrides["max_depth"] = maxDepth
	}
	
	return overrides
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
	// 实现XML生成逻辑
	output, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("XML生成失败: %w", err)
	}
	return string(output), nil
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
	default:
		return nil, fmt.Errorf("不支持的配置文件格式: %s", ext)
	}

	return &config, nil
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *types.Config {
	return &types.Config{
		Formats: types.FormatsConfig{
			XML: types.FormatConfig{
				Enabled: true,
				Structure: map[string]interface{}{
					"root":  "context",
					"file":  "file",
					"files": "files",
					"folder": "folder",
				},
				Fields: map[string]string{
					"path":     "path",
					"content":  "content",
					"filename": "filename",
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
					"separator": "\n\n",
					"add_toc":     false,
					"code_language": true,
				},
			},
		},
		Fields: types.FieldsConfig{
			CustomNames: map[string]string{
				"filepath":   "path",
				"filecontent": "content",
				"filename":   "name",
			},
			Filter: struct {
				Include []string `yaml:"include" json:"include" toml:"include"`
				Exclude []string `yaml:"exclude" json:"exclude" toml:"exclude"`
			}{
				Include: []string{},
				Exclude: []string{},
			},
			Processing: struct {
				MaxLength      int  `yaml:"max_length" json:"max_length" toml:"max_length"`
				AddLineNumbers bool `yaml:"add_line_numbers" json:"add_line_numbers" toml:"add_line_numbers"`
				TrimWhitespace bool `yaml:"trim_whitespace" json:"trim_whitespace" toml:"trim_whitespace"`
				CodeHighlight  bool `yaml:"code_highlight" json:"code_highlight" toml:"code_highlight"`
			}{
				MaxLength:      0,
				AddLineNumbers:   false,
				TrimWhitespace: true,
				CodeHighlight:  false,
			},
		},
		Filters: types.FiltersConfig{
			MaxFileSize:     "10MB",
			ExcludePatterns: constants.DefaultExcludePatterns,
			IncludePatterns: []string{},
			MaxDepth:        constants.DefaultMaxDepth,
			FollowSymlinks:  false,
		},
		Output: types.OutputConfig{
			DefaultFormat:    constants.DefaultFormat,
			OutputDir:        constants.DefaultOutputDir,
			FilenameTemplate: constants.DefaultFilenameTemplate,
			TimestampFormat:  constants.DefaultTimestampFormat,
		},
		UI: types.UIConfig{
			Selector: struct {
				ShowHidden   bool `yaml:"show_hidden" json:"show_hidden" toml:"show_hidden"`
				ShowSize     bool `yaml:"show_size" json:"show_size" toml:"show_size"`
				ShowModified bool `yaml:"show_modified" json:"show_modified" toml:"show_modified"`
			}{
				ShowHidden:   constants.DefaultShowHidden,
				ShowSize:     constants.DefaultShowSize,
				ShowModified: constants.DefaultShowModified,
			},
			Autocomplete: struct {
				Enabled        bool `yaml:"enabled" json:"enabled" toml:"enabled"`
				MinChars       int  `yaml:"min_chars" json:"min_chars" toml:"min_chars"`
				MaxSuggestions int  `yaml:"max_suggestions" json:"max_suggestions" toml:"max_suggestions"`
			}{
				Enabled:        true,
				MinChars:       constants.DefaultMinChars,
				MaxSuggestions: constants.DefaultMaxSuggestions,
			},
		},
	}
}