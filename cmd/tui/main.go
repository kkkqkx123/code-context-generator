// Package main TUI应用程序主入口
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"code-context-generator/cmd/tui/models"
	"code-context-generator/internal/filesystem"
	"code-context-generator/pkg/types"

	tea "github.com/charmbracelet/bubbletea"
)

// 全局变量用于在回调中访问当前程序
var currentProgram *tea.Program

// 全局变量用于在回调中访问当前模型
var currentModel tea.Model

var (
	cfg        *types.Config
	configPath string
	version    = "1.0.0"
)

// AppState 应用程序状态
type AppState int

const (
	StateInit AppState = iota
	StateInput
	StateSelect
	StateProcessing
	StateResult
	StateConfig
	StateError
)

// ViewType 视图类型
type ViewType int

const (
	ViewMain ViewType = iota
	ViewSelect
	ViewProgress
	ViewResult
	ViewConfig
)

// MainModel 主模型
type MainModel struct {
	state           AppState
	currentView     ViewType
	pathInput       string
	outputFormat    string
	outputPath      string
	excludePatterns []string
	includePatterns []string
	options         types.WalkOptions
	result          *types.WalkResult
	err             error
	width           int
	height          int
	// 子模型
	fileSelector *models.FileSelectorModel
	progressBar  *models.ProgressModel
	resultViewer *models.ResultViewerModel
	configEditor *models.ConfigEditorModel
	// 自动补全
	autocomplete    *models.AutocompleteModel
	showAutocomplete bool
}

// 初始化函数
func init() {
	// 初始化模型
	models.SetConfig(nil)
}

// main 主函数
func main() {
	// 创建模型
	model := initialModel()

	// 创建程序
	p := tea.NewProgram(model, tea.WithAltScreen())
	
	// 设置全局程序引用
	currentProgram = p

	// 运行程序
	if _, err := p.Run(); err != nil {
		fmt.Printf("运行程序失败: %v\n", err)
		os.Exit(1)
	}
}

// initialModel 创建初始模型
func initialModel() MainModel {
	return MainModel{
		state:           StateInit,
		currentView:     ViewMain,
		pathInput:       ".",
		outputFormat:    "json",
		outputPath:      "",
		excludePatterns: []string{},
		includePatterns: []string{},
		options: types.WalkOptions{
			MaxDepth:        1,
			MaxFileSize:     10 * 1024 * 1024,
			ExcludePatterns: []string{},
			IncludePatterns: []string{},
			FollowSymlinks:  false,
			ShowHidden:      false,
		},
		// 创建初始模型
		fileSelector:     models.NewFileSelectorModel("."),
		progressBar:      models.NewProgressModel(),
		resultViewer:     models.NewResultViewerModel(),
		configEditor:     models.NewConfigEditorModel(nil),
		autocomplete:     models.NewAutocompleteModel(),
		showAutocomplete: false,
	}
}

// Init 初始化
func (m MainModel) Init() tea.Cmd {
	// 初始化时不需要做任何特殊操作
	return nil
}

// Update 更新模型
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == StateProcessing {
			// 处理中状态只响应 Ctrl+C
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			return m, nil
		}
		return m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.autocomplete != nil {
			m.autocomplete.SetSize(msg.Width, msg.Height)
		}
		return m, nil
	case *models.ProgressMsg:
		if m.progressBar != nil {
			m.progressBar.SetProgress(msg.Progress)
			m.progressBar.SetStatus(msg.Status)
		}
		return m, nil
	case *models.ResultMsg:
		m.result = msg.Result
		m.state = StateResult
		m.currentView = ViewResult
		if m.resultViewer != nil {
			m.resultViewer.SetResult(m.result)
		}
		return m, nil
	case *models.ErrorMsg:
		m.err = msg.Err
		m.state = StateError
		return m, nil
	case *models.FileSelectionMsg:
		m.options.SelectedFiles = msg.Selected
		m.state = StateInput
		m.currentView = ViewMain
		return m, nil
	case *models.ConfigUpdateMsg:
		cfg = msg.Config
		m.state = StateInput
		m.currentView = ViewMain
		return m, nil
	case *models.UpdateSuggestionsMsg:
		if m.autocomplete != nil {
			m.autocomplete.SetSuggestions(msg.Suggestions)
			m.showAutocomplete = len(msg.Suggestions) > 0
		}
		return m, nil
	case *models.ApplySuggestionMsg:
		if msg.Suggestion != "" {
			m.pathInput = msg.Suggestion
			m.showAutocomplete = false
			if m.autocomplete != nil {
				m.autocomplete.Hide()
			}
		}
		return m, nil
	default:
		// 更新子模型
		switch m.currentView {
		case ViewSelect:
			if m.fileSelector != nil {
				newModel, cmd := m.fileSelector.Update(msg)
				m.fileSelector = newModel.(*models.FileSelectorModel)
				return m, cmd
			}
		case ViewProgress:
			if m.progressBar != nil {
				newModel, cmd := m.progressBar.Update(msg)
				m.progressBar = newModel.(*models.ProgressModel)
				return m, cmd
			}
		case ViewResult:
			if m.resultViewer != nil {
				newModel, cmd := m.resultViewer.Update(msg)
				m.resultViewer = newModel.(*models.ResultViewerModel)
				return m, cmd
			}
		case ViewConfig:
			if m.configEditor != nil {
				newModel, cmd := m.configEditor.Update(msg)
				m.configEditor = newModel.(*models.ConfigEditorModel)
				return m, cmd
			}
		}
	}

	return m, nil
}

