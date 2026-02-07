package apps

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fiber/config"
	"fiber/internal/database"
	"fiber/pkg/logx"

	"github.com/gofiber/fiber/v3"
)

func StartAPIServer(cfg *config.Config) {
	app := fiber.New(fiber.Config{
		AppName:     "MyFiberApp",
		ReadTimeout: 10 * time.Second,
	})

	// 1. 初始化数据库 (增加错误处理)
	logx.Info("Connecting to database...")
	// 注意：这里的 context 仅用于控制"连接超时"，连接成功后 dbPool 依然可用
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 5*time.Second)
	dbPool, err := database.NewPool(dbCtx, cfg)
	dbCancel() // 及时释放 context 资源
	if err != nil {
		logx.WithError(err).Fatal("Failed to connect to Database")
	}

	// 2. 初始化 Redis
	logx.Info("Connecting to redis...")
	redisCtx, redisCancel := context.WithTimeout(context.Background(), 5*time.Second)
	rdb, err := database.NewRedis(redisCtx, cfg)
	redisCancel()
	if err != nil {
		logx.WithError(err).Fatal("Failed to connect to Redis")
	}

	// 3. 准备启动 HTTP 服务
	// 创建一个专门接收 Server 启动错误的通道
	serverShutdown := make(chan struct{}) // 用于标记是否需要执行后续清理

	go func() {
		logx.Infof("Server is starting on %s", cfg.ServerPort)
		if err := app.Listen(cfg.ServerPort); err != nil {
			// 只有非关闭导致的错误才认为是启动失败
			logx.Errorf("Server Listen error: %v", err)
			// 这里不需要处理退出，因为 error 会导致这个协程结束
			// 关键是：主协程怎么知道？见下文 select
			// 发送一个系统中断信号给自己，触发主协程的 quit
			p, _ := os.FindProcess(os.Getpid())
			p.Signal(syscall.SIGINT)
		}
		close(serverShutdown) // 标记 server 协程已退出
	}()

	// 4. 监听信号 (主阻塞点)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// 阻塞在这里，等待 (Ctrl+C) 或者 (Server启动失败发的信号)
	<-quit
	logx.Info("Shutdown signal received...")

	// 5. 优雅关闭流程
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// A. 关闭 Fiber (停止接收新请求)
	// ShutdownWithContext 会通知上面的 go func 停止 Listen，并处理完当前请求
	logx.Info("Shutting down Fiber...")
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logx.WithError(err).Error("Fiber shutdown error")
	}

	// 等待子协程完全退出 (可选，但推荐)
	<-serverShutdown

	// B. 关闭数据库
	logx.Info("Closing database pool...")
	dbPool.Close()

	// C. 关闭 Redis
	logx.Info("Closing Redis connection...")
	if err := rdb.Close(); err != nil {
		logx.WithError(err).Error("Failed to close Redis connection")
	}

	logx.Info("Server exited properly.")
}
