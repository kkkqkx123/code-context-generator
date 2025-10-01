# TUI界面自动补全集成方案

## 概述

本文档描述了如何将现有的自动补全模块集成到TUI（终端用户界面）中，以提升用户体验。当前自动补全模块仅在CLI中可用，TUI界面尚未集成此功能。

## 当前状态分析

### 自动补全模块现状
- ✅ 已实现完整的文件路径自动补全功能
- ✅ 支持多种补全类型：文件路径、目录、扩展名、模式匹配
- ✅ 具备缓存机制，性能良好
- ✅ 已拆分为模块化结构（interfaces.go, filepath_autocompleter.go, command_autocompleter.go, suggestion_provider.go, utils.go）
- ✅ CLI已集成，通过`code-context-generator autocomplete`命令可用

### TUI界面现状
- ✅ 基于Bubble Tea框架构建
- ✅ 包含多个输入字段：路径输入、输出格式选择、排除模式等
- ❌ **未集成自动补全功能**
- ❌ 用户需要手动输入完整路径

## 集成方案

### 1. 集成目标

在TUI界面的以下输入场景中提供自动补全：
- 路径输入（主输入字段）
- 文件选择器中的路径输入
- 配置编辑器中的路径相关字段
- 排除/包含模式输入

### 2. 技术实现方案

#### 2.1 创建TUI自动补全模型

创建新的模型文件：`cmd/tui/models/autocomplete.go`

```go
package models

import (
    "code-context-generator/internal/autocomplete"
    "code-context-generator/pkg/types"
    tea "github.com/charmbracelet/bubbletea"
)

// AutocompleteModel 自动补全模型
type AutocompleteModel struct {
    autocompleter autocomplete.Autocompleter
    suggestions   []string
    selectedIndex int
    input         string
    visible       bool
    completeType  types.CompleteType
    maxSuggestions int
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
    }
}

// Update 更新模型状态
func (a *AutocompleteModel) Update(msg tea.Msg) (*AutocompleteModel, tea.Cmd) {
    // 实现更新逻辑
    return a, nil
}

// View 渲染视图
func (a *AutocompleteModel) View() string {
    if !a.visible || len(a.suggestions) == 0 {
        return ""
    }
    // 实现渲染逻辑
    return ""
}

// UpdateSuggestions 更新建议列表
func (a *AutocompleteModel) UpdateSuggestions(input string) {
    context := &types.CompleteContext{
        Type: a.completeType,
        Data: make(map[string]interface{}),
    }
    
    suggestions, err := a.autocompleter.Complete(input, context)
    if err != nil {
        a.suggestions = []string{}
    } else {
        a.suggestions = suggestions
    }
    a.selectedIndex = 0
}
```

#### 2.2 集成到现有模型

更新主模型`MainModel`，集成自动补全功能：

```go
type MainModel struct {
    // ... 现有字段
    
    // 自动补全
    autocomplete    *AutocompleteModel
    showAutocomplete bool
}

// 在Update方法中处理自动补全
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // 处理Tab键触发自动补全
        if msg.String() == "tab" && m.showAutocomplete {
            // 应用当前选中的建议
            return m.applyAutocompleteSuggestion(), nil
        }
        
        // 处理其他键，更新输入时触发自动补全
        if m.state == StateInput {
            m.updateAutocomplete(msg.String())
        }
    }
    
    // ... 其他处理
}

// updateAutocomplete 更新自动补全
func (m *MainModel) updateAutocomplete(input string) {
    m.autocomplete.UpdateSuggestions(input)
    m.showAutocomplete = len(m.autocomplete.suggestions) > 0
}
```

#### 2.3 视图集成

在主视图中显示自动补全建议：

```go
func (m MainModel) View() string {
    var b strings.Builder
    
    // ... 现有视图代码
    
    // 显示自动补全建议
    if m.showAutocomplete {
        b.WriteString(m.autocomplete.View())
    }
    
    return b.String()
}
```

### 3. 用户交互设计

#### 3.1 快捷键
- `Tab`：触发/应用自动补全
- `↑/↓`：在建议列表中导航
- `Esc`：关闭自动补全
- `Enter`：应用选中建议并关闭

