# GitHub 数据同步功能说明

## 概述

本系统支持将渠道（Channels）、用户（Users）和令牌（Tokens）数据自动同步到 GitHub 私有仓库，实现数据的版本控制和备份。

**注意**：日志数据仍然使用 MySQL 数据库存储，不会同步到 GitHub。

## 功能特性

- ✅ 自动定时同步数据到 GitHub 私有仓库
- ✅ 支持启动时从 GitHub 加载数据
- ✅ 使用 JSON 格式存储数据
- ✅ 支持自定义同步间隔
- ✅ 自动清除敏感信息（如密码、访问令牌）
- ✅ 不影响日志系统（日志继续使用 MySQL）

## 配置步骤

### 1. 创建 GitHub 私有仓库

1. 登录 GitHub
2. 创建一个新的私有仓库（例如：`my-api-data`）
3. 记录仓库的 URL，格式为：`https://github.com/用户名/仓库名`

### 2. 生成 GitHub Personal Access Token (PAT)

1. 访问 GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. 点击 "Generate new token (classic)"
3. 设置 Token 名称（如：`api-data-sync`）
4. 选择权限范围：
   - ✅ `repo` (完整的仓库访问权限)
5. 点击 "Generate token"
6. **重要**：复制生成的 token（格式如：`ghp_xxxxxxxxxxxxxxxxxxxx`）

### 3. 配置环境变量

在项目根目录创建或编辑 `.env` 文件，添加以下配置：

```env
# GitHub 数据同步配置
GITHUB_SYNC_TOKEN=ghp_your_github_personal_access_token_here
GITHUB_SYNC_REPO=https://github.com/username/repository
GITHUB_SYNC_INTERVAL=300
```

**配置说明**：

- `GITHUB_SYNC_TOKEN`: GitHub Personal Access Token
- `GITHUB_SYNC_REPO`: GitHub 仓库地址（支持完整 URL 或 `owner/repo` 格式）
- `GITHUB_SYNC_INTERVAL`: 同步间隔（秒），默认 300 秒（5分钟）

### 4. 启动系统

配置完成后，重启应用程序。系统将自动：

1. 启动时从 GitHub 加载最新数据
2. 每隔指定间隔自动同步数据到 GitHub

## 数据文件说明

同步到 GitHub 的数据文件：

- `channels.json` - 渠道配置数据
- `users.json` - 用户数据（已移除密码和访问令牌）
- `tokens.json` - API 令牌数据

## 示例配置

以下是一个完整的配置示例：

```env
# GitHub 同步配置
GITHUB_SYNC_TOKEN=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
GITHUB_SYNC_REPO=https://github.com/kiki0501/keykey
GITHUB_SYNC_INTERVAL=300
```

## 安全建议

1. **使用私有仓库**：确保 GitHub 仓库设置为私有
2. **Token 安全**：
   - 不要将 `.env` 文件提交到版本控制
   - 定期更换 Personal Access Token
   - 使用完毕后及时撤销不需要的 Token
3. **权限最小化**：只授予必要的仓库访问权限
4. **敏感数据**：系统会自动过滤密码等敏感信息

## 日志说明

系统启动时会显示以下日志：

```
GitHub sync service initialized: owner/repo, interval: 5m0s
Starting GitHub auto sync...
Loading data from GitHub...
Loaded X channels from GitHub
Loaded X users from GitHub
Loaded X tokens from GitHub
Data loaded from GitHub successfully
```

定时同步时会显示：

```
Syncing data to GitHub...
Data synced to GitHub successfully
```

## 故障排除

### 问题：系统提示 "GitHub sync not configured"

**解决方案**：检查 `.env` 文件中是否正确配置了 `GITHUB_SYNC_TOKEN` 和 `GITHUB_SYNC_REPO`

### 问题：同步失败，提示权限错误

**解决方案**：
1. 确认 PAT 具有 `repo` 权限
2. 确认仓库 URL 正确
3. 确认 PAT 未过期

### 问题：数据未同步到 GitHub

**解决方案**：
1. 检查系统日志，查看是否有错误信息
2. 确认网络连接正常
3. 确认 GitHub 服务可访问

## API 接口

系统还提供了手动触发同步的能力（需要在代码中实现相应的 API 接口）：

- 手动同步到 GitHub：调用 `service.GetGitHubSyncService().SyncToGitHub()`
- 手动加载数据：调用 `service.GetGitHubSyncService().LoadFromGitHub()`

## 注意事项

1. **首次启动**：首次启动时会从 GitHub 加载数据，如果仓库为空，将使用数据库中的数据
2. **数据冲突**：系统会以最新的同步数据为准
3. **大量数据**：如果数据量很大，同步可能需要一些时间
4. **网络依赖**：需要稳定的网络连接到 GitHub

## 更新记录

- 2025-01-16: 初始版本，支持基本的 GitHub 数据同步功能