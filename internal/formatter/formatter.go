// Package formatter 提供多种格式的输出转换功能
package formatter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"code-context-generator/internal/formatter/encoding"
	"code-context-generator/pkg/constants"
	"code-context-generator/pkg/types"
)

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

// applyCustomStructure 应用自定义结构
func (f *BaseFormatter) applyCustomStructure(data types.ContextData) interface{} {
	// 根据配置应用自定义结构
	if f.config != nil {
		// 尝试将配置转换为FormatConfig
		if formatConfig, ok := f.config.(*types.FormatConfig); ok && formatConfig.Structure != nil {
			// 创建基于实际数据的自定义结构
			result := make(map[string]interface{})

			// 首先复制所有自定义字段（除了已知的结构字段）
			for key, value := range formatConfig.Structure {
				switch key {
				case "root", "files", "folders":
					// 这些字段稍后单独处理
				default:
					// 复制自定义字段
					result[key] = value
				}
			}

			// 应用结构映射
			if rootTag, ok := formatConfig.Structure["root"].(string); ok && rootTag != "" {
				result["XMLName"] = xml.Name{Local: rootTag}
			} else {
				result["XMLName"] = xml.Name{Local: "context"}
			}

			// 检查是否包含元信息
			includeMetadata := false
			// 尝试从完整配置中获取IncludeMetadata设置
			if f.config != nil {
				if _, ok := f.config.(*types.FormatConfig); ok {
					// FormatConfig没有Output字段，需要从父配置获取
					// 这里使用默认的false值
					includeMetadata = false
				}
			}

			// 映射文件和文件夹数据
			if filesTag, ok := formatConfig.Structure["files"].(string); ok && filesTag != "" {
				if includeMetadata {
					result[filesTag] = map[string]interface{}{
						"file": data.Files,
					}
				} else {
					// 简化文件信息
					simplifiedFiles := make([]SimplifiedFileInfo, len(data.Files))
					for i, file := range data.Files {
						simplifiedFiles[i] = SimplifiedFileInfo{
							Path:    file.Path,
							Name:    file.Name,
							Size:    file.Size,
							Content: file.Content,
						}
					}
					result[filesTag] = map[string]interface{}{
						"file": simplifiedFiles,
					}
				}
			} else {
				if includeMetadata {
					result["files"] = map[string]interface{}{
						"file": data.Files,
					}
				} else {
					// 简化文件信息
					simplifiedFiles := make([]SimplifiedFileInfo, len(data.Files))
					for i, file := range data.Files {
						simplifiedFiles[i] = SimplifiedFileInfo{
							Path:    file.Path,
							Name:    file.Name,
							Size:    file.Size,
							Content: file.Content,
						}
					}
					result["files"] = map[string]interface{}{
						"file": simplifiedFiles,
					}
				}
			}

			if foldersTag, ok := formatConfig.Structure["folders"].(string); ok && foldersTag != "" {
				if includeMetadata {
					result[foldersTag] = map[string]interface{}{
						"folder": data.Folders,
					}
				} else {
					// 简化文件夹信息
					simplifiedFolders := make([]SimplifiedFolderInfo, len(data.Folders))
					for i, folder := range data.Folders {
						simplifiedFolders[i] = SimplifiedFolderInfo{
							Path:  folder.Path,
							Name:  folder.Name,
							Size:  folder.Size,
							Count: folder.Count,
						}
					}
					result[foldersTag] = map[string]interface{}{
						"folder": simplifiedFolders,
					}
				}
			} else {
				if includeMetadata {
					result["folders"] = map[string]interface{}{
						"folder": data.Folders,
					}
				} else {
					// 简化文件夹信息
					simplifiedFolders := make([]SimplifiedFolderInfo, len(data.Folders))
					for i, folder := range data.Folders {
						simplifiedFolders[i] = SimplifiedFolderInfo{
							Path:  folder.Path,
							Name:  folder.Name,
							Size:  folder.Size,
							Count: folder.Count,
						}
					}
					result["folders"] = map[string]interface{}{
						"folder": simplifiedFolders,
					}
				}
			}

			// 添加统计信息
			result["file_count"] = data.FileCount
			result["folder_count"] = data.FolderCount
			result["total_size"] = data.TotalSize

			return result
		}
	}

	// 检查是否包含元信息
	includeMetadata := false
	// 尝试从完整配置中获取IncludeMetadata设置
	if f.config != nil {
		if _, ok := f.config.(*types.FormatConfig); ok {
			// FormatConfig没有Output字段，需要从父配置获取
			// 这里使用默认的false值
			includeMetadata = false
		}
	}

	if includeMetadata {
		// 返回可序列化的结构，避免map[string]interface{}
		return struct {
			Files       []types.FileInfo       `json:"files"`
			Folders     []types.FolderInfo     `json:"folders"`
			FileCount   int                    `json:"file_count"`
			FolderCount int                    `json:"folder_count"`
			TotalSize   int64                  `json:"total_size"`
			Metadata    map[string]interface{} `json:"metadata"`
		}{
			Files:       data.Files,
			Folders:     data.Folders,
			FileCount:   data.FileCount,
			FolderCount: data.FolderCount,
			TotalSize:   data.TotalSize,
			Metadata:    data.Metadata,
		}
	} else {
		// 不包含元信息的简化结构
		simplifiedFiles := make([]SimplifiedFileInfo, len(data.Files))
		for i, file := range data.Files {
			simplifiedFiles[i] = SimplifiedFileInfo{
				Path:    file.Path,
				Name:    file.Name,
				Size:    file.Size,
				Content: file.Content,
			}
		}
		
		simplifiedFolders := make([]SimplifiedFolderInfo, len(data.Folders))
		for i, folder := range data.Folders {
			simplifiedFolders[i] = SimplifiedFolderInfo{
				Path:  folder.Path,
				Name:  folder.Name,
				Size:  folder.Size,
				Count: folder.Count,
			}
		}
		
		return struct {
			Files       []SimplifiedFileInfo   `json:"files"`
			Folders     []SimplifiedFolderInfo `json:"folders"`
			FileCount   int                    `json:"file_count"`
			FolderCount int                    `json:"folder_count"`
			TotalSize   int64                  `json:"total_size"`
		}{
			Files:       simplifiedFiles,
			Folders:     simplifiedFolders,
			FileCount:   data.FileCount,
			FolderCount: data.FolderCount,
			TotalSize:   data.TotalSize,
		}
	}

	// 默认返回简化结构
	simplifiedFiles := make([]SimplifiedFileInfo, len(data.Files))
	for i, file := range data.Files {
		simplifiedFiles[i] = SimplifiedFileInfo{
			Path:    file.Path,
			Name:    file.Name,
			Size:    file.Size,
			Content: file.Content,
		}
	}
	
	simplifiedFolders := make([]SimplifiedFolderInfo, len(data.Folders))
	for i, folder := range data.Folders {
		simplifiedFolders[i] = SimplifiedFolderInfo{
			Path:  folder.Path,
			Name:  folder.Name,
			Size:  folder.Size,
			Count: folder.Count,
		}
	}
	
	return struct {
		Files       []SimplifiedFileInfo   `json:"files"`
		Folders     []SimplifiedFolderInfo `json:"folders"`
		FileCount   int                    `json:"file_count"`
		FolderCount int                    `json:"folder_count"`
		TotalSize   int64                  `json:"total_size"`
	}{
		Files:       simplifiedFiles,
		Folders:     simplifiedFolders,
		FileCount:   data.FileCount,
		FolderCount: data.FolderCount,
		TotalSize:   data.TotalSize,
	}
}

