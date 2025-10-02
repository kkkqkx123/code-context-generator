@echo off
REM 代码上下文生成器 - Windows使用示例脚本
REM 该脚本展示了如何使用不同的配置文件和参数

echo 🚀 代码上下文生成器 - Windows使用示例
echo ==================================

REM 检查是否已编译工具
if not exist "code-context-generator.exe" (
    echo ❌ 请先编译代码上下文生成器:
    echo    go build -o code-context-generator.exe cmd\cli\main.go
    exit /b 1
)

REM 创建输出目录
if not exist "output" mkdir output
if not exist "output\logs" mkdir output\logs

echo.
echo 1️⃣ 基础扫描 - 使用默认配置
echo ------------------------
code-context-generator.exe generate -o output\basic-scan.json

echo.
echo 2️⃣ 使用基础配置文件
echo ------------------
code-context-generator.exe generate -c examples\basic-config.toml -o output\basic-config-output.json

echo.
echo 3️⃣ 生成项目文档（包含文件内容）
echo --------------------------------
code-context-generator.exe generate -c examples\project-documentation.toml -o output\project-documentation.md

echo.
echo 4️⃣ 高性能扫描（适合大项目）
echo ----------------------------
code-context-generator.exe generate -c examples\performance-optimized.toml -o output\fast-scan.json

echo.
echo 5️⃣ 自定义参数扫描
echo ----------------
code-context-generator.exe generate -f xml -e "node_modules" -e ".git" -e "*.log" -s 1048576 -d 3 -o output\custom-scan.xml

echo.
echo 6️⃣ 自动文件扫描
echo ----------------
echo 📝 这将自动扫描当前目录...
code-context-generator.exe select -m -f markdown -o output\selected-files.md

echo.
echo 7️⃣ 生成配置文件
echo --------------
code-context-generator.exe config init -o output\my-config.toml

echo.
echo 8️⃣ 验证配置文件
echo --------------
code-context-generator.exe config validate -c examples\basic-config.toml

echo.
echo 9️⃣ 显示当前配置
echo --------------
code-context-generator.exe config show

echo.
echo 🔟 性能测试
echo ----------
echo 📊 扫描性能测试...
echo 开始时间: %date% %time%
code-context-generator.exe generate -c examples\performance-optimized.toml -o output\performance-test.json
echo 结束时间: %date% %time%

echo.
echo 📋 批处理示例 - 扫描常见目录
echo =============================

REM 扫描src目录（如果存在）
if exist "src" (
    echo 📁 扫描目录: src
    code-context-generator.exe generate src -f json -e "*.log" -e "*.tmp" -o output\scan_src.json
)

REM 扫描internal目录（如果存在）
if exist "internal" (
    echo 📁 扫描目录: internal
    code-context-generator.exe generate internal -f json -e "*.log" -e "*.tmp" -o output\scan_internal.json
)

REM 扫描pkg目录（如果存在）
if exist "pkg" (
    echo 📁 扫描目录: pkg
    code-context-generator.exe generate pkg -f json -e "*.log" -e "*.tmp" -o output\scan_pkg.json
)

REM 扫描cmd目录（如果存在）
if exist "cmd" (
    echo 📁 扫描目录: cmd
    code-context-generator.exe generate cmd -f json -e "*.log" -e "*.tmp" -o output\scan_cmd.json
)

echo.
echo 🔄 定时任务示例
echo =============

REM 创建定时任务脚本
echo @echo off > output\scheduled-scan.bat
echo REM 定时扫描脚本 >> output\scheduled-scan.bat
echo set DATE=%%date:~-4,4%%%%date:~-10,2%%%%date:~-7,2%%_%%time:~0,2%%%%time:~3,2%%%%time:~6,2%% >> output\scheduled-scan.bat
echo set DATE=%%DATE: =0%% >> output\scheduled-scan.bat
echo for %%%%i in ("%%cd%%") do set PROJECT_NAME=%%%%~nxi >> output\scheduled-scan.bat
echo. >> output\scheduled-scan.bat
echo code-context-generator.exe generate -c examples\project-documentation.toml -o "backup\%%PROJECT_NAME%%_documentation_%%DATE%%.md" >> output\scheduled-scan.bat
echo. >> output\scheduled-scan.bat
echo echo ✅ 备份完成: backup\%%PROJECT_NAME%%_documentation_%%DATE%%.md >> output\scheduled-scan.bat

echo ✅ 定时任务脚本已创建: output\scheduled-scan.bat
echo    可以添加到Windows任务计划程序中实现定时备份

echo.
echo 📊 结果统计
echo ==========
echo 生成的文件:
dir output\*.json output\*.xml output\*.md output\*.toml 2>nul

echo.
echo 文件大小统计:
du -h output\* 2>nul | sort /R

REM 如果没有du命令，使用替代方案
if %errorlevel% neq 0 (
    echo 使用dir命令显示文件大小:
    dir output\ /-C | findstr /R "^[0-9].*[0-9]$"
)

echo.
echo ✨ 示例完成！
echo ============
echo 📁 输出文件保存在: output\
echo 📝 日志文件保存在: output\logs\
echo.
echo 💡 提示:
echo    - 使用 '-c' 参数指定配置文件
echo    - 使用 '-f' 参数指定输出格式
echo    - 使用 '-e' 参数排除文件/目录
echo    - 使用 '-s' 参数限制文件大小
echo    - 使用 '-d' 参数限制扫描深度
echo    - 使用 '--debug' 参数启用调试模式
echo.
echo 📚 更多帮助:
echo    code-context-generator.exe --help
echo    code-context-generator.exe generate --help
echo    type docs\quickstart.md

echo.
pause