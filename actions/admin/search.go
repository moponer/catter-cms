package admin

import (
	"net/http"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"catter-cms/models"
	"catter-cms/services"
)

// GlobalSearchAPI 全局搜索API
// GET /api/search?q=keyword&limit=20&page=1&type=post
func GlobalSearchAPI(c buffalo.Context) error {
	query := c.Param("q")
	if query == "" {
		query = c.Request().URL.Query().Get("q")
	}
	
	// 获取搜索参数
	searchType := c.Request().URL.Query().Get("type") // post, page, comment, all
	limitStr := c.Request().URL.Query().Get("limit")
	pageStr := c.Request().URL.Query().Get("page")
	
	// 设置默认值
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	offset := (page - 1) * limit
	
	if query == "" {
		return c.Render(http.StatusOK, render.JSON(map[string]interface{}{
			"success": false,
			"message": "搜索关键词不能为空",
			"data": map[string]interface{}{
				"results": []interface{}{},
				"total":   0,
				"page":    page,
				"limit":   limit,
			},
		}))
	}
	
	// 获取搜索服务
	searchService := services.GetSearchService()

	// 执行搜索
	results, total, err := searchService.Search(query, limit, offset)
	if err != nil {
		return c.Render(http.StatusInternalServerError, render.JSON(map[string]interface{}{
			"success": false,
			"message": "搜索失败: " + err.Error(),
		}))
	}
	
	// 如果指定了类型,过滤结果
	if searchType != "" && searchType != "all" {
		filtered := make([]services.SearchResult, 0)
		for _, result := range results {
			if result.Type == searchType {
				filtered = append(filtered, result)
			}
		}
		results = filtered
	}
	
	// 获取详细数据
	detailedResults := make([]map[string]interface{}, 0)
	for _, result := range results {
		detail := map[string]interface{}{
			"id":    result.ID,
			"title": result.Title,
			"type":  result.Type,
			"score": result.Score,
		}
		
		// 根据类型获取详细数据
		switch result.Type {
		case "post":
			var post models.Post
			if err := models.DB.Find(&post, result.ID); err == nil && post.ID != 0 {
				detail["slug"] = post.Slug
				detail["excerpt"] = post.Excerpt
				detail["status"] = post.Status
				detail["created_at"] = post.CreatedAt
				detail["updated_at"] = post.UpdatedAt
				if post.CategoryID != nil {
					var category models.Category
					if err := models.DB.Find(&category, *post.CategoryID); err == nil {
						detail["category"] = category.Name
					}
				}
				// 获取标签
				tags := []models.Tag{}
				models.DB.Q().Where("post_id = ?", post.ID).All(&tags)
				tagNames := make([]string, len(tags))
				for i, tag := range tags {
					tagNames[i] = tag.Name
				}
				detail["tags"] = tagNames
			}
		case "page":
			var page models.Page
			if err := models.DB.Find(&page, result.ID); err == nil && page.ID != 0 {
				detail["slug"] = page.Slug
				detail["excerpt"] = page.Excerpt
				detail["status"] = page.Status
				detail["created_at"] = page.CreatedAt
				detail["updated_at"] = page.UpdatedAt
			}
		case "comment":
			var comment models.Comment
			if err := models.DB.Find(&comment, result.ID); err == nil && comment.ID != "" {
				detail["author_name"] = comment.AuthorName
				detail["content"] = comment.Content
				detail["status"] = comment.Status
				detail["created_at"] = comment.CreatedAt
				
				// 获取关联的文章
				if comment.PostID != "" {
					var post models.Post
					if err := models.DB.Find(&post, comment.PostID); err == nil {
						detail["post_title"] = post.Title
						detail["post_slug"] = post.Slug
					}
				}
			}
		}
		
		detailedResults = append(detailedResults, detail)
	}
	
	return c.Render(http.StatusOK, render.JSON(map[string]interface{}{
		"success": true,
		"message": "搜索成功",
		"data": map[string]interface{}{
			"query":   query,
			"results": detailedResults,
			"total":   total,
			"page":    page,
			"limit":   limit,
			"type":    searchType,
		},
	}))
}

// RebuildSearchIndexAPI 重建搜索索引API
// POST /api/admin/search/rebuild
func RebuildSearchIndexAPI(c buffalo.Context) error {
	// 获取搜索服务
	searchService := services.GetSearchService()

	// 获取所有已发布的文章
	posts := []models.Post{}
	if err := models.DB.Q().Where("status = ?", "published").All(&posts); err != nil {
		return c.Render(http.StatusInternalServerError, render.JSON(map[string]interface{}{
			"success": false,
			"message": "获取文章失败: " + err.Error(),
		}))
	}
	
	// 索引文章
	for _, post := range posts {
		content := post.Content
		if post.Excerpt != "" {
			content = post.Excerpt + "\n" + content
		}
		if err := searchService.IndexPost(
			strconv.Itoa(post.ID),
			post.Title,
			content,
			"",
		); err != nil {
			return c.Render(http.StatusInternalServerError, render.JSON(map[string]interface{}{
				"success": false,
				"message": "索引文章失败: " + err.Error(),
			}))
		}
	}

	// 获取所有已发布的页面
	pages := []models.Page{}
	if err := models.DB.Q().Where("status = ?", "published").All(&pages); err != nil {
		return c.Render(http.StatusInternalServerError, render.JSON(map[string]interface{}{
			"success": false,
			"message": "获取页面失败: " + err.Error(),
		}))
	}

	// 索引页面
	for _, page := range pages {
		content := page.Content
		if page.Excerpt != "" {
			content = page.Excerpt + "\n" + content
		}
		if err := searchService.IndexPage(
			strconv.Itoa(page.ID),
			page.Title,
			content,
		); err != nil {
			return c.Render(http.StatusInternalServerError, render.JSON(map[string]interface{}{
				"success": false,
				"message": "索引页面失败: " + err.Error(),
			}))
		}
	}

	// 获取所有已审核的评论
	comments := []models.Comment{}
	if err := models.DB.Q().Where("status = ?", "approved").All(&comments); err != nil {
		return c.Render(http.StatusInternalServerError, render.JSON(map[string]interface{}{
			"success": false,
			"message": "获取评论失败: " + err.Error(),
		}))
	}

	// 索引评论
	for _, comment := range comments {
		if err := searchService.IndexComment(
			comment.ID,
			comment.Content,
			comment.AuthorName,
		); err != nil {
			return c.Render(http.StatusInternalServerError, render.JSON(map[string]interface{}{
				"success": false,
				"message": "索引评论失败: " + err.Error(),
			}))
		}
	}

	return c.Render(http.StatusOK, render.JSON(map[string]interface{}{
		"success": true,
		"message": "索引重建成功",
		"data": map[string]interface{}{
			"posts_count":    len(posts),
			"pages_count":    len(pages),
			"comments_count": len(comments),
		},
	}))
}
