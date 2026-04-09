package api

import (
	"bytes"
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
	"catter-cms-gin/internal/storage"
)

// 全局存储提供者（由 main.go 初始化）
var storageProvider storage.StorageProvider

// InitStorage 初始化存储提供者
func InitStorage(cfg *config.Config, localBasePath string) error {
	// 将 config.Storage 转换为 storage.StorageConfig
	sc := storage.StorageConfig{
		Type: cfg.Storage.Type,
		CDN:  cfg.Storage.CDN,
		Qiniu: storage.QiniuConfig{
			AccessKey: cfg.Storage.Qiniu.AccessKey,
			SecretKey: cfg.Storage.Qiniu.SecretKey,
			Bucket:    cfg.Storage.Qiniu.Bucket,
			Zone:      cfg.Storage.Qiniu.Zone,
			Domain:    cfg.Storage.Qiniu.Domain,
			UseHTTPS:  cfg.Storage.Qiniu.UseHTTPS,
		},
		AliOSS: storage.AliOSSConfig{
			Endpoint:        cfg.Storage.AliOSS.Endpoint,
			AccessKeyID:     cfg.Storage.AliOSS.AccessKeyID,
			AccessKeySecret: cfg.Storage.AliOSS.AccessKeySecret,
			Bucket:          cfg.Storage.AliOSS.Bucket,
			CustomDomain:    cfg.Storage.AliOSS.CustomDomain,
		},
		TencentCOS: storage.TencentCOSConfig{
			BucketURL:    cfg.Storage.TencentCOS.BucketURL,
			SecretID:     cfg.Storage.TencentCOS.SecretID,
			SecretKey:    cfg.Storage.TencentCOS.SecretKey,
			Region:       cfg.Storage.TencentCOS.Region,
			CustomDomain: cfg.Storage.TencentCOS.CustomDomain,
		},
		S3: storage.S3Config{
			Endpoint:        cfg.Storage.S3.Endpoint,
			Region:          cfg.Storage.S3.Region,
			AccessKeyID:     cfg.Storage.S3.AccessKeyID,
			SecretAccessKey: cfg.Storage.S3.SecretAccessKey,
			Bucket:          cfg.Storage.S3.Bucket,
			ForcePathStyle:  cfg.Storage.S3.ForcePathStyle,
			CustomDomain:    cfg.Storage.S3.CustomDomain,
		},
	}

	provider, err := storage.NewProvider(sc, localBasePath)
	if err != nil {
		return err
	}
	storageProvider = provider
	// 同时设置 storage 包的全局变量（供模板函数使用）
	storage.SetProvider(provider)
	return nil
}

// GetStorage 获取全局存储提供者
func GetStorage() storage.StorageProvider {
	return storageProvider
}

