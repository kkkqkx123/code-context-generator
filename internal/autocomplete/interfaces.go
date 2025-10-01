// Package autocomplete 提供自动补全功能
package autocomplete

import (
	"code-context-generator/pkg/types"
)

// Autocompleter 自动补全器接口
type Autocompleter interface {
	Complete(input string, context *types.CompleteContext) ([]string, error)
	GetSuggestions(input string, maxSuggestions int) []string
	UpdateCache(path string) error
	ClearCache()
	GetCacheSize() int
}

// SuggestionProvider 建议提供者接口
type SuggestionProvider interface {
	GetSuggestions(input string, context *types.CompleteContext) ([]Suggestion, error)
}

// Suggestion 建议项
type Suggestion struct {
	Text        string
	Description string
	Type        string
	Icon        string
}