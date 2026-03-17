package admin

import (
	"catter-cms/actions/api"
	"catter-cms/models"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gofrs/uuid"
)

type MediaHandler struct{}

// MediaListAPI 获取媒体列表API
func MediaListAPI(c buffalo.Context) error {
	// 从查询字符串中获取参数
	page, _ := strconv.Atoi(c.Request().URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(c.Request().URL.Query().Get("page_size"))
	keyword := c.Request().URL.Query().Get("keyword")
	fileType := c.Request().URL.Query().Get("type")

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 30
	}

	media := &models.MediaItems{}
	query := models.DB.Paginate(page, pageSize)
	query = query.Order("created_at DESC")

	if keyword != "" {
		query = query.Where("filename LIKE ?", "%"+keyword+"%")
	}

	if fileType != "" {
		mimeType := ""
		switch fileType {
		case "image":
			mimeType = "image/%"
		case "video":
			mimeType = "video/%"
		case "document":
			mimeType = "application/%"
		}
		if mimeType != "" {
			query = query.Where("mime_type LIKE ?", mimeType)
		}
	}

	if err := query.All(media); err != nil {
		return api.InternalError(c, "获取媒体列表失败")
	}

	total, _ := models.DB.Count(&models.Media{})

	return api.PagedSuccess(c, media, total, page, pageSize)
}

// MediaUploadAPI 上传文件API
func MediaUploadAPI(c buffalo.Context) error {
	file, header, err := c.Request().FormFile("file")
	if err != nil {
		return api.BadRequest(c, "获取文件失败")
	}
	defer file.Close()

	// 验证文件大小 (限制 10MB)
	if header.Size > 10*1024*1024 {
		return api.BadRequest(c, "文件大小不能超过10MB")
	}

	mimeType := getMimeType(header.Filename)

	uploadDir := "public/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return api.InternalError(c, "创建上传目录失败")
	}

	ext := filepath.Ext(header.Filename)
	newFilename := uuid.Must(uuid.NewV4()).String() + ext
	filePath := filepath.Join(uploadDir, newFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		return api.InternalError(c, "创建文件失败")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return api.InternalError(c, "保存文件失败")
	}

	media := &models.Media{
		ID:        uuid.Must(uuid.NewV4()),
		Filename:  header.Filename,
		Path:      "/uploads/" + newFilename,
		Size:      header.Size,
		MimeType:  mimeType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := models.DB.Create(media); err != nil {
		os.Remove(filePath)
		return api.InternalError(c, "保存媒体记录失败")
	}

	return api.SuccessWithMessage(c, "文件上传成功", media)
}

// MediaDeleteAPI 删除媒体文件API
func MediaDeleteAPI(c buffalo.Context) error {
	mediaID := c.Param("id")
	media := &models.Media{}

	if err := models.DB.Find(media, mediaID); err != nil {
		return api.NotFound(c, "媒体文件不存在")
	}

	filePath := "public" + media.Path
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return api.InternalError(c, "删除文件失败")
	}

	if err := models.DB.Destroy(media); err != nil {
		return api.InternalError(c, "删除媒体记录失败")
	}

	return api.SuccessWithMessage(c, "文件删除成功", nil)
}

// MediaDetailAPI 获取媒体详情API
func MediaDetailAPI(c buffalo.Context) error {
	mediaID := c.Param("id")
	media := &models.Media{}

	if err := models.DB.Find(media, mediaID); err != nil {
		return api.NotFound(c, "媒体文件不存在")
	}

	return api.Success(c, media)
}

// MediaBatchDeleteAPI 批量删除媒体文件API
func MediaBatchDeleteAPI(c buffalo.Context) error {
	ids := []string{}
	if err := c.Bind(&ids); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	for _, id := range ids {
		media := &models.Media{}
		if err := models.DB.Find(media, id); err != nil {
			continue
		}

		filePath := "public" + media.Path
		os.Remove(filePath)
		models.DB.Destroy(media)
	}

	return api.SuccessWithMessage(c, "批量删除成功", nil)
}

// getFileType 根据文件名获取文件类型
func getFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	videoExts := []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv"}
	documentExts := []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt"}

	for _, imgExt := range imageExts {
		if ext == imgExt {
			return "image"
		}
	}

	for _, vidExt := range videoExts {
		if ext == vidExt {
			return "video"
		}
	}

	for _, docExt := range documentExts {
		if ext == docExt {
			return "document"
		}
	}

	return "other"
}

// getMimeType 根据文件名获取MIME类型
func getMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	mimeMap := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
		".webp": "image/webp",
		".mp4":  "video/mp4",
		".avi":  "video/x-msvideo",
		".mov":  "video/quicktime",
		".wmv":  "video/x-ms-wmv",
		".flv":  "video/x-flv",
		".mkv":  "video/x-matroska",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".txt":  "text/plain",
	}

	if mimeType, ok := mimeMap[ext]; ok {
		return mimeType
	}

	return "application/octet-stream"
}

// MediaIndexHandler 媒体文件列表页面 (已弃用，SPA模式)
func MediaIndexHandler(c buffalo.Context) error {
	return c.Render(200, render.HTML("<html><body><p>Media SPA</p></body></html>"))
}

// MediaUploadHandler 上传文件 (已弃用，使用MediaUploadAPI)
func MediaUploadHandler(c buffalo.Context) error {
	return MediaUploadAPI(c)
}

// MediaDeleteHandler 删除媒体文件 (已弃用，使用MediaDeleteAPI)
func MediaDeleteHandler(c buffalo.Context) error {
	return MediaDeleteAPI(c)
}

// unused - avoid compile error
var _ = fmt.Sprintf
