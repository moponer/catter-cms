-- 插入测试用户
-- 如果用户已存在则跳过

INSERT INTO users (name, email, password_hash, role, created_at, updated_at)
SELECT 'Admin User', 'admin@catter.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', datetime('now'), datetime('now')
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'admin@catter.com'
);

INSERT INTO users (name, email, password_hash, role, created_at, updated_at)
SELECT 'Editor User', 'editor@catter.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'editor', datetime('now'), datetime('now')
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'editor@catter.com'
);

-- 注意: 两个用户的密码都是 admin123
