-- 创建评论表
CREATE TABLE comments (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    user_id UUID,
    parent_id UUID DEFAULT NULL,
    author_name VARCHAR(255),
    author_email VARCHAR(255),
    author_url VARCHAR(500),
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, spam, trash
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_status ON comments(status);
CREATE INDEX idx_comments_created_at ON comments(created_at);

-- 添加外键约束
-- 注意: SQLite不支持ALTER TABLE ADD CONSTRAINT,这些约束需要在创建表时指定
-- 或者通过应用层保证数据完整性
