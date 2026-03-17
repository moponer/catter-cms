package dto

import (
	"time"
)

// PostDTO 文章数据传输对象
type PostDTO struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Slug          string     `json:"slug"`
	Content       string     `json:"content"`
	Excerpt       string     `json:"excerpt,omitempty"`
	FeaturedImage string     `json:"featured_image,omitempty"`
	Status        string     `json:"status"`
	UserID        int        `json:"user_id"`
	CategoryID    *int       `json:"category_id,omitempty"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// PostListResponse 文章列表响应
type PostListResponse struct {
	Posts []PostDTO `json:"posts"`
	Meta  MetaDTO    `json:"meta"`
}

// MetaDTO 元数据
type MetaDTO struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	Excerpt    string `json:"excerpt,omitempty"`
	CategoryID *int    `json:"category_id,omitempty"`
	FeaturedImage string `json:"featured_image,omitempty"`
	Status     string `json:"status,omitempty"`
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title      string `json:"title,omitempty"`
	Content    string `json:"content,omitempty"`
	Excerpt    string `json:"excerpt,omitempty"`
	CategoryID *int    `json:"category_id,omitempty"`
	FeaturedImage string `json:"featured_image,omitempty"`
	Status     string `json:"status,omitempty"`
}

// PostDetailResponse 文章详情响应
type PostDetailResponse struct {
	Post    PostDTO `json:"post"`
	Category *CategoryDTO `json:"category,omitempty"`
	Tags    []TagDTO `json:"tags,omitempty"`
}

// CategoryDTO 分类数据传输对象
type CategoryDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	ParentID    *int   `json:"parent_id,omitempty"`
}

// TagDTO 标签数据传输对象
type TagDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// CategoryListResponse 分类列表响应
type CategoryListResponse struct {
	Categories []CategoryDTO `json:"categories"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
