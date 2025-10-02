# AGENTS.md

该项目需要构建一个简单的cli工具。该项目的目的是使用go语言实现一个高性能能方便地生成代码上下文的工具。

## 环境
windows11
需要兼容powershell和git bash

### 编程语言
go 1.24.5

## 项目目的
本项目的目的是使用go语言实现一个能方便的通过终端选择文件/文件夹，
并将选中的文件的相对路径与内容打包为结构化的文件（如xml/json/md等），快速整合文件内容，跨文件构建上下文，方便用户将多个文件的内容快速转为提示词。

## 项目功能
1. 能方便地通过终端选择文件/文件夹。
2. 能将选中的文件的相对路径与内容打包为单个xml/json/md文件，并输出到指定目录。如果不指定就输出到当前目录。

## 额外要求
1. 支持在cli界面中临时选择使用哪种导出格式
2. 支持CLI命令方式使用。配置项也应当支持在执行cli命令时通过参数的形式指定
**CLI参数支持**：
```bash
# 基本使用
code-context-generator --format xml --output output.xml

# 使用配置文件
code-context-generator --config config.yaml --format json

# 覆盖配置选项
code-context-generator --format markdown --max-depth 3 --exclude "*.log,*.tmp"

# 指定输出目录和文件名模板
code-context-generator --output-dir ./outputs --filename-template "project_{{.timestamp}}.md"

# 显示配置验证信息
code-context-generator --validate-config --config config.yaml
```

3. 必须能够处理中文路径、文件名
4. 必须拥有高性能
5. 必须支持windows、linux的文件系统。生成文件中的路径统一使用正斜杠（/）作为路径分隔符
6. 必须正确忽略选中的文件夹中的隐藏文件（如.git, .vscode, node_modules等），且在遍历路径前读取.gitignore的规则，忽略这些文件与目录
7. 必须支持递归遍历子文件夹，且在遍历子文件夹时必须正确处理符号链接（symbolic link）
8. 是否遍历所有子目录(默认只遍历1层)、符号链接功能需要支持在.env文件中配置。使用默认值均为false。
11. 采用统一配置文件方案，支持YAML、JSON、TOML三种配置文件格式，使用Go标准库进行解析和生成：

**配置文件格式（config.yaml）**：
```yaml
# 统一配置文件 - 支持多种输出格式
formats:
  xml:
    enabled: true
    structure:
      root: "context"
      file: "file"
      folder: "folder"
      files: "files"
    fields:
      path: "path"
      content: "content"
      filename: "filename"

  json:
    enabled: true
    structure:
      file: "file"
      folder: "folder"
      files: "files"
    fields:
      path: "path"
      content: "content"
      filename: "filename"
    formatting:
      indent: "  "
      sort_keys: false

  toml:
    enabled: true
    structure:
      file_section: "file"
      folder_section: "folder"
    fields:
      path: "path"
      content: "content"
      filename: "filename"

  markdown:
    enabled: true
    structure:
      file_header: "##"
      folder_header: "###"
      code_block: "```"
    formatting:
      separator: "\n\n"
      add_toc: false
      code_language: true

# 通用字段配置
fields:
  custom_names:
    filepath: "path"
    filecontent: "content"
    filename: "name"
  
  filter:
    include: []  # 只包含这些字段，空数组表示包含所有
    exclude: []  # 排除这些字段
  
  processing:
    max_length: 0  # 0表示不限制
    add_line_numbers: false
    trim_whitespace: true
    code_highlight: false

# 文件过滤配置
filters:
  max_file_size: "10MB"
  exclude_patterns:
    - "*.tmp"
    - "*.log"
    - "*.swp"
    - ".*"  # 隐藏文件
  include_patterns: []
  max_depth: 0  # 0表示无限制
  follow_symlinks: false

# 输出配置
output:
  default_format: "xml"
  output_dir: ""  # 空表示当前目录
  filename_template: "context_{{.timestamp}}.{{.extension}}"
  timestamp_format: "20060102_150405"


```

**配置说明**：
- 支持YAML、JSON、TOML三种配置文件格式，使用Go标准库进行解析
- 统一配置结构，支持多种输出格式的灵活配置
- 字段名称可自定义，支持字段过滤和内容预处理
- 配置文件可通过命令行参数指定，支持运行时切换格式
- 保持与原有rule.xml/rule.json的兼容性，支持平滑迁移

**Go标准库对齐**：
- YAML格式：使用`github.com/goccy/go-yaml`（兼容encoding/json接口）
- JSON格式：使用标准库`encoding/json`
- TOML格式：使用`github.com/BurntSushi/toml`（兼容encoding/json接口）
- XML格式：使用标准库`encoding/xml`

支持在.env中选择默认导出格式，也支持CLI参数临时指定格式

12. 支持在cli界面中临时选择使用哪种导出格式

**配置管理**：
- 提供统一的配置管理器（config_manager.go），支持配置文件的加载、解析和格式转换
- 支持配置验证和默认值处理
- 提供配置热重载功能（可选）
- 支持环境变量覆盖配置文件中的设置

**配置方案**：
- 统一使用新的YAML/JSON/TOML配置文件格式
- 不再支持原有的rule.xml/rule.json格式
- 提供一次性迁移工具，帮助用户从旧格式迁移到新格式

## 实现准则

1. 使用Go语言开发，利用其跨平台特性和丰富的标准库
2. 采用模块化设计，将文件处理、格式转换、配置管理等功能分离
3. 考虑使用并发处理来提高性能
4. 实现完善的错误处理和日志记录机制

