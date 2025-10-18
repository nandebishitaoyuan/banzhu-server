package middleware

import (
	"httpServerTest/internal/config"
	"httpServerTest/internal/database"
	"httpServerTest/internal/model"
	"httpServerTest/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少 Authorization 头"})
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.ParseAccessToken(tokenStr, cfg)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}

		var user model.User
		if err := database.DB.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		// 保存到上下文，后续 handler 可以直接获取当前用户
		c.Set("user", user)
		c.Next()
	}
}

// GetCurrentUser 获取当前登录用户
func GetCurrentUser(c *gin.Context) model.User {
	val, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		c.Abort()
		return model.User{}
	}
	user, ok := val.(model.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		c.Abort()
		return model.User{}
	}
	return user
}
