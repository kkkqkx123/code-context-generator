// Package autocomplete 提供自动补全功能
package autocomplete

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"code-context-generator/pkg/constants"
	"code-context-generator/pkg/types"
)

// Autocompleter 自动补全器接口
type Autocompleter interface {
	Complete(input string, context *types.CompleteContext) ([]string, error)
	GetSuggestions(input string, maxSuggestions int) []string
	UpdateCache(path string) error
	ClearCache()
	GetCacheSize() int
}

// FilePathAutocompleter 文件路径自动补全器
type FilePathAutocompleter struct {
	cache    map[string][]string
	mu       sync.RWMutex
	config   *types.AutocompleteConfig
	maxDepth int
}

// NewAutocompleter 创建新的自动补全器
func NewAutocompleter(config *types.AutocompleteConfig) Autocompleter {
	if config == nil {
		config = &types.AutocompleteConfig{
			Enabled:        true,
			MinChars:       constants.DefaultMinChars,
			MaxSuggestions: constants.DefaultMaxSuggestions,
		}
	}

	return &FilePathAutocompleter{
		cache:    make(map[string][]string),
		config:   config,
		maxDepth: constants.DefaultMaxDepth,
	}
}

// Complete 执行自动补全
func (a *FilePathAutocompleter) Complete(input string, context *types.CompleteContext) ([]string, error) {
	if !a.config.Enabled {
		return []string{}, nil
	}

	if len(input) < a.config.MinChars {
		return []string{}, nil
	}

	switch context.Type {
	case types.CompleteFilePath:
		return a.completeFilePath(input, context)
	case types.CompleteDirectory:
		return a.completeDirectory(input, context)
	case types.CompleteExtension:
		return a.completeExtension(input, context)
	case types.CompletePattern:
		return a.completePattern(input, context)
	default:
		return a.completeGeneric(input, context)
	}
}

// GetSuggestions 获取建议列表
func (a *FilePathAutocompleter) GetSuggestions(input string, maxSuggestions int) []string {
	if !a.config.Enabled {
		return []string{}
	}

	if maxSuggestions <= 0 {
		maxSuggestions = a.config.MaxSuggestions
	}

	suggestions := a.getMatchingItems(input)

	if len(suggestions) > maxSuggestions {
		suggestions = suggestions[:maxSuggestions]
	}

	return suggestions
}

// UpdateCache 更新缓存
func (a *FilePathAutocompleter) UpdateCache(path string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 清除旧缓存
	delete(a.cache, path)

	// 获取新缓存数据
	items, err := a.scanDirectory(path)
	if err != nil {
		return fmt.Errorf("扫描目录失败: %w", err)
	}

	a.cache[path] = items
	return nil
}

// ClearCache 清除缓存
func (a *FilePathAutocompleter) ClearCache() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.cache = make(map[string][]string)
}

// GetCacheSize 获取缓存大小
func (a *FilePathAutocompleter) GetCacheSize() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return len(a.cache)
}

// 自动补全方法

func (a *FilePathAutocompleter) completeFilePath(input string, context *types.CompleteContext) ([]string, error) {
	dir := filepath.Dir(input)
	base := filepath.Base(input)

	// 如果目录不存在，尝试补全目录
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return a.completeDirectory(input, context)
	}

	// 获取目录内容
	items, err := a.getDirectoryItems(dir)
	if err != nil {
		return nil, err
	}

	// 过滤匹配的文件
	var matches []string
	for _, item := range items {
		if strings.HasPrefix(item, base) {
			fullPath := filepath.Join(dir, item)
			if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
				matches = append(matches, fullPath)
			}
		}
	}

	return matches, nil
}

func (a *FilePathAutocompleter) completeDirectory(input string, _ *types.CompleteContext) ([]string, error) {
	// 尝试不同的目录级别
	parts := strings.Split(input, string(os.PathSeparator))

	for i := len(parts); i > 0; i-- {
		partialPath := strings.Join(parts[:i], string(os.PathSeparator))

		if partialPath == "" {
			partialPath = "."
		}

		if _, err := os.Stat(partialPath); err == nil {
			// 找到存在的目录
			remaining := strings.Join(parts[i:], string(os.PathSeparator))

			items, err := a.getDirectoryItems(partialPath)
			if err != nil {
				continue
			}

			var matches []string
			for _, item := range items {
				if strings.HasPrefix(item, remaining) {
					fullPath := filepath.Join(partialPath, item)
					if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
						matches = append(matches, fullPath+string(os.PathSeparator))
					}
				}
			}

			if len(matches) > 0 {
				return matches, nil
			}
		}
	}

	return []string{}, nil
}

func (a *FilePathAutocompleter) completeExtension(input string, _ *types.CompleteContext) ([]string, error) {
	// 获取常见文件扩展名
	commonExtensions := []string{
		".go", ".py", ".js", ".ts", ".java", ".cpp", ".c", ".h",
		".json", ".xml", ".yaml", ".yml", ".toml",
		".md", ".txt", ".rst",
		".html", ".css", ".scss", ".sass",
		".sql", ".sh", ".bat", ".ps1",
	}

	var matches []string
	for _, ext := range commonExtensions {
		if strings.HasPrefix(ext, input) {
			matches = append(matches, ext)
		}
	}

	return matches, nil
}