// View 渲染视图
func (m MainModel) View() string {
	if m.err != nil {
		return m.renderError()
	}

	switch m.currentView {
	case ViewMain:
		return m.renderMainView()
	case ViewSelect:
		if m.fileSelector != nil {
			return m.fileSelector.View()
		}
	case ViewProgress:
		if m.progressBar != nil {
			return m.progressBar.View()
		}
	case ViewResult:
		if m.resultViewer != nil {
			return m.resultViewer.View()
		}
	case ViewConfig:
		if m.configEditor != nil {
			return m.configEditor.View()
		}
	}

	return "未知视图"
}

// handleKeyMsg 处理键盘消息
func (m MainModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// 全局退出快捷键
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	// ESC键返回上一级
	if msg.String() == "esc" {
		return m.handleEscKey()
	}

	switch m.state {
	case StateInit, StateInput:
		return m.handleMainKeys(msg)
	case StateError:
		return m.handleErrorKeys(msg)
	case StateSelect:
		return m.handleSelectKeys(msg)
	case StateProcessing:
		return m.handleProcessingKeys(msg)
	case StateResult:
		return m.handleResultKeys(msg)
	case StateConfig:
		return m.handleConfigKeys(msg)
	default:
		return m, nil
	}
}

// handleMainKeys 处理主界面按键
func (m MainModel) handleMainKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q": // 仅在主界面按q退出
		return m, tea.Quit
	case "enter":
		return m.startProcessing()
	case "s":
		m.state = StateSelect
		m.currentView = ViewSelect
		// 重新初始化文件选择器以确保加载文件列表
		if m.fileSelector != nil {
			return m, m.fileSelector.Init()
		}
		return m, nil
	case "c":
		m.state = StateConfig
		m.currentView = ViewConfig
		return m, nil
	case "r":
		if m.options.MaxDepth == 0 {
			m.options.MaxDepth = 1
		} else {
			m.options.MaxDepth = 0
		}
		return m, nil
	case "h":
		m.options.ShowHidden = !m.options.ShowHidden
		return m, nil
	case "tab":
		// Tab键触发自动补全
		if m.showAutocomplete && m.autocomplete != nil {
			// 应用当前选中的建议
			suggestion := m.autocomplete.GetSelectedSuggestion()
			if suggestion != "" {
				m.pathInput = suggestion
				m.showAutocomplete = false
				m.autocomplete.Hide()
			}
		}
		return m, nil
	case "up", "down":
		// 如果自动补全可见，则导航建议列表
		if m.showAutocomplete && m.autocomplete != nil {
			newModel, cmd := m.autocomplete.Update(msg)
			m.autocomplete = newModel
			return m, cmd
		}
		return m, nil
	default:
		// 处理输入
		if m.state == StateInput {
			return m.handleInput(msg)
		}
	}

	return m, nil
}

// handleErrorKeys 处理错误界面按键
func (m MainModel) handleErrorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter":
		m.state = StateInput
		m.err = nil
		return m, nil
	}
	return m, nil
}

