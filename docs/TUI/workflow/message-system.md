# 消息系统工作流

## 概述
TUI应用使用Bubble Tea框架的消息系统实现组件间的通信，通过定义自定义消息类型来传递状态和事件信息。

## 消息类型定义

### ProgressMsg - 进度消息
```go
type ProgressMsg struct {
    Progress float64  // 进度值 (0.0-1.0)
    Status   string   // 状态描述文本
}
```
**用途**: 在文件处理过程中传递进度更新信息
**发送方**: 文件处理逻辑
**接收方**: 进度条模型

### ResultMsg - 结果消息
```go
type ResultMsg struct {
    Result *types.WalkResult  // 扫描结果数据
}
```
**用途**: 传递文件扫描的最终结果
**发送方**: 文件扫描逻辑
**接收方**: 结果查看器模型、主模型

### ErrorMsg - 错误消息
```go
type ErrorMsg struct {
    Err error  // 错误信息
}
```
**用途**: 传递处理过程中发生的错误
**发送方**: 任何可能出错的组件
**接收方**: 主模型、错误处理逻辑

### FileSelectionMsg - 文件选择消息
```go
type FileSelectionMsg struct {
    Selected []string  // 选中的文件路径列表
}
```
**用途**: 传递用户在文件选择器中选中的文件
**发送方**: 文件选择器模型
**接收方**: 主模型

### ConfigUpdateMsg - 配置更新消息
```go
type ConfigUpdateMsg struct {
    Config *types.Config  // 更新后的配置对象
}
```
**用途**: 传递配置更新信息
**发送方**: 配置编辑器模型
**接收方**: 主模型

### FileListMsg - 文件列表消息
```go
type FileListMsg struct {
    Items []selector.FileItem  // 文件项列表
}
```
**用途**: 传递目录扫描得到的文件列表
**发送方**: 文件加载逻辑
**接收方**: 文件选择器模型

## 消息流工作流程

### 1. 文件扫描流程
```
文件扫描开始 -> ProgressMsg(进度更新) -> ... -> ResultMsg(最终结果) -> 主模型
                                      ↓
                               ErrorMsg(如果出错)
```

### 2. 文件选择流程
```
用户选择文件 -> FileSelectionMsg -> 主模型 -> 更新配置 -> 返回主界面
```

### 3. 配置更新流程
```
配置编辑完成 -> ConfigUpdateMsg -> 主模型 -> 更新状态 -> 返回主界面
```

### 4. 文件加载流程
```
目录扫描 -> FileListMsg -> 文件选择器 -> 更新文件列表 -> 重新渲染
```

## 消息处理机制

### 主模型消息分发
主模型作为中央协调器，负责将消息分发到合适的子模型：

```go
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyMsg(msg)
    case *models.ProgressMsg:
        // 更新进度条
        if m.progressBar != nil {
            m.progressBar.SetProgress(msg.Progress)
            m.progressBar.SetStatus(msg.Status)
        }
    case *models.ResultMsg:
        // 处理结果，切换状态
        m.result = msg.Result
        m.state = StateResult
        m.currentView = ViewResult
        if m.resultViewer != nil {
            m.resultViewer.SetResult(m.result)
        }
    // ... 其他消息处理
    }
}
```

### 子模型消息处理
各子模型独立处理自己的消息：

```go
// 文件选择器处理文件列表更新
func (m *FileSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case FileListMsg:
        m.items = msg.Items
        m.updateViewport()
    // ... 其他处理
    }
}
```

## 消息传递模式

### 命令模式 (tea.Cmd)
使用命令模式异步执行操作并返回消息：

```go
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
```

### 直接消息传递
在主模型中直接处理和转发消息：

```go
case *models.ConfigUpdateMsg:
    cfg = msg.Config
    m.state = StateInput
    m.currentView = ViewMain
    return m, nil
```

## 错误处理

### 错误消息传播
```go
case *models.ErrorMsg:
    m.err = msg.Err
    m.state = StateError
    return m, nil
```

### 错误恢复机制
- 记录错误信息
- 切换到错误状态
- 提供用户友好的错误提示
- 支持从错误状态恢复

## 性能考虑

### 消息频率控制
- 避免过于频繁的进度更新消息
- 使用节流机制控制消息发送频率
- 合并相关的状态更新消息

### 内存管理
- 及时清理不需要的消息数据
- 避免在消息中传递大对象
- 使用指针传递复杂数据结构

## 测试策略

### 消息流测试
- 验证消息的正确传递路径
- 测试消息处理的时序性
- 验证状态更新的正确性

### 边界条件测试
- 测试空消息处理
- 验证nil指针消息
- 测试并发消息处理

### 错误场景测试
- 测试错误消息的传递
- 验证错误恢复机制
- 测试异常情况下的消息处理

## 调试和监控

### 消息日志
- 记录重要消息的收发情况
- 跟踪消息处理时间
- 监控消息队列长度

### 状态跟踪
- 记录关键状态变化
- 跟踪消息触发的状态转换
- 监控内存使用情况