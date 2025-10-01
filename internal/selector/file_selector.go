// Package selector 提供文件选择器具体实现
package selector

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"code-context-generator/pkg/constants"
	"code-context-generator/pkg/types"
)

// FileSelector 文件选择器实现
type FileSelector struct {
	config *types.Config
}

// NewSelector 创建新的选择器
func NewSelector(config *types.Config) Selector {
	return &FileSelector{
		config: config,
	}
}

// NewFileSelector 创建新的文件选择器
func NewFileSelector(config *types.Config) *FileSelector {
	return &FileSelector{
		config: config,
	}
}

// SelectFiles 选择文件
func (s *FileSelector) SelectFiles(rootPath string, options *types.SelectOptions) ([]string, error) {
	if options == nil {
		options = &types.SelectOptions{
			Recursive:       true,
			IncludePatterns: []string{},
			ExcludePatterns: constants.DefaultExcludePatterns,
			MaxDepth:        constants.DefaultMaxDepth,
			ShowHidden:      constants.DefaultShowHidden,
			SortBy:          "name",
		}
	}

	var files []string
	
	// 如果不递归，只处理当前目录
	if !options.Recursive {
		entries, err := os.ReadDir(rootPath)
		if err != nil {
			return nil, fmt.Errorf("读取目录失败: %w", err)
		}
		
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			
			fullPath := filepath.Join(rootPath, entry.Name())
			info, err := entry.Info()
			if err != nil {
				continue
			}
			
			if s.shouldIncludeFile(fullPath, info, options) {
				files = append(files, fullPath)
			}
		}
	} else {
		// 递归遍历
		err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // 继续遍历
			}

			// 跳过根目录
			if path == rootPath {
				return nil
			}

			// 跳过目录
			if info.IsDir() {
				// 检查深度限制
				relPath, err := filepath.Rel(rootPath, path)
				if err != nil {
					return nil
				}
				
				depth := strings.Count(relPath, string(os.PathSeparator))
				// MaxDepth 为 0 表示无限制，MaxDepth 为 1 表示只处理根目录下的文件
				if options.MaxDepth > 0 && depth >= options.MaxDepth {
					return filepath.SkipDir
				}
				return nil
			}

			// 检查文件深度限制
			relPath, err := filepath.Rel(rootPath, path)
			if err != nil {
				return nil
			}
			
			depth := strings.Count(relPath, string(os.PathSeparator))
			// MaxDepth 为 0 表示无限制，MaxDepth 为 1 表示只处理根目录下的文件
			if options.MaxDepth > 0 && depth >= options.MaxDepth {
				return nil
			}

			// 应用过滤器
			if s.shouldIncludeFile(path, info, options) {
				files = append(files, path)
			}

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("遍历文件失败: %w", err)
		}
	}

	// 排序
	files = s.SortItems(files, options.SortBy)

	return files, nil
}

// SelectFolders 选择文件夹
func (s *FileSelector) SelectFolders(rootPath string, options *types.SelectOptions) ([]string, error) {
	if options == nil {
		options = &types.SelectOptions{
			Recursive:       true,
			IncludePatterns: []string{},
			ExcludePatterns: []string{},
			MaxDepth:        constants.DefaultMaxDepth,
			ShowHidden:      constants.DefaultShowHidden,
			SortBy:          "name",
		}
	}

	var folders []string
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 继续遍历
		}

		// 跳过文件和根目录
		if !info.IsDir() || path == rootPath {
			return nil
		}

		// 检查深度限制
		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return nil
		}

		depth := strings.Count(relPath, string(os.PathSeparator))
		if depth > options.MaxDepth {
			return filepath.SkipDir
		}

		// 应用过滤器
		if s.shouldIncludeFolder(path, info, options) {
			folders = append(folders, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历文件夹失败: %w", err)
	}

	// 排序
	folders = s.SortItems(folders, options.SortBy)

	return folders, nil
}

// InteractiveSelect 交互式选择
func (s *FileSelector) InteractiveSelect(items []string, prompt string) ([]string, error) {
	if len(items) == 0 {
		return []string{}, nil
	}

	// 这里可以集成TUI选择器
	// 暂时返回所有项目
	return items, nil
}

// FilterItems 过滤项目
func (s *FileSelector) FilterItems(items []string, filter string) []string {
	if filter == "" {
		return items
	}

	var filtered []string
	filter = strings.ToLower(filter)

	for _, item := range items {
		if strings.Contains(strings.ToLower(item), filter) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// SortItems 排序项目
func (s *FileSelector) SortItems(items []string, sortBy string) []string {
	sorted := make([]string, len(items))
	copy(sorted, items)

	switch sortBy {
	case "name":
		sort.Strings(sorted)
	case "size":
		sort.Slice(sorted, func(i, j int) bool {
			info1, err1 := os.Stat(sorted[i])
			info2, err2 := os.Stat(sorted[j])
			if err1 != nil || err2 != nil {
				return sorted[i] < sorted[j]
			}
			return info1.Size() < info2.Size()
		})
	case "modified":
		sort.Slice(sorted, func(i, j int) bool {
			info1, err1 := os.Stat(sorted[i])
			info2, err2 := os.Stat(sorted[j])
			if err1 != nil || err2 != nil {
				return sorted[i] < sorted[j]
			}
			return info1.ModTime().Before(info2.ModTime())
		})
	default:
		sort.Strings(sorted)
	}

	return sorted
}

// 辅助方法

func (s *FileSelector) shouldIncludeFile(path string, info os.FileInfo, options *types.SelectOptions) bool {
	filename := filepath.Base(path)

	// 检查隐藏文件
	if !options.ShowHidden && strings.HasPrefix(filename, ".") {
		return false
	}

	// 检查包含模式
	if len(options.IncludePatterns) > 0 {
		included := false
		for _, pattern := range options.IncludePatterns {
			matched, err := filepath.Match(pattern, filename)
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
		matched, err := filepath.Match(pattern, filename)
		if err == nil && matched {
			return false
		}
	}

	// 检查文件大小
	if s.config != nil && s.config.Filters.MaxFileSize != "" {
		// 解析文件大小字符串为字节数
		maxSize, err := parseFileSize(s.config.Filters.MaxFileSize)
		if err == nil && info.Size() > maxSize {
			return false
		}
	}

	return true
}

func (s *FileSelector) shouldIncludeFolder(path string, _ os.FileInfo, options *types.SelectOptions) bool {
	foldername := filepath.Base(path)

	// 检查隐藏文件夹
	if !options.ShowHidden && strings.HasPrefix(foldername, ".") {
		return false
	}

	// 检查包含模式
	if len(options.IncludePatterns) > 0 {
		included := false
		for _, pattern := range options.IncludePatterns {
			matched, err := filepath.Match(pattern, foldername)
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
		matched, err := filepath.Match(pattern, foldername)
		if err == nil && matched {
			return false
		}
	}

	return true
}