// Package filesystem 提供文件系统遍历和过滤功能
package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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
			MaxDepth:       constants.DefaultMaxDepth,
			MaxFileSize:    10 * 1024 * 1024,
			ExcludePatterns: constants.DefaultExcludePatterns,
			IncludePatterns: []string{},
			FollowSymlinks: false,
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
		if depth > options.MaxDepth {
			if info.IsDir() {
				return filepath.SkipDir
			}
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

// GetFileInfo 获取文件信息
func (w *FileSystemWalker) GetFileInfo(path string) (*types.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("获取文件状态失败: %w", err)
	}

	// 读取文件内容
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取文件内容失败: %w", err)
	}

	return &types.FileInfo{
		Path:     path,
		Name:     info.Name(),
		Size:     info.Size(),
		ModTime:  info.ModTime(),
		IsDir:    info.IsDir(),
		Content:  string(content),
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

// 辅助方法

func (w *FileSystemWalker) shouldIncludeFile(path string, options *types.WalkOptions) bool {
	// 检查文件大小
	if !w.FilterBySize(path, options.MaxFileSize) {
		return false
	}

	// 检查包含模式
	if len(options.IncludePatterns) > 0 {
		included := false
		for _, pattern := range options.IncludePatterns {
			matched, err := filepath.Match(pattern, filepath.Base(path))
			if err == nil && matched {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}

	// 检查排除模式
	for _, pattern := range options.ExcludePatterns {
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err == nil && matched {
			return false
		}
	}

	// 检查隐藏文件
	if !options.ShowHidden && strings.HasPrefix(filepath.Base(path), ".") {
		return false
	}

	return true
}

func (w *FileSystemWalker) calculateFolderSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

// GetFileExtension 获取文件扩展名
func GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}

// IsHiddenFile 检查是否为隐藏文件
func IsHiddenFile(filename string) bool {
	return strings.HasPrefix(filename, ".")
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