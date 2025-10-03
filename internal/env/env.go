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

// 环境变量常量定义
const (
	// 格式配置
	EnvDefaultFormat = constants.EnvPrefix + "DEFAULT_FORMAT"
	
	// 输出配置
	EnvOutputDir        = constants.EnvPrefix + "OUTPUT_DIR"
	EnvFilenameTemplate = constants.EnvPrefix + "FILENAME_TEMPLATE"
	EnvTimestampFormat  = constants.EnvPrefix + "TIMESTAMP_FORMAT"
	
	// 文件处理配置
	EnvMaxFileSize     = constants.EnvPrefix + "MAX_FILE_SIZE"
	EnvMaxDepth        = constants.EnvPrefix + "MAX_DEPTH"
	// EnvRecursive       = constants.EnvPrefix + "RECURSIVE" // 已移除recursive参数
	EnvIncludeHidden   = constants.EnvPrefix + "INCLUDE_HIDDEN"
	EnvFollowSymlinks  = constants.EnvPrefix + "FOLLOW_SYMLINKS"
	EnvExcludeBinary   = constants.EnvPrefix + "EXCLUDE_BINARY"
	EnvExcludePatterns = constants.EnvPrefix + "EXCLUDE_PATTERNS"
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
			unit = strings.TrimSpace(sizeStr[i:])
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
	envVars[EnvDefaultFormat] = GetEnvWithDefault(EnvDefaultFormat, "xml")
	
	// 输出配置
	envVars[EnvOutputDir] = GetEnvWithDefault(EnvOutputDir, "")
	envVars[EnvFilenameTemplate] = GetEnvWithDefault(EnvFilenameTemplate, "")
	envVars[EnvTimestampFormat] = GetEnvWithDefault(EnvTimestampFormat, "")
	
	// 文件处理配置
	envVars[EnvMaxFileSize] = GetEnvWithDefault(EnvMaxFileSize, "")
	envVars[EnvMaxDepth] = GetEnvWithDefault(EnvMaxDepth, "")
	// envVars[EnvRecursive] = strconv.FormatBool(GetEnvBool(EnvRecursive, false)) // 已移除recursive参数
	envVars[EnvIncludeHidden] = strconv.FormatBool(GetEnvBool(EnvIncludeHidden, false))
	envVars[EnvFollowSymlinks] = strconv.FormatBool(GetEnvBool(EnvFollowSymlinks, false))
	envVars[EnvExcludeBinary] = strconv.FormatBool(GetEnvBool(EnvExcludeBinary, true))
	envVars[EnvExcludePatterns] = GetEnvWithDefault(EnvExcludePatterns, "")
	
	return envVars
}

// 获取默认格式配置
func GetDefaultFormat() string {
	return GetEnvWithDefault(EnvDefaultFormat, "xml")
}

// 获取输出目录配置
func GetOutputDir() string {
	return GetEnvWithDefault(EnvOutputDir, "")
}

// 获取文件名模板配置
func GetFilenameTemplate() string {
	return GetEnvWithDefault(EnvFilenameTemplate, "")
}

// 获取时间戳格式配置
func GetTimestampFormat() string {
	return GetEnvWithDefault(EnvTimestampFormat, "")
}

// 获取最大文件大小配置
func GetMaxFileSize() string {
	return GetEnvWithDefault(EnvMaxFileSize, "10MB")
}

// 获取最大深度配置
func GetMaxDepth() int {
	return GetEnvInt(EnvMaxDepth, 0)
}

// 获取是否递归配置 - 已移除，使用max-depth代替
// func GetRecursive() bool {
// 	return GetEnvBool(EnvRecursive, false)
// }

// 获取是否包含隐藏文件配置
func GetIncludeHidden() bool {
	return GetEnvBool(EnvIncludeHidden, false)
}

// 获取是否跟随符号链接配置
func GetFollowSymlinks() bool {
	return GetEnvBool(EnvFollowSymlinks, false)
}

// 获取是否排除二进制文件配置
func GetExcludeBinary() bool {
	return GetEnvBool(EnvExcludeBinary, true)
}

// 获取排除模式配置
func GetExcludePatterns() string {
	return GetEnvWithDefault(EnvExcludePatterns, "")
}
func ApplyEnvOverrides(config map[string]interface{}) {
	envVars := GetAllEnvVars()
	
	for key, value := range envVars {
		if value != "" {
			config[key] = value
		}
	}
}