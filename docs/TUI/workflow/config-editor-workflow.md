# 配置编辑器工作流

## 概述
配置编辑器提供可视化界面用于编辑应用程序的各种配置参数，支持分类标签页展示和实时编辑功能。

## 核心功能

### 多标签页配置分类
- **输出标签页**: 管理输出相关的配置项
- **文件处理标签页**: 管理文件处理相关的配置项
- **UI标签页**: 管理界面相关的配置项
- **性能标签页**: 管理性能相关的配置项

### 配置项展示
- 以键值对形式显示配置项
- 支持当前值格式化显示
- 提供配置项分类组织
- 显示配置项说明

### 交互功能
- Tab键切换配置分类标签页
- 上下键选择配置项
- Enter键编辑当前配置项
- s键保存配置
- Esc键返回主界面

## 工作流程

### 1. 初始化阶段
```
NewConfigEditorModel(config) -> 初始化currentTab为0 -> 设置focus为0 -> 等待用户交互
```

### 2. 标签页切换流程
```
Tab键按下 -> currentTab递增 -> 模运算确保范围(4个标签页) -> 重新渲染对应内容
```

### 3. 配置选择流程
```
上下键按下 -> focus增减 -> 边界检查 -> 高亮显示当前选中项
```

### 4. 配置编辑流程
```
Enter键按下 -> 进入编辑模式 -> 修改配置值 -> 保存或取消
```

## 状态管理

### 内部状态
- `config`: 配置对象引用
- `currentTab`: 当前标签页索引
- `width/height`: 视图尺寸
- `focus`: 当前焦点配置项索引

### 标签页定义
```go
tabs := []string{"输出", "文件处理", "UI", "性能"}
```

## 配置分类详情

### 输出配置 (renderOutputConfig)
- 默认格式: `config.Output.DefaultFormat`
- 输出目录: `config.Output.OutputDir`
- 文件名模板: `config.Output.FilenameTemplate`
- 时间戳格式: `config.Output.TimestampFormat`

### 文件处理配置 (renderFileProcessingConfig)
- 最大文件大小: `config.Filters.MaxFileSize`
- 最大深度: `config.Filters.MaxDepth`
- 跟随符号链接: `config.Filters.FollowSymlinks`
- 排除二进制文件: `config.Filters.ExcludeBinary`

### UI配置 (renderUIConfig)
- 主题: `config.UI.Theme`
- 显示进度: `config.UI.ShowProgress`
- 显示大小: `config.UI.ShowSize`
- 显示日期: `config.UI.ShowDate`
- 显示预览: `config.UI.ShowPreview`

### 性能配置 (renderPerformanceConfig)
- 最大工作线程: `config.Performance.MaxWorkers`
- 缓冲区大小: `config.Performance.BufferSize`
- 缓存启用: `config.Performance.CacheEnabled`
- 缓存大小: `config.Performance.CacheSize`

## 已知问题分析

### 配置器管理界面按键无效
**问题表现：**
- 除Esc和Ctrl+C外，其他按键无响应
- Tab键无法切换标签页
- 上下键无法选择配置项
- Enter键无法进入编辑模式

**可能原因：**
1. 键盘事件处理逻辑不完整
2. 按键映射配置错误
3. 焦点管理系统有问题
4. 事件分发机制异常

**代码分析：**
```go
case "tab":
    m.currentTab = (m.currentTab + 1) % 4 // 假设有4个标签页
case "up", "k":
    m.focus--
    if m.focus < 0 {
        m.focus = 0
    }
case "down", "j":
    m.focus++
case "enter":
    // 编辑当前项
```

**修复建议：**
1. 检查所有按键事件是否正确捕获
2. 验证按键字符串匹配逻辑
3. 确保焦点更新后触发视图重渲染
4. 添加按键响应的调试日志

### 配置编辑功能缺失
**当前问题：**
- Enter键处理为空实现
- 没有实际的配置编辑界面
- 保存功能未完整实现

## 键盘映射

| 按键 | 功能 |
|------|------|
| Tab | 切换配置分类标签页 |
| ↑/k | 向上选择配置项 |
| ↓/j | 向下选择配置项 |
| Enter | 编辑当前配置项 |
| s | 保存配置 |
| Esc | 返回主界面 |
| Ctrl+C/q | 退出程序 |

## 测试要点

### 标签页切换测试
- 测试Tab键正常切换
- 验证4个标签页的循环切换
- 测试内容正确更新
- 验证焦点保持逻辑

### 配置选择测试
- 测试上下键导航
- 验证边界条件处理
- 测试焦点高亮显示
- 验证选择状态保持

### 配置编辑测试
- 测试进入编辑模式
- 验证配置值修改
- 测试保存功能
- 验证配置持久化

### 交互响应测试
- 测试所有按键响应
- 验证界面刷新机制
- 测试配置分类切换
- 验证返回功能

### 配置完整性测试
- 验证所有配置项显示
- 测试配置值格式化
- 验证配置分类正确性
- 测试配置依赖关系