// formatFileWithContentHandling 根据内容处理选项格式化文件
func formatFileWithContentHandling(file types.FileInfo, contentHandling types.XMLContentHandling) (string, error) {
	// 如果是二进制文件，不显示内容
	if file.IsBinary {
		file.Content = "[二进制文件 - 内容未显示]"
	}

	switch contentHandling {
	case types.XMLContentCDATA:
		// 使用CDATA包装内容
		type FileWithCDATA struct {
			XMLName  xml.Name `xml:"file"`
			Path     string   `xml:"path"`
			Name     string   `xml:"name"`
			Size     int64    `xml:"size"`
			Content  string   `xml:",cdata"`
			ModTime  string   `xml:"mod_time"`
			IsDir    bool     `xml:"is_dir"`
			IsHidden bool     `xml:"is_hidden"`
			IsBinary bool     `xml:"is_binary"`
		}

		fileWithCDATA := FileWithCDATA{
			Path:     file.Path,
			Name:     file.Name,
			Size:     file.Size,
			Content:  file.Content,
			ModTime:  file.ModTime.Format(time.RFC3339),
			IsDir:    file.IsDir,
			IsHidden: file.IsHidden,
			IsBinary: file.IsBinary,
		}

		output, err := xml.MarshalIndent(fileWithCDATA, "", "  ")
		if err != nil {
			return "", fmt.Errorf("XML文件格式化失败: %w", err)
		}
		return xml.Header + string(output), nil

	case types.XMLContentRaw:
		// 使用最小转义
		type FileWithRaw struct {
			XMLName  xml.Name `xml:"file"`
			Path     string   `xml:"path"`
			Name     string   `xml:"name"`
			Size     int64    `xml:"size"`
			Content  RawText  `xml:"content"`
			ModTime  string   `xml:"mod_time"`
			IsDir    bool     `xml:"is_dir"`
			IsHidden bool     `xml:"is_hidden"`
			IsBinary bool     `xml:"is_binary"`
		}

		fileWithRaw := FileWithRaw{
			Path:     file.Path,
			Name:     file.Name,
			Size:     file.Size,
			Content:  RawText{Text: file.Content},
			ModTime:  file.ModTime.Format(time.RFC3339),
			IsDir:    file.IsDir,
			IsHidden: file.IsHidden,
			IsBinary: file.IsBinary,
		}

		output, err := xml.MarshalIndent(fileWithRaw, "", "  ")
		if err != nil {
			return "", fmt.Errorf("XML文件格式化失败: %w", err)
		}
		return xml.Header + string(output), nil

	default:
		// 默认使用标准XML序列化（转义）
		output, err := xml.MarshalIndent(file, "", "  ")
		if err != nil {
			return "", fmt.Errorf("XML文件格式化失败: %w", err)
		}
		return xml.Header + string(output), nil
	}
}

