// Package config 提供配置管理功能的单元测试
package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	
	"code-context-generator/pkg/constants"
	"code-context-generator/pkg/types"
)

// TestNewManager 测试创建新的配置管理器
func TestNewManager(t *testing.T) {
	manager := NewManager()
	if manager == nil {
		t.Fatal("NewManager() 返回 nil")
	}

	cm, ok := manager.(*ConfigManager)
	if !ok {
		t.Fatal("NewManager() 返回的类型不是 *ConfigManager")
	}

	if cm.config == nil {
		t.Fatal("ConfigManager.config 为 nil")
	}
}

// TestConfigManager_Load 测试加载配置文件
func TestConfigManager_Load(t *testing.T) {
	// 创建临时配置文件
	tempDir := t.TempDir()
	yamlConfig := filepath.Join(tempDir, "test.yaml")

	// 创建测试配置数据
	testConfig := GetDefaultConfig()
	testConfig.Output.DefaultFormat = "json"
	testConfig.Output.OutputDir = "./test_output"

	// 保存为YAML配置文件
	err := os.WriteFile(yamlConfig, []byte(`formats:
  xml:
    enabled: true
  json:
    enabled: true
  toml:
    enabled: true
  markdown:
    enabled: true
fields:
  custom_names: {}
  filter:
    include: []
    exclude: []
  processing:
    max_length: 0
    add_line_numbers: false
    trim_whitespace: true
    code_highlight: false
filters:
  max_file_size: "10MB"
  exclude_patterns: []
  include_patterns: []
  max_depth: 0
  follow_symlinks: false
output:
  default_format: "json"
  output_dir: "./test_output"
  filename_template: "context_{{.timestamp}}.{{.extension}}"
  timestamp_format: "20060102_150405"
ui:
  selector:
    show_hidden: false
    show_size: true
    show_modified: false
  autocomplete:
    enabled: true
    min_chars: 1
    max_suggestions: 10
`), 0644)
	if err != nil {
		t.Fatalf("创建测试配置文件失败: %v", err)
	}

	tests := []struct {
		name       string
		configPath string
		wantErr    bool
	}{
		{"加载YAML配置", yamlConfig, false},
		{"加载不存在的文件", filepath.Join(tempDir, "nonexistent.yaml"), false}, // 应该创建默认配置
		{"空路径", "", false}, // 应该使用默认配置
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager()
			err := manager.Load(tt.configPath)
			
			if tt.wantErr && err == nil {
				t.Errorf("Load() 期望错误但没有得到错误")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Load() 出现意外错误: %v", err)
			}
		})
	}
}

// TestConfigManager_Get 测试获取配置
func TestConfigManager_Get(t *testing.T) {
	manager := NewManager()
	config := manager.Get()
	
	if config == nil {
		t.Fatal("Get() 返回 nil")
	}

	// 验证默认配置值
	if config.Output.DefaultFormat != constants.DefaultFormat {
		t.Errorf("默认格式不匹配: 期望 %s, 得到 %s", constants.DefaultFormat, config.Output.DefaultFormat)
	}
	
	if config.Output.FilenameTemplate != constants.DefaultFilenameTemplate {
		t.Errorf("文件名模板不匹配: 期望 %s, 得到 %s", constants.DefaultFilenameTemplate, config.Output.FilenameTemplate)
	}
}

// TestConfigManager_Validate 测试配置验证
func TestConfigManager_Validate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() Manager
		wantErr bool
	}{
		{
			name: "有效配置",
			setup: func() Manager {
				manager := NewManager()
				return manager
			},
			wantErr: false,
		},
		{
			name: "空配置",
			setup: func() Manager {
				return &ConfigManager{config: nil}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := tt.setup()
			err := manager.Validate()
			
			if tt.wantErr && err == nil {
				t.Errorf("Validate() 期望错误但没有得到错误")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Validate() 出现意外错误: %v", err)
			}
		})
	}
}

// TestConfigManager_Save 测试保存配置
func TestConfigManager_Save(t *testing.T) {
	tempDir := t.TempDir()
	
	tests := []struct {
		name      string
		format    string
		wantErr   bool
	}{
		{"保存为YAML", "yaml", false},
		{"保存为JSON", "json", false},
		{"保存为TOML", "toml", false},
		{"保存为不支持的格式", "unsupported", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager()
			configPath := filepath.Join(tempDir, "test."+tt.format)
			
			err := manager.Save(configPath, tt.format)
			
			if tt.wantErr && err == nil {
				t.Errorf("Save() 期望错误但没有得到错误")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Save() 出现意外错误: %v", err)
			}
			
			if !tt.wantErr {
				// 验证文件已创建
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					t.Errorf("Save() 未创建文件: %s", configPath)
				}
			}
		})
	}
}

