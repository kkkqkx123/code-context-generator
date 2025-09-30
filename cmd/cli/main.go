// Package main CLI应用程序主入口
package main

import (
	"fmt"
	"os"
	"strings"

	"code-context-generator/internal/autocomplete"
	"code-context-generator/internal/filesystem"
	"code-context-generator/internal/formatter"
	"code-context-generator/internal/selector"
	"code-context-generator/internal/utils"
	"code-context-generator/pkg/types"

	"github.com/goccy/go-yaml"

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

支持多种输出格式（JSON、XML、TOML、Markdown），提供交互式文件选择，
自动补全功能，以及丰富的配置选项。`,
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 加载配置
		// 加载配置
		cfg = &types.Config{
			Output: types.OutputConfig{
				Format:   "json",
				Encoding: "utf-8",
			},
			FileProcessing: types.FileProcessingConfig{
				IncludeHidden:   false,
				MaxFileSize:     10 * 1024 * 1024,
				MaxDepth:        0,
				ExcludePatterns: []string{},
				IncludePatterns: []string{},
				IncludeContent:  false,
				IncludeHash:     false,
			},
			UI: types.UIConfig{
				Theme:        "default",
				ShowProgress: true,
				ShowSize:     true,
				ShowDate:     true,
				ShowPreview:  true,
			},
			Performance: types.PerformanceConfig{
				MaxWorkers:   4,
				BufferSize:   1024,
				CacheEnabled: true,
				CacheSize:    100,
			},
			Logging: types.LoggingConfig{
				Level:      "info",
				FilePath:   "",
				MaxSize:    10,
				MaxBackups: 3,
				MaxAge:     7,
			},
		}
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

// selectCmd 选择命令
var selectCmd = &cobra.Command{
	Use:   "select [路径]",
	Short: "交互式选择文件",
	Long:  "使用交互式界面选择要包含的文件和文件夹",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSelect,
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

// configInitCmd 初始化配置
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化配置文件",
	Long:  "创建默认配置文件",
	RunE:  runConfigInit,
}

// autocompleteCmd 自动补全命令
var autocompleteCmd = &cobra.Command{
	Use:   "autocomplete [路径]",
	Short: "文件路径自动补全",
	Long:  "提供文件路径自动补全建议",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runAutocomplete,
}

// init 初始化函数
func init() {
	// 添加子命令
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(selectCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(autocompleteCmd)

	// 配置命令子命令
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configInitCmd)

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
	generateCmd.Flags().BoolP("content", "C", false, "包含文件内容")
	generateCmd.Flags().BoolP("hash", "H", false, "包含文件哈希")

	// select命令标志
	selectCmd.Flags().StringP("output", "o", "", "输出文件路径")
	selectCmd.Flags().StringP("format", "f", "json", "输出格式")
	selectCmd.Flags().BoolP("multi", "m", true, "允许多选")
	selectCmd.Flags().StringP("filter", "F", "", "文件过滤器")

	// autocomplete命令标志
	autocompleteCmd.Flags().IntP("limit", "l", 10, "最大建议数量")
	autocompleteCmd.Flags().StringP("type", "t", "file", "补全类型 (file, dir, ext, pattern)")
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

	// 验证格式
	if !isValidFormat(format) {
		return fmt.Errorf("无效的输出格式: %s", format)
	}

	// 创建文件系统遍历器
	walker := filesystem.NewFileSystemWalker(types.WalkOptions{
		MaxDepth:        maxDepth,
		MaxFileSize:     int64(maxSize),
		ExcludePatterns: exclude,
		IncludePatterns: include,
		FollowSymlinks:  false,
		ShowHidden:      hidden,
	})

	// 如果递归选项被禁用，设置最大深度为1
	if !recursive && maxDepth == 0 {
		maxDepth = 1
	}

	// 执行遍历
	if verbose {
		fmt.Printf("正在扫描路径: %s (递归: %v)\n", path, recursive)
	}

	result, err := walker.Walk(path, nil)
	if err != nil {
		return fmt.Errorf("扫描失败: %w", err)
	}

	if verbose {
		fmt.Printf("扫描完成: %d 个文件, %d 个目录\n", result.FileCount, result.FolderCount)
	}

	// 创建格式化器
	formatter, err := formatter.NewFormatter(format)
	if err != nil {
		return fmt.Errorf("创建格式化器失败: %w", err)
	}

	// ContextData 已经包含了所有需要的信息
	// 添加根路径到metadata
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
		if err := os.WriteFile(output, []byte(outputData), 0644); err != nil {
			return fmt.Errorf("写入输出文件失败: %w", err)
		}
		if verbose {
			fmt.Println(utils.SuccessColor("输出已写入:"), output)
		}
	} else {
		fmt.Println(outputData)
	}

	return nil
}

// runSelect 运行选择命令
func runSelect(cmd *cobra.Command, args []string) error {
	// 获取路径
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// 解析标志
	output, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")
	multi, _ := cmd.Flags().GetBool("multi")
	filter, _ := cmd.Flags().GetString("filter")

	// 创建选择器配置
	config := &types.Config{
		FileProcessing: types.FileProcessingConfig{
			IncludeHidden: false,
		},
	}
	fileSelector := selector.NewFileSelector(config)

	// 执行选择
	if verbose {
		fmt.Printf("启动交互式选择器... (多选: %v, 过滤器: %s)\n", multi, filter)
	}

	// 选择文件和目录
	selectOptions := &types.SelectOptions{
		Recursive:       true,
		ShowHidden:      false,
		MaxDepth:        0,
		IncludePatterns: []string{},
		ExcludePatterns: []string{},
	}

	files, err := fileSelector.SelectFiles(path, selectOptions)
	if err != nil {
		return fmt.Errorf("选择文件失败: %w", err)
	}

	folders, err := fileSelector.SelectFolders(path, selectOptions)
	if err != nil {
		return fmt.Errorf("选择目录失败: %w", err)
	}

	// 合并选择结果
	allItems := append(files, folders...)

	// 交互式选择
	selected, err := fileSelector.InteractiveSelect(allItems, "选择文件和目录:")
	if err != nil {
		return fmt.Errorf("选择失败: %w", err)
	}

	if len(selected) == 0 {
		fmt.Println("未选择任何文件")
		return nil
	}

	if verbose {
		fmt.Printf("已选择 %d 个项目\n", len(selected))
	}

	// 创建结果
	result := &types.WalkResult{
		Files:    []types.FileInfo{},
		Folders:  []types.FolderInfo{},
		RootPath: path,
	}

	// 添加选择的文件和目录
	for _, item := range selected {
		info, err := os.Stat(item)
		if err != nil {
			continue
		}

		if info.IsDir() {
			result.Folders = append(result.Folders, types.FolderInfo{
				Path:  item,
				Name:  info.Name(),
				Size:  0,
				Count: 0,
			})
		} else {
			result.Files = append(result.Files, types.FileInfo{
				Path:     item,
				Name:     info.Name(),
				Size:     info.Size(),
				ModTime:  info.ModTime(),
				IsBinary: utils.IsBinaryFile(item),
			})
		}
	}

	// 更新统计信息
	result.FileCount = len(result.Files)
	result.FolderCount = len(result.Folders)

	// 创建格式化器
	formatter, err := formatter.NewFormatter(format)
	if err != nil {
		return fmt.Errorf("创建格式化器失败: %w", err)
	}

	// 将 WalkResult 转换为 ContextData
	contextData := types.ContextData{
		Files:       result.Files,
		Folders:     result.Folders,
		FileCount:   result.FileCount,
		FolderCount: result.FolderCount,
		TotalSize:   result.TotalSize,
		Metadata:    make(map[string]interface{}),
	}

	// 格式化输出
	outputData, err := formatter.Format(contextData)
	if err != nil {
		return fmt.Errorf("格式化输出失败: %w", err)
	}

	// 输出结果
	if output != "" {
		if err := os.WriteFile(output, []byte(outputData), 0644); err != nil {
			return fmt.Errorf("写入输出文件失败: %w", err)
		}
		if verbose {
			fmt.Println(utils.SuccessColor("输出已写入:"), output)
		}
	} else {
		fmt.Println(outputData)
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

// runConfigInit 运行配置初始化命令
func runConfigInit(cmd *cobra.Command, args []string) error {
	// 创建默认配置
	defaultConfig := &types.Config{
		Output: types.OutputConfig{
			Format:   "json",
			FilePath: "",
			Encoding: "utf-8",
		},
		FileProcessing: types.FileProcessingConfig{
			IncludeHidden: false,
			MaxFileSize:   10 * 1024 * 1024, // 10MB
			MaxDepth:      0,
			ExcludePatterns: []string{
				"*.exe", "*.dll", "*.so", "*.dylib",
				"*.pyc", "*.pyo", "*.pyd",
				"*.class", "*.jar",
				"*.o", "*.a", "*.lib",
				"*.iso", "*.img", "*.dmg",
				"node_modules", ".git", ".svn", ".hg",
				"__pycache__", "*.egg-info", "dist", "build",
			},
			IncludePatterns: []string{},
			IncludeContent:  false,
			IncludeHash:     false,
		},
		UI: types.UIConfig{
			Theme:        "default",
			ShowProgress: true,
			ShowSize:     true,
			ShowDate:     true,
			ShowPreview:  true,
		},
		Performance: types.PerformanceConfig{
			MaxWorkers:   4,
			BufferSize:   1024,
			CacheEnabled: true,
			CacheSize:    100,
		},
		Logging: types.LoggingConfig{
			Level:      "info",
			FilePath:   "",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
		},
	}

	// 保存配置
	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	if err := os.WriteFile("config.toml", data, 0644); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	fmt.Println(utils.SuccessColor("配置文件已创建: config.toml"))
	return nil
}

// runAutocomplete 运行自动补全命令
func runAutocomplete(cmd *cobra.Command, args []string) error {
	// 获取路径
	path := ""
	if len(args) > 0 {
		path = args[0]
	}

	// 解析标志
	limit, _ := cmd.Flags().GetInt("limit")
	completeType, _ := cmd.Flags().GetString("type")

	// 创建自动补全器
	autocompleter := autocomplete.NewAutocompleter(&types.AutocompleteConfig{
		Enabled:        true,
		MinChars:       1,
		MaxSuggestions: limit,
	})

	// 获取建议
	completeTypeEnum := types.CompleteFilePath
	switch completeType {
	case "dir":
		completeTypeEnum = types.CompleteDirectory
	case "ext":
		completeTypeEnum = types.CompleteExtension
	case "pattern":
		completeTypeEnum = types.CompletePattern
	case "generic":
		completeTypeEnum = types.CompleteGeneric
	}

	context := &types.CompleteContext{
		Type: completeTypeEnum,
		Data: make(map[string]interface{}),
	}
	suggestions, err := autocompleter.Complete(path, context)
	if err != nil {
		return fmt.Errorf("自动补全失败: %w", err)
	}

	// 输出建议
	for _, suggestion := range suggestions {
		fmt.Println(suggestion)
	}

	return nil
}

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

	output.WriteString(fmt.Sprintf("输出格式: %s\n", cfg.Output.Format))
	output.WriteString(fmt.Sprintf("编码: %s\n", cfg.Output.Encoding))
	if cfg.Output.FilePath != "" {
		output.WriteString(fmt.Sprintf("输出文件: %s\n", cfg.Output.FilePath))
	}

	output.WriteString("\n文件处理:\n")
	output.WriteString(fmt.Sprintf("  包含隐藏文件: %v\n", cfg.FileProcessing.IncludeHidden))
	output.WriteString(fmt.Sprintf("  最大文件大小: %d 字节\n", cfg.FileProcessing.MaxFileSize))
	output.WriteString(fmt.Sprintf("  最大深度: %d\n", cfg.FileProcessing.MaxDepth))
	output.WriteString(fmt.Sprintf("  包含内容: %v\n", cfg.FileProcessing.IncludeContent))
	output.WriteString(fmt.Sprintf("  包含哈希: %v\n", cfg.FileProcessing.IncludeHash))

	if len(cfg.FileProcessing.ExcludePatterns) > 0 {
		output.WriteString("  排除模式:\n")
		for _, pattern := range cfg.FileProcessing.ExcludePatterns {
			output.WriteString(fmt.Sprintf("    - %s\n", pattern))
		}
	}

	if len(cfg.FileProcessing.IncludePatterns) > 0 {
		output.WriteString("  包含模式:\n")
		for _, pattern := range cfg.FileProcessing.IncludePatterns {
			output.WriteString(fmt.Sprintf("    - %s\n", pattern))
		}
	}

	output.WriteString("\n性能:\n")
	output.WriteString(fmt.Sprintf("  最大工作线程: %d\n", cfg.Performance.MaxWorkers))
	output.WriteString(fmt.Sprintf("  缓冲区大小: %d\n", cfg.Performance.BufferSize))
	output.WriteString(fmt.Sprintf("  缓存启用: %v\n", cfg.Performance.CacheEnabled))
	output.WriteString(fmt.Sprintf("  缓存大小: %d\n", cfg.Performance.CacheSize))

	return output.String()
}