// formatFolderWithContentHandling 根据内容处理选项格式化文件夹
func formatFolderWithContentHandling(folder types.FolderInfo, contentHandling types.XMLContentHandling) (string, error) {
	switch contentHandling {
	case types.XMLContentCDATA:
		// 使用CDATA包装内容（主要用于文件内容，文件夹较少使用）
		// 这里仍然使用标准序列化，因为文件夹主要包含元数据
		output, err := xml.MarshalIndent(folder, "", "  ")
		if err != nil {
			return "", fmt.Errorf("XML文件夹格式化失败: %w", err)
		}
		return xml.Header + string(output), nil

	case types.XMLContentRaw:
		// 使用最小转义
		type FolderWithRaw struct {
			XMLName  xml.Name           `xml:"folder"`
			Path     string             `xml:"path"`
			Name     string             `xml:"name"`
			Files    []types.FileInfo   `xml:"files"`
			Folders  []types.FolderInfo `xml:"folders"`
			ModTime  string             `xml:"mod_time"`
			IsHidden bool               `xml:"is_hidden"`
			Size     int64              `xml:"size"`
			Count    int                `xml:"count"`
		}

		folderWithRaw := FolderWithRaw{
			Path:     folder.Path,
			Name:     folder.Name,
			Files:    folder.Files,
			Folders:  folder.Folders,
			ModTime:  folder.ModTime.Format(time.RFC3339),
			IsHidden: folder.IsHidden,
			Size:     folder.Size,
			Count:    folder.Count,
		}

		output, err := xml.MarshalIndent(folderWithRaw, "", "  ")
		if err != nil {
			return "", fmt.Errorf("XML文件夹格式化失败: %w", err)
		}
		return xml.Header + string(output), nil

	default:
		// 默认使用标准XML序列化（转义）
		output, err := xml.MarshalIndent(folder, "", "  ")
		if err != nil {
			return "", fmt.Errorf("XML文件夹格式化失败: %w", err)
		}
		return xml.Header + string(output), nil
	}
}

// applyCustomFields 应用自定义字段映射
func (f *BaseFormatter) applyCustomFields(file types.FileInfo) interface{} {
	// 根据配置应用自定义字段映射
	if f.config != nil {
		// 尝试将配置转换为FormatConfig
		if formatConfig, ok := f.config.(*types.FormatConfig); ok && formatConfig.Fields != nil {
			// 这里可以实现字段映射逻辑
			return formatConfig.Fields
		}
	}
	return file
}

// JSONFormatter JSON格式转换器
type JSONFormatter struct {
	BaseFormatter
	config *types.Config // 修改为接收完整的Config
}

// NewJSONFormatter 创建JSON格式转换器
func NewJSONFormatter(config *types.Config) Formatter { // 修改参数类型
	var formatConfig *types.FormatConfig
	if config != nil {
		formatConfig = &config.Formats.JSON
	}
	return &JSONFormatter{
		BaseFormatter: BaseFormatter{
			name:        "JSON",
			description: "JavaScript Object Notation format",
			config:      formatConfig,
		},
		config: config, // 保存完整的配置
	}
}

