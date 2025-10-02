package models

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

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
				return &FileSelectionMsg{Selected: []string{}}
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
		var content strings.Builder
		content.WriteString(TitleStyle.Render("文件选择器"))
		content.WriteString("\n\n")
		content.WriteString(NormalStyle.Render("正在加载文件列表..."))
		content.WriteString("\n\n")
		content.WriteString(HelpStyle.Render("操作: Esc返回主界面"))
		return content.String()
	}

	var content strings.Builder

	// 标题
	content.WriteString(TitleStyle.Render("文件选择器"))
	content.WriteString("\n")

	// 当前路径
	content.WriteString(NormalStyle.Render(fmt.Sprintf("当前路径: %s", m.path)))
	content.WriteString("\n")

	// 分页信息
	visibleHeight := m.height - 6
	if visibleHeight < 3 {
		visibleHeight = 3
	}
	if visibleHeight > 20 {
		visibleHeight = 20
	}
	
	startItem := m.scrollOffset + 1
	endItem := m.scrollOffset + visibleHeight
	if endItem > len(m.items) {
		endItem = len(m.items)
	}
	
	content.WriteString(NormalStyle.Render(fmt.Sprintf("显示: %d-%d / %d 个项目", startItem, endItem, len(m.items))))
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
	if len(m.items) == 0 {
		return
	}
	
	m.cursor += direction
	if m.cursor < 0 {
		m.cursor = len(m.items) - 1 // 循环到末尾
	}
	if m.cursor >= len(m.items) {
		m.cursor = 0 // 循环到开头
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
	// 如果项目为空，不需要滚动
	if len(m.items) == 0 {
		return
	}
	
	// 确保光标在有效范围内
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	
	// 计算可见区域高度（与getVisibleItems保持一致）
	visibleHeight := m.height - 6
	if visibleHeight < 3 {
		visibleHeight = 3
	}
	if visibleHeight > 20 {
		visibleHeight = 20
	}
	
	// 确保滚动偏移在有效范围内
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
	maxScrollOffset := len(m.items) - visibleHeight
	if maxScrollOffset < 0 {
		maxScrollOffset = 0
	}
	if m.scrollOffset > maxScrollOffset {
		m.scrollOffset = maxScrollOffset
	}
	
	// 调整滚动位置以保持光标可见
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}
	if m.cursor >= m.scrollOffset+visibleHeight {
		m.scrollOffset = m.cursor - visibleHeight + 1
	}
}

func (m *FileSelectorModel) getVisibleItems() []selector.FileItem {
	// 如果项目为空，返回空切片
	if len(m.items) == 0 {
		return []selector.FileItem{}
	}
	
	// 计算可见区域高度（减去标题、路径和帮助信息）
	// 标题1行 + 路径1行 + 分页信息1行 + 文件列表 + 帮助信息1行 + 边距2行
	visibleHeight := m.height - 6
	if visibleHeight < 3 {
		visibleHeight = 3 // 最小显示3个项目
	}
	if visibleHeight > 20 {
		visibleHeight = 20 // 最大显示20个项目，实现分页
	}
	
	start := m.scrollOffset
	end := start + visibleHeight
	if end > len(m.items) {
		end = len(m.items)
	}
	
	// 确保start不会超出范围
	if start >= len(m.items) {
		start = len(m.items) - visibleHeight
		if start < 0 {
			start = 0
		}
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
		// 创建超时上下文，防止文件系统操作卡死
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		// 使用goroutine处理文件加载，支持超时
		type result struct {
			items []selector.FileItem
			err   error
		}
		
		resultChan := make(chan result, 1)
		
		go func() {
			// 获取配置
			config := GetConfig()
			showHidden := false
			if config != nil {
				showHidden = config.FileProcessing.IncludeHidden
			}
			
			// 检查路径是否存在
			if _, err := os.Stat(m.path); err != nil {
				resultChan <- result{items: []selector.FileItem{}, err: fmt.Errorf("路径不存在: %s，错误: %v", m.path, err)}
				return
			}
			
			// 获取目录内容
			contents, err := selector.GetDirectoryContents(m.path, showHidden)
			if err != nil {
				resultChan <- result{items: []selector.FileItem{}, err: err}
				return
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
			
			resultChan <- result{items: items, err: nil}
		}()
		
		// 等待结果或超时
		select {
		case res := <-resultChan:
			if res.err != nil {
				return FileListMsg{Items: []selector.FileItem{}}
			}
			
			// 按名称排序
			items := res.items
			sort.Slice(items, func(i, j int) bool {
				// 目录优先，然后按名称排序
				if items[i].IsDir != items[j].IsDir {
					return items[i].IsDir
				}
				return items[i].Name < items[j].Name
			})
			
			return FileListMsg{Items: items}
		case <-ctx.Done():
			// 超时，返回空列表
			return FileListMsg{Items: []selector.FileItem{}}
		}
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
		return &FileSelectionMsg{Selected: selected}
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

	// 显示完整路径（相对于当前目录）
	relPath := item.Path
	if strings.HasPrefix(item.Path, m.path) {
		relPath = strings.TrimPrefix(item.Path, m.path)
		if strings.HasPrefix(relPath, "/") || strings.HasPrefix(relPath, "\\") {
			relPath = relPath[1:]
		}
	}
	if relPath == "" {
		relPath = name
	} else {
		relPath = name
	}

	line := fmt.Sprintf("%s%s %s", prefix, icon, relPath)
	return style.Render(line)
}