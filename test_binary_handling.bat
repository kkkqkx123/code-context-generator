@echo off
setlocal enabledelayedexpansion

echo 正在测试二进制文件处理功能...
echo.

REM 创建测试目录结构
echo 创建二进制测试文件...
mkdir test_binary_files 2>nul
echo 这是文本文件 > test_binary_files\text_file.txt
echo 这是Go源代码 > test_binary_files\source_code.go
echo 这是Python脚本 > test_binary_files\script.py

REM 创建真正的二进制文件
echo 创建二进制文件...
echo ÿÿ > test_binary_files\binary_file.bin
echo 这是PDF文件头 > test_binary_files\document.pdf
echo MZ > test_binary_files\executable.exe

echo 测试文件创建完成，开始测试二进制文件处理...
echo.

REM 测试1: 默认行为（排除二进制文件）
echo 1. 测试默认行为（排除二进制文件）:
echo    命令: go run cmd/cli/main.go generate test_binary_files -f json -o test_output_default.json
go run cmd/cli/main.go generate test_binary_files -f json -o test_output_default.json
if !errorlevel! equ 0 (
    echo    成功执行
) else (
    echo    执行失败
)
echo.

REM 测试2: 显式排除二进制文件
echo 2. 测试显式排除二进制文件:
echo    命令: go run cmd/cli/main.go generate test_binary_files --exclude-binary=true -f markdown -o test_output_exclude.md
go run cmd/cli/main.go generate test_binary_files --exclude-binary=true -f markdown -o test_output_exclude.md
if !errorlevel! equ 0 (
    echo    成功执行
) else (
    echo    执行失败
)
echo.

REM 测试3: 包含二进制文件
echo 3. 测试包含二进制文件:
echo    命令: go run cmd/cli/main.go generate test_binary_files --exclude-binary=false -f xml -o test_output_include.xml
go run cmd/cli/main.go generate test_binary_files --exclude-binary=false -f xml -o test_output_include.xml
if !errorlevel! equ 0 (
    echo    成功执行
) else (
    echo    执行失败
)
echo.

echo 测试结果分析:
echo.
echo 检查输出文件中的二进制文件处理情况:
echo.

echo 默认输出 (JSON):
if exist test_output_default.json (
    findstr /C:"binary_file.bin" test_output_default.json >nul
    if !errorlevel! equ 0 (
        echo    ❌ 默认输出包含二进制文件（应该排除）
    ) else (
        echo    ✅ 默认输出正确排除二进制文件
    )
) else (
    echo    ❌ 默认输出文件不存在
)

echo.
echo 排除二进制文件 (Markdown):
if exist test_output_exclude.md (
    findstr /C:"binary_file.bin" test_output_exclude.md >nul
    if !errorlevel! equ 0 (
        echo    ❌ 排除输出包含二进制文件（应该排除）
    ) else (
        echo    ✅ 排除输出正确排除二进制文件
    )
) else (
    echo    ❌ 排除输出文件不存在
)

echo.
echo 包含二进制文件 (XML):
if exist test_output_include.xml (
    findstr /C:"binary_file.bin" test_output_include.xml >nul
    if !errorlevel! equ 0 (
        echo    ✅ 包含输出包含二进制文件（正确）
    ) else (
        echo    ❌ 包含输出不包含二进制文件（应该包含）
    )
) else (
    echo    ❌ 包含输出文件不存在
)

echo.
echo 检查文件统计信息:
echo.

REM 统计目录中的文件数量
set total_files=0
for %%f in (test_binary_files\*) do set /a total_files+=1
echo 目录中的总文件数: !total_files!

echo.
if exist test_output_default.json (
    for /f "tokens=3 delims=:" %%a in ('findstr /C:"file_count" test_output_default.json') do (
        set json_count=%%a
        set json_count=!json_count:,=!
        set json_count=!json_count: =!
        echo JSON输出中的文件数: !json_count!
        
        if !json_count! lss !total_files! (
            echo    ✅ JSON输出文件数少于总数（可能排除了二进制文件）
        ) else (
            echo    ❌ JSON输出文件数等于总数（可能没有正确过滤）
        )
    )
)

echo.
if exist test_output_include.xml (
    for /f "tokens=3 delims=:" %%a in ('findstr /C:"file_count" test_output_include.xml') do (
        set xml_count=%%a
        set xml_count=!xml_count:,=!
        set xml_count=!xml_count: =!
        echo XML输出中的文件数: !xml_count!
        
        if !xml_count! equ !total_files! (
            echo    ✅ XML输出文件数等于总数（正确包含所有文件）
        ) else (
            echo    ❌ XML输出文件数不等于总数（可能有过滤问题）
        )
    )
)

echo.
echo 测试完成！
echo 输出文件:
echo - test_output_default.json （默认行为）
echo - test_output_exclude.md （显式排除二进制文件）
echo - test_output_include.xml （包含二进制文件）
echo.

set /p cleanup=是否清理测试文件？(Y/N) 
if /i "!cleanup!"=="Y" (
    echo 清理测试文件...
    rmdir /s /q test_binary_files 2>nul
    del test_output_default.json 2>nul
    del test_output_exclude.md 2>nul
    del test_output_include.xml 2>nul
    echo 测试文件已清理。
)

endlocal