// Format 格式化上下文数据
func (f *JSONFormatter) Format(data types.ContextData) (string, error) {
	var outputData interface{}
	
	// 检查是否有真正的自定义结构配置（不只是默认的file/folder映射）
	hasCustomStructure := false
	if f.config != nil && f.config.Formats.JSON.Structure != nil {
		// 如果有超过2个字段，或者字段值不是默认的file/folder，则认为有自定义结构
		if len(f.config.Formats.JSON.Structure) > 2 {
			hasCustomStructure = true
		} else {
			// 检查是否有非默认的字段映射
			for key, value := range f.config.Formats.JSON.Structure {
				if key != "file" && key != "folder" {
					hasCustomStructure = true
					break
				}
				if key == "file" && value != "file" {
					hasCustomStructure = true
					break
				}
				if key == "folder" && value != "folder" {
					hasCustomStructure = true
					break
				}
			}
		}
	}
	
	if hasCustomStructure {
		// 使用自定义结构
		outputData = f.applyCustomStructure(data)
	} else {
		// 检查是否包含元信息
		includeMetadata := false // 默认不包含元信息
		if f.config != nil {
			includeMetadata = f.config.Output.IncludeMetadata
		}

		if includeMetadata {
			// 包含元信息
			outputData = data
		} else {
			// 不包含元信息的简化结构
			simplifiedFiles := make([]SimplifiedFileInfo, len(data.Files))
			for i, file := range data.Files {
				simplifiedFiles[i] = SimplifiedFileInfo{
					Path:    file.Path,
					Name:    file.Name,
					Size:    file.Size,
					Content: file.Content,
				}
			}
			
			simplifiedFolders := make([]SimplifiedFolderInfo, len(data.Folders))
			for i, folder := range data.Folders {
				simplifiedFolders[i] = SimplifiedFolderInfo{
					Path:  folder.Path,
					Name:  folder.Name,
					Size:  folder.Size,
					Count: folder.Count,
				}
			}
			
			outputData = struct {
				Files   []SimplifiedFileInfo   `json:"files"`
				Folders []SimplifiedFolderInfo `json:"folders"`
				FileCount   int                    `json:"file_count"`
				FolderCount int                    `json:"folder_count"`
				TotalSize   int64                  `json:"total_size"`
			}{
				Files:       simplifiedFiles,
				Folders:     simplifiedFolders,
				FileCount:   data.FileCount,
				FolderCount: data.FolderCount,
				TotalSize:   data.TotalSize,
			}
		}
	}

	output, err := json.MarshalIndent(outputData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON格式化失败: %w", err)
	}

	// 检查配置中是否指定了编码格式
	if f.config != nil && f.config.Formats.JSON.Encoding != "" && f.config.Formats.JSON.Encoding != "utf-8" {
		// 转换编码
		encodedOutput, err := encoding.ConvertEncoding(string(output), f.config.Formats.JSON.Encoding)
		if err != nil {
			return "", fmt.Errorf("编码转换失败: %w", err)
		}
		return encodedOutput, nil
	}

	return string(output), nil
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

// FormatFile 格式化单个文件
func (f *JSONFormatter) FormatFile(file types.FileInfo) (string, error) {
	// 如果是二进制文件，不显示内容
	if file.IsBinary {
		file.Content = "[二进制文件 - 内容未显示]"
	}

	// 检查是否包含元信息
	includeMetadata := false // 默认不包含元信息
	if f.config != nil {
		includeMetadata = f.config.Output.IncludeMetadata
	}

	var output []byte
	var err error

	if includeMetadata {
		// 使用自定义字段映射
		customFile := f.applyCustomFields(file)
		output, err = json.MarshalIndent(customFile, "", "  ")
	} else {
		// 不包含元信息的简化结构
		simplifiedFile := SimplifiedFileInfo{
			Path:    file.Path,
			Name:    file.Name,
			Size:    file.Size,
			Content: file.Content,
		}
		output, err = json.MarshalIndent(simplifiedFile, "", "  ")
	}

	if err != nil {
		return "", fmt.Errorf("JSON文件格式化失败: %w", err)
	}
	return string(output), nil
}

// FormatFolder 格式化文件夹
func (f *JSONFormatter) FormatFolder(folder types.FolderInfo) (string, error) {
	// 检查是否包含元信息
	includeMetadata := false // 默认不包含元信息
	if f.config != nil {
		includeMetadata = f.config.Output.IncludeMetadata
	}

	var output []byte
	var err error

	if includeMetadata {
		output, err = json.MarshalIndent(folder, "", "  ")
	} else {
		// 不包含元信息的简化结构
		simplifiedFolder := SimplifiedFolderInfo{
			Path:  folder.Path,
			Name:  folder.Name,
			Size:  folder.Size,
			Count: folder.Count,
		}
		output, err = json.MarshalIndent(simplifiedFolder, "", "  ")
	}

	if err != nil {
		return "", fmt.Errorf("JSON文件夹格式化失败: %w", err)
	}
	return string(output), nil
}

// XMLFormatter XML格式转换器
type XMLFormatter struct {
	BaseFormatter
	config *types.Config
}

// NewXMLFormatter 创建XML格式转换器
func NewXMLFormatter(config *types.Config) Formatter {
	var formatConfig *types.FormatConfig
	if config != nil {
		formatConfig = &config.Formats.XML.FormatConfig
	}
	return &XMLFormatter{
		BaseFormatter: BaseFormatter{
			name:        "XML",
			description: "Extensible Markup Language format",
			config:      formatConfig,
		},
		config: config,
	}
}

// Format 格式化上下文数据
func (f *XMLFormatter) Format(data types.ContextData) (string, error) {
	// 检查是否包含元信息
	includeMetadata := false // 默认不包含元信息
	if f.config != nil {
		includeMetadata = f.config.Output.IncludeMetadata
	}

	if f.config != nil && f.config.Formats.XML.Structure != nil {
		// 使用自定义结构
		customData := f.applyCustomStructure(data)
		output, err := xml.MarshalIndent(customData, "", "  ")
		if err != nil {
			return "", fmt.Errorf("XML格式化失败: %w", err)
		}
		return xml.Header + string(output), nil
	}

	// 根据是否包含元信息创建不同的数据结构
	var output []byte
	var err error

	if includeMetadata {
		// 定义元数据项结构
		type MetadataItem struct {
			Key   string      `xml:"key,attr"`
			Value interface{} `xml:"value"`
		}

		// 包含元信息的默认结构
		type SerializableContextData struct {
			XMLName     xml.Name               `xml:"context"`
			Files       []types.FileInfo       `xml:"files>file"`
			Folders     []types.FolderInfo     `xml:"folders>folder"`
			FileCount   int                    `xml:"file_count"`
			FolderCount int                    `xml:"folder_count"`
			TotalSize   int64                  `xml:"total_size"`
			Metadata    []MetadataItem         `xml:"metadata>item"`
		}

		// 转换metadata为可序列化的结构
		var metadataItems []MetadataItem
		for key, value := range data.Metadata {
			metadataItems = append(metadataItems, MetadataItem{
				Key:   key,
				Value: value,
			})
		}

		serializableData := SerializableContextData{
			Files:       data.Files,
			Folders:     data.Folders,
			FileCount:   data.FileCount,
			FolderCount: data.FolderCount,
			TotalSize:   data.TotalSize,
			Metadata:    metadataItems,
		}
		output, err = xml.MarshalIndent(serializableData, "", "  ")
	} else {
		// 不包含元信息的简化结构
		type SimplifiedContextData struct {
			XMLName     xml.Name               `xml:"context"`
			Files       []SimplifiedFileInfo   `xml:"files>file"`
			Folders     []SimplifiedFolderInfo `xml:"folders>folder"`
			FileCount   int                    `xml:"file_count"`
			FolderCount int                    `xml:"folder_count"`
			TotalSize   int64                  `xml:"total_size"`
		}
		
		simplifiedData := SimplifiedContextData{
			Files:       f.simplifyFiles(data.Files),
			Folders:     f.simplifyFolders(data.Folders),
			FileCount:   data.FileCount,
			FolderCount: data.FolderCount,
			TotalSize:   data.TotalSize,
		}
		output, err = xml.MarshalIndent(simplifiedData, "", "  ")
	}

	if err != nil {
		return "", fmt.Errorf("XML格式化失败: %w", err)
	}

	result := xml.Header + string(output)

	// 检查配置中是否指定了编码格式
	if f.config != nil && f.config.Formats.XML.FormatConfig.Encoding != "" && f.config.Formats.XML.FormatConfig.Encoding != "utf-8" {
		// 转换编码
		encodedOutput, err := encoding.ConvertEncoding(result, f.config.Formats.XML.FormatConfig.Encoding)
		if err != nil {
			return "", fmt.Errorf("编码转换失败: %w", err)
		}
		return encodedOutput, nil
	}

	return result, nil
}

// FormatFile 格式化单个文件
func (f *XMLFormatter) FormatFile(file types.FileInfo) (string, error) {
	// 如果是二进制文件，不显示内容
	if file.IsBinary {
		file.Content = "[二进制文件 - 内容未显示]"
	}

	// 检查是否包含元信息
	includeMetadata := false // 默认不包含元信息
	if f.config != nil {
		includeMetadata = f.config.Output.IncludeMetadata
	}

	// 根据内容处理选项处理文件内容
	if f.config != nil && f.config.Formats.XML.Formatting.ContentHandling != "" {
		return formatFileWithContentHandling(file, f.config.Formats.XML.Formatting.ContentHandling)
	}

	// 根据是否包含元信息创建不同的数据结构
	var output []byte
	var err error

	if includeMetadata {
		// 包含元信息的默认结构
		output, err = xml.MarshalIndent(file, "", "  ")
	} else {
		// 不包含元信息的简化结构
		simplifiedFile := SimplifiedFileInfo{
			Path:    file.Path,
			Name:    file.Name,
			Size:    file.Size,
			Content: file.Content,
		}
		output, err = xml.MarshalIndent(simplifiedFile, "", "  ")
	}

	if err != nil {
		return "", fmt.Errorf("XML文件格式化失败: %w", err)
	}
	return xml.Header + string(output), nil
}

// FormatFolder 格式化文件夹
func (f *XMLFormatter) FormatFolder(folder types.FolderInfo) (string, error) {
	// 检查是否包含元信息
	includeMetadata := false // 默认不包含元信息
	if f.config != nil {
		includeMetadata = f.config.Output.IncludeMetadata
	}

	// 如果配置了内容处理选项，使用相应的处理方式
	if f.config != nil && f.config.Formats.XML.Formatting.ContentHandling != "" {
		return formatFolderWithContentHandling(folder, f.config.Formats.XML.Formatting.ContentHandling)
	}

	// 根据是否包含元信息创建不同的数据结构
	var output []byte
	var err error

	if includeMetadata {
		// 包含元信息的默认结构
		output, err = xml.MarshalIndent(folder, "", "  ")
	} else {
		// 不包含元信息的简化结构
		simplifiedFolder := SimplifiedFolderInfo{
			Path:  folder.Path,
			Name:  folder.Name,
			Size:  folder.Size,
			Count: folder.Count,
		}
		output, err = xml.MarshalIndent(simplifiedFolder, "", "  ")
	}

	if err != nil {
		return "", fmt.Errorf("XML文件夹格式化失败: %w", err)
	}
	return xml.Header + string(output), nil
}

// TOMLFormatter TOML格式转换器
type TOMLFormatter struct {
	BaseFormatter
	encoding string
	config   *types.Config // 添加完整配置引用
}

// NewTOMLFormatter 创建TOML格式转换器
func NewTOMLFormatter(config *types.Config) Formatter { // 修改参数类型
	encoding := "utf-8"
	var formatConfig *types.FormatConfig
	if config != nil {
		formatConfig = &config.Formats.TOML
		if formatConfig.Encoding != "" {
			encoding = formatConfig.Encoding
		}
	}
	return &TOMLFormatter{
		BaseFormatter: BaseFormatter{
			name:        "TOML",
			description: "Tom's Obvious, Minimal Language format",
			config:      formatConfig,
		},
		encoding: encoding,
		config:   config, // 保存完整配置
	}
}

// Format 格式化上下文数据
func (f *TOMLFormatter) Format(data types.ContextData) (string, error) {
	var buf strings.Builder

	// 检查是否包含元信息
	includeMetadata := false // 默认不包含元信息
	if f.config != nil {
		includeMetadata = f.config.Output.IncludeMetadata
	}

	// 写入文件部分
	if len(data.Files) > 0 {
		buf.WriteString("[files]\n")
		for i, file := range data.Files {
			buf.WriteString("  [[files.file]]\n")
			buf.WriteString(fmt.Sprintf("    path = \"%s\"\n", file.Path))
			buf.WriteString(fmt.Sprintf("    name = \"%s\"\n", file.Name))
			buf.WriteString(fmt.Sprintf("    size = %d\n", file.Size))
			buf.WriteString(fmt.Sprintf("    content = \"%s\"\n", encoding.EscapeTOMLString(file.Content)))
			
			if includeMetadata {
				buf.WriteString(fmt.Sprintf("    mod_time = \"%s\"\n", file.ModTime.Format(time.RFC3339)))
			}
			
			if i < len(data.Files)-1 {
				buf.WriteString("\n")
			}
		}
	}

	// 写入文件夹部分
	if len(data.Folders) > 0 {
		buf.WriteString("\n[folders]\n")
		for i, folder := range data.Folders {
			buf.WriteString("  [[folders.folder]]\n")
			buf.WriteString(fmt.Sprintf("    path = \"%s\"\n", folder.Path))
			buf.WriteString(fmt.Sprintf("    name = \"%s\"\n", folder.Name))
			buf.WriteString(fmt.Sprintf("    file_count = %d\n", len(folder.Files)))
			
			if includeMetadata {
				buf.WriteString(fmt.Sprintf("    mod_time = \"%s\"\n", folder.ModTime.Format(time.RFC3339)))
			}
			
			if i < len(data.Folders)-1 {
				buf.WriteString("\n")
			}
		}
	}

	// 写入元信息部分
	if includeMetadata && len(data.Metadata) > 0 {
		buf.WriteString("\n[metadata]\n")
		
		// 处理root_path
		if rootPath, exists := data.Metadata["root_path"]; exists {
			buf.WriteString(fmt.Sprintf("root_path = \"%s\"\n", rootPath))
		}
		
		// 处理git信息
		if gitInfo, exists := data.Metadata["git"]; exists {
			buf.WriteString("\n[metadata.git]\n")
			// 将git信息转换为JSON字符串以便在TOML中显示
			gitJSON, err := json.Marshal(gitInfo)
			if err == nil {
				buf.WriteString(fmt.Sprintf("git_info = %s\n", string(gitJSON)))
			} else {
				buf.WriteString(fmt.Sprintf("git_info = \"%v\"\n", gitInfo))
			}
		}
	}

	result := buf.String()

	// 检查是否需要编码转换
	if f.encoding != "" && f.encoding != "utf-8" {
		encodedResult, err := encoding.ConvertEncoding(result, f.encoding)
		if err != nil {
			return "", fmt.Errorf("编码转换失败: %w", err)
		}
		return encodedResult, nil
	}

	return result, nil
}

// FormatFile 格式化单个文件
func (f *TOMLFormatter) FormatFile(file types.FileInfo) (string, error) {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("path = \"%s\"\n", file.Path))
	buf.WriteString(fmt.Sprintf("name = \"%s\"\n", file.Name))
	buf.WriteString(fmt.Sprintf("size = %d\n", file.Size))

	// 如果是二进制文件，不显示内容
	if file.IsBinary {
		buf.WriteString("content = \"[二进制文件 - 内容未显示]\"\n")
	} else {
		buf.WriteString(fmt.Sprintf("content = \"%s\"\n", encoding.EscapeTOMLString(file.Content)))
	}

	// 检查是否包含元信息
	includeMetadata := false // 默认不包含元信息
	if f.config != nil {
		includeMetadata = f.config.Output.IncludeMetadata
	}

	if includeMetadata {
		buf.WriteString(fmt.Sprintf("mod_time = \"%s\"\n", file.ModTime.Format(time.RFC3339)))
	}

	return buf.String(), nil
}

