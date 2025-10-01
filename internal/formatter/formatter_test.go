package formatter

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"code-context-generator/pkg/types"
)

// 测试辅助函数
func createTestFileInfo() types.FileInfo {
	return types.FileInfo{
		Path:    "test/file.go",
		Name:    "file.go",
		Size:    1024,
		ModTime: time.Now(),
		Content: "package main\n\nfunc main() {\n\tprintln(\"Hello World\")\n}",
	}
}

func createTestFolderInfo() types.FolderInfo {
	return types.FolderInfo{
		Path:     "test/folder",
		Name:     "folder",
		ModTime:  time.Now(),
		Files:    []types.FileInfo{createTestFileInfo()},
	}
}

func createTestContextData() types.ContextData {
	return types.ContextData{
		Files:       []types.FileInfo{createTestFileInfo()},
		Folders:     []types.FolderInfo{createTestFolderInfo()},
		FileCount:   1,
		FolderCount: 1,
		TotalSize:   1024,
	}
}

// JSONFormatter 测试
func TestJSONFormatter_Format(t *testing.T) {
	formatter := NewJSONFormatter(nil)
	data := createTestContextData()

	result, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// 验证结果是有效的JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// 验证基本字段存在
	if _, exists := parsed["files"]; !exists {
		t.Error("Missing 'files' field in JSON output")
	}
	if _, exists := parsed["folders"]; !exists {
		t.Error("Missing 'folders' field in JSON output")
	}
}

func TestJSONFormatter_FormatFile(t *testing.T) {
	formatter := NewJSONFormatter(nil)
	file := createTestFileInfo()

	result, err := formatter.FormatFile(file)
	if err != nil {
		t.Fatalf("FormatFile failed: %v", err)
	}

	// 验证结果是有效的JSON
	var parsed types.FileInfo
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// 验证字段
	if parsed.Name != file.Name {
		t.Errorf("Expected name %s, got %s", file.Name, parsed.Name)
	}
	if parsed.Size != file.Size {
		t.Errorf("Expected size %d, got %d", file.Size, parsed.Size)
	}
}

func TestJSONFormatter_FormatFolder(t *testing.T) {
	formatter := NewJSONFormatter(nil)
	folder := createTestFolderInfo()

	result, err := formatter.FormatFolder(folder)
	if err != nil {
		t.Fatalf("FormatFolder failed: %v", err)
	}

	// 验证结果是有效的JSON
	var parsed types.FolderInfo
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// 验证字段
	if parsed.Name != folder.Name {
		t.Errorf("Expected name %s, got %s", folder.Name, parsed.Name)
	}
	if len(parsed.Files) != len(folder.Files) {
		t.Errorf("Expected %d files, got %d", len(folder.Files), len(parsed.Files))
	}
}

// XMLFormatter 测试
func TestXMLFormatter_Format(t *testing.T) {
	formatter := NewXMLFormatter(nil)
	data := createTestContextData()

	result, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// 验证结果包含XML头
	if !strings.HasPrefix(result, xml.Header) {
		t.Error("XML result should start with XML header")
	}

	// 验证包含基本元素
	if !strings.Contains(result, "<context>") {
		t.Error("Missing <context> element in XML output")
	}
	if !strings.Contains(result, "<files>") {
		t.Error("Missing <files> element in XML output")
	}
	if !strings.Contains(result, "<folders>") {
		t.Error("Missing <folders> element in XML output")
	}
}

func TestXMLFormatter_FormatFile(t *testing.T) {
	formatter := NewXMLFormatter(nil)
	file := createTestFileInfo()

	result, err := formatter.FormatFile(file)
	if err != nil {
		t.Fatalf("FormatFile failed: %v", err)
	}

	// 验证结果包含XML头
	if !strings.HasPrefix(result, xml.Header) {
		t.Error("XML result should start with XML header")
	}

	// 验证包含文件元素
	if !strings.Contains(result, "<Path>") {
		t.Error("Missing <Path> element in XML output")
	}
	if !strings.Contains(result, "<Name>") {
		t.Error("Missing <Name> element in XML output")
	}
}

func TestXMLFormatter_FormatFolder(t *testing.T) {
	formatter := NewXMLFormatter(nil)
	folder := createTestFolderInfo()

	result, err := formatter.FormatFolder(folder)
	if err != nil {
		t.Fatalf("FormatFolder failed: %v", err)
	}

	// 验证结果包含XML头
	if !strings.HasPrefix(result, xml.Header) {
		t.Error("XML result should start with XML header")
	}

	// 验证包含文件夹元素
	if !strings.Contains(result, "<Path>") {
		t.Error("Missing <Path> element in XML output")
	}
	if !strings.Contains(result, "<Name>") {
		t.Error("Missing <Name> element in XML output")
	}
}

