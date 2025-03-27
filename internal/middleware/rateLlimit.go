package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

// Limiter 限流器接口
type Limiter interface {
	Allow(c *gin.Context) bool
}

// RateLimiterConfig 限流器配置
type RateLimiterConfig struct {
	Rate        rate.Limit    // 每秒允许的请求数
	Burst       int           // 允许的突发请求数
	Redis       *redis.Client // Redis客户端实例
	KeyPrefix   string        // Redis键前缀
	LimiterType string        // 限流器类型：ip/global
}

// ipRateLimiter 基于IP的限流器
type ipRateLimiter struct {
	sync.Mutex
	limiters map[string]*rate.Limiter
	config   RateLimiterConfig
}

// globalRateLimiter 全局限流器
type globalRateLimiter struct {
	limiter *rate.Limiter
	config  RateLimiterConfig
}

// redisRateLimiter Redis分布式限流器
type redisRateLimiter struct {
	config RateLimiterConfig
}

func NewRateLimiter(config RateLimiterConfig) Limiter {
	if config.Redis != nil {
		return &redisRateLimiter{config: config}
	}

	switch config.LimiterType {
	case "ip":
		return &ipRateLimiter{
			limiters: make(map[string]*rate.Limiter),
			config:   config,
		}
	default: // 全局限流
		return &globalRateLimiter{
			limiter: rate.NewLimiter(config.Rate, config.Burst),
			config:  config,
		}
	}
}

// Allow 实现Limiter接口
func (i *ipRateLimiter) Allow(c *gin.Context) bool {
	ip := c.ClientIP()
	i.Lock()
	defer i.Unlock()

	limiter, exists := i.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(i.config.Rate, i.config.Burst)
		i.limiters[ip] = limiter
	}

	return limiter.Allow()
}

func (g *globalRateLimiter) Allow(c *gin.Context) bool {
	return g.limiter.Allow()
}

func (r *redisRateLimiter) Allow(c *gin.Context) bool {
	key := r.config.KeyPrefix
	if r.config.LimiterType == "ip" {
		key += c.ClientIP()
	}

	// 使用Redis事务实现原子操作
	script := redis.NewScript(`
		local current = redis.call('incr', KEYS[1])
		if tonumber(current) == 1 then
			redis.call('expire', KEYS[1], 1)
		end
		return current
	`)

	result, err := script.Run(c, r.config.Redis, []string{key}).Int()
	if err != nil {
		return true // Redis不可用时放行请求
	}

	return result <= r.config.Burst
}

// RateLimitMiddleware 生成限流中间件
func RateLimitMiddleware(config RateLimiterConfig) gin.HandlerFunc {
	limiter := NewRateLimiter(config)

	return func(c *gin.Context) {
		if !limiter.Allow(c) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			return
		}
		c.Next()
	}
}
