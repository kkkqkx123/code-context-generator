package formatter

import (
	"encoding/json"
	"fmt"

	"code-context-generator/internal/formatter/encoding"
	"code-context-generator/pkg/types"
)

// JSONFormatter JSON格式转换器
type JSONFormatter struct {
	BaseFormatter
	config *types.Config // 接收完整的Config
}

// NewJSONFormatter 创建JSON格式转换器
func NewJSONFormatter(config *types.Config) Formatter {
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
	if f.config != nil {
		// 只有在配置存在时才检查结构配置
		if f.config.Formats.JSON.Structure != nil {
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
			simplifiedFiles := f.simplifyFiles(data.Files)
			simplifiedFolders := f.simplifyFolders(data.Folders)
			
			outputData = struct {
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
	}

	// 尝试序列化数据，捕获可能的panic
	var output []byte
	var err error
	
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("JSON序列化失败: %v", r)
			}
		}()
		output, err = json.MarshalIndent(outputData, "", "  ")
	}()
	
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

// applyCustomStructure 应用自定义结构到数据
func (f *JSONFormatter) applyCustomStructure(data types.ContextData) interface{} {
	// 如果配置中有自定义结构，使用它
	if f.config != nil && f.config.Formats.JSON.Structure != nil {
		// 创建自定义结构映射
		customStructure := make(map[string]interface{})
		
		// 应用配置的结构映射
		for key, value := range f.config.Formats.JSON.Structure {
			customStructure[key] = value
		}
		
		// 添加一些基本的数据字段（如果适用）
		if _, exists := customStructure["files"]; !exists {
			customStructure["files"] = data.Files
		}
		if _, exists := customStructure["folders"]; !exists {
			customStructure["folders"] = data.Folders
		}
		if _, exists := customStructure["file_count"]; !exists {
			customStructure["file_count"] = data.FileCount
		}
		if _, exists := customStructure["folder_count"]; !exists {
			customStructure["folder_count"] = data.FolderCount
		}
		if _, exists := customStructure["total_size"]; !exists {
			customStructure["total_size"] = data.TotalSize
		}
		if _, exists := customStructure["metadata"]; !exists && f.config != nil && f.config.Output.IncludeMetadata {
			customStructure["metadata"] = data.Metadata
		}
		
		return customStructure
	}
	
	// 如果没有自定义结构配置，返回原始数据
	return data
}

// applyCustomFields 应用自定义字段映射
func (f *JSONFormatter) applyCustomFields(data interface{}) interface{} {
	// 检查是否有自定义字段配置
	if f.config != nil && f.config.Formats.JSON.Fields != nil {
		// 创建自定义字段映射
		customFields := make(map[string]interface{})
		
		// 应用配置的字段映射
		for key, value := range f.config.Formats.JSON.Fields {
			customFields[key] = value
		}
		
		// 添加原始数据的字段（如果适用）
		if fileInfo, ok := data.(types.FileInfo); ok {
			// 如果是文件信息，添加一些基本字段
			customFields["path"] = fileInfo.Path
			customFields["name"] = fileInfo.Name
			customFields["size"] = fileInfo.Size
		}
		
		return customFields
	}
	
	// 如果没有自定义字段配置，返回原始数据
	return data
}