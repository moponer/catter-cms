-- 创建定时任务表
CREATE TABLE IF NOT EXISTS scheduled_tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    task_type VARCHAR(50) NOT NULL,
    cron_expr VARCHAR(100) NOT NULL,
    command TEXT,
    enabled BOOLEAN DEFAULT 1,
    last_run_at DATETIME,
    next_run_at DATETIME,
    run_count INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'idle',
    error_msg TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_scheduled_tasks_enabled ON scheduled_tasks(enabled);
CREATE INDEX IF NOT EXISTS idx_scheduled_tasks_task_type ON scheduled_tasks(task_type);
CREATE INDEX IF NOT EXISTS idx_scheduled_tasks_status ON scheduled_tasks(status);
CREATE INDEX IF NOT EXISTS idx_scheduled_tasks_next_run_at ON scheduled_tasks(next_run_at);

-- 插入默认定时任务（如果表为空）
INSERT OR IGNORE INTO scheduled_tasks (id, name, task_type, cron_expr, command, enabled, status) VALUES
(1, '定时发布文章', 'publish', '0 * * * *', '', 1, 'idle'),
(2, '每日数据备份', 'backup', '0 2 * * *', '', 1, 'idle'),
(3, '缓存清理', 'cache_clear', '0 */6 * * *', '', 1, 'idle');
