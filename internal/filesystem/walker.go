// Package filesystem 提供文件系统遍历和过滤功能
package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"code-context-generator/internal/utils"
	"code-context-generator/pkg/constants"
	"code-context-generator/pkg/types"
)

// Walker 文件系统遍历器接口
type Walker interface {
	Walk(rootPath string, options *types.WalkOptions) (*types.ContextData, error)
	WalkWithProgress(rootPath string, options *types.WalkOptions, progressCallback func(processed, total int, currentFile string)) (*types.ContextData, error)
	GetFileInfo(path string) (*types.FileInfo, error)
	GetFolderInfo(path string) (*types.FolderInfo, error)
	FilterFiles(files []string, patterns []string) []string
	FilterBySize(path string, maxSize int64) bool
	SetConfig(config *types.Config)
}

// FileSystemWalker 文件系统遍历器实现
type FileSystemWalker struct {
	mu           sync.RWMutex
	maxWorkers   int
	maxFileCount int
	maxDepth     int
	timeout      time.Duration
	config       *types.Config // 添加配置引用
}

// NewWalker 创建遍历器
func NewWalker() Walker {
	return &FileSystemWalker{
		maxWorkers:   10,               // 限制并发worker数量
		maxFileCount: 1000,             // 限制最大文件数量
		maxDepth:     5,                // 限制最大深度
		timeout:      30 * time.Second, // 30秒超时
		config:       nil,              // 默认无配置
	}
}

// NewFileSystemWalker 创建新的文件系统遍历器（别名）
func NewFileSystemWalker(options types.WalkOptions) Walker {
	return &FileSystemWalker{
		maxWorkers:   10,               // 限制并发worker数量
		maxFileCount: 1000,             // 限制最大文件数量
		maxDepth:     5,                // 限制最大深度
		timeout:      30 * time.Second, // 30秒超时
		config:       nil,              // 默认无配置
	}
}

// SetConfig 设置配置
func (w *FileSystemWalker) SetConfig(config *types.Config) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.config = config
}

// Walk 遍历文件系统
func (w *FileSystemWalker) Walk(rootPath string, options *types.WalkOptions) (*types.ContextData, error) {
	return w.WalkWithProgress(rootPath, options, nil)
}

