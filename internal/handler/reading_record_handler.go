package handler

import (
	"httpServerTest/internal/config"
	"httpServerTest/internal/middleware"
	"httpServerTest/internal/model"
	"httpServerTest/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterReadingRecordRoute(cfg *config.Config, router *gin.Engine) {
	readingRecordHandler := NewReadingRecordHandler()
	readingRecordRouter := router.Group("/readingRecord")
	readingRecordRouter.Use(middleware.AuthMiddleware(cfg))

	readingRecordRouter.POST("/update", middleware.DebounceMiddleware(), readingRecordHandler.Update)
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