// FormatFolder 格式化文件夹
func (f *TOMLFormatter) FormatFolder(folder types.FolderInfo) (string, error) {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("path = \"%s\"\n", folder.Path))
	buf.WriteString(fmt.Sprintf("name = \"%s\"\n", folder.Name))
	buf.WriteString(fmt.Sprintf("file_count = %d\n", len(folder.Files)))

	// 检查是否包含元信息
	includeMetadata := true // 默认包含元信息
	if f.config != nil {
		includeMetadata = f.config.Output.IncludeMetadata
	}

	if includeMetadata {
		buf.WriteString(fmt.Sprintf("mod_time = \"%s\"\n", folder.ModTime.Format(time.RFC3339)))
	}

	return buf.String(), nil
}

// MarkdownFormatter Markdown格式转换器
type MarkdownFormatter struct {
	BaseFormatter
	encoding string
	config   *types.Config // 添加完整配置引用
}

// NewMarkdownFormatter 创建Markdown格式转换器
func NewMarkdownFormatter(config *types.Config) Formatter { // 修改参数类型
	encoding := "utf-8"
	var formatConfig *types.FormatConfig
	if config != nil {
		formatConfig = &config.Formats.Markdown
		if formatConfig.Encoding != "" {
			encoding = formatConfig.Encoding
		}
	}
	return &MarkdownFormatter{
		BaseFormatter: BaseFormatter{
			name:        "Markdown",
			description: "Markdown format with code blocks",
			config:      formatConfig,
		},
		encoding: encoding,
		config:   config, // 保存完整配置
	}
}

