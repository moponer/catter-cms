-- 创建页面表
CREATE TABLE IF NOT EXISTS pages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    content TEXT NOT NULL,
    excerpt TEXT,
    featured_image TEXT,
    template TEXT DEFAULT 'default',
    status TEXT NOT NULL DEFAULT 'draft',
    user_id INTEGER NOT NULL,
    parent_id INTEGER,
    "order" INTEGER NOT NULL DEFAULT 0,
    meta_title TEXT,
    meta_description TEXT,
    published_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES pages(id) ON DELETE CASCADE
);

-- 创建索引
CREATE INDEX idx_pages_status ON pages(status);
CREATE INDEX idx_pages_slug ON pages(slug);
CREATE INDEX idx_pages_parent_id ON pages(parent_id);
CREATE INDEX idx_pages_order ON pages("order");
CREATE INDEX idx_pages_published_at ON pages(published_at);

-- 插入示例页面（关于我们、联系我们）
INSERT INTO pages (title, slug, content, status, user_id, template, "order") VALUES 
('关于我们', 'about', '# 关于我们\n\n这是关于我们页面。', 'published', 1, 'default', 1),
('联系我们', 'contact', '# 联系我们\n\n欢迎通过以下方式联系我们：\n\n- 邮箱：contact@example.com\n- 电话：123-456-7890', 'published', 1, 'default', 2);
