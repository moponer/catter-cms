package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage 本地文件存储
type LocalStorage struct {
	basePath string // 上传根目录（如 "uploads"）
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{basePath: basePath}
}

// Put 保存文件到本地目录，返回 /uploads/xxx 格式的 URL
func (s *LocalStorage) Put(objectKey string, reader io.Reader, size int64, contentType string) (string, error) {
	// 构建完整文件路径
	fullPath := filepath.Join(s.basePath, objectKey)

	// 确保目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建文件
	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer f.Close()

	// 写入内容
	written, err := io.Copy(f, reader)
	if err != nil {
		os.Remove(fullPath) // 清理失败文件
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	_ = size    // 本地存储不强制校验 size
	_ = written // 已成功写入

	return "/uploads/" + objectKey, nil
}

// Delete 删除本地文件
func (s *LocalStorage) Delete(objectKey string) error {
	fullPath := filepath.Join(s.basePath, objectKey)
	if _, err := os.Stat(fullPath); err == nil {
		return os.Remove(fullPath)
	}
	return nil // 文件不存在视为成功
}

// GetURL 获取本地文件的 URL
func (s *LocalStorage) GetURL(objectKey string) string {
	return "/uploads/" + objectKey
}

// GetThumbURL 获取缩略图 URL（本地存储使用 .thumbs 目录）
func (s *LocalStorage) GetThumbURL(objectKey string) string {
	if objectKey == "" {
		return ""
	}
	ext := filepath.Ext(objectKey)
	base := strings.TrimSuffix(objectKey, ext)
	return "/uploads/.thumbs/thumb_" + base + ".jpg"
}

// Type 返回存储类型
func (s *LocalStorage) Type() string {
	return "local"
}
