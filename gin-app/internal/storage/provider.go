package storage

import "io"

// StorageProvider 存储提供者接口
// 所有存储后端（本地、七牛、阿里 OSS、腾讯 COS、S3）都必须实现此接口
type StorageProvider interface {
	// Put 上传文件，返回访问 URL 和错误
	// objectKey: 存储路径（如 "2026/04/abc123.jpg"）
	Put(objectKey string, reader io.Reader, size int64, contentType string) (url string, err error)

	// Delete 删除文件
	Delete(objectKey string) error

	// GetURL 获取文件的访问 URL
	GetURL(objectKey string) string

	// GetThumbURL 获取缩略图 URL
	// 对象存储通过 URL 参数实现（如七牛的 ?imageView2/w/400）
	// 本地存储通过 .thumbs 目录实现
	GetThumbURL(objectKey string) string

	// Type 返回存储类型标识（local, qiniu, alioss, tencentcos, s3）
	Type() string
}

// StorageConfig 存储配置（统一入口）
type StorageConfig struct {
	Type        string          `yaml:"type"`                  // local, qiniu, alioss, tencentcos, s3
	CDN         string          `yaml:"cdn"`                   // CDN 域名（对象存储时使用）
	Qiniu       QiniuConfig     `yaml:"qiniu,omitempty"`
	AliOSS      AliOSSConfig    `yaml:"alioss,omitempty"`
	TencentCOS  TencentCOSConfig `yaml:"tencentcos,omitempty"`
	S3          S3Config        `yaml:"s3,omitempty"`
}

// QiniuConfig 七牛云配置
type QiniuConfig struct {
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
	Bucket     string `yaml:"bucket"`
	Zone       string `yaml:"zone"`        // 存储区域: huadong, huabei, huanan, beimei, xinjiapo
	Domain     string `yaml:"domain"`      // 访问域名（CDN 域名或七牛域名）
	UseHTTPS   bool   `yaml:"use_https"`   // 是否使用 HTTPS
}

// AliOSSConfig 阿里云 OSS 配置
type AliOSSConfig struct {
	Endpoint        string `yaml:"endpoint"`         // 如 https://oss-cn-hangzhou.aliyuncs.com
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Bucket          string `yaml:"bucket"`
	CustomDomain    string `yaml:"custom_domain"`    // 自定义域名（可选，如 CDN 域名）
}

// TencentCOSConfig 腾讯云 COS 配置
type TencentCOSConfig struct {
	BucketURL       string `yaml:"bucket_url"`       // 如 https://bucket-1234567890.cos.ap-guangzhou.myqcloud.com
	SecretID        string `yaml:"secret_id"`
	SecretKey       string `yaml:"secret_key"`
	Region          string `yaml:"region"`           // 如 ap-guangzhou
	CustomDomain    string `yaml:"custom_domain"`    // 自定义域名（可选）
}

// S3Config AWS S3 / 兼容存储配置（Cloudflare R2, MinIO 等）
type S3Config struct {
	Endpoint        string `yaml:"endpoint"`         // 如 https://s3.amazonaws.com 或 MinIO 地址
	Region          string `yaml:"region"`           // 如 us-east-1
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	Bucket          string `yaml:"bucket"`
	ForcePathStyle  bool   `yaml:"force_path_style"` // MinIO 等兼容存储需要 true
	CustomDomain    string `yaml:"custom_domain"`    // 自定义域名（可选，如 Cloudflare R2 的公开访问域名）
}
