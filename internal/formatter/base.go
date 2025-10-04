// Package formatter 提供多种格式的输出转换功能
package formatter

import (
	"encoding/xml"
	"strings"

	"code-context-generator/pkg/types"
)

// Formatter 格式转换器接口
type Formatter interface {
	Format(data types.ContextData) (string, error)
	FormatFile(file types.FileInfo) (string, error)
	FormatFolder(folder types.FolderInfo) (string, error)
	GetName() string
	GetDescription() string
}

// BaseFormatter 基础格式转换器
type BaseFormatter struct {
	name        string
	description string
	config      interface{}
}

// GetName 获取格式名称
func (f *BaseFormatter) GetName() string {
	return f.name
}

// GetDescription 获取格式描述
func (f *BaseFormatter) GetDescription() string {
	return f.description
}

// CDATAText 包装CDATA文本的类型
type CDATAText struct {
	Text string `xml:",cdata"`
}

// RawText 包装原始文本的类型（最小转义）
type RawText struct {
	Text string `xml:",innerxml"`
}

// MarshalXML 自定义CDATA文本的XML序列化
func (c CDATAText) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// 直接输出CDATA包装的内容
	return e.EncodeElement(struct {
		Text string `xml:",innerxml"`
	}{Text: "<![CDATA[" + c.Text + "]]>"}, start)
}

// MarshalXML 自定义原始文本的XML序列化
func (r RawText) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// 使用innerxml来避免转义，但需要确保内容是有效的XML
	// 这里我们只转义最基本的XML字符
	safeContent := r.Text
	safeContent = strings.ReplaceAll(safeContent, "&", "&amp;")
	safeContent = strings.ReplaceAll(safeContent, "<", "&lt;")
	safeContent = strings.ReplaceAll(safeContent, ">", "&gt;")

	return e.EncodeElement(struct {
		Text string `xml:",innerxml"`
	}{Text: safeContent}, start)
}

// SimplifiedFileInfo 简化的文件信息结构（不包含元信息）
type SimplifiedFileInfo struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Content string `json:"content"`
}

// SimplifiedFolderInfo 简化的文件夹信息结构（不包含元信息）
type SimplifiedFolderInfo struct {
	Path  string `json:"path"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Count int    `json:"count"`
}

// simplifyFiles 简化文件信息，移除元信息字段
func (f *BaseFormatter) simplifyFiles(files []types.FileInfo) []SimplifiedFileInfo {
	simplified := make([]SimplifiedFileInfo, len(files))
	for i, file := range files {
		simplified[i] = SimplifiedFileInfo{
			Path:    file.Path,
			Name:    file.Name,
			Size:    file.Size,
			Content: file.Content,
		}
	}
	return simplified
}

// simplifyFolders 简化文件夹信息，移除元信息字段
func (f *BaseFormatter) simplifyFolders(folders []types.FolderInfo) []SimplifiedFolderInfo {
	simplified := make([]SimplifiedFolderInfo, len(folders))
	for i, folder := range folders {
		simplified[i] = SimplifiedFolderInfo{
			Path:  folder.Path,
			Name:  folder.Name,
			Size:  folder.Size,
			Count: folder.Count,
		}
	}
	return simplified
}

// applyCustomStructure 应用自定义结构到数据
func (f *BaseFormatter) applyCustomStructure(data types.ContextData) interface{} {
	// 基础实现，直接返回原始数据
	// 各格式转换器可以重写此方法以提供特定的自定义结构处理
	return data
}

// applyCustomFields 应用自定义字段映射
func (f *BaseFormatter) applyCustomFields(data interface{}) interface{} {
	// 基础实现，直接返回原始数据
	// 各格式转换器可以重写此方法以提供特定的字段映射处理
	return data
}