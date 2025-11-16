@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ==================================
echo GitHub 数据同步配置向导
echo ==================================
echo.

REM 检查 .env 文件是否存在
if exist .env (
    echo ⚠️  发现现有 .env 文件
    set /p overwrite="是否要覆盖现有配置？(y/n): "
    if /i not "!overwrite!"=="y" (
        echo ❌ 配置已取消
        exit /b 0
    )
    copy .env .env.backup >nul
    echo ✅ 已备份现有配置到 .env.backup
)

echo.
echo 请输入以下信息：
echo.

REM 获取 GitHub Token
set /p github_token="1. GitHub Personal Access Token (ghp_...): "
if "!github_token!"=="" (
    echo ❌ Token 不能为空
    exit /b 1
)

REM 获取仓库地址
set /p github_repo="2. GitHub 仓库地址 (https://github.com/user/repo): "
if "!github_repo!"=="" (
    echo ❌ 仓库地址不能为空
    exit /b 1
)

REM 获取同步间隔
set /p sync_interval="3. 同步间隔（秒，默认 300）: "
if "!sync_interval!"=="" set sync_interval=300

echo.
echo ==================================
echo 正在生成配置文件...
echo ==================================

REM 创建 .env 文件
(
echo # GitHub 数据同步配置
echo GITHUB_SYNC_TOKEN=!github_token!
echo GITHUB_SYNC_REPO=!github_repo!
echo GITHUB_SYNC_INTERVAL=!sync_interval!
) > .env

echo.
echo ✅ 配置文件已生成！
echo.
echo 配置内容：
echo -----------------------------------
type .env
echo -----------------------------------
echo.
echo 📝 后续步骤：
echo 1. 重启应用程序
echo 2. 查看日志确认同步服务启动
echo 3. 等待 !sync_interval! 秒后检查 GitHub 仓库
echo.
echo ⚠️  安全提醒：
echo - 请确保仓库为私有
echo - 使用后及时撤销测试用的 Token
echo - 不要将 .env 文件提交到版本控制
echo.
echo 📚 详细文档：
echo - 快速指南: README_GITHUB_SYNC_CN.md
echo - 完整文档: docs\GITHUB_SYNC.md
echo.
pause