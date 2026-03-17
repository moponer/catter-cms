package services

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/index/scorch"
	"github.com/gobuffalo/buffalo"
)

// SearchService 搜索服务
type SearchService struct {
	index bleve.Index
	mu    sync.RWMutex
	app   *buffalo.App
}

// SearchResult 搜索结果
type SearchResult struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Type  string  `json:"type"` // post, page, comment, user
	Score float64 `json:"score"`
}

var (
	searchInstance *SearchService
	searchOnce     sync.Once
)

// GetSearchService 获取搜索服务单例
func GetSearchService() *SearchService {
	searchOnce.Do(func() {
		searchInstance = &SearchService{}
	})
	return searchInstance
}

// Init 初始化搜索索引
func (s *SearchService) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 创建索引路径
	indexPath := filepath.Join(".", "indexes", "content")
	
	// 尝试打开现有索引
	index, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		// 索引不存在,创建新索引
		indexMapping := bleve.NewIndexMapping()
		
		// 文档映射
		documentMapping := bleve.NewDocumentMapping()
		
		// 标题字段映射(权重更高)
		titleFieldMapping := bleve.NewTextFieldMapping()
		titleFieldMapping.Analyzer = "standard"
		documentMapping.AddFieldMappingsAt("title", titleFieldMapping)
		
		// 内容字段映射
		contentFieldMapping := bleve.NewTextFieldMapping()
		contentFieldMapping.Analyzer = "standard"
		documentMapping.AddFieldMappingsAt("content", contentFieldMapping)
		
		// 作者字段映射
		authorFieldMapping := bleve.NewTextFieldMapping()
		documentMapping.AddFieldMappingsAt("author", authorFieldMapping)
		
		indexMapping.AddDocumentMapping("post", documentMapping)
		indexMapping.AddDocumentMapping("page", documentMapping)
		indexMapping.AddDocumentMapping("comment", documentMapping)
		
		// 使用Scorch存储引擎
		index, err = bleve.NewUsing(indexPath, indexMapping, scorch.Name, scorch.Name, nil)
		if err != nil {
			return fmt.Errorf("failed to create search index: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to open search index: %w", err)
	}
	
	s.index = index
	return nil
}

// IndexPost 索引文章
func (s *SearchService) IndexPost(id, title, content, author string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.index == nil {
		return fmt.Errorf("search index not initialized")
	}
	
	doc := map[string]interface{}{
		"id":      id,
		"title":   title,
		"content": content,
		"author":  author,
		"type":    "post",
	}
	
	return s.index.Index(id, doc)
}

// IndexPage 索引页面
func (s *SearchService) IndexPage(id, title, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.index == nil {
		return fmt.Errorf("search index not initialized")
	}
	
	doc := map[string]interface{}{
		"id":      id,
		"title":   title,
		"content": content,
		"type":    "page",
	}
	
	return s.index.Index(id, doc)
}

// IndexComment 索引评论
func (s *SearchService) IndexComment(id, content, author string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.index == nil {
		return fmt.Errorf("search index not initialized")
	}
	
	doc := map[string]interface{}{
		"id":      id,
		"content": content,
		"author":  author,
		"type":    "comment",
	}
	
	return s.index.Index(id, doc)
}

// DeleteIndex 删除索引
func (s *SearchService) DeleteIndex(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.index == nil {
		return fmt.Errorf("search index not initialized")
	}
	
	return s.index.Delete(id)
}

// Search 执行搜索
func (s *SearchService) Search(query string, limit, offset int) ([]SearchResult, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if s.index == nil {
		return nil, 0, fmt.Errorf("search index not initialized")
	}
	
	// 创建查询
	queryObj := bleve.NewMatchQuery(query)
	searchRequest := bleve.NewSearchRequest(queryObj)
	searchRequest.Size = limit
	searchRequest.From = offset
	searchRequest.Fields = []string{"id", "title", "content", "type", "author"}
	
	// 执行搜索
	searchResult, err := s.index.Search(searchRequest)
	if err != nil {
		return nil, 0, fmt.Errorf("search failed: %w", err)
	}
	
	// 转换结果
	results := make([]SearchResult, 0, len(searchResult.Hits))
	for _, hit := range searchResult.Hits {
		result := SearchResult{
			ID:    hit.ID,
			Score: hit.Score,
		}
		
		if title, ok := hit.Fields["title"].(string); ok {
			result.Title = title
		}
		if docType, ok := hit.Fields["type"].(string); ok {
			result.Type = docType
		}
		
		results = append(results, result)
	}
	
	return results, int(searchResult.Total), nil
}

// RebuildAllIndex 重建所有索引
func (s *SearchService) RebuildAllIndex() error {
	// TODO: 从数据库重新索引所有内容
	// 1. 清空现有索引
	// 2. 重新索引所有文章
	// 3. 重新索引所有页面
	// 4. 重新索引所有评论
	
	return nil
}

// Close 关闭索引
func (s *SearchService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.index != nil {
		return s.index.Close()
	}
	return nil
}
