package api

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
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
	out.Close()

	// 图片压缩处理
	finalPath := filePath
	finalSize := header.Size
	url := "/uploads/" + filename

	if cfg.Upload.ImageCompress && isImageFile(ext) {
		compressedPath, compressedSize, compressErr := compressImage(filePath, ext, cfg.Upload.MaxWidth, cfg.Upload.JPEGQuality)
		if compressErr == nil && compressedPath != "" {
			// 压缩成功，删除原文件
			os.Remove(filePath)
			finalPath = compressedPath
			finalSize = compressedSize
			filename = filepath.Base(compressedPath)
			url = "/uploads/" + filename
		}
		// 压缩失败时静默使用原图，不阻断上传流程
	}

	// 创建媒体记录
	media := models.Media{
		Filename:   header.Filename,
		Path:       finalPath,
		URL:        url,
		FileSize:   finalSize,
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

// UploadCover 上传封面图（自动压缩缩放）
// 封面图会自动缩放到 1200px 宽度，保持比例，质量压缩到 85%
func UploadCover(c *gin.Context) {
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
	allowedImageExts := []string{".jpg", ".jpeg", ".png", ".webp"}
	isAllowed := false
	for _, e := range allowedImageExts {
		if e == ext {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		BadRequest(c, "封面图仅支持 jpg、png、webp 格式")
		return
	}

	// 解码图片
	var img image.Image
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	case ".webp":
		// webp 需要额外库支持，这里简化处理
		img, _, err = image.Decode(file)
	default:
		img, _, err = image.Decode(file)
	}
	if err != nil {
		BadRequest(c, "图片格式无效，请上传有效的图片文件")
		return
	}

	// 缩放图片：使用配置的最大宽度
	coverMaxWidth := cfg.Upload.MaxWidth
	if coverMaxWidth <= 0 {
		coverMaxWidth = 1920
	}
	if img.Bounds().Dx() > coverMaxWidth {
		img = imaging.Resize(img, coverMaxWidth, 0, imaging.Lanczos)
	}

	// JPEG 质量
	jpegQ := cfg.Upload.JPEGQuality
	if jpegQ < 60 {
		jpegQ = 85
	}
	if jpegQ > 100 {
		jpegQ = 85
	}

	// 生成唯一文件名（统一保存为 jpg 格式以获得更好的压缩）
	filename := fmt.Sprintf("cover_%d_%d.jpg", time.Now().Unix(), header.Size)
	filePath := filepath.Join(uploadPath, filename)

	// 创建输出文件
	out, err := os.Create(filePath)
	if err != nil {
		InternalError(c, "保存文件失败")
		return
	}
	defer out.Close()

	// 以配置质量编码为 JPEG
	if err := jpeg.Encode(out, img, &jpeg.Options{Quality: jpegQ}); err != nil {
		InternalError(c, "图片压缩失败")
		return
	}

	// 获取压缩后的文件大小
	fileInfo, _ := os.Stat(filePath)
	compressedSize := fileInfo.Size()

	// 生成 URL
	url := "/uploads/" + filename

	// 创建媒体记录
	media := models.Media{
		Filename:   header.Filename,
		Path:       filePath,
		URL:        url,
		FileSize:   compressedSize,
		FileType:   "image",
		UploadedBy: currentUser.ID,
	}

	if err := models.DB.Create(&media).Error; err != nil {
		InternalError(c, "保存媒体记录失败")
		return
	}

	SuccessWithMessage(c, "封面图上传成功", gin.H{
		"id":              media.ID,
		"filename":        media.Filename,
		"url":             media.URL,
		"file_size":       media.FileSize,
		"original_size":   header.Size,
		"compressed_size": compressedSize,
		"width":           img.Bounds().Dx(),
		"height":          img.Bounds().Dy(),
	})
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

// isImageFile 判断文件扩展名是否为图片格式
func isImageFile(ext string) bool {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg", ".png", ".webp":
		return true
	default:
		return false
	}
}

// compressImage 压缩图片：超宽缩放 + JPEG 质量压缩
// 返回压缩后的路径、大小、错误（如果不需要压缩则返回原路径和 nil）
func compressImage(srcPath, ext string, maxWidth, quality int) (string, int64, error) {
	// 打开原文件
	f, err := os.Open(srcPath)
	if err != nil {
		return srcPath, 0, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	// 解码图片
	var img image.Image
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(f)
	case ".png":
		img, err = png.Decode(f)
	default:
		img, _, err = image.Decode(f)
	}
	if err != nil {
		return srcPath, 0, fmt.Errorf("decode image: %w", err)
	}

	originalW := img.Bounds().Dx()
	needsResize := maxWidth > 0 && originalW > maxWidth

	// 缩放（如果需要）
	if needsResize {
		img = imaging.Resize(img, maxWidth, 0, imaging.Lanczos)
	}

	// 质量校验
	if quality < 60 {
		quality = 60
	}
	if quality > 100 {
		quality = 100
	}

	// 生成输出路径（统一保存为 jpg 以获得最佳压缩）
	outExt := ".jpg"
	outPath := srcPath
	if strings.ToLower(ext) != ".jpg" && strings.ToLower(ext) != ".jpeg" {
		outPath = strings.TrimSuffix(srcPath, ext) + outExt
	}

	// 创建输出文件
	out, err := os.Create(outPath)
	if err != nil {
		return srcPath, 0, fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	// 编码为 JPEG
	if err := jpeg.Encode(out, img, &jpeg.Options{Quality: quality}); err != nil {
		os.Remove(outPath) // 清理失败文件
		return srcPath, 0, fmt.Errorf("encode jpeg: %w", err)
	}

	// 获取压缩后文件信息
	info, err := os.Stat(outPath)
	if err != nil {
		return outPath, 0, nil
	}

	return outPath, info.Size(), nil
}
