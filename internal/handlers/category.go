package handlers

import (
	"strconv"

	"fiber/pkg/errorx"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) ListCategories(c fiber.Ctx) error {
	rawType := c.Query("type", "1")
	categoryType, err := strconv.ParseInt(rawType, 10, 16)
	if err != nil {
		return errorx.ErrParamsInvalid
	}

	items, err := h.S.Category.ListSystemCategories(c.Context(), int16(categoryType))
	if err != nil {
		return err
	}

	return c.JSON(items)
}
