package handlers

import "github.com/gofiber/fiber/v3"

type PingResponse struct {
	Message string `json:"message"`
}

// Ping godoc
// @Summary Ping endpoint
// @Tags Health
// @Produce json
// @Success 200 {object} PingResponse
// @Router /api/ping [get]
func (h *Handler) Ping(c fiber.Ctx) error {
	return c.JSON(PingResponse{
		Message: "pong",
	})
}
