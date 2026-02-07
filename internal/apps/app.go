package apps

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fiber/config"
	"fiber/pkg/logx"

	"github.com/gofiber/fiber/v3"
)

func StartAPIServer(cfg *config.Config) {
	app := fiber.New(fiber.Config{
		ReadTimeout: 10 * time.Second,
	})

	// 1. 创建一个通道来监听系统信号
	// 需要监听：SIGINT (Ctrl+C), SIGTERM (由 Kubernetes/Docker 发出的停止信号)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// 2. 在另一个 goroutine 中启动服务器，防止阻塞主线程
	go func() {
		if err := app.Listen(cfg.ServerPort); err != nil {
			logx.Errorf("Failed to start server: %v", err)
		}
	}()
	logx.Infof("Server is running on %s", cfg.ServerPort)

	// 3. 阻塞在这里，直到接收到信号
	<-quit
	logx.Info("Gracefully shutting down...")

	// 4. 设置超时时间，防止某些请求卡死导致程序无法退出
	// 这是给 Fiber 留出的"善后"时间
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 5. 依次关闭各个组件（注意顺序：先关入口，再关存储）

	// A. 先关闭 Fiber (不再接收新请求，处理完旧请求)
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logx.Errorf("Fiber shutdown error: %v", err)
	}

	// B. 关闭数据库连接 (sqlc/sql.DB)
	// if err := db.Close(); err != nil { ... }

	// C. 关闭 Redis
	// if err := rdb.Close(); err != nil { ... }

	logx.Info("Server exited properly.")
}
