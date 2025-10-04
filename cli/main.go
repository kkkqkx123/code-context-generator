// Package main CLI应用程序主入口
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code-context-generator/internal/config"
	"code-context-generator/internal/env"
	"code-context-generator/internal/filesystem"
	"code-context-generator/internal/formatter"
	"code-context-generator/internal/git"
	"code-context-generator/internal/utils"
	"code-context-generator/pkg/security"
	"code-context-generator/pkg/types"

	"github.com/spf13/cobra"
)

var (
	// 全局变量
	cfg        *types.Config
	configPath string
	verbose    bool
	version    = "1.0.0"
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "code-context-generator [路径]",
	Short: "代码上下文生成器",
	Long: `代码上下文生成器 - 智能生成代码项目结构文档

支持多种输出格式（JSON、XML、TOML、Markdown），提供自动文件扫描，
自动补全功能，以及丰富的配置选项。`,
	Version: version,
	Args:    cobra.MaximumNArgs(1), // 接受一个可选的路径参数
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 首先加载.env文件（如果存在）
		if err := env.LoadEnv(""); err != nil {
			fmt.Printf("警告: 加载.env文件失败: %v\n", err)
		}

		// 加载配置
		configManager := config.NewManager()

		// 如果有指定配置文件路径，使用它
		if configPath != "" {
			if err := configManager.Load(configPath); err != nil {
				return fmt.Errorf("加载配置文件失败: %w", err)
			}
		} else {
			// 尝试加载默认配置文件，如果不存在则使用默认配置，不再自动创建
			defaultConfigPath := "config.yaml"
			configManager.Load(defaultConfigPath) // 忽略错误，使用默认配置
		}

		cfg = configManager.Get()
		return nil
	},
	RunE: runGenerate, // 默认执行生成命令
}

// generateCmd 生成命令 (现在为可选命令，保持向后兼容)
var generateCmd = &cobra.Command{
	Use:   "generate [路径]",
	Short: "生成代码上下文 (可选命令)",
	Long:  "扫描指定路径并生成代码项目结构文档。现在可以直接运行程序而不需要此命令。",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runGenerate,
}

// configCmd 配置命令
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置管理",
	Long:  "管理应用程序配置",
}

// configShowCmd 显示配置
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "显示当前配置",
	Long:  "显示当前配置设置",
	RunE:  runConfigShow,
}

// configInitCmd 初始化配置 (已移除 - 不再自动创建配置文件)
// var configInitCmd = &cobra.Command{
// 	Use:   "init",
// 	Short: "初始化配置文件",
// 	Long:  "创建默认配置文件",
// 	RunE:  runConfigInit,
// }