// handleEscKey 处理ESC键返回上一级
func (m MainModel) handleEscKey() (tea.Model, tea.Cmd) {
	// 如果自动补全可见，先隐藏自动补全
	if m.showAutocomplete && m.autocomplete != nil {
		m.showAutocomplete = false
		m.autocomplete.Hide()
		return m, nil
	}
	
	switch m.state {
	case StateSelect:
		// 从文件选择器返回主界面
		m.state = StateInput
		m.currentView = ViewMain
		return m, nil
	case StateConfig:
		// 从配置编辑器返回主界面
		m.state = StateInput
		m.currentView = ViewMain
		return m, nil
	case StateResult:
		// 从结果查看器返回主界面
		m.state = StateInput
		m.currentView = ViewMain
		return m, nil
	case StateProcessing:
		// 处理中不允许返回，可以取消处理
		return m, nil
	case StateError:
		// 错误状态已经在handleErrorKeys中处理
		return m, nil
	default:
		// 主界面按ESC也退出
		return m, tea.Quit
	}
}

// handleSelectKeys 处理文件选择器按键
func (m MainModel) handleSelectKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// 文件选择器的按键处理在FileSelectorModel中
	// 将按键传递给文件选择器处理
	if m.fileSelector != nil {
		newModel, cmd := m.fileSelector.Update(msg)
		m.fileSelector = newModel.(*models.FileSelectorModel)
		return m, cmd
	}
	return m, nil
}

// handleProcessingKeys 处理处理中按键
func (m MainModel) handleProcessingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "c": // 取消处理
		// 这里应该实现取消逻辑
		m.state = StateInput
		m.currentView = ViewMain
		return m, nil
	}
	return m, nil
}

// handleResultKeys 处理结果查看器按键
func (m MainModel) handleResultKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "b": // 返回主界面
		m.state = StateInput
		m.currentView = ViewMain
		return m, nil
	case "s": // 保存结果
		// 这里应该实现保存逻辑
		return m, nil
	case "q": // 退出程序
		return m, tea.Quit
	default:
		// 将所有其他按键传递给结果查看器处理（包括tab、up、down等）
		if m.resultViewer != nil {
			newModel, cmd := m.resultViewer.Update(msg)
			m.resultViewer = newModel.(*models.ResultViewerModel)
			return m, cmd
		}
	}
	return m, nil
}

// handleConfigKeys 处理配置编辑器按键
func (m MainModel) handleConfigKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// 配置编辑器的按键处理在ConfigEditorModel中
	// 将除了全局快捷键外的所有按键传递给配置编辑器处理
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		// ESC键已经在handleEscKey中处理，这里不再处理
		return m, nil
	default:
		// 将所有其他按键传递给配置编辑器处理
		if m.configEditor != nil {
			newModel, cmd := m.configEditor.Update(msg)
			m.configEditor = newModel.(*models.ConfigEditorModel)
			return m, cmd
		}
	}
	return m, nil
}

// handleInput 处理输入
func (m MainModel) handleInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "backspace":
		if len(m.pathInput) > 0 {
			m.pathInput = m.pathInput[:len(m.pathInput)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.pathInput += msg.String()
		}
	}
	
	// 更新自动补全建议
	if m.autocomplete != nil && len(m.pathInput) > 0 {
		// 异步更新建议
		return m, m.autocomplete.UpdateSuggestionsAsync(m.pathInput)
	}
	
	m.options.MaxDepth = 0
	if m.pathInput != "." {
		m.options.MaxDepth = 1
	}
	return m, nil
}

// startProcessing 开始处理
func (m MainModel) startProcessing() (tea.Model, tea.Cmd) {
	m.state = StateProcessing
	m.currentView = ViewProgress
	currentModel = m // 设置全局模型引用

	return m, tea.Batch(
		tea.Tick(0, func(time.Time) tea.Msg {
			return models.ProgressMsg{Progress: 0, Status: "开始扫描..."}
		}),
		m.processFiles(),
	)
}

