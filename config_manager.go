package main

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
)

// Config 统一配置结构体
type Config struct {
	Formats FormatsConfig `yaml:"formats" json:"formats" toml:"formats"`
	Fields  FieldsConfig  `yaml:"fields" json:"fields" toml:"fields"`
	Filters FiltersConfig `yaml:"filters" json:"filters" toml:"filters"`
	Output  OutputConfig  `yaml:"output" json:"output" toml:"output"`
	UI      UIConfig      `yaml:"ui" json:"ui" toml:"ui"`
}

// FormatsConfig 输出格式配置
type FormatsConfig struct {
	XML      FormatConfig `yaml:"xml" json:"xml" toml:"xml"`
	JSON     FormatConfig `yaml:"json" json:"json" toml:"json"`
	TOML     FormatConfig `yaml:"toml" json:"toml" toml:"toml"`
	Markdown FormatConfig `yaml:"markdown" json:"markdown" toml:"markdown"`
}

// FormatConfig 单个格式配置
type FormatConfig struct {
	Enabled   bool                   `yaml:"enabled" json:"enabled" toml:"enabled"`
	Structure map[string]interface{} `yaml:"structure" json:"structure" toml:"structure"`
	Fields    map[string]string      `yaml:"fields" json:"fields" toml:"fields"`
	Template  string                 `yaml:"template" json:"template" toml:"template"`
	Formatting map[string]interface{} `yaml:"formatting" json:"formatting" toml:"formatting"`
}

// FieldsConfig 字段配置
type FieldsConfig struct {
	CustomNames map[string]string `yaml:"custom_names" json:"custom_names" toml:"custom_names"`
	Filter      struct {
		Include []string `yaml:"include" json:"include" toml:"include"`
		Exclude []string `yaml:"exclude" json:"exclude" toml:"exclude"`
	} `yaml:"filter" json:"filter" toml:"filter"`
	Processing struct {
		MaxLength       int  `yaml:"max_length" json:"max_length" toml:"max_length"`
		AddLineNumbers  bool `yaml:"add_line_numbers" json:"add_line_numbers" toml:"add_line_numbers"`
		TrimWhitespace  bool `yaml:"trim_whitespace" json:"trim_whitespace" toml:"trim_whitespace"`
		CodeHighlight   bool `yaml:"code_highlight" json:"code_highlight" toml:"code_highlight"`
	} `yaml:"processing" json:"processing" toml:"processing"`
}

// FiltersConfig 过滤配置
type FiltersConfig struct {
	MaxFileSize     string   `yaml:"max_file_size" json:"max_file_size" toml:"max_file_size"`
	ExcludePatterns []string `yaml:"exclude_patterns" json:"exclude_patterns" toml:"exclude_patterns"`
	IncludePatterns []string `yaml:"include_patterns" json:"include_patterns" toml:"include_patterns"`
	MaxDepth        int      `yaml:"max_depth" json:"max_depth" toml:"max_depth"`
	FollowSymlinks  bool     `yaml:"follow_symlinks" json:"follow_symlinks" toml:"follow_symlinks"`
}

// OutputConfig 输出配置
type OutputConfig struct {
	DefaultFormat   string `yaml:"default_format" json:"default_format" toml:"default_format"`
	OutputDir       string `yaml:"output_dir" json:"output_dir" toml:"output_dir"`
	FilenameTemplate string `yaml:"filename_template" json:"filename_template" toml:"filename_template"`
	TimestampFormat string `yaml:"timestamp_format" json:"timestamp_format" toml:"timestamp_format"`
}

// UIConfig 界面配置
type UIConfig struct {
	Selector struct {
		ShowHidden   bool `yaml:"show_hidden" json:"show_hidden" toml:"show_hidden"`
		ShowSize     bool `yaml:"show_size" json:"show_size" toml:"show_size"`
		ShowModified bool `yaml:"show_modified" json:"show_modified" toml:"show_modified"`
	} `yaml:"selector" json:"selector" toml:"selector"`
	Autocomplete struct {
		Enabled        bool `yaml:"enabled" json:"enabled" toml:"enabled"`
		MinChars       int  `yaml:"min_chars" json:"min_chars" toml:"min_chars"`
		MaxSuggestions int  `yaml:"max_suggestions" json:"max_suggestions" toml:"max_suggestions"`
	} `yaml:"autocomplete" json:"autocomplete" toml:"autocomplete"`
}