// Format 格式化上下文数据
func (f *MarkdownFormatter) Format(data types.ContextData) (string, error) {
	var sb strings.Builder

	// 检查是否包含元信息
	includeMetadata := true // 默认包含元信息
	if f.config != nil {
		includeMetadata = f.config.Output.IncludeMetadata
	}

	// 添加标题
	sb.WriteString("# 代码上下文\n\n")
	sb.WriteString(fmt.Sprintf("生成时间: %s\n\n", time.Now().Format(time.RFC3339)))

	// 添加文件部分
	if len(data.Files) > 0 {
		sb.WriteString("## 文件\n\n")
		for _, file := range data.Files {
			sb.WriteString(fmt.Sprintf("### %s\n\n", file.Name))
			sb.WriteString(fmt.Sprintf("- **路径**: `%s`\n", file.Path))
			sb.WriteString(fmt.Sprintf("- **大小**: %d 字节\n", file.Size))
			if includeMetadata {
				sb.WriteString(fmt.Sprintf("- **修改时间**: %s\n\n", file.ModTime.Format(time.RFC3339)))
			}

			// 添加代码块（只针对文本文件）
			if !file.IsBinary {
				sb.WriteString("```")
				if ext := filepath.Ext(file.Path); ext != "" {
					sb.WriteString(strings.TrimPrefix(ext, "."))
				}
				sb.WriteString("\n")
				sb.WriteString(file.Content)
				sb.WriteString("\n```\n\n")
			} else {
				sb.WriteString("**[二进制文件 - 内容未显示]**\n\n")
			}
		}
	}

	// 添加文件夹部分
	if len(data.Folders) > 0 {
		sb.WriteString("## 文件夹\n\n")
		for _, folder := range data.Folders {
			sb.WriteString(fmt.Sprintf("### %s\n\n", folder.Name))
			sb.WriteString(fmt.Sprintf("- **路径**: `%s`\n", folder.Path))
			sb.WriteString(fmt.Sprintf("- **文件数**: %d\n", len(folder.Files)))
			if includeMetadata {
				sb.WriteString(fmt.Sprintf("- **修改时间**: %s\n\n", folder.ModTime.Format(time.RFC3339)))
			}

			// 添加文件夹中的文件
			if len(folder.Files) > 0 {
				sb.WriteString("#### 文件列表\n\n")
				for _, file := range folder.Files {
					sb.WriteString(fmt.Sprintf("- `%s` (%d 字节)\n", file.Name, file.Size))
				}
				sb.WriteString("\n")
			}
		}
	}

	// 添加元信息部分
	if includeMetadata && len(data.Metadata) > 0 {
		sb.WriteString("## 元信息\n\n")
		
		// 处理root_path
		if rootPath, exists := data.Metadata["root_path"]; exists {
			sb.WriteString(fmt.Sprintf("- **根路径**: `%s`\n", rootPath))
		}
		
		// 处理git信息
		if gitInfo, exists := data.Metadata["git"]; exists {
			sb.WriteString("\n### Git信息\n\n")
			// 将git信息转换为JSON字符串以便显示
			gitJSON, err := json.MarshalIndent(gitInfo, "", "  ")
			if err == nil {
				sb.WriteString("```json\n")
				sb.WriteString(string(gitJSON))
				sb.WriteString("\n```\n")
			} else {
				sb.WriteString(fmt.Sprintf("```\n%v\n```\n", gitInfo))
			}
		}
	}

	result := sb.String()

	// 检查是否需要编码转换
	if f.encoding != "" && f.encoding != "utf-8" {
		encodedResult, err := encoding.ConvertEncoding(result, f.encoding)
		if err != nil {
			return "", fmt.Errorf("编码转换失败: %w", err)
		}
		return encodedResult, nil
	}

	return result, nil
}