// processFiles 处理文件
func (m MainModel) processFiles() tea.Cmd {
	return func() tea.Msg {
		walker := filesystem.NewWalker()
		options := &types.WalkOptions{
			MaxDepth:        m.options.MaxDepth,
			MaxFileSize:     m.options.MaxFileSize,
			ExcludePatterns: m.options.ExcludePatterns,
			IncludePatterns: m.options.IncludePatterns,
			FollowSymlinks:  m.options.FollowSymlinks,
			ShowHidden:      m.options.ShowHidden,
			ExcludeBinary:   true,
			SelectedFiles:   m.options.SelectedFiles,
		}

		// 创建上下文用于控制进度更新
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// 启动进度更新goroutine
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()
			
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					// 进度将在回调中更新
				}
			}
		}()

		// 使用进度回调进行文件处理
		var result *types.ContextData
		var err error
		
		result, err = walker.WalkWithProgress(m.pathInput, options, func(p, t int, file string) {
			// 更新进度信息
			if currentProgram != nil && t > 0 {
				progress := float64(p) / float64(t)
				currentProgram.Send(models.ProgressMsg{
					Progress: progress,
					Status:   fmt.Sprintf("扫描文件中... %d/%d", p, t),
				})
			}
		})
		
		// 取消进度更新goroutine
		cancel()
		
		if err != nil {
			return models.ErrorMsg{Err: err}
		}

		// 转换结果类型
		walkResult := &types.WalkResult{
			Files:       result.Files,
			Folders:     result.Folders,
			FileCount:   result.FileCount,
			FolderCount: result.FolderCount,
			TotalSize:   result.TotalSize,
			RootPath:    m.pathInput,
		}
		
		return models.ResultMsg{Result: walkResult}
	}
}

// renderMainView 渲染主视图
func (m MainModel) renderMainView() string {
	var content strings.Builder

	// 标题
	content.WriteString(models.TitleStyle.Render("代码上下文生成器"))
	content.WriteString("\n\n")

	// 路径输入
	content.WriteString(models.NormalStyle.Render("扫描路径: "))
	content.WriteString(m.pathInput)
	content.WriteString("\n")
	
	// 显示自动补全建议
	if m.showAutocomplete && m.autocomplete != nil {
		suggestionsView := m.autocomplete.View()
		if suggestionsView != "" {
			content.WriteString(suggestionsView)
			content.WriteString("\n")
		}
	}
	
	content.WriteString("\n")

	// 选项
	content.WriteString(models.NormalStyle.Render("选项:\n"))
	recursive := "否"
	if m.options.MaxDepth != 0 {
		recursive = "是"
	}
	content.WriteString(fmt.Sprintf("\n  递归扫描: %s (按 r 切换)\n", recursive))

	hidden := "否"
	if m.options.ShowHidden {
		hidden = "是"
	}
	content.WriteString(fmt.Sprintf("  包含隐藏文件: %s (按 h 切换)\n", hidden))

	content.WriteString(fmt.Sprintf("  输出格式: %s\n", m.outputFormat))

	if m.outputPath != "" {
		content.WriteString(fmt.Sprintf("  输出文件: %s\n", m.outputPath))
	}

	content.WriteString("\n")

	// 操作提示
	content.WriteString(models.HelpStyle.Render("操作:\n"))
	content.WriteString("\n  Enter - 开始扫描\n")
	content.WriteString("  s - 选择文件\n")
	content.WriteString("  c - 配置设置\n")
	content.WriteString("  Tab - 应用自动补全建议\n")
	content.WriteString("  ↑↓ - 选择自动补全建议\n")
	content.WriteString("  ESC - 退出程序\n")
	content.WriteString("  Ctrl+C - 强制退出\n")

	return content.String()
}

// renderError 渲染错误视图
func (m MainModel) renderError() string {
	var content strings.Builder

	content.WriteString(models.ErrorStyle.Render("错误:\n"))
	content.WriteString(m.err.Error())
	content.WriteString("\n\n")
	content.WriteString(models.HelpStyle.Render("按 Esc 或 Enter 返回"))

	return content.String()
}

// getDefaultConfig 获取默认配置
func getDefaultConfig() *types.Config {
	return &types.Config{
		Output: types.OutputConfig{
			Format:   "json",
			Encoding: "utf-8",
		},
		FileProcessing: types.FileProcessingConfig{
			IncludeHidden: false,
			MaxFileSize:   10 * 1024 * 1024,
			MaxDepth:      0,
			ExcludePatterns: []string{
				"*.exe", "*.dll", "*.so", "*.dylib",
				"*.pyc", "*.pyo", "*.pyd",
				"node_modules", ".git", ".svn",
			},
			IncludePatterns: []string{},
			IncludeContent:  false,
			IncludeHash:     false,
		},
		UI: types.UIConfig{
			Theme:        "default",
			ShowProgress: true,
			ShowSize:     true,
			ShowDate:     true,
			ShowPreview:  true,
		},
		Performance: types.PerformanceConfig{
			MaxWorkers:   4,
			BufferSize:   1024,
			CacheEnabled: true,
			CacheSize:    100,
		},
		Logging: types.LoggingConfig{
			Level:      "info",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
		},
	}
}
