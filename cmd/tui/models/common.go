// Package models TUI模型定义
package models

import (
	"code-context-generator/pkg/types"
)

// 全局配置变量，需要在main中设置
var cfg *types.Config

// SetConfig 设置全局配置
func SetConfig(config *types.Config) {
	cfg = config
}

// GetConfig 获取全局配置
func GetConfig() *types.Config {
	return cfg
}
