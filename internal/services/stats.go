package services

import (
	"context"
	"time"

	"fiber/internal/repository/dbgen"
	"fiber/pkg/utils"

	"github.com/jackc/pgx/v5/pgtype"
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

func (s *StatsService) GetHomeStats(ctx context.Context, userIDStr string, dateStr string) (*HomeStatsOutput, error) {
	userID, err := utils.StringToUUID(userIDStr)
	if err != nil {
		return nil, err
	}

	var startTime, endTime time.Time

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
		StartTime: pgtype.Timestamptz{Time: startTime, Valid: true},
		EndTime:   pgtype.Timestamptz{Time: endTime, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// Fetch Daily Stats
	dailyStats, err := s.q.GetDailyStats(ctx, dbgen.GetDailyStatsParams{
		UserID:    userID,
		StartTime: pgtype.Timestamptz{Time: startTime, Valid: true},
		EndTime:   pgtype.Timestamptz{Time: endTime, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	if dailyStats == nil {
		dailyStats = []dbgen.GetDailyStatsRow{}
	}

	return &HomeStatsOutput{
		TotalIncome:  utils.FormatNumeric(monthlyStats.TotalIncome),
		TotalExpense: utils.FormatNumeric(monthlyStats.TotalExpense),
		DailyStats:   dailyStats,
	}, nil
}
