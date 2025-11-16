# AI 对话内容记录功能文档

## 功能概述

本功能允许管理员记录和查看通过 API 中转的 AI 对话内容，包括用户请求和 AI 响应。这对于审计、分析用户使用情况、改进服务质量等场景非常有用。

## ⚠️ 重要提示

### 隐私和安全注意事项

1. **用户隐私**：此功能会记录用户的完整对话内容，请确保符合相关隐私法规（如 GDPR、CCPA 等）
2. **数据安全**：对话内容可能包含敏感信息，请确保数据库安全
3. **存储空间**：对话内容会占用大量存储空间，建议配置合理的保留期限
4. **性能影响**：记录内容会略微增加响应时间和数据库负载
5. **仅管理员可见**：普通用户无法查看任何对话内容，仅管理员（Root 用户）可访问

### 默认配置

- **功能状态**：默认关闭（`CONTENT_LOGGING_ENABLED=false`）
- **保留期限**：7天（超过7天的内容会自动清理）
- **内容长度限制**：10000字符（超过部分会被截断）

## 功能特性

### 1. 内容记录

- ✅ 记录用户的完整请求内容（JSON 格式）
- ✅ 记录 AI 的完整响应内容
- ✅ 支持流式和非流式响应
- ✅ 自动处理超长内容（截断并标记）
- ✅ 可通过配置开关启用/禁用

### 2. 自动清理

- ✅ 每天凌晨2点自动执行清理任务
- ✅ 清理超过保留期限的对话内容
- ✅ 只清理内容字段，保留日志记录本身
- ✅ 清理过程不影响系统运行

### 3. 管理员查询

- ✅ 管理员可查看所有用户的对话内容
- ✅ 支持分页、筛选、搜索
- ✅ 支持按用户、模型、时间范围筛选
- ✅ 可查看单条日志的完整详情

## 配置说明

### 环境变量

在 `.env` 文件中添加以下配置：

```bash
# 是否启用内容记录功能（默认关闭）
CONTENT_LOGGING_ENABLED=false

# 内容保留天数（默认7天）
CONTENT_RETENTION_DAYS=7

# 单次记录的最大内容长度（默认10000字符，超过会被截断）
MAX_CONTENT_LENGTH=10000
```

### 数据库配置

也可以在系统设置中动态修改配置（无需重启）：

1. 登录管理员账号
2. 进入"系统设置"页面
3. 找到对话内容记录相关选项
4. 修改配置并保存

配置项：
- `ContentLoggingEnabled`：是否启用内容记录
- `ContentRetentionDays`：内容保留天数
- `MaxContentLength`：最大内容长度

## 数据库迁移

### 首次启用功能

如果您是首次启用此功能，需要执行数据库迁移：

#### 方法1：使用 SQL 脚本

```bash
# MySQL/MariaDB
mysql -u your_username -p your_database < scripts/migrate_add_content_fields.sql

# Docker 环境
docker exec -i your_container_name mysql -u root -p your_database < scripts/migrate_add_content_fields.sql
```

#### 方法2：GORM 自动迁移

如果您的应用使用了 GORM 的 AutoMigrate，只需重启应用即可自动创建字段。

详细迁移说明请参考：[`scripts/README_MIGRATION.md`](../scripts/README_MIGRATION.md)

## API 接口

### 1. 获取对话内容列表

**请求**：
```http
GET /api/log/content?page=1&page_size=20
```

**权限**：仅管理员

**查询参数**：
- `page`：页码（默认1）
- `page_size`：每页数量（默认10）
- `type`：日志类型（默认2，消费日志）
- `start_timestamp`：开始时间戳
- `end_timestamp`：结束时间戳
- `username`：用户名筛选
- `token_name`：令牌名称筛选
- `model_name`：模型名称筛选
- `channel`：渠道ID筛选
- `group`：用户组筛选

