-- 插入3个内置主题记录

INSERT INTO themes (id, name, slug, version, description, author, thumbnail, is_active, created_at, updated_at)
VALUES 
  ('00000000-0000-0000-0000-000000000001', 'Catter CMS Official', 'catter-cms-official', '1.0.0', 'Catter CMS官方官网主题，现代化设计风格', 'Catter Team', '/themes/catter-cms-official/assets/images/thumbnail.jpg', 0, datetime('now'), datetime('now')),
  ('00000000-0000-0000-0000-000000000002', 'Enterprise Website', 'enterprise-website', '1.0.0', '通用企业官网主题，适合各类企业展示使用', 'Catter Team', '/themes/enterprise-website/assets/images/thumbnail.jpg', 0, datetime('now'), datetime('now')),
  ('00000000-0000-0000-0000-000000000003', 'Catter Blog', 'catter-blog', '1.0.0', '现代化的博客主题，支持响应式设计', 'Catter Team', '/themes/catter-blog/assets/images/thumbnail.jpg', 1, datetime('now'), datetime('now'));

INSERT OR REPLACE INTO settings (key, value, created_at, updated_at)
VALUES ('active_theme', 'catter-blog', datetime('now'), datetime('now'));
