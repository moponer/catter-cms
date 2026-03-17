-- ========================================
-- 主题系统数据库迁移
-- 创建主题配置相关表
-- ========================================

-- 主题表(扩展原有的themes表)
-- 已在之前的迁移中创建,这里添加必要的字段
CREATE TABLE IF NOT EXISTS themes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    version VARCHAR(50) NOT NULL DEFAULT '1.0.0',
    api_version VARCHAR(20) DEFAULT '1.0',
    description TEXT,
    author VARCHAR(255),
    author_email VARCHAR(255),
    author_url VARCHAR(500),
    license VARCHAR(50) DEFAULT 'MIT',
    thumbnail VARCHAR(500),
    is_active BOOLEAN DEFAULT FALSE,
    capabilities JSON,              -- 主题能力(JSON格式)
    settings JSON,                  -- 主题配置(JSON格式)
    created_at DATETIME,
    updated_at DATETIME
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_themes_slug ON themes(slug);
CREATE INDEX IF NOT EXISTS idx_themes_is_active ON themes(is_active);

-- 主题模块配置表
CREATE TABLE IF NOT EXISTS theme_modules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    theme_id INTEGER NOT NULL,
    module_id VARCHAR(100) NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    sort_order INTEGER DEFAULT 0,   -- 排序顺序
    config JSON,                    -- 模块配置(JSON格式)
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (theme_id) REFERENCES themes(id) ON DELETE CASCADE,
    UNIQUE(theme_id, module_id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_theme_modules_theme_id ON theme_modules(theme_id);
CREATE INDEX IF NOT EXISTS idx_theme_modules_enabled ON theme_modules(enabled);
CREATE INDEX IF NOT EXISTS idx_theme_modules_sort_order ON theme_modules(sort_order);

-- 主题小工具配置表
CREATE TABLE IF NOT EXISTS theme_widgets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    theme_id INTEGER NOT NULL,
    widget_id VARCHAR(100) NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    area VARCHAR(100) NOT NULL,     -- 小工具区域(sidebar/footer等)
    sort_order INTEGER DEFAULT 0,   -- 排序顺序
    config JSON,                    -- 小工具配置(JSON格式)
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (theme_id) REFERENCES themes(id) ON DELETE CASCADE,
    UNIQUE(theme_id, widget_id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_theme_widgets_theme_id ON theme_widgets(theme_id);
CREATE INDEX IF NOT EXISTS idx_theme_widgets_area ON theme_widgets(area);
CREATE INDEX IF NOT EXISTS idx_theme_widgets_enabled ON theme_widgets(enabled);
CREATE INDEX IF NOT EXISTS idx_theme_widgets_sort_order ON theme_widgets(sort_order);

-- 主题菜单配置表
CREATE TABLE IF NOT EXISTS theme_menus (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    theme_id INTEGER NOT NULL,
    location VARCHAR(100) NOT NULL, -- 菜单位置(header/footer等)
    label VARCHAR(255) NOT NULL,    -- 菜单标签
    url VARCHAR(500) NOT NULL,      -- 菜单链接
    icon VARCHAR(100),              -- 菜单图标
    target VARCHAR(20) DEFAULT '_self', -- 打开方式
    parent_id INTEGER,              -- 父菜单ID(用于多级菜单)
    sort_order INTEGER DEFAULT 0,   -- 排序顺序
    enabled BOOLEAN DEFAULT TRUE,
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (theme_id) REFERENCES themes(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES theme_menus(id) ON DELETE CASCADE
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_theme_menus_theme_id ON theme_menus(theme_id);
CREATE INDEX IF NOT EXISTS idx_theme_menus_location ON theme_menus(location);
CREATE INDEX IF NOT EXISTS idx_theme_menus_parent_id ON theme_menus(parent_id);
CREATE INDEX IF NOT EXISTS idx_theme_menus_enabled ON theme_menus(enabled);

-- 主题页面配置表
CREATE TABLE IF NOT EXISTS theme_pages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    theme_id INTEGER NOT NULL,
    page_key VARCHAR(100) NOT NULL, -- 页面标识(home/blog等)
    title VARCHAR(255),             -- 页面标题
    template VARCHAR(100),          -- 模板名称
    layout VARCHAR(100),            -- 布局名称
    enabled BOOLEAN DEFAULT TRUE,
    config JSON,                    -- 页面配置(JSON格式)
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (theme_id) REFERENCES themes(id) ON DELETE CASCADE,
    UNIQUE(theme_id, page_key)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_theme_pages_theme_id ON theme_pages(theme_id);
CREATE INDEX IF NOT EXISTS idx_theme_pages_enabled ON theme_pages(enabled);

-- 主题SEO配置表
CREATE TABLE IF NOT EXISTS theme_seo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    theme_id INTEGER NOT NULL,
    title VARCHAR(500),             -- 网站标题
    description TEXT,               -- 网站描述
    keywords TEXT,                  -- 关键词(逗号分隔)
    author VARCHAR(255),            -- 作者
    canonical_url VARCHAR(500),     -- 规范URL
    robots VARCHAR(100),           -- 爬虫规则
    opengraph JSON,                -- Open Graph配置
    twitter JSON,                  -- Twitter Card配置
    schema JSON,                   -- 结构化数据
    sitemap JSON,                  -- 站点地图配置
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (theme_id) REFERENCES themes(id) ON DELETE CASCADE,
    UNIQUE(theme_id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_theme_seo_theme_id ON theme_seo(theme_id);

-- 主题自定义代码表
CREATE TABLE IF NOT EXISTS theme_custom_code (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    theme_id INTEGER NOT NULL,
    css TEXT,                       -- 自定义CSS
    js TEXT,                        -- 自定义JavaScript
    head TEXT,                      -- 自定义Head内容
    body_top TEXT,                  -- Body顶部内容
    body_bottom TEXT,               -- Body底部内容
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (theme_id) REFERENCES themes(id) ON DELETE CASCADE,
    UNIQUE(theme_id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_theme_custom_code_theme_id ON theme_custom_code(theme_id);

-- 主题配置历史表(用于版本控制和回滚)
CREATE TABLE IF NOT EXISTS theme_config_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    theme_id INTEGER NOT NULL,
    version VARCHAR(50) NOT NULL,   -- 配置版本
    settings JSON,                  -- 完整配置快照
    description VARCHAR(500),       -- 变更描述
    created_by INTEGER,             -- 操作用户ID
    created_at DATETIME,
    FOREIGN KEY (theme_id) REFERENCES themes(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_theme_config_history_theme_id ON theme_config_history(theme_id);
CREATE INDEX IF NOT EXISTS idx_theme_config_history_version ON theme_config_history(version);
CREATE INDEX IF NOT EXISTS idx_theme_config_history_created_at ON theme_config_history(created_at);

-- 主题备份表
CREATE TABLE IF NOT EXISTS theme_backups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    theme_id INTEGER NOT NULL,
    file_path VARCHAR(500),         -- 备份文件路径
    file_size INTEGER,              -- 文件大小(字节)
    description VARCHAR(500),       -- 备份描述
    created_by INTEGER,             -- 操作用户ID
    created_at DATETIME,
    FOREIGN KEY (theme_id) REFERENCES themes(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_theme_backups_theme_id ON theme_backups(theme_id);
CREATE INDEX IF NOT EXISTS idx_theme_backups_created_at ON theme_backups(created_at);

-- 插入默认主题数据(如果不存在)
INSERT OR IGNORE INTO themes (id, name, slug, version, description, author, thumbnail, is_active, created_at, updated_at)
VALUES (1, 'Catter Blog', 'catter-blog', '2.0.0', '现代化的博客主题，支持模块化配置', 'Catter Team', '/themes/catter-blog/assets/images/thumbnail.jpg', 1, datetime('now'), datetime('now'));

INSERT OR IGNORE INTO themes (id, name, slug, version, description, author, thumbnail, is_active, created_at, updated_at)
VALUES (2, 'Catter CMS Official', 'catter-cms-official', '2.0.0', 'Catter CMS官方官网主题', 'Catter Team', '/themes/catter-cms-official/assets/images/thumbnail.jpg', 0, datetime('now'), datetime('now'));

-- 插入默认SEO配置
INSERT OR IGNORE INTO theme_seo (theme_id, title, description, keywords, author, robots, created_at, updated_at)
VALUES (1, 'Catter Blog - 技术博客', '一个现代化的技术博客,分享编程、技术和生活', '博客,技术,编程,Web开发', 'Catter Team', 'index, follow', datetime('now'), datetime('now'));

-- 插入默认自定义代码
INSERT OR IGNORE INTO theme_custom_code (theme_id, css, js, head, body_top, body_bottom, created_at, updated_at)
VALUES (1, '', '', '<meta name="theme-color" content="#3B82F6">', '', '', datetime('now'), datetime('now'));
