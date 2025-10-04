# AI优化输出格式扩展计划

## 功能概述

基于repomix的AI优化输出格式，增强当前项目的输出格式，使其更适合AI工具处理。主要改进包括XML标签结构、AI友好的文档组织方式等。

## 当前状态分析

### 现有输出格式
- JSON格式：结构化但缺乏AI优化
- XML格式：基本结构但缺少AI特定标签
- Markdown格式：可读性好但AI解析复杂
- TOML格式：配置友好但AI兼容性差

### 与repomix对比
- ❌ 缺少AI导向的文件摘要
- ❌ 缺少XML标签层次结构
- ❌ 缺少AI使用说明
- ❌ 缺少自定义指令支持

## 扩展目标

### 核心功能
1. **AI优化XML格式** - 类似repomix的XML结构
2. **智能文件摘要** - AI友好的项目概述
3. **自定义指令支持** - 用户可配置AI提示
4. **Token计数集成** - 与Token计算模块集成

### 输出格式增强
- XML格式：增加AI特定标签
- Markdown格式：优化AI可读性
- JSON格式：保持兼容性
- 新增：AI专用格式

## 技术实现方案

### 1. XML格式增强

#### 当前XML结构
```xml
<project>
  <files>
    <file path="path/to/file">
      <content>file content</content>
    </file>
  </files>
</project>
```

#### 目标XML结构
```xml
<file_summary>
  <generation_info>
    <tool>code-context-generator</tool>
    <version>1.0.0</version>
    <timestamp>2024-01-01T00:00:00Z</timestamp>
  </generation_info>
  <ai_instructions>
    <purpose>This file contains a packed representation of the entire repository...</purpose>
    <usage_guidelines>
      - This file should be treated as read-only...
      - Use for AI analysis, code review, or documentation generation
    </usage_guidelines>
  </ai_instructions>
</file_summary>

<directory_structure>
  <!-- 目录树结构 -->
</directory_structure>

<files>
  <file path="src/main.go">
    <metadata>
      <size>1024</size>
      <lines>50</lines>
      <tokens>200</tokens>
      <language>go</language>
    </metadata>
    <content>
      <!-- 文件内容 -->
    </content>
  </file>
</files>

<instruction>
  <!-- 自定义指令内容 -->
</instruction>
```

### 2. 智能文件摘要生成

#### 摘要内容
- 项目概述和目的
- AI使用指南
- 文件组织说明
- Token使用建议

#### 实现逻辑
```go
type AISummary struct {
    GenerationHeader string
    Purpose          string
    FileFormat       string
    UsageGuidelines []string
    Notes           []string
}

func GenerateAISummary(config *Config, fileCount int, totalTokens int) AISummary {
    // 生成AI友好的摘要信息
}
```

### 3. 自定义指令支持

#### 功能设计
- 支持从文件读取自定义指令
- 支持内联指令配置
- 支持不同AI工具的特定指令

#### 配置示例
```yaml
output:
  ai_instructions:
    enabled: true
    file_path: ".ai-instructions.md"
    default_instructions: |
      Please analyze this codebase for:
      - Code quality issues
      - Performance optimization opportunities
      - Security vulnerabilities
```

## 实施步骤

### 第一阶段：核心功能开发（1-2周）
1. **XML格式重构**
   - 设计新的XML结构
   - 实现AI摘要生成
   - 更新XML格式化器

2. **智能摘要生成**
   - 开发摘要模板系统
   - 集成项目元数据
   - 实现多语言支持

### 第二阶段：高级功能（1周）
3. **自定义指令系统**
   - 实现指令文件读取
   - 支持模板变量替换
   - 提供预设指令库

4. **格式一致性**
   - 确保各格式间一致性
   - 优化AI可读性
   - 性能测试和优化

## 代码修改点

### 主要修改文件
1. `internal/formatter/formatter.go` - 核心格式化逻辑
2. `internal/formatter/xml/formatter.go` - XML特定格式化
3. `internal/formatter/markdown/formatter.go` - Markdown优化
4. `pkg/types/config.go` - 配置结构扩展

### 新增文件
1. `internal/formatter/ai_summary.go` - AI摘要生成器
2. `internal/formatter/instruction_loader.go` - 指令加载器
3. `internal/formatter/template_system.go` - 模板系统

## 配置变更

### 新增配置项
```yaml
output:
  ai_optimized: true
  ai_summary:
    enabled: true
    template: "default"  # default, minimal, detailed
  ai_instructions:
    enabled: true
    file_path: ""
    content: ""
```

## 测试策略

### 单元测试
- AI摘要生成测试
- XML结构验证
- 指令加载测试

### 集成测试
- 完整项目输出测试
- AI工具兼容性测试
- 性能基准测试

### 验收标准
- 生成的XML能被主流AI工具正确解析
- Token计数准确
- 自定义指令正常工作
- 性能影响<10%

## 风险评估

### 技术风险
- XML结构复杂性增加
- 性能影响评估

