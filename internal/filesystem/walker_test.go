package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"code-context-generator/pkg/types"
)

// TestShouldIncludeFile 测试文件包含/排除逻辑
func TestShouldIncludeFile(t *testing.T) {
	// 创建临时测试目录结构
	tempDir := t.TempDir()
	
	// 创建测试文件结构
	testFiles := []string{
		"file1.go",
		"file2.txt", 
		"file3.md",
		"test.log",
		".hidden.txt",
		"subdir/file4.go",
		"subdir/file5.txt",
		"subdir/nested/file6.go",
		".git/config",
		".git/HEAD",
		".git/hooks/pre-commit",
		"node_modules/package/index.js",
		"build/output.exe",
		"build/temp.obj",
		"docs/readme.md",
		"docs/guide/install.md",
		"vendor/lib/file.go",
		"coverage.out",
		"file.swp",
	}

	// 创建目录和文件
	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("创建目录失败 %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("创建文件失败 %s: %v", file, err)
		}
	}

	walker := &FileSystemWalker{}

	tests := []struct {
		name     string
		filePath string
		options  *types.WalkOptions
		expected bool
		desc     string
	}{
		{
			name:     "包含所有文件",
			filePath: filepath.Join(tempDir, "file1.go"),
			options: &types.WalkOptions{
				ExcludePatterns: []string{},
				IncludePatterns: []string{},
			},
			expected: true,
			desc:     "没有排除模式时应该包含文件",
		},
		{
			name:     "排除.git目录文件",
			filePath: filepath.Join(tempDir, ".git/config"),
			options: &types.WalkOptions{
				ExcludePatterns: []string{".git/"},
			},
			expected: false,
			desc:     "应该排除.git/目录下的文件",
		},
		{
			name:     "排除.git目录文件-HEAD",
			filePath: filepath.Join(tempDir, ".git/HEAD"),
			options: &types.WalkOptions{
				ExcludePatterns: []string{".git/"},
			},
			expected: false,
			desc:     "应该排除.git/目录下的HEAD文件",
		},
		{
			name:     "排除node_modules目录",
			filePath: filepath.Join(tempDir, "node_modules/package/index.js"),
			options: &types.WalkOptions{
				ExcludePatterns: []string{"node_modules/"},
			},
			expected: false,
			desc:     "应该排除node_modules/目录下的文件",
		},
		{
			name:     "排除build目录",
			filePath: filepath.Join(tempDir, "build/output.exe"),
			options: &types.WalkOptions{
				ExcludePatterns: []string{"build/"},
			},
			expected: false,
			desc:     "应该排除build/目录下的文件",
		},
		{
			name:     "排除vendor目录",
			filePath: filepath.Join(tempDir, "vendor/lib/file.go"),
			options: &types.WalkOptions{
				ExcludePatterns: []string{"vendor/"},
			},
			expected: false,
			desc:     "应该排除vendor/目录下的文件",
		},
		{
			name:     "按扩展名排除",
			filePath: filepath.Join(tempDir, "test.log"),
			options: &types.WalkOptions{
				ExcludePatterns: []string{"*.log"},
			},
			expected: false,
			desc:     "应该排除.log文件",
		},
		{
			name:     "按扩展名排除.swp",
			filePath: filepath.Join(tempDir, "file.swp"),
			options: &types.WalkOptions{
				ExcludePatterns: []string{"*.swp"},
			},
			expected: false,
			desc:     "应该排除.swp文件",
		},
		{
			name:     "包含模式匹配",
			filePath: filepath.Join(tempDir, "file1.go"),
			options: &types.WalkOptions{
				IncludePatterns: []string{"*.go"},
			},
			expected: true,
			desc:     "应该包含匹配包含模式的.go文件",
		},
		{
			name:     "包含模式不匹配",
			filePath: filepath.Join(tempDir, "file2.txt"),
			options: &types.WalkOptions{
				IncludePatterns: []string{"*.go"},
			},
			expected: false,
			desc:     "应该排除不匹配包含模式的.txt文件",
		},
		{
			name:     "目录包含模式",
			filePath: filepath.Join(tempDir, "docs/readme.md"),
			options: &types.WalkOptions{
				IncludePatterns: []string{"docs/*.md"},
			},
			expected: true,
			desc:     "应该包含匹配目录模式的docs/readme.md文件",
		},
		{
			name:     "目录包含模式不匹配",
			filePath: filepath.Join(tempDir, "file3.md"),
			options: &types.WalkOptions{
				IncludePatterns: []string{"docs/*.md"},
			},
			expected: false,
			desc:     "应该排除不匹配目录模式的根目录.md文件",
		},
		{
			name:     "多级目录包含模式",
			filePath: filepath.Join(tempDir, "docs/guide/install.md"),
			options: &types.WalkOptions{
				IncludePatterns: []string{"docs/**/*.md"},
			},
			expected: true,
			desc:     "应该包含匹配多级目录模式的文件",
		},
		{
			name:     "排除和包含模式组合",
			filePath: filepath.Join(tempDir, "subdir/file4.go"),
			options: &types.WalkOptions{
				ExcludePatterns: []string{"*.txt"},
				IncludePatterns: []string{"*.go"},
			},
			expected: true,
			desc:     "应该包含.go文件即使.txt被排除",
		},
		{
			name:     "Windows路径分隔符测试",
			filePath: filepath.Join(tempDir, "subdir\\file4.go"),
			options: &types.WalkOptions{
				IncludePatterns: []string{"subdir/*.go"},
			},
			expected: true,
			desc:     "应该正确处理Windows路径分隔符",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := walker.shouldIncludeFile(tt.filePath, tempDir, tt.options)
			if result != tt.expected {
				t.Errorf("%s: shouldIncludeFile() = %v, want %v\n描述: %s\n文件: %s\n选项: %+v", 
					tt.name, result, tt.expected, tt.desc, tt.filePath, tt.options)
			}
		})
	}
}