// TOMLFormatter 测试
func TestTOMLFormatter_Format(t *testing.T) {
	formatter := NewTOMLFormatter(nil)
	data := createTestContextData()

	result, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// 验证包含基本结构
	if !strings.Contains(result, "[files]") {
		t.Error("Missing [files] section in TOML output")
	}
	if !strings.Contains(result, "[[files.file]]") {
		t.Error("Missing [[files.file]] section in TOML output")
	}
	if !strings.Contains(result, "[folders]") {
		t.Error("Missing [folders] section in TOML output")
	}
	if !strings.Contains(result, "[[folders.folder]]") {
		t.Error("Missing [[folders.folder]] section in TOML output")
	}
}

func TestTOMLFormatter_FormatFile(t *testing.T) {
	formatter := NewTOMLFormatter(nil)
	file := createTestFileInfo()

	result, err := formatter.FormatFile(file)
	if err != nil {
		t.Fatalf("FormatFile failed: %v", err)
	}

	// 验证包含文件字段
	if !strings.Contains(result, "path = \"test/file.go\"") {
		t.Error("Missing or incorrect path field in TOML output")
	}
	if !strings.Contains(result, "name = \"file.go\"") {
		t.Error("Missing or incorrect name field in TOML output")
	}
	if !strings.Contains(result, "size = 1024") {
		t.Error("Missing or incorrect size field in TOML output")
	}
}

func TestTOMLFormatter_FormatFolder(t *testing.T) {
	formatter := NewTOMLFormatter(nil)
	folder := createTestFolderInfo()

	result, err := formatter.FormatFolder(folder)
	if err != nil {
		t.Fatalf("FormatFolder failed: %v", err)
	}

	// 验证包含文件夹字段
	if !strings.Contains(result, "path = \"test/folder\"") {
		t.Error("Missing or incorrect path field in TOML output")
	}
	if !strings.Contains(result, "name = \"folder\"") {
		t.Error("Missing or incorrect name field in TOML output")
	}
	if !strings.Contains(result, "file_count = 1") {
		t.Error("Missing or incorrect file_count field in TOML output")
	}
}

// MarkdownFormatter 测试
func TestMarkdownFormatter_Format(t *testing.T) {
	formatter := NewMarkdownFormatter(nil)
	data := createTestContextData()

	result, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// 验证包含Markdown标题
	if !strings.Contains(result, "# 代码上下文") {
		t.Error("Missing main title in Markdown output")
	}
	if !strings.Contains(result, "## 文件") {
		t.Error("Missing files section title in Markdown output")
	}
	if !strings.Contains(result, "## 文件夹") {
		t.Error("Missing folders section title in Markdown output")
	}

	// 验证包含代码块
	if !strings.Contains(result, "```go") {
		t.Error("Missing code block in Markdown output")
	}
	if !strings.Contains(result, "package main") {
		t.Error("Missing file content in Markdown output")
	}
}

func TestMarkdownFormatter_FormatFile(t *testing.T) {
	formatter := NewMarkdownFormatter(nil)
	file := createTestFileInfo()

	result, err := formatter.FormatFile(file)
	if err != nil {
		t.Fatalf("FormatFile failed: %v", err)
	}

	// 验证包含文件标题
	if !strings.Contains(result, "## file.go") {
		t.Error("Missing file title in Markdown output")
	}

	// 验证包含文件信息
	if !strings.Contains(result, "**路径**") {
		t.Error("Missing path information in Markdown output")
	}
	if !strings.Contains(result, "**大小**") {
		t.Error("Missing size information in Markdown output")
	}

	// 验证包含代码块
	if !strings.Contains(result, "```go") {
		t.Error("Missing code block in Markdown output")
	}
}

func TestMarkdownFormatter_FormatFolder(t *testing.T) {
	formatter := NewMarkdownFormatter(nil)
	folder := createTestFolderInfo()

	result, err := formatter.FormatFolder(folder)
	if err != nil {
		t.Fatalf("FormatFolder failed: %v", err)
	}

	// 验证包含文件夹标题
	if !strings.Contains(result, "## folder") {
		t.Error("Missing folder title in Markdown output")
	}

	// 验证包含文件夹信息
	if !strings.Contains(result, "**路径**") {
		t.Error("Missing path information in Markdown output")
	}
	if !strings.Contains(result, "**文件数**") {
		t.Error("Missing file count information in Markdown output")
	}

	// 验证包含文件列表
	if !strings.Contains(result, "### 文件列表") {
		t.Error("Missing file list title in Markdown output")
	}
}