// generateObjectKey 生成带日期目录的对象存储路径
// 格式: 2026/04/1234567890_12345.jpg
func generateObjectKey(filename string) string {
	now := time.Now()
	ext := filepath.Ext(filename)
	uniqueName := fmt.Sprintf("%d_%d%s", now.Unix(), time.Now().UnixNano()%10000, ext)
	return fmt.Sprintf("%d/%02d/%s", now.Year(), now.Month(), uniqueName)
}

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
	maxSize := cfg.Upload.MaxSize

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

	// 生成对象存储路径
	objectKey := generateObjectKey(header.Filename)

	// 获取文件内容（可能需要压缩处理）
	var uploadReader io.Reader = file
	finalSize := header.Size

	// 图片压缩处理
	if cfg.Upload.ImageCompress && isImageFile(ext) {
		// 重置读取位置
		file.Seek(0, 0)
		compressedReader, compressedSize, compressErr := compressImageToReader(file, ext, cfg.Upload.MaxWidth, cfg.Upload.JPEGQuality)
		if compressErr == nil && compressedReader != nil {
			uploadReader = compressedReader
			finalSize = compressedSize
			objectKey = strings.TrimSuffix(objectKey, ext) + ".jpg"
		}
	}

	// 上传到存储
	url, err := storageProvider.Put(objectKey, uploadReader, finalSize, getContentType(ext))
	if err != nil {
		InternalError(c, "文件上传失败: "+err.Error())
		return
	}

	// 生成缩略图 URL
	var thumbnailURL string
	if isImageFile(ext) || isImageFile(".jpg") {
		if storage.IsObjectStorage(storageProvider) {
			// 对象存储：通过 URL 参数实现缩略图
			thumbnailURL = storageProvider.GetThumbURL(objectKey)
		} else {
			// 本地存储：保存缩略图文件
			thumbnailURL = generateLocalThumbnail(objectKey, cfg.Upload.Path)
		}
	}

	// 创建媒体记录
	media := models.Media{
		Filename:     header.Filename,
		Path:         objectKey,
		URL:          url,
		ThumbnailURL: thumbnailURL,
		FileSize:     finalSize,
		UploadedBy:   currentUser.ID,
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

	// 从存储中删除文件
	if err := storageProvider.Delete(media.Path); err != nil {
		// 对象存储删除失败只记录日志，不阻断
		fmt.Printf("删除存储文件失败: %v\n", err)
	}

	// 本地存储：额外删除缩略图文件
	if storageProvider.Type() == "local" && media.ThumbnailURL != "" {
		thumbRelPath := strings.TrimPrefix(media.ThumbnailURL, "/")
		if _, err := os.Stat(thumbRelPath); err == nil {
			os.Remove(thumbRelPath)
		}
	}

	// 删除数据库记录
	if err := models.DB.Delete(&media).Error; err != nil {
		InternalError(c, "删除媒体失败")
		return
	}

	SuccessWithMessage(c, "媒体删除成功", nil)
}

// BatchDeleteMedia 批量删除媒体
func BatchDeleteMedia(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		BadRequest(c, "请提供要删除的媒体ID列表")
		return
	}

	var medias []models.Media
	if err := models.DB.Where("id IN ?", req.IDs).Find(&medias).Error; err != nil {
		InternalError(c, "查询媒体失败")
		return
	}

	successCount := 0
	failCount := 0
	for _, media := range medias {
		// 从存储中删除文件
		if err := storageProvider.Delete(media.Path); err != nil {
			fmt.Printf("删除存储文件失败: %v\n", err)
		}

		// 本地存储：额外删除缩略图文件
		if storageProvider.Type() == "local" && media.ThumbnailURL != "" {
			thumbRelPath := strings.TrimPrefix(media.ThumbnailURL, "/")
			if _, err := os.Stat(thumbRelPath); err == nil {
				os.Remove(thumbRelPath)
			}
		}

		// 删除数据库记录
		if err := models.DB.Delete(&media).Error; err != nil {
			failCount++
			continue
		}
		successCount++
	}

	if successCount == 0 {
		InternalError(c, "删除失败，资源不存在")
		return
	}

	msg := fmt.Sprintf("成功删除 %d 个文件", successCount)
	if failCount > 0 {
		msg += fmt.Sprintf("，%d 个删除失败", failCount)
	}
	SuccessWithMessage(c, msg, gin.H{"success": successCount, "failed": failCount})
}