// WalkWithProgress 带进度回调的遍历文件系统
func (w *FileSystemWalker) WalkWithProgress(rootPath string, options *types.WalkOptions, progressCallback func(processed, total int, currentFile string)) (*types.ContextData, error) {
	if options == nil {
		options = &types.WalkOptions{
			MaxDepth:        constants.DefaultMaxDepth,
			MaxFileSize:     10 * 1024 * 1024,
			ExcludePatterns: constants.DefaultExcludePatterns,
			IncludePatterns: []string{},
			FollowSymlinks:  false,
		}
	}

	var contextData types.ContextData
	var wg sync.WaitGroup
	var mu sync.Mutex
	var walkErrors []error

	// 初始化contextData的统计信息
	contextData.Files = []types.FileInfo{}
	contextData.Folders = []types.FolderInfo{}
	contextData.Metadata = make(map[string]interface{})
	// 设置根路径 - 对于多个文件，使用第一个文件的路径
	if len(options.MultipleFiles) > 0 {
		contextData.Metadata["root_path"] = filepath.Dir(options.MultipleFiles[0])
		contextData.Metadata["multiple_files"] = options.MultipleFiles
	} else {
		contextData.Metadata["root_path"] = rootPath
	}

	// 如果指定了多个文件，直接处理这些文件而不遍历目录
	if len(options.MultipleFiles) > 0 {
		return w.processMultipleFiles(options.MultipleFiles, options, progressCallback)
	}

	// 验证根路径
	if _, err := os.Stat(rootPath); err != nil {
		return nil, fmt.Errorf("根路径不存在: %w", err)
	}

	// 首先统计总文件数
	totalFiles := 0
	processedFiles := 0
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if w.shouldIncludeFile(path, rootPath, options) {
			totalFiles++
		}
		return nil
	})

	// 限制文件数量
	if totalFiles > w.maxFileCount {
		return nil, fmt.Errorf("文件数量超过限制: %d > %d", totalFiles, w.maxFileCount)
	}

	semaphore := make(chan struct{}, w.maxWorkers) // 限制并发数量
	progressMu := sync.Mutex{}                     // 保护进度更新

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
		// 新的max-depth逻辑：
		// 0: 只扫描当前目录（不递归）
		// 1: 递归1层
		// -1: 无限递归
		if options.MaxDepth >= 0 && depth >= options.MaxDepth {
			if info.IsDir() {
				return filepath.SkipDir
			}
			// 跳过深度超过限制的文件
			return nil
		}

		// 处理文件
		if !info.IsDir() && w.shouldIncludeFile(path, rootPath, options) {
			semaphore <- struct{}{} // 获取信号量
			wg.Add(1)
			go func(filePath string, rootPath string) {
				defer func() {
					<-semaphore // 释放信号量
					wg.Done()
				}()

				// 获取文件信息
				fileInfo, err := w.GetFileInfo(filePath)
				if err != nil {
					mu.Lock()
					walkErrors = append(walkErrors, fmt.Errorf("获取文件信息失败 %s: %w", filePath, err))
					mu.Unlock()
					return
				}

				mu.Lock()
				// 检查文件是否已经存在，避免重复
				fileExists := false
				for _, existingFile := range contextData.Files {
					if existingFile.Path == fileInfo.Path {
						fileExists = true
						break
					}
				}
				
				if !fileExists {
					contextData.Files = append(contextData.Files, *fileInfo)
					contextData.FileCount++
					// 只有在包含元信息时才累加文件大小
					if w.config != nil && w.config.Output.IncludeMetadata {
						contextData.TotalSize += fileInfo.Size
					}
				}
				mu.Unlock()

				// 更新进度
				progressMu.Lock()
				processedFiles++
				currentProcessed := processedFiles
				progressMu.Unlock()

				if progressCallback != nil && currentProcessed%10 == 0 { // 每10个文件更新一次进度
					progressCallback(currentProcessed, totalFiles, filepath.Base(filePath))
				}
			}(path, rootPath)
		} else if info.IsDir() {
			// 处理文件夹 - 只统计符合过滤条件的文件夹
			if path != rootPath { // 跳过根路径
				// 检查文件夹是否应该被包含（基于排除模式）
				shouldInclude := true
				if len(options.ExcludePatterns) > 0 {
					folderName := filepath.Base(path)
					relPath, _ := filepath.Rel(rootPath, path)
					relPath = filepath.ToSlash(relPath)
					
					for _, pattern := range options.ExcludePatterns {
						// 检查文件夹名是否匹配排除模式
						if matched, _ := filepath.Match(pattern, folderName); matched {
							shouldInclude = false
							break
						}
						// 检查相对路径是否匹配排除模式
						if strings.Contains(pattern, "/") {
							if matched, _ := filepath.Match(pattern, relPath); matched {
								shouldInclude = false
								break
							}
							// 对于目录模式（以/结尾），检查当前文件夹是否匹配
							if strings.HasSuffix(pattern, "/") {
								dirPattern := strings.TrimSuffix(pattern, "/")
								if matched, _ := filepath.Match(dirPattern, folderName); matched {
									shouldInclude = false
									break
								}
							}
						}
					}
				}
				
				if shouldInclude {
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
				// 只有在包含元信息时才累加文件夹大小
				if w.config != nil && w.config.Output.IncludeMetadata {
					contextData.TotalSize += folderInfo.Size
				}
				mu.Unlock()
				}
			}
		}

		return nil
	})

	wg.Wait()

	// 最终进度更新
	if progressCallback != nil {
		progressCallback(totalFiles, totalFiles, "完成")
	}

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
func (w *FileSystemWalker) shouldIncludeFile(path string, rootPath string, options *types.WalkOptions) bool {
	// 如果指定了多个文件，只包含这些文件
	if len(options.MultipleFiles) > 0 {
		// 将路径转换为绝对路径进行比较
		absPath, err := filepath.Abs(path)
		if err != nil {
			return false
		}

		for _, selectedFile := range options.MultipleFiles {
			absSelectedFile, err := filepath.Abs(selectedFile)
			if err != nil {
				continue
			}
			if absPath == absSelectedFile {
				return true
			}
		}
		return false
	}

	// 如果指定了选中的文件，只包含这些文件（向后兼容）
	if len(options.SelectedFiles) > 0 {
		// 将路径转换为绝对路径进行比较
		absPath, err := filepath.Abs(path)
		if err != nil {
			return false
		}

		for _, selectedFile := range options.SelectedFiles {
			absSelectedFile, err := filepath.Abs(selectedFile)
			if err != nil {
				continue
			}
			if absPath == absSelectedFile {
				return true
			}
		}
		return false
	}

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
			// 尝试匹配相对路径（用于目录模式如 *.go）
			if strings.Contains(pattern, "/") {
				rel, _ := filepath.Rel(rootPath, path)
				// 将Windows路径分隔符转换为正斜杠以匹配模式
				rel = filepath.ToSlash(rel)
				if matchedPattern, _ := filepath.Match(pattern, rel); matchedPattern {
					matched = true
					break
				}
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
			// 尝试匹配文件名
			if matchedPattern, _ := filepath.Match(pattern, filename); matchedPattern {
				return false
			}
			// 尝试匹配相对路径（用于目录模式如 .git/）
			if strings.Contains(pattern, "/") {
				rel, _ := filepath.Rel(rootPath, path)
				// 将Windows路径分隔符转换为正斜杠以匹配模式
				rel = filepath.ToSlash(rel)

				if matchedPattern, _ := filepath.Match(pattern, rel); matchedPattern {
					return false
				}
				// 对于目录模式（以/结尾），检查文件是否在匹配目录下
				if strings.HasSuffix(pattern, "/") {
					dirPattern := strings.TrimSuffix(pattern, "/")
					// 检查相对路径是否以目录模式开头
					if strings.HasPrefix(rel, dirPattern+"/") {
						return false
					}
					// 检查路径中的任何目录部分是否匹配
					pathDirs := strings.Split(rel, "/")
					for i, dir := range pathDirs {
						if matchedDir, _ := filepath.Match(dirPattern, dir); matchedDir {
							// 确保这是完整目录名匹配，而不是部分匹配
							if i < len(pathDirs)-1 || rel == dirPattern {
								return false
							}
						}
					}
				}
			}
		}
	}

	return true
}

