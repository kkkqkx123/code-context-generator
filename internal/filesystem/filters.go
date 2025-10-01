// Package filesystem 提供文件系统遍历和过滤功能
package filesystem

import (
	"os"
	"path/filepath"
)

// FilterFiles 根据模式过滤文件
func (w *FileSystemWalker) FilterFiles(files []string, patterns []string) []string {
	if len(patterns) == 0 {
		return files
	}

	var filtered []string
	for _, file := range files {
		for _, pattern := range patterns {
			matched, err := filepath.Match(pattern, filepath.Base(file))
			if err == nil && matched {
				filtered = append(filtered, file)
				break
			}
		}
	}

	return filtered
}

// FilterBySize 根据文件大小过滤
func (w *FileSystemWalker) FilterBySize(path string, maxSize int64) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if maxSize <= 0 {
		return true // 没有大小限制
	}

	return info.Size() <= maxSize
}