// init 初始化函数
func init() {
	// 添加子命令
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(configCmd)

	// 配置命令子命令
	configCmd.AddCommand(configShowCmd)
	// configCmd.AddCommand(configInitCmd) // 已移除 - 不再提供配置文件初始化功能

	// 全局标志
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "配置文件路径")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "详细输出")

	// 根命令的生成标志（与generate命令相同）
	rootCmd.Flags().StringP("output", "o", "", "输出文件路径")
	rootCmd.Flags().StringP("format", "f", "json", "输出格式 (json, xml, toml, markdown)")
	rootCmd.Flags().StringSliceP("exclude", "e", []string{}, "排除的文件/目录模式")
	rootCmd.Flags().StringSliceP("include", "i", []string{}, "包含的文件/目录模式")
	// 注意：recursive 参数已被移除，使用 max-depth 控制递归行为
	// max-depth = 0: 只扫描当前目录
	// max-depth = 1: 递归1层
	// max-depth = -1 或很大值: 无限递归
	rootCmd.Flags().Bool("hidden", false, "包含隐藏文件")
	rootCmd.Flags().IntP("max-depth", "d", 0, "最大扫描深度 (0表示只扫描当前目录，1表示递归1层，-1表示无限制)")
	rootCmd.Flags().IntP("max-size", "s", 0, "最大文件大小 (字节, 0表示无限制)")
	rootCmd.Flags().BoolP("content", "C", true, "包含文件内容")
	rootCmd.Flags().BoolP("hash", "H", false, "包含文件哈希")
	rootCmd.Flags().Bool("exclude-binary", true, "排除二进制文件")
	rootCmd.Flags().String("encoding", "utf-8", "输出文件编码格式")
	rootCmd.Flags().StringSliceP("multiple-files", "m", []string{}, "多个文件路径（可多次使用）")
	rootCmd.Flags().StringP("pattern-file", "p", "", "从文件读取模式（支持.gitignore格式，兼容Windows/Linux路径分隔符）")

	// generate命令标志（保持向后兼容）
	generateCmd.Flags().StringP("output", "o", "", "输出文件路径")
	generateCmd.Flags().StringP("format", "f", "json", "输出格式 (json, xml, toml, markdown)")
	generateCmd.Flags().StringSliceP("exclude", "e", []string{}, "排除的文件/目录模式")
	generateCmd.Flags().StringSliceP("include", "i", []string{}, "包含的文件/目录模式")
	// 注意：recursive 参数已被移除，使用 max-depth 控制递归行为
	generateCmd.Flags().Bool("hidden", false, "包含隐藏文件")
	generateCmd.Flags().IntP("max-depth", "d", 0, "最大扫描深度 (0表示只扫描当前目录，1表示递归1层，-1表示无限制)")
	generateCmd.Flags().IntP("max-size", "s", 0, "最大文件大小 (字节, 0表示无限制)")
	generateCmd.Flags().BoolP("content", "C", true, "包含文件内容")
	generateCmd.Flags().BoolP("hash", "H", false, "包含文件哈希")
	generateCmd.Flags().Bool("exclude-binary", true, "排除二进制文件")
	generateCmd.Flags().String("encoding", "utf-8", "输出文件编码格式")
	generateCmd.Flags().StringSliceP("multiple-files", "m", []string{}, "多个文件路径（可多次使用）")
	generateCmd.Flags().StringP("pattern-file", "p", "", "从文件读取模式（支持.gitignore格式，兼容Windows/Linux路径分隔符）")

	// Git集成相关标志
	generateCmd.Flags().Bool("git-enabled", false, "启用Git集成功能")
	generateCmd.Flags().Bool("git-logs", false, "包含Git提交历史")
	generateCmd.Flags().Int("git-log-count", 50, "Git提交历史记录数量")
	generateCmd.Flags().Bool("git-diffs", false, "包含Git差异信息")
	generateCmd.Flags().String("git-diff-format", "unified", "Git差异格式 (unified, context)")
	generateCmd.Flags().Bool("git-stats", false, "包含Git统计信息")
	generateCmd.Flags().String("git-time-period", "1y", "Git统计时间周期 (1y, 6m, 3m, 1m, 1w)")
	generateCmd.Flags().StringSlice("git-authors", []string{}, "过滤特定作者（可多次使用）")
	generateCmd.Flags().StringSlice("git-paths", []string{}, "过滤特定路径（可多次使用）")
	generateCmd.Flags().String("git-since", "", "Git提交开始时间 (YYYY-MM-DD)")
	generateCmd.Flags().String("git-until", "", "Git提交结束时间 (YYYY-MM-DD)")

	// 元信息标志
	generateCmd.Flags().Bool("include-metadata", false, "包含元信息（如Git数据等）")
}

// main 主函数
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, utils.ErrorColor("错误:"), err)
		os.Exit(1)
	}
}

