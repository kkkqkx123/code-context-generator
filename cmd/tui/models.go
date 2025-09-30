// Package main TUI模型定义
package main

import (
	"fmt"
	"strings"

	"code-context-generator/internal/selector"
	"code-context-generator/pkg/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// 消息类型定义

type ProgressMsg struct {
	Progress float64
	Status   string
}

type ResultMsg struct {
	Result *types.WalkResult
}

type ErrorMsg struct {
	Err error
}

type FileSelectionMsg struct {
	Selected []string
}

type ConfigUpdateMsg struct {
	Config *types.Config
}

// FileSelectorModel 文件选择器模型
type FileSelectorModel struct {
	path          string
	items         []selector.FileItem
	selected      map[int]bool
	cursor        int
	scrollOffset  int
	multiSelect   bool
	filter        string
	height        int
	width         int
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
	content.WriteString(titleStyle.Render("文件选择器"))
	content.WriteString("\n\n")
	
	// 路径
	content.WriteString(normalStyle.Render(fmt.Sprintf("路径: %s", m.path)))
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
	content.WriteString(helpStyle.Render("操作: ↑↓移动, 空格选择, Enter确认, Esc取消, a全选, n取消全选, i反选"))
	
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
	visibleHeight := m.height - 10 // 减去标题和帮助信息的高度
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}
	if m.cursor >= m.scrollOffset+visibleHeight {
		m.scrollOffset = m.cursor - visibleHeight + 1
	}
}

func (m *FileSelectorModel) getVisibleItems() []selector.FileItem {
	visibleHeight := m.height - 10
	start := m.scrollOffset
	end := start + visibleHeight
	if end > len(m.items) {
		end = len(m.items)
	}
	return m.items[start:end]
}

func (m *FileSelectorModel) renderFileItem(item selector.FileItem, isSelected, isCursor bool) string {
	var style lipgloss.Style
	
	if isCursor {
		style = selectedStyle
	} else {
		style = normalStyle
	}
	
	prefix := "  "
	if isSelected {
		prefix = "✓ "
	}
	
	icon := "📄"
	if item.IsDir {
		icon = "📁"
	}
	
	name := item.Name
	if item.IsDir {
		name += "/"
	}
	
	line := fmt.Sprintf("%s%s %s", prefix, icon, name)
	return style.Render(line)
}

func (m *FileSelectorModel) updateViewport() {
	// 更新视口大小
}