// FileInfo 文件信息
type FileInfo struct {
	Name    string `yaml:"name" json:"name" toml:"name"`
	Path    string `yaml:"path" json:"path" toml:"path"`
	Content string `yaml:"content" json:"content" toml:"content"`
	Size    int64  `yaml:"size" json:"size" toml:"size"`
}

// FolderInfo 文件夹信息
type FolderInfo struct {
	Name  string     `yaml:"name" json:"name" toml:"name"`
	Path  string     `yaml:"path" json:"path" toml:"path"`
	Files []FileInfo `yaml:"files" json:"files" toml:"files"`
}

// ContextData 上下文数据
type ContextData struct {
	Files   []FileInfo   `yaml:"file" json:"file" toml:"file"`
	Folders []FolderInfo `yaml:"folder" json:"folder" toml:"folder"`
}

// ConfigManager 配置管理器
type ConfigManager struct {
	config     *Config
	mu         sync.RWMutex
	configPath string // 配置文件路径
}

// NewConfigManager 创建配置管理器
func NewConfigManager(configPath ...string) *ConfigManager {
	if len(configPath) > 0 && configPath[0] != "" {
		config, err := LoadConfig(configPath[0])
		if err != nil {
			// 如果加载失败，返回默认配置
			return &ConfigManager{
				config:     loadDefaultConfig(),
				configPath: configPath[0],
			}
		}
		return &ConfigManager{
			config:     config,
			configPath: configPath[0],
		}
	}
	return &ConfigManager{
		config: loadDefaultConfig(),
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config

	// 根据文件扩展名选择解析器
	switch ext := strings.ToLower(filepath.Ext(filename)); ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &config)
	case ".json":
		err = json.Unmarshal(data, &config)
	case ".toml":
		err = toml.Unmarshal(data, &config)
	default:
		return nil, fmt.Errorf("不支持的配置文件格式: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 设置默认值（配置已经包含默认值，无需额外设置）
	// setDefaults(&config)

	return &config, nil
}



// GetConfig 获取当前配置
func (cm *ConfigManager) GetConfig() *Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// Reload 重新加载配置
func (cm *ConfigManager) Reload() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	config, err := LoadConfig(cm.configPath)
	if err != nil {
		return err
	}

	cm.config = config
	return nil
}

