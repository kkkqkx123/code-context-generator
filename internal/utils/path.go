// Package utils 提供通用工具函数
package utils

import (
	"path/filepath"
	"strings"
)

// NormalizePath 规范化路径
func NormalizePath(path string) string {
	return filepath.Clean(path)
}

// GetRelativePath 获取相对路径
func GetRelativePath(base, target string) (string, error) {
	return filepath.Rel(base, target)
}

// GetAbsolutePath 获取绝对路径
func GetAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

// IsSubPath 检查是否为子路径
func IsSubPath(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	// 如果相对路径是 "." 或空字符串，说明是同一个路径，不算子路径
	if rel == "." || rel == "" {
		return false
	}
	return !strings.HasPrefix(rel, "..")
}

// GetCommonPath 获取共同路径
func GetCommonPath(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	
	if len(paths) == 1 {
		return filepath.Dir(paths[0])
	}

	// 转换为绝对路径并清理
	absPaths := make([]string, 0, len(paths))
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			continue // 跳过无效路径
		}
		absPaths = append(absPaths, filepath.Clean(absPath))
	}
	
	if len(absPaths) == 0 {
		return ""
	}

	// 找到最短的路径
	minPath := absPaths[0]
	for _, path := range absPaths {
		if len(path) < len(minPath) {
			minPath = path
		}
	}

	// 从最短路径开始，逐步向上查找共同路径
	for {
		common := true
		for _, path := range absPaths {
			// 使用 filepath.HasPrefix 来处理路径分隔符问题
			if !filepath.HasPrefix(path, minPath) {
				common = false
				break
			}
		}
		
		if common {
			return minPath
		}
		
		parent := filepath.Dir(minPath)
		if parent == minPath {
			break
		}
		minPath = parent
	}

	return ""
}