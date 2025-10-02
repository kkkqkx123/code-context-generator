package utils

import (
	"runtime"
	"testing"
)

func TestNormalizeLineEndings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Unix 换行符",
			input:    "第一行\n第二行\n第三行",
			expected: "第一行\n第二行\n第三行",
		},
		{
			name:     "Windows 换行符",
			input:    "第一行\r\n第二行\r\n第三行",
			expected: "第一行\n第二行\n第三行",
		},
		{
			name:     "混合换行符",
			input:    "第一行\n第二行\r\n第三行",
			expected: "第一行\n第二行\n第三行",
		},
		{
			name:     "无换行符",
			input:    "单行文本",
			expected: "单行文本",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
	}

	if runtime.GOOS == "windows" {
		// Windows 系统下，预期结果是 \r\n
		for i := range tests {
			tests[i].expected = normalizeToWindows(tests[i].expected)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeLineEndings(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeLineEndings() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestNormalizeLineEndingsBytes(t *testing.T) {
	input := []byte("第一行\n第二行\r\n第三行")
	result := NormalizeLineEndingsBytes(input)
	
	expected := "第一行\n第二行\n第三行"
	if runtime.GOOS == "windows" {
		expected = "第一行\r\n第二行\r\n第三行"
	}
	
	if string(result) != expected {
		t.Errorf("NormalizeLineEndingsBytes() = %q, want %q", string(result), expected)
	}
}

// 辅助函数：将字符串中的 \n 转换为 \r\n
func normalizeToWindows(s string) string {
	return string(NormalizeLineEndingsBytes([]byte(s)))
}