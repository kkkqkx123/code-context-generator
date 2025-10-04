// Package security 实现安全扫描功能
package security

import (
	"regexp"
	"strings"

	"code-context-generator/pkg/types"
)

// BaseDetector 基础检测器结构体
type BaseDetector struct {
	name               string
	supportedLanguages []string
}

// NewBaseDetector 创建基础检测器
func NewBaseDetector(name string, languages []string) *BaseDetector {
	return &BaseDetector{
		name:               name,
		supportedLanguages: languages,
	}
}

// GetName 返回检测器名称
func (d *BaseDetector) GetName() string {
	return d.name
}

// GetSupportedLanguages 返回支持的语言
func (d *BaseDetector) GetSupportedLanguages() []string {
	return d.supportedLanguages
}

// CredentialsDetector 硬编码凭证检测器
type CredentialsDetector struct {
	*BaseDetector
	patterns []*regexp.Regexp
}

// NewCredentialsDetector 创建凭证检测器
func NewCredentialsDetector() *CredentialsDetector {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(password|passwd|pwd)\s*=\s*['\"][^'\"]+['\"]`),
		regexp.MustCompile(`(?i)(api[_-]?key|apikey)\s*=\s*['\"][^'\"]+['\"]`),
		regexp.MustCompile(`(?i)(secret|token)\s*=\s*['\"][^'\"]+['\"]`),
		regexp.MustCompile(`(?i)(aws[_-]?access[_-]?key|aws[_-]?secret[_-]?key)\s*=\s*['\"][^'\"]+['\"]`),
		regexp.MustCompile(`(?i)(database|db)[_-]?(password|passwd)\s*=\s*['\"][^'\"]+['\"]`),
	}

	return &CredentialsDetector{
		BaseDetector: NewBaseDetector("hardcoded_credentials", []string{"go", "python", "javascript", "java", "php", "ruby"}),
		patterns:     patterns,
	}
}

// Detect 检测硬编码凭证
func (d *CredentialsDetector) Detect(filePath string, content string) []types.SecurityIssue {
	var issues []types.SecurityIssue

	lines := strings.Split(content, "\n")
	for lineNum, line := range lines {
		for _, pattern := range d.patterns {
			if pattern.MatchString(line) {
				// 检查是否已经有相同行的凭证检测问题
				found := false
				for _, existing := range issues {
					if existing.Line == lineNum+1 {
						found = true
						break
					}
				}
				if !found {
					issue := types.SecurityIssue{
						ID:             "CREDENTIALS_001",
						Type:           "HardcodedCredentials",
						Severity:       types.SeverityHigh,
						Message:        "检测到硬编码的敏感凭证信息",
						File:           filePath,
						Line:           lineNum + 1,
						Column:         strings.Index(line, pattern.FindString(line)) + 1,
						Snippet:        strings.TrimSpace(line),
						Recommendation: "使用环境变量或安全的配置管理系统存储敏感信息",
						Confidence:     0.85,
					}
					issues = append(issues, issue)
				}
			}
		}
	}

	return issues
}

// SQLInjectionDetector SQL注入检测器
type SQLInjectionDetector struct {
	*BaseDetector
}

// NewSQLInjectionDetector 创建SQL注入检测器
func NewSQLInjectionDetector() *SQLInjectionDetector {
	return &SQLInjectionDetector{
		BaseDetector: NewBaseDetector("sql_injection", []string{"go", "python", "javascript", "java", "php", "ruby"}),
	}
}

