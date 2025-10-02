// Package main CLI应用程序主入口
package main

import (
	"fmt"
	"os"
	"strings"

	"code-context-generator/internal/config"
	"code-context-generator/internal/env"
	"code-context-generator/internal/filesystem"
	"code-context-generator/internal/formatter"
	"code-context-generator/internal/utils"
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
	Use:   "code-context-generator",
	Short: "代码上下文生成器",
	Long: `代码上下文生成器 - 智能生成代码项目结构文档

支持多种输出格式（JSON、XML、TOML、Markdown），提供自动文件扫描，
自动补全功能，以及丰富的配置选项。`,
	Version: version,
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
}

// generateCmd 生成命令
var generateCmd = &cobra.Command{
	Use:   "generate [路径]",
	Short: "生成代码上下文",
	Long:  "扫描指定路径并生成代码项目结构文档",
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

	// generate命令标志
	generateCmd.Flags().StringP("output", "o", "", "输出文件路径")
	generateCmd.Flags().StringP("format", "f", "json", "输出格式 (json, xml, toml, markdown)")
	generateCmd.Flags().StringSliceP("exclude", "e", []string{}, "排除的文件/目录模式")
	generateCmd.Flags().StringSliceP("include", "i", []string{}, "包含的文件/目录模式")
	generateCmd.Flags().BoolP("recursive", "r", true, "递归扫描")
	generateCmd.Flags().Bool("hidden", false, "包含隐藏文件")
	generateCmd.Flags().IntP("max-depth", "d", 0, "最大扫描深度 (0表示无限制)")
	generateCmd.Flags().IntP("max-size", "s", 0, "最大文件大小 (字节, 0表示无限制)")
	generateCmd.Flags().BoolP("content", "C", true, "包含文件内容")
	generateCmd.Flags().BoolP("hash", "H", false, "包含文件哈希")
	generateCmd.Flags().Bool("exclude-binary", true, "排除二进制文件")
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
	// 获取路径
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// 解析标志
	output, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	include, _ := cmd.Flags().GetStringSlice("include")
	recursive, _ := cmd.Flags().GetBool("recursive")
	hidden, _ := cmd.Flags().GetBool("hidden")
	maxDepth, _ := cmd.Flags().GetInt("max-depth")
	maxSize, _ := cmd.Flags().GetInt("max-size")
	content, _ := cmd.Flags().GetBool("content")
	hash, _ := cmd.Flags().GetBool("hash")
	excludeBinary, _ := cmd.Flags().GetBool("exclude-binary")

	// 合并配置文件设置（命令行参数优先）
	if len(exclude) == 0 && len(cfg.Filters.ExcludePatterns) > 0 {
		exclude = cfg.Filters.ExcludePatterns
	}
	if len(include) == 0 && len(cfg.Filters.IncludePatterns) > 0 {
		include = cfg.Filters.IncludePatterns
	}
	if maxDepth == 0 && cfg.Filters.MaxDepth > 0 {
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

	// 验证格式
	if !isValidFormat(format) {
		return fmt.Errorf("无效的输出格式: %s", format)
	}

	// 创建文件系统遍历器
	walker := filesystem.NewFileSystemWalker(types.WalkOptions{})

	// 如果递归选项被禁用，设置最大深度为1
	if !recursive && maxDepth == 0 {
		maxDepth = 1
	}

	// 执行遍历
	if verbose {
		fmt.Printf("正在扫描路径: %s (递归: %v)\n", path, recursive)
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
	}

	result, err := walker.Walk(path, walkOptions)
	if err != nil {
		return fmt.Errorf("扫描失败: %w", err)
	}

	if verbose {
		fmt.Printf("扫描完成: %d 个文件, %d 个目录\n", result.FileCount, result.FolderCount)
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

	// 输出结果
	if output != "" {
		// 标准化换行符为当前操作系统格式
		normalizedData := utils.NormalizeLineEndings(outputData)
		if err := os.WriteFile(output, []byte(normalizedData), 0644); err != nil {
			return fmt.Errorf("写入输出文件失败: %w", err)
		}
		if verbose {
			fmt.Println(utils.SuccessColor("输出已写入:"), output)
		}
	} else {
		// 控制台输出也标准化换行符
		fmt.Println(utils.NormalizeLineEndings(outputData))
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

	return output.String()
}
