// Package autocomplete 提供自动补全功能
package autocomplete

import (
	"time"
)

// AutocompleterOptions 自动补全选项
type AutocompleterOptions struct {
	Enabled        bool
	MinChars       int
	MaxSuggestions int
	CacheSize      int
	Timeout        time.Duration
}

// removeDuplicates 移除重复项
func removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}