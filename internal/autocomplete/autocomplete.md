// Package autocomplete 提供自动补全功能
package autocomplete

// 此文件现在只包含包级别的文档说明
// 具体实现已拆分到以下文件：
// - interfaces.go: 接口定义
// - filepath_autocompleter.go: 文件路径自动补全器
// - command_autocompleter.go: 命令自动补全器
// - suggestion_provider.go: 建议提供者
// - utils.go: 工具函数

// 使用示例：
// import "code-context-generator/internal/autocomplete"
//
// // 创建文件路径自动补全器
// autocompleter := autocomplete.NewAutocompleter(config)
//
// // 创建命令自动补全器
// cmdAutocompleter := autocomplete.NewCommandAutocompleter()
//
// // 创建组合建议提供者
// provider := autocomplete.NewCompositeSuggestionProvider(provider1, provider2)
