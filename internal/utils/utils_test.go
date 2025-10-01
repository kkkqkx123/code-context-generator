// Package utils 单元测试
package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestFileUtils 文件工具函数测试
func TestFileExists(t *testing.T) {
	// 测试存在的文件
	if !FileExists("utils.go") {
		t.Error("FileExists 应该返回 true 对于存在的文件")
	}

	// 测试不存在的文件
	if FileExists("nonexistent.go") {
		t.Error("FileExists 应该返回 false 对于不存在的文件")
	}
}

func TestDirectoryExists(t *testing.T) {
	// 测试存在的目录
	if !DirectoryExists(".") {
		t.Error("DirectoryExists 应该返回 true 对于存在的目录")
	}

	// 测试不存在的目录
	if DirectoryExists("nonexistent_dir") {
		t.Error("DirectoryExists 应该返回 false 对于不存在的目录")
	}

	// 测试文件而不是目录
	if DirectoryExists("utils.go") {
		t.Error("DirectoryExists 应该返回 false 对于文件")
	}
}

func TestGetFileHash(t *testing.T) {
	// 创建测试文件
	testFile := "test_hash.txt"
	content := "test content"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	defer os.Remove(testFile)

	// 测试获取文件哈希
	hash, err := GetFileHash(testFile)
	if err != nil {
		t.Errorf("GetFileHash 返回错误: %v", err)
	}
	if hash == "" {
		t.Error("GetFileHash 应该返回非空哈希值")
	}

	// 测试不存在的文件
	_, err = GetFileHash("nonexistent.txt")
	if err == nil {
		t.Error("GetFileHash 应该对不存在的文件返回错误")
	}
}

func TestGetFileSize(t *testing.T) {
	// 创建测试文件
	testFile := "test_size.txt"
	content := "test content"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	defer os.Remove(testFile)

	// 测试获取文件大小
	size, err := GetFileSize(testFile)
	if err != nil {
		t.Errorf("GetFileSize 返回错误: %v", err)
	}
	if size != int64(len(content)) {
		t.Errorf("GetFileSize 返回的大小不正确: 期望 %d, 实际 %d", len(content), size)
	}

	// 测试不存在的文件
	_, err = GetFileSize("nonexistent.txt")
	if err == nil {
		t.Error("GetFileSize 应该对不存在的文件返回错误")
	}
}

func TestGetFileModTime(t *testing.T) {
	// 创建测试文件
	testFile := "test_modtime.txt"
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	defer os.Remove(testFile)

	// 测试获取文件修改时间
	modTime, err := GetFileModTime(testFile)
	if err != nil {
		t.Errorf("GetFileModTime 返回错误: %v", err)
	}
	if modTime.IsZero() {
		t.Error("GetFileModTime 应该返回非零时间")
	}

	// 测试不存在的文件
	_, err = GetFileModTime("nonexistent.txt")
	if err == nil {
		t.Error("GetFileModTime 应该对不存在的文件返回错误")
	}
}

func TestIsTextFile(t *testing.T) {
	testCases := []struct {
		filename string
		expected bool
	}{
		{"test.txt", true},
		{"test.md", true},
		{"test.json", true},
		{"test.go", true},
		{"test.py", true},
		{"test.js", true},
		{"test.html", true},
		{"test.css", true},
		{"test.exe", false},
		{"test.bin", false},
		{"test.jpg", false},
		{"test.png", false},
		{"test.pdf", false},
		{"test", false},
		{"", false},
	}

	for _, tc := range testCases {
		result := IsTextFile(tc.filename)
		if result != tc.expected {
			t.Errorf("IsTextFile(%s) = %v, 期望 %v", tc.filename, result, tc.expected)
		}
	}
}

func TestIsBinaryFile(t *testing.T) {
	// IsBinaryFile 应该返回与 IsTextFile 相反的结果
	testFiles := []string{"test.txt", "test.exe", "test.jpg"}
	
	for _, filename := range testFiles {
		textResult := IsTextFile(filename)
		binaryResult := IsBinaryFile(filename)
		if textResult == binaryResult {
			t.Errorf("IsBinaryFile(%s) = %v, 应该与 IsTextFile 相反", filename, binaryResult)
		}
	}
}

// TestStringUtils 字符串工具函数测试
func TestTruncateString(t *testing.T) {
	testCases := []struct {
		input     string
		maxLength int
		expected  string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "he..."},
		{"hello", 3, "hel"},
		{"hello", 0, ""},
		{"", 5, ""},
	}

	for _, tc := range testCases {
		result := TruncateString(tc.input, tc.maxLength)
		if result != tc.expected {
			t.Errorf("TruncateString(%q, %d) = %q, 期望 %q", tc.input, tc.maxLength, result, tc.expected)
		}
	}
}