func (m *FileSelectorModel) loadFiles() tea.Cmd {
	return func() tea.Msg {
		// 这里应该实际加载文件列表
		// 为了简化，返回空列表
		return FileListMsg{Items: []selector.FileItem{}}
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
	
	content.WriteString(titleStyle.Render("处理中..."))
	content.WriteString("\n\n")
	
	// 进度条
	barWidth := m.width - 4
	if barWidth > 0 {
		filled := int(float64(barWidth) * m.progress)
		empty := barWidth - filled
		
		bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
		content.WriteString(normalStyle.Render(fmt.Sprintf("[%s] %.1f%%", bar, m.progress*100)))
		content.WriteString("\n\n")
	}
	
	// 状态信息
	content.WriteString(normalStyle.Render(m.status))
	content.WriteString("\n\n")
	
	// 帮助信息
	content.WriteString(helpStyle.Render("操作: Ctrl+C 取消"))
	
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
				return ConfigUpdateMsg{Config: cfg}
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
				return ConfigUpdateMsg{Config: cfg}
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
	content.WriteString(titleStyle.Render("扫描结果"))
	content.WriteString("\n\n")
	
	// 标签页
	tabs := []string{"概览", "文件", "目录"}
	for i, tab := range tabs {
		if i == m.currentTab {
			content.WriteString(selectedStyle.Render(fmt.Sprintf("[%s]", tab)))
		} else {
			content.WriteString(normalStyle.Render(fmt.Sprintf(" %s ", tab)))
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
	content.WriteString(helpStyle.Render("操作: Tab切换标签, ↑↓滚动, b返回主界面, s保存, ESC返回主界面"))
	
	return content.String()
}

// SetResult 设置结果
func (m *ResultViewerModel) SetResult(result *types.WalkResult) {
	m.result = result
}

// 辅助方法
func (m *ResultViewerModel) renderOverview() string {
	var content strings.Builder
	
	content.WriteString(normalStyle.Render(fmt.Sprintf("根路径: %s\n", m.result.RootPath)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("文件数量: %d\n", m.result.FileCount)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("目录数量: %d\n", m.result.FolderCount)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("总大小: %d 字节\n", m.result.TotalSize)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("扫描时间: %v\n", m.result.ScanDuration)))
	
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
		content.WriteString(normalStyle.Render(fmt.Sprintf("%s (%d bytes)\n", file.Name, file.Size)))
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
		content.WriteString(normalStyle.Render(fmt.Sprintf("%s/\n", folder.Name)))
	}
	
	return content.String()
}

func (m *ResultViewerModel) saveResult() tea.Cmd {
	return func() tea.Msg {
		// 这里应该实现保存逻辑
		return nil
	}
}

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
	content.WriteString(titleStyle.Render("配置编辑器"))
	content.WriteString("\n\n")
	
	// 标签页
	tabs := []string{"输出", "文件处理", "UI", "性能"}
	for i, tab := range tabs {
		if i == m.currentTab {
			content.WriteString(selectedStyle.Render(fmt.Sprintf("[%s]", tab)))
		} else {
			content.WriteString(normalStyle.Render(fmt.Sprintf(" %s ", tab)))
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
	content.WriteString(helpStyle.Render("操作: Tab切换标签, ↑↓选择, Enter编辑, s保存, ESC返回主界面"))
	
	return content.String()
}

// 辅助方法
func (m *ConfigEditorModel) renderOutputConfig() string {
	var content strings.Builder
	
	content.WriteString(normalStyle.Render(fmt.Sprintf("格式: %s\n", m.config.Output.Format)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("编码: %s\n", m.config.Output.Encoding)))
	if m.config.Output.FilePath != "" {
		content.WriteString(normalStyle.Render(fmt.Sprintf("输出文件: %s\n", m.config.Output.FilePath)))
	}
	
	return content.String()
}

func (m *ConfigEditorModel) renderFileProcessingConfig() string {
	var content strings.Builder
	
	content.WriteString(normalStyle.Render(fmt.Sprintf("包含隐藏文件: %v\n", m.config.FileProcessing.IncludeHidden)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("最大文件大小: %d\n", m.config.FileProcessing.MaxFileSize)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("最大深度: %d\n", m.config.FileProcessing.MaxDepth)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("包含内容: %v\n", m.config.FileProcessing.IncludeContent)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("包含哈希: %v\n", m.config.FileProcessing.IncludeHash)))
	
	return content.String()
}

func (m *ConfigEditorModel) renderUIConfig() string {
	var content strings.Builder
	
	content.WriteString(normalStyle.Render(fmt.Sprintf("主题: %s\n", m.config.UI.Theme)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("显示进度: %v\n", m.config.UI.ShowProgress)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("显示大小: %v\n", m.config.UI.ShowSize)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("显示日期: %v\n", m.config.UI.ShowDate)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("显示预览: %v\n", m.config.UI.ShowPreview)))
	
	return content.String()
}

func (m *ConfigEditorModel) renderPerformanceConfig() string {
	var content strings.Builder
	
	content.WriteString(normalStyle.Render(fmt.Sprintf("最大工作线程: %d\n", m.config.Performance.MaxWorkers)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("缓冲区大小: %d\n", m.config.Performance.BufferSize)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("缓存启用: %v\n", m.config.Performance.CacheEnabled)))
	content.WriteString(normalStyle.Render(fmt.Sprintf("缓存大小: %d\n", m.config.Performance.CacheSize)))
	
	return content.String()
}

func (m *ConfigEditorModel) saveConfig() tea.Cmd {
	return func() tea.Msg {
		// 这里应该实现保存配置逻辑
		return ConfigUpdateMsg{Config: m.config}
	}
}

// FileListMsg 文件列表消息
type FileListMsg struct {
	Items []selector.FileItem
}