# TUI 测试指南

## 概述
本文档提供TUI界面各个功能模块的测试策略、测试用例设计和测试执行指导，用于验证代码实现的完整性和正确性。

## 测试环境准备

### 依赖安装
```bash
go mod download
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install github.com/onsi/gomega/...@latest
```

### 测试目录结构
```
tests/
├── tui/
│   ├── main_test.go           # 主模型测试
│   ├── file_selector_test.go  # 文件选择器测试
│   ├── result_viewer_test.go  # 结果查看器测试
│   ├── config_editor_test.go  # 配置编辑器测试
│   ├── progress_test.go       # 进度条测试
│   └── integration_test.go    # 集成测试
```

## 单元测试策略

### 主模型测试 (MainModel)

#### 状态转换测试
```go
Describe("MainModel状态转换", func() {
    Context("从StateInput到StateSelect", func() {
        It("应该正确切换到文件选择状态", func() {
            model := initialModel()
            // 模拟切换到文件选择
            Expect(model.state).To(Equal(StateSelect))
            Expect(model.currentView).To(Equal(ViewSelect))
        })
    })
    
    Context("从StateProcessing到StateResult", func() {
        It("接收到ResultMsg后应该切换到结果状态", func() {
            // 模拟接收结果消息
            Expect(model.state).To(Equal(StateResult))
            Expect(model.currentView).To(Equal(ViewResult))
        })
    })
})
```

#### 消息处理测试
```go
Describe("消息处理", func() {
    Context("处理ProgressMsg", func() {
        It("应该更新进度条状态", func() {
            msg := &ProgressMsg{Progress: 0.5, Status: "处理中..."}
            // 验证进度条状态更新
            Expect(model.progressBar.GetProgress()).To(Equal(0.5))
            Expect(model.progressBar.GetStatus()).To(Equal("处理中..."))
        })
    })
})
```

### 文件选择器测试 (FileSelectorModel)

#### 文件加载测试
```go
Describe("文件加载功能", func() {
    Context("加载目录文件", func() {
        It("应该正确加载文件列表", func() {
            model := NewFileSelectorModel("./test_data")
            // 模拟文件加载
            Expect(len(model.items)).To(BeNumerically(">", 0))
            Expect(model.items[0].Name).ToNot(BeEmpty())
        })
    })
    
    Context("空目录处理", func() {
        It("应该正确处理空目录", func() {
            model := NewFileSelectorModel("./empty_dir")
            // 验证空目录处理
            Expect(len(model.items)).To(Equal(0))
        })
    })
})
```

#### 选择功能测试
```go
Describe("文件选择功能", func() {
    Context("单文件选择", func() {
        It("应该正确切换选择状态", func() {
            model.toggleSelection()
            Expect(model.selected[0]).To(BeTrue())
            model.toggleSelection()
            Expect(model.selected[0]).To(BeFalse())
        })
    })
    
    Context("全选功能", func() {
        It("应该选中所有文件", func() {
            model.selectAll()
            for i := range model.items {
                Expect(model.selected[i]).To(BeTrue())
            }
        })
    })
})
```

### 结果查看器测试 (ResultViewerModel)

#### Tab切换测试
```go
Describe("标签页切换功能", func() {
    Context("Tab键切换", func() {
        It("应该正确切换标签页", func() {
            model := NewResultViewerModel()
            initialTab := model.GetCurrentTab()
            
            // 模拟Tab键按下
            model.Update(tea.KeyMsg{Type: tea.KeyTab})
            
            Expect(model.GetCurrentTab()).To(Equal((initialTab + 1) % 3))
        })
    })
})
```

#### 内容渲染测试
```go
Describe("内容渲染", func() {
    Context("概览页面", func() {
        It("应该正确显示统计信息", func() {
            result := &types.WalkResult{
                RootPath: "./test",
                FileCount: 10,
                FolderCount: 2,
                TotalSize: 1024,
                ScanDuration: time.Second,
            }
            model.SetResult(result)
            
            content := model.View()
            Expect(content).To(ContainSubstring("根路径: ./test"))
            Expect(content).To(ContainSubstring("文件数量: 10"))
        })
    })
})
```

### 配置编辑器测试 (ConfigEditorModel)

#### 配置显示测试
```go
Describe("配置显示功能", func() {
    Context("输出配置标签页", func() {
        It("应该正确显示输出配置", func() {
            config := &types.Config{
                Output: types.OutputConfig{
                    DefaultFormat: "json",
                    OutputDir: "./output",
                },
            }
            model := NewConfigEditorModel(config)
            
            content := model.View()
            Expect(content).To(ContainSubstring("默认格式: json"))
            Expect(content).To(ContainSubstring("输出目录: ./output"))
        })
    })
})
```

#### 按键响应测试
```go
Describe("按键响应功能", func() {
    Context("Tab键切换标签页", func() {
        It("应该响应Tab键切换", func() {
            model := NewConfigEditorModel(&types.Config{})
            initialTab := model.GetCurrentTab()
            
            // 模拟Tab键
            model.Update(tea.KeyMsg{Type: tea.KeyTab})
            
            Expect(model.GetCurrentTab()).ToNot(Equal(initialTab))
        })
    })
})
```

