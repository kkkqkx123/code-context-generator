package selector

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"code-context-generator/pkg/types"
)

// TestNewSelector 测试创建新的选择器
func TestNewSelector(t *testing.T) {
	// 测试默认配置
	selector := NewSelector(nil)
	if selector == nil {
		t.Fatal("NewSelector returned nil")
	}

	// 测试自定义配置
	config := &types.Config{
		Filters: types.FiltersConfig{
			MaxFileSize: "10MB",
		},
	}
	selector = NewSelector(config)
	if selector == nil {
		t.Fatal("NewSelector with config returned nil")
	}
}

// TestFileSelector_SelectFiles 测试文件选择功能
func TestFileSelector_SelectFiles(t *testing.T) {
	// 创建临时目录结构
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 创建测试文件
	testFiles := []string{
		filepath.Join(tempDir, "test.txt"),
		filepath.Join(tempDir, "main.go"),
		filepath.Join(subDir, "subtest.txt"),
		filepath.Join(subDir, "hidden.txt"),
	}

	for _, file := range testFiles {
		if err := os.WriteFile(file, []byte("test content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// 创建隐藏文件
	hiddenFile := filepath.Join(tempDir, ".hidden.txt")
	if err := os.WriteFile(hiddenFile, []byte("hidden content"), 0644); err != nil {
		t.Fatal(err)
	}

	config := &types.Config{
		Filters: types.FiltersConfig{
			MaxFileSize: "1MB",
		},
	}
	selector := NewSelector(config).(*FileSelector)

	tests := []struct {
		name           string
		rootPath       string
		options        *types.SelectOptions
		expectedMin    int
		expectedMax    int
		shouldContain  []string
		shouldNotContain []string
	}{
		{
			name:     "select all files recursively",
			rootPath: tempDir,
			options: &types.SelectOptions{
				Recursive:       true,
				IncludePatterns: []string{},
				ExcludePatterns: []string{},
				MaxDepth:        0,
				ShowHidden:      false,
				SortBy:          "name",
			},
			expectedMin: 4,
			expectedMax: 4,
			shouldContain: []string{"test.txt", "main.go", "subtest.txt", "hidden.txt"},
		},
		{
			name:     "select files with include pattern",
			rootPath: tempDir,
			options: &types.SelectOptions{
				Recursive:       true,
				IncludePatterns: []string{"*.txt"},
				ExcludePatterns: []string{},
				MaxDepth:        0,
				ShowHidden:      false,
				SortBy:          "name",
			},
			expectedMin: 3,
			expectedMax: 3,
			shouldContain: []string{"test.txt", "subtest.txt", "hidden.txt"},
			shouldNotContain: []string{"main.go"},
		},
		{
			name:     "select files with exclude pattern",
			rootPath: tempDir,
			options: &types.SelectOptions{
				Recursive:       true,
				IncludePatterns: []string{},
				ExcludePatterns: []string{"*.go"},
				MaxDepth:        0,
				ShowHidden:      false,
				SortBy:          "name",
			},
			expectedMin: 3,
			expectedMax: 3,
			shouldContain: []string{"test.txt", "subtest.txt", "hidden.txt"},
			shouldNotContain: []string{"main.go"},
		},
		{
			name:     "select files with max depth",
			rootPath: tempDir,
			options: &types.SelectOptions{
				Recursive:       true,
				IncludePatterns: []string{},
				ExcludePatterns: []string{},
				MaxDepth:        1,
				ShowHidden:      false,
				SortBy:          "name",
			},
			expectedMin: 2,
			expectedMax: 2,
			shouldContain: []string{"test.txt", "main.go"},
		},
		{
			name:     "select files with show hidden",
			rootPath: tempDir,
			options: &types.SelectOptions{
				Recursive:       true,
				IncludePatterns: []string{},
				ExcludePatterns: []string{},
				MaxDepth:        0,
				ShowHidden:      true,
				SortBy:          "name",
			},
			expectedMin: 5,
			expectedMax: 5,
			shouldContain: []string{".hidden.txt", "test.txt", "main.go", "subtest.txt", "hidden.txt"},
		},
		{
			name:     "select files with nil options (default)",
			rootPath: tempDir,
			options:  nil,
			expectedMin: 4,
			expectedMax: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := selector.SelectFiles(tt.rootPath, tt.options)
			if err != nil {
				t.Errorf("SelectFiles() error = %v", err)
				return
			}

			if len(files) < tt.expectedMin || len(files) > tt.expectedMax {
				t.Errorf("SelectFiles() got %d files, expected between %d and %d", len(files), tt.expectedMin, tt.expectedMax)
			}

			// 检查应该包含的文件
			for _, shouldContain := range tt.shouldContain {
				found := false
				for _, file := range files {
					if filepath.Base(file) == shouldContain {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("SelectFiles() should contain file %s", shouldContain)
				}
			}

			// 检查不应该包含的文件
			for _, shouldNotContain := range tt.shouldNotContain {
				found := false
				for _, file := range files {
					if filepath.Base(file) == shouldNotContain {
						found = true
						break
					}
				}
				if found {
					t.Errorf("SelectFiles() should not contain file %s", shouldNotContain)
				}
			}
		})
	}
}

// TestFileSelector_SelectFolders 测试文件夹选择功能
func TestFileSelector_SelectFolders(t *testing.T) {
	// 创建临时目录结构
	tempDir := t.TempDir()
	subDir1 := filepath.Join(tempDir, "subdir1")
	subDir2 := filepath.Join(tempDir, "subdir2")
	hiddenDir := filepath.Join(tempDir, ".hidden")

	for _, dir := range []string{subDir1, subDir2, hiddenDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}

	config := &types.Config{}
	selector := NewSelector(config).(*FileSelector)

	tests := []struct {
		name           string
		rootPath       string
		options        *types.SelectOptions
		expectedMin    int
		expectedMax    int
		shouldContain  []string
		shouldNotContain []string
	}{
		{
			name:     "select all folders recursively",
			rootPath: tempDir,
			options: &types.SelectOptions{
				Recursive:       true,
				IncludePatterns: []string{},
				ExcludePatterns: []string{},
				MaxDepth:        0,
				ShowHidden:      false,
				SortBy:          "name",
			},
			expectedMin: 2,
			expectedMax: 2,
			shouldContain: []string{"subdir1", "subdir2"},
		},
		{
			name:     "select folders with show hidden",
			rootPath: tempDir,
			options: &types.SelectOptions{
				Recursive:       true,
				IncludePatterns: []string{},
				ExcludePatterns: []string{},
				MaxDepth:        0,
				ShowHidden:      true,
				SortBy:          "name",
			},
			expectedMin: 3,
			expectedMax: 3,
			shouldContain: []string{"subdir1", "subdir2", ".hidden"},
		},
		{
			name:     "select folders with max depth",
			rootPath: tempDir,
			options: &types.SelectOptions{
				Recursive:       true,
				IncludePatterns: []string{},
				ExcludePatterns: []string{},
				MaxDepth:        1,
				ShowHidden:      false,
				SortBy:          "name",
			},
			expectedMin: 2,
			expectedMax: 2,
		},
		{
			name:     "select folders with include pattern",
			rootPath: tempDir,
			options: &types.SelectOptions{
				Recursive:       true,
				IncludePatterns: []string{"sub*"},
				ExcludePatterns: []string{},
				MaxDepth:        0,
				ShowHidden:      false,
				SortBy:          "name",
			},
			expectedMin: 2,
			expectedMax: 2,
			shouldContain: []string{"subdir1", "subdir2"},
		},
		{
			name:     "select folders with nil options (default)",
			rootPath: tempDir,
			options:  nil,
			expectedMin: 2,
			expectedMax: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			folders, err := selector.SelectFolders(tt.rootPath, tt.options)
			if err != nil {
				t.Errorf("SelectFolders() error = %v", err)
				return
			}

			if len(folders) < tt.expectedMin || len(folders) > tt.expectedMax {
				t.Errorf("SelectFolders() got %d folders, expected between %d and %d", len(folders), tt.expectedMin, tt.expectedMax)
			}

			// 检查应该包含的文件夹
			for _, shouldContain := range tt.shouldContain {
				found := false
				for _, folder := range folders {
					if filepath.Base(folder) == shouldContain {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("SelectFolders() should contain folder %s", shouldContain)
				}
			}

			// 检查不应该包含的文件夹
			for _, shouldNotContain := range tt.shouldNotContain {
				found := false
				for _, folder := range folders {
					if filepath.Base(folder) == shouldNotContain {
						found = true
						break
					}
				}
				if found {
					t.Errorf("SelectFolders() should not contain folder %s", shouldNotContain)
				}
			}
		})
	}
}

// TestFileSelector_InteractiveSelect 测试交互式选择功能
func TestFileSelector_InteractiveSelect(t *testing.T) {
	selector := NewSelector(nil).(*FileSelector)

	tests := []struct {
		name     string
		items    []string
		prompt   string
		expected int
	}{
		{
			name:     "interactive select with items",
			items:    []string{"item1", "item2", "item3"},
			prompt:   "Select items",
			expected: 3,
		},
		{
			name:     "interactive select with empty items",
			items:    []string{},
			prompt:   "Select items",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := selector.InteractiveSelect(tt.items, tt.prompt)
			if err != nil {
				t.Errorf("InteractiveSelect() error = %v", err)
				return
			}

			if len(result) != tt.expected {
				t.Errorf("InteractiveSelect() got %d items, expected %d", len(result), tt.expected)
			}
		})
	}
}

// TestFileSelector_FilterItems 测试项目过滤功能
func TestFileSelector_FilterItems(t *testing.T) {
	selector := NewSelector(nil).(*FileSelector)

	items := []string{
		"test.txt",
		"main.go",
		"README.md",
		"config.yaml",
		"test_backup.txt",
	}

	tests := []struct {
		name     string
		items    []string
		filter   string
		expected int
		contains []string
	}{
		{
			name:     "filter with matching pattern",
			items:    items,
			filter:   "test",
			expected: 2,
			contains: []string{"test.txt", "test_backup.txt"},
		},
		{
			name:     "filter with no match",
			items:    items,
			filter:   "nomatch",
			expected: 0,
		},
		{
			name:     "filter with empty filter",
			items:    items,
			filter:   "",
			expected: 5,
		},
		{
			name:     "filter with case insensitive",
			items:    items,
			filter:   "TEST",
			expected: 2,
			contains: []string{"test.txt", "test_backup.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := selector.FilterItems(tt.items, tt.filter)

			if len(result) != tt.expected {
				t.Errorf("FilterItems() got %d items, expected %d", len(result), tt.expected)
			}

			// 检查应该包含的项目
			for _, shouldContain := range tt.contains {
				found := false
				for _, item := range result {
					if item == shouldContain {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("FilterItems() should contain item %s", shouldContain)
				}
			}
		})
	}
}

// TestFileSelector_SortItems 测试项目排序功能
func TestFileSelector_SortItems(t *testing.T) {
	// 创建临时文件用于测试按大小和修改时间排序
	tempDir := t.TempDir()
	
	// 创建测试文件
	files := []string{
		filepath.Join(tempDir, "a.txt"),
		filepath.Join(tempDir, "c.txt"),
		filepath.Join(tempDir, "b.txt"),
	}

	for i, file := range files {
		content := []byte("content")
		if i == 1 {
			content = []byte("larger content for testing")
		}
		if err := os.WriteFile(file, content, 0644); err != nil {
			t.Fatal(err)
		}
		// 修改文件时间
		if i == 2 {
			time.Sleep(10 * time.Millisecond) // 确保时间不同
		}
	}

	selector := NewSelector(nil).(*FileSelector)

	tests := []struct {
		name     string
		items    []string
		sortBy   string
		validate func([]string) bool
	}{
		{
			name:   "sort by name",
			items:  files,
			sortBy: "name",
			validate: func(result []string) bool {
				return filepath.Base(result[0]) == "a.txt" &&
					filepath.Base(result[1]) == "b.txt" &&
					filepath.Base(result[2]) == "c.txt"
			},
		},
		{
			name:   "sort by size",
			items:  files,
			sortBy: "size",
			validate: func(result []string) bool {
				// a.txt 和 b.txt 大小相同，c.txt 更大
				return len(result) == 3
			},
		},
		{
			name:   "sort by modified time",
			items:  files,
			sortBy: "modified",
			validate: func(result []string) bool {
				return len(result) == 3
			},
		},
		{
			name:   "sort by unknown (defaults to name)",
			items:  files,
			sortBy: "unknown",
			validate: func(result []string) bool {
				return len(result) == 3
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := selector.SortItems(tt.items, tt.sortBy)

			if len(result) != len(tt.items) {
				t.Errorf("SortItems() got %d items, expected %d", len(result), len(tt.items))
				return
			}

			if !tt.validate(result) {
				t.Errorf("SortItems() validation failed for sortBy=%s", tt.sortBy)
			}
		})
	}
}

// TestPatternMatcher 测试模式匹配器
func TestPatternMatcher(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		path     string
		expected bool
	}{
		{
			name:     "match single pattern",
			patterns: []string{"*.txt"},
			path:     "test.txt",
			expected: true,
		},
		{
			name:     "match multiple patterns",
			patterns: []string{"*.txt", "*.go"},
			path:     "main.go",
			expected: true,
		},
		{
			name:     "no match",
			patterns: []string{"*.txt"},
			path:     "main.go",
			expected: false,
		},
		{
			name:     "match with wildcard",
			patterns: []string{"test*"},
			path:     "test123.txt",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := NewPatternMatcher(tt.patterns)
			result := pm.Match(tt.path)

			if result != tt.expected {
				t.Errorf("Match() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestPatternMatcher_MatchAny 测试模式匹配器的MatchAny方法
func TestPatternMatcher_MatchAny(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		path     string
		expected bool
	}{
		{
			name:     "match with wildcard",
			patterns: []string{"test*"},
			path:     "test123.txt",
			expected: true,
		},
		{
			name:     "match with contains",
			patterns: []string{"test"},
			path:     "mytestfile.txt",
			expected: true,
		},
		{
			name:     "no match",
			patterns: []string{"nomatch"},
			path:     "test.txt",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := NewPatternMatcher(tt.patterns)
			result := pm.MatchAny(tt.path)

			if result != tt.expected {
				t.Errorf("MatchAny() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestDefaultGetFileType 测试默认文件类型获取函数
func TestDefaultGetFileType(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "get go file type",
			path:     "main.go",
			expected: "go",
		},
		{
			name:     "get txt file type",
			path:     "test.txt",
			expected: "txt",
		},
		{
			name:     "get file type without extension",
			path:     "Makefile",
			expected: "unknown",
		},
		{
			name:     "get file type with multiple extensions",
			path:     "archive.tar.gz",
			expected: "gz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultGetFileType(tt.path)

			if result != tt.expected {
				t.Errorf("DefaultGetFileType() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestParseFileSize 测试文件大小解析功能
func TestParseFileSize(t *testing.T) {
	tests := []struct {
		name     string
		sizeStr  string
		expected int64
		wantErr  bool
	}{
		{
			name:     "parse bytes",
			sizeStr:  "1024",
			expected: 1024,
			wantErr:  false,
		},
		{
			name:     "parse KB",
			sizeStr:  "1KB",
			expected: 1024,
			wantErr:  false,
		},
		{
			name:     "parse MB",
			sizeStr:  "1MB",
			expected: 1024 * 1024,
			wantErr:  false,
		},
		{
			name:     "parse GB",
			sizeStr:  "1GB",
			expected: 1024 * 1024 * 1024,
			wantErr:  false,
		},
		{
			name:     "parse with space",
			sizeStr:  "1 MB",
			expected: 1024 * 1024,
			wantErr:  false,
		},
		{
			name:     "parse with lowercase",
			sizeStr:  "1mb",
			expected: 1024 * 1024,
			wantErr:  false,
		},
		{
			name:     "parse invalid format",
			sizeStr:  "invalid",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "parse with decimal",
			sizeStr:  "1.5MB",
			expected: int64(1.5 * 1024 * 1024),
			wantErr:  false,
		},
		{
			name:     "parse unsupported unit",
			sizeStr:  "1TB",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseFileSize(tt.sizeStr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseFileSize() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("parseFileSize() unexpected error = %v", err)
				}
				if result != tt.expected {
					t.Errorf("parseFileSize() = %v, expected %v", result, tt.expected)
				}
			}
		})
	}
}

// TestDefaultGetFileIcon 测试默认文件图标获取函数
func TestDefaultGetFileIcon(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "get go file icon",
			path:     "main.go",
			expected: "🐹",
		},
		{
			name:     "get python file icon",
			path:     "script.py",
			expected: "🐍",
		},
		{
			name:     "get javascript file icon",
			path:     "app.js",
			expected: "📜",
		},
		{
			name:     "get markdown file icon",
			path:     "README.md",
			expected: "📝",
		},
		{
			name:     "get default file icon",
			path:     "unknown.xyz",
			expected: "📄",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultGetFileIcon(tt.path)

			if result != tt.expected {
				t.Errorf("DefaultGetFileIcon() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestGetFileInfo 测试文件信息获取功能
func TestGetFileInfo(t *testing.T) {
	// 创建临时文件
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	content := []byte("test content for file info")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	result, err := GetFileInfo(testFile)
	if err != nil {
		t.Errorf("GetFileInfo() error = %v", err)
		return
	}

	if result.Path != testFile {
		t.Errorf("GetFileInfo() Path = %v, expected %v", result.Path, testFile)
	}

	if result.Name != "test.txt" {
		t.Errorf("GetFileInfo() Name = %v, expected test.txt", result.Name)
	}

	if result.Size != int64(len(content)) {
		t.Errorf("GetFileInfo() Size = %v, expected %v", result.Size, len(content))
	}

	if result.IsDir {
		t.Errorf("GetFileInfo() IsDir = true, expected false")
	}

	if result.IsHidden {
		t.Errorf("GetFileInfo() IsHidden = true, expected false")
	}

	if result.Type != "txt" {
		t.Errorf("GetFileInfo() Type = %v, expected txt", result.Type)
	}
}

// TestGetDirectoryContents 测试目录内容获取功能
func TestGetDirectoryContents(t *testing.T) {
	// 创建临时目录结构
	tempDir := t.TempDir()
	
	// 创建测试文件
	testFiles := []string{
		filepath.Join(tempDir, "file1.txt"),
		filepath.Join(tempDir, "file2.go"),
		filepath.Join(tempDir, ".hidden"),
	}

	for _, file := range testFiles {
		if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name        string
		path        string
		showHidden  bool
		expectedMin int
		expectedMax int
	}{
		{
			name:        "get directory contents without hidden",
			path:        tempDir,
			showHidden:  false,
			expectedMin: 2,
			expectedMax: 2,
		},
		{
			name:        "get directory contents with hidden",
			path:        tempDir,
			showHidden:  true,
			expectedMin: 3,
			expectedMax: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetDirectoryContents(tt.path, tt.showHidden)
			if err != nil {
				t.Errorf("GetDirectoryContents() error = %v", err)
				return
			}

			if len(result) < tt.expectedMin || len(result) > tt.expectedMax {
				t.Errorf("GetDirectoryContents() got %d items, expected between %d and %d", len(result), tt.expectedMin, tt.expectedMax)
			}
		})
	}
}

// TestFileSelector_shouldIncludeFile 测试文件包含逻辑
func TestFileSelector_shouldIncludeFile(t *testing.T) {
	// 创建临时文件用于测试
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	config := &types.Config{
		Filters: types.FiltersConfig{
			MaxFileSize: "1MB",
		},
	}
	selector := NewSelector(config).(*FileSelector)

	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		path     string
		info     os.FileInfo
		options  *types.SelectOptions
		expected bool
	}{
		{
			name:    "include normal file",
			path:    testFile,
			info:    info,
			options: &types.SelectOptions{
				ShowHidden:      false,
				IncludePatterns: []string{},
				ExcludePatterns: []string{},
			},
			expected: true,
		},
		{
			name:    "exclude hidden file",
			path:    filepath.Join(tempDir, ".hidden.txt"),
			info:    info, // 复用文件信息，实际测试中应该创建真实文件
			options: &types.SelectOptions{
				ShowHidden:      false,
				IncludePatterns: []string{},
				ExcludePatterns: []string{},
			},
			expected: false,
		},
		{
			name:    "include with matching pattern",
			path:    testFile,
			info:    info,
			options: &types.SelectOptions{
				ShowHidden:      false,
				IncludePatterns: []string{"*.txt"},
				ExcludePatterns: []string{},
			},
			expected: true,
		},
		{
			name:    "exclude with matching pattern",
			path:    testFile,
			info:    info,
			options: &types.SelectOptions{
				ShowHidden:      false,
				IncludePatterns: []string{},
				ExcludePatterns: []string{"*.txt"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := selector.shouldIncludeFile(tt.path, tt.info, tt.options)

			if result != tt.expected {
				t.Errorf("shouldIncludeFile() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestFileSelector_shouldIncludeFolder 测试文件夹包含逻辑
func TestFileSelector_shouldIncludeFolder(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()
	
	config := &types.Config{}
	selector := NewSelector(config).(*FileSelector)

	tests := []struct {
		name     string
		path     string
		options  *types.SelectOptions
		expected bool
	}{
		{
			name:    "include normal folder",
			path:    filepath.Join(tempDir, "normal"),
			options: &types.SelectOptions{
				ShowHidden:      false,
				IncludePatterns: []string{},
				ExcludePatterns: []string{},
			},
			expected: true,
		},
		{
			name:    "exclude hidden folder",
			path:    filepath.Join(tempDir, ".hidden"),
			options: &types.SelectOptions{
				ShowHidden:      false,
				IncludePatterns: []string{},
				ExcludePatterns: []string{},
			},
			expected: false,
		},
		{
			name:    "include with matching pattern",
			path:    filepath.Join(tempDir, "test_folder"),
			options: &types.SelectOptions{
				ShowHidden:      false,
				IncludePatterns: []string{"test_*"},
				ExcludePatterns: []string{},
			},
			expected: true,
		},
		{
			name:    "exclude with matching pattern",
			path:    filepath.Join(tempDir, "test_folder"),
			options: &types.SelectOptions{
				ShowHidden:      false,
				IncludePatterns: []string{},
				ExcludePatterns: []string{"test_*"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := selector.shouldIncludeFolder(tt.path, nil, tt.options)

			if result != tt.expected {
				t.Errorf("shouldIncludeFolder() = %v, expected %v", result, tt.expected)
			}
		})
	}
}