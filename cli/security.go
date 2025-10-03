// Package main CLI安全扫描命令
package main

import (
	"fmt"
	"os"
	"strings"

	"code-context-generator/internal/config"
	"code-context-generator/pkg/security"
	"code-context-generator/pkg/types"

	"github.com/spf13/cobra"
)

// securityCmd 安全扫描命令
var securityCmd = &cobra.Command{
	Use:   "security [路径]",
	Short: "执行代码安全扫描",
	Long: `执行代码安全扫描，检测潜在的安全漏洞和代码质量问题

支持检测：
- 硬编码凭证
- SQL注入漏洞
- XSS漏洞
- 路径遍历漏洞
- 代码质量问题

支持多种编程语言：Go, Python, JavaScript, Java, PHP, Ruby等。`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSecurityScan,
}

// securityConfigCmd 安全配置命令
var securityConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "安全扫描配置",
	Long:  "管理安全扫描配置",
}

// securityConfigShowCmd 显示安全配置
var securityConfigShowCmd = &cobra.Command{
	Use:   "show",
	Short: "显示当前安全配置",
	Long:  "显示当前安全扫描配置设置",
	RunE:  runSecurityConfigShow,
}

// securityConfigInitCmd 初始化安全配置
var securityConfigInitCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化安全配置",
	Long:  "创建默认安全扫描配置文件",
	RunE:  runSecurityConfigInit,
}

// initSecurityCommands 初始化安全扫描命令
func initSecurityCommands() {
	// 添加安全扫描命令
	rootCmd.AddCommand(securityCmd)

	// 安全配置子命令
	securityCmd.AddCommand(securityConfigCmd)
	securityConfigCmd.AddCommand(securityConfigShowCmd)
	securityConfigCmd.AddCommand(securityConfigInitCmd)

	// 安全扫描命令标志
	securityCmd.Flags().Bool("enabled", true, "启用安全扫描")
	securityCmd.Flags().Bool("fail-on-critical", false, "发现严重问题时退出码为非零")
	securityCmd.Flags().String("scan-level", "standard", "扫描级别 (basic, standard, comprehensive)")
	securityCmd.Flags().String("report-format", "text", "报告格式 (text, json, xml, html)")
	securityCmd.Flags().String("output-file", "", "输出报告文件路径")
	securityCmd.Flags().Bool("include-details", true, "包含详细问题信息")
	securityCmd.Flags().Bool("show-statistics", true, "显示扫描统计信息")

	// 检测器配置标志
	securityCmd.Flags().Bool("detect-credentials", true, "检测硬编码凭证")
	securityCmd.Flags().Bool("detect-sql-injection", true, "检测SQL注入漏洞")
	securityCmd.Flags().Bool("detect-xss", true, "检测XSS漏洞")
	securityCmd.Flags().Bool("detect-path-traversal", true, "检测路径遍历漏洞")
	securityCmd.Flags().Bool("detect-quality", true, "检测代码质量问题")

	// 排除配置标志
	securityCmd.Flags().StringSlice("exclude-files", []string{}, "排除的文件列表")
	securityCmd.Flags().StringSlice("exclude-patterns", []string{}, "排除的文件模式")
}

// runSecurityScan 运行安全扫描
func runSecurityScan(cmd *cobra.Command, args []string) error {
	// 解析标志
	enabled, _ := cmd.Flags().GetBool("enabled")
	failOnCritical, _ := cmd.Flags().GetBool("fail-on-critical")
	scanLevelStr, _ := cmd.Flags().GetString("scan-level")
	reportFormat, _ := cmd.Flags().GetString("report-format")
	outputFile, _ := cmd.Flags().GetString("output-file")
	includeDetails, _ := cmd.Flags().GetBool("include-details")
	showStatistics, _ := cmd.Flags().GetBool("show-statistics")

	detectCredentials, _ := cmd.Flags().GetBool("detect-credentials")
	detectSQLInjection, _ := cmd.Flags().GetBool("detect-sql-injection")
	detectXSS, _ := cmd.Flags().GetBool("detect-xss")
	detectPathTraversal, _ := cmd.Flags().GetBool("detect-path-traversal")
	detectQuality, _ := cmd.Flags().GetBool("detect-quality")

	excludeFiles, _ := cmd.Flags().GetStringSlice("exclude-files")
	excludePatterns, _ := cmd.Flags().GetStringSlice("exclude-patterns")

	// 获取扫描路径
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// 检查路径是否存在
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("路径不存在: %v", err)
	}

	// 创建安全配置
	securityConfig := createSecurityConfig(
		enabled, failOnCritical, scanLevelStr, reportFormat, outputFile,
		includeDetails, showStatistics, detectCredentials, detectSQLInjection,
		detectXSS, detectPathTraversal, detectQuality, excludeFiles, excludePatterns,
	)

	// 创建安全管理器
	manager := security.NewSecurityManager(securityConfig)

	// 执行扫描
	fmt.Println("开始安全扫描...")
	report, err := manager.RunScan(path)
	if err != nil {
		return fmt.Errorf("安全扫描失败: %v", err)
	}

	// 生成报告
	reportContent, err := manager.GenerateReport(report)
	if err != nil {
		return fmt.Errorf("生成报告失败: %v", err)
	}

	// 输出报告
	if outputFile != "" {
		// 写入文件
		err = os.WriteFile(outputFile, reportContent, 0644)
		if err != nil {
			return fmt.Errorf("写入报告文件失败: %v", err)
		}
		fmt.Printf("安全报告已保存到: %s\n", outputFile)
	} else {
		// 输出到控制台
		fmt.Println(string(reportContent))
	}

	// 检查是否需要因严重问题而退出
	if failOnCritical && report.Summary.CriticalIssues > 0 {
		fmt.Printf("发现 %d 个严重问题，退出码为1\n", report.Summary.CriticalIssues)
		os.Exit(1)
	}

	return nil
}

