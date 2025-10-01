// Package selector 提供文件工具函数
package selector

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// GetFileType 获取文件类型
type GetFileType func(path string) string

// DefaultGetFileType 默认文件类型获取函数
func DefaultGetFileType(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return "unknown"
	}
	return strings.TrimPrefix(ext, ".")
}

// parseFileSize 解析文件大小字符串为字节数
func parseFileSize(sizeStr string) (int64, error) {
	sizeStr = strings.TrimSpace(strings.ToUpper(sizeStr))

	// 提取数字和单位
	var numStr string
	var unit string

	for i, char := range sizeStr {
		if char >= '0' && char <= '9' || char == '.' {
			numStr += string(char)
		} else {
			unit = sizeStr[i:]
			break
		}
	}

	if numStr == "" {
		return 0, fmt.Errorf("无效的文件大小格式: %s", sizeStr)
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("解析数字失败: %w", err)
	}

	// 根据单位计算字节数
	switch strings.TrimSpace(unit) {
	case "", "B":
		return int64(num), nil
	case "K", "KB":
		return int64(num * 1024), nil
	case "M", "MB":
		return int64(num * 1024 * 1024), nil
	case "G", "GB":
		return int64(num * 1024 * 1024 * 1024), nil
	default:
		return 0, fmt.Errorf("不支持的大小单位: %s", unit)
	}
}

// GetFileIcon 获取文件图标
type GetFileIcon func(path string) string

// DefaultGetFileIcon 默认文件图标获取函数
func DefaultGetFileIcon(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".go":
		return "🐹"
	case ".py":
		return "🐍"
	case ".js":
		return "📜"
	case ".ts":
		return "📘"
	case ".json":
		return "📋"
	case ".xml":
		return "📄"
	case ".yaml", ".yml":
		return "📋"
	case ".toml":
		return "⚙️"
	case ".md":
		return "📝"
	case ".txt":
		return "📄"
	default:
		return "📄"
	}
}

// FileInfo 文件信息结构
type FileInfo struct {
	Path     string
	Name     string
	Size     int64
	ModTime  time.Time
	IsDir    bool
	IsHidden bool
	Icon     string
	Type     string
}

// GetFileInfo 获取文件信息
func GetFileInfo(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Path:     path,
		Name:     info.Name(),
		Size:     info.Size(),
		ModTime:  info.ModTime(),
		IsDir:    info.IsDir(),
		IsHidden: strings.HasPrefix(info.Name(), "."),
		Icon:     DefaultGetFileIcon(path),
		Type:     DefaultGetFileType(path),
	}, nil
}

// GetDirectoryContents 获取目录内容
func GetDirectoryContents(path string, showHidden bool) ([]FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var contents []FileInfo
	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		// 检查隐藏文件
		if !showHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := GetFileInfo(fullPath)
		if err != nil {
			continue
		}

		contents = append(contents, *info)
	}

	return contents, nil
}