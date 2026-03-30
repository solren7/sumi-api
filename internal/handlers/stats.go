package handlers

import (
	"strconv"

	"fiber/middleware"
	"fiber/pkg/errorx"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) GetMonthlyStats(c fiber.Ctx) error {
	userID, err := middleware.UserID(c)
	if err != nil {
		return err
	}

	stats, err := h.S.Stats.GetMonthlyStats(c.Context(), userID, c.Query("month"))
	if err != nil {
		return err
	}
	return c.JSON(stats)
}

func (h *Handler) GetDailyStats(c fiber.Ctx) error {
	userID, err := middleware.UserID(c)
	if err != nil {
		return err
	}

	stats, err := h.S.Stats.GetDailyStats(c.Context(), userID, c.Query("month"))
	if err != nil {
		return err
	}
	return c.JSON(stats)
}

func (h *Handler) GetCategoryStats(c fiber.Ctx) error {
	userID, err := middleware.UserID(c)
	if err != nil {
		return err
	}

	rawType := c.Query("type", "1")
	parsedType, err := strconv.ParseInt(rawType, 10, 16)
	if err != nil {
		return errorx.ErrParamsInvalid
	}

	stats, err := h.S.Stats.GetCategoryStats(c.Context(), userID, c.Query("month"), int16(parsedType))
	if err != nil {
		return err
	}
	return c.JSON(stats)
}