func TestPadString(t *testing.T) {
	testCases := []struct {
		input    string
		length   int
		padChar  rune
		expected string
	}{
		{"hello", 10, '-', "hello-----"},
		{"hello", 5, '-', "hello"},
		{"hello", 3, '-', "hello"},
		{"", 5, '-', "-----"},
	}

	for _, tc := range testCases {
		result := PadString(tc.input, tc.length, tc.padChar)
		if result != tc.expected {
			t.Errorf("PadString(%q, %d, %q) = %q, 期望 %q", tc.input, tc.length, tc.padChar, result, tc.expected)
		}
	}
}

func TestPadLeft(t *testing.T) {
	testCases := []struct {
		input    string
		length   int
		padChar  rune
		expected string
	}{
		{"hello", 10, '-', "-----hello"},
		{"hello", 5, '-', "hello"},
		{"hello", 3, '-', "hello"},
		{"", 5, '-', "-----"},
	}

	for _, tc := range testCases {
		result := PadLeft(tc.input, tc.length, tc.padChar)
		if result != tc.expected {
			t.Errorf("PadLeft(%q, %d, %q) = %q, 期望 %q", tc.input, tc.length, tc.padChar, result, tc.expected)
		}
	}
}

func TestPadCenter(t *testing.T) {
	testCases := []struct {
		input    string
		length   int
		padChar  rune
		expected string
	}{
		{"hello", 10, '-', "--hello---"},
		{"hello", 9, '-', "--hello--"},
		{"hello", 5, '-', "hello"},
		{"hello", 3, '-', "hello"},
		{"", 5, '-', "-----"},
	}

	for _, tc := range testCases {
		result := PadCenter(tc.input, tc.length, tc.padChar)
		if result != tc.expected {
			t.Errorf("PadCenter(%q, %d, %q) = %q, 期望 %q", tc.input, tc.length, tc.padChar, result, tc.expected)
		}
	}
}

func TestRemoveDuplicates(t *testing.T) {
	testCases := []struct {
		input    []string
		expected []string
	}{
		{[]string{"a", "b", "c", "b", "a"}, []string{"a", "b", "c"}},
		{[]string{"a", "a", "a"}, []string{"a"}},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{}, []string{}},
		{[]string{""}, []string{""}},
	}

	for _, tc := range testCases {
		result := RemoveDuplicates(tc.input)
		if len(result) != len(tc.expected) {
			t.Errorf("RemoveDuplicates 返回的长度不正确: 期望 %d, 实际 %d", len(tc.expected), len(result))
			continue
		}
		for i := range result {
			if result[i] != tc.expected[i] {
				t.Errorf("RemoveDuplicates 返回的结果不匹配: 期望 %v, 实际 %v", tc.expected, result)
				break
			}
		}
	}
}

func TestSplitLines(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{"line1\nline2\nline3", []string{"line1", "line2", "line3"}},
		{"line1\r\nline2\r\nline3", []string{"line1", "line2", "line3"}},
		{"single", []string{"single"}},
		{"", []string{""}},
		{"line1\n", []string{"line1", ""}},
	}

	for _, tc := range testCases {
		result := SplitLines(tc.input)
		if len(result) != len(tc.expected) {
			t.Errorf("SplitLines 返回的长度不正确: 期望 %d, 实际 %d", len(tc.expected), len(result))
			continue
		}
		for i := range result {
			if result[i] != tc.expected[i] {
				t.Errorf("SplitLines 返回的结果不匹配: 期望 %v, 实际 %v", tc.expected, result)
				break
			}
		}
	}
}

func TestJoinLines(t *testing.T) {
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{"line1", "line2", "line3"}, "line1\nline2\nline3"},
		{[]string{"single"}, "single"},
		{[]string{}, ""},
		{[]string{"line1", "", "line3"}, "line1\n\nline3"},
	}

	for _, tc := range testCases {
		result := JoinLines(tc.input)
		if result != tc.expected {
			t.Errorf("JoinLines(%v) = %q, 期望 %q", tc.input, result, tc.expected)
		}
	}
}

