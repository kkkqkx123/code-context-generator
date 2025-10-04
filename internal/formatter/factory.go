package formatter

import (
	"fmt"
	"strings"

	"code-context-generator/pkg/types"
)

// FormatterFactory 格式转换器工厂
type FormatterFactory struct {
	formatters map[string]Formatter
}

// NewFormatterFactory 创建格式转换器工厂
func NewFormatterFactory() *FormatterFactory {
	factory := &FormatterFactory{
		formatters: make(map[string]Formatter),
	}
	
	// 注册所有格式转换器
	factory.RegisterFormatters()
	
	return factory
}

// RegisterFormatters 注册所有格式转换器
func (f *FormatterFactory) RegisterFormatters() {
	// 注册JSON格式转换器
	f.Register("json", NewJSONFormatter(nil))
	
	// 注册XML格式转换器
	f.Register("xml", NewXMLFormatter(nil))
	
	// 注册TOML格式转换器
	f.Register("toml", NewTOMLFormatter(nil))
	
	// 注册Markdown格式转换器
	f.Register("markdown", NewMarkdownFormatter(nil))
}

// Register 注册格式转换器
func (f *FormatterFactory) Register(name string, formatter Formatter) {
	f.formatters[name] = formatter
}

// Get 获取格式转换器（兼容旧接口）
func (f *FormatterFactory) Get(name string) (Formatter, error) {
	return f.GetFormatter(name, nil)
}

// GetFormatter 获取格式转换器
func (f *FormatterFactory) GetFormatter(name string, config *types.Config) (Formatter, error) {
	// 大小写不敏感处理
	name = strings.ToLower(name)
	
	// 根据名称创建对应的格式转换器实例
	switch name {
	case "json":
		return NewJSONFormatter(config), nil
	case "xml":
		return NewXMLFormatter(config), nil
	case "toml":
		return NewTOMLFormatter(config), nil
	case "markdown", "md":
		return NewMarkdownFormatter(config), nil
	default:
		return nil, fmt.Errorf("不支持的格式: %s", name)
	}
}

// GetSupportedFormats 获取所有支持的格式（兼容旧接口）
func (f *FormatterFactory) GetSupportedFormats() []string {
	return f.GetAvailableFormats()
}

// GetAvailableFormats 获取所有可用的格式
func (f *FormatterFactory) GetAvailableFormats() []string {
	formats := make([]string, 0, len(f.formatters))
	for name := range f.formatters {
		formats = append(formats, name)
	}
	return formats
}

// GetFormatterInfo 获取格式转换器信息
func (f *FormatterFactory) GetFormatterInfo(name string) (string, string, error) {
	formatter, err := f.GetFormatter(name, nil)
	if err != nil {
		return "", "", err
	}
	return formatter.GetName(), formatter.GetDescription(), nil
}

// NewFormatter 创建格式转换器（兼容旧接口）
func NewFormatter(format string, config *types.Config) (Formatter, error) {
	factory := NewFormatterFactory()
	return factory.GetFormatter(format, config)
}

// CreateDefaultFactory 创建默认工厂（兼容旧接口）
func CreateDefaultFactory(config *types.Config) *FormatterFactory {
	factory := NewFormatterFactory()
	return factory
}