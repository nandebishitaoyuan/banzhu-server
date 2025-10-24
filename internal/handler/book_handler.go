package handler

import (
	"httpServerTest/internal/config"
	"httpServerTest/internal/middleware"
	"httpServerTest/internal/model"
	"httpServerTest/internal/service"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

func RegisterBookRoute(cfg *config.Config, router *gin.Engine) {
	bookHandler := NewBookHandler(cfg)

	router.GET("/book/sync", bookHandler.SyncBook)

	bookRouter := router.Group("/book")
	bookRouter.Use(middleware.AuthMiddleware(cfg))

	bookRouter.POST("/page", bookHandler.GetPage)
	bookRouter.GET("/delete", bookHandler.DeleteBook)
	bookHandler.SyncBook(nil)
}

type BookHandler struct {
	svc *service.BookService
	cfg *config.Config
	mu  sync.Mutex
}

func NewBookHandler(cfg *config.Config) *BookHandler {
	return &BookHandler{
		svc: service.NewBookService(cfg),
		cfg: cfg,
	}
}

func (h *BookHandler) GetPage(c *gin.Context) {
	var param model.KeywordPageParam[string]
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(400, err.Error())
		return
	}
	page, err := h.svc.GetPage(param)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, page)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "id不能为空")
		return
	}
	// 将字符串转换为uint
	id64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "id必须是数字"+err.Error())
		return
	}

	// 调用服务
	if err = h.svc.DeleteBook(id64); err != nil {
		c.JSON(http.StatusBadRequest, "删除书籍失败！"+err.Error())
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, "删除章节失败！"+err.Error())
	}

	c.JSON(http.StatusOK, nil)
}

// SyncBook 同步本地已存在的书籍到数据库
func (h *BookHandler) SyncBook(c *gin.Context) {
	h.mu.Lock()         // 加锁
	defer h.mu.Unlock() // 确保出错也能解锁

	err := h.svc.SyncBook()
	if err != nil {
		if c != nil {
			c.JSON(500, err.Error())
		}
		println(err.Error())
		return
	}

	if c != nil {
		c.JSON(200, "同步成功")
	}
}
