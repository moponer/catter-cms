package api

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/config"
	"catter-cms-gin/internal/models"
)

// ListMedia 获取媒体列表
func ListMedia(c *gin.Context) {
	var media []models.Media
	query := models.DB

	// 文件类型筛选
	if fileType := c.Query("file_type"); fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}

	// 搜索
	if search := c.Query("search"); search != "" {
		query = query.Where("filename LIKE ? OR title LIKE ? OR alt LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	var total int64
	query.Model(&models.Media{}).Count(&total)

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&media).Error; err != nil {
		InternalError(c, "获取媒体列表失败")
		return
	}

	SuccessResponse(c, gin.H{
		"list":      media,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetMedia 获取单个媒体
func GetMedia(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的媒体ID")
		return
	}

	var media models.Media
	if err := models.DB.First(&media, id).Error; err != nil {
		NotFound(c, "媒体不存在")
		return
	}

	SuccessResponse(c, media)
}

// UploadMedia 上传媒体
func UploadMedia(c *gin.Context) {
	// 获取当前用户
	user, exists := c.Get("current_user")
	if !exists || user == nil {
		Unauthorized(c, "未登录")
		return
	}
	currentUser := user.(*models.User)

	// 获取上传配置
	cfg := config.Get()
	uploadPath := cfg.Upload.Path
	maxSize := cfg.Upload.MaxSize

	// 确保上传目录存在
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		InternalError(c, "创建上传目录失败")
		return
	}

	// 获取文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		BadRequest(c, "请选择要上传的文件")
		return
	}
	defer file.Close()

	// 检查文件大小
	if maxSize > 0 && header.Size > maxSize {
		BadRequest(c, fmt.Sprintf("文件大小不能超过 %d MB", maxSize/1024/1024))
		return
	}

	// 检查文件扩展名
	ext := filepath.Ext(header.Filename)
	allowedExts := cfg.Upload.AllowedExt
	if allowedExts != "" {
		found := false
		for _, e := range splitAndTrim(allowedExts, ",") {
			if e == ext {
				found = true
				break
			}
		}
		if !found {
			BadRequest(c, "不支持的文件类型")
			return
		}
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%d_%d%s", time.Now().Unix(), header.Size, ext)
	filePath := filepath.Join(uploadPath, filename)

	// 写入文件
	out, err := os.Create(filePath)
	if err != nil {
		InternalError(c, "保存文件失败")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		InternalError(c, "保存文件失败")
		return
	}

	// 生成 URL
	url := "/uploads/" + filename

	// 创建媒体记录
	media := models.Media{
		Filename:   header.Filename,
		Path:       filePath,
		URL:        url,
		FileSize:   header.Size,
		UploadedBy: currentUser.ID,
	}

	if err := models.DB.Create(&media).Error; err != nil {
		InternalError(c, "保存媒体记录失败")
		return
	}

	SuccessWithMessage(c, "文件上传成功", media)
}

// DeleteMedia 删除媒体
func DeleteMedia(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的媒体ID")
		return
	}

	var media models.Media
	if err := models.DB.First(&media, id).Error; err != nil {
		NotFound(c, "媒体不存在")
		return
	}

	// 删除物理文件
	if _, err := os.Stat(media.Path); err == nil {
		os.Remove(media.Path)
	}

	// 删除数据库记录
	if err := models.DB.Delete(&media).Error; err != nil {
		InternalError(c, "删除媒体失败")
		return
	}

	SuccessWithMessage(c, "媒体删除成功", nil)
}

// splitAndTrim 分割并去除空格
func splitAndTrim(s, sep string) []string {
	parts := make([]string, 0)
	for _, p := range split(s, sep) {
		if trimmed := trim(p); trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func split(s, sep string) []string {
	if s == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trim(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
