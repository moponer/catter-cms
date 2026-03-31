package api

import (
	"strings"
)

// escapeLikeWildcards 转义 LIKE 查询中的通配符，防止 SQL 注入
// 转义顺序：反斜杠 -> 百分号 -> 下划线
func escapeLikeWildcards(s string) string {
	// 转义反斜杠（必须最先处理）
	s = strings.ReplaceAll(s, "\\", "\\\\")
	// 转义百分号
	s = strings.ReplaceAll(s, "%", "\\%")
	// 转义下划线
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}
