-- 创建菜单表
CREATE TABLE IF NOT EXISTS menus (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    location TEXT NOT NULL UNIQUE,
    description TEXT,
    "order" INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建菜单项表（支持多级菜单）
CREATE TABLE IF NOT EXISTS menu_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    menu_id INTEGER NOT NULL,
    parent_id INTEGER,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    target TEXT NOT NULL DEFAULT '_self',
    classes TEXT,
    icon TEXT,
    "order" INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES menu_items(id) ON DELETE CASCADE
);

-- 创建索引
CREATE INDEX idx_menus_location ON menus(location);
CREATE INDEX idx_menus_is_active ON menus(is_active);
CREATE INDEX idx_menu_items_menu_id ON menu_items(menu_id);
CREATE INDEX idx_menu_items_parent_id ON menu_items(parent_id);
CREATE INDEX idx_menu_items_order ON menu_items("order");
CREATE INDEX idx_menu_items_is_active ON menu_items(is_active);

-- 插入默认菜单
INSERT INTO menus (name, location, description, "order") VALUES 
('主导航', 'primary', '网站主导航菜单', 1),
('顶部菜单', 'top', '顶部辅助菜单', 2),
('底部菜单', 'footer', '底部菜单', 3);

-- 插入默认菜单项（主导航）
INSERT INTO menu_items (menu_id, title, url, "order") VALUES 
(1, '首页', '/', 1),
(1, '博客', '/blog', 2),
(1, '关于我们', '/about', 3),
(1, '联系我们', '/contact', 4);
