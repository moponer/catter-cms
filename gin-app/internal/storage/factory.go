package storage

import (
	"fmt"
	"log"
	"strings"
)

// currentProvider 全局存储提供者实例
var currentProvider StorageProvider

// SetProvider 设置全局存储提供者（由 main.go 调用）
func SetProvider(p StorageProvider) {
	currentProvider = p
}

// GetProvider 获取全局存储提供者
func GetProvider() StorageProvider {
	return currentProvider
}

// GetGlobalType 获取全局存储类型
func GetGlobalType() string {
	if currentProvider == nil {
		return "local"
	}
	return currentProvider.Type()
}

// GlobalGetThumbURL 全局缩略图 URL 生成（供模板函数使用）
// 本地存储: /uploads/xxx.jpg → /uploads/.thumbs/thumb_xxx.jpg
// 对象存储: https://cdn.example.com/xxx.jpg → https://cdn.example.com/xxx.jpg?imageView2/w/400/q/75
func GlobalGetThumbURL(originalURL string) string {
	if originalURL == "" {
		return ""
	}
	if currentProvider == nil {
		return originalURL
	}

	if currentProvider.Type() == "local" {
		// 本地存储：路径转换
		if !strings.HasPrefix(originalURL, "/uploads/") {
			return originalURL
		}
		filename := strings.TrimPrefix(originalURL, "/uploads/")
		// 文件名生成规则必须与 generateLocalThumbnail 一致：
		// thumb_ + objectKey 中 / 替换为 _ + 去掉原扩展名 + .jpg
		thumbFilename := "thumb_" + strings.ReplaceAll(filename, "/", "_")
		thumbFilename = strings.TrimSuffix(thumbFilename, extOf(thumbFilename)) + ".jpg"
		return "/uploads/.thumbs/" + thumbFilename
	}

	// 对象存储：从完整 URL 提取 objectKey，再传给 provider
	objectKey := extractObjectKey(originalURL)
	if objectKey == "" {
		// 无法提取 objectKey，直接在 URL 后追加缩略图参数
		return originalURL + "?imageView2/w/400/q/75"
	}
	return currentProvider.GetThumbURL(objectKey)
}

// extractObjectKey 从完整 URL 中提取 objectKey
// https://cdn.example.com/cover/2026/04/xxx.jpg → cover/2026/04/xxx.jpg
func extractObjectKey(url string) string {
	// 移除协议
	noProto := url
	if idx := strings.Index(noProto, "://"); idx >= 0 {
		noProto = noProto[idx+3:]
	}
	// 跳过域名，取第一个 / 之后的部分
	if idx := strings.Index(noProto, "/"); idx >= 0 {
		return noProto[idx+1:]
	}
	return ""
}

func extOf(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}

// NewProvider 根据配置创建对应的存储提供者
// 如果 storageType 为空或 "local"，使用本地存储
func NewProvider(cfg StorageConfig, localBasePath string) (StorageProvider, error) {
	switch cfg.Type {
	case "", "local":
		return NewLocalStorage(localBasePath), nil

	case "qiniu":
		provider, err := NewQiniuStorage(cfg.Qiniu)
		if err != nil {
			return nil, fmt.Errorf("初始化七牛云存储失败: %w", err)
		}
		log.Printf("Storage: 七牛云 (bucket=%s, zone=%s)", cfg.Qiniu.Bucket, cfg.Qiniu.Zone)
		return provider, nil

	case "alioss":
		provider, err := NewAliOSSStorage(cfg.AliOSS)
		if err != nil {
			return nil, fmt.Errorf("初始化阿里云 OSS 失败: %w", err)
		}
		log.Printf("Storage: 阿里云 OSS (bucket=%s, endpoint=%s)", cfg.AliOSS.Bucket, cfg.AliOSS.Endpoint)
		return provider, nil

	case "tencentcos":
		provider, err := NewTencentCOSStorage(cfg.TencentCOS)
		if err != nil {
			return nil, fmt.Errorf("初始化腾讯云 COS 失败: %w", err)
		}
		log.Printf("Storage: 腾讯云 COS (bucket_url=%s)", cfg.TencentCOS.BucketURL)
		return provider, nil

	case "s3":
		provider, err := NewS3Storage(cfg.S3)
		if err != nil {
			return nil, fmt.Errorf("初始化 S3 存储失败: %w", err)
		}
		log.Printf("Storage: S3 兼容存储 (bucket=%s, endpoint=%s)", cfg.S3.Bucket, cfg.S3.Endpoint)
		return provider, nil

	default:
		return nil, fmt.Errorf("不支持的存储类型: %s（支持: local, qiniu, alioss, tencentcos, s3）", cfg.Type)
	}
}

// IsObjectStorage 判断是否为对象存储（非本地）
func IsObjectStorage(provider StorageProvider) bool {
	return provider != nil && provider.Type() != "local"
}

// GenerateObjectKey 生成带日期分目录的对象存储路径
// 格式: 2026/04/abc123.jpg（年/月/文件名）
func GenerateObjectKey(filename string) string {
	// 简单实现：直接返回文件名，日期目录由调用者拼接
	// 这样更灵活，上传者可以自行决定目录结构
	return filename
}