// TestConfigManager_GetEnvOverrides 测试获取环境变量覆盖
func TestConfigManager_GetEnvOverrides(t *testing.T) {
	// 设置测试环境变量
	os.Setenv(constants.EnvPrefix+"DEFAULT_FORMAT", "json")
	os.Setenv(constants.EnvPrefix+"OUTPUT_DIR", "/test/output")
	os.Setenv(constants.EnvPrefix+"MAX_DEPTH", "5")
	
	// 清理环境变量
	defer func() {
		os.Unsetenv(constants.EnvPrefix + "DEFAULT_FORMAT")
		os.Unsetenv(constants.EnvPrefix + "OUTPUT_DIR")
		os.Unsetenv(constants.EnvPrefix + "MAX_DEPTH")
	}()
	
	manager := NewManager()
	overrides := manager.GetEnvOverrides()
	
	if overrides == nil {
		t.Fatal("GetEnvOverrides() 返回 nil")
	}
	
	// 验证环境变量覆盖
	if overrides["default_format"] != "json" {
		t.Errorf("期望 default_format = json, 得到 %s", overrides["default_format"])
	}
	if overrides["output_dir"] != "/test/output" {
		t.Errorf("期望 output_dir = /test/output, 得到 %s", overrides["output_dir"])
	}
	if overrides["max_depth"] != "5" {
		t.Errorf("期望 max_depth = 5, 得到 %s", overrides["max_depth"])
	}
}

// TestConfigManager_GetOutputFilename 测试生成输出文件名
func TestConfigManager_GetOutputFilename(t *testing.T) {
	manager := NewManager()
	
	filename := manager.GetOutputFilename("txt")
	
	if filename == "" {
		t.Error("GetOutputFilename() 返回空文件名")
	}
	
	// 验证文件名包含扩展名
	if !strings.Contains(filename, "txt") {
		t.Errorf("文件名 %s 不包含扩展名 txt", filename)
	}
	
	// 验证文件名包含时间戳占位符
	if !strings.Contains(filename, "{{.timestamp}}") && !strings.Contains(filename, "20") {
		t.Errorf("文件名 %s 不包含时间戳信息", filename)
	}
}

// TestConfigManager_Reload 测试重新加载配置
func TestConfigManager_Reload(t *testing.T) {
	manager := NewManager()
	
	// 测试未设置路径时的重载
	err := manager.Reload()
	if err == nil {
		t.Error("期望Reload()在未设置路径时返回错误")
	}
	
	// 创建临时配置文件
	tempFile := filepath.Join(t.TempDir(), "config_test.yaml")
	err = manager.Save(tempFile, "yaml")
	if err != nil {
		t.Fatalf("保存配置文件失败: %v", err)
	}
	
	// 加载配置
	err = manager.Load(tempFile)
	if err != nil {
		t.Fatalf("加载配置文件失败: %v", err)
	}
	
	// 修改配置
	config := manager.Get()
	originalFormat := config.Output.DefaultFormat
	config.Output.DefaultFormat = "json"
	
	// 重新加载配置
	err = manager.Reload()
	if err != nil {
		t.Errorf("Reload() 失败: %v", err)
	}
	
	// 验证配置已恢复
	config = manager.Get()
	if config.Output.DefaultFormat != originalFormat {
		t.Errorf("配置未正确重载: 期望 %s, 得到 %s", originalFormat, config.Output.DefaultFormat)
	}
}

// TestConfigManager_GenerateOutput 测试生成输出内容
func TestConfigManager_GenerateOutput(t *testing.T) {
	manager := NewManager()
	
	// 创建简单的测试数据（避免XML序列化问题）
	testData := types.ContextData{
		Files: []types.FileInfo{
			{
				Name:    "test.go",
				Path:    "test.go",
				Content: "package main\n\nfunc main() {}",
				Size:    30,
			},
		},
		Folders: []types.FolderInfo{
			{
				Name:  "src",
				Path:  "src",
				Files: []types.FileInfo{
					{
						Name:    "main.go",
						Path:    "src/main.go",
						Content: "package main",
						Size:    20,
					},
				},
			},
		},
		FileCount:   1,
		FolderCount: 1,
	}
	
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"生成JSON", "json", false},
		{"生成TOML", "toml", false},
		{"生成Markdown", "markdown", false},
		{"不支持的格式", "unsupported", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := manager.GenerateOutput(testData, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateOutput() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && output == "" {
				t.Error("GenerateOutput() 返回空字符串")
			}
		})
	}
}

// TestGetDefaultConfig 测试获取默认配置
func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()
	
	if config == nil {
		t.Fatal("GetDefaultConfig() 返回 nil")
	}

	// 验证默认配置的关键字段
	if config.Output.DefaultFormat == "" {
		t.Error("默认格式不能为空")
	}
	
	if config.Output.FilenameTemplate == "" {
		t.Error("文件名模板不能为空")
	}
}

// TestLoadConfig 测试加载配置文件
func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		extension string
		wantErr   bool
	}{
		{
			name: "加载YAML配置",
			content: `output:
  default_format: json
  filename_template: "{{.timestamp}}.{{.extension}}"
`,
			extension: ".yaml",
			wantErr: false,
		},
		{
			name: "加载JSON配置",
			content: `{
  "output": {
	"default_format": "json",
	"filename_template": "{{.timestamp}}.{{.extension}}"
  }
}`,
			extension: ".json",
			wantErr: false,
		},
		{
			name: "加载TOML配置",
			content: `[output]
default_format = "json"
filename_template = "{{.timestamp}}.{{.extension}}"
`,
			extension: ".toml",
			wantErr: false,
		},
		{
			name:      "不支持的格式",
			content:   `test content`,
			extension: ".txt",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile := filepath.Join(t.TempDir(), "config"+tt.extension)
			err := os.WriteFile(tempFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("写入测试文件失败: %v", err)
			}
			
			config, err := LoadConfig(tempFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			if !tt.wantErr && config == nil {
				t.Error("LoadConfig() 返回 nil 配置")
			}
		})
	}
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}