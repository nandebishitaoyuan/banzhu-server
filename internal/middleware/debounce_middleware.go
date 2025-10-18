package middleware

import (
	"sync"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// 用于存储每个用户或请求标识符的时间戳
	lastRequestTimes = make(map[string]time.Time)
	mu               sync.Mutex
	debounceDuration = 2 * time.Second // 防抖间隔时间
)

// DebounceMiddleware 防抖中间件
func DebounceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户标识符，可以是 IP 或用户 ID
		identifier := c.ClientIP() // 如果是基于用户 ID，可以换成其他标识符

		// 获取当前时间
		currentTime := time.Now()

		mu.Lock()
		defer mu.Unlock()

		// 检查用户最近一次请求的时间
		lastRequestTime, exists := lastRequestTimes[identifier]
		if exists && currentTime.Sub(lastRequestTime) < debounceDuration {
			// 如果请求时间过于接近上次请求时间，拒绝处理请求
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "请求太频繁了，请稍后再试！"})
			c.Abort()
			return
		}

		// 更新用户最后一次请求时间
		lastRequestTimes[identifier] = currentTime

		// 继续处理请求
		c.Next()
	}
}
