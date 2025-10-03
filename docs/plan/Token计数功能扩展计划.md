# Token计数功能扩展计划

## 功能概述

为项目添加Token计数功能，支持多种Token化算法，提供文件级和项目级的Token统计，帮助用户管理LLM上下文限制。

## 当前状态分析

### 现有功能
- ❌ 无Token计数功能
- ❌ 无LLM上下文管理
- ❌ 无Token优化建议

### 与repomix对比
- ❌ 缺少Token计数显示
- ❌ 缺少Token树状视图
- ❌ 缺少Token优化功能

## 扩展目标

### 核心功能
1. **多算法Token计数** - 支持GPT、Claude、通用算法
2. **分层Token统计** - 文件、目录、项目级统计
3. **Token优化建议** - 基于Token使用情况
4. **Token树状视图** - 可视化Token分布

### 用户体验
- 实时Token计数显示
- Token使用警告
- 优化建议提示

## 技术实现方案

### 1. Token计算引擎

#### 支持的算法
- **GPT系列**: GPT-2, GPT-3, GPT-4 Tokenizer
- **Claude系列**: Anthropic Tokenizer
- **通用算法**: 基于字符/单词的简单计数
- **自定义**: 用户可配置算法

#### 核心接口设计
```go
type Tokenizer interface {
    CountTokens(text string) int
    Name() string
    Version() string
}

type TokenCounter struct {
    tokenizers map[string]Tokenizer
    defaultAlgorithm string
}

func (tc *TokenCounter) CountFileTokens(filePath string, algorithm string) (int, error)
func (tc *TokenCounter) CountProjectTokens(projectPath string, algorithm string) (map[string]int, error)
```

### 2. Token统计系统

#### 统计维度
- **文件级**: 单个文件的Token数量
- **目录级**: 目录下所有文件的Token总和
- **项目级**: 整个项目的Token总数
- **类型级**: 按文件类型统计

#### 数据结构
```go
type TokenStats struct {
    TotalTokens int                    `json:"total_tokens"`
    FileCount   int                    `json:"file_count"`
    AvgTokens   float64               `json:"avg_tokens"`
    MaxTokens   int                    `json:"max_tokens"`
    MinTokens   int                    `json:"min_tokens"`
    ByFile      map[string]int         `json:"by_file"`
    ByType      map[string]FileTypeStats `json:"by_type"`
    TreeView    *TokenTreeNode         `json:"tree_view"`
}

type FileTypeStats struct {
    FileCount int `json:"file_count"`
    TotalTokens int `json:"total_tokens"`
    AvgTokens float64 `json:"avg_tokens"`
}
```

### 3. Token树状视图

#### 功能特性
- 分层显示Token分布
- 支持最小Token阈值过滤
- 可交互的树状结构
- 导出为多种格式

#### 命令行界面
```bash
# 显示Token树状视图
c-gen --token-tree

# 只显示超过1000Token的文件/目录
c-gen --token-tree 1000

# 使用特定算法
c-gen --token-tree --token-algorithm claude

# 导出为JSON
c-gen --token-tree --output token-tree.json
```

### 4. Token优化建议

#### 建议类型
- **文件拆分**: 大文件拆分建议
- **代码压缩**: 冗余代码识别
- **格式优化**: AI友好格式建议
- **选择性包含**: 关键文件推荐

#### 建议算法
```go
type OptimizationSuggestion struct {
    Type        SuggestionType `json:"type"`
    FilePath    string         `json:"file_path"`
    CurrentTokens int          `json:"current_tokens"`
    ExpectedSavings int        `json:"expected_savings"`
    Description string         `json:"description"`
    Priority    PriorityLevel  `json:"priority"`
}

func GenerateSuggestions(stats TokenStats, maxContext int) []OptimizationSuggestion
```

## 实施步骤

### 第一阶段：基础功能（2周）
1. **Token计算引擎**
   - 实现基础Tokenizer接口
   - 集成GPT-2 Tokenizer
   - 实现简单字符计数算法

2. **基本统计功能**
   - 文件级Token计数
   - 项目级汇总统计
   - 命令行输出格式

### 第二阶段：高级功能（2周）
3. **多算法支持**
   - 集成Claude Tokenizer
   - 支持自定义算法
   - 算法性能优化

4. **树状视图功能**
   - 分层统计实现
   - 树状结构生成
   - 可视化输出格式

### 第三阶段：优化功能（1周）
5. **优化建议系统**
   - 建议算法实现
   - 优先级计算
   - 用户交互界面

6. **性能优化**
   - 缓存机制
   - 并行处理
   - 内存优化

## 代码架构

### 新增包结构
```
internal/tokenizer/
├── tokenizer.go          # 核心接口
├── gpt_tokenizer.go     # GPT系列
├── claude_tokenizer.go  # Claude系列
├── simple_tokenizer.go  # 简单算法
└── manager.go           # Token管理器

internal/tokenstats/
├── collector.go         # 统计收集器
├── analyzer.go          # 分析器
├── tree_builder.go      # 树状构建器
└── optimizer.go         # 优化器

pkg/types/
├── token_stats.go       # 统计数据结构
└── suggestions.go       # 建议数据结构
```

### 主要修改文件
1. `cli/main.go` - 添加Token相关命令行参数
2. `internal/formatter/formatter.go` - 集成Token统计输出
3. `internal/filesystem/walker.go` - 添加Token计数回调

## 配置变更

### 新增配置项
```yaml
token:
  enabled: true
  default_algorithm: "gpt2"  # gpt2, claude, simple
  algorithms:
    gpt2:
      enabled: true
    claude:
      enabled: false
      # Claude特定配置
    simple:
      enabled: true
  optimization:
    enabled: true
    max_context_size: 16000  # 默认LLM上下文大小
    suggestions_threshold: 0.8  # 建议阈值
```

## 命令行参数

### 新增参数
```bash
# Token计数相关
--token-count              # 显示Token计数
--token-algorithm          # 指定Token算法
--token-tree [threshold]   # 显示Token树状视图
--token-optimize           # 生成优化建议

# 输出控制
--show-tokens              # 在输出中包含Token计数
--token-warning-threshold  # Token警告阈值
```

## 测试策略

### 单元测试
- Token算法准确性测试
- 统计计算正确性测试
- 边界条件处理测试

### 集成测试
- 完整项目Token计数测试
- 多算法一致性测试
- 性能基准测试

### 验收测试
- 与repomix Token计数对比
- LLM工具兼容性测试
- 用户场景测试

## 性能考虑

### 优化策略
- **懒加载**: Tokenizer按需初始化
- **缓存**: 文件内容Token计数缓存
- **并行**: 多文件并行Token计数
- **流式**: 大文件流式处理

### 内存管理
- 控制并发数量
- 及时释放资源
- 监控内存使用

## 风险评估

### 技术风险
- Token算法准确性
- 性能影响
- 第三方依赖稳定性

### 缓解措施
- 多算法验证
- 性能监控
- 备用简单算法

## 成功指标

### 功能指标
- 支持3种以上Token算法
- Token计数准确率>99%
- 性能影响<15%

### 用户体验
- 响应时间<2秒（中等项目）
- 内存使用可控
- 命令行界面友好

## 后续扩展

### 短期扩展
- 更多LLM模型支持
- 实时Token监控
- 智能文件选择

### 长期扩展
- AI驱动的Token优化
- 预测性Token管理
- 集成代码分析工具