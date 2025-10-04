package formatter

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"code-context-generator/internal/formatter/encoding"
	"code-context-generator/pkg/types"
)

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
	// 检查是否启用AI优化
	if f.config != nil && f.config.Output.AIOptimized {
		return f.formatAIOptimized(data)
	}

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
			XMLName     xml.Name           `xml:"context"`
			Files       []types.FileInfo   `xml:"files>file"`
			Folders     []types.FolderInfo `xml:"folders>folder"`
			FileCount   int                `xml:"file_count"`
			FolderCount int                `xml:"folder_count"`
			TotalSize   int64              `xml:"total_size"`
			Metadata    []MetadataItem     `xml:"metadata>item"`
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

	// 根据内容处理选项处理文件内容
	if f.config != nil && f.config.Formats.XML.Formatting.ContentHandling != "" {
		return f.formatFileWithContentHandling(file, f.config.Formats.XML.Formatting.ContentHandling)
	}

	// 默认使用标准XML序列化
	output, err := xml.MarshalIndent(file, "", "  ")
	if err != nil {
		return "", fmt.Errorf("XML文件格式化失败: %w", err)
	}
	return xml.Header + string(output), nil
}

// formatAIOptimized AI优化的XML格式化
func (f *XMLFormatter) formatAIOptimized(data types.ContextData) (string, error) {
	var result strings.Builder

	// 生成AI摘要
	languages := f.detectLanguages(data.Files)
	summaryGenerator := NewAISummaryGenerator(f.config)
	summary := summaryGenerator.GenerateSummary(data.FileCount, data.TotalSize, languages)
	
	// 写入文件摘要
	result.WriteString(summary.FormatAsXML())
	result.WriteString("\n")

	// 生成目录结构
	directoryStructure := f.generateDirectoryStructure(data.Folders)
	result.WriteString(directoryStructure)
	result.WriteString("\n")

	// 生成文件内容
	result.WriteString("<files>\n")
	for _, file := range data.Files {
		fileXML, err := f.formatFileAIOptimized(file)
		if err != nil {
			return "", fmt.Errorf("格式化文件失败 %s: %w", file.Path, err)
		}
		result.WriteString(fileXML)
		result.WriteString("\n")
	}
	result.WriteString("</files>\n")

	// 添加自定义指令
	if f.config != nil && f.config.Output.AIInstructions.Enabled {
		instructionLoader := NewInstructionLoader(f.config)
		instructions, err := instructionLoader.LoadInstructions()
		if err != nil {
			return "", fmt.Errorf("加载AI指令失败: %w", err)
		}
		if instructions != "" {
			result.WriteString("\n<instruction>\n")
			result.WriteString(fmt.Sprintf("  <![CDATA[\n%s\n  ]]>", instructions))
			result.WriteString("\n</instruction>\n")
		}
	}

	return xml.Header + "<project>\n" + result.String() + "</project>", nil
}

// formatFileAIOptimized AI优化的文件格式化
func (f *XMLFormatter) formatFileAIOptimized(file types.FileInfo) (string, error) {
	// 如果是二进制文件，不显示内容
	content := file.Content
	if file.IsBinary {
		content = "[二进制文件 - 内容未显示]"
	}

	// 计算文件元数据
	lines := len(strings.Split(content, "\n"))
	tokens := f.estimateTokens(content)
	language := f.detectLanguage(file.Path, content)

	fileXML := fmt.Sprintf(`  <file path="%s">
    <metadata>
      <size>%d</size>
      <lines>%d</lines>
      <tokens>%d</tokens>
      <language>%s</language>
    </metadata>
    <content>
      <![CDATA[%s]]>
    </content>
  </file>`,
		escapeXMLAttribute(file.Path),
		file.Size,
		lines,
		tokens,
		language,
		content,
	)

	return fileXML, nil
}

// generateDirectoryStructure 生成目录结构
func (f *XMLFormatter) generateDirectoryStructure(folders []types.FolderInfo) string {
	var result strings.Builder
	result.WriteString("<directory_structure>\n")
	
	for _, folder := range folders {
		f.formatFolderStructure(&result, folder, 1)
	}
	
	result.WriteString("</directory_structure>")
	return result.String()
}

