# 双入口点使用指南

## 概述

代码上下文生成器项目提供两个不同的程序入口点，分别针对不同的使用场景和用户需求设计。本指南帮助用户选择最适合的入口点，并提供最佳实践建议。

## 入口点对比

| 特性 | 根目录 main.go | cmd/cli/main.go |
|------|---------------|-----------------|
| **设计理念** | 简单直观，零学习成本 | 功能完整，专业级工具 |
| **使用复杂度** | ⭐（简单） | ⭐⭐⭐（中等） |
| **功能丰富度** | ⭐⭐（基础） | ⭐⭐⭐⭐⭐（完整） |
| **配置支持** | ⚠️（需要改进） | ✅（完整支持） |
| **自动化支持** | ⚠️（基础支持） | ✅（完整支持） |
| **交互体验** | ✅（强制交互） | ✅（可选交互） |

## 使用场景指南

### 🎯 何时使用根目录 main.go

#### 适用场景
- **快速体验**：第一次使用，想快速了解功能
- **简单需求**：只需要基本的代码上下文生成
- **交互偏好**：喜欢通过交互式界面选择文件和格式
- **临时使用**：偶尔使用，不想记忆复杂命令
- **教学演示**：向他人展示基本功能

#### 推荐命令
```bash
# 快速生成（进入交互模式）
go run main.go

# 指定格式和输出文件（跳过交互）
go run main.go -f markdown -o project-docs.md

# 查看帮助
go run main.go -h
```

#### 优点
- ✅ 零学习成本
- ✅ 强制交互，避免遗漏重要选择
- ✅ 代码量少，启动快速
- ✅ 适合新用户上手

#### 局限性
- ⚠️ 配置复用有限（需要改进）
- ⚠️ 高级功能缺失
- ⚠️ 自动化脚本支持有限

### 🎯 何时使用 cmd/cli/main.go

#### 适用场景
- **专业使用**：需要完整的功能集和配置管理
- **自动化集成**：CI/CD 流程、构建脚本
- **批量处理**：需要处理多个项目或复杂过滤规则
- **团队协作**：需要标准化配置和一致输出
- **高级需求**：需要递归控制、文件过滤、内容包含等高级功能

#### 推荐命令
```bash
# 基本生成（使用配置文件）
go run cmd/cli/main.go generate -f markdown -o docs.md

# 高级过滤（排除测试文件，限制深度）
go run cmd/cli/main.go generate \
  -f markdown \
  -o docs.md \
  --exclude-patterns "*_test.go,*.tmp" \
  --max-depth 3

# 配置管理
go run cmd/cli/main.go config show
go run cmd/cli/main.go config set output.format markdown

```

#### 优点
- ✅ 完整的配置体系支持
- ✅ 丰富的命令和参数选项
- ✅ 强大的过滤和处理能力
- ✅ 适合自动化和脚本集成
- ✅ 团队协作标准化

#### 学习成本
- ⚠️ 需要理解命令结构
- ⚠️ 参数选项较多

## 配置使用最佳实践

### 根目录 main.go（改进后）

1. **创建配置文件**（`config.yaml`）：
```yaml
# 文件处理设置
file_processing:
  recursive: true
  include_hidden: false

# 过滤规则
filters:
  max_depth: 0  # 0表示无限制
  max_file_size: 10485760  # 10MB
  exclude_patterns:
    - "*.log"
    - "*.tmp"
    - "node_modules"
    - ".git"
    - "*.exe"
    - "*.dll"
  include_patterns: []
  exclude_binary: true
  follow_symlinks: false

# 输出设置
output:
  format: markdown
  include_content: true
  max_content_size: 1048576  # 1MB
  include_metadata: true
```

2. **使用环境变量**：
```bash
# Windows (PowerShell)
$env:CCG_MAX_DEPTH="3"
$env:CCG_EXCLUDE_BINARY="true"

# Linux/Mac
export CCG_MAX_DEPTH=3
export CCG_EXCLUDE_BINARY=true
```

### cmd/cli/main.go

1. **配置文件位置**：
   - 默认：`config.yaml`（当前目录）
   - 自定义：`--config path/to/config.yaml`
   - 全局：`~/.code-context-generator/config.yaml`

2. **配置管理命令**：
```bash
# 查看当前配置
go run cmd/cli/main.go config show

# 设置配置项
go run cmd/cli/main.go config set filters.max_depth 3
go run cmd/cli/main.go config set output.format json

# 重置配置
go run cmd/cli/main.go config reset
```

## 迁移指南

### 从简单入口迁移到专业入口

如果你发现需要更多功能，可以平滑迁移：

1. **保持配置文件**：两个入口点使用相同的配置文件格式
2. **参数映射**：
   ```bash
   # 简单版本
   go run main.go -f markdown -o docs.md
   
   # 专业版本（等效）
   go run cmd/cli/main.go generate -f markdown -o docs.md
   ```

3. **逐步采用高级功能**：
   ```bash
   # 添加过滤
   go run cmd/cli/main.go generate -f markdown -o docs.md --exclude-patterns "*.log"
   
   # 使用配置
   go run cmd/cli/main.go generate --config my-config.yaml
   ```

## 常见问题

### Q: 两个入口点生成的输出有区别吗？

**A**: 在配置复用改进完成后，两个入口点使用相同的配置和处理逻辑，生成的输出应该是一致的。当前根目录版本需要改进配置加载。

### Q: 应该选择哪个入口点？

**A**: 
- **新手用户**：从根目录 `main.go` 开始
- **简单使用**：根目录 `main.go` 足够
- **专业需求**：使用 `cmd/cli/main.go`
- **自动化场景**：必须使用 `cmd/cli/main.go`

### Q: 可以同时使用两个入口点吗？

**A**: 可以，它们是完全独立的程序。但要注意：
- 使用相同的配置文件避免重复配置
- 输出文件路径避免冲突
- 理解各自的参数差异

### Q: 如何确保配置一致性？

**A**: 
1. 使用版本控制管理配置文件
2. 在团队内标准化配置模板
3. 使用 `cmd/cli/main.go` 的配置管理功能
4. 定期检查和同步配置

## 性能对比

| 场景 | 根目录 main.go | cmd/cli/main.go |
|------|---------------|-----------------|
| **启动时间** | 快（代码简单） | 稍慢（框架加载） |
| **小项目处理** | 高效 | 高效 |
| **大项目处理** | 中等 | 优化更好 |
| **内存使用** | 较低 | 中等 |
| **并发处理** | 基础 | 支持并发 |

## 总结建议

### 🟢 推荐使用根目录 main.go 的情况
- 个人项目文档生成
- 快速原型和验证
- 教学和演示场景
- 不喜欢复杂命令的用户

### 🟢 推荐使用 cmd/cli/main.go 的情况
- 企业级项目和团队协作
- 自动化构建和部署
- 需要精细控制的场景
- 专业开发者和高级用户

### 🔧 最佳实践
1. **新手起步**：从根目录 `main.go` 开始，快速上手
2. **需求增长**：逐步迁移到 `cmd/cli/main.go` 获得更多功能
3. **团队协作**：统一使用 `cmd/cli/main.go` 和标准化配置
4. **配置管理**：使用版本控制管理配置文件，确保一致性

通过合理选择入口点，可以在简单性和功能性之间找到最佳平衡，提高开发效率和文档质量。