// runGenerate 运行生成命令
func runGenerate(cmd *cobra.Command, args []string) error {
	// 解析标志
	output, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	include, _ := cmd.Flags().GetStringSlice("include")
	// recursive 参数已被移除
	hidden, _ := cmd.Flags().GetBool("hidden")
	maxDepth, _ := cmd.Flags().GetInt("max-depth")
	maxSize, _ := cmd.Flags().GetInt("max-size")
	content, _ := cmd.Flags().GetBool("content")
	hash, _ := cmd.Flags().GetBool("hash")
	excludeBinary, _ := cmd.Flags().GetBool("exclude-binary")
	encoding, _ := cmd.Flags().GetString("encoding")
	multipleFiles, _ := cmd.Flags().GetStringSlice("multiple-files")
	patternFile, _ := cmd.Flags().GetString("pattern-file")

	// Git集成相关标志
	gitEnabled, _ := cmd.Flags().GetBool("git-enabled")
	gitLogs, _ := cmd.Flags().GetBool("git-logs")
	gitLogCount, _ := cmd.Flags().GetInt("git-log-count")
	gitDiffs, _ := cmd.Flags().GetBool("git-diffs")
	gitDiffFormat, _ := cmd.Flags().GetString("git-diff-format")
	gitStats, _ := cmd.Flags().GetBool("git-stats")
	gitTimePeriod, _ := cmd.Flags().GetString("git-time-period")
	gitAuthors, _ := cmd.Flags().GetStringSlice("git-authors")
	gitPaths, _ := cmd.Flags().GetStringSlice("git-paths")
	gitSince, _ := cmd.Flags().GetString("git-since")
	gitUntil, _ := cmd.Flags().GetString("git-until")

	// 元信息标志
	includeMetadata, _ := cmd.Flags().GetBool("include-metadata")

	// 如果指定了多个文件，使用第一个文件作为路径参数
	path := "."
	if len(multipleFiles) > 0 {
		path = multipleFiles[0] // 使用第一个文件作为基础路径
	} else if len(args) > 0 {
		path = args[0]
	}

	// 如果指定了模式文件，读取并解析模式
	if patternFile != "" {
		patterns, err := readPatternFile(patternFile)
		if err != nil {
			return fmt.Errorf("读取模式文件失败: %w", err)
		}
		if len(exclude) == 0 {
			exclude = patterns
		} else {
			exclude = append(exclude, patterns...)
		}
	}

	// 合并配置文件设置（命令行参数优先）
	if len(exclude) == 0 && len(cfg.Filters.ExcludePatterns) > 0 {
		exclude = cfg.Filters.ExcludePatterns
	}
	if len(include) == 0 && len(cfg.Filters.IncludePatterns) > 0 {
		include = cfg.Filters.IncludePatterns
	}
	// 修复：当命令行maxDepth为0时，使用配置中的值（包括0）
	if maxDepth == 0 {
		maxDepth = cfg.Filters.MaxDepth
	}
	if maxSize == 0 && cfg.Filters.MaxFileSize != "" {
		// 解析配置文件中的文件大小字符串
		parsedSize := env.ParseFileSize(cfg.Filters.MaxFileSize)
		if parsedSize > 0 {
			maxSize = int(parsedSize)
		}
	}
	if !hidden && cfg.FileProcessing.IncludeHidden {
		hidden = cfg.FileProcessing.IncludeHidden
	}
	if !excludeBinary && cfg.Filters.ExcludeBinary {
		excludeBinary = cfg.Filters.ExcludeBinary
	}

	// 应用编码设置（命令行参数优先）
	if encoding != "" && encoding != "utf-8" {
		cfg.Output.Encoding = encoding
	}

	// 合并Git配置（命令行参数优先）
	if gitEnabled {
		cfg.Git.Enabled = true
	}
	if gitLogs {
		cfg.Git.IncludeLogs = true
	}
	if gitLogCount > 0 && gitLogCount != 50 {
		cfg.Git.LogCount = gitLogCount
	}
	if gitDiffs {
		cfg.Git.IncludeDiffs = true
	}
	if gitDiffFormat != "" && gitDiffFormat != "unified" {
		cfg.Git.DiffFormat = gitDiffFormat
	}
	if gitStats {
		cfg.Git.Stats.Enabled = true
	}
	if gitTimePeriod != "" && gitTimePeriod != "1y" {
		cfg.Git.Stats.TimePeriod = gitTimePeriod
	}
	if len(gitAuthors) > 0 {
		cfg.Git.Filters.Authors = gitAuthors
	}
	if len(gitPaths) > 0 {
		cfg.Git.Filters.Paths = gitPaths
	}
	if gitSince != "" {
		cfg.Git.Filters.Since = gitSince
	}
	if gitUntil != "" {
		cfg.Git.Filters.Until = gitUntil
	}

	// 合并元信息配置
	if includeMetadata {
		cfg.Output.IncludeMetadata = true
	}

	// 验证格式
	if !isValidFormat(format) {
		return fmt.Errorf("无效的输出格式: %s", format)
	}

	// 创建文件系统遍历器
	walker := filesystem.NewFileSystemWalker(types.WalkOptions{})

	// 设置walker的配置
	if fsWalker, ok := walker.(*filesystem.FileSystemWalker); ok {
		fsWalker.SetConfig(cfg)
	}

	// 新的max-depth逻辑：
	// 0: 只扫描当前目录（不递归）
	// 1: 递归1层
	// -1 或很大值: 无限递归
	// 如果maxDepth为0，保持为0（只扫描当前目录）

	// 执行遍历
	if verbose {
		if len(multipleFiles) > 0 {
			fmt.Printf("正在处理指定文件: %v\n", multipleFiles)
		} else {
			fmt.Printf("正在扫描路径: %s (最大深度: %d)\n", path, maxDepth)
		}
		fmt.Printf("排除模式: %v\n", exclude)
		fmt.Printf("最大深度: %d, 最大文件大小: %d\n", maxDepth, maxSize)
	}

	// 创建遍历选项
	walkOptions := &types.WalkOptions{
		MaxDepth:        maxDepth,
		MaxFileSize:     int64(maxSize),
		ExcludePatterns: exclude,
		IncludePatterns: include,
		FollowSymlinks:  false,
		ShowHidden:      hidden,
		ExcludeBinary:   excludeBinary,
		MultipleFiles:   multipleFiles,
		PatternFile:     patternFile,
	}

	var result *types.ContextData
	var err error

	if len(multipleFiles) > 0 {
		// 处理多个指定文件
		result, err = walker.Walk(multipleFiles[0], walkOptions)
	} else {
		// 正常遍历目录
		result, err = walker.Walk(path, walkOptions)
	}

	if err != nil {
		return fmt.Errorf("扫描失败: %w", err)
	}

	if verbose {
		fmt.Printf("扫描完成: %d 个文件, %d 个目录\n", result.FileCount, result.FolderCount)
	}

	// 执行安全扫描
	if cfg.Security.Enabled {
		fmt.Println(utils.InfoColor("🔍 开始安全扫描..."))
		securityIntegration := security.NewSecurityIntegration(&cfg.Security)

		// 收集要扫描的文件路径
		var filesToScan []string
		for _, file := range result.Files {
			filesToScan = append(filesToScan, file.Path)
		}
		for _, folder := range result.Folders {
			for _, file := range folder.Files {
				filesToScan = append(filesToScan, file.Path)
			}
		}

		securityReport, err := securityIntegration.ScanFiles(filesToScan)
		if err != nil {
			fmt.Printf("安全扫描失败: %v\n", err)
		} else {
			securityIntegration.PrintSummary(securityReport)

			// 如果启用了失败选项且有关键问题，则退出
			if cfg.Security.FailOnCritical && securityIntegration.HasCriticalIssues(securityReport) {
				return fmt.Errorf("发现严重安全问题，扫描终止")
			}

			// 生成安全报告文件
			if cfg.Security.ReportFormat != "" {
				securityReportFile := fmt.Sprintf("security_report_%s.%s",
					filepath.Base(path), cfg.Security.ReportFormat)
				if cfg.Security.ReportFormat == "text" {
					securityReportFile = fmt.Sprintf("security_report_%s.txt", filepath.Base(path))
				}

				err = securityIntegration.GenerateReport(securityReport, securityReportFile)
				if err != nil {
					fmt.Printf("生成安全报告失败: %v\n", err)
				} else {
					fmt.Printf("安全报告已生成: %s\n", securityReportFile)
				}
			}
		}
	}

	// 执行Git集成
	if cfg.Git.Enabled {
		fmt.Println(utils.InfoColor("🔍 开始Git集成分析..."))
		gitIntegration, err := git.NewIntegration(path, &cfg.Git)
		if err != nil {
			fmt.Printf("Git集成初始化失败: %v\n", err)
			// Git集成失败不终止整个流程，只是警告
		} else {
			// 获取Git集成数据
			gitData, err := gitIntegration.GetGitIntegrationData()
			if err != nil {
				fmt.Printf("Git集成失败: %v\n", err)
				// Git集成失败不终止整个流程，只是警告
			} else if gitData != nil {
				// 将Git数据添加到结果中
				if result.Metadata == nil {
					result.Metadata = make(map[string]interface{})
				}
				result.Metadata["git"] = gitData
				
				if verbose {
					fmt.Printf("Git仓库: %s\n", gitData.GitInfo.RepositoryPath)
					if gitData.GitInfo.IsGitRepo {
						fmt.Printf("分支: %s\n", gitData.GitInfo.CurrentBranch)
						if cfg.Git.IncludeLogs && gitData.GitHistory != nil {
							fmt.Printf("提交数量: %d\n", len(gitData.GitHistory.Commits))
						}
						if cfg.Git.Stats.Enabled && gitData.GitStats != nil {
							fmt.Printf("统计信息已生成\n")
						}
					}
				}
			}
		}
	}

	// 创建格式化器
	formatter, err := formatter.NewFormatter(format, cfg)
	if err != nil {
		return fmt.Errorf("创建格式化器失败: %w", err)
	}

	// ContextData 已经包含了所有需要的信息
	// 初始化metadata map并添加根路径
	if result.Metadata == nil {
		result.Metadata = make(map[string]interface{})
	}
	result.Metadata["root_path"] = path
	contextData := *result

	// 格式化输出
	outputData, err := formatter.Format(contextData)
	if err != nil {
		return fmt.Errorf("格式化输出失败: %w", err)
	}

	// 添加额外信息
	if content || hash {
		// 创建 WalkResult 用于 addFileContent
		walkResult := &types.WalkResult{
			Files:       result.Files,
			Folders:     result.Folders,
			FileCount:   result.FileCount,
			FolderCount: result.FolderCount,
			TotalSize:   result.TotalSize,
			RootPath:    path,
		}
		outputData = addFileContent(outputData, walkResult, content, hash)
	}

	// 输出结果 - 默认写入文件，控制台输出仅在明确指定时
	if output != "" {
		// 使用指定的输出文件
		// 标准化换行符为当前操作系统格式
		normalizedData := utils.NormalizeLineEndings(outputData)
		if err := os.WriteFile(output, []byte(normalizedData), 0644); err != nil {
			return fmt.Errorf("写入输出文件失败: %w", err)
		}
		if verbose {
			fmt.Println(utils.SuccessColor("输出已写入:"), output)
		}
	} else {
		// 自动生成默认输出文件名
		var defaultOutput string
		if len(multipleFiles) > 0 {
			// 使用第一个文件名作为基础名称
			baseName := filepath.Base(multipleFiles[0])
			ext := filepath.Ext(baseName)
			baseName = strings.TrimSuffix(baseName, ext)
			defaultOutput = fmt.Sprintf("context_%s.%s", baseName, format)
			if format == "markdown" {
				defaultOutput = fmt.Sprintf("context_%s.md", baseName)
			}
		} else {
			defaultOutput = fmt.Sprintf("context_%s.%s", filepath.Base(path), format)
			if format == "markdown" {
				defaultOutput = fmt.Sprintf("context_%s.md", filepath.Base(path))
			}
		}

		// 标准化换行符为当前操作系统格式
		normalizedData := utils.NormalizeLineEndings(outputData)
		if err := os.WriteFile(defaultOutput, []byte(normalizedData), 0644); err != nil {
			return fmt.Errorf("写入默认输出文件失败: %w", err)
		}
		fmt.Println(utils.SuccessColor("✅ 成功生成代码上下文文件:"), defaultOutput)
		fmt.Printf("📊 包含 %d 个文件，%d 个文件夹\n", result.FileCount, result.FolderCount)
		fmt.Printf("💾 总大小: %s\n", utils.FormatFileSize(result.TotalSize))

		// 显示安全扫描状态
		if cfg.Security.Enabled {
			fmt.Println(utils.SuccessColor("🔒 安全扫描已启用"))
		} else {
			fmt.Println(utils.InfoColor("🔓 安全扫描已禁用"))
		}

		// 显示Git集成状态
		if cfg.Git.Enabled {
			fmt.Println(utils.SuccessColor("🔀 Git集成已启用"))
			if result.Metadata != nil {
				if gitData, ok := result.Metadata["git"].(*types.GitIntegrationData); ok && gitData.GitInfo != nil && gitData.GitInfo.IsGitRepo {
					fmt.Printf("📋 Git仓库: %s\n", gitData.GitInfo.CurrentBranch)
					if cfg.Git.IncludeLogs && gitData.GitHistory != nil {
						fmt.Printf("📝 提交历史: %d条记录\n", len(gitData.GitHistory.Commits))
					}
					if cfg.Git.Stats.Enabled && gitData.GitStats != nil {
						fmt.Printf("📊 Git统计信息已生成\n")
					}
				}
			}
		} else {
			fmt.Println(utils.InfoColor("🔀 Git集成已禁用"))
		}
	}

	return nil
}