func TestCountLines(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{"line1\nline2\nline3", 3},
		{"line1\r\nline2\r\nline3", 3},
		{"single", 1},
		{"", 1},
		{"line1\n", 2},
	}

	for _, tc := range testCases {
		result := CountLines(tc.input)
		if result != tc.expected {
			t.Errorf("CountLines(%q) = %d, 期望 %d", tc.input, result, tc.expected)
		}
	}
}

// TestPathUtils 路径工具函数测试
func TestNormalizePath(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"path/to/file", filepath.Join("path", "to", "file")},
		{"path/../file", "file"},
		{"./path/to/file", filepath.Join("path", "to", "file")},
		{"/absolute/path", filepath.Join("/", "absolute", "path")},
	}

	for _, tc := range testCases {
		result := NormalizePath(tc.input)
		if result != tc.expected {
			t.Errorf("NormalizePath(%q) = %q, 期望 %q", tc.input, result, tc.expected)
		}
	}
}

func TestGetRelativePath(t *testing.T) {
	// 创建临时目录结构
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}

	testFile := filepath.Join(subDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 测试获取相对路径
	relPath, err := GetRelativePath(tempDir, testFile)
	if err != nil {
		t.Errorf("GetRelativePath 返回错误: %v", err)
	}
	expected := filepath.Join("subdir", "test.txt")
	if filepath.ToSlash(relPath) != filepath.ToSlash(expected) {
		t.Errorf("GetRelativePath 返回的路径不正确: 期望 %q, 实际 %q", expected, relPath)
	}

	// 测试无效路径
	_, err = GetRelativePath("/nonexistent", testFile)
	if err == nil {
		t.Error("GetRelativePath 应该对无效路径返回错误")
	}
}

func TestGetAbsolutePath(t *testing.T) {
	// 测试相对路径
	relPath := "utils.go"
	absPath, err := GetAbsolutePath(relPath)
	if err != nil {
		t.Errorf("GetAbsolutePath 返回错误: %v", err)
	}
	if !filepath.IsAbs(absPath) {
		t.Error("GetAbsolutePath 应该返回绝对路径")
	}

	// 测试已经存在的绝对路径
	if _, err := GetAbsolutePath(absPath); err != nil {
		t.Errorf("GetAbsolutePath 对绝对路径返回错误: %v", err)
	}
}

func TestIsSubPath(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}

	testCases := []struct {
		parent   string
		child    string
		expected bool
	}{
		{tempDir, subDir, true},
		{tempDir, tempDir, false}, // 相同路径不算子路径
		{tempDir, "/other", false},
		{tempDir, filepath.Join(tempDir, "..", "other"), false},
	}

	for _, tc := range testCases {
		result := IsSubPath(tc.parent, tc.child)
		if result != tc.expected {
			t.Errorf("IsSubPath(%q, %q) = %v, 期望 %v", tc.parent, tc.child, result, tc.expected)
		}
	}
}

func TestGetCommonPath(t *testing.T) {
	// 创建临时目录结构用于测试
	tempDir := t.TempDir()
	dir1 := filepath.Join(tempDir, "a", "b", "c")
	dir2 := filepath.Join(tempDir, "a", "b", "d")
	dir3 := filepath.Join(tempDir, "a", "b", "e")
	dir4 := filepath.Join(tempDir, "a", "d", "e")
	
	// 创建目录
	for _, dir := range []string{dir1, dir2, dir3, dir4} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("创建测试目录失败: %v", err)
		}
	}

	testCases := []struct {
		paths    []string
		expected string
	}{
		{[]string{dir1, dir2, dir3}, filepath.Join(tempDir, "a", "b")},
		{[]string{dir1, dir4}, filepath.Join(tempDir, "a")},
		{[]string{dir1}, filepath.Join(tempDir, "a", "b")},
		{[]string{}, ""},
	}

	for _, tc := range testCases {
		result := GetCommonPath(tc.paths)
		if result != tc.expected {
			t.Errorf("GetCommonPath(%v) = %q, 期望 %q", tc.paths, result, tc.expected)
		}
	}
}

// TestRegexUtils 正则表达式工具函数测试
func TestMatchPattern(t *testing.T) {
	testCases := []struct {
		pattern  string
		text     string
		expected bool
		hasError bool
	}{
		{"hello", "hello world", true, false},
		{"world", "hello world", true, false},
		{"^hello", "hello world", true, false},
		{"world$", "hello world", true, false},
		{"xyz", "hello world", false, false},
		{"[", "hello", false, true}, // 无效的正则表达式
	}

	for _, tc := range testCases {
		result, err := MatchPattern(tc.pattern, tc.text)
		if tc.hasError {
			if err == nil {
				t.Errorf("MatchPattern(%q, %q) 应该返回错误", tc.pattern, tc.text)
			}
		} else {
			if err != nil {
				t.Errorf("MatchPattern(%q, %q) 返回意外错误: %v", tc.pattern, tc.text, err)
			}
			if result != tc.expected {
				t.Errorf("MatchPattern(%q, %q) = %v, 期望 %v", tc.pattern, tc.text, result, tc.expected)
			}
		}
	}
}

