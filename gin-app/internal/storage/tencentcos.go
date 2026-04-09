package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// TencentCOSStorage 腾讯云 COS 存储
type TencentCOSStorage struct {
	cfg    TencentCOSConfig
	client *cos.Client
}

// NewTencentCOSStorage 创建腾讯云 COS 存储实例
func NewTencentCOSStorage(cfg TencentCOSConfig) (*TencentCOSStorage, error) {
	if cfg.BucketURL == "" || cfg.SecretID == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("腾讯云 COS 配置不完整: 需要 bucket_url, secret_id, secret_key")
	}

	bucketURL, err := url.Parse(cfg.BucketURL)
	if err != nil {
		return nil, fmt.Errorf("bucket_url 格式无效: %w", err)
	}

	client := cos.NewClient(&cos.BaseURL{BucketURL: bucketURL}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})

	return &TencentCOSStorage{
		cfg:    cfg,
		client: client,
	}, nil
}

// Put 上传文件到腾讯云 COS
func (s *TencentCOSStorage) Put(objectKey string, reader io.Reader, size int64, contentType string) (string, error) {
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
	}

	_, err := s.client.Object.Put(context.Background(), objectKey, reader, opt)
	if err != nil {
		return "", fmt.Errorf("COS 上传失败: %w", err)
	}

	return s.GetURL(objectKey), nil
}

// Delete 从腾讯云 COS 删除文件
func (s *TencentCOSStorage) Delete(objectKey string) error {
	_, err := s.client.Object.Delete(context.Background(), objectKey)
	return err
}

// GetURL 获取文件的访问 URL
func (s *TencentCOSStorage) GetURL(objectKey string) string {
	if s.cfg.CustomDomain != "" {
		domain := strings.TrimRight(s.cfg.CustomDomain, "/")
		if !strings.HasPrefix(domain, "http") {
			domain = "https://" + domain
		}
		return domain + "/" + objectKey
	}
	// 使用 BucketURL
	bucketURL := strings.TrimRight(s.cfg.BucketURL, "/")
	return bucketURL + "/" + objectKey
}

// GetThumbURL 获取缩略图 URL（使用 COS 数据万象图片处理）
// 通过 URL 参数实现，如 ?imageView2/w/400/q/75
func (s *TencentCOSStorage) GetThumbURL(objectKey string) string {
	if objectKey == "" {
		return ""
	}
	return s.GetURL(objectKey) + "?imageView2/w/400/q/75"
}

// Type 返回存储类型
func (s *TencentCOSStorage) Type() string {
	return "tencentcos"
}
