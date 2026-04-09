package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// QiniuStorage 七牛云存储
type QiniuStorage struct {
	cfg      QiniuConfig
	mac      *qbox.Mac
	bucket   *storage.BucketManager
	formUploader *storage.FormUploader
}

// NewQiniuStorage 创建七牛云存储实例
func NewQiniuStorage(cfg QiniuConfig) (*QiniuStorage, error) {
	if cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("七牛云配置不完整: 需要 access_key, secret_key, bucket")
	}

	mac := qbox.NewMac(cfg.AccessKey, cfg.SecretKey)
	cfgZone, err := getQiniuZone(cfg.Zone)
	if err != nil {
		return nil, err
	}

	bucket := storage.NewBucketManager(mac, &storage.Config{
		Zone:          cfgZone,
		UseHTTPS:      cfg.UseHTTPS,
		UseCdnDomains: false,
	})

	formUploader := storage.NewFormUploader(&storage.Config{
		Zone:          cfgZone,
		UseHTTPS:      cfg.UseHTTPS,
		UseCdnDomains: false,
	})

	return &QiniuStorage{
		cfg:          cfg,
		mac:          mac,
		bucket:       bucket,
		formUploader: formUploader,
	}, nil
}

// Put 上传文件到七牛云
func (s *QiniuStorage) Put(objectKey string, reader io.Reader, size int64, contentType string) (string, error) {
	ret := storage.PutRet{}
	putPolicy := storage.PutPolicy{
		Scope: s.cfg.Bucket,
	}
	upToken := putPolicy.UploadToken(s.mac)

	err := s.formUploader.Put(context.Background(), &ret, upToken, objectKey, reader, size, &storage.PutExtra{
		MimeType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("七牛云上传失败: %w", err)
	}

	return s.GetURL(objectKey), nil
}

// Delete 从七牛云删除文件
func (s *QiniuStorage) Delete(objectKey string) error {
	return s.bucket.Delete(s.cfg.Bucket, objectKey)
}

// GetURL 获取文件的访问 URL
func (s *QiniuStorage) GetURL(objectKey string) string {
	domain := s.cfg.Domain
	if domain == "" {
		// 如果没配置域名，使用公开空间默认域名格式
		return fmt.Sprintf("https://%s.qiniudn.com/%s", s.cfg.Bucket, objectKey)
	}
	// 确保域名格式正确
	domain = strings.TrimRight(domain, "/")
	if !strings.HasPrefix(domain, "http") {
		if s.cfg.UseHTTPS {
			domain = "https://" + domain
		} else {
			domain = "http://" + domain
		}
	}
	encodedKey := url.PathEscape(objectKey)
	return domain + "/" + encodedKey
}

// GetThumbURL 获取缩略图 URL（使用七牛图片处理）
// 通过 URL 参数实现，如 ?imageView2/2/w/400/q/75
func (s *QiniuStorage) GetThumbURL(objectKey string) string {
	if objectKey == "" {
		return ""
	}
	return s.GetURL(objectKey) + "?imageView2/2/w/400/q/75"
}

// Type 返回存储类型
func (s *QiniuStorage) Type() string {
	return "qiniu"
}

// getQiniuZone 根据区域名称获取七牛区域配置
func getQiniuZone(zone string) (*storage.Zone, error) {
	switch zone {
	case "huadong", "":
		return &storage.ZoneHuadong, nil
	case "huabei":
		return &storage.ZoneHuabei, nil
	case "huanan":
		return &storage.ZoneHuanan, nil
	case "beimei":
		return &storage.ZoneBeimei, nil
	case "xinjiapo":
		return &storage.ZoneXinjiapo, nil
	default:
		return nil, fmt.Errorf("不支持的七牛区域: %s（支持: huadong, huabei, huanan, beimei, xinjiapo）", zone)
	}
}
