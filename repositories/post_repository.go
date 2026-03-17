package repositories

import (
	"catter-cms/errors"
	"catter-cms/models"
)

// PostRepository 文章仓储
type PostRepository struct {
	db *pop.Connection
}

// NewPostRepository 创建新的文章仓储
func NewPostRepository(db *pop.Connection) *PostRepository {
	return &PostRepository{db: db}
}

// GetPublishedPosts 获取已发布的文章
func (r *PostRepository) GetPublishedPosts(page, perPage int) ([]models.Post, int, error) {
	posts := models.Posts{}
	query := r.db.Where("status = ?", "published")
	query = query.Order("published_at DESC")
	
	offset := (page - 1) * perPage
	if offset < 0 {
		offset = 0
	}
	
	query = query.Limit(perPage)
	query = query.Offset(offset)
	
	if err := query.All(&posts); err != nil {
		return nil, 0, errors.NewInternalServerError("获取文章列表失败")
	}
	
	total, err := r.db.Where("status = ?", "published").Count()
	if err != nil {
		return nil, 0, errors.NewInternalServerError("获取文章总数失败")
	}
	
	return posts, total, nil
}

// GetByID 根据ID获取文章
func (r *PostRepository) GetByID(id int) (*models.Post, error) {
	post := &models.Post{}
	if err := r.db.Find(post, id); err != nil {
		return nil, errors.NewNotFoundError("文章", id)
	}
	return post, nil
}

// Create 创建新文章
func (r *PostRepository) Create(post *models.Post) error {
	if err := r.db.Create(post); err != nil {
		return errors.NewInternalServerError("创建文章失败")
	}
	return nil
}

// Update 更新文章
func (r *PostRepository) Update(post *models.Post) error {
	if err := r.db.Update(post); err != nil {
		return errors.NewInternalServerError("更新文章失败")
	}
	return nil
}

// Delete 删除文章
func (r *PostRepository) Delete(id int) error {
	post := &models.Post{}
	if err := r.db.Find(post, id); err != nil {
		return errors.NewNotFoundError("文章", id)
	}
	
	if err := r.db.Destroy(post); err != nil {
		return errors.NewInternalServerError("删除文章失败")
	}
	
	return nil
}

// Count 获取文章总数
func (r *PostRepository) Count() (int, error) {
	count, err := r.db.Where("status = ?", "published").Count()
	if err != nil {
		return 0, errors.NewInternalServerError("获取文章总数失败")
	}
	return count, nil
}
