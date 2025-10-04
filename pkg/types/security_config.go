// Package types 定义安全扫描相关的类型和接口
package types

import (
	"time"
)

// SecurityConfig 安全扫描配置结构体
type SecurityConfig struct {
	Enabled        bool     `yaml:"enabled"`
	FailOnCritical bool     `yaml:"fail_on_critical"`
	ScanLevel      string   `yaml:"scan_level"`
	ReportFormat   string   `yaml:"report_format"`
	
	Detectors      DetectorConfig    `yaml:"detectors"`
	Exclusions     ExclusionConfig   `yaml:"exclusions"`
	Reporting      ReportingConfig   `yaml:"reporting"`
}

// ScanLevel 扫描级别
type ScanLevel int

const (
	ScanLevelBasic ScanLevel = iota // 基础扫描
	ScanLevelStandard               // 标准扫描
	ScanLevelComprehensive          // 全面扫描
)

// DetectorConfig 检测器配置
type DetectorConfig struct {
	Credentials    bool `yaml:"credentials"`
	SQLInjection   bool `yaml:"sql_injection"`
	XSS            bool `yaml:"xss"`
	PathTraversal  bool `yaml:"path_traversal"`
	Quality        bool `yaml:"quality"`
}

// HardcodedCredentialsConfig 硬编码凭证检测配置
type HardcodedCredentialsConfig struct {
	Enabled          bool     `yaml:"enabled"`
	SeverityThreshold SeverityLevel `yaml:"severity_threshold"`
	Patterns         []string `yaml:"patterns"`
}

// VulnerabilityConfig 安全漏洞检测配置
type VulnerabilityConfig struct {
	Enabled          bool     `yaml:"enabled"`
	SeverityThreshold SeverityLevel `yaml:"severity_threshold"`
}

// QualityConfig 代码质量检测配置
type QualityConfig struct {
	Enabled          bool     `yaml:"enabled"`
	SeverityThreshold SeverityLevel `yaml:"severity_threshold"`
}

// ExclusionConfig 排除配置
type ExclusionConfig struct {
	Files    []string         `yaml:"files"`
	Patterns []string         `yaml:"patterns"`
	Rules    []ExclusionRule `yaml:"rules"`
}

// ExclusionRule 排除规则
type ExclusionRule struct {
	Pattern string `yaml:"pattern"`
	Reason  string `yaml:"reason"`
}

// ReportingConfig 报告配置
type ReportingConfig struct {
	Format         string `yaml:"format"`
	OutputFile     string `yaml:"output_file"`
	IncludeDetails bool   `yaml:"include_details"`
	ShowStatistics bool   `yaml:"show_statistics"`
}

// SecurityReport 安全报告结构体
type SecurityReport struct {
	ScanID       string          `json:"scan_id"`
	Timestamp    time.Time       `json:"timestamp"`
	ScanDuration time.Duration   `json:"scan_duration"`
	
	Summary      ScanSummary     `json:"summary"`
	Issues       []SecurityIssue `json:"issues"`
	Statistics   ScanStatistics  `json:"statistics"`
	
	Config       SecurityConfig  `json:"config"`
}

// ScanSummary 扫描摘要
type ScanSummary struct {
	TotalFiles     int `json:"total_files"`
	ScannedFiles   int `json:"scanned_files"`
	IssuesFound    int `json:"issues_found"`
	CriticalIssues int `json:"critical_issues"`
	HighIssues     int `json:"high_issues"`
	MediumIssues   int `json:"medium_issues"`
	LowIssues      int `json:"low_issues"`
}

// ScanStatistics 扫描统计
type ScanStatistics struct {
	TotalTime    time.Duration `json:"total_time"`
	AverageTime  time.Duration `json:"average_time"`
	FilesPerSec  float64       `json:"files_per_sec"`
	MemoryUsage  int64         `json:"memory_usage"`
}

// SecurityIssue 安全问题
type SecurityIssue struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Severity      SeverityLevel  `json:"severity"`
	Message       string        `json:"message"`
	File          string        `json:"file"`
	Line          int           `json:"line"`
	Column        int           `json:"column"`
	Snippet       string        `json:"snippet"`
	Recommendation string        `json:"recommendation"`
	Confidence    float64       `json:"confidence"`
}

// SeverityLevel 严重性级别
type SeverityLevel int

const (
	SeverityLow SeverityLevel = iota
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

// String 返回严重性级别的字符串表示
func (s SeverityLevel) String() string {
	switch s {
	case SeverityLow:
		return "low"
	case SeverityMedium:
		return "medium"
	case SeverityHigh:
		return "high"
	case SeverityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// SecurityDetector 安全检测器接口
type SecurityDetector interface {
	Detect(filePath string, content string) []SecurityIssue
	GetName() string
	GetSupportedLanguages() []string
}

// SecurityReporter 安全报告器接口
type SecurityReporter interface {
	Generate(report *SecurityReport) ([]byte, error)
	GetSupportedFormats() []string
}

// SecurityScanner 安全扫描器接口
type SecurityScanner interface {
	Scan(path string) (*SecurityReport, error)
	GetConfig() *SecurityConfig
	SetConfig(config *SecurityConfig)
}