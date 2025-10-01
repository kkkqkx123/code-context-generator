// Package filesystem 提供文件系统遍历和过滤功能
package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"code-context-generator/internal/utils"
	"code-context-generator/pkg/constants"
	"code-context-generator/pkg/types"
)

// Walker 文件系统遍历器接口
type Walker interface {
	Walk(rootPath string, options *types.WalkOptions) (*types.ContextData, error)
	GetFileInfo(path string) (*types.FileInfo, error)
	GetFolderInfo(path string) (*types.FolderInfo, error)
	FilterFiles(files []string, patterns []string) []string
	FilterBySize(path string, maxSize int64) bool
}

// FileSystemWalker 文件系统遍历器实现
type FileSystemWalker struct {
	mu sync.RWMutex
}

// NewWalker 创建新的文件系统遍历器
func NewWalker() Walker {
	return &FileSystemWalker{}
}

// NewFileSystemWalker 创建新的文件系统遍历器（别名）
func NewFileSystemWalker(options types.WalkOptions) Walker {
	return &FileSystemWalker{}
}

// Walk 遍历文件系统
func (w *FileSystemWalker) Walk(rootPath string, options *types.WalkOptions) (*types.ContextData, error) {
	if options == nil {
		options = &types.WalkOptions{
			MaxDepth:        constants.DefaultMaxDepth,
			MaxFileSize:     10 * 1024 * 1024,
			ExcludePatterns: constants.DefaultExcludePatterns,
			IncludePatterns: []string{},
			FollowSymlinks:  false,
		}
	}

	// 验证根路径
	if _, err := os.Stat(rootPath); err != nil {
		return nil, fmt.Errorf("根路径不存在: %w", err)
	}

	var contextData types.ContextData
	var wg sync.WaitGroup
	var mu sync.Mutex
	var walkErrors []error
	
	// 初始化contextData的统计信息
	contextData.Files = []types.FileInfo{}
	contextData.Folders = []types.FolderInfo{}
	contextData.Metadata = make(map[string]interface{})

	// 遍历文件系统
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			walkErrors = append(walkErrors, err)
			return nil // 继续遍历
		}

		// 检查深度限制
		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return err
		}

		depth := strings.Count(relPath, string(os.PathSeparator))
		if options.MaxDepth > 0 && depth >= options.MaxDepth {
			if info.IsDir() {
				return filepath.SkipDir
			}
			// 跳过深度超过限制的文件
			return nil
		}

		// 处理文件
		if !info.IsDir() {
			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()

				// 应用过滤器
				if !w.shouldIncludeFile(filePath, options) {
					return
				}

				// 获取文件信息
				fileInfo, err := w.GetFileInfo(filePath)
				if err != nil {
					mu.Lock()
					walkErrors = append(walkErrors, fmt.Errorf("获取文件信息失败 %s: %w", filePath, err))
					mu.Unlock()
					return
				}

				mu.Lock()
				contextData.Files = append(contextData.Files, *fileInfo)
				contextData.FileCount++
				contextData.TotalSize += fileInfo.Size
				mu.Unlock()
			}(path)
		} else {
			// 处理文件夹
			if path != rootPath { // 跳过根路径
				folderInfo, err := w.GetFolderInfo(path)
				if err != nil {
					mu.Lock()
					walkErrors = append(walkErrors, fmt.Errorf("获取文件夹信息失败 %s: %w", path, err))
					mu.Unlock()
					return nil
				}

				mu.Lock()
				contextData.Folders = append(contextData.Folders, *folderInfo)
				contextData.FolderCount++
				mu.Unlock()
			}
		}

		return nil
	})

	wg.Wait()

	if err != nil {
		return nil, fmt.Errorf("遍历文件系统失败: %w", err)
	}

	if len(walkErrors) > 0 {
		// 记录错误但不中断流程
		fmt.Printf("遍历过程中遇到 %d 个错误\n", len(walkErrors))
		for _, e := range walkErrors {
			fmt.Printf("  - %v\n", e)
		}
	}

	return &contextData, nil
}

// shouldIncludeFile 检查是否应该包含文件
func (w *FileSystemWalker) shouldIncludeFile(path string, options *types.WalkOptions) bool {
	// 检查文件大小
	if !w.FilterBySize(path, options.MaxFileSize) {
		return false
	}

	// 检查是否为二进制文件（如果启用了二进制文件排除）
	if options.ExcludeBinary && utils.IsBinaryFile(path) {
		return false
	}

	// 检查包含模式
	if len(options.IncludePatterns) > 0 {
		matched := false
		filename := filepath.Base(path)
		for _, pattern := range options.IncludePatterns {
			if matchedPattern, _ := filepath.Match(pattern, filename); matchedPattern {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 检查排除模式
	if len(options.ExcludePatterns) > 0 {
		filename := filepath.Base(path)
		for _, pattern := range options.ExcludePatterns {
			if matchedPattern, _ := filepath.Match(pattern, filename); matchedPattern {
				return false
			}
		}
	}
	
	return true
}