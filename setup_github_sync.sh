#!/bin/bash

# GitHub 数据同步配置脚本
# 用途：快速配置 GitHub 同步功能

echo "=================================="
echo "GitHub 数据同步配置向导"
echo "=================================="
echo ""

# 检查 .env 文件是否存在
if [ -f .env ]; then
    echo "⚠️  发现现有 .env 文件"
    read -p "是否要覆盖现有配置？(y/n): " overwrite
    if [ "$overwrite" != "y" ]; then
        echo "❌ 配置已取消"
        exit 0
    fi
    cp .env .env.backup
    echo "✅ 已备份现有配置到 .env.backup"
fi

echo ""
echo "请输入以下信息："
echo ""

# 获取 GitHub Token
read -p "1. GitHub Personal Access Token (ghp_...): " github_token
if [ -z "$github_token" ]; then
    echo "❌ Token 不能为空"
    exit 1
fi

# 获取仓库地址
read -p "2. GitHub 仓库地址 (https://github.com/user/repo): " github_repo
if [ -z "$github_repo" ]; then
    echo "❌ 仓库地址不能为空"
    exit 1
fi

# 获取同步间隔
read -p "3. 同步间隔（秒，默认 300）: " sync_interval
if [ -z "$sync_interval" ]; then
    sync_interval=300
fi

echo ""
echo "=================================="
echo "正在生成配置文件..."
echo "=================================="

# 创建 .env 文件
cat > .env << EOF
# GitHub 数据同步配置
GITHUB_SYNC_TOKEN=$github_token
GITHUB_SYNC_REPO=$github_repo
GITHUB_SYNC_INTERVAL=$sync_interval
EOF

echo ""
echo "✅ 配置文件已生成！"
echo ""
echo "配置内容："
echo "-----------------------------------"
cat .env
echo "-----------------------------------"
echo ""
echo "📝 后续步骤："
echo "1. 重启应用程序"
echo "2. 查看日志确认同步服务启动"
echo "3. 等待 $sync_interval 秒后检查 GitHub 仓库"
echo ""
echo "⚠️  安全提醒："
echo "- 请确保仓库为私有"
echo "- 使用后及时撤销测试用的 Token"
echo "- 不要将 .env 文件提交到版本控制"
echo ""
echo "📚 详细文档："
echo "- 快速指南: README_GITHUB_SYNC_CN.md"
echo "- 完整文档: docs/GITHUB_SYNC.md"
echo ""