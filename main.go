package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"
)

func main() {
	// 创建示例配置
	config := createSampleConfig()
	
	// 保存配置到文件
	fmt.Println("保存配置到文件...")
	
	// 保存为YAML
	yamlData, _ := yaml.Marshal(config)
	os.WriteFile("config.yaml", yamlData, 0644)
	
	// 保存为JSON
	jsonData, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile("config.json", jsonData, 0644)
	
	// 保存为TOML
	var buf bytes.Buffer
	tomlEncoder := toml.NewEncoder(&buf)
	tomlEncoder.Encode(config)
	os.WriteFile("config.toml", buf.Bytes(), 0644)
	
	fmt.Println("配置文件已保存")
	
	// 创建配置管理器
	cm, err := NewConfigManager("config.yaml")
	if err != nil {
		log.Fatalf("创建配置管理器失败: %v", err)
	}
	
	// 演示配置加载
	fmt.Println("\n=== 加载YAML配置 ===")
	yamlConfig, err := LoadConfig("config.yaml")
	if err != nil {
		log.Printf("加载YAML配置失败: %v", err)
	} else {
		fmt.Printf("默认格式: %s\n", yamlConfig.Output.DefaultFormat)
	}
	
	fmt.Println("\n=== 加载JSON配置 ===")
	jsonConfig, err := LoadConfig("config.json")
	if err != nil {
		log.Printf("加载JSON配置失败: %v", err)
	} else {
		fmt.Printf("默认格式: %s\n", jsonConfig.Output.DefaultFormat)
	}
	
	fmt.Println("\n=== 加载TOML配置 ===")
	tomlConfig, err := LoadConfig("config.toml")
	if err != nil {
		log.Printf("加载TOML配置失败: %v", err)
	} else {
		fmt.Printf("默认格式: %s\n", tomlConfig.Output.DefaultFormat)
	}
	
	// 演示配置管理器功能
	fmt.Println("\n=== 配置管理器功能演示 ===")
	currentConfig := cm.GetConfig()
	fmt.Printf("当前默认格式: %s\n", currentConfig.Output.DefaultFormat)
	
	// 演示输出生成
	fmt.Println("\n=== 生成XML输出 ===")
	xmlOutput, err := cm.GenerateXMLOutput(sampleFiles, sampleFolders)
	if err != nil {
		log.Printf("XML输出生成失败: %v", err)
	} else {
		fmt.Println(xmlOutput)
	}
	
	fmt.Println("\n=== 生成JSON输出 ===")
	jsonOutput, err := cm.GenerateJSONOutput(sampleFiles, sampleFolders)
	if err != nil {
		log.Printf("JSON输出生成失败: %v", err)
	} else {
		fmt.Println(jsonOutput)
	}
	
	fmt.Println("\n=== 生成Markdown输出 ===")
	mdOutput, err := cm.GenerateMarkdownOutput(sampleFiles, sampleFolders)
	if err != nil {
		log.Printf("Markdown输出生成失败: %v", err)
	} else {
		fmt.Println(mdOutput)
	}
}