// UploadCover 上传封面图（自动压缩缩放）
// 封面图会自动缩放到配置的最大宽度，保持比例，质量压缩
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
	maxSize := cfg.Upload.MaxSize

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
		img, _, err = image.Decode(file)
	default:
		img, _, err = image.Decode(file)
	}
	if err != nil {
		BadRequest(c, "图片格式无效，请上传有效的图片文件")
		return
	}

	// 缩放图片
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

	// 压缩为 JPEG 到内存 buffer
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: jpegQ}); err != nil {
		InternalError(c, "图片压缩失败")
		return
	}

	compressedSize := int64(buf.Len())

	// 生成对象存储路径（封面图前缀 + 日期目录）
	objectKey := fmt.Sprintf("cover/%d/%02d/cover_%d.jpg",
		time.Now().Year(), time.Now().Month(), time.Now().UnixNano())

	// 上传到存储
	url, err := storageProvider.Put(objectKey, &buf, compressedSize, "image/jpeg")
	if err != nil {
		InternalError(c, "封面图上传失败: "+err.Error())
		return
	}

	// 生成缩略图 URL
	var thumbnailURL string
	if storage.IsObjectStorage(storageProvider) {
		thumbnailURL = storageProvider.GetThumbURL(objectKey)
	} else {
		thumbnailURL = generateLocalThumbnail(objectKey, cfg.Upload.Path)
	}

	// 创建媒体记录
	media := models.Media{
		Filename:     header.Filename,
		Path:         objectKey,
		URL:          url,
		ThumbnailURL: thumbnailURL,
		FileSize:     compressedSize,
		FileType:     "image",
		UploadedBy:   currentUser.ID,
	}

	if err := models.DB.Create(&media).Error; err != nil {
		InternalError(c, "保存媒体记录失败")
		return
	}

	SuccessWithMessage(c, "封面图上传成功", gin.H{
		"id":              media.ID,
		"filename":        media.Filename,
		"url":             media.URL,
		"thumbnail_url":   media.ThumbnailURL,
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

// getContentType 根据文件扩展名获取 Content-Type
func getContentType(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
}

// compressImageToReader 压缩图片到内存 reader
// 返回压缩后的 reader、大小、错误
func compressImageToReader(r io.Reader, ext string, maxWidth, quality int) (io.Reader, int64, error) {
	// 解码图片
	var img image.Image
	var err error
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(r)
	case ".png":
		img, err = png.Decode(r)
	default:
		img, _, err = image.Decode(r)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("decode image: %w", err)
	}

	// 缩放
	if maxWidth > 0 && img.Bounds().Dx() > maxWidth {
		img = imaging.Resize(img, maxWidth, 0, imaging.Lanczos)
	}

	// 质量校验
	if quality < 60 {
		quality = 60
	}
	if quality > 100 {
		quality = 100
	}

	// 编码到 buffer
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, 0, fmt.Errorf("encode jpeg: %w", err)
	}

	return &buf, int64(buf.Len()), nil
}

// generateLocalThumbnail 生成本地缩略图文件（仅本地存储使用）
func generateLocalThumbnail(objectKey, uploadPath string) string {
	// 从 objectKey 构建本地文件路径
	localPath := filepath.Join(uploadPath, objectKey)
	f, err := os.Open(localPath)
	if err != nil {
		return ""
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return ""
	}

	// 无论原始图尺寸，都生成缩略图（小图直接复制）
	if img.Bounds().Dx() <= thumbnailWidth {
		// 原始图已经很小，直接复制为缩略图
		thumb := imaging.Clone(img)

		thumbDir := filepath.Join(uploadPath, ".thumbs")
		if err := os.MkdirAll(thumbDir, 0755); err != nil {
			return ""
		}

		thumbFilename := "thumb_" + strings.ReplaceAll(objectKey, "/", "_")
		thumbFilename = strings.TrimSuffix(thumbFilename, filepath.Ext(thumbFilename)) + ".jpg"
		thumbPath := filepath.Join(thumbDir, thumbFilename)

		out, err := os.Create(thumbPath)
		if err != nil {
			return ""
		}
		defer out.Close()

		if err := jpeg.Encode(out, thumb, &jpeg.Options{Quality: 75}); err != nil {
			os.Remove(thumbPath)
			return ""
		}

		return "/uploads/.thumbs/" + thumbFilename
	}

	thumb := imaging.Resize(img, thumbnailWidth, 0, imaging.Lanczos)

	thumbDir := filepath.Join(uploadPath, ".thumbs")
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return ""
	}

	thumbFilename := "thumb_" + strings.ReplaceAll(objectKey, "/", "_")
	thumbFilename = strings.TrimSuffix(thumbFilename, filepath.Ext(thumbFilename)) + ".jpg"
	thumbPath := filepath.Join(thumbDir, thumbFilename)

	out, err := os.Create(thumbPath)
	if err != nil {
		return ""
	}
	defer out.Close()

	if err := jpeg.Encode(out, thumb, &jpeg.Options{Quality: 75}); err != nil {
		os.Remove(thumbPath)
		return ""
	}

	return "/uploads/.thumbs/" + thumbFilename
}

const thumbnailWidth = 400
