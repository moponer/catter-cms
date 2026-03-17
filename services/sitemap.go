package services

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"
	"time"

	"catter-cms/models"

	"github.com/gobuffalo/pop/v6"
)

// Sitemap priority constants
const (
	PriorityHomepage = 1.0
	PriorityPost     = 0.8
	PriorityPage     = 0.6
	PriorityCategory = 0.5
)

// SitemapURL represents a URL entry in the sitemap
type SitemapURL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float64 `xml:"priority,omitempty"`
}

// SitemapXML represents the sitemap XML structure
type SitemapXML struct {
	XMLName xml.Name   `xml:"urlset"`
	Xmlns   string     `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

// SitemapConfig represents the sitemap configuration
type SitemapConfig struct {
	ChangeFrequency   string
	IncludePosts      bool
	IncludePages      bool
	IncludeCategories bool
	PostPriority      float64
	PagePriority      float64
}

// DefaultSitemapConfig returns the default sitemap configuration
func DefaultSitemapConfig() SitemapConfig {
	return SitemapConfig{
		ChangeFrequency:   "weekly",
		IncludePosts:      true,
		IncludePages:      true,
		IncludeCategories: true,
		PostPriority:      PriorityPost,
		PagePriority:      PriorityPage,
	}
}

// SitemapService handles sitemap generation
type SitemapService struct {
	db      *pop.Connection
	baseURL string
	config  SitemapConfig
}

var sitemapService *SitemapService

// GetSitemapService returns the singleton sitemap service instance
func GetSitemapService() *SitemapService {
	if sitemapService == nil {
		sitemapService = &SitemapService{
			db:      models.DB,
			baseURL: "http://localhost:3000", // Default base URL, should be configurable
			config:  DefaultSitemapConfig(),
		}
	}
	return sitemapService
}

// SetBaseURL sets the base URL for the sitemap
func (s *SitemapService) SetBaseURL(url string) error {
	if url == "" {
		return fmt.Errorf("base URL cannot be empty")
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("base URL must start with http:// or https://")
	}
	s.baseURL = strings.TrimSuffix(url, "/")
	return nil
}

// GetBaseURL returns the current base URL
func (s *SitemapService) GetBaseURL() string {
	return s.baseURL
}

// SetConfig sets the sitemap configuration
func (s *SitemapService) SetConfig(config SitemapConfig) {
	s.config = config
}

// GetConfig returns the current sitemap configuration
func (s *SitemapService) GetConfig() SitemapConfig {
	return s.config
}

// Generate generates the sitemap XML
func (s *SitemapService) Generate() (string, error) {
	var urls []SitemapURL

	// Add homepage
	urls = append(urls, SitemapURL{
		Loc:        s.baseURL + "/",
		LastMod:    time.Now().Format(time.RFC3339),
		ChangeFreq: s.config.ChangeFrequency,
		Priority:   PriorityHomepage,
	})

	// Add published posts if enabled
	if s.config.IncludePosts {
		posts, err := s.getPublishedPosts()
		if err != nil {
			return "", fmt.Errorf("failed to get posts: %w", err)
		}
		for _, post := range posts {
			urls = append(urls, SitemapURL{
				Loc:        fmt.Sprintf("%s/posts/%s", s.baseURL, post.Slug),
				LastMod:    post.UpdatedAt.Format(time.RFC3339),
				ChangeFreq: s.config.ChangeFrequency,
				Priority:   s.config.PostPriority,
			})
		}
	}

	// Add published pages if enabled
	if s.config.IncludePages {
		pages, err := s.getPublishedPages()
		if err != nil {
			return "", fmt.Errorf("failed to get pages: %w", err)
		}
		for _, page := range pages {
			urls = append(urls, SitemapURL{
				Loc:        fmt.Sprintf("%s/pages/%s", s.baseURL, page.Slug),
				LastMod:    page.UpdatedAt.Format(time.RFC3339),
				ChangeFreq: s.config.ChangeFrequency,
				Priority:   s.config.PagePriority,
			})
		}
	}

	// Add categories if enabled
	if s.config.IncludeCategories {
		categories, err := s.getCategories()
		if err != nil {
			return "", fmt.Errorf("failed to get categories: %w", err)
		}
		for _, category := range categories {
			urls = append(urls, SitemapURL{
				Loc:        fmt.Sprintf("%s/categories/%s", s.baseURL, category.Slug),
				LastMod:    category.UpdatedAt.Format(time.RFC3339),
				ChangeFreq: s.config.ChangeFrequency,
				Priority:   PriorityCategory,
			})
		}
	}

	// Create sitemap structure
	sitemap := SitemapXML{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	// Generate XML
	output, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal XML: %w", err)
	}

	// Add XML header
	xmlHeader := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	return xmlHeader + string(output), nil
}

// getPublishedPosts retrieves all published posts
func (s *SitemapService) getPublishedPosts() ([]models.Post, error) {
	var posts []models.Post
	err := s.db.Where("status = ?", "published").All(&posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// getPublishedPages retrieves all published pages
func (s *SitemapService) getPublishedPages() ([]models.Page, error) {
	var pages []models.Page
	err := s.db.Where("status = ?", "published").All(&pages)
	if err != nil {
		return nil, err
	}
	return pages, nil
}

// getCategories retrieves all categories
func (s *SitemapService) getCategories() ([]models.Category, error) {
	var categories []models.Category
	err := s.db.All(&categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}



// GetURLCount returns the total number of URLs in the sitemap
func (s *SitemapService) GetURLCount() (int, error) {
	postCount, err := s.getPostCount()
	if err != nil {
		return 0, err
	}

	pageCount, err := s.getPageCount()
	if err != nil {
		return 0, err
	}

	categoryCount, err := s.getCategoryCount()
	if err != nil {
		return 0, err
	}

	// +1 for homepage
	return postCount + pageCount + categoryCount + 1, nil
}

// getPostCount returns the count of published posts
func (s *SitemapService) getPostCount() (int, error) {
	count, err := s.db.Where("status = ?", "published").Count(&models.Post{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// getPageCount returns the count of published pages
func (s *SitemapService) getPageCount() (int, error) {
	count, err := s.db.Where("status = ?", "published").Count(&models.Page{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// getCategoryCount returns the count of categories
func (s *SitemapService) getCategoryCount() (int, error) {
	count, err := s.db.Count(&models.Category{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Rebuild regenerates the sitemap with fresh data
func (s *SitemapService) Rebuild() (string, error) {
	log.Println("Rebuilding sitemap...")
	xmlStr, err := s.Generate()
	if err != nil {
		return "", err
	}
	log.Println("Sitemap rebuilt successfully")
	return xmlStr, nil
}
