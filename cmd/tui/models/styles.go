package models

import "github.com/charmbracelet/lipgloss"

// UI样式定义
var (
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginBottom(1)

	SelectedStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	NormalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA"))

	HelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)

	WarningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Bold(true)
)