#### 3.2 视觉反馈
- 建议列表显示在输入字段下方
- 当前选中项高亮显示
- 文件和目录使用不同图标/颜色区分

#### 3.3 智能触发
- 输入达到一定字符数后自动触发（可配置）
- 根据输入内容智能判断补全类型
- 支持手动触发（Tab键）

### 4. 性能优化

#### 4.1 缓存策略
- 利用现有`FilePathAutocompleter`的缓存机制
- 预加载常用目录
- 异步更新缓存，避免阻塞UI

#### 4.2 异步处理
```go
// 异步更新建议
type updateSuggestionsMsg struct {
    suggestions []string
    err         error
}

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
```

#### 4.3 限制建议数量
- 默认最多显示10个建议（可配置）
- 优先显示最匹配的项目
- 支持滚动查看更多建议

### 5. 配置集成

#### 5.1 自动补全配置
扩展现有配置结构：

```go
type TUIConfig struct {
    // ... 现有配置
    Autocomplete AutocompleteUIConfig `yaml:"autocomplete" json:"autocomplete"`
}

type AutocompleteUIConfig struct {
    Enabled        bool `yaml:"enabled" json:"enabled"`
    MinChars       int  `yaml:"min_chars" json:"min_chars"`
    MaxSuggestions int  `yaml:"max_suggestions" json:"max_suggestions"`
    ShowIcons      bool `yaml:"show_icons" json:"show_icons"`
    TriggerDelay   int  `yaml:"trigger_delay_ms" json:"trigger_delay_ms"`
}
```

#### 5.2 环境变量支持
```bash
CODE_CONTEXT_TUI_AUTOCOMPLETE_ENABLED=true
CODE_CONTEXT_TUI_AUTOCOMPLETE_MIN_CHARS=2
CODE_CONTEXT_TUI_AUTOCOMPLETE_MAX_SUGGESTIONS=15
```

### 6. 错误处理

#### 6.1 容错机制
- 自动补全失败时不影响主功能
- 显示友好的错误提示
- 记录错误日志供调试

#### 6.2 边界情况处理
- 空输入处理
- 无效路径处理
- 权限不足处理
- 网络文件系统超时处理

### 7. 测试策略

#### 7.1 单元测试
- 测试自动补全模型的各种状态
- 测试用户交互逻辑
- 测试配置加载和应用

#### 7.2 集成测试
- 测试与主模型的集成
- 测试视图渲染
- 测试性能表现

#### 7.3 用户测试
- 收集用户反馈
- 优化交互体验
- 调整默认配置

## 实施计划

### 第一阶段：基础集成（1-2周）
1. 创建AutocompleteModel
2. 集成到MainModel
3. 实现基本视图渲染
4. 添加Tab键触发功能

### 第二阶段：增强功能（1-2周）
1. 添加上下键导航
2. 实现智能类型判断
3. 添加视觉反馈
4. 集成配置系统

### 第三阶段：性能优化（1周）
1. 实现异步建议更新
2. 添加缓存机制
3. 优化大目录处理
4. 添加性能监控

### 第四阶段：测试和完善（1周）
1. 编写单元测试
2. 用户测试和反馈
3. 修复bug和优化体验
4. 更新文档

## 预期效果

### 用户体验提升
- 减少输入错误
- 提高输入效率
- 降低学习成本
- 增强专业感

### 功能完整性
- 与CLI功能保持一致
- 支持所有补全类型
- 可配置性强
- 跨平台兼容

### 性能表现
- 响应时间 < 100ms
- 内存占用合理
- 不影响主功能性能
- 支持大目录处理

## 风险评估

### 技术风险
- **风险**：Bubble Tea框架限制
- **缓解**：提前验证技术可行性，准备备选方案

### 性能风险
- **风险**：大目录导致卡顿
- **缓解**：异步处理、限制建议数量、添加超时机制

### 用户体验风险
- **风险**：自动补全干扰正常输入
- **缓解**：提供关闭选项、可配置触发方式、渐进式推出

## 后续优化

1. **历史记录**：记住用户常用路径
2. **模糊匹配**：支持模糊搜索
3. **多语言**：支持多语言路径
4. **自定义图标**：允许用户自定义文件类型图标
5. **插件系统**：支持第三方补全源