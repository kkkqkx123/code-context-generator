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
1. 必须能够处理中文路径、文件名
2. 必须拥有高性能
3. 采取简单的TUI设计，避免任何复杂的UI设计，以免影响性能，降低可靠性
4. 必须支持windows、linux的文件系统。生成文件中的路径统一使用正斜杠（/）作为路径分隔符
5. 提供交互式选择功能和基于前缀匹配的自动补全功能，且该功能必须支持windows和linux两种环境
6. 必须正确忽略选中的文件夹中的隐藏文件（如.git, .vscode, node_modules等），且在遍历路径前读取.gitignore的规则，忽略这些文件与目录
7. 必须支持递归遍历子文件夹，且在遍历子文件夹时必须正确处理符号链接（symbolic link）
8. 是否遍历所有子目录(默认只遍历1层)、开启自动补全、符号链接功能需要支持在.env文件中配置。使用默认值均为false。
9. 支持在rule.xml中设定生成的xml文件的格式，默认格式为：
```xml
<context>
    <file>
        <path>relative/path/to/file</path>
        <content>file content</content>
    </file>
    <folder>
        <path>relative/path/to/folder</path>
        <files>
            <file>
                <filename>file1</filename>
                <path>relative/path/to/file1</path>
                <content>file1 content</content>
            </file>
        </files>
    </folder>
</context>
```
支持在rule.json中设定生成的json文件的格式，默认格式为：
```json
{
    "file": [
        {
            "path": "relative/path/to/file",
            "content": "file content"
        }
    ],
    "folder": [
        {
            "path": "relative/path/to/folder",
            "files": [
                {
                    "filename": "file1",
                    "path": "relative/path/to/file1",
                    "content": "file1 content"
                }
            ]
        }
    ]   
}
```
目录内的子目录即目录格式的嵌套。

可修改性：层次固定，但具体名称可以修改，内部可以新增键值
需要专门的模块用于解析2种规则的内容

支持在.env中选择使用哪种导出格式，默认格式为xml。

10. 支持在cli界面中临时选择使用哪种导出格式
11. 支持TUI界面、cli命令2种方式使用。配置项也应当支持在执行cli命令时通过参数的形式指定