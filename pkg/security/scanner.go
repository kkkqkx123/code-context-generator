// Package security 实现安全扫描功能
package security

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"code-context-generator/pkg/types"
)

// SecurityScanner 安全扫描器实现
type SecurityScanner struct {
	config   *types.SecurityConfig
	registry *DetectorRegistry
}

// NewSecurityScanner 创建安全扫描器
func NewSecurityScanner(config *types.SecurityConfig) *SecurityScanner {
	return &SecurityScanner{
		config:   config,
		registry: NewDetectorRegistry(),
	}
}

// Scan 执行安全扫描
func (s *SecurityScanner) Scan(path string) (*types.SecurityReport, error) {
	startTime := time.Now()

	// 生成扫描ID
	scanID, err := generateScanID()
	if err != nil {
		return nil, fmt.Errorf("生成扫描ID失败: %v", err)
	}

	// 检查路径是否存在
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("路径不存在: %v", err)
	}

	var files []string
	if fileInfo.IsDir() {
		// 扫描目录
		files, err = s.scanDirectory(path)
		if err != nil {
			return nil, fmt.Errorf("扫描目录失败: %v", err)
		}
	} else {
		// 扫描单个文件
		files = []string{path}
	}

	// 执行扫描
	issues := s.scanFiles(files)

	// 计算扫描时间
	scanDuration := time.Since(startTime)

	// 生成报告
	report := s.generateReport(scanID, files, issues, scanDuration)

	return report, nil
}

// scanDirectory 扫描目录
func (s *SecurityScanner) scanDirectory(dirPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查文件是否在排除列表中
		if s.isExcluded(path) {
			return nil
		}

		// 检查文件扩展名
		if s.isSupportedFile(path) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// isExcluded 检查文件是否被排除
func (s *SecurityScanner) isExcluded(filePath string) bool {
	// 检查文件排除列表
	for _, excludedFile := range s.config.Exclusions.Files {
		if strings.Contains(filePath, excludedFile) {
			return true
		}
	}

	// 检查模式排除
	for _, pattern := range s.config.Exclusions.Patterns {
		matched, err := filepath.Match(pattern, filepath.Base(filePath))
		if err == nil && matched {
			return true
		}
	}

	return false
}

// isSupportedFile 检查文件是否支持
func (s *SecurityScanner) isSupportedFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	supportedExtensions := map[string]bool{
		".go":    true,
		".py":    true,
		".js":    true,
		".ts":    true,
		".java":  true,
		".php":   true,
		".rb":    true,
		".cpp":   true,
		".c":     true,
		".h":     true,
		".cs":    true,
		".swift": true,
		".rs":    true,
		".yml":   true,
		".yaml":  true,
		".json":  true,
		".xml":   true,
		".toml":  true,
		".ini":   true,
		".cfg":   true,
		".conf":  true,
	}

	return supportedExtensions[ext]
}

// scanFiles 扫描文件
func (s *SecurityScanner) scanFiles(files []string) []types.SecurityIssue {
	var allIssues []types.SecurityIssue

	for _, filePath := range files {
		// 读取文件内容
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			// 记录错误但继续扫描其他文件
			fmt.Printf("警告: 无法读取文件 %s: %v\n", filePath, err)
			continue
		}

		// 获取文件语言
		language := s.getFileLanguage(filePath)

		// 获取适用于该语言的检测器
		detectors := s.registry.GetDetectorsForLanguage(language)

		// 执行检测
		for _, detector := range detectors {
			issues := detector.Detect(filePath, string(content))
			allIssues = append(allIssues, issues...)
		}
	}

	return allIssues
}

// getFileLanguage 获取文件语言
func (s *SecurityScanner) getFileLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	languageMap := map[string]string{
		".go":    "go",
		".py":    "python",
		".js":    "javascript",
		".ts":    "javascript",
		".java":  "java",
		".php":   "php",
		".rb":    "ruby",
		".cpp":   "cpp",
		".c":     "c",
		".h":     "c",
		".cs":    "csharp",
		".swift": "swift",
		".rs":    "rust",
		".yml":   "yaml",
		".yaml":  "yaml",
		".json":  "json",
		".xml":   "xml",
		".toml":  "toml",
		".ini":   "ini",
		".cfg":   "config",
		".conf":  "config",
	}

	if lang, exists := languageMap[ext]; exists {
		return lang
	}

	return "unknown"
}

// generateReport 生成安全报告
func (s *SecurityScanner) generateReport(scanID string, files []string, issues []types.SecurityIssue, duration time.Duration) *types.SecurityReport {
	// 统计问题数量
	var critical, high, medium, low int
	for _, issue := range issues {
		switch issue.Severity {
		case types.SeverityCritical:
			critical++
		case types.SeverityHigh:
			high++
		case types.SeverityMedium:
			medium++
		case types.SeverityLow:
			low++
		}
	}

	// 计算统计信息
	totalFiles := len(files)
	scannedFiles := len(files)
	totalIssues := len(issues)

	// 计算性能统计
	filesPerSec := float64(scannedFiles) / duration.Seconds()

	summary := types.ScanSummary{
		TotalFiles:     totalFiles,
		ScannedFiles:   scannedFiles,
		IssuesFound:    totalIssues,
		CriticalIssues: critical,
		HighIssues:     high,
		MediumIssues:   medium,
		LowIssues:      low,
	}

	statistics := types.ScanStatistics{
		TotalTime:   duration,
		AverageTime: duration / time.Duration(scannedFiles),
		FilesPerSec: filesPerSec,
		MemoryUsage: 0, // 暂时不实现内存统计
	}

	report := &types.SecurityReport{
		ScanID:       scanID,
		Timestamp:    time.Now(),
		ScanDuration: duration,
		Summary:      summary,
		Issues:       issues,
		Statistics:   statistics,
		Config:       *s.config,
	}

	return report
}

