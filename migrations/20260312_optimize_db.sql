-- 添加索引以优化查询性能

-- 文章表索引
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_category_id ON posts(category_id);
CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);

-- 分类表索引
CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug);
CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id);

-- 标签表索引
CREATE INDEX IF NOT EXISTS idx_tags_slug ON tags(slug);

-- 用户表索引
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- 文章标签关联表索引
CREATE INDEX IF NOT EXISTS idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX IF NOT EXISTS idx_post_tags_tag_id ON post_tags(tag_id);

-- 媒体表索引
CREATE INDEX IF NOT EXISTS idx_media_post_id ON media(post_id);
CREATE INDEX IF NOT EXISTS idx_media_type ON media(type);

-- 主题表索引
CREATE INDEX IF NOT EXISTS idx_themes_name ON themes(name);
CREATE INDEX IF NOT EXISTS idx_themes_active ON themes(active);

-- 设置表索引
CREATE INDEX IF NOT EXISTS idx_settings_key ON settings(key);