## 集成测试策略

### 端到端流程测试
```go
Describe("端到端流程", func() {
    Context("完整扫描流程", func() {
        It("应该完成从输入到结果的完整流程", func() {
            // 1. 初始化主模型
            model := initialModel()
            
            // 2. 模拟文件选择
            model.currentView = ViewSelect
            // ... 模拟文件选择过程
            
            // 3. 模拟处理过程
            model.currentView = ViewProgress
            // ... 模拟进度更新
            
            // 4. 验证结果
            model.currentView = ViewResult
            Expect(model.result).ToNot(BeNil())
        })
    })
})
```

### 状态一致性测试
```go
Describe("状态一致性", func() {
    Context("模型状态同步", func() {
        It("主模型和子模型状态应该保持一致", func() {
            model := initialModel()
            
            // 验证初始状态一致性
            Expect(model.state).To(Equal(StateInput))
            Expect(model.currentView).To(Equal(ViewMain))
            
            // 测试状态转换后的同步性
            // ...
        })
    })
})
```

## 已知问题测试验证

### Tab切换问题测试
```go
Describe("Tab切换功能修复验证", func() {
    Context("结果查看器Tab切换", func() {
        It("Tab键应该正确切换标签页", func() {
            model := NewResultViewerModel()
            
            // 连续按Tab键
            for i := 0; i < 5; i++ {
                model.Update(tea.KeyMsg{Type: tea.KeyTab})
                Expect(model.GetCurrentTab()).To(Equal(i % 3))
            }
        })
    })
    
    Context("配置编辑器Tab切换", func() {
        It("Tab键应该正确切换配置标签页", func() {
            model := NewConfigEditorModel(&types.Config{})
            
            // 测试4个标签页的循环切换
            for i := 0; i < 8; i++ {
                model.Update(tea.KeyMsg{Type: tea.KeyTab})
                Expect(model.GetCurrentTab()).To(Equal(i % 4))
            }
        })
    })
})
```

### 按键响应问题测试
```go
Describe("按键响应功能修复验证", func() {
    Context("配置编辑器按键响应", func() {
        It("应该响应除Esc和Ctrl+C外的其他按键", func() {
            model := NewConfigEditorModel(&types.Config{})
            
            // 测试Tab键
            _, cmd := model.Update(tea.KeyMsg{Type: tea.KeyTab})
            Expect(cmd).ToNot(BeNil())
            
            // 测试上下键
            _, cmd = model.Update(tea.KeyMsg{Type: tea.KeyUp})
            Expect(cmd).ToNot(BeNil())
            
            _, cmd = model.Update(tea.KeyMsg{Type: tea.KeyDown})
            Expect(cmd).ToNot(BeNil())
        })
    })
})
```

## 性能测试

### 大文件列表处理
```go
Describe("性能测试", func() {
    Context("大量文件处理", func() {
        It("应该高效处理大量文件列表", func() {
            model := NewFileSelectorModel("./large_dir")
            
            start := time.Now()
            // 模拟加载大量文件
            model.loadFiles()
            duration := time.Since(start)
            
            Expect(duration).To(BeNumerically("<", time.Second))
        })
    })
})
```

### 内存使用测试
```go
Describe("内存使用", func() {
    Context("内存泄漏检查", func() {
        It("不应该存在内存泄漏", func() {
            initialMem := getMemoryUsage()
            
            // 多次创建和销毁模型
            for i := 0; i < 100; i++ {
                model := NewFileSelectorModel("./test")
                model = nil // 释放引用
            }
            
            finalMem := getMemoryUsage()
            memoryGrowth := finalMem - initialMem
            
            Expect(memoryGrowth).To(BeNumerically("<", 1024*1024)) // 1MB
        })
    })
})
```

## 测试执行指南

### 运行所有测试
```bash
cd tests/tui
go test -v ./...
```

### 运行特定测试
```bash
# 运行主模型测试
go test -v -run TestMainModel

# 运行文件选择器测试
go test -v -run TestFileSelector

# 运行集成测试
go test -v -run TestIntegration
```

### 覆盖率测试
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 基准测试
```bash
go test -bench=. -benchmem
```

## 测试最佳实践

### 1. 测试隔离
- 每个测试用例应该独立运行
- 避免测试间的状态依赖
- 使用Setup和Teardown清理环境

### 2. 测试数据管理
- 使用专用的测试数据目录
- 创建临时文件和目录
- 测试完成后清理测试数据

### 3. 断言清晰
- 使用描述性的断言消息
- 验证所有相关的状态变化
- 测试正向和负向场景

### 4. 模拟和桩
- 使用接口模拟外部依赖
- 创建测试用的模拟数据
- 避免在单元测试中访问真实文件系统

## 持续集成

### 自动化测试配置
```yaml
# .github/workflows/test.yml
name: TUI Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
    - run: go test -v ./tests/tui/...
    - run: go test -coverprofile=coverage.out ./tests/tui/...
    - uses: codecov/codecov-action@v1
```

### 测试报告
- 生成HTML格式的测试报告
- 集成代码覆盖率报告
- 设置测试失败通知