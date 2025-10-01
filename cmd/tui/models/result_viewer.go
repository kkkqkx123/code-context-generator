package models

import (
	"code-context-generator/pkg/types"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ResultViewerModel 结果查看器模型
type ResultViewerModel struct {
	result       *types.WalkResult
	scrollOffset int
	width        int
	height       int
	currentTab   int
}

// NewResultViewerModel 创建结果查看器模型
func NewResultViewerModel() *ResultViewerModel {
	return &ResultViewerModel{
		currentTab: 0,
	}
}

// Init 初始化
func (m *ResultViewerModel) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m *ResultViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg {
				return ConfigUpdateMsg{Config: GetConfig()}
			}
		case "tab":
			m.currentTab = (m.currentTab + 1) % 3 // 假设有3个标签页
		case "up", "k":
			m.scrollOffset--
			if m.scrollOffset < 0 {
				m.scrollOffset = 0
			}
		case "down", "j":
			m.scrollOffset++
		case "b": // b键返回主界面
			return m, func() tea.Msg {
				return ConfigUpdateMsg{Config: GetConfig()}
			}
		case "s":
			// 保存结果
			return m, m.saveResult()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View 渲染视图
func (m *ResultViewerModel) View() string {
	if m.result == nil {
		return "没有结果可显示"
	}

	var content strings.Builder

	// 标题
	content.WriteString(TitleStyle.Render("扫描结果"))
	content.WriteString("\n\n")

	// 标签页
	tabs := []string{"概览", "文件", "目录"}
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
	case 0: // 概览
		content.WriteString(m.renderOverview())
	case 1: // 文件
		content.WriteString(m.renderFiles())
	case 2: // 目录
		content.WriteString(m.renderDirectories())
	}

	// 帮助信息
	content.WriteString("\n")
	content.WriteString(HelpStyle.Render("操作: Tab切换标签, ↑↓滚动, b返回主界面, s保存, ESC返回主界面"))

	return content.String()
}

// SetResult 设置结果
func (m *ResultViewerModel) SetResult(result *types.WalkResult) {
	m.result = result
}

// GetCurrentTab 获取当前标签页（用于测试）
func (m *ResultViewerModel) GetCurrentTab() int {
	return m.currentTab
}

// 辅助方法
func (m *ResultViewerModel) renderOverview() string {
	var content strings.Builder

	content.WriteString(NormalStyle.Render(fmt.Sprintf("根路径: %s\n", m.result.RootPath)))
		content.WriteString(NormalStyle.Render(fmt.Sprintf("文件数量: %d\n", m.result.FileCount)))
		content.WriteString(NormalStyle.Render(fmt.Sprintf("目录数量: %d\n", m.result.FolderCount)))
		content.WriteString(NormalStyle.Render(fmt.Sprintf("总大小: %s\n", formatFileSize(m.result.TotalSize))))
		content.WriteString(NormalStyle.Render(fmt.Sprintf("扫描时间: %v\n", m.result.ScanDuration)))

	return content.String()
}

func (m *ResultViewerModel) renderFiles() string {
	var content strings.Builder

	start := m.scrollOffset
	end := start + m.height - 10
	if end > len(m.result.Files) {
		end = len(m.result.Files)
	}

	for i := start; i < end; i++ {
		file := m.result.Files[i]
		icon := getFileIcon(file.Name, false) // 文件不是目录
		content.WriteString(NormalStyle.Render(fmt.Sprintf("%s %s (%s)\n", icon, file.Name, formatFileSize(file.Size))))
	}

	return content.String()
}

func (m *ResultViewerModel) renderDirectories() string {
	var content strings.Builder

	start := m.scrollOffset
	end := start + m.height - 10
	if end > len(m.result.Folders) {
		end = len(m.result.Folders)
	}

	for i := start; i < end; i++ {
		folder := m.result.Folders[i]
		content.WriteString(NormalStyle.Render(fmt.Sprintf("📂 %s/\n", folder.Name)))
	}

	return content.String()
}

func (m *ResultViewerModel) saveResult() tea.Cmd {
	return func() tea.Msg {
		// 这里应该实现保存逻辑
		return nil
	}
}
