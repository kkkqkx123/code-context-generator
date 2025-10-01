// Package filesystem 提供文件系统遍历和过滤功能
package filesystem

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"code-context-generator/internal/utils"
)

// GetFileExtension 获取文件扩展名
func GetFileExtension(filename string) string {
	// 隐藏文件（以.开头）没有扩展名
	if strings.HasPrefix(filename, ".") && len(filename) > 1 {
		// 检查是否包含另一个点（如.gitignore）
		lastDotIndex := strings.LastIndex(filename, ".")
		if lastDotIndex == 0 {
			// 只有开头的点，没有扩展名
			return ""
		}
	}
	return filepath.Ext(filename)
}

// IsHiddenFile 检查是否为隐藏文件
func IsHiddenFile(filename string) bool {
	return strings.HasPrefix(filename, ".")
}

// GetFileSize 获取文件大小
func GetFileSize(path string) (int64, error) {
	return utils.GetFileSize(path)
}

// GetFileModTime 获取文件修改时间
func GetFileModTime(path string) (time.Time, error) {
	return utils.GetFileModTime(path)
}

// IsDirectory 检查是否为目录
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsSymlink 检查是否为符号链接
func IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

// GetSymlinkTarget 获取符号链接目标
func GetSymlinkTarget(path string) (string, error) {
	target, err := os.Readlink(path)
	if err != nil {
		return "", err
	}

	// 如果是相对路径，转换为绝对路径
	if !filepath.IsAbs(target) {
		dir := filepath.Dir(path)
		target = filepath.Join(dir, target)
	}

	return filepath.Abs(target)
}

// CreateDirectory 创建目录
func CreateDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// RemoveDirectory 删除目录
func RemoveDirectory(path string) error {
	return os.RemoveAll(path)
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// MoveFile 移动文件
func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

// GetDirectorySize 获取目录大小
func GetDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// GetDirectoryFileCount 获取目录中的文件数量
func GetDirectoryFileCount(path string) (int, error) {
	count := 0
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count, err
}