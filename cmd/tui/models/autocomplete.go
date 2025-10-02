// Package models TUI模型定义
package models

import (
	"strings"
	
	"code-context-generator/internal/autocomplete"
	"code-context-generator/pkg/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AutocompleteModel 自动补全模型
type AutocompleteModel struct {
	autocompleter   autocomplete.Autocompleter
	suggestions     []string
	selectedIndex   int
	input          string
	visible        bool
	completeType   types.CompleteType
	maxSuggestions int
	width          int
	height         int
}

// NewAutocompleteModel 创建自动补全模型
func NewAutocompleteModel() *AutocompleteModel {
	return &AutocompleteModel{
		autocompleter:  autocomplete.NewAutocompleter(nil),
		suggestions:    []string{},
		selectedIndex:  0,
		visible:        false,
		completeType:   types.CompleteFilePath,
		maxSuggestions: 10,
		width:          40,
		height:         10,
	}
}

// Init 初始化
func (a *AutocompleteModel) Init() tea.Cmd {
	return nil
}

// Update 更新模型状态
func (a *AutocompleteModel) Update(msg tea.Msg) (*AutocompleteModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !a.visible {
			return a, nil
		}

		switch msg.String() {
		case "up", "k":
			if a.selectedIndex > 0 {
				a.selectedIndex--
			}
		case "down", "j":
			if a.selectedIndex < len(a.suggestions)-1 {
				a.selectedIndex++
			}
		case "tab", "enter":
			// 应用选中建议
			if a.selectedIndex >= 0 && a.selectedIndex < len(a.suggestions) {
				return a, a.applySuggestion()
			}
		case "esc":
			a.visible = false
			a.selectedIndex = 0
		}
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
	}

	return a, nil
}

// View 渲染视图
func (a *AutocompleteModel) View() string {
	if !a.visible || len(a.suggestions) == 0 {
		return ""
	}

	var b strings.Builder
	
	// 样式定义
	suggestionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Background(lipgloss.Color("235"))
	
	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("236")).
		Bold(true)

	// 渲染建议列表
	maxDisplay := min(a.maxSuggestions, len(a.suggestions))
	for i := 0; i < maxDisplay; i++ {
		style := suggestionStyle
		if i == a.selectedIndex {
			style = selectedStyle
		}
		
		suggestion := a.suggestions[i]
		if len(suggestion) > a.width-2 {
			suggestion = suggestion[:a.width-5] + "..."
		}
		
		b.WriteString(style.Render("  " + suggestion))
		if i < maxDisplay-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// UpdateSuggestions 更新建议列表
func (a *AutocompleteModel) UpdateSuggestions(input string) {
	a.input = input
	if input == "" {
		a.suggestions = []string{}
		a.visible = false
		return
	}

	context := &types.CompleteContext{
		Type: a.completeType,
		Data: make(map[string]interface{}),
	}

	suggestions, err := a.autocompleter.Complete(input, context)
	if err != nil {
		a.suggestions = []string{}
		a.visible = false
		return
	}

	a.suggestions = suggestions
	a.selectedIndex = 0
	a.visible = len(suggestions) > 0
}

// GetSuggestions 获取建议列表
func (a *AutocompleteModel) GetSuggestions() []string {
	return a.suggestions
}

// SetSuggestions 设置建议列表
func (a *AutocompleteModel) SetSuggestions(suggestions []string) {
	a.suggestions = suggestions
	a.visible = len(suggestions) > 0
	a.selectedIndex = 0
}

// GetSelectedIndex 获取选中索引
func (a *AutocompleteModel) GetSelectedIndex() int {
	return a.selectedIndex
}

// SetSelectedIndex 设置选中索引
func (a *AutocompleteModel) SetSelectedIndex(index int) {
	a.selectedIndex = index
}

// UpdateSuggestionsAsync 异步更新建议
func (a *AutocompleteModel) UpdateSuggestionsAsync(input string) tea.Cmd {
	return func() tea.Msg {
		context := &types.CompleteContext{
			Type: a.completeType,
			Data: make(map[string]interface{}),
		}

		suggestions, err := a.autocompleter.Complete(input, context)
		return updateSuggestionsMsg{
			suggestions: suggestions,
			err:         err,
		}
	}
}

// GetSelectedSuggestion 获取当前选中的建议
func (a *AutocompleteModel) GetSelectedSuggestion() string {
	if a.selectedIndex >= 0 && a.selectedIndex < len(a.suggestions) {
		return a.suggestions[a.selectedIndex]
	}
	return ""
}

// IsVisible 检查是否可见
func (a *AutocompleteModel) IsVisible() bool {
	return a.visible
}

// Hide 隐藏自动补全
func (a *AutocompleteModel) Hide() {
	a.visible = false
	a.selectedIndex = 0
}

// Show 显示自动补全
func (a *AutocompleteModel) Show() {
	a.visible = true
}

// SetCompleteType 设置补全类型
func (a *AutocompleteModel) SetCompleteType(completeType types.CompleteType) {
	a.completeType = completeType
}

// SetMaxSuggestions 设置最大建议数量
func (a *AutocompleteModel) SetMaxSuggestions(max int) {
	a.maxSuggestions = max
}

// SetSize 设置窗口大小
func (a *AutocompleteModel) SetSize(width, height int) {
	a.width = width
	a.height = height
}

// applySuggestion 应用建议的消息
func (a *AutocompleteModel) applySuggestion() tea.Cmd {
	return func() tea.Msg {
		return applySuggestionMsg{
			suggestion: a.GetSelectedSuggestion(),
		}
	}
}

// updateSuggestionsMsg 更新建议消息
type updateSuggestionsMsg struct {
	suggestions []string
	err         error
}

// applySuggestionMsg 应用建议消息
type applySuggestionMsg struct {
	suggestion string
}

// min 返回最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}