package main

import (
	"fmt"
	"log"

	"code-context-generator/internal/config"
	"code-context-generator/internal/env"
)

func main() {
	// 首先加载.env文件（如果存在）
	if err := env.LoadEnv(""); err != nil {
		log.Printf("警告: 加载.env文件失败: %v", err)
	}

	// 创建配置管理器
	cm := config.NewManager()

	// 获取默认配置
	fmt.Printf("默认配置: %+v\n", cm.Get())

	// 保存配置为YAML格式
	if err := cm.Save("config.yaml", "yaml"); err != nil {
		log.Fatalf("保存配置失败: %v", err)
	}
	fmt.Println("配置已保存为YAML格式")
}
