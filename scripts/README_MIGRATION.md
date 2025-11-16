# 数据库迁移说明

## 对话内容记录功能迁移

### 概述
此迁移为 `logs` 表添加两个新字段，用于记录 AI 对话的请求和响应内容。

### 迁移文件
- `migrate_add_content_fields.sql` - 添加请求和响应内容字段

### 执行步骤

#### 方法1：使用 MySQL 命令行

```bash
# 登录到 MySQL
mysql -u your_username -p your_database

# 执行迁移脚本
source scripts/migrate_add_content_fields.sql;

# 或者直接执行
mysql -u your_username -p your_database < scripts/migrate_add_content_fields.sql
```

#### 方法2：使用 Docker（如果使用 Docker 部署）

```bash
# 复制脚本到容器
docker cp scripts/migrate_add_content_fields.sql your_container_name:/tmp/

# 在容器内执行
docker exec -i your_container_name mysql -u root -p your_database < /tmp/migrate_add_content_fields.sql
```

#### 方法3：使用 GORM 自动迁移（推荐）

如果应用使用了 GORM 的 AutoMigrate，模型更新后重启应用即可自动创建新字段：

```go
// 在 main.go 或初始化代码中
db.AutoMigrate(&model.Log{})
```

### 新增字段说明

| 字段名 | 类型 | 说明 | 可为空 |
|--------|------|------|--------|
| `request_content` | TEXT | 用户请求的完整内容（JSON 格式） | 是 |
| `response_content` | TEXT | AI 返回的响应内容 | 是 |

### 验证迁移

执行以下 SQL 验证字段是否添加成功：

```sql
DESCRIBE logs;

-- 或者
SHOW COLUMNS FROM logs LIKE '%content%';

-- 或者
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'logs' 
AND column_name IN ('request_content', 'response_content');
```

### 回滚方案

如果需要回滚此迁移：

```sql
ALTER TABLE `logs` 
DROP COLUMN `request_content`,
DROP COLUMN `response_content`;
```

### 注意事项

1. **备份数据**：执行迁移前请务必备份数据库
2. **停机时间**：对于大表，添加字段可能需要一些时间
3. **存储空间**：TEXT 字段会占用额外存储空间，建议配置自动清理策略
4. **性能影响**：新字段为 NULL 类型，不会影响现有查询性能
5. **数据库类型**：脚本默认为 MySQL/MariaDB，其他数据库需要调整语法

### 配置说明

迁移完成后，需要在 `.env` 文件中配置以下选项：

```bash
# 是否启用内容记录（默认关闭）
CONTENT_LOGGING_ENABLED=false

# 内容保留天数（默认7天）
CONTENT_RETENTION_DAYS=7

# 单次记录的最大内容长度（默认10000字符）
MAX_CONTENT_LENGTH=10000
```

### 清理策略

系统会自动清理超过保留期限的对话内容。你也可以手动执行清理：

```sql
-- 清理7天前的对话内容
UPDATE logs 
SET request_content = NULL, response_content = NULL 
WHERE created_at < DATE_SUB(NOW(), INTERVAL 7 DAY)
AND (request_content IS NOT NULL OR response_content IS NOT NULL);
```

### 相关文件

- `model/log.go` - Log 模型定义
- `common/constants.go` - 配置常量
- `model/option.go` - 配置管理
- `relay/compatible_handler.go` - 请求内容捕获
- `relay/channel/openai/relay-openai.go` - 响应内容捕获