// Package filesystem 提供文件系统遍历和过滤功能
package filesystem

// 文件系统功能已拆分到以下文件：
// - walker.go: 文件系统遍历器接口和主要遍历逻辑
// - fileinfo.go: 文件信息获取功能
// - filters.go: 文件过滤功能
// - utils.go: 通用文件系统工具函数
//
// 使用示例：
//   walker := filesystem.NewWalker()
//   contextData, err := walker.Walk(rootPath, options)
//
// 原有的所有功能仍然可用，只是按职责分离到了不同的文件中
