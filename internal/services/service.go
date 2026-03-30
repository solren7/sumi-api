package services

import (
	"fiber/config"
	"fiber/internal/repository/dbgen"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	Auth     *AuthService
	APIKey   *APIKeyService
	Category *CategoryService
	Bill     *BillService
	Stats    *StatsService
}

func NewService(q *dbgen.Queries, cfg *config.Config, rdb redis.UniversalClient) *Service {
	statsSvc := NewStatsService(q, cfg, rdb)

	return &Service{
		Auth:     NewAuthService(q, cfg, rdb),
		APIKey:   NewAPIKeyService(q, cfg, rdb),
		Category: NewCategoryService(q, cfg, rdb),
		Bill:     NewBillService(q, cfg, rdb, statsSvc),
		Stats:    statsSvc,
	}
}