// FormatFile 格式化单个文件
func (f *MarkdownFormatter) FormatFile(file types.FileInfo) (string, error) {
	var sb strings.Builder

	// 检查是否包含元信息
	includeMetadata := true // 默认包含元信息
	if f.config != nil {
		includeMetadata = f.config.Output.IncludeMetadata
	}

	sb.WriteString(fmt.Sprintf("## %s\n\n", file.Name))
	sb.WriteString(fmt.Sprintf("- **路径**: `%s`\n", file.Path))
	sb.WriteString(fmt.Sprintf("- **大小**: %d 字节\n", file.Size))
	if includeMetadata {
		sb.WriteString(fmt.Sprintf("- **修改时间**: %s\n\n", file.ModTime.Format(time.RFC3339)))
	}

	// 添加代码块（只针对文本文件）
	if !file.IsBinary {
		sb.WriteString("```")
		if ext := filepath.Ext(file.Path); ext != "" {
			sb.WriteString(strings.TrimPrefix(ext, "."))
		}
		sb.WriteString("\n")
		sb.WriteString(file.Content)
		sb.WriteString("\n```\n")
	} else {
		sb.WriteString("**[二进制文件 - 内容未显示]**\n")
	}

	return sb.String(), nil
}

// FormatFolder 格式化文件夹
func (f *MarkdownFormatter) FormatFolder(folder types.FolderInfo) (string, error) {
	var sb strings.Builder

	// 检查是否包含元信息
	includeMetadata := true // 默认包含元信息
	if f.config != nil {
		includeMetadata = f.config.Output.IncludeMetadata
	}

	sb.WriteString(fmt.Sprintf("## %s\n\n", folder.Name))
	sb.WriteString(fmt.Sprintf("- **路径**: `%s`\n", folder.Path))
	sb.WriteString(fmt.Sprintf("- **文件数**: %d\n", len(folder.Files)))
	if includeMetadata {
		sb.WriteString(fmt.Sprintf("- **修改时间**: %s\n\n", folder.ModTime.Format(time.RFC3339)))
	}

	// 添加文件列表（始终显示）
	if len(folder.Files) > 0 {
		sb.WriteString("### 文件列表\n\n")
		for _, file := range folder.Files {
			sb.WriteString(fmt.Sprintf("- `%s` (%d 字节)\n", file.Name, file.Size))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// FormatterFactory 格式转换器工厂
type FormatterFactory struct {
	formatters map[string]Formatter
}

// NewFormatterFactory 创建格式转换器工厂
func NewFormatterFactory() *FormatterFactory {
	return &FormatterFactory{
		formatters: make(map[string]Formatter),
	}
}

// Register 注册格式转换器
func (ff *FormatterFactory) Register(format string, formatter Formatter) {
	ff.formatters[strings.ToLower(format)] = formatter
}

// Get 获取格式转换器
func (ff *FormatterFactory) Get(format string) (Formatter, error) {
	formatter, exists := ff.formatters[strings.ToLower(format)]
	if !exists {
		return nil, fmt.Errorf("不支持的格式: %s", format)
	}
	return formatter, nil
}

// GetSupportedFormats 获取支持的格式列表
func (ff *FormatterFactory) GetSupportedFormats() []string {
	formats := make([]string, 0, len(ff.formatters))
	for format := range ff.formatters {
		formats = append(formats, format)
	}
	return formats
}

// NewFormatter 创建格式转换器
func NewFormatter(format string, config *types.Config) (Formatter, error) {
	factory := CreateDefaultFactory(config)
	return factory.Get(format)
}

// CreateDefaultFactory 创建默认的格式转换器工厂
func CreateDefaultFactory(config *types.Config) *FormatterFactory {
	factory := NewFormatterFactory()
	factory.Register(constants.FormatJSON, NewJSONFormatter(config)) // 修复调用
	factory.Register(constants.FormatXML, NewXMLFormatter(config))
	factory.Register(constants.FormatTOML, NewTOMLFormatter(config))
	factory.Register(constants.FormatMarkdown, NewMarkdownFormatter(config))
	return factory
}
