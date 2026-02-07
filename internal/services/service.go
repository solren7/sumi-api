package services

import (
	"fiber/config"
	"fiber/internal/repository/dbgen"
)

type Service struct {
	Auth  *AuthService
	Bill  *BillService
	Stats *StatsService
}

func NewService(q *dbgen.Queries, cfg *config.Config) *Service {
	return &Service{
		Auth:  NewAuthService(q, cfg),
		Bill:  NewBillService(q),
		Stats: NewStatsService(q),
	}
}
