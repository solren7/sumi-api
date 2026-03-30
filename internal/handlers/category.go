package handlers

import (
	"strconv"

	"fiber/pkg/errorx"

	"github.com/gofiber/fiber/v3"
)

// ListCategories godoc
// @Summary List system categories
// @Tags Categories
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param type query int false "Category type: 1 expense, 2 income" default(1)
// @Success 200 {array} services.CategoryNode
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/categories [get]
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
