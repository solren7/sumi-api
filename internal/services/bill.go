package services

import (
	"context"
	"time"

	"fiber/internal/repository/dbgen"
	"fiber/pkg/errorx"
	"fiber/pkg/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

type BillService struct {
	q *dbgen.Queries
}

func NewBillService(q *dbgen.Queries) *BillService {
	return &BillService{q: q}
}

type CreateBillInput struct {
	UserID      string
	Amount      string
	Description string
	BillType    int16
	Category    int32
	RecordDate  string
}

func (s *BillService) CreateBill(ctx context.Context, input CreateBillInput) (*dbgen.Bill, error) {
	userID, err := utils.StringToUUID(input.UserID)
	if err != nil {
		return nil, err
	}

	amountNumeric, err := utils.StringToNumeric(input.Amount)
	if err != nil {
		return nil, err
	}

	// Validate numeric
	val, _ := amountNumeric.Value()
	if val == nil {
		return nil, errorx.ErrParamsInvalid
	}

	var recordDate time.Time
	if input.RecordDate == "" {
		recordDate = time.Now()
	} else {
		recordDate, err = time.Parse("2006-01-02 15:04:05", input.RecordDate)
		if err != nil {
			recordDate, err = time.Parse("2006-01-02", input.RecordDate)
			if err != nil {
				return nil, err
			}
		}
	}

	bill, err := s.q.CreateBill(ctx, dbgen.CreateBillParams{
		UserID:      userID,
		Amount:      amountNumeric,
		Description: input.Description,
		BillType:    input.BillType,
		Category:    input.Category,
		RecordDate:  pgtype.Timestamptz{Time: recordDate, Valid: true},
	})
	return &bill, err
}

type ListBillsInput struct {
	UserID string
	Limit  int
	Offset int
}

func (s *BillService) ListBills(ctx context.Context, input ListBillsInput) ([]dbgen.Bill, error) {
	userID, err := utils.StringToUUID(input.UserID)
	if err != nil {
		return nil, err
	}

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
	UserID      string
	Amount      string
	Description string
	BillType    int16
	Category    int32
	RecordDate  string
}

func (s *BillService) UpdateBill(ctx context.Context, input UpdateBillInput) (dbgen.Bill, error) {
	userID, err := utils.StringToUUID(input.UserID)
	if err != nil {
		return dbgen.Bill{}, err
	}

	amountNumeric, err := utils.StringToNumeric(input.Amount)
	if err != nil {
		return dbgen.Bill{}, err
	}

	var recordDate time.Time
	if input.RecordDate == "" {
		recordDate = time.Now()
	} else {
		recordDate, err = time.Parse("2006-01-02 15:04:05", input.RecordDate)
		if err != nil {
			recordDate, err = time.Parse("2006-01-02", input.RecordDate)
			if err != nil {
				return dbgen.Bill{}, err
			}
		}
	}

	bill, err := s.q.UpdateBill(ctx, dbgen.UpdateBillParams{
		ID:          input.ID,
		UserID:      userID,
		Amount:      amountNumeric,
		Description: input.Description,
		BillType:    input.BillType,
		Category:    input.Category,
		RecordDate:  pgtype.Timestamptz{Time: recordDate, Valid: true},
	})
	return bill, err
}

func (s *BillService) DeleteBill(ctx context.Context, userIDStr string, billID int64) error {
	userID, err := utils.StringToUUID(userIDStr)
	if err != nil {
		return err
	}

	return s.q.DeleteBill(ctx, dbgen.DeleteBillParams{
		ID:     billID,
		UserID: userID,
	})
}
