package storage

import (
	"fmt"
	"io"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// AliOSSStorage 阿里云 OSS 存储
type AliOSSStorage struct {
	cfg    AliOSSConfig
	client *oss.Client
	bucket *oss.Bucket
}

// NewAliOSSStorage 创建阿里云 OSS 存储实例
func NewAliOSSStorage(cfg AliOSSConfig) (*AliOSSStorage, error) {
	if cfg.Endpoint == "" || cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("阿里云 OSS 配置不完整: 需要 endpoint, access_key_id, access_key_secret, bucket")
	}

	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建 OSS 客户端失败: %w", err)
	}

	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("获取 Bucket 失败: %w", err)
	}

	return &AliOSSStorage{
		cfg:    cfg,
		client: client,
		bucket: bucket,
	}, nil
}

// Put 上传文件到阿里云 OSS
func (s *AliOSSStorage) Put(objectKey string, reader io.Reader, size int64, contentType string) (string, error) {
	options := []oss.Option{
		oss.ContentType(contentType),
		oss.ObjectACL(oss.ACLPublicRead),
	}
	if size > 0 {
		options = append(options, oss.ContentLength(size))
	}

	err := s.bucket.PutObject(objectKey, reader, options...)
	if err != nil {
		return "", fmt.Errorf("OSS 上传失败: %w", err)
	}

	return s.GetURL(objectKey), nil
}

// Delete 从阿里云 OSS 删除文件
func (s *AliOSSStorage) Delete(objectKey string) error {
	return s.bucket.DeleteObject(objectKey)
}

// GetURL 获取文件的访问 URL
func (s *AliOSSStorage) GetURL(objectKey string) string {
	if s.cfg.CustomDomain != "" {
		domain := strings.TrimRight(s.cfg.CustomDomain, "/")
		if !strings.HasPrefix(domain, "http") {
			domain = "https://" + domain
		}
		return domain + "/" + objectKey
	}
	// 使用 OSS 默认域名
	// 格式: https://bucket.endpoint/key
	endpoint := strings.TrimPrefix(s.cfg.Endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	return fmt.Sprintf("https://%s.%s/%s", s.cfg.Bucket, endpoint, objectKey)
}

// GetThumbURL 获取缩略图 URL（使用 OSS 图片处理）
// 通过 URL 参数实现，如 ?x-oss-process=image/resize,w_400/quality,q_75
func (s *AliOSSStorage) GetThumbURL(objectKey string) string {
	if objectKey == "" {
		return ""
	}
	return s.GetURL(objectKey) + "?x-oss-process=image/resize,w_400/quality,q_75"
}

// Type 返回存储类型
func (s *AliOSSStorage) Type() string {
	return "alioss"
}