func (a *FilePathAutocompleter) completePattern(input string, _ *types.CompleteContext) ([]string, error) {
	// 支持通配符模式匹配
	dir := filepath.Dir(input)
	pattern := filepath.Base(input)

	items, err := a.getDirectoryItems(dir)
	if err != nil {
		return nil, err
	}

	var matches []string
	for _, item := range items {
		if matched, _ := filepath.Match(pattern, item); matched {
			matches = append(matches, filepath.Join(dir, item))
		}
	}

	return matches, nil
}

func (a *FilePathAutocompleter) completeGeneric(input string, _ *types.CompleteContext) ([]string, error) {
	// 通用补全：尝试文件和目录
	dir := filepath.Dir(input)
	base := filepath.Base(input)

	items, err := a.getDirectoryItems(dir)
	if err != nil {
		return nil, err
	}

	var matches []string
	for _, item := range items {
		if strings.HasPrefix(item, base) {
			fullPath := filepath.Join(dir, item)
			if info, err := os.Stat(fullPath); err == nil {
				if info.IsDir() {
					matches = append(matches, fullPath+string(os.PathSeparator))
				} else {
					matches = append(matches, fullPath)
				}
			}
		}
	}

	return matches, nil
}

// 辅助方法

func (a *FilePathAutocompleter) getMatchingItems(input string) []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var allItems []string

	// 从缓存中获取匹配项
	for _, items := range a.cache {
		for _, item := range items {
			if strings.Contains(item, input) {
				allItems = append(allItems, item)
			}
		}
	}

	// 去重和排序
	uniqueItems := removeDuplicates(allItems)
	sort.Strings(uniqueItems)

	return uniqueItems
}

func (a *FilePathAutocompleter) getDirectoryItems(dir string) ([]string, error) {
	// 检查缓存
	a.mu.RLock()
	if items, exists := a.cache[dir]; exists {
		a.mu.RUnlock()
		return items, nil
	}
	a.mu.RUnlock()

	// 扫描目录
	items, err := a.scanDirectory(dir)
	if err != nil {
		return nil, err
	}

	// 更新缓存
	a.mu.Lock()
	a.cache[dir] = items
	a.mu.Unlock()

	return items, nil
}

func (a *FilePathAutocompleter) scanDirectory(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var items []string
	for _, entry := range entries {
		name := entry.Name()

		// 跳过隐藏文件
		if strings.HasPrefix(name, ".") {
			continue
		}

		items = append(items, name)
	}

	return items, nil
}

func removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// CommandAutocompleter 命令自动补全器
type CommandAutocompleter struct {
	commands map[string]*CommandInfo
}

// CommandInfo 命令信息
type CommandInfo struct {
	Name        string
	Description string
	Aliases     []string
	Subcommands []string
	Options     []string
}

// NewCommandAutocompleter 创建命令自动补全器
func NewCommandAutocompleter() *CommandAutocompleter {
	return &CommandAutocompleter{
		commands: make(map[string]*CommandInfo),
	}
}

// RegisterCommand 注册命令
func (c *CommandAutocompleter) RegisterCommand(info *CommandInfo) {
	c.commands[info.Name] = info
}

// Complete 补全命令
func (c *CommandAutocompleter) Complete(input string) []string {
	var matches []string

	for name, info := range c.commands {
		if strings.HasPrefix(name, input) {
			matches = append(matches, name)
		}

		// 检查别名
		for _, alias := range info.Aliases {
			if strings.HasPrefix(alias, input) {
				matches = append(matches, alias)
			}
		}
	}

	sort.Strings(matches)
	return matches
}

// GetCommandInfo 获取命令信息
func (c *CommandAutocompleter) GetCommandInfo(command string) (*CommandInfo, bool) {
	info, exists := c.commands[command]
	return info, exists
}

// Suggestion 建议项
type Suggestion struct {
	Text        string
	Description string
	Type        string
	Icon        string
}

// SuggestionProvider 建议提供者接口
type SuggestionProvider interface {
	GetSuggestions(input string, context *types.CompleteContext) ([]Suggestion, error)
}

// CompositeSuggestionProvider 组合建议提供者
type CompositeSuggestionProvider struct {
	providers []SuggestionProvider
}

// NewCompositeSuggestionProvider 创建组合建议提供者
func NewCompositeSuggestionProvider(providers ...SuggestionProvider) *CompositeSuggestionProvider {
	return &CompositeSuggestionProvider{
		providers: providers,
	}
}

// GetSuggestions 获取建议
func (c *CompositeSuggestionProvider) GetSuggestions(input string, context *types.CompleteContext) ([]Suggestion, error) {
	var allSuggestions []Suggestion

	for _, provider := range c.providers {
		suggestions, err := provider.GetSuggestions(input, context)
		if err != nil {
			continue // 跳过出错的提供者
		}
		allSuggestions = append(allSuggestions, suggestions...)
	}

	// 去重和限制数量
	uniqueSuggestions := removeDuplicateSuggestions(allSuggestions)
	if len(uniqueSuggestions) > constants.DefaultMaxSuggestions {
		uniqueSuggestions = uniqueSuggestions[:constants.DefaultMaxSuggestions]
	}

	return uniqueSuggestions, nil
}

func removeDuplicateSuggestions(suggestions []Suggestion) []Suggestion {
	seen := make(map[string]bool)
	var result []Suggestion

	for _, suggestion := range suggestions {
		if !seen[suggestion.Text] {
			seen[suggestion.Text] = true
			result = append(result, suggestion)
		}
	}

	return result
}

// AutocompleterOptions 自动补全选项
type AutocompleterOptions struct {
	Enabled        bool
	MinChars       int
	MaxSuggestions int
	CacheSize      int
	Timeout        time.Duration
}
