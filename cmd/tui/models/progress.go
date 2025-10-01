package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ProgressModel 进度条模型
type ProgressModel struct {
	progress float64
	status   string
	width    int
	height   int
}

// NewProgressModel 创建进度条模型
func NewProgressModel() *ProgressModel {
	return &ProgressModel{
		progress: 0,
		status:   "准备中...",
	}
}

// Init 初始化
func (m *ProgressModel) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m *ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case ProgressMsg:
		m.progress = msg.Progress
		m.status = msg.Status
	}
	return m, nil
}

// View 渲染视图
func (m *ProgressModel) View() string {
	var content strings.Builder

	content.WriteString(TitleStyle.Render("处理中..."))
	content.WriteString("\n\n")

	// 进度条
	barWidth := m.width - 4
	if barWidth > 0 {
		filled := int(float64(barWidth) * m.progress)
		empty := barWidth - filled

		bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
		content.WriteString(NormalStyle.Render(fmt.Sprintf("[%s] %.1f%%", bar, m.progress*100)))
		content.WriteString("\n\n")
	}

	// 状态信息
	content.WriteString(NormalStyle.Render(m.status))
	content.WriteString("\n\n")

	// 帮助信息
	content.WriteString(HelpStyle.Render("操作: Ctrl+C 取消"))

	return content.String()
}

// SetProgress 设置进度
func (m *ProgressModel) SetProgress(progress float64) {
	m.progress = progress
}

// SetStatus 设置状态
func (m *ProgressModel) SetStatus(status string) {
	m.status = status
}