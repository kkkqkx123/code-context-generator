package utils

import (
	"runtime"
	"strings"
)

// NormalizeLineEndings 将文本中的换行符标准化为当前操作系统的换行符
// 支持 Windows (\r\n) 和 Linux (\n) 系统
func NormalizeLineEndings(text string) string {
	if runtime.GOOS == "windows" {
		// Windows 系统：将所有换行符转换为 \r\n
		// 首先将 \r\n 转换为 \n，避免重复转换
		text = strings.ReplaceAll(text, "\r\n", "\n")
		// 然后将所有 \n 转换为 \r\n
		text = strings.ReplaceAll(text, "\n", "\r\n")
	} else {
		// Linux/Unix 系统：将所有换行符转换为 \n
		// 将 \r\n 转换为 \n
		text = strings.ReplaceAll(text, "\r\n", "\n")
		// 移除单独的 \r
		text = strings.ReplaceAll(text, "\r", "")
	}
	return text
}

// NormalizeLineEndingsBytes 将字节数组中的换行符标准化为当前操作系统的换行符
func NormalizeLineEndingsBytes(data []byte) []byte {
	return []byte(NormalizeLineEndings(string(data)))
}