// createSecurityConfig 创建安全配置
func createSecurityConfig(
	enabled, failOnCritical bool,
	scanLevelStr, reportFormat, outputFile string,
	includeDetails, showStatistics bool,
	detectCredentials, detectSQLInjection, detectXSS, detectPathTraversal, detectQuality bool,
	excludeFiles, excludePatterns []string,
) *types.SecurityConfig {
	// 扫描级别字符串直接使用

	// 创建检测器配置
	detectorConfig := types.DetectorConfig{
		Credentials:   detectCredentials,
		SQLInjection:  detectSQLInjection,
		XSS:           detectXSS,
		PathTraversal: detectPathTraversal,
		Quality:       detectQuality,
	}

	// 创建排除配置
	exclusionConfig := types.ExclusionConfig{
		Files:    excludeFiles,
		Patterns: excludePatterns,
		Rules:    []types.ExclusionRule{},
	}

	// 创建报告配置
	reportingConfig := types.ReportingConfig{
		Format:         reportFormat,
		OutputFile:     outputFile,
		IncludeDetails: includeDetails,
		ShowStatistics: showStatistics,
	}

	return &types.SecurityConfig{
		Enabled:        enabled,
		FailOnCritical: failOnCritical,
		ScanLevel:      scanLevelStr,
		Detectors:      detectorConfig,
		Exclusions:     exclusionConfig,
		Reporting:      reportingConfig,
	}
}

// runSecurityConfigShow 显示安全配置
func runSecurityConfigShow(cmd *cobra.Command, args []string) error {
	fmt.Println("当前安全配置:")
	fmt.Printf("启用安全扫描: %v\n", cfg.Security.Enabled)
	fmt.Printf("发现严重问题时退出: %v\n", cfg.Security.FailOnCritical)
	fmt.Printf("扫描级别: %v\n", cfg.Security.ScanLevel)

	fmt.Println("\n检测器配置:")
	fmt.Printf("硬编码凭证检测: %v\n", cfg.Security.Detectors.Credentials)
	fmt.Printf("SQL注入检测: %v\n", cfg.Security.Detectors.SQLInjection)
	fmt.Printf("XSS漏洞检测: %v\n", cfg.Security.Detectors.XSS)
	fmt.Printf("路径遍历检测: %v\n", cfg.Security.Detectors.PathTraversal)
	fmt.Printf("代码质量检测: %v\n", cfg.Security.Detectors.Quality)

	fmt.Println("\n报告配置:")
	fmt.Printf("报告格式: %s\n", cfg.Security.Reporting.Format)
	fmt.Printf("输出文件: %s\n", cfg.Security.Reporting.OutputFile)
	fmt.Printf("包含详细信息: %v\n", cfg.Security.Reporting.IncludeDetails)
	fmt.Printf("显示统计信息: %v\n", cfg.Security.Reporting.ShowStatistics)

	return nil
}

// runSecurityConfigInit 初始化安全配置
func runSecurityConfigInit(cmd *cobra.Command, args []string) error {
	// 创建默认安全配置
	defaultSecurityConfig := &types.SecurityConfig{
		Enabled:        true,
		FailOnCritical: false,
		ScanLevel:      "standard",
		Detectors: types.DetectorConfig{
			Credentials:   true,
			SQLInjection:  true,
			XSS:           true,
			PathTraversal: true,
			Quality:       true,
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

	// 更新配置
	cfg.Security = *defaultSecurityConfig

	// 保存配置
	configManager := config.NewManager()

	// 确定配置文件路径和格式
	savePath := configPath
	if savePath == "" {
		savePath = "config.yaml"
	}

	// 确定文件格式
	format := "yaml"
	if strings.HasSuffix(savePath, ".json") {
		format = "json"
	} else if strings.HasSuffix(savePath, ".toml") {
		format = "toml"
	}

	// 保存配置
	if err := configManager.Save(savePath, format); err != nil {
		return fmt.Errorf("保存配置文件失败: %v", err)
	}

	fmt.Printf("安全配置已初始化并保存到: %s\n", savePath)
	return nil
}

// init 初始化函数 - 添加安全扫描命令
func init() {
	initSecurityCommands()
}
