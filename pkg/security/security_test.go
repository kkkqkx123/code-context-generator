// Package security 安全扫描测试
package security

import (
	"os"
	"path/filepath"
	"testing"
	
	"code-context-generator/pkg/types"
)

// TestCredentialsDetector 测试凭证检测器
func TestCredentialsDetector(t *testing.T) {
	detector := NewCredentialsDetector()
	
	// 测试用例
	testCases := []struct {
		name     string
		content  string
		expected int // 期望发现的问题数量
	}{
		{
			name: "检测硬编码密码",
			content: `
				password = "secret123"
				pwd = "mypassword"
				passwd = "123456"
			`,
			expected: 3,
		},
		{
			name: "检测API密钥",
			content: `
				api_key = "sk-1234567890"
				apiKey = "secret-key-123"
				api-key = "test-key"
			`,
			expected: 3,
		},
		{
			name: "检测AWS凭证",
			content: `
				aws_access_key = "AKIAIOSFODNN7EXAMPLE"
				aws_secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
			`,
			expected: 2,
		},
		{
			name: "检测数据库密码",
			content: `
				database_password = "db_secret"
				db_password = "db_pass123"
			`,
			expected: 2,
		},
		{
			name: "无敏感信息",
			content: `
				username = "user123"
				port = 8080
				host = "localhost"
			`,
			expected: 0,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			issues := detector.Detect("test.go", tc.content)
			if len(issues) != tc.expected {
				t.Errorf("期望发现 %d 个问题，实际发现 %d 个问题", tc.expected, len(issues))
			}
		})
	}
}

// TestSQLInjectionDetector 测试SQL注入检测器
func TestSQLInjectionDetector(t *testing.T) {
	detector := NewSQLInjectionDetector()
	
	testCases := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name: "检测SQL注入",
			content: `
				query := "SELECT * FROM users WHERE id = " + userId
				query = fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)
			`,
			expected: 2,
		},
		{
			name: "检测数据库查询注入",
			content: `
				db.query("SELECT * FROM users WHERE id = " + request.id)
				db.exec("DELETE FROM users WHERE name = " + userName)
			`,
			expected: 2,
		},
		{
			name: "安全查询",
			content: `
				query := "SELECT * FROM users WHERE id = ?"
				stmt, err := db.Prepare(query)
			`,
			expected: 0,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			issues := detector.Detect("test.go", tc.content)
			if len(issues) != tc.expected {
				t.Errorf("期望发现 %d 个问题，实际发现 %d 个问题", tc.expected, len(issues))
			}
		})
	}
}

// TestXSSDetector 测试XSS检测器
func TestXSSDetector(t *testing.T) {
	detector := NewXSSDetector()
	
	testCases := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name: "检测XSS漏洞",
			content: `
				element.innerHTML = userInput
				document.write(htmlContent)
			`,
			expected: 2,
		},
		{
			name: "检测eval使用",
			content: `
				eval(userCode)
				setTimeout("alert('" + message + "')", 1000)
			`,
			expected: 2,
		},
		{
			name: "安全代码",
			content: `
				element.textContent = userInput
				element.innerText = sanitizedContent
			`,
			expected: 0,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			issues := detector.Detect("test.js", tc.content)
			if len(issues) != tc.expected {
				t.Errorf("期望发现 %d 个问题，实际发现 %d 个问题", tc.expected, len(issues))
			}
		})
	}
}

// TestPathTraversalDetector 测试路径遍历检测器
func TestPathTraversalDetector(t *testing.T) {
	detector := NewPathTraversalDetector()
	
	testCases := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name: "检测路径遍历",
			content: `
				filePath := "../config/secret.txt"
				path := "..\\..\\etc\\passwd"
			`,
			expected: 2,
		},
		{
			name: "检测文件操作注入",
			content: `
				os.Open("data/" + userInput)
				readFile("logs/" + fileName)
			`,
			expected: 2,
		},
		{
			name: "安全文件操作",
			content: `
				os.Open(filepath.Join("data", sanitizedPath))
				readFile(filepath.Clean(userInput))
			`,
			expected: 0,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			issues := detector.Detect("test.go", tc.content)
			if len(issues) != tc.expected {
				t.Errorf("期望发现 %d 个问题，实际发现 %d 个问题", tc.expected, len(issues))
			}
		})
	}
}

