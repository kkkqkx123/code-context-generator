// Package types 定义项目的核心类型和接口
package types

import (
	"fmt"
	"time"
)

// FileInfo 文件信息结构体
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

// FolderInfo 文件夹信息结构体
type FolderInfo struct {
	Name     string       `yaml:"name" json:"name" toml:"name"`
	Path     string       `yaml:"path" json:"path" toml:"path"`
	Files    []FileInfo   `yaml:"files" json:"files" toml:"files"`
	Folders  []FolderInfo `yaml:"folders" json:"folders" toml:"folders"`
	ModTime  time.Time    `yaml:"mod_time" json:"mod_time" toml:"mod_time"`
	IsHidden bool         `yaml:"is_hidden" json:"is_hidden" toml:"is_hidden"`
	Size     int64        `yaml:"size" json:"size" toml:"size"`
	Count    int          `yaml:"count" json:"count" toml:"count"`
}

// ContextData 上下文数据结构
type ContextData struct {
	Files       []FileInfo             `yaml:"files" json:"files" toml:"files"`
	Folders     []FolderInfo           `yaml:"folders" json:"folders" toml:"folders"`
	FileCount   int                    `yaml:"file_count" json:"file_count" toml:"file_count"`
	FolderCount int                    `yaml:"folder_count" json:"folder_count" toml:"folder_count"`
	TotalSize   int64                  `yaml:"total_size" json:"total_size" toml:"total_size"`
	Metadata    map[string]interface{} `yaml:"metadata" json:"metadata" toml:"metadata"`
}

// WalkResult 遍历结果
type WalkResult struct {
	Files       []FileInfo   `yaml:"files" json:"files" toml:"files"`
	Folders     []FolderInfo `yaml:"folders" json:"folders" toml:"folders"`
	FileCount   int          `yaml:"file_count" json:"file_count" toml:"file_count"`
	FolderCount int          `yaml:"folder_count" json:"folder_count" toml:"folder_count"`
	TotalSize   int64        `yaml:"total_size" json:"total_size" toml:"total_size"`
	RootPath    string       `yaml:"root_path" json:"root_path" toml:"root_path"`
	ScanDuration string      `yaml:"scan_duration" json:"scan_duration" toml:"scan_duration"`
}

// Config 统一配置结构体
type Config struct {
	Formats       FormatsConfig       `yaml:"formats" json:"formats" toml:"formats"`
	Fields        FieldsConfig        `yaml:"fields" json:"fields" toml:"fields"`
	Filters       FiltersConfig       `yaml:"filters" json:"filters" toml:"filters"`
	Output        OutputConfig        `yaml:"output" json:"output" toml:"output"`
	UI            UIConfig            `yaml:"ui" json:"ui" toml:"ui"`
	FileProcessing FileProcessingConfig `yaml:"file_processing" json:"file_processing" toml:"file_processing"`
	Performance   PerformanceConfig   `yaml:"performance" json:"performance" toml:"performance"`
	Logging       LoggingConfig       `yaml:"logging" json:"logging" toml:"logging"`
}

// FormatsConfig 输出格式配置
type FormatsConfig struct {
	XML      XMLFormatConfig `yaml:"xml" json:"xml" toml:"xml"`
	JSON     FormatConfig `yaml:"json" json:"json" toml:"json"`
	TOML     FormatConfig `yaml:"toml" json:"toml" toml:"toml"`
	Markdown FormatConfig `yaml:"markdown" json:"markdown" toml:"markdown"`
}

// FormatConfig 单个格式配置
type FormatConfig struct {
	Enabled    bool                   `yaml:"enabled" json:"enabled" toml:"enabled"`
	Structure  map[string]interface{} `yaml:"structure" json:"structure" toml:"structure"`
	Fields     map[string]string      `yaml:"fields" json:"fields" toml:"fields"`
	Template   string                 `yaml:"template" json:"template" toml:"template"`
	Formatting map[string]interface{} `yaml:"formatting" json:"formatting" toml:"formatting"`
}

// XMLFormatConfig XML格式专用配置
type XMLFormatConfig struct {
	Enabled    bool              `yaml:"enabled" json:"enabled" toml:"enabled"`
	RootTag    string            `yaml:"root_tag" json:"root_tag" toml:"root_tag"`
	FileTag    string            `yaml:"file_tag" json:"file_tag" toml:"file_tag"`
	FolderTag  string            `yaml:"folder_tag" json:"folder_tag" toml:"folder_tag"`
	FilesTag   string            `yaml:"files_tag" json:"files_tag" toml:"files_tag"`
	Fields     map[string]string   `yaml:"fields" json:"fields" toml:"fields"`
	Template   string            `yaml:"template" json:"template" toml:"template"`
	Formatting XMLFormattingConfig `yaml:"formatting" json:"formatting" toml:"formatting"`
}

