// Package selector 提供模式匹配功能
package selector

import (
	"path/filepath"
	"strings"
)

// PatternMatcher 模式匹配器
type PatternMatcher struct {
	patterns []string
}

// NewPatternMatcher 创建模式匹配器
func NewPatternMatcher(patterns []string) *PatternMatcher {
	return &PatternMatcher{
		patterns: patterns,
	}
}

// Match 检查是否匹配任何模式
func (pm *PatternMatcher) Match(path string) bool {
	filename := filepath.Base(path)
	for _, pattern := range pm.patterns {
		matched, err := filepath.Match(pattern, filename)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// MatchAny 检查是否匹配任何模式（支持通配符）
func (pm *PatternMatcher) MatchAny(path string) bool {
	filename := filepath.Base(path)
	for _, pattern := range pm.patterns {
		// 支持通配符匹配
		if matched, _ := filepath.Match(pattern, filename); matched {
			return true
		}
		// 支持包含匹配
		if strings.Contains(filename, pattern) {
			return true
		}
	}
	return false
}