package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/time/rate"
)

// IPRateLimiter IP rate limiter
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
	// lastSeen 记录每个 key 最近一次访问时间，用于淘汰长期不活跃的条目
	lastSeen map[string]time.Time
}

// evictInterval 淘汰周期与不活跃阈值：超过该时间未访问的条目会被清理
const evictInterval = time.Hour

// NewIPRateLimiter creates a new IP rate limiter
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips:      make(map[string]*rate.Limiter),
		mu:       &sync.RWMutex{},
		r:        r,
		b:        b,
		lastSeen: make(map[string]time.Time),
	}
	// 后台定时清理，防止恶意 IP 遍历导致 map 无限增长
	go func() {
		for range time.Tick(evictInterval) {
			i.evictStale()
		}
	}()
	return i
}

// evictStale 清理长期不活跃的条目
func (i *IPRateLimiter) evictStale() {
	i.mu.Lock()
	defer i.mu.Unlock()
	cutoff := time.Now().Add(-evictInterval)
	for ip, seen := range i.lastSeen {
		if seen.Before(cutoff) {
			delete(i.ips, ip)
			delete(i.lastSeen, ip)
		}
	}
}

// AddIP creates a new rate limiter and adds it to the map
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter
	i.lastSeen[ip] = time.Now()
	return limiter
}

// GetLimiter retrieves the rate limiter from the map
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	if limiter, exists := i.ips[ip]; exists {
		i.lastSeen[ip] = time.Now()
		return limiter
	}
	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter
	i.lastSeen[ip] = time.Now()
	return limiter
}

// LoginRateLimiter creates a login rate limit middleware.
// Limits: 10 requests per minute per IP, 5 requests per minute per email.
func LoginRateLimiter() gin.HandlerFunc {
	ipLimiter := NewIPRateLimiter(rate.Every(time.Minute), 10)
	userLimiter := NewIPRateLimiter(rate.Every(time.Minute), 5)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !ipLimiter.GetLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "当前 IP 登录尝试过于频繁，请1分钟后再试 (限制: 10次/分钟)"})
			c.Abort()
			return
		}

		var req struct {
			Email string `json:"email"`
		}
		// ShouldBindBodyWith 缓存 body，以便后续的 handler 还能继续读取
		if err := c.ShouldBindBodyWith(&req, binding.JSON); err == nil && req.Email != "" {
			if !userLimiter.GetLimiter(req.Email).Allow() {
				c.JSON(http.StatusTooManyRequests, gin.H{"message": "该账号登录尝试过于频繁，请1分钟后再试 (限制: 5次/分钟)"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
