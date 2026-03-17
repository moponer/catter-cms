-- 添加 views 字段到 posts 表
ALTER TABLE posts ADD COLUMN views INTEGER DEFAULT 0;

-- 为现有文章设置默认阅读量
UPDATE posts SET views = 0 WHERE views IS NULL;