// runConfigShow 运行配置显示命令
func runConfigShow(cmd *cobra.Command, args []string) error {
	// 生成配置输出
	configOutput := generateConfigOutput(cfg)
	fmt.Println(configOutput)
	return nil
}

// runConfigInit 运行配置初始化命令 (已移除 - 不再自动创建配置文件)
// func runConfigInit(cmd *cobra.Command, args []string) error {
// 	// 初始化配置
// 	configManager := config.NewManager()
// 	cfg = configManager.Get()

// 	// 保存配置到文件
// 	if err := configManager.Save("config.yaml", "yaml"); err != nil {
// 		return fmt.Errorf("保存配置文件失败: %w", err)
// 	}

// 	fmt.Println(utils.SuccessColor("配置文件已创建: config.yaml"))
// 	return nil
// }

// isValidFormat 检查格式是否有效
func isValidFormat(format string) bool {
	validFormats := []string{"json", "xml", "toml", "markdown", "md"}
	for _, valid := range validFormats {
		if format == valid {
			return true
		}
	}
	return false
}

// addFileContent 添加文件内容
func addFileContent(outputData string, _ *types.WalkResult, includeContent, includeHash bool) string {
	// 如果不需要包含内容和哈希，直接返回原始数据
	if !includeContent && !includeHash {
		return outputData
	}

	// 这里可以根据需要添加文件内容和哈希处理逻辑
	// 目前保持简化实现，后续可以根据具体需求扩展
	if verbose {
		fmt.Println(utils.InfoColor("注意: 文件内容和哈希功能暂未完全实现"))
	}

	return outputData
}

