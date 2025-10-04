package formatter

import (
	"fmt"

	"code-context-generator/internal/formatter/encoding"
	"code-context-generator/pkg/types"
	"github.com/BurntSushi/toml"
)

// TOMLFormatter TOML格式转换器
type TOMLFormatter struct {
	BaseFormatter
	config *types.Config
}

// NewTOMLFormatter 创建TOML格式转换器
func NewTOMLFormatter(config *types.Config) Formatter {
	return &TOMLFormatter{
		BaseFormatter: BaseFormatter{
			name:        "TOML",
			description: "Tom's Obvious, Minimal Language format",
			config:      nil,
		},
		config: config,
	}
}

// Format 格式化上下文数据
func (f *TOMLFormatter) Format(data types.ContextData) (string, error) {
	// 检查是否包含元信息
	includeMetadata := false // 默认不包含元信息
	if f.config != nil {
		includeMetadata = f.config.Output.IncludeMetadata
	}

	if f.config != nil && f.config.Formats.TOML.Structure != nil {
		// 使用自定义结构
		customData := f.applyCustomStructure(data)
		output, err := toml.Marshal(customData)
		if err != nil {
			return "", fmt.Errorf("TOML格式化失败: %w", err)
		}
		return string(output), nil
	}

	// 根据是否包含元信息创建不同的数据结构
	var output []byte
	var err error

	if includeMetadata {
		// 包含元信息的默认结构
		type SerializableContextData struct {
			Files       []types.FileInfo   `toml:"files"`
			Folders     []types.FolderInfo `toml:"folders"`
			FileCount   int                `toml:"file_count"`
			FolderCount int                `toml:"folder_count"`
			TotalSize   int64              `toml:"total_size"`
			Metadata    map[string]interface{} `toml:"metadata"`
		}

		serializableData := SerializableContextData{
			Files:       data.Files,
			Folders:     data.Folders,
			FileCount:   data.FileCount,
			FolderCount: data.FolderCount,
			TotalSize:   data.TotalSize,
			Metadata:    data.Metadata,
		}
		output, err = toml.Marshal(serializableData)
	} else {
		// 不包含元信息的简化结构
		type SimplifiedContextData struct {
			Files       []SimplifiedFileInfo   `toml:"files"`
			Folders     []SimplifiedFolderInfo `toml:"folders"`
			FileCount   int                    `toml:"file_count"`
			FolderCount int                    `toml:"folder_count"`
			TotalSize   int64                  `toml:"total_size"`
		}
		
		simplifiedData := SimplifiedContextData{
			Files:       f.simplifyFiles(data.Files),
			Folders:     f.simplifyFolders(data.Folders),
			FileCount:   data.FileCount,
			FolderCount: data.FolderCount,
			TotalSize:   data.TotalSize,
		}
		output, err = toml.Marshal(simplifiedData)
	}

	if err != nil {
		return "", fmt.Errorf("TOML格式化失败: %w", err)
	}

	result := string(output)

	// 检查配置中是否指定了编码格式
	if f.config != nil && f.config.Formats.TOML.Encoding != "" && f.config.Formats.TOML.Encoding != "utf-8" {
		// 转换编码
		encodedOutput, err := encoding.ConvertEncoding(result, f.config.Formats.TOML.Encoding)
		if err != nil {
			return "", fmt.Errorf("编码转换失败: %w", err)
		}
		return encodedOutput, nil
	}

	return result, nil
}

// FormatFile 格式化单个文件
func (f *TOMLFormatter) FormatFile(file types.FileInfo) (string, error) {
	// 如果是二进制文件，不显示内容
	if file.IsBinary {
		file.Content = "[二进制文件 - 内容未显示]"
	}

	// 使用TOML序列化文件信息
	output, err := toml.Marshal(file)
	if err != nil {
		return "", fmt.Errorf("TOML文件格式化失败: %w", err)
	}

	result := string(output)

	// 检查配置中是否指定了编码格式
	if f.config != nil && f.config.Formats.TOML.Encoding != "" && f.config.Formats.TOML.Encoding != "utf-8" {
		// 转换编码
		encodedOutput, err := encoding.ConvertEncoding(result, f.config.Formats.TOML.Encoding)
		if err != nil {
			return "", fmt.Errorf("编码转换失败: %w", err)
		}
		return encodedOutput, nil
	}

	return result, nil
}

// FormatFolder 格式化文件夹
func (f *TOMLFormatter) FormatFolder(folder types.FolderInfo) (string, error) {
	// 使用TOML序列化文件夹信息
	output, err := toml.Marshal(folder)
	if err != nil {
		return "", fmt.Errorf("TOML文件夹格式化失败: %w", err)
	}

	result := string(output)

	// 检查配置中是否指定了编码格式
	if f.config != nil && f.config.Formats.TOML.Encoding != "" && f.config.Formats.TOML.Encoding != "utf-8" {
		// 转换编码
		encodedOutput, err := encoding.ConvertEncoding(result, f.config.Formats.TOML.Encoding)
		if err != nil {
			return "", fmt.Errorf("编码转换失败: %w", err)
		}
		return encodedOutput, nil
	}

	return result, nil
}