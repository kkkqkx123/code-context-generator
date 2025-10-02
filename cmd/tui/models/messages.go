package models

import (
	"code-context-generator/internal/selector"
	"code-context-generator/pkg/types"
)

// ProgressMsg 进度消息
type ProgressMsg struct {
	Progress float64
	Status   string
}

// ProcessingUpdateMsg 处理更新消息
type ProcessingUpdateMsg struct {
	FilesProcessed   int
	TotalFiles       int
	CurrentFile      string
	Status          string
}

// ResultMsg 结果消息
type ResultMsg struct {
	Result *types.WalkResult
}

// ErrorMsg 错误消息
type ErrorMsg struct {
	Err error
}

// FileSelectionMsg 文件选择消息
type FileSelectionMsg struct {
	Selected []string
}

// ConfigUpdateMsg 配置更新消息
type ConfigUpdateMsg struct {
	Config *types.Config
}

// FileListMsg 文件列表消息
type FileListMsg struct {
	Items []selector.FileItem
}

// UpdateSuggestionsMsg 更新建议消息
type UpdateSuggestionsMsg struct {
	Suggestions []string
	Err         error
}

// ApplySuggestionMsg 应用建议消息
type ApplySuggestionMsg struct {
	Suggestion string
}