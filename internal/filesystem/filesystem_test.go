package filesystem

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"code-context-generator/pkg/types"
)

func TestNewWalker(t *testing.T) {
	walker := NewWalker()
	if walker == nil {
		t.Fatal("NewWalker() returned nil")
	}

	// 检查类型
	if _, ok := walker.(*FileSystemWalker); !ok {
		t.Errorf("NewWalker() returned wrong type: %T", walker)
	}
}

func TestNewFileSystemWalker(t *testing.T) {
	options := types.WalkOptions{
		MaxDepth:        3,
		MaxFileSize:     1024 * 1024,
		ExcludePatterns: []string{"*.tmp"},
		IncludePatterns: []string{"*.go"},
		FollowSymlinks:  false,
	}

	walker := NewFileSystemWalker(options)
	if walker == nil {
		t.Fatal("NewFileSystemWalker() returned nil")
	}

	// 检查类型
	if _, ok := walker.(*FileSystemWalker); !ok {
		t.Errorf("NewFileSystemWalker() returned wrong type: %T", walker)
	}
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"test.go", ".go"},
		{"test.txt", ".txt"},
		{"test", ""},
		{"test.tar.gz", ".gz"},
		{"", ""},
		{".hidden", ""},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := GetFileExtension(tt.filename)
			if result != tt.expected {
				t.Errorf("GetFileExtension(%q) = %q, want %q", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestIsHiddenFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{".hidden", true},
		{"normal.txt", false},
		{"", false},
		{"..", true},
		{".git", true},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := IsHiddenFile(tt.filename)
			if result != tt.expected {
				t.Errorf("IsHiddenFile(%q) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestIsDirectory(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "test_dir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "test_file")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	tests := []struct {
		path     string
		expected bool
	}{
		{tempDir, true},
		{tempFile.Name(), false},
		{"nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := IsDirectory(tt.path)
			if result != tt.expected {
				t.Errorf("IsDirectory(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetFileSize(t *testing.T) {
	// 创建临时文件
	tempFile, err := os.CreateTemp("", "test_file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// 写入测试数据
	testData := []byte("Hello, World!")
	if _, err := tempFile.Write(testData); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	size, err := GetFileSize(tempFile.Name())
	if err != nil {
		t.Fatalf("GetFileSize() error = %v", err)
	}

	if size != int64(len(testData)) {
		t.Errorf("GetFileSize() = %v, want %v", size, len(testData))
	}
}

func TestGetFileModTime(t *testing.T) {
	// 创建临时文件
	tempFile, err := os.CreateTemp("", "test_file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	modTime, err := GetFileModTime(tempFile.Name())
	if err != nil {
		t.Fatalf("GetFileModTime() error = %v", err)
	}

	// 检查时间是否合理（应该在过去1分钟内）
	now := time.Now()
	if modTime.After(now) || modTime.Before(now.Add(-time.Minute)) {
		t.Errorf("GetFileModTime() = %v, expected recent time", modTime)
	}
}

func TestCreateDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_create")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	newDir := filepath.Join(tempDir, "new_directory")
	err = CreateDirectory(newDir)
	if err != nil {
		t.Fatalf("CreateDirectory() error = %v", err)
	}

	// 检查目录是否存在
	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		t.Errorf("CreateDirectory() directory was not created")
	}

	// 测试创建已存在的目录（应该不报错）
	err = CreateDirectory(newDir)
	if err != nil {
		t.Errorf("CreateDirectory() failed for existing directory: %v", err)
	}
}

func TestRemoveDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_remove")
	if err != nil {
		t.Fatal(err)
	}

	err = RemoveDirectory(tempDir)
	if err != nil {
		t.Fatalf("RemoveDirectory() error = %v", err)
	}

	// 检查目录是否被删除
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Errorf("RemoveDirectory() directory still exists")
	}
}

func TestCopyFile(t *testing.T) {
	// 创建源文件
	srcFile, err := os.CreateTemp("", "src_file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(srcFile.Name())

	// 写入测试数据
	testData := []byte("Hello, Copy Test!")
	if _, err := srcFile.Write(testData); err != nil {
		t.Fatal(err)
	}
	srcFile.Close()

	// 创建目标文件路径
	dstFile, err := os.CreateTemp("", "dst_file")
	if err != nil {
		t.Fatal(err)
	}
	dstPath := dstFile.Name()
	dstFile.Close()
	defer os.Remove(dstPath)

	// 复制文件
	err = CopyFile(srcFile.Name(), dstPath)
	if err != nil {
		t.Fatalf("CopyFile() error = %v", err)
	}

	// 验证内容
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(dstContent) != string(testData) {
		t.Errorf("CopyFile() content mismatch: got %q, want %q", string(dstContent), string(testData))
	}
}

func TestMoveFile(t *testing.T) {
	// 创建源文件
	srcFile, err := os.CreateTemp("", "src_file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(srcFile.Name())

	// 写入测试数据
	testData := []byte("Hello, Move Test!")
	if _, err := srcFile.Write(testData); err != nil {
		t.Fatal(err)
	}
	srcFile.Close()

	// 创建目标目录
	tempDir, err := os.MkdirTemp("", "move_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	dstPath := filepath.Join(tempDir, "moved_file.txt")

	// 移动文件
	err = MoveFile(srcFile.Name(), dstPath)
	if err != nil {
		t.Fatalf("MoveFile() error = %v", err)
	}

	// 验证源文件不存在
	if _, err := os.Stat(srcFile.Name()); !os.IsNotExist(err) {
		t.Errorf("MoveFile() source file still exists")
	}

	// 验证目标文件存在且内容正确
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(dstContent) != string(testData) {
		t.Errorf("MoveFile() content mismatch: got %q, want %q", string(dstContent), string(testData))
	}
}

func TestGetDirectorySize(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "size_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试文件
	testFiles := []struct {
		name string
		size int
	}{
		{"file1.txt", 100},
		{"file2.txt", 200},
		{"subdir/file3.txt", 150},
	}

	totalSize := 0
	for _, tf := range testFiles {
		filePath := filepath.Join(tempDir, tf.name)
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}

		data := make([]byte, tf.size)
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			t.Fatal(err)
		}
		totalSize += tf.size
	}

	size, err := GetDirectorySize(tempDir)
	if err != nil {
		t.Fatalf("GetDirectorySize() error = %v", err)
	}

	if size != int64(totalSize) {
		t.Errorf("GetDirectorySize() = %v, want %v", size, totalSize)
	}
}

func TestGetDirectoryFileCount(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "count_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试文件
	testFiles := []string{
		"file1.txt",
		"file2.txt",
		"subdir/file3.txt",
		"subdir/nested/file4.txt",
	}

	expectedCount := len(testFiles)
	for _, tf := range testFiles {
		filePath := filepath.Join(tempDir, tf)
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	count, err := GetDirectoryFileCount(tempDir)
	if err != nil {
		t.Fatalf("GetDirectoryFileCount() error = %v", err)
	}

	if count != expectedCount {
		t.Errorf("GetDirectoryFileCount() = %v, want %v", count, expectedCount)
	}
}

func TestFileSystemWalker_GetFileInfo(t *testing.T) {
	walker := NewWalker()
	fsWalker, ok := walker.(*FileSystemWalker)
	if !ok {
		t.Fatal("NewWalker() did not return *FileSystemWalker")
	}
	
	// 设置配置以启用元信息
	config := &types.Config{
		Output: types.OutputConfig{
			IncludeMetadata: true,
		},
	}
	fsWalker.SetConfig(config)

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "test_file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// 写入测试数据
	testData := []byte("Test file content")
	if _, err := tempFile.Write(testData); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	fileInfo, err := fsWalker.GetFileInfo(tempFile.Name())
	if err != nil {
		t.Fatalf("GetFileInfo() error = %v", err)
	}

	// 验证文件信息
	if fileInfo.Path != tempFile.Name() {
		t.Errorf("GetFileInfo() Path = %v, want %v", fileInfo.Path, tempFile.Name())
	}

	if fileInfo.Size != int64(len(testData)) {
		t.Errorf("GetFileInfo() Size = %v, want %v", fileInfo.Size, len(testData))
	}

	if fileInfo.Content != string(testData) {
		t.Errorf("GetFileInfo() Content = %v, want %v", fileInfo.Content, string(testData))
	}

	if fileInfo.IsDir {
		t.Error("GetFileInfo() IsDir should be false for file")
	}
}

func TestFileSystemWalker_GetFolderInfo(t *testing.T) {
	walker := &FileSystemWalker{}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "folder_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// 在目录中创建一些文件
	testFiles := []string{"file1.txt", "file2.go"}
	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	folderInfo, err := walker.GetFolderInfo(tempDir)
	if err != nil {
		t.Fatalf("GetFolderInfo() error = %v", err)
	}

	// 验证文件夹信息
	if folderInfo.Path != tempDir {
		t.Errorf("GetFolderInfo() Path = %v, want %v", folderInfo.Path, tempDir)
	}

	if len(folderInfo.Files) != len(testFiles) {
		t.Errorf("GetFolderInfo() Files count = %v, want %v", len(folderInfo.Files), len(testFiles))
	}
}

func TestFileSystemWalker_FilterFiles(t *testing.T) {
	walker := &FileSystemWalker{}

	files := []string{
		"/path/to/file1.txt",
		"/path/to/file2.go",
		"/path/to/test.log",
		"/path/to/config.yaml",
	}

	patterns := []string{"*.txt", "*.go"}

	filtered := walker.FilterFiles(files, patterns)

	expected := []string{
		"/path/to/file1.txt",
		"/path/to/file2.go",
	}

	if len(filtered) != len(expected) {
		t.Errorf("FilterFiles() returned %d files, want %d", len(filtered), len(expected))
	}

	for i, file := range filtered {
		if file != expected[i] {
			t.Errorf("FilterFiles()[%d] = %v, want %v", i, file, expected[i])
		}
	}
}

func TestFileSystemWalker_FilterBySize(t *testing.T) {
	walker := &FileSystemWalker{}

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "size_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// 写入测试数据
	testData := []byte("Test data for size filtering")
	if _, err := tempFile.Write(testData); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	tests := []struct {
		name     string
		maxSize  int64
		expected bool
	}{
		{"within limit", int64(len(testData) + 10), true},
		{"exact size", int64(len(testData)), true},
		{"exceeds limit", int64(len(testData) - 1), false},
		{"no limit", 0, true},
		{"negative limit", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := walker.FilterBySize(tempFile.Name(), tt.maxSize)
			if result != tt.expected {
				t.Errorf("FilterBySize() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFileSystemWalker_Walk(t *testing.T) {
	walker := NewWalker()

	// 创建临时目录结构
	tempDir, err := os.MkdirTemp("", "walk_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试文件结构
	testStructure := map[string]string{
		"file1.txt":              "content1",
		"file2.go":               "content2",
		"subdir/file3.txt":       "content3",
		"subdir/nested/file4.go": "content4",
	}

	for path, content := range testStructure {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// 测试基本遍历
	t.Run("basic walk", func(t *testing.T) {
		options := &types.WalkOptions{
			MaxDepth:        3,
			MaxFileSize:     1024 * 1024,
			ExcludePatterns: []string{},
			IncludePatterns: []string{},
			ShowHidden:      false,
		}

		contextData, err := walker.Walk(tempDir, options)
		if err != nil {
			t.Fatalf("Walk() error = %v", err)
		}

		// 验证文件数量
		if len(contextData.Files) != len(testStructure) {
			t.Errorf("Walk() found %d files, want %d", len(contextData.Files), len(testStructure))
		}

		// 验证文件夹数量
		if len(contextData.Folders) != 2 { // subdir 和 subdir/nested
			t.Errorf("Walk() found %d folders, want 2", len(contextData.Folders))
		}
	})

	// 测试深度限制
	t.Run("depth limit", func(t *testing.T) {
		options := &types.WalkOptions{
			MaxDepth:        1,
			MaxFileSize:     1024 * 1024,
			ExcludePatterns: []string{},
			IncludePatterns: []string{},
			ShowHidden:      false,
		}

		contextData, err := walker.Walk(tempDir, options)
		if err != nil {
			t.Fatalf("Walk() error = %v", err)
		}

		// 应该只找到根目录的文件
		expectedRootFiles := 2 // file1.txt 和 file2.go
		if len(contextData.Files) != expectedRootFiles {
			t.Errorf("Walk() with depth limit found %d files, want %d", len(contextData.Files), expectedRootFiles)
		}
	})

	// 测试包含模式
	t.Run("include patterns", func(t *testing.T) {
		options := &types.WalkOptions{
			MaxDepth:        3,
			MaxFileSize:     1024 * 1024,
			ExcludePatterns: []string{},
			IncludePatterns: []string{"*.txt"},
			ShowHidden:      false,
		}

		contextData, err := walker.Walk(tempDir, options)
		if err != nil {
			t.Fatalf("Walk() error = %v", err)
		}

		// 应该只找到.txt文件
		expectedTxtFiles := 2 // file1.txt 和 subdir/file3.txt
		if len(contextData.Files) != expectedTxtFiles {
			t.Errorf("Walk() with include patterns found %d files, want %d", len(contextData.Files), expectedTxtFiles)
		}
	})

	// 测试排除模式
	t.Run("exclude patterns", func(t *testing.T) {
		options := &types.WalkOptions{
			MaxDepth:        3,
			MaxFileSize:     1024 * 1024,
			ExcludePatterns: []string{"*.go"},
			IncludePatterns: []string{},
			ShowHidden:      false,
		}

		contextData, err := walker.Walk(tempDir, options)
		if err != nil {
			t.Fatalf("Walk() error = %v", err)
		}

		// 应该只找到非.go文件
		expectedNonGoFiles := 2 // file1.txt 和 subdir/file3.txt
		if len(contextData.Files) != expectedNonGoFiles {
			t.Errorf("Walk() with exclude patterns found %d files, want %d", len(contextData.Files), expectedNonGoFiles)
		}
	})

	// 测试大小限制
	t.Run("size limit", func(t *testing.T) {
		options := &types.WalkOptions{
			MaxDepth:        3,
			MaxFileSize:     5, // 很小的限制（小于8字节）
			ExcludePatterns: []string{},
			IncludePatterns: []string{},
			ShowHidden:      false,
		}

		contextData, err := walker.Walk(tempDir, options)
		if err != nil {
			t.Fatalf("Walk() error = %v", err)
		}

		// 应该没有找到文件（所有文件都超过5字节）
		if len(contextData.Files) != 0 {
			t.Errorf("Walk() with size limit found %d files, want 0", len(contextData.Files))
		}
	})
}

func TestFileSystemWalker_shouldIncludeFile(t *testing.T) {
	walker := &FileSystemWalker{}

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "include_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// 写入测试数据
	testData := []byte("Test inclusion")
	if _, err := tempFile.Write(testData); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	tests := []struct {
		name     string
		options  *types.WalkOptions
		expected bool
	}{
		{
			name: "include all",
			options: &types.WalkOptions{
				MaxFileSize:     1024 * 1024,
				IncludePatterns: []string{},
				ExcludePatterns: []string{},
				ShowHidden:      false,
			},
			expected: true,
		},
		{
			name: "exclude by pattern",
			options: &types.WalkOptions{
				MaxFileSize:     1024 * 1024,
				IncludePatterns: []string{},
				ExcludePatterns: []string{"*.tmp"},
				ShowHidden:      false,
			},
			expected: true, // 我们的测试文件不是.tmp
		},
		{
			name: "include by pattern",
			options: &types.WalkOptions{
				MaxFileSize:     1024 * 1024,
				IncludePatterns: []string{"*.tmp"},
				ExcludePatterns: []string{},
				ShowHidden:      false,
			},
			expected: false, // 我们的测试文件不是.tmp
		},
		{
			name: "size exceeded",
			options: &types.WalkOptions{
				MaxFileSize:     5, // 小于文件大小
				IncludePatterns: []string{},
				ExcludePatterns: []string{},
				ShowHidden:      false,
			},
			expected: false,
		},
		{
			name: "hidden file",
			options: &types.WalkOptions{
				MaxFileSize:     1024 * 1024,
				IncludePatterns: []string{},
				ExcludePatterns: []string{},
				ShowHidden:      false,
			},
			expected: true, // 我们的测试文件不是隐藏文件
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := walker.shouldIncludeFile(tempFile.Name(), "", tt.options)
			if result != tt.expected {
				t.Errorf("shouldIncludeFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// 基准测试
func BenchmarkGetFileExtension(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetFileExtension("test.file.extension.txt")
	}
}

func BenchmarkIsHiddenFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsHiddenFile(".hidden_file")
	}
}

func BenchmarkFilterFiles(b *testing.B) {
	walker := &FileSystemWalker{}
	files := []string{
		"file1.txt", "file2.go", "file3.log", "file4.yaml",
		"file5.json", "file6.md", "file7.py", "file8.rs",
	}
	patterns := []string{"*.txt", "*.go", "*.md"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		walker.FilterFiles(files, patterns)
	}
}
