// Package types 定义项目的核心类型和接口
package types

import (
	"fmt"
	"time"
)

// FileInfo 文件信息结构体
type FileInfo struct {
	Name     string    `yaml:"name"`
	Path     string    `yaml:"path"`
	Content  string    `yaml:"content"`
	Size     int64     `yaml:"size,omitempty"`
	ModTime  time.Time `yaml:"mod_time,omitempty"`
	IsDir    bool      `yaml:"is_dir,omitempty"`
	IsHidden bool      `yaml:"is_hidden,omitempty"`
	IsBinary bool      `yaml:"is_binary,omitempty"`
}

// FolderInfo 文件夹信息结构体
type FolderInfo struct {
	Name     string       `yaml:"name"`
	Path     string       `yaml:"path"`
	Files    []FileInfo   `yaml:"files"`
	Folders  []FolderInfo `yaml:"folders"`
	ModTime  time.Time    `yaml:"mod_time"`
	IsHidden bool         `yaml:"is_hidden"`
	Size     int64        `yaml:"size"`
	Count    int          `yaml:"count"`
}

// ContextData 上下文数据结构
type ContextData struct {
	Files       []FileInfo             `yaml:"files"`
	Folders     []FolderInfo           `yaml:"folders"`
	FileCount   int                    `yaml:"file_count"`
	FolderCount int                    `yaml:"folder_count"`
	TotalSize   int64                  `yaml:"total_size"`
	Metadata    map[string]interface{} `yaml:"metadata"`
}

// WalkResult 遍历结果
type WalkResult struct {
	Files       []FileInfo   `yaml:"files"`
	Folders     []FolderInfo `yaml:"folders"`
	FileCount   int          `yaml:"file_count"`
	FolderCount int          `yaml:"folder_count"`
	TotalSize   int64        `yaml:"total_size"`
	RootPath    string       `yaml:"root_path"`
	ScanDuration string      `yaml:"scan_duration"`
}

// Config 统一配置结构体
type Config struct {
	Formats       FormatsConfig       `yaml:"formats"`
	Fields        FieldsConfig        `yaml:"fields"`
	Filters       FiltersConfig       `yaml:"filters"`
	Output        OutputConfig        `yaml:"output"`
	UI            UIConfig            `yaml:"ui"`
	FileProcessing FileProcessingConfig `yaml:"file_processing"`
	Performance   PerformanceConfig   `yaml:"performance"`
	Logging       LoggingConfig       `yaml:"logging"`
	Security      SecurityConfig      `yaml:"security"`
	Git           GitIntegrationConfig `yaml:"git"`
}

// FormatsConfig 输出格式配置
type FormatsConfig struct {
	XML      XMLFormatConfig `yaml:"xml"`
	JSON     FormatConfig `yaml:"json"`
	TOML     FormatConfig `yaml:"toml"`
	Markdown FormatConfig `yaml:"markdown"`
}

// FormatConfig 单个格式配置
type FormatConfig struct {
	Enabled    bool                   `yaml:"enabled"`
	Structure  map[string]interface{} `yaml:"structure"`
	Fields     map[string]string      `yaml:"fields"`
	Template   string                 `yaml:"template"`
	Formatting map[string]interface{} `yaml:"formatting"`
	Encoding   string                 `yaml:"encoding"`
}

// XMLFormatConfig XML格式专用配置
type XMLFormatConfig struct {
	FormatConfig `yaml:",inline"`
	RootTag      string            `yaml:"root_tag"`
	FileTag      string            `yaml:"file_tag"`
	FolderTag    string            `yaml:"folder_tag"`
	FilesTag     string            `yaml:"files_tag"`
	Formatting   XMLFormattingConfig `yaml:"formatting"`
}

// XMLFormattingConfig XML格式化配置
type XMLFormattingConfig struct {
	Indent      string           `yaml:"indent"`
	Declaration bool             `yaml:"declaration"`
	Encoding    string           `yaml:"encoding"`
	ContentHandling XMLContentHandling `yaml:"content_handling"`
}

// XMLContentHandling XML内容处理方式
type XMLContentHandling string

const (
	// XMLContentEscaped 使用XML实体转义（默认）
	XMLContentEscaped XMLContentHandling = "escaped"
	// XMLContentCDATA 使用CDATA包装
	XMLContentCDATA XMLContentHandling = "cdata"
	// XMLContentRaw 保留原始格式（最小转义）
	XMLContentRaw XMLContentHandling = "raw"
)

