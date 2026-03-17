-- 删除索引

DROP INDEX IF EXISTS idx_posts_slug;
DROP INDEX IF EXISTS idx_posts_status;
DROP INDEX IF EXISTS idx_posts_published_at;
DROP INDEX IF EXISTS idx_posts_user_id;
DROP INDEX IF EXISTS idx_posts_category_id;
DROP INDEX IF EXISTS idx_posts_created_at;

DROP INDEX IF EXISTS idx_categories_slug;
DROP INDEX IF EXISTS idx_categories_parent_id;

DROP INDEX IF EXISTS idx_tags_slug;

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;

DROP INDEX IF EXISTS idx_comments_post_id;
DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS idx_comments_created_at;

DROP INDEX IF EXISTS idx_post_tags_post_id;
DROP INDEX IF EXISTS idx_post_tags_tag_id;
