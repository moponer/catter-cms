-- ========================================
-- 主题系统数据库迁移 - 回滚
-- ========================================

-- 删除主题备份表
DROP TABLE IF EXISTS theme_backups;

-- 删除主题配置历史表
DROP TABLE IF EXISTS theme_config_history;

-- 删除主题自定义代码表
DROP TABLE IF EXISTS theme_custom_code;

-- 删除主题SEO配置表
DROP TABLE IF EXISTS theme_seo;

-- 删除主题页面配置表
DROP TABLE IF EXISTS theme_pages;

-- 删除主题菜单配置表
DROP TABLE IF EXISTS theme_menus;

-- 删除主题小工具配置表
DROP TABLE IF EXISTS theme_widgets;

-- 删除主题模块配置表
DROP TABLE IF EXISTS theme_modules;

-- 删除themes表(谨慎操作,如果其他功能依赖此表,请注释掉)
-- DROP TABLE IF EXISTS themes;

-- 如果保留了themes表,删除新增的字段
-- ALTER TABLE themes DROP COLUMN IF EXISTS api_version;
-- ALTER TABLE themes DROP COLUMN IF EXISTS capabilities;
-- ALTER TABLE themes DROP COLUMN IF EXISTS settings;