// Detect 检测SQL注入漏洞
func (d *SQLInjectionDetector) Detect(filePath string, content string) []types.SecurityIssue {
	var issues []types.SecurityIssue

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE).*\+.*`),
		regexp.MustCompile(`(?i)fmt\.Sprintf`),
		regexp.MustCompile(`(?i)query\(.*\+.*`),
		regexp.MustCompile(`(?i)exec\(.*\+.*`),
	}

	lines := strings.Split(content, "\n")
	for lineNum, line := range lines {
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				// 检查是否已经有相同行的SQL注入问题
				found := false
				for _, existing := range issues {
					if existing.Line == lineNum+1 {
						found = true
						break
					}
				}
				if !found {
					issue := types.SecurityIssue{
						ID:             "SQL_INJECTION_001",
						Type:           "SQLInjection",
						Severity:       types.SeverityHigh,
						Message:        "检测到可能的SQL注入漏洞",
						File:           filePath,
						Line:           lineNum + 1,
						Column:         strings.Index(line, pattern.FindString(line)) + 1,
						Snippet:        strings.TrimSpace(line),
						Recommendation: "使用参数化查询或预编译语句来防止SQL注入",
						Confidence:     0.75,
					}
					issues = append(issues, issue)
				}
			}
		}
	}

	return issues
}

// XSSDetector XSS漏洞检测器
type XSSDetector struct {
	*BaseDetector
}

// NewXSSDetector 创建XSS检测器
func NewXSSDetector() *XSSDetector {
	return &XSSDetector{
		BaseDetector: NewBaseDetector("xss_vulnerabilities", []string{"go", "python", "javascript", "java", "php", "ruby"}),
	}
}

// Detect 检测XSS漏洞
func (d *XSSDetector) Detect(filePath string, content string) []types.SecurityIssue {
	var issues []types.SecurityIssue

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)innerHTML\s*=\s*`),
		regexp.MustCompile(`(?i)document\.write\(`),
		regexp.MustCompile(`(?i)eval\(`),
		regexp.MustCompile(`(?i)setTimeout\(.*\+`),
		regexp.MustCompile(`(?i)setInterval\(.*\+`),
	}

	lines := strings.Split(content, "\n")
	for lineNum, line := range lines {
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				issue := types.SecurityIssue{
					ID:             "XSS_001",
					Type:           "XSSVulnerability",
					Severity:       types.SeverityMedium,
					Message:        "检测到可能的XSS漏洞",
					File:           filePath,
					Line:           lineNum + 1,
					Column:         strings.Index(line, pattern.FindString(line)) + 1,
					Snippet:        strings.TrimSpace(line),
					Recommendation: "对用户输入进行适当的转义和验证",
					Confidence:     0.70,
				}
				issues = append(issues, issue)
			}
		}
	}

	return issues
}

// PathTraversalDetector 路径遍历检测器
type PathTraversalDetector struct {
	*BaseDetector
}

// NewPathTraversalDetector 创建路径遍历检测器
func NewPathTraversalDetector() *PathTraversalDetector {
	return &PathTraversalDetector{
		BaseDetector: NewBaseDetector("path_traversal", []string{"go", "python", "javascript", "java", "php", "ruby"}),
	}
}