// GetConfig 获取配置
func (s *SecurityScanner) GetConfig() *types.SecurityConfig {
	return s.config
}

// SetConfig 设置配置
func (s *SecurityScanner) SetConfig(config *types.SecurityConfig) {
	s.config = config
}

// generateScanID 生成扫描ID
func generateScanID() (string, error) {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SecurityReporterImpl 安全报告器实现
type SecurityReporterImpl struct{}

// NewSecurityReporter 创建安全报告器
func NewSecurityReporter() *SecurityReporterImpl {
	return &SecurityReporterImpl{}
}

// Generate 生成报告
func (r *SecurityReporterImpl) Generate(report *types.SecurityReport) ([]byte, error) {
	// 这里实现报告生成逻辑
	// 暂时返回简单文本报告
	return r.generateTextReport(report), nil
}

// generateTextReport 生成文本报告
func (r *SecurityReporterImpl) generateTextReport(report *types.SecurityReport) []byte {
	var builder strings.Builder

	builder.WriteString("=== 安全扫描报告 ===\n")
	builder.WriteString(fmt.Sprintf("扫描ID: %s\n", report.ScanID))
	builder.WriteString(fmt.Sprintf("扫描时间: %s\n", report.Timestamp.Format("2006-01-02 15:04:05")))
	builder.WriteString(fmt.Sprintf("扫描耗时: %v\n", report.ScanDuration))
	builder.WriteString("\n")

	// 摘要信息
	builder.WriteString("=== 扫描摘要 ===\n")
	builder.WriteString(fmt.Sprintf("总文件数: %d\n", report.Summary.TotalFiles))
	builder.WriteString(fmt.Sprintf("已扫描文件: %d\n", report.Summary.ScannedFiles))
	builder.WriteString(fmt.Sprintf("发现问题: %d\n", report.Summary.IssuesFound))
	builder.WriteString(fmt.Sprintf("严重问题: %d\n", report.Summary.CriticalIssues))
	builder.WriteString(fmt.Sprintf("高危问题: %d\n", report.Summary.HighIssues))
	builder.WriteString(fmt.Sprintf("中危问题: %d\n", report.Summary.MediumIssues))
	builder.WriteString(fmt.Sprintf("低危问题: %d\n", report.Summary.LowIssues))
	builder.WriteString("\n")

	// 统计信息
	builder.WriteString("=== 性能统计 ===\n")
	builder.WriteString(fmt.Sprintf("总扫描时间: %v\n", report.Statistics.TotalTime))
	builder.WriteString(fmt.Sprintf("平均文件扫描时间: %v\n", report.Statistics.AverageTime))
	builder.WriteString(fmt.Sprintf("文件扫描速度: %.2f 文件/秒\n", report.Statistics.FilesPerSec))
	builder.WriteString("\n")

	// 问题详情
	if len(report.Issues) > 0 {
		builder.WriteString("=== 问题详情 ===\n")
		for i, issue := range report.Issues {
			builder.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, issue.Severity.String(), issue.Type))
			builder.WriteString(fmt.Sprintf("   文件: %s:%d:%d\n", issue.File, issue.Line, issue.Column))
			builder.WriteString(fmt.Sprintf("   描述: %s\n", issue.Message))
			builder.WriteString(fmt.Sprintf("   代码: %s\n", issue.Snippet))
			builder.WriteString(fmt.Sprintf("   建议: %s\n", issue.Recommendation))
			builder.WriteString(fmt.Sprintf("   置信度: %.2f\n", issue.Confidence))
			builder.WriteString("\n")
		}
	} else {
		builder.WriteString("=== 扫描结果 ===\n")
		builder.WriteString("未发现安全问题\n")
	}

	return []byte(builder.String())
}

// GetSupportedFormats 获取支持的格式
func (r *SecurityReporterImpl) GetSupportedFormats() []string {
	return []string{"text", "json", "xml", "html"}
}

// SecurityManager 安全管理器
type SecurityManager struct {
	scanner  *SecurityScanner
	reporter *SecurityReporterImpl
}

// NewSecurityManager 创建安全管理器
func NewSecurityManager(config *types.SecurityConfig) *SecurityManager {
	return &SecurityManager{
		scanner:  NewSecurityScanner(config),
		reporter: NewSecurityReporter(),
	}
}

// RunScan 运行安全扫描
func (m *SecurityManager) RunScan(path string) (*types.SecurityReport, error) {
	report, err := m.scanner.Scan(path)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// GenerateReport 生成报告
func (m *SecurityManager) GenerateReport(report *types.SecurityReport) ([]byte, error) {
	return m.reporter.Generate(report)
}

// GetScanner 获取扫描器
func (m *SecurityManager) GetScanner() *SecurityScanner {
	return m.scanner
}

// GetReporter 获取报告器
func (m *SecurityManager) GetReporter() *SecurityReporterImpl {
	return m.reporter
}