// processMultipleFiles 处理多个指定文件而不遍历目录
func (w *FileSystemWalker) processMultipleFiles(files []string, options *types.WalkOptions, progressCallback func(processed, total int, currentFile string)) (*types.ContextData, error) {
	var contextData types.ContextData
	contextData.Files = []types.FileInfo{}
	contextData.Folders = []types.FolderInfo{}
	contextData.Metadata = make(map[string]interface{})

	// 用于跟踪已处理的文件夹，避免重复
	processedFolders := make(map[string]bool)

	// 过滤并处理指定的文件
	processedFiles := 0
	totalFiles := len(files)
	
	// 用于跟踪已处理的文件，避免重复
	processedFilesMap := make(map[string]bool)

	for _, filePath := range files {
		// 验证文件存在
		info, err := os.Stat(filePath)
		if err != nil {
			continue // 跳过不存在的文件
		}

		// 只处理文件，跳过目录
		if info.IsDir() {
			continue
		}

		// 检查文件是否已经处理过（去重）
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			continue
		}
		if processedFilesMap[absPath] {
			continue // 跳过已处理的文件
		}

		// 检查是否应该包含此文件
		if !w.shouldIncludeFile(filePath, filepath.Dir(filePath), options) {
			continue
		}

		// 获取文件信息
		fileInfo, err := w.GetFileInfo(filePath)
		if err != nil {
			continue // 跳过无法获取信息的文件
		}

		contextData.Files = append(contextData.Files, *fileInfo)
		contextData.FileCount++
		contextData.TotalSize += fileInfo.Size
		processedFilesMap[absPath] = true
		processedFiles++

		// 获取文件所在目录的路径
		dirPath := filepath.Dir(filePath)
		
		// 如果目录还未处理过，添加文件夹信息
		if !processedFolders[dirPath] {
			// 获取目录信息
			dirInfo, err := os.Stat(dirPath)
			if err == nil && dirInfo.IsDir() {
				// 获取目录中的文件列表（用于计算文件数量）
				filesInDir, _ := os.ReadDir(dirPath)
				fileCount := 0
				for _, entry := range filesInDir {
					if !entry.IsDir() {
						fileCount++
					}
				}

				folderInfo := types.FolderInfo{
					Name:     filepath.Base(dirPath),
					Path:     dirPath,
					Files:    []types.FileInfo{}, // 这里不填充具体文件，保持简洁
					Folders:  nil,
					ModTime:  dirInfo.ModTime(),
					IsHidden: strings.HasPrefix(filepath.Base(dirPath), "."),
					Size:     0, // 目录大小计算复杂，这里设为0
					Count:    fileCount,
				}
				
				contextData.Folders = append(contextData.Folders, folderInfo)
				contextData.FolderCount++
				processedFolders[dirPath] = true
			}
		}

		// 更新进度
		if progressCallback != nil {
			progressCallback(processedFiles, totalFiles, filepath.Base(filePath))
		}
	}

	// 最终进度更新
	if progressCallback != nil {
		progressCallback(processedFiles, totalFiles, "完成")
	}

	return &contextData, nil
}
