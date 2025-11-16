-- 为 logs 表添加请求和响应内容字段
-- 用于记录 AI 对话内容，便于管理员查看和分析

-- MySQL/MariaDB 语法
ALTER TABLE `logs` 
ADD COLUMN `request_content` TEXT NULL COMMENT '请求内容' AFTER `content`,
ADD COLUMN `response_content` TEXT NULL COMMENT '响应内容' AFTER `request_content`;

-- 添加索引以提升查询性能（可选）
-- 注意：TEXT 字段不能直接创建索引，如需要可以使用前缀索引
-- CREATE INDEX idx_logs_created_at ON `logs`(`created_at`);

-- 如果使用 SQLite（开发环境）
-- SQLite 不支持 ADD COLUMN 后的 AFTER 子句，需要使用以下语法：
-- ALTER TABLE `logs` ADD COLUMN `request_content` TEXT;
-- ALTER TABLE `logs` ADD COLUMN `response_content` TEXT;

-- PostgreSQL 语法（如果使用 PostgreSQL）
-- ALTER TABLE "logs" 
-- ADD COLUMN "request_content" TEXT,
-- ADD COLUMN "response_content" TEXT;

-- 验证迁移
-- SELECT column_name, data_type, character_maximum_length 
-- FROM information_schema.columns 
-- WHERE table_name = 'logs' 
-- AND column_name IN ('request_content', 'response_content');