// TestQualityDetector 测试代码质量检测器
func TestQualityDetector(t *testing.T) {
	detector := NewQualityDetector()
	
	testCases := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name: "检测未使用变量",
			content: `
				var unusedVar = "test"
				let unusedLet = 123
				const unusedConst = true
			`,
			expected: 3,
		},
		{
			name: "检测错误处理",
			content: `
				if err != nil {
					// 错误被忽略
				}
			`,
			expected: 1,
		},
		{
			name: "良好代码",
			content: `
				var usedVar = "test"
				fmt.Println(usedVar)
				if err != nil {
					return err
				}
			`,
			expected: 0,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			issues := detector.Detect("test.go", tc.content)
			if len(issues) != tc.expected {
				t.Errorf("期望发现 %d 个问题，实际发现 %d 个问题", tc.expected, len(issues))
			}
		})
	}
}

// TestSecurityScanner 测试安全扫描器
func TestSecurityScanner(t *testing.T) {
	// 创建临时测试文件
	tempDir, err := os.MkdirTemp("", "security-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.go")
	testContent := `
package main

import "fmt"

func main() {
	password = "secret123"
	query := "SELECT * FROM users WHERE id = " + userId
	var unusedVar = "test"
	if err != nil {
		// 错误被忽略
	}
}
`
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatal(err)
	}
	
	// 创建安全配置
	config := &types.SecurityConfig{
		Enabled:        true,
		FailOnCritical: false,
		ScanLevel:      "standard",
		ReportFormat:   "text",
		Detectors: types.DetectorConfig{
			HardcodedCredentials: types.HardcodedCredentialsConfig{
				Enabled:          true,
				SeverityThreshold: types.SeverityLow,
				Patterns:         []string{},
			},
			SQLInjection: types.VulnerabilityConfig{
				Enabled:          true,
				SeverityThreshold: types.SeverityLow,
			},
			XSSVulnerabilities: types.VulnerabilityConfig{
				Enabled:          true,
				SeverityThreshold: types.SeverityLow,
			},
			PathTraversal: types.VulnerabilityConfig{
				Enabled:          true,
				SeverityThreshold: types.SeverityLow,
			},
			UnusedVariables: types.QualityConfig{
				Enabled:          true,
				SeverityThreshold: types.SeverityLow,
			},
			ErrorHandling: types.QualityConfig{
				Enabled:          true,
				SeverityThreshold: types.SeverityLow,
			},
		},
		Exclusions: types.ExclusionConfig{
			Files:    []string{},
			Patterns: []string{},
			Rules:    []types.ExclusionRule{},
		},
		Reporting: types.ReportingConfig{
			Format:         "text",
			OutputFile:     "",
			IncludeDetails: true,
			ShowStatistics: true,
		},
	}
	
	// 创建扫描器
	scanner := NewSecurityScanner(config)
	
	// 执行扫描
	report, err := scanner.Scan(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	
	// 验证报告
	if report.Summary.TotalFiles == 0 {
		t.Error("扫描文件数应该大于0")
	}
	
	if report.Summary.IssuesFound == 0 {
		t.Error("应该发现安全问题")
	}
	
	if report.Summary.CriticalIssues == 0 {
		t.Error("应该发现严重问题")
	}
	
	// 验证报告内容
	if len(report.Issues) == 0 {
		t.Error("报告应该包含问题详情")
	}
}

// TestDetectorRegistry 测试检测器注册表
func TestDetectorRegistry(t *testing.T) {
	registry := NewDetectorRegistry()
	
	// 测试获取所有检测器
	detectors := registry.GetAllDetectors()
	if len(detectors) == 0 {
		t.Error("注册表应该包含检测器")
	}
	
	// 测试获取特定语言的检测器
	goDetectors := registry.GetDetectorsForLanguage("go")
	if len(goDetectors) == 0 {
		t.Error("应该找到支持Go语言的检测器")
	}
	
	// 测试获取特定检测器
	detector := registry.GetDetector("hardcoded_credentials")
	if detector == nil {
		t.Error("应该找到硬编码凭证检测器")
	}
}

// TestSeverityLevelString 测试严重性级别字符串表示
func TestSeverityLevelString(t *testing.T) {
	testCases := []struct {
		severity types.SeverityLevel
		expected string
	}{
		{types.SeverityLow, "low"},
		{types.SeverityMedium, "medium"},
		{types.SeverityHigh, "high"},
		{types.SeverityCritical, "critical"},
	}
	
	for _, tc := range testCases {
		result := tc.severity.String()
		if result != tc.expected {
			t.Errorf("期望 %s，实际 %s", tc.expected, result)
		}
	}
}