**响应示例**：
```json
{
  "success": true,
  "message": "",
  "data": {
    "items": [
      {
        "id": 12345,
        "user_id": 1,
        "username": "user@example.com",
        "created_at": 1700000000,
        "model_name": "gpt-4",
        "prompt_tokens": 100,
        "completion_tokens": 200,
        "request_content": "{\"messages\":[...]}",
        "response_content": "AI的回复内容..."
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 2. 获取单条日志详情

**请求**：
```http
GET /api/log/content/:id
```

**权限**：仅管理员

**响应示例**：
```json
{
  "success": true,
  "message": "",
  "data": {
    "id": 12345,
    "user_id": 1,
    "username": "user@example.com",
    "created_at": 1700000000,
    "type": 2,
    "content": "消费描述",
    "model_name": "gpt-4",
    "prompt_tokens": 100,
    "completion_tokens": 200,
    "quota": 300,
    "request_content": "{\"messages\":[{\"role\":\"user\",\"content\":\"你好\"}]}",
    "response_content": "你好！有什么我可以帮助你的吗？",
    "token_name": "my-token",
    "channel_id": 1,
    "is_stream": true
  }
}
```

## 使用场景

### 1. 审计和合规

- 记录所有 API 调用的完整内容
- 用于合规审计和问题追溯
- 分析用户使用模式

### 2. 质量监控

- 检查用户提问质量
- 分析 AI 响应质量
- 发现和改进系统问题

### 3. 滥用检测

- 检测异常使用行为
- 发现违规内容
- 及时采取措施

### 4. 数据分析

- 分析用户需求趋势
- 优化模型选择
- 改进服务策略

## 最佳实践

### 1. 配置建议

```bash
# 生产环境推荐配置
CONTENT_LOGGING_ENABLED=true
CONTENT_RETENTION_DAYS=7          # 平衡存储和审计需求
MAX_CONTENT_LENGTH=5000           # 避免存储超大内容
```

### 2. 存储优化

- 定期监控数据库大小
- 根据实际需求调整保留天数
- 考虑使用数据库压缩功能
- 对于超大量日志，可考虑分表存储

### 3. 性能优化

- 为 `created_at` 字段建立索引（已默认创建）
- 避免频繁查询超大时间范围
- 使用分页查询，避免一次性加载大量数据

### 4. 安全建议

- 确保只有受信任的管理员可以访问
- 定期审查访问日志
- 考虑对敏感内容进行加密存储
- 遵守相关隐私法规

### 5. 隐私保护

- 在用户协议中明确说明数据记录策略
- 提供用户数据删除机制
- 考虑匿名化处理
- 确保数据传输加密（HTTPS）

## 故障排查

### 1. 内容未被记录

检查以下项：
- [ ] `CONTENT_LOGGING_ENABLED` 是否为 `true`
- [ ] 数据库迁移是否已执行
- [ ] `request_content` 和 `response_content` 字段是否存在
- [ ] 应用是否已重启（如果修改了环境变量）

### 2. 内容被截断

- 检查 `MAX_CONTENT_LENGTH` 配置
- 如需记录更长内容，增大此值
- 注意：过长的内容会影响数据库性能

### 3. 清理任务未执行

- 检查应用日志，查看清理任务是否启动
- 确认 `CONTENT_RETENTION_DAYS` 配置正确
- 验证定时任务是否正常运行

### 4. 查询接口返回403

- 确认使用管理员账号登录
- 检查用户角色是否为 Root User
- 查看应用日志了解详细错误信息

## 相关文件

### 后端代码

- [`model/log.go`](../model/log.go) - Log 模型定义和清理函数
- [`controller/log.go`](../controller/log.go) - 对话内容查询接口
- [`router/api-router.go`](../router/api-router.go) - API 路由配置
- [`relay/compatible_handler.go`](../relay/compatible_handler.go) - 请求内容捕获
- [`relay/channel/openai/relay-openai.go`](../relay/channel/openai/relay-openai.go) - 响应内容捕获
- [`main.go`](../main.go) - 定时清理任务启动

### 配置文件

- [`common/constants.go`](../common/constants.go) - 配置常量定义
- [`model/option.go`](../model/option.go) - 配置管理
- [`.env.example`](../.env.example) - 环境变量示例

### 数据库脚本

- [`scripts/migrate_add_content_fields.sql`](../scripts/migrate_add_content_fields.sql) - 数据库迁移脚本
- [`scripts/README_MIGRATION.md`](../scripts/README_MIGRATION.md) - 迁移说明文档

## 技术实现

### 数据流程

```
用户请求 → Gin Context → 请求内容捕获 → 转发到上游
                               ↓
                        存储到 context
                               ↓
上游响应 → 响应内容捕获 → 存储到 context
                               ↓
                        RecordConsumeLog
                               ↓
                        保存到数据库
```

### 清理流程

```
应用启动 → 启动定时任务 → 等待到凌晨2点
                               ↓
                        执行清理函数
                               ↓
                    更新超期记录内容为NULL
                               ↓
                        记录清理日志
                               ↓
                        等待下次执行
```

## 常见问题

### Q1: 启用此功能会影响性能吗？

A: 会有轻微影响，主要体现在：
- 数据库写入增加（每次请求多保存两个TEXT字段）
- 略微增加响应时间（毫秒级）
- 建议在生产环境测试后再全面启用

### Q2: 可以只记录某些用户的对话吗？

A: 当前版本不支持按用户筛选记录，是全局开关。如有需求，可以修改 `RecordConsumeLog` 函数添加用户筛选逻辑。

### Q3: 清理后的数据可以恢复吗？

A: 不可以。清理操作会将内容字段设为 NULL，数据无法恢复。请在配置保留期限时谨慎考虑。

### Q4: 可以导出对话内容吗？

A: 可以通过 API 接口获取数据后自行导出，或直接从数据库导出。未来版本可能会添加导出功能。

### Q5: 支持其他 AI 提供商吗？

A: 当前主要支持 OpenAI 格式的请求/响应。其他格式需要在对应的 channel handler 中添加内容捕获逻辑。

## 更新日志

### v1.0.0 (2025-01-16)

- ✅ 初始版本发布
- ✅ 支持请求和响应内容记录
- ✅ 自动清理过期内容
- ✅ 管理员查询接口
- ✅ 配置管理功能

## 支持

如有问题或建议，请：
1. 查看本文档的故障排查部分
2. 检查应用日志
3. 提交 Issue 到项目仓库

## 许可证

本功能遵循项目的开源许可证。