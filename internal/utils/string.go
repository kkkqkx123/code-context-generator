// Package utils 提供通用工具函数
package utils

import (
	"strings"
)

// TruncateString 截断字符串
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	if maxLength <= 3 {
		return s[:maxLength]
	}
	return s[:maxLength-3] + "..."
}

// PadString 填充字符串
func PadString(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(string(padChar), length-len(s))
	return s + padding
}

// PadLeft 左填充
func PadLeft(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(string(padChar), length-len(s))
	return padding + s
}

// PadCenter 居中填充
func PadCenter(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	totalPadding := length - len(s)
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding
	return strings.Repeat(string(padChar), leftPadding) + s + strings.Repeat(string(padChar), rightPadding)
}

// RemoveDuplicates 移除字符串切片中的重复项
func RemoveDuplicates(strings []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(strings))
	
	for _, s := range strings {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	
	return result
}

// SplitLines 分割字符串为多行
func SplitLines(s string) []string {
	return strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
}

// JoinLines 连接多行为字符串
func JoinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

// CountLines 计算行数
func CountLines(s string) int {
	return len(SplitLines(s))
}