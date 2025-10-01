# TUI 主界面工作流

## 概述
TUI主界面是整个应用程序的入口点，负责管理不同视图状态之间的切换和协调各个子模型的交互。

## 状态管理

### 应用状态 (AppState)
- `StateInit`: 初始化状态
- `StateInput`: 输入状态（主界面）
- `StateSelect`: 文件选择状态
- `StateProcessing`: 处理中状态
- `StateResult`: 结果显示状态
- `StateConfig`: 配置编辑状态
- `StateError`: 错误状态

### 视图类型 (ViewType)
- `ViewMain`: 主视图
- `ViewSelect`: 文件选择视图
- `ViewProgress`: 进度条视图
- `ViewResult`: 结果视图
- `ViewConfig`: 配置视图

## 主工作流

### 1. 初始化流程
```
main() -> initialModel() -> tea.NewProgram() -> p.Run()
```

### 2. 状态转换流程
```
StateInput -> StateSelect -> StateProcessing -> StateResult
     ↓           ↓              ↓              ↓
StateConfig   返回           错误处理       返回主界面
```

### 3. 消息处理机制

#### 键盘消息处理
- `ctrl+c`, `q`: 退出程序
- `esc`: 返回上一状态
- 其他按键根据当前视图状态分发到对应子模型

#### 窗口大小消息
- 更新主模型和子模型的宽高信息
- 触发视图重新渲染

#### 自定义消息处理
- `ProgressMsg`: 更新进度条状态
- `ResultMsg`: 接收处理结果，切换到结果视图
- `ErrorMsg`: 处理错误，切换到错误状态
- `FileSelectionMsg`: 接收文件选择结果
- `ConfigUpdateMsg`: 接收配置更新

## 子模型协调

### 文件选择器模型 (FileSelectorModel)
- 管理文件列表的加载和显示
- 处理文件多选逻辑
- 返回选中的文件路径列表

### 进度条模型 (ProgressModel)
- 显示处理进度
- 显示当前状态信息
- 处理取消操作

### 结果查看器模型 (ResultViewerModel)
- 显示扫描结果概览
- 提供标签页切换功能
- 支持结果保存

### 配置编辑器模型 (ConfigEditorModel)
- 显示和编辑配置项
- 提供分类标签页
- 支持配置保存

## 关键问题识别

### 已知问题
1. **Tab切换标签未成功**: 结果查看器和配置编辑器中的Tab键切换功能可能存在实现问题
2. **TUI界面加载文件列表失败**: 文件选择器可能存在文件加载或显示问题
3. **配置器管理界面按键无效**: 配置编辑器中除esc、ctrl+C外的其他按键可能未正确响应

## 测试要点

### 状态转换测试
- 验证各状态之间的正确转换
- 测试异常状态下的错误处理
- 验证消息传递的完整性

### 子模型集成测试
- 测试子模型与主模型的消息交互
- 验证子模型生命周期管理
- 测试窗口大小变化时的响应

### 用户交互测试
- 测试键盘快捷键的响应
- 验证视图切换的流畅性
- 测试边界条件下的稳定性