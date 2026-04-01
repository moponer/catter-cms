package models

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"gorm.io/gorm"
)

// AIToken AI 访问令牌模型（catter-cms-key）
type AIToken struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;size:255"`
	Token       string    `json:"-" gorm:"uniqueIndex;not null;size:255"` // cat_xxx，json:"-" 不输出完整 token
	TokenPrefix string    `json:"token_prefix" gorm:"size:10"`            // cat_abc1... 显示用前缀
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	User        *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Permissions string    `json:"permissions" gorm:"type:text;not null;default:'[\"read\"]'"` // JSON: ["read","write","admin"]
	Scopes      string    `json:"scopes" gorm:"type:text;not null;default:'[\"posts\",\"pages\",\"categories\",\"tags\",\"media\",\"comments\",\"settings\"]'"` // JSON: 资源范围
	LastUsedAt  *time.Time `json:"last_used_at"`
	ExpiresAt   time.Time  `json:"expires_at" gorm:"not null;default:'9999-12-31 23:59:59'"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (AIToken) TableName() string {
	return "ai_tokens"
}

// BeforeCreate 创建前钩子 - 生成 token 和前缀
func (t *AIToken) BeforeCreate(tx *gorm.DB) error {
	if t.Token == "" {
		token, err := GenerateAIToken()
		if err != nil {
			return err
		}
		t.Token = token
	}
	if t.TokenPrefix == "" {
		// 显示前8个字符 + ...
		if len(t.Token) > 8 {
			t.TokenPrefix = t.Token[:8] + "..."
		} else {
			t.TokenPrefix = t.Token
		}
	}
	if t.ExpiresAt.IsZero() {
		t.ExpiresAt = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	}
	return nil
}

// IsExpired 检查是否已过期
func (t *AIToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// HasPermission 检查是否拥有指定权限
func (t *AIToken) HasPermission(perm string) bool {
	perms := t.GetPermissions()
	for _, p := range perms {
		if p == perm || p == "admin" {
			return true
		}
	}
	return false
}

// HasScope 检查是否拥有指定资源范围
func (t *AIToken) HasScope(scope string) bool {
	scopes := t.GetScopes()
	for _, s := range scopes {
		if s == scope || s == "all" {
			return true
		}
	}
	return false
}

// GetPermissions 解析权限列表
func (t *AIToken) GetPermissions() []string {
	return parseJSONStringArray(t.Permissions)
}

// GetScopes 解析资源范围列表
func (t *AIToken) GetScopes() []string {
	return parseJSONStringArray(t.Scopes)
}

// AIAuditLog AI 审计日志
type AIAuditLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TokenID   uint      `json:"token_id" gorm:"index"`
	Token     *AIToken  `json:"token,omitempty" gorm:"foreignKey:TokenID"`
	Action    string    `json:"action" gorm:"not null;size:255"`  // GET, POST, PUT, DELETE
	Resource  string    `json:"resource" gorm:"size:255"`          // posts, pages, categories...
	Path      string    `json:"path" gorm:"size:1000"`             // 请求路径
	Status    int       `json:"status"`                            // HTTP 状态码
	IPAddress string    `json:"ip_address" gorm:"size:100"`
	UserAgent string    `json:"user_agent" gorm:"size:500"`
	RequestID string    `json:"request_id" gorm:"size:100"`        // 关联请求 ID
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (AIAuditLog) TableName() string {
	return "ai_audit_logs"
}

// GenerateAIToken 生成 AI 令牌（cat_ 前缀 + 32 字节随机 hex）
func GenerateAIToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "cat_" + hex.EncodeToString(bytes), nil
}

// parseJSONStringArray 简易 JSON 字符串数组解析（避免引入额外依赖）
func parseJSONStringArray(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return []string{}
	}
	// 移除 [ ]
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	s = strings.TrimSpace(s)
	if s == "" {
		return []string{}
	}
	// 按 , 分割，清理引号
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, "\"")
		p = strings.Trim(p, "'")
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
