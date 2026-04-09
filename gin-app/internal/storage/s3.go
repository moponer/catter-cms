package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Storage AWS S3 / 兼容存储（Cloudflare R2, MinIO 等）
type S3Storage struct {
	cfg    S3Config
	client *s3.Client
}

// NewS3Storage 创建 S3 存储实例
func NewS3Storage(cfg S3Config) (*S3Storage, error) {
	if cfg.Endpoint == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("S3 配置不完整: 需要 endpoint, access_key_id, secret_access_key, bucket")
	}

	// 自定义配置加载
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           cfg.Endpoint,
			SigningRegion: cfg.Region,
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("创建 S3 配置失败: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.ForcePathStyle {
			o.UsePathStyle = true
		}
	})

	return &S3Storage{
		cfg:    cfg,
		client: client,
	}, nil
}

// Put 上传文件到 S3
func (s *S3Storage) Put(objectKey string, reader io.Reader, size int64, contentType string) (string, error) {
	contentTypePtr := &contentType
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.cfg.Bucket),
		Key:         aws.String(objectKey),
		Body:        reader,
		ContentType: contentTypePtr,
		ACL:         "public-read",
	}
	if size > 0 {
		input.ContentLength = aws.Int64(size)
	}

	_, err := s.client.PutObject(context.Background(), input)
	if err != nil {
		return "", fmt.Errorf("S3 上传失败: %w", err)
	}

	return s.GetURL(objectKey), nil
}

// Delete 从 S3 删除文件
func (s *S3Storage) Delete(objectKey string) error {
	_, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(objectKey),
	})
	return err
}

// GetURL 获取文件的访问 URL
func (s *S3Storage) GetURL(objectKey string) string {
	if s.cfg.CustomDomain != "" {
		domain := strings.TrimRight(s.cfg.CustomDomain, "/")
		if !strings.HasPrefix(domain, "http") {
			domain = "https://" + domain
		}
		return domain + "/" + objectKey
	}

	// 使用 S3 默认格式
	endpoint := strings.TrimPrefix(s.cfg.Endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	if s.cfg.ForcePathStyle {
		// Path style: https://endpoint/bucket/key
		return fmt.Sprintf("https://%s/%s/%s", endpoint, s.cfg.Bucket, objectKey)
	}
	// Virtual-hosted style: https://bucket.endpoint/key
	return fmt.Sprintf("https://%s.%s/%s", s.cfg.Bucket, endpoint, objectKey)
}

// GetThumbURL 获取缩略图 URL
// S3 本身不支持 URL 图片处理，返回原图 URL
// 如果使用 Cloudflare R2 等兼容存储，可通过其 CDN 的 Image Resizing 功能实现
func (s *S3Storage) GetThumbURL(objectKey string) string {
	if objectKey == "" {
		return ""
	}
	// S3 标准不支持图片处理参数，直接返回原图
	// 如果用户使用 Cloudflare R2 + Image Resizing，可通过自定义域名 + /cdn-cgi/image/width=400/... 实现
	return s.GetURL(objectKey)
}

// Type 返回存储类型
func (s *S3Storage) Type() string {
	return "s3"
}
