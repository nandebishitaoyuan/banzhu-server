package handler

import (
	"httpServerTest/internal/config"
	"httpServerTest/internal/middleware"
	"httpServerTest/internal/model"
	"httpServerTest/internal/service"
	"httpServerTest/pkg/stringUtil"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func RegisterChapterRoute(cfg *config.Config, router *gin.Engine) {
	chapterHandler := NewChapterHandler(cfg)
	bookRouter := router.Group("/chapter")
	bookRouter.Use(middleware.AuthMiddleware(cfg))

	bookRouter.POST("/page", chapterHandler.GetPage)
	bookRouter.GET("/list", chapterHandler.GetList)
	bookRouter.GET("/delete", chapterHandler.Delete)
	bookRouter.GET("/downloadChapter", chapterHandler.DownloadChapter)
	bookRouter.GET("/downloadImage", chapterHandler.DownloadImage)
}

type ChapterHandler struct {
	svc *service.ChapterService
	cfg *config.Config
}

func NewChapterHandler(cfg *config.Config) *ChapterHandler {
	return &ChapterHandler{
		svc: service.NewChapterService(cfg),
		cfg: cfg,
	}
}

func (h *ChapterHandler) GetPage(c *gin.Context) {
	var param model.KeywordPageParam[uint64]
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.svc.GetPage(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, res)
}

func (h *ChapterHandler) GetList(c *gin.Context) {
	bookId := c.Query("bookId")

	// 将字符串转换为uint
	id64, err := stringUtil.StringToUint64(bookId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.svc.GetList(id64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, res)
}

func (h *ChapterHandler) Delete(c *gin.Context) {
	id := c.Query("id")

	// 将字符串转换为uint
	id64, err := stringUtil.StringToUint64(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 调用服务
	if h.svc.DeleteChapter(id64) != nil {
		c.JSON(http.StatusBadRequest, "删除失败！")
	}

	c.JSON(http.StatusOK, nil)
}

func (h *ChapterHandler) DownloadChapter(c *gin.Context) {
	id := c.Query("id")

	// 将字符串转换为uint
	id64, err := stringUtil.StringToUint64(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	chapter, err := h.svc.GetChapterById(id64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(*chapter.Path))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("media-type", "text/plain")
	c.File(*chapter.Path)
}

func (h *ChapterHandler) DownloadImage(c *gin.Context) {
	id := c.Query("id")

	path := "/data/images/" + id + ".png"

	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(path))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("media-type", "image/png")
	c.File(path)
}
