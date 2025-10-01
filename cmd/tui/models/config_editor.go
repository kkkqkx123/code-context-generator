package models

import (
	"fmt"
	"strings"

	"code-context-generator/pkg/types"
	tea "github.com/charmbracelet/bubbletea"
)

// ConfigEditorModel 配置编辑器模型
type ConfigEditorModel struct {
	config     *types.Config
	currentTab int
	width      int
	height     int
	focus      int
}

// NewConfigEditorModel 创建配置编辑器模型
func NewConfigEditorModel(config *types.Config) *ConfigEditorModel {
	return &ConfigEditorModel{
		config: config,
		focus:  0,
	}
}

// Init 初始化
func (m *ConfigEditorModel) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m *ConfigEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg {
				return ConfigUpdateMsg{Config: m.config}
			}
		case "tab":
			m.currentTab = (m.currentTab + 1) % 4 // 假设有4个配置标签页
		case "up", "k":
			m.focus--
			if m.focus < 0 {
				m.focus = 0
			}
		case "down", "j":
			m.focus++
		case "enter":
			// 编辑当前项
			return m, nil
		case "s":
			// 保存配置
			return m, m.saveConfig()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View 渲染视图
func (m *ConfigEditorModel) View() string {
	var content strings.Builder

	// 标题
	content.WriteString(TitleStyle.Render("配置编辑器"))
	content.WriteString("\n\n")

	// 标签页
	tabs := []string{"输出", "文件处理", "UI", "性能"}
	for i, tab := range tabs {
		if i == m.currentTab {
			content.WriteString(SelectedStyle.Render(fmt.Sprintf("[%s]", tab)))
		} else {
			content.WriteString(NormalStyle.Render(fmt.Sprintf(" %s ", tab)))
		}
		content.WriteString(" ")
	}
	content.WriteString("\n\n")

	// 内容
	switch m.currentTab {
	case 0: // 输出
		content.WriteString(m.renderOutputConfig())
	case 1: // 文件处理
		content.WriteString(m.renderFileProcessingConfig())
	case 2: // UI
		content.WriteString(m.renderUIConfig())
	case 3: // 性能
		content.WriteString(m.renderPerformanceConfig())
	}

	// 帮助信息
	content.WriteString("\n")
	content.WriteString(HelpStyle.Render("操作: Tab切换标签, ↑↓选择, Enter编辑, s保存, ESC返回主界面"))

	return content.String()
}

// 辅助方法
func (m *ConfigEditorModel) renderOutputConfig() string {
	var content strings.Builder

	content.WriteString(NormalStyle.Render(fmt.Sprintf("默认格式: %s\n", m.config.Output.DefaultFormat)))
	content.WriteString(NormalStyle.Render(fmt.Sprintf("输出目录: %s\n", m.config.Output.OutputDir)))
	content.WriteString(NormalStyle.Render(fmt.Sprintf("文件名模板: %s\n", m.config.Output.FilenameTemplate)))
	content.WriteString(NormalStyle.Render(fmt.Sprintf("时间戳格式: %s\n", m.config.Output.TimestampFormat)))

	return content.String()
}

func (m *ConfigEditorModel) renderFileProcessingConfig() string {
	var content strings.Builder

	content.WriteString(NormalStyle.Render(fmt.Sprintf("最大文件大小: %s\n", m.config.Filters.MaxFileSize)))
	content.WriteString(NormalStyle.Render(fmt.Sprintf("最大深度: %d\n", m.config.Filters.MaxDepth)))
	content.WriteString(NormalStyle.Render(fmt.Sprintf("跟随符号链接: %v\n", m.config.Filters.FollowSymlinks)))
	content.WriteString(NormalStyle.Render(fmt.Sprintf("排除二进制文件: %v\n", m.config.Filters.ExcludeBinary)))

	return content.String()
}

func (m *ConfigEditorModel) renderUIConfig() string {
	var content strings.Builder

	content.WriteString(NormalStyle.Render(fmt.Sprintf("主题: %s\n", m.config.UI.Theme)))
		content.WriteString(NormalStyle.Render(fmt.Sprintf("显示进度: %v\n", m.config.UI.ShowProgress)))
		content.WriteString(NormalStyle.Render(fmt.Sprintf("显示大小: %v\n", m.config.UI.ShowSize)))
		content.WriteString(NormalStyle.Render(fmt.Sprintf("显示日期: %v\n", m.config.UI.ShowDate)))
		content.WriteString(NormalStyle.Render(fmt.Sprintf("显示预览: %v\n", m.config.UI.ShowPreview)))

	return content.String()
}

func (m *ConfigEditorModel) renderPerformanceConfig() string {
	var content strings.Builder

	content.WriteString(NormalStyle.Render(fmt.Sprintf("最大工作线程: %d\n", m.config.Performance.MaxWorkers)))
		content.WriteString(NormalStyle.Render(fmt.Sprintf("缓冲区大小: %d\n", m.config.Performance.BufferSize)))
		content.WriteString(NormalStyle.Render(fmt.Sprintf("缓存启用: %v\n", m.config.Performance.CacheEnabled)))
		content.WriteString(NormalStyle.Render(fmt.Sprintf("缓存大小: %d\n", m.config.Performance.CacheSize)))

	return content.String()
}

func (m *ConfigEditorModel) saveConfig() tea.Cmd {
	return func() tea.Msg {
		// 这里应该实现保存配置逻辑
		return ConfigUpdateMsg{Config: m.config}
	}
}

// GetCurrentTab 获取当前标签页（用于测试）
func (m *ConfigEditorModel) GetCurrentTab() int {
	return m.currentTab
}