func TestFindMatches(t *testing.T) {
	testCases := []struct {
		pattern  string
		text     string
		expected []string
		hasError bool
	}{
		{"l", "hello world", []string{"l", "l", "l"}, false},
		{"o", "hello world", []string{"o", "o"}, false},
		{"xyz", "hello world", []string{}, false},
		{"[", "hello", nil, true}, // 无效的正则表达式
	}

	for _, tc := range testCases {
		result, err := FindMatches(tc.pattern, tc.text)
		if tc.hasError {
			if err == nil {
				t.Errorf("FindMatches(%q, %q) 应该返回错误", tc.pattern, tc.text)
			}
		} else {
			if err != nil {
				t.Errorf("FindMatches(%q, %q) 返回意外错误: %v", tc.pattern, tc.text, err)
			}
			if len(result) != len(tc.expected) {
				t.Errorf("FindMatches(%q, %q) 返回的匹配数量不正确: 期望 %d, 实际 %d", tc.pattern, tc.text, len(tc.expected), len(result))
				continue
			}
			for i := range result {
				if result[i] != tc.expected[i] {
					t.Errorf("FindMatches(%q, %q) 返回的结果不匹配: 期望 %v, 实际 %v", tc.pattern, tc.text, tc.expected, result)
					break
				}
			}
		}
	}
}

func TestReplacePattern(t *testing.T) {
	testCases := []struct {
		pattern     string
		replacement string
		text        string
		expected    string
		hasError    bool
	}{
		{"world", "Go", "hello world", "hello Go", false},
		{"l", "L", "hello", "heLLo", false},
		{"xyz", "ABC", "hello world", "hello world", false},
		{"[", "X", "hello", "", true}, // 无效的正则表达式
	}

	for _, tc := range testCases {
		result, err := ReplacePattern(tc.pattern, tc.replacement, tc.text)
		if tc.hasError {
			if err == nil {
				t.Errorf("ReplacePattern(%q, %q, %q) 应该返回错误", tc.pattern, tc.replacement, tc.text)
			}
		} else {
			if err != nil {
				t.Errorf("ReplacePattern(%q, %q, %q) 返回意外错误: %v", tc.pattern, tc.replacement, tc.text, err)
			}
			if result != tc.expected {
				t.Errorf("ReplacePattern(%q, %q, %q) = %q, 期望 %q", tc.pattern, tc.replacement, tc.text, result, tc.expected)
			}
		}
	}
}

// TestTimeUtils 时间工具函数测试
func TestFormatDuration(t *testing.T) {
	testCases := []struct {
		duration time.Duration
		expected string
	}{
		{500 * time.Millisecond, "0.5s"},
		{1500 * time.Millisecond, "1.5s"},
		{30 * time.Second, "30.0s"},
		{90 * time.Second, "1.5m"},
		{2 * time.Minute, "2.0m"},
		{90 * time.Minute, "1.5h"},
		{3 * time.Hour, "3.0h"},
	}

	for _, tc := range testCases {
		result := FormatDuration(tc.duration)
		if result != tc.expected {
			t.Errorf("FormatDuration(%v) = %q, 期望 %q", tc.duration, result, tc.expected)
		}
	}
}

func TestParseTime(t *testing.T) {
	testCases := []struct {
		timeStr  string
		hasError bool
	}{
		{"2023-01-01T12:00:00Z", false},     // RFC3339
		{"2023-01-01 12:00:00", false},     // 2006-01-02 15:04:05
		{"2023-01-01", false},              // 2006-01-02
		{"12:00:00", false},                // 15:04:05
		{"2023/01/01", false},              // 2006/01/02
		{"invalid", true},                  // 无效格式
		{"2023-13-01", true},               // 无效日期
	}

	for _, tc := range testCases {
		result, err := ParseTime(tc.timeStr)
		if tc.hasError {
			if err == nil {
				t.Errorf("ParseTime(%q) 应该返回错误", tc.timeStr)
			}
		} else {
			if err != nil {
				t.Errorf("ParseTime(%q) 返回意外错误: %v", tc.timeStr, err)
			}
			if result.IsZero() {
				t.Errorf("ParseTime(%q) 返回零时间", tc.timeStr)
			}
		}
	}
}

