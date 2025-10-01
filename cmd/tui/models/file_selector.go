package models

import (
	"fmt"
	"sort"
	"strings"

	"code-context-generator/internal/selector"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FileSelectorModel 文件选择器模型
type FileSelectorModel struct {
	path         string
	items        []selector.FileItem
	selected     map[int]bool
	cursor       int
	scrollOffset int
	multiSelect  bool
	filter       string
	height       int
	width        int
}

// NewFileSelectorModel 创建文件选择器模型
func NewFileSelectorModel(path string) *FileSelectorModel {
	return &FileSelectorModel{
		path:         path,
		items:        []selector.FileItem{},
		selected:     make(map[int]bool),
		cursor:       0,
		scrollOffset: 0,
		multiSelect:  true,
		filter:       "",
	}
}

// Init 初始化
func (m *FileSelectorModel) Init() tea.Cmd {
	return m.loadFiles()
}

// Update 更新模型
func (m *FileSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg {
				return FileSelectionMsg{Selected: []string{}}
			}
		case "enter":
			return m, m.confirmSelection()
		case "up", "k":
			m.moveCursor(-1)
		case "down", "j":
			m.moveCursor(1)
		case " ":
			if m.multiSelect {
				m.toggleSelection()
			}
		case "a":
			if m.multiSelect {
				m.selectAll()
			}
		case "n":
			if m.multiSelect {
				m.selectNone()
			}
		case "i":
			m.invertSelection()
		case "/":
			// 进入搜索模式
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateViewport()
	case FileListMsg:
		m.items = msg.Items
		m.updateViewport()
	}
	return m, nil
}

// View 渲染视图
func (m *FileSelectorModel) View() string {
	if len(m.items) == 0 {
		return "正在加载文件列表..."
	}

	var content strings.Builder

	// 标题
	content.WriteString(TitleStyle.Render("文件选择器"))
	content.WriteString("\n\n")

	// 路径
	content.WriteString(NormalStyle.Render(fmt.Sprintf("路径: %s", m.path)))
	content.WriteString("\n\n")

	// 文件列表
	visibleItems := m.getVisibleItems()
	for i, item := range visibleItems {
		actualIndex := m.scrollOffset + i
		isSelected := m.selected[actualIndex]
		isCursor := actualIndex == m.cursor

		line := m.renderFileItem(item, isSelected, isCursor)
		content.WriteString(line)
		content.WriteString("\n")
	}

	// 帮助信息
	content.WriteString("\n")
	content.WriteString(HelpStyle.Render("操作: ↑↓移动, 空格选择, Enter确认, Esc取消, a全选, n取消全选, i反选"))

	return content.String()
}

// 辅助方法
func (m *FileSelectorModel) moveCursor(direction int) {
	m.cursor += direction
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	m.updateScroll()
}

func (m *FileSelectorModel) toggleSelection() {
	m.selected[m.cursor] = !m.selected[m.cursor]
}

func (m *FileSelectorModel) selectAll() {
	for i := range m.items {
		m.selected[i] = true
	}
}

func (m *FileSelectorModel) selectNone() {
	for i := range m.items {
		m.selected[i] = false
	}
}

func (m *FileSelectorModel) invertSelection() {
	for i := range m.items {
		m.selected[i] = !m.selected[i]
	}
}

func (m *FileSelectorModel) updateScroll() {
	// 确保最小高度为10，避免负数
	if m.height <= 10 {
		return
	}
	
	visibleHeight := m.height - 10 // 减去标题和帮助信息的高度
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}
	if m.cursor >= m.scrollOffset+visibleHeight {
		m.scrollOffset = m.cursor - visibleHeight + 1
	}
}

func (m *FileSelectorModel) getVisibleItems() []selector.FileItem {
	// 确保最小高度为10，避免负数
	if m.height <= 10 {
		return m.items
	}
	
	visibleHeight := m.height - 10
	start := m.scrollOffset
	end := start + visibleHeight
	if end > len(m.items) {
		end = len(m.items)
	}
	
	// 确保start不会超出范围
	if start >= len(m.items) {
		start = len(m.items)
	}
	if start < 0 {
		start = 0
	}
	
	// 确保end不小于start
	if end < start {
		end = start
	}
	
	return m.items[start:end]
}

func (m *FileSelectorModel) updateViewport() {
	// 更新视口大小
}

func (m *FileSelectorModel) loadFiles() tea.Cmd {
	return func() tea.Msg {
		// 获取目录内容
		contents, err := selector.GetDirectoryContents(m.path, GetConfig().FileProcessing.IncludeHidden)
		if err != nil {
			// 如果出错，返回空列表
			return FileListMsg{Items: []selector.FileItem{}}
		}

		// 将FileInfo转换为FileItem
		items := make([]selector.FileItem, 0, len(contents))
		for _, info := range contents {
			item := selector.FileItem{
				Path:     info.Path,
				Name:     info.Name,
				Size:     info.Size,
				ModTime:  info.ModTime,
				IsDir:    info.IsDir,
				IsHidden: info.IsHidden,
				Icon:     info.Icon,
				Type:     info.Type,
				Selected: false,
			}
			items = append(items, item)
		}

		// 按名称排序
		sort.Slice(items, func(i, j int) bool {
			// 目录优先，然后按名称排序
			if items[i].IsDir != items[j].IsDir {
				return items[i].IsDir
			}
			return items[i].Name < items[j].Name
		})

		return FileListMsg{Items: items}
	}
}

func (m *FileSelectorModel) confirmSelection() tea.Cmd {
	return func() tea.Msg {
		var selected []string
		for i, item := range m.items {
			if m.selected[i] {
				selected = append(selected, item.Path)
			}
		}
		return FileSelectionMsg{Selected: selected}
	}
}

func (m *FileSelectorModel) renderFileItem(item selector.FileItem, isSelected, isCursor bool) string {
	var style lipgloss.Style

	if isCursor {
		style = SelectedStyle
	} else if item.IsDir {
		// 目录使用特殊的样式
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00BFFF")). // 深蓝色
			Bold(true)
	} else {
		style = NormalStyle
	}

	prefix := "  "
	if isSelected {
		prefix = "✓ "
	}

	icon := getFileIcon(item.Name, item.IsDir)

	name := item.Name
	if item.IsDir {
		name += "/"
	}

	line := fmt.Sprintf("%s%s %s", prefix, icon, name)
	return style.Render(line)
}