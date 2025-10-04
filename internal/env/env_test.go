// Package env 提供.env文件加载和环境变量管理功能的单元测试
package env

import (
	"os"
	"strings"
	"testing"
)

// TestLoadEnv 测试加载.env文件功能
func TestLoadEnv(t *testing.T) {
	// 保存原始环境变量
	originalEnv := make(map[string]string)
	for _, key := range []string{"TEST_KEY_1", "TEST_KEY_2", "CODE_CONTEXT_DEFAULT_FORMAT"} {
		originalEnv[key] = os.Getenv(key)
	}

	// 测试用例结束后恢复原始环境变量
	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
		// 删除测试文件
		os.Remove(".env")
		os.Remove("test.env")
	}()

	tests := []struct {
		name          string
		envPath       string
		envContent    string
		expectedError bool
	}{
		{
			name:          "加载不存在的.env文件",
			envPath:       "",
			envContent:    "",
			expectedError: false, // 不应该报错
		},
		{
			name:          "加载存在的.env文件",
			envPath:       ".env",
			envContent:    "TEST_KEY_1=value1\nTEST_KEY_2=value2\n",
			expectedError: false,
		},
		{
			name:          "加载指定路径的.env文件",
			envPath:       "test.env",
			envContent:    "CODE_CONTEXT_DEFAULT_FORMAT=json\n",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理环境
			os.Unsetenv("TEST_KEY_1")
			os.Unsetenv("TEST_KEY_2")
			os.Unsetenv("CODE_CONTEXT_DEFAULT_FORMAT")

			// 如果指定了内容，创建.env文件
			if tt.envContent != "" && tt.envPath != "" {
				err := os.WriteFile(tt.envPath, []byte(tt.envContent), 0644)
				if err != nil {
					t.Fatalf("创建测试文件失败: %v", err)
				}
			}

			err := LoadEnv(tt.envPath)
			if (err != nil) != tt.expectedError {
				t.Errorf("LoadEnv() error = %v, expectedError %v", err, tt.expectedError)
			}

			// 验证环境变量是否正确设置
			if tt.envContent != "" && err == nil {
				lines := strings.Split(tt.envContent, "\n")
				for _, line := range lines {
					if line != "" && !strings.HasPrefix(line, "#") {
						parts := strings.SplitN(line, "=", 2)
						if len(parts) == 2 {
							key, expectedValue := parts[0], parts[1]
							actualValue := os.Getenv(key)
							if actualValue != expectedValue {
								t.Errorf("环境变量 %s = %v, 期望 %v", key, actualValue, expectedValue)
							}
						}
					}
				}
			}

			// 清理测试文件
			if tt.envPath != "" && tt.envPath != ".env" {
				os.Remove(tt.envPath)
			}
		})
	}
}

