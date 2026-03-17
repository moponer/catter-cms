package services

import (
	"strings"

	"catter-cms/errors"
	"catter-cms/models"
	"github.com/gobuffalo/pop/v6"
)

// PostService 文章服务
type PostService struct {
	db *pop.Connection
}

// NewPostService 创建新的文章服务
func NewPostService(db *pop.Connection) *PostService {
	return &PostService{db: db}
}

// GetPublishedPosts 获取已发布的文章
// 参数:
//   - page: 页码，从1开始
//   - perPage: 每页数量，默认10
// 返回:
//   - []models.Post: 文章列表
//   - int: 总数
//   - error: 错误信息
func (s *PostService) GetPublishedPosts(page, perPage int) ([]models.Post, int, error) {
	posts := models.Posts{}
	query := s.db.Where("status = ?", "published")
	query = query.Order("published_at DESC")
	
	// 计算偏移量
	offset := (page - 1) * perPage
	if offset < 0 {
		offset = 0
	}
	
	// 执行分页查询
	query = query.Limit(perPage)
	query = query.Paginate(page, perPage)
	
	if err := query.All(&posts); err != nil {
		return nil, 0, errors.NewInternalServerError("获取文章列表失败")
	}
	
	// 获取总数
	total, err := s.db.Where("status = ?", "published").Count(&models.Post{})
	if err != nil {
		return nil, 0, errors.NewInternalServerError("获取文章总数失败")
	}
	
	return posts, total, nil
}

// GetPostByID 根据ID获取文章
// 参数:
//   - id: 文章ID
// 返回:
//   - *models.Post: 文章对象
//   - error: 错误信息
func (s *PostService) GetPostByID(id int) (*models.Post, error) {
	post := &models.Post{}
	if err := s.db.Find(post, id); err != nil {
		return nil, errors.NewNotFoundError("文章", id)
	}
	return post, nil
}

// CreatePost 创建新文章
// 参数:
//   - post: 文章对象
// 返回:
//   - *models.Post: 创建的文章
//   - error: 错误信息
func (s *PostService) CreatePost(post *models.Post) (*models.Post, error) {
	// 验证文章数据
	if post.Title == "" {
		return nil, errors.NewBadRequestError("标题不能为空", nil)
	}
	if post.Content == "" {
		return nil, errors.NewBadRequestError("内容不能为空", nil)
	}
	
	// 设置默认状态
	if post.Status == "" {
		post.Status = "draft"
	}
	
	// 生成slug
	if post.Slug == "" {
		post.Slug = generateSlug(post.Title)
	}
	
	// 保存到数据库
	if err := s.db.Create(post); err != nil {
		return nil, errors.NewInternalServerError("创建文章失败")
	}
	
	return post, nil
}

// UpdatePost 更新文章
// 参数:
//   - id: 文章ID
//   - post: 更新的文章对象
// 返回:
//   - *models.Post: 更新后的文章
//   - error: 错误信息
func (s *PostService) UpdatePost(id int, post *models.Post) (*models.Post, error) {
	// 查找现有文章
	existingPost := &models.Post{}
	if err := s.db.Find(existingPost, id); err != nil {
		return nil, errors.NewNotFoundError("文章", id)
	}
	
	// 验证文章数据
	if post.Title == "" {
		return nil, errors.NewBadRequestError("标题不能为空", nil)
	}
	
	// 更新字段
	existingPost.Title = post.Title
	existingPost.Content = post.Content
	existingPost.Excerpt = post.Excerpt
	existingPost.FeaturedImage = post.FeaturedImage
	existingPost.CategoryID = post.CategoryID
	existingPost.Status = post.Status
	
	// 如果提供了slug，更新slug
	if post.Slug != "" {
		existingPost.Slug = post.Slug
	}
	
	// 更新到数据库
	if err := s.db.Update(existingPost); err != nil {
		return nil, errors.NewInternalServerError("更新文章失败")
	}
	
	return existingPost, nil
}

// DeletePost 删除文章
// 参数:
//   - id: 文章ID
// 返回:
//   - error: 错误信息
func (s *PostService) DeletePost(id int) error {
	post := &models.Post{}
	if err := s.db.Find(post, id); err != nil {
		return errors.NewNotFoundError("文章", id)
	}
	
	if err := s.db.Destroy(post); err != nil {
		return errors.NewInternalServerError("删除文章失败")
	}
	
	return nil
}

// generateSlug 生成slug
func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "[^a-z0-9-]", "")
	return slug
}