// Detect 检测路径遍历漏洞
func (d *PathTraversalDetector) Detect(filePath string, content string) []types.SecurityIssue {
	var issues []types.SecurityIssue

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)\.\./`),
		regexp.MustCompile(`(?i)\.\.\\\\`),
		regexp.MustCompile(`(?i)file:.*\+`),
		regexp.MustCompile(`(?i)open\(.*\+`),
		regexp.MustCompile(`(?i)readFile\(.*\+`),
	}

	lines := strings.Split(content, "\n")
	for lineNum, line := range lines {
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				issue := types.SecurityIssue{
					ID:             "PATH_TRAVERSAL_001",
					Type:           "PathTraversal",
					Severity:       types.SeverityHigh,
					Message:        "检测到可能的路径遍历漏洞",
					File:           filePath,
					Line:           lineNum + 1,
					Column:         strings.Index(line, pattern.FindString(line)) + 1,
					Snippet:        strings.TrimSpace(line),
					Recommendation: "对文件路径进行规范化处理，并限制访问范围",
					Confidence:     0.80,
				}
				issues = append(issues, issue)
			}
		}
	}

	return issues
}

// QualityDetector 代码质量检测器
type QualityDetector struct {
	*BaseDetector
}

// NewQualityDetector 创建代码质量检测器
func NewQualityDetector() *QualityDetector {
	return &QualityDetector{
		BaseDetector: NewBaseDetector("code_quality", []string{"go", "python", "javascript", "java", "php", "ruby"}),
	}
}

// Detect 检测代码质量问题
func (d *QualityDetector) Detect(filePath string, content string) []types.SecurityIssue {
	var issues []types.SecurityIssue

	// 检测未使用的变量
	unusedVarPattern := regexp.MustCompile(`(?i)(var|let|const)\s+(\w+)\s*=`)
	lines := strings.Split(content, "\n")

	for lineNum, line := range lines {
		// 检测未使用的变量
		if unusedVarPattern.MatchString(line) {
			matches := unusedVarPattern.FindStringSubmatch(line)
			if len(matches) > 2 {
				varName := matches[2]
				// 简单检查变量是否被使用
				if !strings.Contains(content, varName) || strings.Count(content, varName) == 1 {
					// 检查是否已经有相同的未使用变量问题
					found := false
					for _, existing := range issues {
						if existing.ID == "QUALITY_001" && existing.Line == lineNum+1 {
							found = true
							break
						}
					}
					if !found {
						issue := types.SecurityIssue{
							ID:             "QUALITY_001",
							Type:           "UnusedVariable",
							Severity:       types.SeverityLow,
							Message:        "检测到未使用的变量",
							File:           filePath,
							Line:           lineNum + 1,
							Column:         strings.Index(line, varName) + 1,
							Snippet:        strings.TrimSpace(line),
							Recommendation: "移除未使用的变量或确保其被正确使用",
							Confidence:     0.65,
						}
						issues = append(issues, issue)
					}
				}
			}
		}

		// 检测空的错误处理 - 检测注释掉的错误处理
		if strings.Contains(line, "//") && strings.Contains(line, "错误") && strings.Contains(line, "忽略") {
			// 检查是否已经有相同的错误处理问题
			found := false
			for _, existing := range issues {
				if existing.ID == "QUALITY_002" && existing.Line == lineNum+1 {
					found = true
					break
				}
			}
			if !found {
				issue := types.SecurityIssue{
					ID:             "QUALITY_002",
					Type:           "IncompleteErrorHandling",
					Severity:       types.SeverityMedium,
					Message:        "检测到不完整的错误处理",
					File:           filePath,
					Line:           lineNum + 1,
					Column:         1,
					Snippet:        strings.TrimSpace(line),
					Recommendation: "确保错误被正确处理，避免忽略错误",
					Confidence:     0.70,
				}
				issues = append(issues, issue)
			}
		}
	}

	return issues
}

// DetectorRegistry 检测器注册表
type DetectorRegistry struct {
	detectors map[string]types.SecurityDetector
}

// NewDetectorRegistry 创建检测器注册表
func NewDetectorRegistry() *DetectorRegistry {
	registry := &DetectorRegistry{
		detectors: make(map[string]types.SecurityDetector),
	}

	// 注册所有检测器
	registry.Register(NewCredentialsDetector())
	registry.Register(NewSQLInjectionDetector())
	registry.Register(NewXSSDetector())
	registry.Register(NewPathTraversalDetector())
	registry.Register(NewQualityDetector())

	return registry
}

// Register 注册检测器
func (r *DetectorRegistry) Register(detector types.SecurityDetector) {
	r.detectors[detector.GetName()] = detector
}

// GetDetector 获取检测器
func (r *DetectorRegistry) GetDetector(name string) types.SecurityDetector {
	return r.detectors[name]
}

// GetAllDetectors 获取所有检测器
func (r *DetectorRegistry) GetAllDetectors() []types.SecurityDetector {
	detectors := make([]types.SecurityDetector, 0, len(r.detectors))
	for _, detector := range r.detectors {
		detectors = append(detectors, detector)
	}
	return detectors
}

// GetDetectorsForLanguage 获取支持特定语言的检测器
func (r *DetectorRegistry) GetDetectorsForLanguage(language string) []types.SecurityDetector {
	var supported []types.SecurityDetector

	for _, detector := range r.detectors {
		for _, lang := range detector.GetSupportedLanguages() {
			if strings.EqualFold(lang, language) {
				supported = append(supported, detector)
				break
			}
		}
	}

	return supported
}
