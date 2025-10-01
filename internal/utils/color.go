// Package utils 提供通用工具函数
package utils

// ColorCode 颜色代码
type ColorCode string

const (
	ColorReset  ColorCode = "\033[0m"
	ColorRed    ColorCode = "\033[31m"
	ColorGreen  ColorCode = "\033[32m"
	ColorYellow ColorCode = "\033[33m"
	ColorBlue   ColorCode = "\033[34m"
	ColorPurple ColorCode = "\033[35m"
	ColorCyan   ColorCode = "\033[36m"
	ColorWhite  ColorCode = "\033[37m"
)

// Colorize 给文本添加颜色
func Colorize(text string, color ColorCode) string {
	return string(color) + text + string(ColorReset)
}

// ErrorColor 错误颜色
func ErrorColor(text string) string {
	return Colorize(text, ColorRed)
}

// SuccessColor 成功颜色
func SuccessColor(text string) string {
	return Colorize(text, ColorGreen)
}

// WarningColor 警告颜色
func WarningColor(text string) string {
	return Colorize(text, ColorYellow)
}

// InfoColor 信息颜色
func InfoColor(text string) string {
	return Colorize(text, ColorBlue)
}