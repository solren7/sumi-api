package apps

import (
	"fiber/internal/handlers"
	"fiber/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func RegisterRoutes(app *fiber.App, handler *handlers.Handler, _ any) {
	// Middlewares
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(middleware.LogxMeta)

	// Routes
	api := app.Group("/api")

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/check-email", handler.CheckEmail)
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
	auth.Post("/refresh", handler.Refresh)
	auth.Post("/logout", handler.Logout)
	auth.Get("/me", middleware.JWTOnly(handler.S.Auth), handler.Me)

	apiKeys := api.Group("/api-keys", middleware.JWTOnly(handler.S.Auth))
	apiKeys.Post("/", handler.CreateAPIKey)
	apiKeys.Get("/", handler.ListAPIKeys)
	apiKeys.Post("/:id/revoke", handler.RevokeAPIKey)
	apiKeys.Delete("/:id", handler.DeleteAPIKey)

	categories := api.Group("/categories", middleware.JWTOrAPIKey(handler.S.Auth, handler.S.APIKey, "categories:read"))
	categories.Get("/", handler.ListCategories)

	authenticatedBills := middleware.JWTOrAPIKey(handler.S.Auth, handler.S.APIKey, "transactions:read")
	writeBills := middleware.JWTOrAPIKey(handler.S.Auth, handler.S.APIKey, "transactions:write")
	updateBills := middleware.JWTOrAPIKey(handler.S.Auth, handler.S.APIKey, "transactions:update")
	deleteBills := middleware.JWTOrAPIKey(handler.S.Auth, handler.S.APIKey, "transactions:delete")

	// Bills Routes
	bills := api.Group("/bills")
	bills.Get("/", authenticatedBills, handler.ListBills)
	bills.Get("/:id", authenticatedBills, handler.GetBill)
	bills.Post("/", writeBills, handler.CreateBill)
	bills.Put("/:id", updateBills, handler.UpdateBill)
	bills.Delete("/:id", deleteBills, handler.DeleteBill)

	transactions := api.Group("/transactions")
	transactions.Get("/", authenticatedBills, handler.ListBills)
	transactions.Get("/:id", authenticatedBills, handler.GetBill)
	transactions.Post("/", writeBills, handler.CreateBill)
	transactions.Put("/:id", updateBills, handler.UpdateBill)
	transactions.Delete("/:id", deleteBills, handler.DeleteBill)

	stats := api.Group("/stats", middleware.JWTOrAPIKey(handler.S.Auth, handler.S.APIKey, "stats:read"))
	stats.Get("/monthly", handler.GetMonthlyStats)
	stats.Get("/daily", handler.GetDailyStats)
	stats.Get("/category", handler.GetCategoryStats)

	// Backward-compatible alias for the original home stats endpoint.
	homeStats := api.Group("/stats")
	homeStats.Get("/home", middleware.JWTOrAPIKey(handler.S.Auth, handler.S.APIKey, "stats:read"), handler.GetMonthlyStats)
}
