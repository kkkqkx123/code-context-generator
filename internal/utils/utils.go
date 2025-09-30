// Package utils 提供通用工具函数
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// FileUtils 文件工具函数

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// DirectoryExists 检查目录是否存在
func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// GetFileHash 获取文件哈希值
func GetFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetFileSize 获取文件大小
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetFileModTime 获取文件修改时间
func GetFileModTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// IsTextFile 检查是否为文本文件
func IsTextFile(path string) bool {
	// 简单的文本文件检测
	ext := strings.ToLower(filepath.Ext(path))
	textExtensions := []string{
		".txt", ".md", ".json", ".xml", ".yaml", ".yml", ".toml",
		".go", ".py", ".js", ".ts", ".java", ".cpp", ".c", ".h",
		".html", ".css", ".scss", ".sass", ".sql", ".sh", ".bat",
		".ps1", ".rb", ".php", ".rs", ".swift", ".kt", ".scala",
	}

	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}
	return false
}

// IsBinaryFile 检查是否为二进制文件
func IsBinaryFile(path string) bool {
	return !IsTextFile(path)
}

// StringUtils 字符串工具函数

// TruncateString 截断字符串
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	if maxLength <= 3 {
		return s[:maxLength]
	}
	return s[:maxLength-3] + "..."
}

// PadString 填充字符串
func PadString(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(string(padChar), length-len(s))
	return s + padding
}

// PadLeft 左填充
func PadLeft(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(string(padChar), length-len(s))
	return padding + s
}

// PadCenter 居中填充
func PadCenter(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	totalPadding := length - len(s)
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding
	return strings.Repeat(string(padChar), leftPadding) + s + strings.Repeat(string(padChar), rightPadding)
}

// RemoveDuplicates 移除字符串切片中的重复项
func RemoveDuplicates(strings []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(strings))
	
	for _, s := range strings {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	
	return result
}

// SplitLines 分割字符串为多行
func SplitLines(s string) []string {
	return strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
}

// JoinLines 连接多行为字符串
func JoinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

// CountLines 计算行数
func CountLines(s string) int {
	return len(SplitLines(s))
}

// PathUtils 路径工具函数

// NormalizePath 规范化路径
func NormalizePath(path string) string {
	return filepath.Clean(path)
}

// GetRelativePath 获取相对路径
func GetRelativePath(base, target string) (string, error) {
	return filepath.Rel(base, target)
}

// GetAbsolutePath 获取绝对路径
func GetAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

// IsSubPath 检查是否为子路径
func IsSubPath(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..")
}

// GetCommonPath 获取共同路径
func GetCommonPath(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	
	if len(paths) == 1 {
		return filepath.Dir(paths[0])
	}

	// 转换为绝对路径
	absPaths := make([]string, 0, len(paths))
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			continue // 跳过无效路径
		}
		absPaths = append(absPaths, absPath)
	}
	
	if len(absPaths) == 0 {
		return ""
	}

	// 找到最短的路径
	minPath := absPaths[0]
	for _, path := range absPaths {
		if len(path) < len(minPath) {
			minPath = path
		}
	}

	// 从最短路径开始，逐步向上查找共同路径
	for {
		common := true
		for _, path := range absPaths {
			if !strings.HasPrefix(path, minPath) {
				common = false
				break
			}
		}
		
		if common {
			return minPath
		}
		
		parent := filepath.Dir(minPath)
		if parent == minPath {
			break
		}
		minPath = parent
	}

	return ""
}

// RegexUtils 正则表达式工具函数

// MatchPattern 匹配模式
func MatchPattern(pattern, text string) (bool, error) {
	matched, err := regexp.MatchString(pattern, text)
	if err != nil {
		return false, fmt.Errorf("正则表达式匹配失败: %w", err)
	}
	return matched, nil
}

// FindMatches 查找所有匹配
func FindMatches(pattern, text string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("编译正则表达式失败: %w", err)
	}
	return re.FindAllString(text, -1), nil
}

// ReplacePattern 替换模式
func ReplacePattern(pattern, replacement, text string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("编译正则表达式失败: %w", err)
	}
	return re.ReplaceAllString(text, replacement), nil
}

// TimeUtils 时间工具函数

// FormatDuration 格式化持续时间
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

// ParseTime 解析时间字符串
func ParseTime(timeStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
		"15:04:05",
		"2006/01/02",
	}

	for _, format := range formats {
		t, err := time.Parse(format, timeStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析时间字符串: %s", timeStr)
}

// FormatFileSize 格式化文件大小
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// ValidationUtils 验证工具函数

// IsValidFilename 检查文件名是否有效
func IsValidFilename(filename string) bool {
	if filename == "" {
		return false
	}
	
	// 检查是否包含非法字符
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(filename, char) {
			return false
		}
	}
	
	// 检查是否以点或空格开头/结尾
	if strings.HasPrefix(filename, ".") || strings.HasSuffix(filename, ".") ||
	   strings.HasPrefix(filename, " ") || strings.HasSuffix(filename, " ") {
		return false
	}
	
	return true
}

// IsValidPath 检查路径是否有效
func IsValidPath(path string) bool {
	if path == "" {
		return false
	}
	
	// 检查路径长度
	if len(path) > 260 { // Windows路径长度限制
		return false
	}
	
	// 检查是否包含空字符
	if strings.Contains(path, "\x00") {
		return false
	}
	
	return true
}

// SafePathJoin 安全地连接路径
func SafePathJoin(base, elem string) (string, error) {
	// 检查路径遍历攻击
	if strings.Contains(elem, "..") {
		return "", fmt.Errorf("路径包含非法字符: %s", elem)
	}
	
	joined := filepath.Join(base, elem)
	
	// 确保结果仍在基础路径内
	if !strings.HasPrefix(filepath.Clean(joined), filepath.Clean(base)) {
		return "", fmt.Errorf("路径超出基础目录范围")
	}
	
	return joined, nil
}

// ColorUtils 颜色工具函数

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