package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"code-context-generator/internal/config"
	"code-context-generator/internal/env"
	"code-context-generator/internal/filesystem"
	"code-context-generator/internal/formatter"
	"code-context-generator/internal/utils"
	"code-context-generator/pkg/types"
)

func main() {
	// 手动解析命令行参数
	var format, outputFile string
	var showHelp bool

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-f", "--format":
			if i+1 < len(args) {
				format = args[i+1]
				i++
			}
		case "-o", "--output":
			if i+1 < len(args) {
				outputFile = args[i+1]
				i++
			}
		case "-h", "--help":
			showHelp = true
		case "generate":
			// 忽略generate参数，兼容用户的命令格式
		}
	}

	// 显示帮助信息
	if showHelp {
		fmt.Println("=== 代码上下文生成器 ===")
		fmt.Println("使用方式: go run main.go [generate] [选项]")
		fmt.Println()
		fmt.Println("选项:")
		fmt.Println("  -f string    输出格式 (json, xml, markdown)")
		fmt.Println("  -o string    输出文件路径")
		fmt.Println("  -h           显示帮助信息")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  go run main.go -f markdown -o output.md")
		fmt.Println("  go run main.go generate -f markdown -o output.md")
		fmt.Println("  go run main.go -f json")
		fmt.Println()
		fmt.Println("如果不指定格式，将使用默认格式 (markdown)")
		return
	}

	// 首先加载.env文件（如果存在）
	if err := env.LoadEnv(""); err != nil {
		log.Printf("警告: 加载.env文件失败: %v", err)
	}

	// 创建配置管理器
	cm := config.NewManager()

	// 获取当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前目录失败: %v", err)
	}

	fmt.Println("=== 代码上下文生成器 ===")
	fmt.Printf("当前目录: %s\n", currentDir)
	fmt.Println()

	// 创建文件系统遍历器
	walker := filesystem.NewWalker()
	
	// 设置walker的配置
	if fsWalker, ok := walker.(*filesystem.FileSystemWalker); ok {
		fsWalker.SetConfig(cm.Get())
	}

	// 设置遍历选项
	walkOptions := &types.WalkOptions{
		MaxDepth:        3,       // 限制子目录深度
		MaxFileSize:     1048576, // 1MB
		ExcludePatterns: []string{".git", "node_modules", "*.exe", "*.dll"},
		ExcludeBinary:   false,
		ShowHidden:      false,
	}

	// 遍历当前目录
	contextData, err := walker.Walk(currentDir, walkOptions)
	if err != nil {
		log.Fatalf("遍历目录失败: %v", err)
	}

	if len(contextData.Files) == 0 {
		fmt.Println("未找到任何文件")
		return
	}

	fmt.Printf("找到 %d 个文件, %d 个目录\n", contextData.FileCount, contextData.FolderCount)

	// 转换为上下文数据
	contextData.Metadata = map[string]interface{}{
		"root_path":    currentDir,
		"generated_at": "现在",
	}

	// 确定输出格式
	var selectedFormat string
	if format != "" {
		// 使用命令行指定的格式
		selectedFormat = format
		// 验证格式是否有效
		validFormats := map[string]bool{
			"json": true, "xml": true, "markdown": true, "md": true,
		}
		if !validFormats[selectedFormat] {
			log.Fatalf("无效的输出格式: %s，支持的格式: json, xml, markdown", selectedFormat)
		}
		if selectedFormat == "md" {
			selectedFormat = "markdown" // 统一处理
		}
		fmt.Printf("使用指定的输出格式: %s\n", selectedFormat)
	} else {
		// 默认格式
		fmt.Println("\n选择输出格式:")
		fmt.Println("1. JSON")
		fmt.Println("2. XML")
		fmt.Println("3. Markdown")
		fmt.Print("请选择 (1-3): ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			selectedFormat = "json"
		case 2:
			selectedFormat = "xml"
		case 3:
			selectedFormat = "markdown"
		default:
			selectedFormat = "json"
		}
	}

	// 创建格式化器
	formatter, err := formatter.NewFormatter(selectedFormat, cm.Get())
	if err != nil {
		log.Fatalf("创建格式化器失败: %v", err)
	}

	// 格式化输出
	outputData, err := formatter.Format(*contextData)
	if err != nil {
		log.Fatalf("格式化输出失败: %v", err)
	}

	// 确定输出文件路径
	var finalOutputFile string
	if outputFile != "" {
		finalOutputFile = outputFile
	} else {
		// 生成默认输出文件名
		finalOutputFile = fmt.Sprintf("context_%s.%s",
			filepath.Base(currentDir), selectedFormat)
	}

	// 保存到文件
	if err := os.WriteFile(finalOutputFile, []byte(outputData), 0644); err != nil {
		log.Fatalf("写入输出文件失败: %v", err)
	}

	fmt.Printf("\n✅ 成功生成代码上下文文件: %s\n", finalOutputFile)
	fmt.Printf("📊 包含 %d 个文件，%d 个文件夹\n", contextData.FileCount, contextData.FolderCount)
	fmt.Printf("💾 总大小: %s\n", utils.FormatFileSize(contextData.TotalSize))
}
