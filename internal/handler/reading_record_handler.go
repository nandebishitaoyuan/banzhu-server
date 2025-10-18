package handler

import (
	"httpServerTest/internal/config"
	"httpServerTest/internal/middleware"
	"httpServerTest/internal/model"
	"httpServerTest/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterReadingRecordRoute(cfg *config.Config, router *gin.Engine) {
	readingRecordHandler := NewReadingRecordHandler()
	readingRecordRouter := router.Group("/readingRecord")
	readingRecordRouter.Use(middleware.AuthMiddleware(cfg))

	readingRecordRouter.POST("/update", middleware.DebounceMiddleware(), readingRecordHandler.Update)
	readingRecordRouter.GET("/getListByBookId", readingRecordHandler.GetListByBookId)
	readingRecordRouter.POST("/get", readingRecordHandler.Get)
}

type ReadingRecordHandler struct {
	svc *service.ReadingRecordService
}

func NewReadingRecordHandler() *ReadingRecordHandler {
	return &ReadingRecordHandler{
		svc: service.NewReadingRecordService(),
	}
}

func (h *ReadingRecordHandler) Update(context *gin.Context) {
	var param model.ReadingRecordDataParam
	if err := context.ShouldBindJSON(&param); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := middleware.GetCurrentUser(context)
	if err := h.svc.Update(param, user.ID); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}
}

func (h *ReadingRecordHandler) Get(context *gin.Context) {
	var param model.ReadingRecordQueryParam
	if err := context.ShouldBindJSON(&param); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := middleware.GetCurrentUser(context)
	res, err := h.svc.Get(param, user.ID)
	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}
	context.JSON(http.StatusOK, res)
}

func (h *ReadingRecordHandler) GetListByBookId(context *gin.Context) {
	id := context.Query("bookId")
	if id == "" {
		context.JSON(http.StatusBadRequest, "id不能为空")
		return
	}
	// 将字符串转换为uint
	id64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, "id必须是数字"+err.Error())
		return
	}

	user := middleware.GetCurrentUser(context)
	res, err := h.svc.GetListByBookId(id64, user.ID)
	if err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}
	context.JSON(http.StatusOK, res)
}
