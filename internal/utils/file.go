// Package utils 提供通用工具函数
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// DirectoryExists 检查目录是否存在
func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// GetFileHash 获取文件哈希值
func GetFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetFileSize 获取文件大小
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetFileModTime 获取文件修改时间
func GetFileModTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// IsTextFile 检查是否为文本文件
func IsTextFile(path string) bool {
	// 首先检查文件扩展名
	ext := strings.ToLower(filepath.Ext(path))
	textExtensions := []string{
		".txt", ".md", ".json", ".xml", ".yaml", ".yml", ".toml",
		".go", ".py", ".js", ".ts", ".java", ".cpp", ".c", ".h",
		".html", ".css", ".scss", ".sass", ".sql", ".sh", ".bat",
		".ps1", ".rb", ".php", ".rs", ".swift", ".kt", ".scala",
	}

	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}

	// 如果没有扩展名，尝试读取文件内容来判断
	if ext == "" {
		file, err := os.Open(path)
		if err != nil {
			return false // 无法打开文件，假设为二进制文件
		}
		defer file.Close()

		// 读取前512字节来判断是否为文本文件
		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return false // 读取错误，假设为二进制文件
		}

		// 检查是否包含null字节（二进制文件的标志）
		for i := 0; i < n; i++ {
			if buffer[i] == 0 {
				return false // 包含null字节，是二进制文件
			}
		}

		// 检查是否包含可打印字符
		printableCount := 0
		for i := 0; i < n; i++ {
			b := buffer[i]
			if b >= 32 && b <= 126 { // 可打印ASCII字符
				printableCount++
			} else if b == 9 || b == 10 || b == 13 { // tab, newline, carriage return
				printableCount++
			}
		}

		// 如果大部分字符都是可打印的，认为是文本文件
		if n > 0 && float64(printableCount)/float64(n) > 0.8 {
			return true
		}
	}

	return false
}

// IsBinaryFile 检查是否为二进制文件
func IsBinaryFile(path string) bool {
	return !IsTextFile(path)
}

// ReadFileContent 读取文件内容（带大小限制）
func ReadFileContent(path string, maxSize int64) (string, bool, error) {
	// 使用新的编码感知函数
	return ReadFileContentWithEncoding(path, maxSize)
}

// ReadFileContentWithEncoding 智能编码读取文件内容
func ReadFileContentWithEncoding(path string, maxSize int64) (string, bool, error) {
	// 获取文件信息
	info, err := os.Stat(path)
	if err != nil {
		return "", false, err
	}

	// 检查文件大小
	if maxSize > 0 && info.Size() > maxSize {
		return "", false, fmt.Errorf("文件大小超过限制: %d > %d", info.Size(), maxSize)
	}

	// 读取文件内容
	content, err := os.ReadFile(path)
	if err != nil {
		return "", false, fmt.Errorf("读取文件失败: %w", err)
	}

	// 检测是否为二进制文件
	isBinary := !IsTextFile(path)
	if isBinary {
		return "[二进制文件]", isBinary, nil
	}

	// 检测编码并转换
	encoding, cleanData := DetectEncoding(content)
	utf8Content, err := ConvertToUTF8(cleanData, encoding)
	if err != nil {
		return "", false, fmt.Errorf("编码转换失败: %w", err)
	}

	return utf8Content, isBinary, nil
}