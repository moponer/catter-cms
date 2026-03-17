package container

import (
	"catter-cms/repositories"
	"catter-cms/services"
)

// Container 依赖注入容器
type Container struct {
	db       *pop.Connection
	postRepo repositories.PostRepository
	postService services.PostService
}

// NewContainer 创建新的依赖注入容器
func NewContainer(db *pop.Connection) *Container {
	return &Container{
		db:          db,
		postRepo:    repositories.NewPostRepository(db),
		postService: services.NewPostService(db),
	}
}

// GetPostService 获取文章服务
func (c *Container) GetPostService() *services.PostService {
	return c.postService
}

// GetPostRepository 获取文章仓储
func (c *Container) GetPostRepository() *repositories.PostRepository {
	return c.postRepo
}