// SaveConfig 保存配置到文件
func (cm *ConfigManager) SaveConfig(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	var data []byte
	var err error

	switch ext {
	case ".yaml", ".yml":
		data, err = yaml.Marshal(cm.config)
	case ".json":
		data, err = json.MarshalIndent(cm.config, "", "  ")
	case ".toml":
		data, err = toml.Marshal(cm.config)
	default:
		return fmt.Errorf("不支持的配置文件格式: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}



// GenerateOutput 生成输出内容
func (cm *ConfigManager) GenerateOutput(data ContextData, format string) (string, error) {
	switch strings.ToLower(format) {
	case "xml":
		return cm.generateXML(data)
	case "json":
		return cm.generateJSON(data)
	case "toml":
		return cm.generateTOML(data)
	case "markdown", "md":
		return cm.generateMarkdown(data)
	default:
		return "", fmt.Errorf("不支持的输出格式: %s", format)
	}
}

// generateXML 生成XML格式
func (cm *ConfigManager) generateXML(data ContextData) (string, error) {
	xmlData := struct {
		XMLName xml.Name     `xml:"context"`
		Files   []FileInfo   `xml:"file"`
		Folders []FolderInfo `xml:"folder"`
	}{
		Files:   data.Files,
		Folders: data.Folders,
	}

	output, err := xml.MarshalIndent(xmlData, "", "  ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(output), nil
}

// generateJSON 生成JSON格式
func (cm *ConfigManager) generateJSON(data ContextData) (string, error) {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// generateTOML 生成TOML格式
func (cm *ConfigManager) generateTOML(data ContextData) (string, error) {
	output, err := toml.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// generateMarkdown 生成Markdown格式
func (cm *ConfigManager) generateMarkdown(data ContextData) (string, error) {
	var sb strings.Builder

	// 添加文件部分
	for _, file := range data.Files {
		sb.WriteString(fmt.Sprintf("## 文件: %s\n\n", file.Path))
		if cm.config.Fields.Processing.CodeHighlight {
			// 这里可以添加代码高亮逻辑
			sb.WriteString("```\n")
		} else {
			sb.WriteString("```\n")
		}
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

// GetOutputFilename 生成输出文件名
func (cm *ConfigManager) GetOutputFilename(format string) string {
	template := cm.config.Output.FilenameTemplate
	if template == "" {
		template = "context_{{.timestamp}}.{{.extension}}"
	}

	timestamp := time.Now().Format(cm.config.Output.TimestampFormat)
	if timestamp == "" {
		timestamp = time.Now().Format("20060102_150405")
	}

	filename := strings.ReplaceAll(template, "{{.timestamp}}", timestamp)
	filename = strings.ReplaceAll(filename, "{{.extension}}", format)

	return filename
}

// loadDefaultConfig 加载默认配置
func loadDefaultConfig() *Config {
	return &Config{
		Formats: FormatsConfig{
			XML: FormatConfig{
				Enabled: true,
				Structure: map[string]interface{}{
					"root":  "context",
					"file":  "file",
					"folder": "folder",
					"files": "files",
				},
				Fields: map[string]string{
					"path":     "path",
					"content":  "content",
					"filename": "filename",
				},
			},
			JSON: FormatConfig{
				Enabled: true,
				Structure: map[string]interface{}{
					"file":  "file",
					"folder": "folder",
				},
				Fields: map[string]string{
					"path":     "path",
					"content":  "content",
					"filename": "filename",
				},
			},
			TOML: FormatConfig{
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
			Markdown: FormatConfig{
				Enabled: true,
				Structure: map[string]interface{}{
					"file_header":   "##",
					"folder_header": "###",
					"code_block":    "```",
				},
				Fields: map[string]string{
					"path":    "",
					"content": "",
				},
				Formatting: map[string]interface{}{
					"separator":     "\n\n",
					"add_toc":       false,
					"code_language": true,
				},
			},
		},
		Fields: FieldsConfig{
			CustomNames: map[string]string{
				"filepath":    "path",
				"filecontent": "content",
				"filename":    "name",
			},
			Filter: struct {
				Include []string `yaml:"include" json:"include" toml:"include"`
				Exclude []string `yaml:"exclude" json:"exclude" toml:"exclude"`
			}{
				Include: []string{},
				Exclude: []string{},
			},
			Processing: struct {
				MaxLength       int  `yaml:"max_length" json:"max_length" toml:"max_length"`
				AddLineNumbers  bool `yaml:"add_line_numbers" json:"add_line_numbers" toml:"add_line_numbers"`
				TrimWhitespace  bool `yaml:"trim_whitespace" json:"trim_whitespace" toml:"trim_whitespace"`
				CodeHighlight   bool `yaml:"code_highlight" json:"code_highlight" toml:"code_highlight"`
			}{
				MaxLength:       0,
				AddLineNumbers:  false,
				TrimWhitespace:  true,
				CodeHighlight:   false,
			},
		},
		Filters: FiltersConfig{
			MaxFileSize:     "10MB",
			ExcludePatterns: []string{"*.tmp", "*.log", "*.swp", ".*"},
			IncludePatterns: []string{},
			MaxDepth:        0,
			FollowSymlinks:  false,
		},
		Output: OutputConfig{
			DefaultFormat:    "xml",
			OutputDir:        "",
			FilenameTemplate: "context_{{.timestamp}}.{{.extension}}",
			TimestampFormat:  "20060102_150405",
		},
		UI: UIConfig{
			Selector: struct {
				ShowHidden   bool `yaml:"show_hidden" json:"show_hidden" toml:"show_hidden"`
				ShowSize     bool `yaml:"show_size" json:"show_size" toml:"show_size"`
				ShowModified bool `yaml:"show_modified" json:"show_modified" toml:"show_modified"`
			}{
				ShowHidden:   false,
				ShowSize:     true,
				ShowModified: false,
			},
			Autocomplete: struct {
				Enabled        bool `yaml:"enabled" json:"enabled" toml:"enabled"`
				MinChars       int  `yaml:"min_chars" json:"min_chars" toml:"min_chars"`
				MaxSuggestions int  `yaml:"max_suggestions" json:"max_suggestions" toml:"max_suggestions"`
			}{
				Enabled:        true,
				MinChars:       1,
				MaxSuggestions: 10,
			},
		},
	}
}