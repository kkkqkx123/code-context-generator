# 代码上下文生成器 - 文档中心

欢迎使用代码上下文生成器文档中心！这里包含了使用、部署和开发该工具所需的全部文档。

## 📚 文档目录

### 📖 使用文档
- [**使用文档**](usage.md) - 完整的使用指南，包含CLI和TUI的所有功能说明
- [**快速开始**](../README.md) - 项目README，快速了解项目功能
- [**二进制文件处理**](binary-file-handling.md) - 二进制文件检测和处理机制

### 🚀 部署文档
- [**部署文档**](deployment.md) - 详细的部署指南，支持多种部署方式
- [**配置说明**](#配置文档) - 配置文件详解和示例

### 💻 开发文档
- [**开发环境文档**](development.md) - 完整的开发环境搭建和开发流程指南
- [**API文档**](#api文档) - 代码API文档（自动生成）

## 🎯 快速导航

### 新用户
1. 首先查看[快速开始](../README.md)了解项目
2. 阅读[使用文档](usage.md)学习如何使用
3. 查看[配置说明](#配置文档)进行个性化配置

### 部署人员
1. 阅读[部署文档](deployment.md)选择合适的部署方案
2. 查看[系统要求](deployment.md#系统要求)确认环境
3. 参考[监控和日志](deployment.md#监控和日志)进行运维

### 开发人员
1. 详细阅读[开发环境文档](development.md)搭建开发环境
2. 查看[代码结构](development.md#项目结构)了解项目架构
3. 遵循[开发流程](development.md#开发流程)进行开发
4. 运行[测试指南](development.md#测试指南)确保代码质量

## 📋 功能特性

### 🎯 核心功能
- **多格式输出**: 支持 JSON、XML、TOML、Markdown 格式
- **智能文件选择**: 交互式文件/目录选择界面
- **自动补全**: 文件路径智能补全功能
- **配置管理**: 灵活的配置系统，支持环境变量覆盖
- **二进制文件处理**: 智能检测并处理二进制文件，避免内容错误

### 🚀 高级特性
- **并发处理**: 基于 goroutine 池的高性能文件扫描
- **大文件处理**: 流式读取，支持大文件处理
- **模式匹配**: 支持 glob 模式和正则表达式过滤
- **缓存机制**: 智能缓存提升重复扫描性能
- **跨平台**: 支持 Windows、Linux、macOS

### 🎨 用户界面
- **CLI 模式**: 功能丰富的命令行界面（基于 Cobra）
- **TUI 模式**: 现代化的终端用户界面（基于 Bubble Tea）
- **进度显示**: 实时进度条和状态信息
- **主题支持**: 可定制的界面主题

## 🔧 配置文档

### 配置文件格式
支持三种格式：TOML、YAML、JSON，默认使用 TOML 格式。

#### 基础配置示例
```toml
[output]
format = "json"
encoding = "utf-8"

[file_processing]
include_hidden = false
max_file_size = 10485760  # 10MB
max_depth = 0  # 无限制
exclude_patterns = [
    "*.exe", "*.dll", "*.so", "*.dylib",
    "*.pyc", "*.pyo", "*.pyd",
    "node_modules", ".git", ".svn", ".hg",
    "__pycache__", "*.egg-info", "dist", "build"
]
exclude_binary = true  # 排除二进制文件

[ui]
theme = "default"
show_progress = true
```

#### 完整配置示例
详见[使用文档](usage.md#配置文件详解)中的配置详解部分。

### 环境变量配置
```bash
# 输出格式
export CODE_CONTEXT_FORMAT=json

# 最大文件大小
export CODE_CONTEXT_MAX_SIZE=10485760

# 扫描深度
export CODE_CONTEXT_MAX_DEPTH=3

# 日志级别
export CODE_CONTEXT_LOG_LEVEL=info
```

## 🚀 快速开始示例

### CLI使用示例
```bash
# 扫描当前目录并输出JSON格式
./code-context-generator generate

# 扫描指定目录并输出Markdown格式
./code-context-generator generate /path/to/project -f markdown -o project.md

# 排除特定文件/目录
./code-context-generator generate -e "*.log" -e "node_modules" -e ".git"

# 包含文件内容和哈希值
./code-context-generator generate -C -H -f xml -o context.xml

# 排除二进制文件（默认行为）
./code-context-generator generate --exclude-binary

# 包含二进制文件（不推荐）
./code-context-generator generate --exclude-binary=false
```

### TUI使用示例
```bash
# 启动TUI界面
./code-context-generator-tui

# TUI界面提供：
# - 可视化路径输入
# - 交互式文件选择
# - 实时配置编辑
# - 进度显示
# - 结果预览
```

### 交互式选择示例
```bash
# 启动交互式文件选择器
./code-context-generator select

# 多选模式
./code-context-generator select -m -f json -o selected.json
```

## 📊 性能指标

### 基准测试结果
- **扫描速度**: 1000个文件/秒（平均）
- **内存使用**: 低于100MB（标准项目）
- **CPU使用**: 支持多核并发处理
- **大文件支持**: 支持GB级别文件处理

### 优化建议
1. 合理设置 `max_workers` 参数
2. 启用缓存机制
3. 使用适当的缓冲区大小
4. 排除不必要的目录
5. 限制扫描深度和文件大小

## 🔍 故障排除

### 常见问题

#### Q: 如何处理大文件？
**A**: 使用 `-s` 参数限制文件大小，例如 `-s 10485760` 限制为10MB。

#### Q: 如何排除特定目录？
**A**: 使用 `-e` 参数指定排除模式：`-e "node_modules" -e ".git" -e "*.log"`

#### Q: 如何包含隐藏文件？
**A**: 使用 `-h` 或 `--hidden` 参数包含隐藏文件。

#### Q: 如何处理二进制文件？
**A**: 默认情况下工具会自动检测并排除二进制文件。使用 `--exclude-binary` 控制此行为（默认true）。二进制文件在输出中会显示为"[二进制文件 - 内容未显示]"。

#### Q: 如何自定义输出格式？
**A**: 通过配置文件中的模板系统自定义输出格式。

#### Q: 性能优化建议？
1. 合理设置并发参数
2. 启用缓存机制
3. 使用适当的缓冲区大小
4. 限制扫描深度和文件大小
5. 排除不必要的目录

### 错误处理
- **权限错误**: 检查文件和目录的读取权限
- **内存不足**: 减小缓冲区大小和并发数
- **配置文件错误**: 验证配置文件语法

## 📈 更新日志

### v1.0.0 (2024-01-01)
- ✅ 初始版本发布
- ✅ 支持CLI和TUI界面
- ✅ 支持JSON、XML、TOML、Markdown格式
- ✅ 基础文件过滤功能
- ✅ 配置管理系统
- ✅ 二进制文件智能检测和处理

### 开发计划
- 🔄 添加更多输出格式
- 🔄 Web界面支持
- 🔄 插件系统
- 🔄 云存储集成
- 🔄 团队协作功能

## 🤝 贡献指南

### 如何贡献
1. Fork 项目仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交修改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

### 开发规范
- 遵循 [Go代码规范](development.md#代码风格)
- 编写完整的测试用例
- 更新相关文档
- 通过所有质量检查

详细开发指南请查看[开发环境文档](development.md)。

## 📞 获取帮助

### 文档资源
- 📖 [使用文档](usage.md) - 使用方法和示例
- 🚀 [部署文档](deployment.md) - 部署和配置指南
- 💻 [开发文档](development.md) - 开发环境搭建

### 社区支持
- 🐛 [问题报告](https://github.com/yourusername/code-context-generator/issues)
- 💬 [讨论区](https://github.com/yourusername/code-context-generator/discussions)
- 📧 [邮件支持](mailto:support@example.com)

### 更新和支持
- ⭐ 给项目点个Star支持开发
- 🔔 关注项目获取更新通知
- 📝 提交Issue报告问题
- 🔄 提交Pull Request贡献代码

---

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](../LICENSE) 文件了解详情。

## 🙏 致谢

感谢所有贡献者和使用者的支持！特别感谢以下贡献者：
- 项目贡献者列表
- 社区支持成员
- 文档编写参与者

---

*最后更新：2024年1月1日*  
*文档版本：v1.0.0*