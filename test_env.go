package main

import (
	"fmt"
	"log"

	"code-context-generator/internal/config"
	"code-context-generator/internal/env"
)

func main() {
	fmt.Println("=== .env配置测试 ===")
	
	// 首先加载.env文件
	if err := env.LoadEnv(""); err != nil {
		log.Printf("加载.env文件失败: %v", err)
	} else {
		fmt.Println("✓ .env文件加载成功")
	}
	
	// 显示环境变量
	envVars := env.GetAllEnvVars()
	fmt.Println("\n=== 环境变量配置 ===")
	for key, value := range envVars {
		fmt.Printf("%s: %s\n", key, value)
	}
	
	// 创建配置管理器并加载配置
	cm := config.NewManager()
	if err := cm.Load(""); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	
	cfg := cm.Get()
	fmt.Println("\n=== 最终配置 ===")
	fmt.Printf("默认格式: %s\n", cfg.Output.DefaultFormat)
	fmt.Printf("输出目录: %s\n", cfg.Output.OutputDir)
	fmt.Printf("最大深度: %d\n", cfg.Filters.MaxDepth)
	fmt.Printf("跟随符号链接: %v\n", cfg.Filters.FollowSymlinks)
	fmt.Printf("排除二进制文件: %v\n", cfg.Filters.ExcludeBinary)
	fmt.Printf("自动补全: %v\n", cfg.UI.Autocomplete.Enabled)
	
	fmt.Println("\n=== 测试完成 ===")
}