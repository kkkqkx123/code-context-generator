package main

import (
	"fmt"
	"os"
)

func main() {
	// 创建配置管理器
	cm := NewConfigManager()

	// 加载配置文件
	if err := cm.LoadConfig("config.yaml"); err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		// 使用默认配置继续
	}

	// 示例数据
	data := ContextData{
		Files: []FileInfo{
			{
				Name:    "main.go",
				Path:    "src/main.go",
				Content: "package main\n\nfunc main() {\n    fmt.Println(\"Hello World\")\n}",
				Size:    50,
			},
		},
		Folders: []FolderInfo{
			{
				Name: "src",
				Path: "src",
				Files: []FileInfo{
					{
						Name:    "utils.go",
						Path:    "src/utils.go",
						Content: "package src\n\nfunc Helper() {\n}",
						Size:    30,
					},
				},
			},
		},
	}

	// 生成不同格式的输出
	formats := []string{"xml", "json", "toml", "markdown"}
	for _, format := range formats {
		output, err := cm.GenerateOutput(data, format)
		if err != nil {
			fmt.Printf("生成%s格式失败: %v\n", format, err)
			continue
		}

		filename := cm.GetOutputFilename(format)
		if err := os.WriteFile(filename, []byte(output), 0644); err != nil {
			fmt.Printf("保存%s文件失败: %v\n", format, err)
			continue
		}

		fmt.Printf("已生成%s格式文件: %s\n", format, filename)
	}
}