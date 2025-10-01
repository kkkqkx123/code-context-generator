// Package env 提供.env文件加载和环境变量管理功能
package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"code-context-generator/pkg/constants"
)

// LoadEnv 加载.env文件到环境变量中
func LoadEnv(envPath string) error {
	// 如果没有指定路径，使用默认的.env文件
	if envPath == "" {
		envPath = ".env"
	}

	// 检查文件是否存在
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// 如果文件不存在，不报错，直接返回
		return nil
	}

	// 加载.env文件
	if err := godotenv.Load(envPath); err != nil {
		return fmt.Errorf("加载.env文件失败: %w", err)
	}

	return nil
}

// GetEnvWithDefault 获取环境变量，如果不存在则返回默认值
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvBool 获取布尔类型的环境变量
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// GetEnvInt 获取整数类型的环境变量
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvInt64 获取int64类型的环境变量
func GetEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// ParseFileSize 解析文件大小字符串 (例如: "10MB", "1KB")
func ParseFileSize(sizeStr string) int64 {
	sizeStr = strings.TrimSpace(strings.ToUpper(sizeStr))
	
	if sizeStr == "" {
		return 0
	}

	// 提取数字部分和单位部分
	var numStr string
	var unit string
	
	for i, char := range sizeStr {
		if char >= '0' && char <= '9' {
			numStr += string(char)
		} else {
			unit = sizeStr[i:]
			break
		}
	}

	if numStr == "" {
		return 0
	}

	size, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0
	}

	// 根据单位转换
	switch unit {
	case "B", "":
		return size
	case "KB":
		return size * 1024
	case "MB":
		return size * 1024 * 1024
	case "GB":
		return size * 1024 * 1024 * 1024
	default:
		return size // 默认按字节处理
	}
}

// GetAllEnvVars 获取所有环境变量配置
func GetAllEnvVars() map[string]string {
	envVars := make(map[string]string)
	
	// 格式配置
	envVars["default_format"] = GetEnvWithDefault(constants.EnvPrefix+"DEFAULT_FORMAT", "xml")
	
	// 输出配置
	envVars["output_dir"] = GetEnvWithDefault(constants.EnvPrefix+"OUTPUT_DIR", "")
	
	// 文件处理配置
	envVars["max_depth"] = GetEnvWithDefault(constants.EnvPrefix+"MAX_DEPTH", "0")
	envVars["recursive"] = strconv.FormatBool(GetEnvBool(constants.EnvPrefix+"RECURSIVE", false))
	envVars["include_hidden"] = strconv.FormatBool(GetEnvBool(constants.EnvPrefix+"INCLUDE_HIDDEN", false))
	envVars["follow_symlinks"] = strconv.FormatBool(GetEnvBool(constants.EnvPrefix+"FOLLOW_SYMLINKS", false))
	envVars["exclude_binary"] = strconv.FormatBool(GetEnvBool(constants.EnvPrefix+"EXCLUDE_BINARY", true))
	
	// 文件大小配置
	maxFileSize := GetEnvWithDefault(constants.EnvPrefix+"MAX_FILE_SIZE", "10MB")
	envVars["max_file_size"] = strconv.FormatInt(ParseFileSize(maxFileSize), 10)
	
	// 自动补全配置
	envVars["autocomplete_enabled"] = strconv.FormatBool(GetEnvBool(constants.EnvPrefix+"AUTOCOMPLETE_ENABLED", true))
	
	return envVars
}

// ApplyEnvOverrides 将环境变量应用到配置中
func ApplyEnvOverrides(config map[string]interface{}) {
	envVars := GetAllEnvVars()
	
	for key, value := range envVars {
		if value != "" {
			config[key] = value
		}
	}
}