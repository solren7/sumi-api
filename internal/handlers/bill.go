package handlers

import (
	"strconv"
	"time"

	"fiber/internal/repository/dbgen"
	"fiber/internal/services"
	"fiber/pkg/errorx"

	"github.com/gofiber/fiber/v3"
	"github.com/shopspring/decimal"
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
		Amount:      b.Amount.String(),
		Description: b.Description,
		BillType:    b.BillType,
		Category:    b.Category,
		RecordDate:  b.RecordDate.Format(time.RFC3339),
		CreatedAt:   b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   b.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *Handler) CreateBill(c fiber.Ctx) error {
	req := new(CreateBillRequest)
	if err := c.Bind().Body(req); err != nil {
		return errorx.ErrParamsInvalid
	}

	if req.Amount == "" {
		return errorx.New(fiber.StatusBadRequest, "Amount is required")
	}

	amountNumeric, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return err
	}

	var recordDate time.Time
	if req.RecordDate == "" {
		recordDate = time.Now()
	} else {
		recordDate, err = time.Parse("2006-01-02 15:04:05", req.RecordDate)
		if err != nil {
			recordDate, err = time.Parse("2006-01-02", req.RecordDate)
			if err != nil {
				return err
			}
		}
	}

	bill, err := h.S.Bill.CreateBill(c, services.CreateBillInput{
		Amount:      amountNumeric,
		Description: req.Description,
		BillType:    req.BillType,
		Category:    req.Category,
		RecordDate:  recordDate,
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

	bills, err := h.S.Bill.ListBills(c, services.ListBillsInput{
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
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errorx.ErrParamsInvalid
	}

	err = h.S.Bill.DeleteBill(c, id)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) UpdateBill(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errorx.ErrParamsInvalid
	}

	req := new(UpdateBillRequest)
	if err := c.Bind().Body(req); err != nil {
		return errorx.ErrParamsInvalid
	}

	amountNumeric, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return err
	}

	var recordDate time.Time
	if req.RecordDate == "" {
		recordDate = time.Now()
	} else {
		recordDate, err = time.Parse("2006-01-02 15:04:05", req.RecordDate)
		if err != nil {
			recordDate, err = time.Parse("2006-01-02", req.RecordDate)
			if err != nil {
				return err
			}
		}
	}

	bill, err := h.S.Bill.UpdateBill(c, services.UpdateBillInput{
		ID:          id,
		Amount:      amountNumeric,
		Description: req.Description,
		BillType:    req.BillType,
		Category:    req.Category,
		RecordDate:  recordDate,
	})
	if err != nil {
		return err
	}

	return c.JSON(toBillResponse(bill))
}
