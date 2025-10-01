// Package autocomplete 提供自动补全功能
package autocomplete

import (
	"sort"
	"strings"
)

// CommandAutocompleter 命令自动补全器
type CommandAutocompleter struct {
	commands map[string]*CommandInfo
}

// CommandInfo 命令信息
type CommandInfo struct {
	Name        string
	Description string
	Aliases     []string
	Subcommands []string
	Options     []string
}

// NewCommandAutocompleter 创建命令自动补全器
func NewCommandAutocompleter() *CommandAutocompleter {
	return &CommandAutocompleter{
		commands: make(map[string]*CommandInfo),
	}
}

// RegisterCommand 注册命令
func (c *CommandAutocompleter) RegisterCommand(info *CommandInfo) {
	c.commands[info.Name] = info
}

// Complete 补全命令
func (c *CommandAutocompleter) Complete(input string) []string {
	var matches []string

	for name, info := range c.commands {
		if strings.HasPrefix(name, input) {
			matches = append(matches, name)
		}

		// 检查别名
		for _, alias := range info.Aliases {
			if strings.HasPrefix(alias, input) {
				matches = append(matches, alias)
			}
		}
	}

	sort.Strings(matches)
	return matches
}

// GetCommandInfo 获取命令信息
func (c *CommandAutocompleter) GetCommandInfo(command string) (*CommandInfo, bool) {
	info, exists := c.commands[command]
	return info, exists
}