// formatFolderStructure 递归格式化文件夹结构
func (f *XMLFormatter) formatFolderStructure(result *strings.Builder, folder types.FolderInfo, depth int) {
	indent := strings.Repeat("  ", depth)
	result.WriteString(fmt.Sprintf("%s<directory name=\"%s\" path=\"%s\" files=\"%d\" size=\"%d\" />\n",
		indent,
		escapeXMLAttribute(folder.Name),
		escapeXMLAttribute(folder.Path),
		len(folder.Files),
		folder.Size,
	))

	for _, subFolder := range folder.Folders {
		f.formatFolderStructure(result, subFolder, depth+1)
	}
}

// detectLanguages 检测编程语言
func (f *XMLFormatter) detectLanguages(files []types.FileInfo) []string {
	languageMap := make(map[string]bool)
	for _, file := range files {
		lang := f.detectLanguage(file.Path, file.Content)
		if lang != "unknown" {
			languageMap[lang] = true
		}
	}

	languages := make([]string, 0, len(languageMap))
	for lang := range languageMap {
		languages = append(languages, lang)
	}
	return languages
}

// detectLanguage 检测文件语言
func (f *XMLFormatter) detectLanguage(filepath, content string) string {
	ext := strings.ToLower(getFileExtension(filepath))
	switch ext {
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".java":
		return "java"
	case ".cpp", ".cc", ".cxx":
		return "cpp"
	case ".c":
		return "c"
	case ".cs":
		return "csharp"
	case ".php":
		return "php"
	case ".rb":
		return "ruby"
	case ".rs":
		return "rust"
	case ".swift":
		return "swift"
	case ".kt":
		return "kotlin"
	case ".scala":
		return "scala"
	case ".html", ".htm":
		return "html"
	case ".css":
		return "css"
	case ".xml":
		return "xml"
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".md":
		return "markdown"
	case ".sh", ".bash":
		return "shell"
	case ".sql":
		return "sql"
	default:
		return "unknown"
	}
}

// estimateTokens 估算token数量（简化实现）
func (f *XMLFormatter) estimateTokens(content string) int {
	// 简单的估算：假设每个token大约4个字符
	return len(content) / 4
}