// FieldsConfig 字段配置
type FieldsConfig struct {
	CustomNames map[string]string `yaml:"custom_names"`
	Filter      struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	} `yaml:"filter"`
	Processing struct {
		MaxLength      int  `yaml:"max_length"`
		AddLineNumbers bool `yaml:"add_line_numbers"`
		TrimWhitespace bool `yaml:"trim_whitespace"`
		CodeHighlight  bool `yaml:"code_highlight"`
	} `yaml:"processing"`
}

// FiltersConfig 文件过滤配置
type FiltersConfig struct {
	MaxFileSize     string   `yaml:"max_file_size"`
	ExcludePatterns []string `yaml:"exclude_patterns"`
	IncludePatterns []string `yaml:"include_patterns"`
	MaxDepth        int      `yaml:"max_depth"`
	FollowSymlinks  bool     `yaml:"follow_symlinks"`
	ExcludeBinary   bool     `yaml:"exclude_binary"`
}

// OutputConfig 输出配置
type OutputConfig struct {
	Format       string `yaml:"format"`
	FilePath     string `yaml:"file_path"`
	Encoding     string `yaml:"encoding"`
	DefaultFormat    string `yaml:"default_format"`
	OutputDir        string `yaml:"output_dir"`
	FilenameTemplate string `yaml:"filename_template"`
	TimestampFormat  string `yaml:"timestamp_format"`
	IncludeMetadata  bool   `yaml:"include_metadata"`
}

// UIConfig 界面配置
type UIConfig struct {
	Theme         string `yaml:"theme"`
	ShowProgress  bool   `yaml:"show_progress"`
	ShowSize      bool   `yaml:"show_size"`
	ShowDate      bool   `yaml:"show_date"`
	ShowPreview   bool   `yaml:"show_preview"`
	Selector struct {
		ShowHidden   bool `yaml:"show_hidden"`
		ShowSize     bool `yaml:"show_size"`
		ShowModified bool `yaml:"show_modified"`
	} `yaml:"selector"`
}

// SelectOptions 选择选项
type SelectOptions struct {
	// Recursive       bool // 已移除，使用max-depth控制递归
	IncludePatterns []string
	ExcludePatterns []string
	MaxDepth        int
	ShowHidden      bool
	SortBy          string
}

// WalkOptions 文件遍历选项
type WalkOptions struct {
	MaxDepth        int
	MaxFileSize     int64
	ExcludePatterns []string
	IncludePatterns []string
	FollowSymlinks  bool
	ShowHidden      bool
	ExcludeBinary   bool
	SelectedFiles   []string // 选中的具体文件路径，如果为空则使用模式匹配
	MultipleFiles   []string // 多个文件路径（-m参数）
	PatternFile     string   // 模式文件路径（-r参数）
}

// FileProcessingConfig 文件处理配置
type FileProcessingConfig struct {
	IncludeHidden   bool     `yaml:"include_hidden"`
	MaxFileSize     int64    `yaml:"max_file_size"`
	MaxDepth        int      `yaml:"max_depth"`
	ExcludePatterns []string `yaml:"exclude_patterns"`
	IncludePatterns []string `yaml:"include_patterns"`
	IncludeContent  bool     `yaml:"include_content"`
	IncludeHash     bool     `yaml:"include_hash"`
}

// PerformanceConfig 性能配置
type PerformanceConfig struct {
	MaxWorkers   int
	BufferSize   int
	CacheEnabled bool
	CacheSize    int
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

// CLIOptions 命令行选项
type CLIOptions struct {
	Format           string
	Output           string
	Config           string
	Exclude          []string
	Include          []string
	MaxDepth         int
	FollowSymlinks   bool
	OutputDir        string
	FilenameTemplate string
	ValidateConfig   bool
}

// AppError 应用错误类型
type AppError struct {
	Type    ErrorType
	Message string
	Cause   error
	Context map[string]interface{}
}

// ErrorType 错误类型枚举
type ErrorType int

const (
	ErrConfig ErrorType = iota
	ErrFileSystem
	ErrFormat
	ErrValidation
	ErrPermission
	ErrNetwork
	ErrUnknown
)

// String 返回错误类型的字符串表示
func (et ErrorType) String() string {
	switch et {
	case ErrConfig:
		return "ConfigError"
	case ErrFileSystem:
		return "FileSystemError"
	case ErrFormat:
		return "FormatError"
	case ErrValidation:
		return "ValidationError"
	case ErrPermission:
		return "PermissionError"
	case ErrNetwork:
		return "NetworkError"
	default:
		return "UnknownError"
	}
}

// Error 实现error接口
func (ae *AppError) Error() string {
	if ae.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", ae.Type, ae.Message, ae.Cause)
	}
	return fmt.Sprintf("%s: %s", ae.Type, ae.Message)
}

// Unwrap 返回底层错误
func (ae *AppError) Unwrap() error {
	return ae.Cause
}
