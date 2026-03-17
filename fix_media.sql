-- 重建 media 表
DROP TABLE IF EXISTS media;

CREATE TABLE media (
    id TEXT PRIMARY KEY,
    filename TEXT NOT NULL,
    path TEXT NOT NULL,
    size INTEGER,
    mime_type TEXT,
    thumbnail TEXT,
    description TEXT,
    user_id TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