// escapeXMLAttribute 转义XML属性值
func escapeXMLAttribute(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// getFileExtension 获取文件扩展名
func getFileExtension(filepath string) string {
	parts := strings.Split(filepath, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
}

// FormatFolder 格式化文件夹
func (f *XMLFormatter) FormatFolder(folder types.FolderInfo) (string, error) {
	// 如果配置了内容处理选项，使用相应的处理方式
	if f.config != nil && f.config.Formats.XML.Formatting.ContentHandling != "" {
		return f.formatFolderWithContentHandling(folder, f.config.Formats.XML.Formatting.ContentHandling)
	}

	// 默认使用标准XML序列化
	output, err := xml.MarshalIndent(folder, "", "  ")
	if err != nil {
		return "", fmt.Errorf("XML文件夹格式化失败: %w", err)
	}
	return xml.Header + string(output), nil
}

// formatFileWithContentHandling 根据内容处理选项格式化文件
func (f *XMLFormatter) formatFileWithContentHandling(file types.FileInfo, contentHandling types.XMLContentHandling) (string, error) {
	switch contentHandling {
	case types.XMLContentCDATA:
		return f.formatFileWithCDATA(file)
	case types.XMLContentRaw:
		return f.formatFileWithRawContent(file)
	default:
		// 默认使用标准XML序列化
		output, err := xml.MarshalIndent(file, "", "  ")
		if err != nil {
			return "", fmt.Errorf("XML文件格式化失败: %w", err)
		}
		return xml.Header + string(output), nil
	}
}

// formatFolderWithContentHandling 根据内容处理选项格式化文件夹
func (f *XMLFormatter) formatFolderWithContentHandling(folder types.FolderInfo, contentHandling types.XMLContentHandling) (string, error) {
	switch contentHandling {
	case types.XMLContentCDATA:
		return f.formatFolderWithCDATA(folder)
	case types.XMLContentRaw:
		return f.formatFolderWithRawContent(folder)
	default:
		// 默认使用标准XML序列化
		output, err := xml.MarshalIndent(folder, "", "  ")
		if err != nil {
			return "", fmt.Errorf("XML文件夹格式化失败: %w", err)
		}
		return xml.Header + string(output), nil
	}
}

// formatFileWithCDATA 使用CDATA包装文件内容
func (f *XMLFormatter) formatFileWithCDATA(file types.FileInfo) (string, error) {
	// 如果是二进制文件，不显示内容
	content := file.Content
	if file.IsBinary {
		content = "[二进制文件 - 内容未显示]"
	}

	// 创建自定义结构体，将内容包装在CDATA中
	fileWithCDATA := struct {
		XMLName  xml.Name  `xml:"File"`
		Name     string    `xml:"Name"`
		Path     string    `xml:"Path"`
		Size     int64     `xml:"Size"`
		ModTime  string    `xml:"ModTime"`
		IsBinary bool      `xml:"IsBinary"`
		IsDir    bool      `xml:"IsDir"`
		IsHidden bool      `xml:"IsHidden"`
		Content  CDATAText `xml:"Content"`
	}{
		Name:     file.Name,
		Path:     file.Path,
		Size:     file.Size,
		ModTime:  file.ModTime.Format("2006-01-02T15:04:05Z"),
		IsBinary: file.IsBinary,
		IsDir:    file.IsDir,
		IsHidden: file.IsHidden,
		Content:  CDATAText{Text: content},
	}

	output, err := xml.MarshalIndent(fileWithCDATA, "", "  ")
	if err != nil {
		return "", fmt.Errorf("XML文件格式化失败: %w", err)
	}
	return xml.Header + string(output), nil
}

// formatFolderWithCDATA 使用CDATA包装文件夹内容
func (f *XMLFormatter) formatFolderWithCDATA(folder types.FolderInfo) (string, error) {
	// 转换文件列表
	files := make([]struct {
		XMLName  xml.Name  `xml:"File"`
		Name     string    `xml:"Name"`
		Path     string    `xml:"Path"`
		Size     int64     `xml:"Size"`
		ModTime  string    `xml:"ModTime"`
		IsBinary bool      `xml:"IsBinary"`
		IsDir    bool      `xml:"IsDir"`
		IsHidden bool      `xml:"IsHidden"`
		Content  CDATAText `xml:"Content"`
	}, len(folder.Files))

	for i, file := range folder.Files {
		// 如果是二进制文件，不显示内容
		content := file.Content
		if file.IsBinary {
			content = "[二进制文件 - 内容未显示]"
		}

		files[i] = struct {
			XMLName  xml.Name  `xml:"File"`
			Name     string    `xml:"Name"`
			Path     string    `xml:"Path"`
			Size     int64     `xml:"Size"`
			ModTime  string    `xml:"ModTime"`
			IsBinary bool      `xml:"IsBinary"`
			IsDir    bool      `xml:"IsDir"`
			IsHidden bool      `xml:"IsHidden"`
			Content  CDATAText `xml:"Content"`
		}{
			XMLName:  xml.Name{Local: "File"},
			Name:     file.Name,
			Path:     file.Path,
			Size:     file.Size,
			ModTime:  file.ModTime.Format("2006-01-02T15:04:05Z"),
			IsBinary: file.IsBinary,
			IsDir:    file.IsDir,
			IsHidden: file.IsHidden,
			Content:  CDATAText{Text: content},
		}
	}

	// 递归转换子文件夹
	subFolders := make([]interface{}, len(folder.Folders))
	for i, subFolder := range folder.Folders {
		subFolderXML, err := f.formatSubFolderWithCDATA(subFolder)
		if err != nil {
			return "", err
		}
		subFolders[i] = subFolderXML
	}

	folderWithCDATA := struct {
		XMLName  xml.Name    `xml:"Folder"`
		Name     string      `xml:"Name"`
		Path     string      `xml:"Path"`
		Size     int64       `xml:"Size"`
		ModTime  string      `xml:"ModTime"`
		IsHidden bool        `xml:"IsHidden"`
		Count    int         `xml:"Count"`
		Files    interface{} `xml:"Files"`
		Folders  interface{} `xml:"Folders"`
	}{
		Name:     folder.Name,
		Path:     folder.Path,
		Size:     folder.Size,
		ModTime:  folder.ModTime.Format(time.RFC3339),
		IsHidden: folder.IsHidden,
		Count:    folder.Count,
		Files:    files,
		Folders:  subFolders,
	}

	output, err := xml.MarshalIndent(folderWithCDATA, "", "  ")
	if err != nil {
		return "", fmt.Errorf("XML文件夹格式化失败: %w", err)
	}
	return xml.Header + string(output), nil
}

// formatSubFolderWithCDATA 递归处理子文件夹
func (f *XMLFormatter) formatSubFolderWithCDATA(folder types.FolderInfo) (interface{}, error) {
	// 转换文件列表
	files := make([]interface{}, len(folder.Files))
	for i, file := range folder.Files {
		content := file.Content
		if file.IsBinary {
			content = "[二进制文件 - 内容未显示]"
		}

		files[i] = struct {
			XMLName  xml.Name  `xml:"File"`
			Name     string    `xml:"Name"`
			Path     string    `xml:"Path"`
			Size     int64     `xml:"Size"`
			ModTime  string    `xml:"ModTime"`
			IsBinary bool      `xml:"IsBinary"`
			IsDir    bool      `xml:"IsDir"`
			IsHidden bool      `xml:"IsHidden"`
			Content  CDATAText `xml:"Content"`
		}{
			Name:     file.Name,
			Path:     file.Path,
			Size:     file.Size,
			ModTime:  file.ModTime.Format("2006-01-02T15:04:05Z"),
			IsBinary: file.IsBinary,
			IsDir:    file.IsDir,
			IsHidden: file.IsHidden,
			Content:  CDATAText{Text: content},
		}
	}

	// 递归转换子文件夹
	subFolders := make([]interface{}, len(folder.Folders))
	for i, subFolder := range folder.Folders {
		subFolderXML, err := f.formatSubFolderWithCDATA(subFolder)
		if err != nil {
			return nil, err
		}
		subFolders[i] = subFolderXML
	}

	return struct {
		XMLName  xml.Name    `xml:"Folder"`
		Name     string      `xml:"Name"`
		Path     string      `xml:"Path"`
		Size     int64       `xml:"Size"`
		ModTime  string      `xml:"ModTime"`
		IsHidden bool        `xml:"IsHidden"`
		Count    int         `xml:"Count"`
		Files    interface{} `xml:"Files"`
		Folders  interface{} `xml:"Folders"`
	}{
		Name:     folder.Name,
		Path:     folder.Path,
		Size:     folder.Size,
		ModTime:  folder.ModTime.Format(time.RFC3339),
		IsHidden: folder.IsHidden,
		Count:    folder.Count,
		Files:    files,
		Folders:  subFolders,
	}, nil
}

// formatFileWithRawContent 使用最小转义格式化文件
func (f *XMLFormatter) formatFileWithRawContent(file types.FileInfo) (string, error) {
	// 如果是二进制文件，不显示内容
	content := file.Content
	if file.IsBinary {
		content = "[二进制文件 - 内容未显示]"
	}

	// 创建自定义结构体，最小化转义
	fileWithRaw := struct {
		XMLName  xml.Name `xml:"File"`
		Name     string   `xml:"Name"`
		Path     string   `xml:"Path"`
		Size     int64    `xml:"Size"`
		ModTime  string   `xml:"ModTime"`
		IsBinary bool     `xml:"IsBinary"`
		IsDir    bool     `xml:"IsDir"`
		IsHidden bool     `xml:"IsHidden"`
		Content  RawText  `xml:"Content"`
	}{
		Name:     file.Name,
		Path:     file.Path,
		Size:     file.Size,
		ModTime:  file.ModTime.Format("2006-01-02T15:04:05Z"),
		IsBinary: file.IsBinary,
		IsDir:    file.IsDir,
		IsHidden: file.IsHidden,
		Content:  RawText{Text: content},
	}

	output, err := xml.MarshalIndent(fileWithRaw, "", "  ")
	if err != nil {
		return "", fmt.Errorf("XML文件格式化失败: %w", err)
	}
	return xml.Header + string(output), nil
}

// formatFolderWithRawContent 使用最小转义格式化文件夹
func (f *XMLFormatter) formatFolderWithRawContent(folder types.FolderInfo) (string, error) {
	// 转换文件列表
	files := make([]struct {
		XMLName  xml.Name `xml:"File"`
		Name     string   `xml:"Name"`
		Path     string   `xml:"Path"`
		Size     int64    `xml:"Size"`
		ModTime  string   `xml:"ModTime"`
		IsBinary bool     `xml:"IsBinary"`
		IsDir    bool     `xml:"IsDir"`
		IsHidden bool     `xml:"IsHidden"`
		Content  RawText  `xml:"Content"`
	}, len(folder.Files))

	for i, file := range folder.Files {
		// 如果是二进制文件，不显示内容
		content := file.Content
		if file.IsBinary {
			content = "[二进制文件 - 内容未显示]"
		}

		files[i] = struct {
			XMLName  xml.Name `xml:"File"`
			Name     string   `xml:"Name"`
			Path     string   `xml:"Path"`
			Size     int64    `xml:"Size"`
			ModTime  string   `xml:"ModTime"`
			IsBinary bool     `xml:"IsBinary"`
			IsDir    bool     `xml:"IsDir"`
			IsHidden bool     `xml:"IsHidden"`
			Content  RawText  `xml:"Content"`
		}{
			Name:     file.Name,
			Path:     file.Path,
			Size:     file.Size,
			ModTime:  file.ModTime.Format("2006-01-02T15:04:05Z"),
			IsBinary: file.IsBinary,
			IsDir:    file.IsDir,
			IsHidden: file.IsHidden,
			Content:  RawText{Text: content},
		}
	}

	// 递归转换子文件夹
	subFolders := make([]interface{}, len(folder.Folders))
	for i, subFolder := range folder.Folders {
		subFolderXML, err := f.formatSubFolderWithRawContent(subFolder)
		if err != nil {
			return "", err
		}
		subFolders[i] = subFolderXML
	}

	folderWithRaw := struct {
		XMLName  xml.Name    `xml:"Folder"`
		Name     string      `xml:"Name"`
		Path     string      `xml:"Path"`
		Size     int64       `xml:"Size"`
		ModTime  string      `xml:"ModTime"`
		IsHidden bool        `xml:"IsHidden"`
		Count    int         `xml:"Count"`
		Files    interface{} `xml:"Files"`
		Folders  interface{} `xml:"Folders"`
	}{
		Name:     folder.Name,
		Path:     folder.Path,
		Size:     folder.Size,
		ModTime:  folder.ModTime.Format("2006-01-02T15:04:05Z"),
		IsHidden: folder.IsHidden,
		Count:    folder.Count,
		Files:    files,
		Folders:  subFolders,
	}

	output, err := xml.MarshalIndent(folderWithRaw, "", "  ")
	if err != nil {
		return "", fmt.Errorf("XML文件夹格式化失败: %w", err)
	}
	return xml.Header + string(output), nil
}

// formatSubFolderWithRawContent 递归处理子文件夹
func (f *XMLFormatter) formatSubFolderWithRawContent(folder types.FolderInfo) (interface{}, error) {
	// 转换文件列表
	files := make([]interface{}, len(folder.Files))
	for i, file := range folder.Files {
		content := file.Content
		if file.IsBinary {
			content = "[二进制文件 - 内容未显示]"
		}

		files[i] = struct {
			XMLName  xml.Name `xml:"File"`
			Name     string   `xml:"Name"`
			Path     string   `xml:"Path"`
			Size     int64    `xml:"Size"`
			ModTime  string   `xml:"ModTime"`
			IsBinary bool     `xml:"IsBinary"`
			IsDir    bool     `xml:"IsDir"`
			IsHidden bool     `xml:"IsHidden"`
			Content  RawText  `xml:"Content"`
		}{
			Name:     file.Name,
			Path:     file.Path,
			Size:     file.Size,
			ModTime:  file.ModTime.Format("2006-01-02T15:04:05Z"),
			IsBinary: file.IsBinary,
			IsDir:    file.IsDir,
			IsHidden: file.IsHidden,
			Content:  RawText{Text: content},
		}
	}

	// 递归转换子文件夹
	subFolders := make([]interface{}, len(folder.Folders))
	for i, subFolder := range folder.Folders {
		subFolderXML, err := f.formatSubFolderWithRawContent(subFolder)
		if err != nil {
			return nil, err
		}
		subFolders[i] = subFolderXML
	}

	return struct {
		XMLName  xml.Name    `xml:"Folder"`
		Name     string      `xml:"Name"`
		Path     string      `xml:"Path"`
		Size     int64       `xml:"Size"`
		ModTime  string      `xml:"ModTime"`
		IsHidden bool        `xml:"IsHidden"`
		Count    int         `xml:"Count"`
		Files    interface{} `xml:"Files"`
		Folders  interface{} `xml:"Folders"`
	}{
		Name:     folder.Name,
		Path:     folder.Path,
		Size:     folder.Size,
		ModTime:  folder.ModTime.Format("2006-01-02T15:04:05Z"),
		IsHidden: folder.IsHidden,
		Count:    folder.Count,
		Files:    files,
		Folders:  subFolders,
	}, nil
}