// FormatterFactory 测试
func TestFormatterFactory(t *testing.T) {
	factory := NewFormatterFactory()

	// 注册格式
	jsonFormatter := NewJSONFormatter(nil)
	factory.Register("json", jsonFormatter)
	factory.Register("JSON", jsonFormatter) // 测试大小写不敏感

	// 测试获取格式
	formatter, err := factory.Get("json")
	if err != nil {
		t.Fatalf("Get formatter failed: %v", err)
	}
	if formatter == nil {
		t.Error("Formatter should not be nil")
	}
	if formatter.GetName() != "JSON" {
		t.Errorf("Expected formatter name 'JSON', got '%s'", formatter.GetName())
	}

	// 测试不存在的格式
	_, err = factory.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent format")
	}

	// 测试获取支持的格式
	supportedFormats := factory.GetSupportedFormats()
	if len(supportedFormats) != 1 { // json 和 JSON 应该被视为同一个
		t.Errorf("Expected 1 supported format, got %d", len(supportedFormats))
	}
}

func TestNewFormatter(t *testing.T) {
	// 测试创建JSON格式
	formatter, err := NewFormatter("json")
	if err != nil {
		t.Fatalf("NewFormatter failed: %v", err)
	}
	if formatter.GetName() != "JSON" {
		t.Errorf("Expected formatter name 'JSON', got '%s'", formatter.GetName())
	}

	// 测试创建XML格式
	formatter, err = NewFormatter("xml")
	if err != nil {
		t.Fatalf("NewFormatter failed: %v", err)
	}
	if formatter.GetName() != "XML" {
		t.Errorf("Expected formatter name 'XML', got '%s'", formatter.GetName())
	}

	// 测试创建TOML格式
	formatter, err = NewFormatter("toml")
	if err != nil {
		t.Fatalf("NewFormatter failed: %v", err)
	}
	if formatter.GetName() != "TOML" {
		t.Errorf("Expected formatter name 'TOML', got '%s'", formatter.GetName())
	}

	// 测试创建Markdown格式
	formatter, err = NewFormatter("markdown")
	if err != nil {
		t.Fatalf("NewFormatter failed: %v", err)
	}
	if formatter.GetName() != "Markdown" {
		t.Errorf("Expected formatter name 'Markdown', got '%s'", formatter.GetName())
	}

	// 测试不存在的格式
	_, err = NewFormatter("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent format")
	}
}

// 测试自定义配置的情况
func TestJSONFormatter_WithCustomConfig(t *testing.T) {
	customConfig := &types.FormatConfig{
		Structure: map[string]interface{}{
			"custom_field": "custom_value",
			"files":        []interface{}{},
		},
	}
	
	formatter := NewJSONFormatter(customConfig)
	data := createTestContextData()

	result, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format with custom config failed: %v", err)
	}

	// 验证结果是有效的JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// 验证自定义字段存在
	if _, exists := parsed["custom_field"]; !exists {
		t.Error("Missing 'custom_field' in custom config output")
	}
}

func TestJSONFormatter_WithCustomFields(t *testing.T) {
	customConfig := &types.FormatConfig{
		Fields: map[string]string{
			"custom_file_field": "custom_value",
		},
	}
	
	formatter := NewJSONFormatter(customConfig)
	file := createTestFileInfo()

	result, err := formatter.FormatFile(file)
	if err != nil {
		t.Fatalf("FormatFile with custom config failed: %v", err)
	}

	// 验证结果是有效的JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// 验证自定义字段存在
	if _, exists := parsed["custom_file_field"]; !exists {
		t.Error("Missing 'custom_file_field' in custom config output")
	}
}

// 测试XMLFormatter的自定义配置
func TestXMLFormatter_WithCustomConfig(t *testing.T) {
	// XMLFormatter不支持自定义配置中的复杂结构
	// 所以我们只测试基本的自定义配置，不设置Structure
	customConfig := &types.FormatConfig{
		Fields: map[string]string{
			"version": "1.0",
		},
	}
	
	formatter := NewXMLFormatter(customConfig)
	data := createTestContextData()

	result, err := formatter.Format(data)
	if err != nil {
		t.Fatalf("Format with custom config failed: %v", err)
	}

	// 验证结果包含XML头
	if !strings.HasPrefix(result, xml.Header) {
		t.Error("XML result should start with XML header")
	}

	// 验证基本的XML结构存在
	if !strings.Contains(result, "<context>") {
		t.Error("XML result should contain context element")
	}
}