// TestWalkWithExcludePatterns 测试实际的Walk函数与排除模式
func TestWalkWithExcludePatterns(t *testing.T) {
	// 创建临时测试目录结构
	tempDir := t.TempDir()
	
	// 创建测试文件结构
	testFiles := []string{
		"main.go",
		"readme.md", 
		".git/config",
		".git/HEAD",
		"node_modules/lib/index.js",
		"build/output.exe",
		"vendor/lib/helper.go",
		"test.log",
		"coverage.out",
	}

	// 创建目录和文件
	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("创建目录失败 %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("创建文件失败 %s: %v", file, err)
		}
	}

	walker := NewWalker()
	
	options := &types.WalkOptions{
		ExcludePatterns: []string{
			".git/",
			"node_modules/", 
			"build/",
			"vendor/",
			"*.log",
			"*.out",
		},
		MaxDepth: 10,
	}

	result, err := walker.Walk(tempDir, options)
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// 验证结果
	if result == nil {
		t.Fatal("Walk() returned nil result")
	}

	// 检查应该包含的文件
	expectedFiles := []string{"main.go", "readme.md"}
	for _, expectedFile := range expectedFiles {
		found := false
		for _, file := range result.Files {
			if filepath.Base(file.Path) == expectedFile {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("应该包含文件 %s", expectedFile)
		}
	}

	// 检查应该排除的文件
	excludedFiles := []string{"config", "HEAD", "index.js", "output.exe", "helper.go", "test.log", "coverage.out"}
	for _, excludedFile := range excludedFiles {
		found := false
		for _, file := range result.Files {
			if filepath.Base(file.Path) == excludedFile {
				found = true
				break
			}
		}
		if found {
			t.Errorf("应该排除文件 %s", excludedFile)
		}
	}

	// 验证文件数量
	expectedFileCount := len(expectedFiles)
	if result.FileCount != expectedFileCount {
		t.Errorf("文件数量不匹配: got %d, want %d", result.FileCount, expectedFileCount)
	}

	t.Logf("测试通过: 找到 %d 个文件, 期望 %d 个文件", result.FileCount, expectedFileCount)
	for _, file := range result.Files {
		t.Logf("包含文件: %s", file.Path)
	}
}

// TestWalkWithIncludePatterns 测试实际的Walk函数与包含模式
func TestWalkWithIncludePatterns(t *testing.T) {
	// 创建临时测试目录结构
	tempDir := t.TempDir()
	
	// 创建测试文件结构
	testFiles := []string{
		"main.go",
		"helper.go", 
		"readme.md",
		"config.json",
		"test.txt",
		"docs/guide.md",
		"docs/api.txt",
		"src/utils.js",
		"src/styles.css",
	}

	// 创建目录和文件
	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("创建目录失败 %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("创建文件失败 %s: %v", file, err)
		}
	}

	walker := NewWalker()
	
	options := &types.WalkOptions{
		IncludePatterns: []string{
			"*.go",
			"docs/*.md",
		},
		MaxDepth: 10,
	}

	result, err := walker.Walk(tempDir, options)
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// 验证结果
	if result == nil {
		t.Fatal("Walk() returned nil result")
	}

	// 检查应该包含的文件
	expectedFiles := []string{"main.go", "helper.go", "guide.md"}
	for _, expectedFile := range expectedFiles {
		found := false
		for _, file := range result.Files {
			if filepath.Base(file.Path) == expectedFile {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("应该包含文件 %s", expectedFile)
		}
	}

	// 检查应该排除的文件
	excludedFiles := []string{"readme.md", "config.json", "test.txt", "api.txt", "utils.js", "styles.css"}
	for _, excludedFile := range excludedFiles {
		found := false
		for _, file := range result.Files {
			if filepath.Base(file.Path) == excludedFile {
				found = true
				break
			}
		}
		if found {
			t.Errorf("应该排除文件 %s", excludedFile)
		}
	}

	t.Logf("测试通过: 找到 %d 个文件", result.FileCount)
	for _, file := range result.Files {
		t.Logf("包含文件: %s", file.Path)
	}
}