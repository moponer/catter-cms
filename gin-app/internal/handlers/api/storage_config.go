package api

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"catter-cms-gin/config"
	"catter-cms-gin/internal/storage"
)

var _ = storage.IsObjectStorage

// storageMu 保护动态切换存储时的并发安全
var storageMu sync.RWMutex

// maskSecret 脱敏密钥，保留前3后3字符
func maskSecret(s string) string {
	if s == "" {
		return ""
	}
	l := len(s)
	if l <= 8 {
		return strings.Repeat("*", l)
	}
	return s[:3] + strings.Repeat("*", l-6) + s[l-3:]
}

// GetStorageConfig 获取存储配置（脱敏）
func GetStorageConfig(c *gin.Context) {
	cfg := config.Get()
	s := cfg.Storage

	result := gin.H{
		"type": s.Type,
		"cdn":  s.CDN,
	}

	switch s.Type {
	case "qiniu":
		result["qiniu"] = gin.H{
			"access_key": maskSecret(s.Qiniu.AccessKey),
			"secret_key": maskSecret(s.Qiniu.SecretKey),
			"bucket":     s.Qiniu.Bucket,
			"zone":       s.Qiniu.Zone,
			"domain":     s.Qiniu.Domain,
			"use_https":  s.Qiniu.UseHTTPS,
		}
	case "alioss":
		result["alioss"] = gin.H{
			"endpoint":         s.AliOSS.Endpoint,
			"access_key_id":     maskSecret(s.AliOSS.AccessKeyID),
			"access_key_secret": maskSecret(s.AliOSS.AccessKeySecret),
			"bucket":           s.AliOSS.Bucket,
			"custom_domain":    s.AliOSS.CustomDomain,
		}
	case "tencentcos":
		result["tencentcos"] = gin.H{
			"bucket_url":    s.TencentCOS.BucketURL,
			"secret_id":     maskSecret(s.TencentCOS.SecretID),
			"secret_key":    maskSecret(s.TencentCOS.SecretKey),
			"region":        s.TencentCOS.Region,
			"custom_domain": s.TencentCOS.CustomDomain,
		}
	case "s3":
		result["s3"] = gin.H{
			"endpoint":          s.S3.Endpoint,
			"region":            s.S3.Region,
			"access_key_id":     maskSecret(s.S3.AccessKeyID),
			"secret_access_key": maskSecret(s.S3.SecretAccessKey),
			"bucket":            s.S3.Bucket,
			"force_path_style":  s.S3.ForcePathStyle,
			"custom_domain":     s.S3.CustomDomain,
		}
	}

	SuccessResponse(c, result)
}

// storageConfigRequest 前端提交的存储配置
type storageConfigRequest struct {
	Type string `json:"type"`
	CDN  string `json:"cdn"`

	Qiniu      *qiniuRequest      `json:"qiniu,omitempty"`
	AliOSS     *aliossRequest     `json:"alioss,omitempty"`
	TencentCOS *tencentcosRequest `json:"tencentcos,omitempty"`
	S3         *s3Request         `json:"s3,omitempty"`
}

type qiniuRequest struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	Zone      string `json:"zone"`
	Domain    string `json:"domain"`
	UseHTTPS  bool   `json:"use_https"`
}

type aliossRequest struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Bucket          string `json:"bucket"`
	CustomDomain    string `json:"custom_domain"`
}

type tencentcosRequest struct {
	BucketURL    string `json:"bucket_url"`
	SecretID     string `json:"secret_id"`
	SecretKey    string `json:"secret_key"`
	Region       string `json:"region"`
	CustomDomain string `json:"custom_domain"`
}

type s3Request struct {
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Bucket          string `json:"bucket"`
	ForcePathStyle  bool   `json:"force_path_style"`
	CustomDomain    string `json:"custom_domain"`
}

// configYAML 用于反序列化/序列化完整 config.yaml
type configYAML struct {
	Server   interface{} `yaml:"server,omitempty"`
	Database interface{} `yaml:"database,omitempty"`
	JWT      interface{} `yaml:"jwt,omitempty"`
	Upload   interface{} `yaml:"upload,omitempty"`
	Storage  interface{} `yaml:"storage,omitempty"`
	Security interface{} `yaml:"security,omitempty"`
	Env      string      `yaml:"env,omitempty"`
}

