package handlers

import (
	"fiber/internal/repository/dbgen"

	"github.com/gofiber/fiber/v3"
)

type HomeStatsResponse struct {
	TotalIncome  string                   `json:"total_income"`
	TotalExpense string                   `json:"total_expense"`
	DailyStats   []dbgen.GetDailyStatsRow `json:"daily_stats"`
}

func (h *Handler) GetHomeStats(c fiber.Ctx) error {
	dateStr := c.Query("date") // Format: YYYY-MM
	stats, err := h.S.Stats.GetHomeStats(c, dateStr)
	if err != nil {
		return err
	}

	return c.JSON(stats)
}