func TestFormatFileSize(t *testing.T) {
	testCases := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1099511627776, "1.0 TB"},
	}

	for _, tc := range testCases {
		result := FormatFileSize(tc.bytes)
		if result != tc.expected {
			t.Errorf("FormatFileSize(%d) = %q, 期望 %q", tc.bytes, result, tc.expected)
		}
	}
}

// TestValidationUtils 验证工具函数测试
func TestIsValidFilename(t *testing.T) {
	testCases := []struct {
		filename string
		expected bool
	}{
		{"valid.txt", true},
		{"file-name_123.go", true},
		{"", false},
		{"file/name.txt", false},
		{"file\\name.txt", false},
		{"file:name.txt", false},
		{"file*name.txt", false},
		{"file?name.txt", false},
		{"file\"name.txt", false},
		{"file<name.txt", false},
		{"file>name.txt", false},
		{"file|name.txt", false},
		{".hidden", false},
		{"file.", false},
		{" file.txt", false},
		{"file.txt ", false},
	}

	for _, tc := range testCases {
		result := IsValidFilename(tc.filename)
		if result != tc.expected {
			t.Errorf("IsValidFilename(%q) = %v, 期望 %v", tc.filename, result, tc.expected)
		}
	}
}

func TestIsValidPath(t *testing.T) {
	testCases := []struct {
		path     string
		expected bool
	}{
		{"/valid/path", true},
		{"relative/path", true},
		{"", false},
		{strings.Repeat("a", 300), false}, // 超过Windows路径长度限制
		{"path\x00withnull", false},     // 包含空字符
	}

	for _, tc := range testCases {
		result := IsValidPath(tc.path)
		if result != tc.expected {
			t.Errorf("IsValidPath(%q) = %v, 期望 %v", tc.path, result, tc.expected)
		}
	}
}

func TestSafePathJoin(t *testing.T) {
	testCases := []struct {
		base     string
		elem     string
		expected string
		hasError bool
	}{
		{"/base", "file.txt", filepath.Join("/base", "file.txt"), false},
		{"/base", "subdir/file.txt", filepath.Join("/base", "subdir", "file.txt"), false},
		{"/base", "../file.txt", "", true}, // 路径遍历攻击
		{"/base", "subdir/../file.txt", "", true}, // 路径遍历攻击
		{"/base", "", filepath.Join("/base", ""), false},
	}

	for _, tc := range testCases {
		result, err := SafePathJoin(tc.base, tc.elem)
		if tc.hasError {
			if err == nil {
				t.Errorf("SafePathJoin(%q, %q) 应该返回错误", tc.base, tc.elem)
			}
		} else {
			if err != nil {
				t.Errorf("SafePathJoin(%q, %q) 返回意外错误: %v", tc.base, tc.elem, err)
			}
			if result != tc.expected {
				t.Errorf("SafePathJoin(%q, %q) = %q, 期望 %q", tc.base, tc.elem, result, tc.expected)
			}
		}
	}
}

// TestColorUtils 颜色工具函数测试
func TestColorize(t *testing.T) {
	text := "test"
	colored := Colorize(text, ColorRed)
	
	if !strings.Contains(colored, string(ColorRed)) {
		t.Error("Colorize 应该包含颜色代码")
	}
	if !strings.Contains(colored, string(ColorReset)) {
		t.Error("Colorize 应该包含重置代码")
	}
	if !strings.Contains(colored, text) {
		t.Error("Colorize 应该包含原始文本")
	}
}

func TestErrorColor(t *testing.T) {
	result := ErrorColor("error")
	if !strings.Contains(result, string(ColorRed)) {
		t.Error("ErrorColor 应该使用红色")
	}
}

func TestSuccessColor(t *testing.T) {
	result := SuccessColor("success")
	if !strings.Contains(result, string(ColorGreen)) {
		t.Error("SuccessColor 应该使用绿色")
	}
}

func TestWarningColor(t *testing.T) {
	result := WarningColor("warning")
	if !strings.Contains(result, string(ColorYellow)) {
		t.Error("WarningColor 应该使用黄色")
	}
}

func TestInfoColor(t *testing.T) {
	result := InfoColor("info")
	if !strings.Contains(result, string(ColorBlue)) {
		t.Error("InfoColor 应该使用蓝色")
	}
}