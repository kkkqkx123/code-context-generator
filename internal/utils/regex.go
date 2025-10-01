// Package utils 提供通用工具函数
package utils

import (
	"fmt"
	"regexp"
)

// MatchPattern 匹配模式
func MatchPattern(pattern, text string) (bool, error) {
	matched, err := regexp.MatchString(pattern, text)
	if err != nil {
		return false, fmt.Errorf("正则表达式匹配失败: %w", err)
	}
	return matched, nil
}

// FindMatches 查找所有匹配
func FindMatches(pattern, text string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("编译正则表达式失败: %w", err)
	}
	return re.FindAllString(text, -1), nil
}

// ReplacePattern 替换模式
func ReplacePattern(pattern, replacement, text string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("编译正则表达式失败: %w", err)
	}
	return re.ReplaceAllString(text, replacement), nil
}