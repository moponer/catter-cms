package validators

import (
	"regexp"
	"strings"
)

// Validator 验证器接口
type Validator interface {
	Validate(input string) (bool, string)
}

// EmailValidator 邮箱验证器
type EmailValidator struct {
	pattern *regexp.Regexp
}

// NewEmailValidator 创建新的邮箱验证器
func NewEmailValidator() *EmailValidator {
	return &EmailValidator{
		pattern: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

// Validate 验证邮箱
func (v *EmailValidator) Validate(email string) (bool, string) {
	if email == "" {
		return false, "邮箱不能为空"
	}
	return v.pattern.MatchString(email), "邮箱格式正确"
}

// PasswordValidator 密码验证器
type PasswordValidator struct {
	minLength int
	requireUppercase bool
	requireLowercase bool
	requireNumber bool
	requireSpecial bool
}

// NewPasswordValidator 创建新的密码验证器
func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		minLength:      8,
		requireUppercase: true,
		requireLowercase: true,
		requireNumber:     true,
		requireSpecial:   true,
	}
}

// Validate 验证密码
func (v *PasswordValidator) Validate(password string) (bool, string) {
	if len(password) < v.minLength {
		return false, "密码长度至少8位"
	}
	
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)
	
	if v.requireUppercase && !hasUpper {
		return false, "密码必须包含大写字母"
	}
	
	if v.requireLowercase && !hasLower {
		return false, "密码必须包含小写字母"
	}
	
	if v.requireNumber && !hasNumber {
		return false, "密码必须包含数字"
	}
	
	if v.requireSpecial && !hasSpecial {
		return false, "密码必须包含特殊字符"
	}
	
	return true, "密码强度符合要求"
}

// XSSSanitizer XSS防护器
type XSSSanitizer struct {
	dangerousTags []string
}

// NewXSSSanitizer 创建新的XSS防护器
func NewXSSSanitizer() *XSSSanitizer {
	return &XSSSanitizer{
		dangerousTags: []string{
			"<script>", "</script>",
			"<iframe>", "</iframe>",
			"<object>", "</object>",
			"<embed>", "</embed>",
			"<form>", "</form>",
			"<input>", "</input>",
			"<style>", "</style>",
			"<link>", "</link>",
			"<meta>", "</meta>",
		},
	}
}

// Sanitize 清理输入
func (s *XSSSanitizer) Sanitize(input string) string {
	result := input
	for _, tag := range s.dangerousTags {
		result = strings.ReplaceAll(result, tag, "")
	}
	return result
}

// ValidateInput 验证输入
func ValidateInput(input string, validators ...Validator) (bool, string) {
	for _, validator := range validators {
		if ok, message := validator.Validate(input); !ok {
			return false, message
		}
	}
	return true, "输入验证通过"
}

// ValidateEmail 验证邮箱
func ValidateEmail(email string) (bool, string) {
	validator := NewEmailValidator()
	return validator.Validate(email)
}

// ValidatePassword 验证密码
func ValidatePassword(password string) (bool, string) {
	validator := NewPasswordValidator()
	return validator.Validate(password)
}

// SanitizeXSS 清理XSS
func SanitizeXSS(input string) string {
	sanitizer := NewXSSSanitizer()
	return sanitizer.Sanitize(input)
}
