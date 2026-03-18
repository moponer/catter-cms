package api

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// SystemInfo 系统信息结构
type SystemInfo struct {
	// Catter CMS 版本
	Version string `json:"version"`
	// 服务器信息
	Server ServerInfo `json:"server"`
	// 数据库信息
	Database DatabaseInfo `json:"database"`
	// 内存信息
	Memory MemoryInfo `json:"memory"`
	// CPU信息
	CPU CPUInfo `json:"cpu"`
	// 运行时间
	Uptime string `json:"uptime"`
	// 启动时间
	StartTime string `json:"start_time"`
}

// ServerInfo 服务器信息
type ServerInfo struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Hostname string `json:"hostname"`
}

// DatabaseInfo 数据库信息
type DatabaseInfo struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

// CPUInfo CPU信息
type CPUInfo struct {
	Cores     int     `json:"cores"`
	Goroutines int    `json:"goroutines"`
}

var startTime = time.Now()

// GetSystemInfo 获取系统信息
func GetSystemInfo(c *gin.Context) {
	// 获取内存统计
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 计算运行时间
	uptime := time.Since(startTime)
	uptimeStr := formatDuration(uptime)

	info := SystemInfo{
		Version:   "2.0.0",
		StartTime: startTime.Format("2006-01-02 15:04:05"),
		Uptime:    uptimeStr,
		Server: ServerInfo{
			OS:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			Hostname: "localhost",
		},
		Database: DatabaseInfo{
			Type:    "SQLite",
			Version: "3.x",
			Status:  "connected",
		},
		Memory: MemoryInfo{
			Total:       m.Sys,
			Used:        m.Alloc,
			Free:        m.Sys - m.Alloc,
			UsedPercent: float64(m.Alloc) / float64(m.Sys) * 100,
		},
		CPU: CPUInfo{
			Cores:      runtime.NumCPU(),
			Goroutines: runtime.NumGoroutine(),
		},
	}

	SuccessResponse(c, info)
}

// formatDuration 格式化持续时间
func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return formatInt(days) + "天" + formatInt(hours) + "小时"
	}
	if hours > 0 {
		return formatInt(hours) + "小时" + formatInt(minutes) + "分钟"
	}
	return formatInt(minutes) + "分钟"
}

func formatInt(n int) string {
	if n < 10 {
		return fmt.Sprintf("0%d", n)
	}
	return fmt.Sprintf("%d", n)
}
