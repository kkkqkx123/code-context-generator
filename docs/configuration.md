# 配置文档

## 概述

代码上下文生成器支持通过配置文件进行灵活的格式和输出设置。配置文件统一使用YAML格式，支持配置YAML、JSON、MarkDown和TOML格式的输出。

## 配置文件格式覆盖功能

### 智能格式识别

系统会根据配置文件名自动识别并应用相应的格式配置：

- **配置文件名包含"xml"** → 自动应用XML格式配置，默认输出格式设为XML
- **配置文件名包含"json"** → 自动应用JSON格式配置，默认输出格式设为JSON  
- **配置文件名包含"toml"** → 自动应用TOML格式配置，默认输出格式设为TOML
- **配置文件名包含"markdown"或"md"** → 自动应用Markdown格式配置，默认输出格式设为Markdown

### 使用示例

```bash
# 使用XML配置（自动应用XML格式）
./test-build.exe -c config-xml.yaml generate

# 使用JSON配置（自动应用JSON格式）
./test-build.exe -c config-json.yaml generate

# 使用TOML配置（自动应用TOML格式）
./test-build.exe -c config-toml.yaml generate
```

## 配置结构

### 格式配置（Formats）

#### XML格式
```yaml
formats:
  xml:
    enabled: true
    root_tag: "context"
    file_tag: "file"
    files_tag: "files"
    folder_tag: "folder"
    fields:
      path: "path"
      content: "content"
      filename: "filename"
    formatting:
      indent: "  "
      declaration: true
      encoding: "UTF-8"
```

#### JSON格式
```yaml
formats:
  json:
    enabled: true
    structure:
      file: "source_file"
      folder: "directory"
      files: "all_files"
    fields:
      path: "relative_path"
      content: "file_content"
      filename: "file_name"
```

#### TOML格式
```yaml
formats:
  toml:
    enabled: true
    structure:
      file_section: "file"
      folder_section: "folder"
    fields:
      path: "path"
      content: "content"
      filename: "filename"
```

#### Markdown格式
```yaml
formats:
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
```

### 输出配置（Output）

```yaml
output:
  default_format: "xml"  # 会被配置文件名自动覆盖
  output_dir: "output"
  filename_template: "context_{{.timestamp}}.{{.extension}}"
  timestamp_format: "20060102_150405"
```
时间戳为参考时间，用于生成唯一的文件名。

### 文件处理配置（FileProcessing）

```yaml
file_processing:
  include_hidden: false
  include_content: true
  include_hash: false
```

### 字段配置（Fields）

```yaml
fields:
  custom_names:
    filepath: "path"
    filecontent: "content"
    filename: "name"
  filter:
    include: []
    exclude: []
  processing:
    max_length: 0
    add_line_numbers: false
    trim_whitespace: true
    code_highlight: false
```

### 过滤配置（Filters）

```yaml
filters:
  max_file_size: "10MB"
  exclude_patterns:
    - "*.tmp"
    - "*.log"
    - "node_modules"
  include_patterns: []
  max_depth: 0
  follow_symlinks: false
  exclude_binary: true
```

### 安全扫描配置（Security）

```yaml
security:
  enabled: false  # 是否启用安全检查（默认禁用）
  fail_on_critical: false  # 是否在发现严重问题时失败
  scan_level: "standard"  # 扫描级别：basic/standard/comprehensive
  report_format: "text"  # 安全报告格式：text/json/xml/html
  detectors:
    credentials: false  # 是否检测硬编码凭证
    sql_injection: false  # 是否检测SQL注入
    xss: false  # 是否检测XSS漏洞
    path_traversal: false  # 是否检测路径遍历
    quality: false  # 是否检测代码质量问题
  exclusions:
    files: []  # 排除的文件列表
    patterns: []  # 排除的文件模式
    rules: []  # 排除规则
  reporting:
    format: "text"  # 报告格式
    output_file: ""  # 输出文件路径
    include_details: true  # 包含详细问题信息
    show_statistics: true  # 显示扫描统计信息
```



