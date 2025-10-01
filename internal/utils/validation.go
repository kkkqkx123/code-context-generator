// Package utils 提供通用工具函数
package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

// IsValidFilename 检查文件名是否有效
func IsValidFilename(filename string) bool {
	if filename == "" {
		return false
	}
	
	// 检查是否包含非法字符
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(filename, char) {
			return false
		}
	}
	
	// 检查是否以点或空格开头/结尾
	if strings.HasPrefix(filename, ".") || strings.HasSuffix(filename, ".") ||
	   strings.HasPrefix(filename, " ") || strings.HasSuffix(filename, " ") {
		return false
	}
	
	return true
}

// IsValidPath 检查路径是否有效
func IsValidPath(path string) bool {
	if path == "" {
		return false
	}
	
	// 检查路径长度
	if len(path) > 260 { // Windows路径长度限制
		return false
	}
	
	// 检查是否包含空字符
	if strings.Contains(path, "\x00") {
		return false
	}
	
	return true
}

// SafePathJoin 安全地连接路径
func SafePathJoin(base, elem string) (string, error) {
	// 检查路径遍历攻击
	if strings.Contains(elem, "..") {
		return "", fmt.Errorf("路径包含非法字符: %s", elem)
	}
	
	joined := filepath.Join(base, elem)
	
	// 确保结果仍在基础路径内
	if !strings.HasPrefix(filepath.Clean(joined), filepath.Clean(base)) {
		return "", fmt.Errorf("路径超出基础目录范围")
	}
	
	return joined, nil
}