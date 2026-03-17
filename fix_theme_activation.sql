-- 修复主题激活状态
-- 将所有主题设为未激活
UPDATE themes SET is_active = 0;

-- 激活 blog 主题（唯一有完整模板的主题）
UPDATE themes SET is_active = 1 WHERE slug = 'blog';

-- 更新系统设置
INSERT OR REPLACE INTO settings (key, value, category, type, created_at, updated_at)
VALUES ('active_theme', 'blog', 'theme', 'string', datetime('now'), datetime('now'));
