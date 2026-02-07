package handlers

import (
	"context"
	"strconv"
	"time"

	"fiber/internal/repository/dbgen"
	"fiber/internal/services"
	"fiber/pkg/errorx"
	"fiber/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

type CreateBillRequest struct {
	Amount      string `json:"amount"`
	Description string `json:"description"`
	BillType    int16  `json:"bill_type"` // 1: Expense, 2: Income
	Category    int32  `json:"category"`
	RecordDate  string `json:"record_date"` // Format: "2006-01-02 15:04:05"
}

type UpdateBillRequest struct {
	Amount      string `json:"amount"`
	Description string `json:"description"`
	BillType    int16  `json:"bill_type"`
	Category    int32  `json:"category"`
	RecordDate  string `json:"record_date"`
}

type BillResponse struct {
	ID          int64  `json:"id"`
	UserID      string `json:"user_id"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
	BillType    int16  `json:"bill_type"`
	Category    int32  `json:"category"`
	RecordDate  string `json:"record_date"` // Keeping as string for JSON
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func toBillResponse(b dbgen.Bill) BillResponse {
	return BillResponse{
		ID:          b.ID,
		UserID:      utils.UUIDToString(b.UserID),
		Amount:      utils.FormatNumeric(b.Amount),
		Description: b.Description,
		BillType:    b.BillType,
		Category:    b.Category,
		RecordDate:  b.RecordDate.Time.Format(time.RFC3339),
		CreatedAt:   b.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:   b.UpdatedAt.Time.Format(time.RFC3339),
	}
}

func (h *Handler) CreateBill(c fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)

	req := new(CreateBillRequest)
	if err := c.Bind().Body(req); err != nil {
		return errorx.ErrParamsInvalid
	}

	if req.Amount == "" {
		return errorx.New(fiber.StatusBadRequest, "Amount is required")
	}

	bill, err := h.S.Bill.CreateBill(context.Background(), services.CreateBillInput{
		UserID:      userIDStr,
		Amount:      req.Amount,
		Description: req.Description,
		BillType:    req.BillType,
		Category:    req.Category,
		RecordDate:  req.RecordDate,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(toBillResponse(*bill))
}

func (h *Handler) ListBills(c fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)

	limitStr := c.Query("limit", "20")
	limit, _ := strconv.Atoi(limitStr)

	offsetStr := c.Query("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)

	bills, err := h.S.Bill.ListBills(context.Background(), services.ListBillsInput{
		UserID: userIDStr,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return err
	}

	response := make([]BillResponse, len(bills))
	for i, b := range bills {
		response[i] = toBillResponse(b)
	}

	return c.JSON(response)
}

func (h *Handler) DeleteBill(c fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errorx.ErrParamsInvalid
	}

	err = h.S.Bill.DeleteBill(context.Background(), userIDStr, id)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) UpdateBill(c fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errorx.ErrParamsInvalid
	}

	req := new(UpdateBillRequest)
	if err := c.Bind().Body(req); err != nil {
		return errorx.ErrParamsInvalid
	}

	bill, err := h.S.Bill.UpdateBill(context.Background(), services.UpdateBillInput{
		ID:          id,
		UserID:      userIDStr,
		Amount:      req.Amount,
		Description: req.Description,
		BillType:    req.BillType,
		Category:    req.Category,
		RecordDate:  req.RecordDate,
	})
	if err != nil {
		return err
	}

	return c.JSON(toBillResponse(bill))
}