// TestGetEnvWithDefault 测试获取环境变量（带默认值）
func TestGetEnvWithDefault(t *testing.T) {
	// 保存原始环境变量
	originalValue := os.Getenv("TEST_ENV_VAR")
	defer func() {
		if originalValue == "" {
			os.Unsetenv("TEST_ENV_VAR")
		} else {
			os.Setenv("TEST_ENV_VAR", originalValue)
		}
	}()

	tests := []struct {
		name         string
		key          string
		defaultValue string
		setValue     string
		expected     string
	}{
		{
			name:         "环境变量存在",
			key:          "TEST_ENV_VAR",
			defaultValue: "default",
			setValue:     "actual",
			expected:     "actual",
		},
		{
			name:         "环境变量不存在",
			key:          "TEST_ENV_VAR_NOT_EXIST",
			defaultValue: "default",
			setValue:     "",
			expected:     "default",
		},
		{
			name:         "环境变量为空",
			key:          "TEST_ENV_VAR",
			defaultValue: "default",
			setValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			if tt.setValue != "" {
				os.Setenv(tt.key, tt.setValue)
			} else {
				os.Unsetenv(tt.key)
			}

			result := GetEnvWithDefault(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvWithDefault() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

// TestGetEnvBool 测试获取布尔类型的环境变量
func TestGetEnvBool(t *testing.T) {
	// 保存原始环境变量
	originalValue := os.Getenv("TEST_BOOL_VAR")
	defer func() {
		if originalValue == "" {
			os.Unsetenv("TEST_BOOL_VAR")
		} else {
			os.Setenv("TEST_BOOL_VAR", originalValue)
		}
	}()

	tests := []struct {
		name         string
		key          string
		defaultValue bool
		setValue     string
		expected     bool
	}{
		{
			name:         "环境变量为true",
			key:          "TEST_BOOL_VAR",
			defaultValue: false,
			setValue:     "true",
			expected:     true,
		},
		{
			name:         "环境变量为false",
			key:          "TEST_BOOL_VAR",
			defaultValue: true,
			setValue:     "false",
			expected:     false,
		},
		{
			name:         "环境变量为1",
			key:          "TEST_BOOL_VAR",
			defaultValue: false,
			setValue:     "1",
			expected:     true,
		},
		{
			name:         "环境变量为0",
			key:          "TEST_BOOL_VAR",
			defaultValue: true,
			setValue:     "0",
			expected:     false,
		},
		{
			name:         "环境变量不存在",
			key:          "TEST_BOOL_VAR_NOT_EXIST",
			defaultValue: true,
			setValue:     "",
			expected:     true,
		},
		{
			name:         "环境变量为无效值",
			key:          "TEST_BOOL_VAR",
			defaultValue: true,
			setValue:     "invalid",
			expected:     true, // 返回默认值
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			if tt.setValue != "" {
				os.Setenv(tt.key, tt.setValue)
			} else {
				os.Unsetenv(tt.key)
			}

			result := GetEnvBool(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvBool() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

// TestGetEnvInt 测试获取整数类型的环境变量
func TestGetEnvInt(t *testing.T) {
	// 保存原始环境变量
	originalValue := os.Getenv("TEST_INT_VAR")
	defer func() {
		if originalValue == "" {
			os.Unsetenv("TEST_INT_VAR")
		} else {
			os.Setenv("TEST_INT_VAR", originalValue)
		}
	}()

	tests := []struct {
		name         string
		key          string
		defaultValue int
		setValue     string
		expected     int
	}{
		{
			name:         "环境变量为有效整数",
			key:          "TEST_INT_VAR",
			defaultValue: 10,
			setValue:     "42",
			expected:     42,
		},
		{
			name:         "环境变量为负数",
			key:          "TEST_INT_VAR",
			defaultValue: 10,
			setValue:     "-5",
			expected:     -5,
		},
		{
			name:         "环境变量不存在",
			key:          "TEST_INT_VAR_NOT_EXIST",
			defaultValue: 10,
			setValue:     "",
			expected:     10,
		},
		{
			name:         "环境变量为无效值",
			key:          "TEST_INT_VAR",
			defaultValue: 10,
			setValue:     "invalid",
			expected:     10, // 返回默认值
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			if tt.setValue != "" {
				os.Setenv(tt.key, tt.setValue)
			} else {
				os.Unsetenv(tt.key)
			}

			result := GetEnvInt(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvInt() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

// TestGetEnvInt64 测试获取int64类型的环境变量
func TestGetEnvInt64(t *testing.T) {
	// 保存原始环境变量
	originalValue := os.Getenv("TEST_INT64_VAR")
	defer func() {
		if originalValue == "" {
			os.Unsetenv("TEST_INT64_VAR")
		} else {
			os.Setenv("TEST_INT64_VAR", originalValue)
		}
	}()

	tests := []struct {
		name         string
		key          string
		defaultValue int64
		setValue     string
		expected     int64
	}{
		{
			name:         "环境变量为有效int64",
			key:          "TEST_INT64_VAR",
			defaultValue: 100,
			setValue:     "9223372036854775807", // MaxInt64
			expected:     9223372036854775807,
		},
		{
			name:         "环境变量为大负数",
			key:          "TEST_INT64_VAR",
			defaultValue: 100,
			setValue:     "-9223372036854775808", // MinInt64
			expected:     -9223372036854775808,
		},
		{
			name:         "环境变量不存在",
			key:          "TEST_INT64_VAR_NOT_EXIST",
			defaultValue: 100,
			setValue:     "",
			expected:     100,
		},
		{
			name:         "环境变量为无效值",
			key:          "TEST_INT64_VAR",
			defaultValue: 100,
			setValue:     "invalid",
			expected:     100, // 返回默认值
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			if tt.setValue != "" {
				os.Setenv(tt.key, tt.setValue)
			} else {
				os.Unsetenv(tt.key)
			}

			result := GetEnvInt64(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvInt64() = %v, 期望 %v", result, tt.expected)
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
	}{
		{
			name:     "空字符串",
			sizeStr:  "",
			expected: 0,
		},
		{
			name:     "纯数字（字节）",
			sizeStr:  "1024",
			expected: 1024,
		},
		{
			name:     "KB单位",
			sizeStr:  "10KB",
			expected: 10 * 1024,
		},
		{
			name:     "MB单位",
			sizeStr:  "5MB",
			expected: 5 * 1024 * 1024,
		},
		{
			name:     "GB单位",
			sizeStr:  "2GB",
			expected: 2 * 1024 * 1024 * 1024,
		},
		{
			name:     "小写单位",
			sizeStr:  "10mb",
			expected: 10 * 1024 * 1024,
		},
		{
			name:     "带空格",
			sizeStr:  "  10 MB  ",
			expected: 10 * 1024 * 1024,
		},
		{
			name:     "无效格式（无数字）",
			sizeStr:  "MB",
			expected: 0,
		},
		{
			name:     "无效单位",
			sizeStr:  "10TB",
			expected: 10, // 默认按字节处理
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseFileSize(tt.sizeStr)
			if result != tt.expected {
				t.Errorf("ParseFileSize(%q) = %v, 期望 %v", tt.sizeStr, result, tt.expected)
			}
		})
	}
}

// TestGetAllEnvVars 测试获取所有环境变量配置
func TestGetAllEnvVars(t *testing.T) {
	// 保存原始环境变量
	originalValues := make(map[string]string)
	envKeys := []string{
		EnvDefaultFormat,
		EnvOutputDir,
		EnvFilenameTemplate,
		EnvTimestampFormat,
		EnvMaxFileSize,
		EnvMaxDepth,
		EnvIncludeHidden,
		EnvFollowSymlinks,
		EnvExcludeBinary,
		EnvExcludePatterns,
	}

	for _, key := range envKeys {
		originalValues[key] = os.Getenv(key)
	}

	defer func() {
		// 恢复原始环境变量
		for key, value := range originalValues {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// 设置一些测试环境变量
	os.Setenv(EnvDefaultFormat, "json")
	os.Setenv(EnvOutputDir, "/tmp/output")
	os.Setenv(EnvMaxFileSize, "20MB")
	os.Setenv(EnvMaxDepth, "5")
	os.Setenv(EnvIncludeHidden, "true")
	os.Setenv(EnvFollowSymlinks, "true")
	os.Setenv(EnvExcludeBinary, "false")
	os.Setenv(EnvExcludePatterns, "*.tmp,*.log")

	result := GetAllEnvVars()

	// 验证结果
	if result[EnvDefaultFormat] != "json" {
		t.Errorf("GetAllEnvVars()[%s] = %v, 期望 %v", EnvDefaultFormat, result[EnvDefaultFormat], "json")
	}

	if result[EnvOutputDir] != "/tmp/output" {
		t.Errorf("GetAllEnvVars()[%s] = %v, 期望 %v", EnvOutputDir, result[EnvOutputDir], "/tmp/output")
	}

	if result[EnvMaxFileSize] != "20MB" {
		t.Errorf("GetAllEnvVars()[%s] = %v, 期望 %v", EnvMaxFileSize, result[EnvMaxFileSize], "20MB")
	}

	if result[EnvMaxDepth] != "5" {
		t.Errorf("GetAllEnvVars()[%s] = %v, 期望 %v", EnvMaxDepth, result[EnvMaxDepth], "5")
	}

	if result[EnvIncludeHidden] != "true" {
		t.Errorf("GetAllEnvVars()[%s] = %v, 期望 %v", EnvIncludeHidden, result[EnvIncludeHidden], "true")
	}

	if result[EnvFollowSymlinks] != "true" {
		t.Errorf("GetAllEnvVars()[%s] = %v, 期望 %v", EnvFollowSymlinks, result[EnvFollowSymlinks], "true")
	}

	if result[EnvExcludeBinary] != "false" {
		t.Errorf("GetAllEnvVars()[%s] = %v, 期望 %v", EnvExcludeBinary, result[EnvExcludeBinary], "false")
	}

	if result[EnvExcludePatterns] != "*.tmp,*.log" {
		t.Errorf("GetAllEnvVars()[%s] = %v, 期望 %v", EnvExcludePatterns, result[EnvExcludePatterns], "*.tmp,*.log")
	}
}

// TestApplyEnvOverrides 测试应用环境变量覆盖
func TestApplyEnvOverrides(t *testing.T) {
	config := make(map[string]interface{})

	// 保存原始环境变量
	originalValues := make(map[string]string)
	envKeys := []string{
		EnvDefaultFormat,
		EnvOutputDir,
		EnvMaxFileSize,
	}

	for _, key := range envKeys {
		originalValues[key] = os.Getenv(key)
	}

	defer func() {
		// 恢复原始环境变量
		for key, value := range originalValues {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// 设置测试环境变量
	os.Setenv(EnvDefaultFormat, "toml")
	os.Setenv(EnvOutputDir, "/test/output")
	os.Setenv(EnvMaxFileSize, "50MB")

	ApplyEnvOverrides(config)

	// 验证配置是否被正确应用
	if config[EnvDefaultFormat] != "toml" {
		t.Errorf("ApplyEnvOverrides() 设置 %s = %v, 期望 %v", EnvDefaultFormat, config[EnvDefaultFormat], "toml")
	}

	if config[EnvOutputDir] != "/test/output" {
		t.Errorf("ApplyEnvOverrides() 设置 %s = %v, 期望 %v", EnvOutputDir, config[EnvOutputDir], "/test/output")
	}

	if config[EnvMaxFileSize] != "50MB" {
		t.Errorf("ApplyEnvOverrides() 设置 %s = %v, 期望 %v", EnvMaxFileSize, config[EnvMaxFileSize], "50MB")
	}
}

// TestConfigGetterFunctions 测试配置获取函数
func TestConfigGetterFunctions(t *testing.T) {
	// 保存原始环境变量
	originalValues := make(map[string]string)
	envKeys := []string{
		EnvDefaultFormat,
		EnvOutputDir,
		EnvFilenameTemplate,
		EnvTimestampFormat,
		EnvMaxFileSize,
		EnvMaxDepth,
		EnvIncludeHidden,
		EnvFollowSymlinks,
		EnvExcludeBinary,
		EnvExcludePatterns,
	}

	for _, key := range envKeys {
		originalValues[key] = os.Getenv(key)
	}

	defer func() {
		// 恢复原始环境变量
		for key, value := range originalValues {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// 测试字符串配置获取函数
	t.Run("字符串配置获取", func(t *testing.T) {
		os.Setenv(EnvDefaultFormat, "markdown")
		os.Setenv(EnvOutputDir, "/custom/output")
		os.Setenv(EnvFilenameTemplate, "custom_{{.timestamp}}.{{.extension}}")
		os.Setenv(EnvTimestampFormat, "2006-01-02")
		os.Setenv(EnvMaxFileSize, "15MB")
		os.Setenv(EnvExcludePatterns, "*.cache,*.temp")

		tests := []struct {
			name     string
			function func() string
			expected string
		}{
			{"GetDefaultFormat", GetDefaultFormat, "markdown"},
			{"GetOutputDir", GetOutputDir, "/custom/output"},
			{"GetFilenameTemplate", GetFilenameTemplate, "custom_{{.timestamp}}.{{.extension}}"},
			{"GetTimestampFormat", GetTimestampFormat, "2006-01-02"},
			{"GetMaxFileSize", GetMaxFileSize, "15MB"},
			{"GetExcludePatterns", GetExcludePatterns, "*.cache,*.temp"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tt.function()
				if result != tt.expected {
					t.Errorf("%s() = %v, 期望 %v", tt.name, result, tt.expected)
				}
			})
		}
	})

	// 测试整数配置获取函数
	t.Run("整数配置获取", func(t *testing.T) {
		os.Setenv(EnvMaxDepth, "10")

		result := GetMaxDepth()
		expected := 10
		if result != expected {
			t.Errorf("GetMaxDepth() = %v, 期望 %v", result, expected)
		}
	})

	// 测试布尔配置获取函数
	t.Run("布尔配置获取", func(t *testing.T) {
		os.Setenv(EnvIncludeHidden, "false")
		os.Setenv(EnvFollowSymlinks, "true")
		os.Setenv(EnvExcludeBinary, "false")

		tests := []struct {
			name     string
			function func() bool
			expected bool
		}{
			{"GetIncludeHidden", GetIncludeHidden, false},
			{"GetFollowSymlinks", GetFollowSymlinks, true},
			{"GetExcludeBinary", GetExcludeBinary, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tt.function()
				if result != tt.expected {
					t.Errorf("%s() = %v, 期望 %v", tt.name, result, tt.expected)
				}
			})
		}
	})
}

// TestDefaultValues 测试默认值
func TestDefaultValues(t *testing.T) {
	// 清理所有相关的环境变量
	envKeys := []string{
		EnvDefaultFormat,
		EnvOutputDir,
		EnvFilenameTemplate,
		EnvTimestampFormat,
		EnvMaxFileSize,
		EnvMaxDepth,
		EnvIncludeHidden,
		EnvFollowSymlinks,
		EnvExcludeBinary,
		EnvExcludePatterns,
	}

	for _, key := range envKeys {
		os.Unsetenv(key)
	}

	// 测试默认值
	tests := []struct {
		name     string
		function interface{}
		expected interface{}
	}{
		{"GetDefaultFormat默认值", GetDefaultFormat(), "xml"},
		{"GetOutputDir默认值", GetOutputDir(), ""},
		{"GetFilenameTemplate默认值", GetFilenameTemplate(), ""},
		{"GetTimestampFormat默认值", GetTimestampFormat(), ""},
		{"GetMaxFileSize默认值", GetMaxFileSize(), "10MB"},
		{"GetMaxDepth默认值", GetMaxDepth(), 0},
		{"GetIncludeHidden默认值", GetIncludeHidden(), false},
		{"GetFollowSymlinks默认值", GetFollowSymlinks(), false},
		{"GetExcludeBinary默认值", GetExcludeBinary(), true},
		{"GetExcludePatterns默认值", GetExcludePatterns(), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			switch f := tt.function.(type) {
			case string:
				result = f
			case int:
				result = f
			case bool:
				result = f
			default:
				t.Fatalf("不支持的函数返回类型")
			}

			if result != tt.expected {
				t.Errorf("%s = %v, 期望 %v", tt.name, result, tt.expected)
			}
		})
	}
}

// TestEnvironmentVariableConstants 测试环境变量常量
func TestEnvironmentVariableConstants(t *testing.T) {
	expectedConstants := map[string]string{
		"EnvDefaultFormat":       "CODE_CONTEXT_DEFAULT_FORMAT",
		"EnvOutputDir":           "CODE_CONTEXT_OUTPUT_DIR",
		"EnvFilenameTemplate":    "CODE_CONTEXT_FILENAME_TEMPLATE",
		"EnvTimestampFormat":     "CODE_CONTEXT_TIMESTAMP_FORMAT",
		"EnvMaxFileSize":         "CODE_CONTEXT_MAX_FILE_SIZE",
		"EnvMaxDepth":            "CODE_CONTEXT_MAX_DEPTH",
		"EnvIncludeHidden":       "CODE_CONTEXT_INCLUDE_HIDDEN",
		"EnvFollowSymlinks":      "CODE_CONTEXT_FOLLOW_SYMLINKS",
		"EnvExcludeBinary":       "CODE_CONTEXT_EXCLUDE_BINARY",
		"EnvExcludePatterns":     "CODE_CONTEXT_EXCLUDE_PATTERNS",
	}

	actualConstants := map[string]string{
		"EnvDefaultFormat":    EnvDefaultFormat,
		"EnvOutputDir":        EnvOutputDir,
		"EnvFilenameTemplate": EnvFilenameTemplate,
		"EnvTimestampFormat":  EnvTimestampFormat,
		"EnvMaxFileSize":      EnvMaxFileSize,
		"EnvMaxDepth":         EnvMaxDepth,
		"EnvIncludeHidden":    EnvIncludeHidden,
		"EnvFollowSymlinks":   EnvFollowSymlinks,
		"EnvExcludeBinary":    EnvExcludeBinary,
		"EnvExcludePatterns":  EnvExcludePatterns,
	}

	for name, expected := range expectedConstants {
		if actualConstants[name] != expected {
			t.Errorf("常量 %s = %v, 期望 %v", name, actualConstants[name], expected)
		}
	}
}
