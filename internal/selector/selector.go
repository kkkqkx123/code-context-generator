// Package selector 提供文件和文件夹选择功能
package selector

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"code-context-generator/pkg/constants"
	"code-context-generator/pkg/types"
)

// Selector 选择器接口
type Selector interface {
	SelectFiles(rootPath string, options *types.SelectOptions) ([]string, error)
	SelectFolders(rootPath string, options *types.SelectOptions) ([]string, error)
	InteractiveSelect(items []string, prompt string) ([]string, error)
	FilterItems(items []string, filter string) []string
	SortItems(items []string, sortBy string) []string
}

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
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 继续遍历
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查深度限制
		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return nil
		}

		depth := strings.Count(relPath, string(os.PathSeparator))
		if depth > options.MaxDepth {
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

// PatternMatcher 模式匹配器
type PatternMatcher struct {
	patterns []string
}

// NewPatternMatcher 创建模式匹配器
func NewPatternMatcher(patterns []string) *PatternMatcher {
	return &PatternMatcher{
		patterns: patterns,
	}
}

// Match 检查是否匹配任何模式
func (pm *PatternMatcher) Match(path string) bool {
	filename := filepath.Base(path)
	for _, pattern := range pm.patterns {
		matched, err := filepath.Match(pattern, filename)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// MatchAny 检查是否匹配任何模式（支持通配符）
func (pm *PatternMatcher) MatchAny(path string) bool {
	filename := filepath.Base(path)
	for _, pattern := range pm.patterns {
		// 支持通配符匹配
		if matched, _ := filepath.Match(pattern, filename); matched {
			return true
		}
		// 支持包含匹配
		if strings.Contains(filename, pattern) {
			return true
		}
	}
	return false
}

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

// NewFileSelector 创建新的文件选择器
func NewFileSelector(config *types.Config) *FileSelector {
	return &FileSelector{
		config: config,
	}
}

// SelectorOptions 选择器选项
type SelectorOptions struct {
	MaxDepth        int
	IncludePatterns []string
	ExcludePatterns []string
	ShowHidden      bool
	SortBy          string
}

// FileItem 文件项
type FileItem struct {
	Path     string
	Name     string
	Size     int64
	ModTime  time.Time
	IsDir    bool
	IsHidden bool
	Icon     string
	Type     string
	Selected bool
}
