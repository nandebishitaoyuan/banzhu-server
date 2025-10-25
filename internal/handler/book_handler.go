package handler

import (
	"httpServerTest/internal/config"
	"httpServerTest/internal/middleware"
	"httpServerTest/internal/model"
	"httpServerTest/internal/service"
	"httpServerTest/pkg/stringUtil"
	"net/http"
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
	bookRouter.GET("/syncById", bookHandler.SyncBookById)
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
	// 将字符串转换为uint
	id64, err := stringUtil.StringToUint64(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 调用服务
	err = h.svc.DeleteBook(id64)
	// 删除章节
	err = service.ChapterService{}.DeleteChapterByBookId(id64)
	user := middleware.GetCurrentUser(c)
	// 删除阅读记录
	err = service.ReadingRecordService{}.DeleteByBookId(id64, user.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, "删除书籍失败！"+err.Error())
	} else {
		c.JSON(http.StatusOK, nil)
	}
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

func (h *BookHandler) SyncBookById(context *gin.Context) {
	id := context.Query("id")
	// 将字符串转换为uint
	id64, err := stringUtil.StringToUint64(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = h.svc.SyncBookById(id64)
	if err != nil {
		context.JSON(http.StatusBadRequest, "同步指定书籍失败！"+err.Error())
		return
	}
	user := middleware.GetCurrentUser(context)
	// 删除阅读记录
	err = service.ReadingRecordService{}.DeleteByBookId(id64, user.ID)
	context.JSON(http.StatusOK, "同步成功")
}
