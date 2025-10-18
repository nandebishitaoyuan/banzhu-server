package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"httpServerTest/internal/config"
	"httpServerTest/internal/database"
	"httpServerTest/internal/model"
	"httpServerTest/internal/server"
	"httpServerTest/pkg/logger"
)

func main() {
	l := logger.New()

	// 加载配置
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		l.Fatalf("配置加载失败: %v", err)
	}

	// 初始化数据库
	if err := database.Init(cfg); err != nil {
		l.Fatalf("数据库连接失败: %v", err)
	}
	if err := model.Migrate(database.DB); err != nil {
		l.Fatalf("数据库迁移失败: %v", err)
	}

	// 启动 HTTP 服务
	srv := server.NewHTTPServer(cfg)
	go func() {
		if err := srv.Start(); err != nil {
			l.Fatalf("服务器启动失败: %v", err)
		}
	}()
	l.Infof("✅ 服务启动成功，监听端口: %s", cfg.Server.Addr)

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info("🛑 收到退出信号，准备关闭服务器...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		l.Fatalf("关闭服务器失败: %v", err)
	}
	l.Info("✅ 服务器关闭。")
}
