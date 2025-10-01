// Package filesystem 提供文件系统遍历和过滤功能
package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"code-context-generator/internal/utils"
	"code-context-generator/pkg/types"
)

// GetFileInfo 获取文件信息
func (w *FileSystemWalker) GetFileInfo(path string) (*types.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("获取文件状态失败: %w", err)
	}

	// 检查是否为二进制文件
	isBinary := !utils.IsTextFile(path)
	
	var content string
	if !isBinary {
		// 使用编码感知的文件读取
		fileContent, _, err := utils.ReadFileContent(path, 0) // 0表示无大小限制
		if err != nil {
			return nil, fmt.Errorf("读取文件内容失败: %w", err)
		}
		content = fileContent
	}

	return &types.FileInfo{
		Path:     path,
		Name:     info.Name(),
		Size:     info.Size(),
		ModTime:  info.ModTime(),
		IsDir:    info.IsDir(),
		Content:  content,
		IsBinary: isBinary,
	}, nil
}

// GetFolderInfo 获取文件夹信息
func (w *FileSystemWalker) GetFolderInfo(path string) (*types.FolderInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("获取文件夹状态失败: %w", err)
	}

	// 读取文件夹内容
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("读取文件夹内容失败: %w", err)
	}

	var files []types.FileInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(path, entry.Name())
			fileInfo, err := w.GetFileInfo(filePath)
			if err != nil {
				continue // 跳过无法读取的文件
			}
			files = append(files, *fileInfo)
		}
	}

	return &types.FolderInfo{
		Path:    path,
		Name:    info.Name(),
		ModTime: info.ModTime(),
		Files:   files,
	}, nil
}