// readPatternFile 读取模式文件，支持.gitignore格式，兼容Windows/Linux路径分隔符
func readPatternFile(patternFile string) ([]string, error) {
	content, err := os.ReadFile(patternFile)
	if err != nil {
		return nil, fmt.Errorf("无法读取模式文件: %w", err)
	}

	var patterns []string
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 统一路径分隔符：将\和/都转换为当前系统的路径分隔符
		// 这样可以支持Windows和Linux格式的路径
		if filepath.Separator == '\\' {
			// Windows系统：将/转换为\，同时处理双反斜杠
			line = strings.ReplaceAll(line, "/", "\\")
			line = strings.ReplaceAll(line, "\\\\", "\\") // 处理双反斜杠
		} else {
			// Unix/Linux系统：将\转换为/
			line = strings.ReplaceAll(line, "\\", "/")
			line = strings.ReplaceAll(line, "//", "/") // 处理双斜杠
		}

		patterns = append(patterns, line)
	}

	return patterns, nil
}

// generateConfigOutput 生成配置输出
func generateConfigOutput(cfg *types.Config) string {
	var output strings.Builder

	output.WriteString("当前配置:\n")
	output.WriteString("==================\n\n")

	output.WriteString(fmt.Sprintf("默认格式: %s\n", cfg.Output.DefaultFormat))
	output.WriteString(fmt.Sprintf("输出目录: %s\n", cfg.Output.OutputDir))
	output.WriteString(fmt.Sprintf("文件名模板: %s\n", cfg.Output.FilenameTemplate))

	output.WriteString("\n文件处理:\n")
	output.WriteString(fmt.Sprintf("  最大文件大小: %s\n", cfg.Filters.MaxFileSize))
	output.WriteString(fmt.Sprintf("  最大深度: %d\n", cfg.Filters.MaxDepth))
	output.WriteString(fmt.Sprintf("  跟随符号链接: %v\n", cfg.Filters.FollowSymlinks))
	output.WriteString(fmt.Sprintf("  排除二进制文件: %v\n", cfg.Filters.ExcludeBinary))

	if len(cfg.Filters.ExcludePatterns) > 0 {
		output.WriteString("  排除模式:\n")
		for _, pattern := range cfg.Filters.ExcludePatterns {
			output.WriteString(fmt.Sprintf("    - %s\n", pattern))
		}
	}

	if len(cfg.Filters.IncludePatterns) > 0 {
		output.WriteString("  包含模式:\n")
		for _, pattern := range cfg.Filters.IncludePatterns {
			output.WriteString(fmt.Sprintf("    - %s\n", pattern))
		}
	}

	output.WriteString("\nGit集成:\n")
	output.WriteString(fmt.Sprintf("  启用状态: %v\n", cfg.Git.Enabled))
	if cfg.Git.Enabled {
		output.WriteString(fmt.Sprintf("  包含提交历史: %v\n", cfg.Git.IncludeLogs))
		if cfg.Git.IncludeLogs {
			output.WriteString(fmt.Sprintf("  提交历史数量: %d\n", cfg.Git.LogCount))
		}
		output.WriteString(fmt.Sprintf("  包含差异信息: %v\n", cfg.Git.IncludeDiffs))
		if cfg.Git.IncludeDiffs {
			output.WriteString(fmt.Sprintf("  差异格式: %s\n", cfg.Git.DiffFormat))
		}
		output.WriteString(fmt.Sprintf("  包含统计信息: %v\n", cfg.Git.Stats.Enabled))
		if cfg.Git.Stats.Enabled {
			output.WriteString(fmt.Sprintf("  统计时间周期: %s\n", cfg.Git.Stats.TimePeriod))
			output.WriteString(fmt.Sprintf("  作者排行数量: %d\n", cfg.Git.Stats.AuthorsTop))
			output.WriteString(fmt.Sprintf("  文件排行数量: %d\n", cfg.Git.Stats.FilesTop))
		}
		if len(cfg.Git.Filters.Authors) > 0 {
			output.WriteString("  作者过滤:\n")
			for _, author := range cfg.Git.Filters.Authors {
				output.WriteString(fmt.Sprintf("    - %s\n", author))
			}
		}
		if len(cfg.Git.Filters.Paths) > 0 {
			output.WriteString("  路径过滤:\n")
			for _, path := range cfg.Git.Filters.Paths {
				output.WriteString(fmt.Sprintf("    - %s\n", path))
			}
		}
		if cfg.Git.Filters.Since != "" {
			output.WriteString(fmt.Sprintf("  开始时间: %s\n", cfg.Git.Filters.Since))
		}
		if cfg.Git.Filters.Until != "" {
			output.WriteString(fmt.Sprintf("  结束时间: %s\n", cfg.Git.Filters.Until))
		}
	}

	return output.String()
}
