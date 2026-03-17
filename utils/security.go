package utils

import (
	"errors"
	"regexp"
	"strings"
)

// ValidatePasswordStrength 验证密码强度
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度至少8位")
	}
	
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)
	
	strength := 0
	if hasUpper {
		strength++
	}
	if hasLower {
		strength++
	}
	if hasNumber {
		strength++
	}
	if hasSpecial {
		strength++
	}
	
	if strength < 3 {
		return errors.New("密码强度不足，请包含大小写字母、数字和特殊字符")
	}
	
	return nil
}

// SanitizeInput 清理输入
func SanitizeInput(input string) string {
	dangerousTags := []string{
		"<script>", "</script>",
		"<iframe>", "</iframe>",
		"<object>", "</object>",
		"<embed>", "</embed>",
		"<form>", "</form>",
		"<input>", "</input>",
		"<style>", "</style>",
		"<link>", "</link>",
		"<meta>", "</meta>",
	}
	
	result := input
	for _, tag := range dangerousTags {
		result = strings.ReplaceAll(result, tag, "")
	}
	return result
}

// MaskEmail 脱敏邮箱
func MaskEmail(email string) string {
	if len(email) < 3 {
		return email
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	username := parts[0]
	if len(username) <= 2 {
		return username
	}
	return username[:2] + "***@" + parts[1]
}

// MaskPassword 脱敏密码
func MaskPassword(password string) string {
	if len(password) < 4 {
		return password
	}
	return password[:2] + "******"
}

// MaskToken 脱敏token
func MaskToken(token string) string {
	if len(token) < 8 {
		return token
	}
	return token[:4] + "..." + token[len(token)-4:]
}
