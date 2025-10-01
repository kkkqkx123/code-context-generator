#!/bin/bash

echo "正在测试二进制文件处理功能..."
echo

# 创建测试目录结构
mkdir -p test_binary_files
cd test_binary_files

# 创建文本文件
echo "这是一个文本文件" > text_file.txt
echo "包含一些文本内容" >> text_file.txt

# 创建二进制文件
echo "创建二进制测试文件..."
echo -e '\x4D\x5A' > binary_file.bin
# 创建一些二进制数据
dd if=/dev/zero of=binary_file.exe bs=1024 count=1 2>/dev/null
dd if=/dev/zero of=binary_file.dll bs=2048 count=1 2>/dev/null

# 创建更多文本文件
echo "这是另一个文本文件" > another_text.txt
echo 'function hello() { console.log("Hello"); }' > script.js

# 返回上级目录
cd ..

echo
echo "测试文件创建完成，开始测试二进制文件处理..."
echo

# 测试默认行为（排除二进制文件）
echo "1. 测试默认行为（排除二进制文件）:"
echo "   命令: ./code-context-generator generate test_binary_files -f json -o test_output_default.json"
./code-context-generator generate test_binary_files -f json -o test_output_default.json

echo
echo "2. 测试显式排除二进制文件:"
echo "   命令: ./code-context-generator generate test_binary_files --exclude-binary=true -f markdown -o test_output_exclude.md"
./code-context-generator generate test_binary_files --exclude-binary=true -f markdown -o test_output_exclude.md

echo
echo "3. 测试包含二进制文件:"
echo "   命令: ./code-context-generator generate test_binary_files --exclude-binary=false -f xml -o test_output_include.xml"
./code-context-generator generate test_binary_files --exclude-binary=false -f xml -o test_output_include.xml

echo
echo "测试结果分析:"
echo
echo "检查输出文件中的二进制文件处理情况:"
echo
echo "默认输出 (JSON):"
grep -i "binary" test_output_default.json || echo "未找到二进制文件（符合预期）"
echo
echo "排除二进制文件 (Markdown):"
grep -i "binary" test_output_exclude.md || echo "未找到二进制文件（符合预期）"
echo
echo "包含二进制文件 (XML):"
grep -i "binary" test_output_include.xml || echo "未找到二进制文件（意外）"

echo
echo "检查文件统计信息:"
echo "目录中的总文件数:"
ls -1 test_binary_files | wc -l
echo
echo "JSON输出中的文件数:"
grep -o '"files": [0-9]*' test_output_default.json || echo "无法统计"
echo
echo "XML输出中的文件数:"
grep -c '<file>' test_output_include.xml || echo "无法统计"

echo
echo "测试完成！"
echo "输出文件:"
echo "- test_output_default.json (默认行为)"
echo "- test_output_exclude.md (显式排除二进制文件)"
echo "- test_output_include.xml (包含二进制文件)"

# 清理测试文件
echo
echo "是否清理测试文件？(Y/N)"
read -r cleanup
if [[ "$cleanup" =~ ^[Yy]$ ]]; then
    rm -rf test_binary_files
    rm -f test_output_default.json test_output_exclude.md test_output_include.xml
    echo "测试文件已清理。"
else
    echo "测试文件保留在当前目录。"
fi