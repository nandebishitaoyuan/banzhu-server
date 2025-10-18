package handler

import (
	"httpServerTest/internal/config"
	"httpServerTest/internal/middleware"
	"httpServerTest/internal/service"
	"httpServerTest/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoute(cfg *config.Config, router *gin.Engine) {
	userHandler := NewUserHandler(cfg)

	router.POST("/user/register", userHandler.Register)
	router.POST("/user/login", userHandler.Login)
	router.GET("/user/refreshToken", userHandler.RefreshToken)

	auth := router.Group("/user")
	auth.Use(middleware.AuthMiddleware(cfg))

	auth.GET("/profile", userHandler.Profile)
}

type UserHandler struct {
	svc *service.UserService
	cfg *config.Config
}

func NewUserHandler(cfg *config.Config) *UserHandler {
	return &UserHandler{
		svc: service.NewUserService(),
		cfg: cfg,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.Register(req.Username, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "注册成功")
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "参数错误")
		return
	}
	user, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	token, err := jwt.GenerateAccessToken(user.ID, h.cfg)
	refreshToken, err := jwt.GenerateRefreshToken(user.ID, h.cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "生成 token 失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  token,
		"refreshToken": refreshToken,
		"userInfo":     user,
	})
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.Query("refreshToken")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, "刷新令牌不能为空")
		return
	}

	claims, err := jwt.ParseRefreshToken(refreshToken, h.cfg)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "刷新令牌无效或过期")
		return
	}

	newAccess, _ := jwt.GenerateAccessToken(claims.UserID, h.cfg)
	newRefresh, _ := jwt.GenerateRefreshToken(claims.UserID, h.cfg)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccess,
		"refresh_token": newRefresh,
	})
}

func (h *UserHandler) Profile(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	c.JSON(http.StatusOK, user)
}
