package apps

import (
	"fiber/config"
	"fiber/internal/handlers"
	"fiber/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func RegisterRoutes(app *fiber.App, handler *handlers.Handler, cfg *config.Config) {
	// Middlewares
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	api := app.Group("/api")

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/check-email", handler.CheckEmail)
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	// Protected Routes
	api.Use(middleware.AuthMiddleware(cfg))

	// Bills Routes
	bills := api.Group("/bills")
	bills.Post("/", handler.CreateBill)
	bills.Get("/", handler.ListBills)
	bills.Put("/:id", handler.UpdateBill)
	bills.Delete("/:id", handler.DeleteBill)

	// Stats Routes
	stats := api.Group("/stats")
	stats.Get("/home", handler.GetHomeStats)
}
