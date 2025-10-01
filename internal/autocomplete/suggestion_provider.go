// Package autocomplete 提供自动补全功能
package autocomplete

import (
	"code-context-generator/pkg/constants"
	"code-context-generator/pkg/types"
)

// CompositeSuggestionProvider 组合建议提供者
type CompositeSuggestionProvider struct {
	providers []SuggestionProvider
}

// NewCompositeSuggestionProvider 创建组合建议提供者
func NewCompositeSuggestionProvider(providers ...SuggestionProvider) *CompositeSuggestionProvider {
	return &CompositeSuggestionProvider{
		providers: providers,
	}
}

// GetSuggestions 获取建议
func (c *CompositeSuggestionProvider) GetSuggestions(input string, context *types.CompleteContext) ([]Suggestion, error) {
	var allSuggestions []Suggestion

	for _, provider := range c.providers {
		suggestions, err := provider.GetSuggestions(input, context)
		if err != nil {
			continue // 跳过出错的提供者
		}
		allSuggestions = append(allSuggestions, suggestions...)
	}

	// 去重和限制数量
	uniqueSuggestions := removeDuplicateSuggestions(allSuggestions)
	if len(uniqueSuggestions) > constants.DefaultMaxSuggestions {
		uniqueSuggestions = uniqueSuggestions[:constants.DefaultMaxSuggestions]
	}

	return uniqueSuggestions, nil
}

// removeDuplicateSuggestions 移除重复建议
func removeDuplicateSuggestions(suggestions []Suggestion) []Suggestion {
	seen := make(map[string]bool)
	var result []Suggestion

	for _, suggestion := range suggestions {
		if !seen[suggestion.Text] {
			seen[suggestion.Text] = true
			result = append(result, suggestion)
		}
	}

	return result
}