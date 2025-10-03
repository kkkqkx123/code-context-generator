package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"code-context-generator/pkg/types"
)

// SecurityIntegration 安全扫描集成器
type SecurityIntegration struct {
	scanner *SecurityScanner
	config  *types.SecurityConfig
	enabled bool
}

// NewSecurityIntegration 创建安全扫描集成器
func NewSecurityIntegration(config *types.SecurityConfig) *SecurityIntegration {
	return &SecurityIntegration{
		scanner: NewSecurityScanner(config),
		config:  config,
		enabled: config.Enabled,
	}
}

// ScanProject 扫描整个项目
func (si *SecurityIntegration) ScanProject(projectPath string) (*types.SecurityReport, error) {
	if !si.enabled {
		scanID, _ := generateScanID()
		return &types.SecurityReport{
			ScanID:    scanID,
			Timestamp: time.Now(),
			Summary: types.ScanSummary{
				TotalFiles:   0,
				ScannedFiles: 0,
				IssuesFound:  0,
			},
			Issues: []types.SecurityIssue{},
		}, nil
	}

	return si.scanner.Scan(projectPath)
}

// ScanFiles 扫描多个文件
func (si *SecurityIntegration) ScanFiles(files []string) (*types.SecurityReport, error) {
	if !si.enabled {
		scanID, _ := generateScanID()
		return &types.SecurityReport{
			ScanID:       scanID,
			Timestamp:    time.Now(),
			ScanDuration: 0,
			Summary: types.ScanSummary{
				TotalFiles:   len(files),
				ScannedFiles: 0,
				IssuesFound:  0,
			},
			Issues: []types.SecurityIssue{},
		}, nil
	}

	var allIssues []types.SecurityIssue
	startTime := time.Now()
	var scannedFiles []string

	for _, file := range files {
		if !si.isSupportedExtension(file) {
			continue
		}

		// 扫描单个文件
		fileReport, err := si.scanner.Scan(file)
		if err != nil {
			continue // 跳过无法扫描的文件
		}

		allIssues = append(allIssues, fileReport.Issues...)
		scannedFiles = append(scannedFiles, file)
	}

	scanDuration := time.Since(startTime)
	scanID, _ := generateScanID()

	// 统计问题数量
	var critical, high, medium, low int
	for _, issue := range allIssues {
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

	summary := types.ScanSummary{
		TotalFiles:     len(files),
		ScannedFiles:   len(scannedFiles),
		IssuesFound:    len(allIssues),
		CriticalIssues: critical,
		HighIssues:     high,
		MediumIssues:   medium,
		LowIssues:      low,
	}

	report := &types.SecurityReport{
		ScanID:       scanID,
		Timestamp:    startTime,
		ScanDuration: scanDuration,
		Summary:      summary,
		Issues:       allIssues,
	}

	return report, nil
}

// shouldScanFile 判断是否应该扫描文件
func (si *SecurityIntegration) shouldScanFile(file string) bool {
	// 检查文件是否存在
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}

	// 检查文件扩展名
	ext := filepath.Ext(file)
	if !si.isSupportedExtension(ext) {
		return false
	}

	return true
}

// isSupportedExtension 检查是否支持的文件扩展名
func (si *SecurityIntegration) isSupportedExtension(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	supportedExts := []string{".go", ".py", ".js", ".ts", ".java", ".php", ".rb", ".cpp", ".c", ".cs", ".swift", ".rs", ".yml", ".yaml", ".json", ".xml", ".toml", ".ini", ".cfg", ".conf"}
	
	for _, supported := range supportedExts {
		if ext == supported {
			return true
		}
	}
	return false
}

// GenerateReport 生成安全报告
func (si *SecurityIntegration) GenerateReport(report *types.SecurityReport, outputPath string) error {
	if !si.enabled {
		return nil
	}

	reporter := NewSecurityReporter()
	reportData, err := reporter.Generate(report)
	if err != nil {
		return fmt.Errorf("生成报告失败: %v", err)
	}

	if outputPath != "" {
		return os.WriteFile(outputPath, reportData, 0644)
	}

	// 如果未指定输出路径，则打印到控制台
	fmt.Print(string(reportData))
	return nil
}

// HasCriticalIssues 检查是否存在严重问题
func (si *SecurityIntegration) HasCriticalIssues(report *types.SecurityReport) bool {
	if !si.enabled {
		return false
	}

	for _, issue := range report.Issues {
		if issue.Severity == types.SeverityCritical || issue.Severity == types.SeverityHigh {
			return true
		}
	}
	return false
}

// PrintSummary 打印扫描摘要
func (si *SecurityIntegration) PrintSummary(report *types.SecurityReport) {
	if !si.enabled {
		return
	}

	fmt.Printf("\n🔒 安全扫描完成\n")
	fmt.Printf("📊 扫描摘要:\n")
	fmt.Printf("  📁 扫描文件: %d\n", report.Summary.ScannedFiles)
	fmt.Printf("  🔍 发现问题: %d\n", report.Summary.IssuesFound)
	fmt.Printf("  ⚠️  严重问题: %d\n", report.Summary.CriticalIssues)
	fmt.Printf("  🔴 高危问题: %d\n", report.Summary.HighIssues)
	fmt.Printf("  🟡 中危问题: %d\n", report.Summary.MediumIssues)
	fmt.Printf("  🟢 低危问题: %d\n", report.Summary.LowIssues)
	fmt.Printf("  ⏱️  扫描时间: %s\n", report.ScanDuration.String())
	fmt.Printf("  📅 扫描时间: %s\n", report.Timestamp.Format("2006-01-02 15:04:05"))

	if report.Summary.IssuesFound > 0 {
		fmt.Printf("\n💡 建议查看详细报告以了解具体问题\n")
	}
}