// SaveStorageConfig 保存存储配置并动态切换
func SaveStorageConfig(c *gin.Context) {
	var req storageConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 合并密钥：如果提交的是脱敏值（含 *），保留原值
	cfg := config.Get()
	mergeStorageConfig(&req, &cfg.Storage)

	// 验证必填字段
	if err := validateStorageConfig(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 先尝试创建 Provider 验证配置是否正确
	testCfg := buildStorageConfig(&req)
	uploadsPath := cfg.GetUploadsPath()
	testProvider, err := storage.NewProvider(testCfg, uploadsPath)
	if err != nil {
		BadRequest(c, "存储配置验证失败: "+err.Error())
		return
	}

	// 验证通过，写入 config.yaml
	if err := writeStorageToConfigYAML(&req); err != nil {
		InternalError(c, "写入配置文件失败: "+err.Error())
		return
	}

	// 动态切换存储 Provider
	storageMu.Lock()
	storageProvider = testProvider
	storage.SetProvider(testProvider)
	storageMu.Unlock()

	// 更新内存中的 config
	cfg.Storage = config.Storage{
		Type: req.Type,
		CDN:  req.CDN,
	}
	if req.Qiniu != nil {
		cfg.Storage.Qiniu = config.QiniuConfig{
			AccessKey: req.Qiniu.AccessKey,
			SecretKey: req.Qiniu.SecretKey,
			Bucket:    req.Qiniu.Bucket,
			Zone:      req.Qiniu.Zone,
			Domain:    req.Qiniu.Domain,
			UseHTTPS:  req.Qiniu.UseHTTPS,
		}
	}
	if req.AliOSS != nil {
		cfg.Storage.AliOSS = config.AliOSSConfig{
			Endpoint:        req.AliOSS.Endpoint,
			AccessKeyID:     req.AliOSS.AccessKeyID,
			AccessKeySecret: req.AliOSS.AccessKeySecret,
			Bucket:          req.AliOSS.Bucket,
			CustomDomain:    req.AliOSS.CustomDomain,
		}
	}
	if req.TencentCOS != nil {
		cfg.Storage.TencentCOS = config.TencentCOSConfig{
			BucketURL:    req.TencentCOS.BucketURL,
			SecretID:     req.TencentCOS.SecretID,
			SecretKey:    req.TencentCOS.SecretKey,
			Region:       req.TencentCOS.Region,
			CustomDomain: req.TencentCOS.CustomDomain,
		}
	}
	if req.S3 != nil {
		cfg.Storage.S3 = config.S3Config{
			Endpoint:        req.S3.Endpoint,
			Region:          req.S3.Region,
			AccessKeyID:     req.S3.AccessKeyID,
			SecretAccessKey: req.S3.SecretAccessKey,
			Bucket:          req.S3.Bucket,
			ForcePathStyle:  req.S3.ForcePathStyle,
			CustomDomain:    req.S3.CustomDomain,
		}
	}

	log.Printf("Storage switched to: %s", req.Type)

	// 返回脱敏后的配置
	GetStorageConfig(c)
}

// TestStorageConnection 测试存储连接
func TestStorageConnection(c *gin.Context) {
	var req storageConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 合并密钥
	cfg := config.Get()
	mergeStorageConfig(&req, &cfg.Storage)

	// 验证必填字段
	if err := validateStorageConfig(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 创建 Provider
	testCfg := buildStorageConfig(&req)
	uploadsPath := cfg.GetUploadsPath()
	testProvider, err := storage.NewProvider(testCfg, uploadsPath)
	if err != nil {
		BadRequest(c, "创建存储连接失败: "+err.Error())
		return
	}

	// 测试上传一个极小的文件
	testKey := fmt.Sprintf("_test/%d/connection_test.txt", time.Now().Unix())
	testContent := "Catter CMS storage connection test"
	reader := strings.NewReader(testContent)

	url, err := testProvider.Put(testKey, reader, int64(len(testContent)), "text/plain")
	if err != nil {
		BadRequest(c, "上传测试文件失败: "+err.Error())
		return
	}

	// 清理测试文件
	_ = testProvider.Delete(testKey)

	SuccessWithMessage(c, "连接测试成功", gin.H{
		"url":         url,
		"storage_type": req.Type,
	})
}

// mergeStorageConfig 合并密钥：如果提交的值包含 *（脱敏标记），则保留原值
func mergeStorageConfig(req *storageConfigRequest, current *config.Storage) {
	switch req.Type {
	case "qiniu":
		if req.Qiniu != nil {
			if strings.Contains(req.Qiniu.AccessKey, "*") {
				req.Qiniu.AccessKey = current.Qiniu.AccessKey
			}
			if strings.Contains(req.Qiniu.SecretKey, "*") {
				req.Qiniu.SecretKey = current.Qiniu.SecretKey
			}
		}
	case "alioss":
		if req.AliOSS != nil {
			if strings.Contains(req.AliOSS.AccessKeyID, "*") {
				req.AliOSS.AccessKeyID = current.AliOSS.AccessKeyID
			}
			if strings.Contains(req.AliOSS.AccessKeySecret, "*") {
				req.AliOSS.AccessKeySecret = current.AliOSS.AccessKeySecret
			}
		}
	case "tencentcos":
		if req.TencentCOS != nil {
			if strings.Contains(req.TencentCOS.SecretID, "*") {
				req.TencentCOS.SecretID = current.TencentCOS.SecretID
			}
			if strings.Contains(req.TencentCOS.SecretKey, "*") {
				req.TencentCOS.SecretKey = current.TencentCOS.SecretKey
			}
		}
	case "s3":
		if req.S3 != nil {
			if strings.Contains(req.S3.AccessKeyID, "*") {
				req.S3.AccessKeyID = current.S3.AccessKeyID
			}
			if strings.Contains(req.S3.SecretAccessKey, "*") {
				req.S3.SecretAccessKey = current.S3.SecretAccessKey
			}
		}
	}
}

// validateStorageConfig 验证存储配置的必填字段
func validateStorageConfig(req *storageConfigRequest) error {
	switch req.Type {
	case "local":
		return nil
	case "qiniu":
		if req.Qiniu == nil {
			return fmt.Errorf("请填写七牛云配置")
		}
		if req.Qiniu.AccessKey == "" || strings.Contains(req.Qiniu.AccessKey, "*") {
			return fmt.Errorf("请填写七牛云 AccessKey")
		}
		if req.Qiniu.SecretKey == "" || strings.Contains(req.Qiniu.SecretKey, "*") {
			return fmt.Errorf("请填写七牛云 SecretKey")
		}
		if req.Qiniu.Bucket == "" {
			return fmt.Errorf("请填写七牛云 Bucket")
		}
		return nil
	case "alioss":
		if req.AliOSS == nil {
			return fmt.Errorf("请填写阿里云 OSS 配置")
		}
		if req.AliOSS.Endpoint == "" {
			return fmt.Errorf("请填写 Endpoint")
		}
		if req.AliOSS.AccessKeyID == "" || strings.Contains(req.AliOSS.AccessKeyID, "*") {
			return fmt.Errorf("请填写 AccessKey ID")
		}
		if req.AliOSS.AccessKeySecret == "" || strings.Contains(req.AliOSS.AccessKeySecret, "*") {
			return fmt.Errorf("请填写 AccessKey Secret")
		}
		if req.AliOSS.Bucket == "" {
			return fmt.Errorf("请填写 Bucket")
		}
		return nil
	case "tencentcos":
		if req.TencentCOS == nil {
			return fmt.Errorf("请填写腾讯云 COS 配置")
		}
		if req.TencentCOS.BucketURL == "" {
			return fmt.Errorf("请填写 Bucket URL")
		}
		if req.TencentCOS.SecretID == "" || strings.Contains(req.TencentCOS.SecretID, "*") {
			return fmt.Errorf("请填写 Secret ID")
		}
		if req.TencentCOS.SecretKey == "" || strings.Contains(req.TencentCOS.SecretKey, "*") {
			return fmt.Errorf("请填写 Secret Key")
		}
		return nil
	case "s3":
		if req.S3 == nil {
			return fmt.Errorf("请填写 S3 配置")
		}
		if req.S3.Endpoint == "" {
			return fmt.Errorf("请填写 Endpoint")
		}
		if req.S3.Region == "" {
			return fmt.Errorf("请填写 Region")
		}
		if req.S3.AccessKeyID == "" || strings.Contains(req.S3.AccessKeyID, "*") {
			return fmt.Errorf("请填写 Access Key ID")
		}
		if req.S3.SecretAccessKey == "" || strings.Contains(req.S3.SecretAccessKey, "*") {
			return fmt.Errorf("请填写 Secret Access Key")
		}
		if req.S3.Bucket == "" {
			return fmt.Errorf("请填写 Bucket")
		}
		return nil
	default:
		return fmt.Errorf("不支持的存储类型: %s", req.Type)
	}
}

// buildStorageConfig 将请求转换为 storage.StorageConfig
func buildStorageConfig(req *storageConfigRequest) storage.StorageConfig {
	sc := storage.StorageConfig{
		Type: req.Type,
		CDN:  req.CDN,
	}
	if req.Qiniu != nil {
		sc.Qiniu = storage.QiniuConfig{
			AccessKey: req.Qiniu.AccessKey,
			SecretKey: req.Qiniu.SecretKey,
			Bucket:    req.Qiniu.Bucket,
			Zone:      req.Qiniu.Zone,
			Domain:    req.Qiniu.Domain,
			UseHTTPS:  req.Qiniu.UseHTTPS,
		}
	}
	if req.AliOSS != nil {
		sc.AliOSS = storage.AliOSSConfig{
			Endpoint:        req.AliOSS.Endpoint,
			AccessKeyID:     req.AliOSS.AccessKeyID,
			AccessKeySecret: req.AliOSS.AccessKeySecret,
			Bucket:          req.AliOSS.Bucket,
			CustomDomain:    req.AliOSS.CustomDomain,
		}
	}
	if req.TencentCOS != nil {
		sc.TencentCOS = storage.TencentCOSConfig{
			BucketURL:    req.TencentCOS.BucketURL,
			SecretID:     req.TencentCOS.SecretID,
			SecretKey:    req.TencentCOS.SecretKey,
			Region:       req.TencentCOS.Region,
			CustomDomain: req.TencentCOS.CustomDomain,
		}
	}
	if req.S3 != nil {
		sc.S3 = storage.S3Config{
			Endpoint:        req.S3.Endpoint,
			Region:          req.S3.Region,
			AccessKeyID:     req.S3.AccessKeyID,
			SecretAccessKey: req.S3.SecretAccessKey,
			Bucket:          req.S3.Bucket,
			ForcePathStyle:  req.S3.ForcePathStyle,
			CustomDomain:    req.S3.CustomDomain,
		}
	}
	return sc
}

// writeStorageToConfigYAML 将存储配置写入 config.yaml（只修改 storage 段）
func writeStorageToConfigYAML(req *storageConfigRequest) error {
	// 找到 config.yaml 路径
	configPath := findConfigPath()
	if configPath == "" {
		configPath = "config.yaml"
	}

	// 读取现有 YAML
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 反序列化为 map 保留所有字段
	var cfgMap map[string]interface{}
	if err := yaml.Unmarshal(data, &cfgMap); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 构建 storage 段
	storageMap := map[string]interface{}{
		"type": req.Type,
	}
	if req.CDN != "" {
		storageMap["cdn"] = req.CDN
	}

	switch req.Type {
	case "local":
		// 本地存储只需 type
	case "qiniu":
		if req.Qiniu != nil {
			storageMap["qiniu"] = map[string]interface{}{
				"access_key": req.Qiniu.AccessKey,
				"secret_key": req.Qiniu.SecretKey,
				"bucket":     req.Qiniu.Bucket,
				"zone":       req.Qiniu.Zone,
				"domain":     req.Qiniu.Domain,
				"use_https":  req.Qiniu.UseHTTPS,
			}
		}
	case "alioss":
		if req.AliOSS != nil {
			storageMap["alioss"] = map[string]interface{}{
				"endpoint":          req.AliOSS.Endpoint,
				"access_key_id":     req.AliOSS.AccessKeyID,
				"access_key_secret": req.AliOSS.AccessKeySecret,
				"bucket":            req.AliOSS.Bucket,
				"custom_domain":     req.AliOSS.CustomDomain,
			}
		}
	case "tencentcos":
		if req.TencentCOS != nil {
			storageMap["tencentcos"] = map[string]interface{}{
				"bucket_url":    req.TencentCOS.BucketURL,
				"secret_id":     req.TencentCOS.SecretID,
				"secret_key":    req.TencentCOS.SecretKey,
				"region":        req.TencentCOS.Region,
				"custom_domain": req.TencentCOS.CustomDomain,
			}
		}
	case "s3":
		if req.S3 != nil {
			storageMap["s3"] = map[string]interface{}{
				"endpoint":          req.S3.Endpoint,
				"region":            req.S3.Region,
				"access_key_id":     req.S3.AccessKeyID,
				"secret_access_key": req.S3.SecretAccessKey,
				"bucket":            req.S3.Bucket,
				"force_path_style":  req.S3.ForcePathStyle,
				"custom_domain":     req.S3.CustomDomain,
			}
		}
	}

	// 更新 storage 段
	cfgMap["storage"] = storageMap

	// 重新序列化
	out, err := yaml.Marshal(cfgMap)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 写回文件
	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// findConfigPath 查找配置文件路径
func findConfigPath() string {
	paths := []string{
		"./config.yaml",
		"../config.yaml",
		"../../config.yaml",
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}


