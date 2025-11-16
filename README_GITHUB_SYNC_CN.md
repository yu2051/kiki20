# GitHub 数据同步快速开始指南

## 🚀 快速配置（3 步完成）

### 步骤 1：准备 GitHub 仓库和 Token

您已经提供了：
- ✅ **GitHub Token**: `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`
- ✅ **仓库地址**: `https://github.com/kiki0501/keykey`

### 步骤 2：配置环境变量

在项目根目录创建 `.env` 文件（如果已存在则编辑），添加以下内容：

```env
# GitHub 数据同步配置
GITHUB_SYNC_TOKEN=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
GITHUB_SYNC_REPO=https://github.com/kiki0501/keykey
GITHUB_SYNC_INTERVAL=300
```

### 步骤 3：重启应用

重启您的应用程序，系统将自动启用 GitHub 同步功能。

## 📊 功能说明

### 自动同步的数据

系统会自动同步以下数据到 GitHub：

1. **渠道数据** (`channels.json`)
   - 渠道配置
   - API 密钥（完整保留）
   - 模型映射等

2. **用户数据** (`users.json`)
   - 用户信息
   - ⚠️ 密码已自动移除
   - ⚠️ 访问令牌已自动移除

3. **令牌数据** (`tokens.json`)
   - API 令牌
   - 额度信息
   - 过期时间等

### 不同步的数据

- ❌ **日志数据** - 继续使用 MySQL 存储

## 🔄 同步机制

- **启动加载**: 应用启动时从 GitHub 加载最新数据
- **定时同步**: 每 5 分钟（可配置）自动同步一次
- **数据格式**: JSON 格式，便于版本控制和查看

## 📝 查看同步日志

启动应用后，您会看到类似以下的日志：

```
[INFO] GitHub sync service initialized: kiki0501/keykey, interval: 5m0s
[INFO] Starting GitHub auto sync...
[INFO] Loading data from GitHub...
[INFO] Loaded 10 channels from GitHub
[INFO] Loaded 5 users from GitHub  
[INFO] Loaded 15 tokens from GitHub
[INFO] Data loaded from GitHub successfully
```

每次定时同步时：

```
[INFO] Syncing data to GitHub...
[INFO] Data synced to GitHub successfully
```

## ⚙️ 高级配置

### 修改同步间隔

在 `.env` 文件中修改 `GITHUB_SYNC_INTERVAL`（单位：秒）：

```env
# 每 10 分钟同步一次
GITHUB_SYNC_INTERVAL=600

# 每 1 分钟同步一次（不推荐，频率太高）
GITHUB_SYNC_INTERVAL=60
```

### 禁用 GitHub 同步

删除或注释掉 `.env` 文件中的相关配置：

```env
# GITHUB_SYNC_TOKEN=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
# GITHUB_SYNC_REPO=https://github.com/kiki0501/keykey
```

## 🔒 安全建议

1. ✅ **已使用私有仓库** - 您的仓库是私有的，数据安全
2. ⚠️ **Token 安全提醒**:
   - 使用完成后，建议到 GitHub 设置中撤销此 Token
   - 生成新的 Token 用于生产环境
   - 不要将 `.env` 文件提交到代码仓库

3. 📌 **Token 管理位置**:
   - GitHub 设置: https://github.com/settings/tokens
   - 找到 Token 并点击 "Revoke" 可以撤销

## 🛠️ 故障排查

### 问题 1: 看不到同步日志

**原因**: 配置未生效

**解决方案**:
```bash
# 1. 检查 .env 文件是否在项目根目录
ls -la .env

# 2. 检查配置内容
cat .env

# 3. 重启应用
```

### 问题 2: 提示权限错误

**原因**: Token 权限不足或已过期

**解决方案**:
1. 确认 Token 具有 `repo` 完整权限
2. 检查 Token 是否过期
3. 重新生成 Token

### 问题 3: 数据未出现在 GitHub 仓库

**原因**: 可能需要等待首次同步

**解决方案**:
1. 等待 5 分钟（一个同步周期）
2. 检查应用日志是否有错误
3. 手动访问 GitHub 仓库检查

## 📦 查看 GitHub 仓库中的数据

访问您的仓库：https://github.com/kiki0501/keykey

您应该能看到以下文件：
- `channels.json` - 渠道配置数据
- `users.json` - 用户数据
- `tokens.json` - 令牌数据

每次同步都会创建新的 commit，您可以查看历史记录。

## 🎯 下一步

配置完成后：

1. ✅ 重启应用
2. ✅ 查看启动日志确认同步服务启动
3. ✅ 等待 5 分钟后检查 GitHub 仓库
4. ✅ 确认数据已成功同步
5. ⚠️ **重要**: 完成测试后撤销当前 Token，生成新的用于生产

## 📞 需要帮助？

如果遇到问题：
1. 查看应用日志
2. 检查网络连接
3. 确认 GitHub Token 权限
4. 参考详细文档: `docs/GITHUB_SYNC.md`

---

**安全提醒**: 本文档中的 Token 仅用于示例，请在配置完成后立即更换为新的 Token！