// Package utils 提供通用工具函数
package utils

import (
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
	
	"golang.org/x/text/encoding/unicode"
)

// DetectEncoding 检测文件编码
func DetectEncoding(data []byte) (string, []byte) {
	if len(data) == 0 {
		return "utf-8", data
	}

	// 检查BOM头
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return "utf-8", data[3:] // 跳过BOM
	}
	
	if len(data) >= 2 {
		// UTF-16 LE BOM
		if data[0] == 0xFF && data[1] == 0xFE {
			return "utf-16le", data[2:]
		}
		// UTF-16 BE BOM
		if data[0] == 0xFE && data[1] == 0xFF {
			return "utf-16be", data[2:]
		}
	}

	// 检查是否为UTF-8
	if isValidUTF8(data) {
		return "utf-8", data
	}

	// 检查是否为UTF-16
	if isValidUTF16(data) {
		return "utf-16le", data // 默认小端
	}

	// 检查是否为GBK
	if isValidGBK(data) {
		return "gbk", data
	}

	// 检查是否为ANSI (Windows-1252)
	if isValidANSI(data) {
		return "ansi", data
	}

	// 默认按UTF-8处理
	return "utf-8", data
}

// isValidUTF8 检查数据是否为有效的UTF-8编码
func isValidUTF8(data []byte) bool {
	for i := 0; i < len(data); {
		r := rune(data[i])
		if r < 0x80 {
			i++
			continue
		}

		// 多字节UTF-8序列
		if i+1 >= len(data) {
			return false
		}

		if r < 0xE0 {
			// 2字节序列: 110xxxxx 10xxxxxx
			if data[i+1]&0xC0 != 0x80 {
				return false
			}
			i += 2
		} else if r < 0xF0 {
			// 3字节序列: 1110xxxx 10xxxxxx 10xxxxxx
			if i+2 >= len(data) || data[i+1]&0xC0 != 0x80 || data[i+2]&0xC0 != 0x80 {
				return false
			}
			i += 3
		} else if r < 0xF8 {
			// 4字节序列: 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
			if i+3 >= len(data) || data[i+1]&0xC0 != 0x80 || data[i+2]&0xC0 != 0x80 || data[i+3]&0xC0 != 0x80 {
				return false
			}
			i += 4
		} else {
			return false
		}
	}
	return true
}

// isValidUTF16 检查数据是否为有效的UTF-16编码
func isValidUTF16(data []byte) bool {
	if len(data)%2 != 0 {
		return false
	}

	// 简单检查：查看是否有大量0字节
	zeroCount := 0
	for _, b := range data {
		if b == 0 {
			zeroCount++
		}
	}

	// 如果超过25%的字符是0，可能是UTF-16
	return float64(zeroCount)/float64(len(data)) > 0.25
}

// isValidGBK 检查数据是否为有效的GBK编码
func isValidGBK(data []byte) bool {
	for i := 0; i < len(data); {
		b := data[i]
		if b < 0x80 {
			i++
			continue
		}

		// GBK双字节字符
		if i+1 < len(data) {
			b2 := data[i+1]
			// GBK范围: 0x8140-0xFEFE
			if (b >= 0x81 && b <= 0xFE) && (b2 >= 0x40 && b2 <= 0xFE && b2 != 0x7F) {
				i += 2
				continue
			}
		}
		return false
	}
	return true
}

// isValidANSI 检查数据是否为有效的ANSI (Windows-1252) 编码
func isValidANSI(data []byte) bool {
	// ANSI (Windows-1252) 是单字节编码，所有字节值都有效
	// 这里我们检查是否主要是可打印字符
	printableCount := 0
	for _, b := range data {
		// 可打印字符和控制字符
		if (b >= 32 && b <= 126) || (b >= 160 && b <= 255) || b == 9 || b == 10 || b == 13 {
			printableCount++
		}
	}
	
	// 如果大部分字符都是可打印的，认为是ANSI编码
	return float64(printableCount)/float64(len(data)) > 0.8
}

// ConvertToUTF8 将数据转换为UTF-8编码
func ConvertToUTF8(data []byte, encoding string) (string, error) {
	switch strings.ToLower(encoding) {
	case "utf-8":
		return string(data), nil
	case "utf-16le":
		return utf16ToUTF8(data, true), nil
	case "utf-16be":
		return utf16ToUTF8(data, false), nil
	case "gbk":
		return gbkToUTF8(data)
	case "ansi", "windows-1252", "cp1252":
		return ansiToUTF8(data)
	default:
		return string(data), nil // 默认按UTF-8处理
	}
}

// utf16ToUTF8 将UTF-16转换为UTF-8
func utf16ToUTF8(data []byte, littleEndian bool) string {
	if len(data)%2 != 0 {
		return string(data) // 如果不是偶数长度，直接返回
	}

	var result strings.Builder
	for i := 0; i < len(data); i += 2 {
		var r rune
		if littleEndian {
			r = rune(data[i]) | rune(data[i+1])<<8
		} else {
			r = rune(data[i])<<8 | rune(data[i+1])
		}

		if r == 0 {
			break // 遇到null字符停止
		}

		if r < 0x80 {
			result.WriteByte(byte(r))
		} else if r < 0x800 {
			result.WriteByte(0xC0 | byte(r>>6))
			result.WriteByte(0x80 | byte(r&0x3F))
		} else {
			result.WriteByte(0xE0 | byte(r>>12))
			result.WriteByte(0x80 | byte((r>>6)&0x3F))
			result.WriteByte(0x80 | byte(r&0x3F))
		}
	}

	return result.String()
}

// gbkToUTF8 将GBK转换为UTF-8
func gbkToUTF8(data []byte) (string, error) {
	decoder := simplifiedchinese.GBK.NewDecoder()
	reader := transform.NewReader(strings.NewReader(string(data)), decoder)
	result, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("GBK解码失败: %w", err)
	}
	return string(result), nil
}

// ansiToUTF8 将ANSI (Windows-1252) 转换为UTF-8
func ansiToUTF8(data []byte) (string, error) {
	decoder := charmap.Windows1252.NewDecoder()
	reader := transform.NewReader(strings.NewReader(string(data)), decoder)
	result, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("ANSI解码失败: %w", err)
	}
	return string(result), nil
}

// GetEncodingDecoder 获取指定编码的解码器
func GetEncodingDecoder(enc string) (encoding.Encoding, error) {
	switch strings.ToLower(enc) {
	case "gbk", "gb2312":
		return simplifiedchinese.GBK, nil
	case "ansi", "windows-1252", "cp1252":
		return charmap.Windows1252, nil
	case "utf-16le":
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM), nil
	case "utf-16be":
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM), nil
	case "utf-8":
		return encoding.Nop, nil
	default:
		return nil, fmt.Errorf("不支持的编码: %s", enc)
	}
}