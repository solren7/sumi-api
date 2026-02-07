package services

import (
	"time"

	"fiber/internal/repository/dbgen"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BillService struct {
	q *dbgen.Queries
}

func NewBillService(q *dbgen.Queries) *BillService {
	return &BillService{q: q}
}

type CreateBillInput struct {
	Amount      decimal.Decimal
	Description string
	BillType    int16
	Category    int32
	RecordDate  time.Time
}

func (s *BillService) CreateBill(ctx fiber.Ctx, input CreateBillInput) (*dbgen.Bill, error) {
	userID := ctx.Locals("user_id").(uuid.UUID)
	bill, err := s.q.CreateBill(ctx, dbgen.CreateBillParams{
		UserID:      userID,
		Amount:      input.Amount,
		Description: input.Description,
		BillType:    input.BillType,
		Category:    input.Category,
		RecordDate:  input.RecordDate,
	})
	return &bill, err
}

type ListBillsInput struct {
	UserID string
	Limit  int
	Offset int
}

func (s *BillService) ListBills(ctx fiber.Ctx, input ListBillsInput) ([]dbgen.Bill, error) {
	userID := ctx.Locals("user_id").(uuid.UUID)

	bills, err := s.q.ListBills(ctx, dbgen.ListBillsParams{
		UserID: userID,
		Limit:  int32(input.Limit),
		Offset: int32(input.Offset),
	})
	if err != nil {
		return nil, err
	}

	if bills == nil {
		return []dbgen.Bill{}, nil
	}
	return bills, nil
}

type UpdateBillInput struct {
	ID          int64
	Amount      decimal.Decimal
	Description string
	BillType    int16
	Category    int32
	RecordDate  time.Time
}

func (s *BillService) UpdateBill(ctx fiber.Ctx, input UpdateBillInput) (dbgen.Bill, error) {
	userID := ctx.Locals("user_id").(uuid.UUID)

	bill, err := s.q.UpdateBill(ctx, dbgen.UpdateBillParams{
		ID:          input.ID,
		UserID:      userID,
		Amount:      input.Amount,
		Description: input.Description,
		BillType:    input.BillType,
		Category:    input.Category,
		RecordDate:  input.RecordDate,
	})
	return bill, err
}

func (s *BillService) DeleteBill(ctx fiber.Ctx, billID int64) error {
	userID := ctx.Locals("user_id").(uuid.UUID)

	return s.q.DeleteBill(ctx, dbgen.DeleteBillParams{
		ID:     billID,
		UserID: userID,
	})
}
