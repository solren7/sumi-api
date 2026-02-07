package services

import (
	"time"

	"fiber/internal/repository/dbgen"
	"fiber/pkg/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type StatsService struct {
	q *dbgen.Queries
}

func NewStatsService(q *dbgen.Queries) *StatsService {
	return &StatsService{q: q}
}

type HomeStatsOutput struct {
	TotalIncome  string
	TotalExpense string
	DailyStats   []dbgen.GetDailyStatsRow
}

func (s *StatsService) GetHomeStats(ctx fiber.Ctx, dateStr string) (*HomeStatsOutput, error) {
	var startTime, endTime time.Time
	userID := ctx.Locals("user_id").(uuid.UUID)
	if dateStr == "" {
		now := time.Now()
		startTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		endTime = startTime.AddDate(0, 1, 0).Add(-time.Second)
	} else {
		parsedDate, err := time.Parse("2006-01", dateStr)
		if err != nil {
			return nil, err
		}
		startTime = time.Date(parsedDate.Year(), parsedDate.Month(), 1, 0, 0, 0, 0, time.Local)
		endTime = startTime.AddDate(0, 1, 0).Add(-time.Second)
	}

	// Fetch Monthly Stats
	monthlyStats, err := s.q.GetMonthlyStats(ctx, dbgen.GetMonthlyStatsParams{
		UserID:    userID,
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		return nil, err
	}

	// Fetch Daily Stats
	dailyStats, err := s.q.GetDailyStats(ctx, dbgen.GetDailyStatsParams{
		UserID:    userID,
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		return nil, err
	}

	if dailyStats == nil {
		dailyStats = []dbgen.GetDailyStatsRow{}
	}

	return &HomeStatsOutput{
		TotalIncome:  utils.NumericToString(monthlyStats.TotalIncome),
		TotalExpense: utils.NumericToString(monthlyStats.TotalExpense),
		DailyStats:   dailyStats,
	}, nil
}