// 测试错误处理
func TestFormatters_ErrorHandling(t *testing.T) {
	// 测试XMLFormatter处理不可序列化的数据
	t.Run("XMLFormatter_InvalidCustomConfig", func(t *testing.T) {
		customConfig := &types.FormatConfig{
			Structure: map[string]interface{}{
				"invalid": make(chan int), // channel不能被XML序列化
			},
		}
		
		formatter := NewXMLFormatter(customConfig)
		data := createTestContextData()

		_, err := formatter.Format(data)
		if err == nil {
			t.Error("Expected error for invalid XML custom config")
		}
		if !strings.Contains(err.Error(), "XML格式化失败") {
			t.Errorf("Expected XML formatting error, got: %v", err)
		}
	})

	// 测试JSONFormatter处理循环引用（虽然不太可能，但测试错误处理）
	t.Run("JSONFormatter_InvalidData", func(t *testing.T) {
		// 这里我们模拟一个不能被JSON序列化的ContextData
		// 由于实际类型是ContextData，我们测试自定义配置的情况
		customConfig := &types.FormatConfig{
			Structure: map[string]interface{}{
				"invalid": func() {},
			},
		}
		
		formatterWithInvalidConfig := NewJSONFormatter(customConfig)
		data := createTestContextData()
		
		_, err := formatterWithInvalidConfig.Format(data)
		if err == nil {
			t.Error("Expected error for invalid JSON custom config")
		}
		if !strings.Contains(err.Error(), "JSON格式化失败") {
			t.Errorf("Expected JSON formatting error, got: %v", err)
		}
	})
}

// 测试空数据和边界情况
func TestFormatters_EmptyData(t *testing.T) {
	emptyData := types.ContextData{
		Files:       []types.FileInfo{},
		Folders:     []types.FolderInfo{},
		FileCount:   0,
		FolderCount: 0,
		TotalSize:   0,
	}

	// 测试JSONFormatter
	jsonFormatter := NewJSONFormatter(nil)
	result, err := jsonFormatter.Format(emptyData)
	if err != nil {
		t.Fatalf("JSON format empty data failed: %v", err)
	}
	if !strings.Contains(result, `"files": null`) && !strings.Contains(result, `"files": []`) {
		t.Error("JSON empty data should contain empty files array")
	}

	// 测试XMLFormatter
	xmlFormatter := NewXMLFormatter(nil)
	result, err = xmlFormatter.Format(emptyData)
	if err != nil {
		t.Fatalf("XML format empty data failed: %v", err)
	}
	if !strings.Contains(result, "<files>") {
		t.Error("XML empty data should contain files element")
	}

	// 测试TOMLFormatter
	tomlFormatter := NewTOMLFormatter(nil)
	result, err = tomlFormatter.Format(emptyData)
	if err != nil {
		t.Fatalf("TOML format empty data failed: %v", err)
	}
	// TOML空数据不应该包含文件部分
	if strings.Contains(result, "[files]") {
		t.Error("TOML empty data should not contain files section")
	}

	// 测试MarkdownFormatter
	markdownFormatter := NewMarkdownFormatter(nil)
	result, err = markdownFormatter.Format(emptyData)
	if err != nil {
		t.Fatalf("Markdown format empty data failed: %v", err)
	}
	// Markdown空数据不应该包含文件部分
	if strings.Contains(result, "## 文件") {
		t.Error("Markdown empty data should not contain files section")
	}
}

// 测试FormatterFactory的大小写不敏感
func TestFormatterFactory_CaseInsensitive(t *testing.T) {
	factory := CreateDefaultFactory(nil)

	// 测试各种大小写变体
	testCases := []string{"json", "JSON", "Json", "jSoN"}
	
	for _, format := range testCases {
		formatter, err := factory.Get(format)
		if err != nil {
			t.Errorf("Get formatter for %s failed: %v", format, err)
		}
		if formatter == nil {
			t.Errorf("Formatter for %s should not be nil", format)
		}
		if formatter.GetName() != "JSON" {
			t.Errorf("Expected formatter name 'JSON' for %s, got '%s'", format, formatter.GetName())
		}
	}
}

// 辅助函数测试
func TestEscapeTOMLString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple text", "simple text"},
		{"text with \"quotes\"", "text with \\\"quotes\\\""},
		{"text with \\ backslash", "text with \\\\ backslash"},
		{"text with\nnewline", "text with\\nnewline"},
		{"text with\ttab", "text with\\ttab"},
		{"text with\rcarriage return", "text with\\rcarriage return"},
	}

	for _, test := range tests {
		result := escapeTOMLString(test.input)
		if result != test.expected {
			t.Errorf("escapeTOMLString(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}