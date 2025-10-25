package server

import (
	"context"
	"httpServerTest/internal/config"
	"httpServerTest/internal/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	srv *http.Server
}

func NewHTTPServer(cfg *config.Config) *HTTPServer {
	router := gin.New()
	router.Use(gin.Recovery())

	// 路由
	handler.RegisterUserRoute(cfg, router)
	handler.RegisterBookRoute(cfg, router)
	handler.RegisterChapterRoute(cfg, router)
	handler.RegisterReadingRecordRoute(cfg, router)

	srv := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: router,
	}
	return &HTTPServer{srv: srv}
}

func (s *HTTPServer) Start() error {
	return s.srv.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
