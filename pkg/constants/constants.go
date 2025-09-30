// Package constants 定义项目的常量
package constants

import "time"

// 应用常量
const (
	AppName        = "code-context-generator"
	AppVersion     = "1.0.0"
	AppDescription = "High-Performance Code Context Generation Tool"
)

// 配置常量
const (
	DefaultConfigFile       = "config.yaml"
	DefaultFormat           = "xml"
	DefaultOutputDir        = ""
	DefaultFilenameTemplate = "context_{{.timestamp}}.{{.extension}}"
	DefaultTimestampFormat  = "20060102_150405"
	MaxFileSizeDefault      = 10 * 1024 * 1024 // 10MB
)

// 文件处理常量
const (
	MaxFileSizeLimit  = 100 * 1024 * 1024 // 100MB
	DefaultMaxDepth   = 0                 // 无限制
	BufferSize        = 32 * 1024         // 32KB
	MaxConcurrency    = 10
	ChannelBufferSize = 100
)

// UI常量
const (
	DefaultMinChars       = 1
	DefaultMaxSuggestions = 10
	DefaultShowHidden     = false
	DefaultShowSize       = true
	DefaultShowModified   = false
)

// 格式常量
const (
	FormatXML      = "xml"
	FormatJSON     = "json"
	FormatTOML     = "toml"
	FormatMarkdown = "markdown"
)

// 错误消息常量
const (
	ErrMsgConfigLoad       = "配置文件加载失败"
	ErrMsgConfigValidate   = "配置验证失败"
	ErrMsgFileRead         = "文件读取失败"
	ErrMsgFileWrite        = "文件写入失败"
	ErrMsgFormatGenerate   = "格式生成失败"
	ErrMsgPathInvalid      = "路径无效"
	ErrMsgPermissionDenied = "权限不足"
	ErrMsgFileTooLarge     = "文件过大"
)

// 时间常量
const (
	DefaultTimeout         = 30 * time.Second
	FileWatchInterval      = 1 * time.Second
	ProgressUpdateInterval = 100 * time.Millisecond
)

// 正则表达式模式
const (
	PatternHiddenFile  = `^\.`
	PatternGitignore   = `^\.gitignore$`
	PatternConfigFile  = `^config\.(yaml|yml|json|toml)$`
	PatternTemplateVar = `\{\{\.(\w+)\}\}`
)

// 环境变量前缀
const (
	EnvPrefix = "CODE_CONTEXT_"
)

// 默认排除模式
var DefaultExcludePatterns = []string{
	"*.tmp",
	"*.log",
	"*.swp",
	".*",
	"node_modules/",
	"target/",
	"dist/",
	"build/",
	".env",
	".git/",
	".vscode/",
	".idea/",
	"__pycache__/",
	"*.pyc",
	".venv",
	"*.class",
}

// 支持的格式列表
var SupportedFormats = []string{
	FormatXML,
	FormatJSON,
	FormatTOML,
	FormatMarkdown,
}