## 环境变量配置

系统支持通过环境变量覆盖配置：

### 基本配置
- `CONTEXT_DEFAULT_FORMAT`: 默认输出格式
- `CONTEXT_OUTPUT_DIR`: 输出目录
- `CONTEXT_FILENAME_TEMPLATE`: 文件名模板
- `CONTEXT_TIMESTAMP_FORMAT`: 时间戳格式
- `CONTEXT_MAX_FILE_SIZE`: 最大文件大小
- `CONTEXT_MAX_DEPTH`: 最大深度
- `CONTEXT_RECURSIVE`: 是否递归
- `CONTEXT_INCLUDE_HIDDEN`: 是否包含隐藏文件
- `CONTEXT_FOLLOW_SYMLINKS`: 是否跟随符号链接
- `CONTEXT_EXCLUDE_BINARY`: 是否排除二进制文件
- `CONTEXT_EXCLUDE_PATTERNS`: 排除模式（逗号分隔）

### 安全扫描配置（默认禁用）
- `CONTEXT_SECURITY_ENABLED`: 是否启用安全检查（默认：false）
- `CONTEXT_SECURITY_SCAN_LEVEL`: 扫描级别（basic/standard/comprehensive，默认：standard）
- `CONTEXT_SECURITY_FAIL_ON_CRITICAL`: 是否在发现严重问题时失败（默认：false）
- `CONTEXT_SECURITY_REPORT_FORMAT`: 安全报告格式（text/json/xml，默认：text）
- `CONTEXT_SECURITY_DETECT_CREDENTIALS`: 是否检测硬编码凭证（默认：false）
- `CONTEXT_SECURITY_DETECT_SQL_INJECTION`: 是否检测SQL注入（默认：false）
- `CONTEXT_SECURITY_DETECT_XSS`: 是否检测XSS漏洞（默认：false）
- `CONTEXT_SECURITY_DETECT_PATH_TRAVERSAL`: 是否检测路径遍历（默认：false）
- `CONTEXT_SECURITY_DETECT_QUALITY`: 是否检测代码质量问题（默认：false）


## 完整配置示例

### XML配置文件（config-xml.yaml）
```yaml
formats:
  xml:
    enabled: true
  json:
    enabled: false
  toml:
    enabled: false
  markdown:
    enabled: false

output:
  output_dir: xml_output
  filename_template: "project_{{.timestamp}}.xml"

filters:
  max_file_size: "10MB"
  exclude_patterns:
    - "*.tmp"
    - "*.log"
    - "node_modules"
```

### JSON配置文件（config-json.yaml）
```yaml
formats:
  xml:
    enabled: false
  json:
    enabled: true
    structure:
      file: source_file
      folder: directory
      files: all_files
    fields:
      path: relative_path
      content: file_content
      filename: file_name
  toml:
    enabled: false
  markdown:
    enabled: false

output:
  output_dir: json_output
  filename_template: "context_{{.timestamp}}.json"

filters:
  max_file_size: "2MB"
  exclude_patterns:
    - "*.tmp"
    - "*.log"
    - "*.swp"
    - ".*"
    - "node_modules/"
    - "target/"
    - "dist/"
    - "build/"
```

## 配置文件加载规则

1. 系统首先加载配置文件内容
2. 应用环境变量覆盖
3. 根据配置文件名应用格式特定的配置覆盖
4. 最终配置 = 基础配置 + 环境变量覆盖 + 格式特定覆盖

## 注意事项

1. 由于go语言的XML解析不支持复杂结构，XML配置文件格式暂不支持，请使用YAML、JSON或TOML格式
2. 格式自动识别基于配置文件名，不依赖于文件扩展名
3. 如果对应格式在配置中未启用，格式覆盖将不会生效
4. 环境变量配置的优先级高于配置文件但低于命令行参数