-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'editor',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建设置表
CREATE TABLE IF NOT EXISTS settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT NOT NULL UNIQUE,
    value TEXT NOT NULL,
    type TEXT DEFAULT 'string',
    category TEXT DEFAULT 'general',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入默认管理员用户（密码: admin123）
INSERT OR IGNORE INTO users (email, name, password_hash, role) VALUES 
('admin@catter-cms.com', '管理员', '$2a$10$fhpkmt7tnz1vSQpM.lxMauRvBEkFz3I5v5pRoIo95eAW8wRTq/L0i', 'admin');

-- 插入默认系统设置
INSERT OR IGNORE INTO settings (key, value, type, category) VALUES 
('site_name', 'Catter CMS', 'string', 'basic'),
('site_description', '现代化的内容管理系统', 'string', 'basic'),
('site_keywords', 'CMS, 内容管理, Buffalo, Vue', 'string', 'basic'),
('active_theme', 'catter-blog', 'string', 'theme'),
('default_locale', 'zh-CN', 'string', 'system'),
('per_page', '10', 'string', 'system'),
('allow_register', 'false', 'string', 'system'),
('require_moderation', 'false', 'string', 'system');

-- 用户表索引
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- 设置表索引
CREATE INDEX IF NOT EXISTS idx_settings_key ON settings(key);
CREATE INDEX IF NOT EXISTS idx_settings_category ON settings(category);
