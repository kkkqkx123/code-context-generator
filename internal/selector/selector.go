// Package selector 提供文件和文件夹选择功能
package selector

import (
	"time"

	"code-context-generator/pkg/types"
)

// Selector 选择器接口
type Selector interface {
	SelectFiles(rootPath string, options *types.SelectOptions) ([]string, error)
	SelectFolders(rootPath string, options *types.SelectOptions) ([]string, error)
	InteractiveSelect(items []string, prompt string) ([]string, error)
	FilterItems(items []string, filter string) []string
	SortItems(items []string, sortBy string) []string
}

// SelectorOptions 选择器选项
type SelectorOptions struct {
	MaxDepth        int
	IncludePatterns []string
	ExcludePatterns []string
	ShowHidden      bool
	SortBy          string
}

// FileItem 文件项
type FileItem struct {
	Path     string
	Name     string
	Size     int64
	ModTime  time.Time
	IsDir    bool
	IsHidden bool
	Icon     string
	Type     string
	Selected bool
}