// XMLFormattingConfig XML格式化配置
type XMLFormattingConfig struct {
	Indent      string `yaml:"indent" json:"indent" toml:"indent"`
	Declaration bool   `yaml:"declaration" json:"declaration" toml:"declaration"`
	Encoding    string `yaml:"encoding" json:"encoding" toml:"encoding"`
}

// FieldsConfig 字段配置
type FieldsConfig struct {
	CustomNames map[string]string `yaml:"custom_names" json:"custom_names" toml:"custom_names"`
	Filter      struct {
		Include []string `yaml:"include" json:"include" toml:"include"`
		Exclude []string `yaml:"exclude" json:"exclude" toml:"exclude"`
	} `yaml:"filter" json:"filter" toml:"filter"`
	Processing struct {
		MaxLength      int  `yaml:"max_length" json:"max_length" toml:"max_length"`
		AddLineNumbers bool `yaml:"add_line_numbers" json:"add_line_numbers" toml:"add_line_numbers"`
		TrimWhitespace bool `yaml:"trim_whitespace" json:"trim_whitespace" toml:"trim_whitespace"`
		CodeHighlight  bool `yaml:"code_highlight" json:"code_highlight" toml:"code_highlight"`
	} `yaml:"processing" json:"processing" toml:"processing"`
}

// FiltersConfig 文件过滤配置
type FiltersConfig struct {
	MaxFileSize     string   `yaml:"max_file_size" json:"max_file_size" toml:"max_file_size"`
	ExcludePatterns []string `yaml:"exclude_patterns" json:"exclude_patterns" toml:"exclude_patterns"`
	IncludePatterns []string `yaml:"include_patterns" json:"include_patterns" toml:"include_patterns"`
	MaxDepth        int      `yaml:"max_depth" json:"max_depth" toml:"max_depth"`
	FollowSymlinks  bool     `yaml:"follow_symlinks" json:"follow_symlinks" toml:"follow_symlinks"`
	ExcludeBinary   bool     `yaml:"exclude_binary" json:"exclude_binary" toml:"exclude_binary"`
}

// OutputConfig 输出配置
type OutputConfig struct {
	Format       string `yaml:"format" json:"format" toml:"format"`
	FilePath     string `yaml:"file_path" json:"file_path" toml:"file_path"`
	Encoding     string `yaml:"encoding" json:"encoding" toml:"encoding"`
	DefaultFormat    string `yaml:"default_format" json:"default_format" toml:"default_format"`
	OutputDir        string `yaml:"output_dir" json:"output_dir" toml:"output_dir"`
	FilenameTemplate string `yaml:"filename_template" json:"filename_template" toml:"filename_template"`
	TimestampFormat  string `yaml:"timestamp_format" json:"timestamp_format" toml:"timestamp_format"`
}

// UIConfig 界面配置
type UIConfig struct {
	Theme         string `yaml:"theme" json:"theme" toml:"theme"`
	ShowProgress  bool   `yaml:"show_progress" json:"show_progress" toml:"show_progress"`
	ShowSize      bool   `yaml:"show_size" json:"show_size" toml:"show_size"`
	ShowDate      bool   `yaml:"show_date" json:"show_date" toml:"show_date"`
	ShowPreview   bool   `yaml:"show_preview" json:"show_preview" toml:"show_preview"`
	Selector struct {
		ShowHidden   bool `yaml:"show_hidden" json:"show_hidden" toml:"show_hidden"`
		ShowSize     bool `yaml:"show_size" json:"show_size" toml:"show_size"`
		ShowModified bool `yaml:"show_modified" json:"show_modified" toml:"show_modified"`
	} `yaml:"selector" json:"selector" toml:"selector"`
	Autocomplete AutocompleteConfig `yaml:"autocomplete" json:"autocomplete" toml:"autocomplete"`
}

// AutocompleteConfig 自动补全配置
type AutocompleteConfig struct {
	Enabled        bool `yaml:"enabled" json:"enabled" toml:"enabled"`
	MinChars       int  `yaml:"min_chars" json:"min_chars" toml:"min_chars"`
	MaxSuggestions int  `yaml:"max_suggestions" json:"max_suggestions" toml:"max_suggestions"`
}

// CompleteContext 补全上下文
type CompleteContext struct {
	Type CompleteType
	Data map[string]interface{}
}

// CompleteType 补全类型
type CompleteType int

const (
	CompleteFilePath CompleteType = iota
	CompleteDirectory
	CompleteExtension
	CompletePattern
	CompleteGeneric
)

// SelectOptions 选择选项
type SelectOptions struct {
	Recursive       bool
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
}

// FileProcessingConfig 文件处理配置
type FileProcessingConfig struct {
	IncludeHidden   bool
	MaxFileSize   int64
	MaxDepth      int
	ExcludePatterns []string
	IncludePatterns []string
	IncludeContent  bool
	IncludeHash     bool
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
	NonInteractive   bool
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
