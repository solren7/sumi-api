package cache

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func RefreshTokenKey(tokenHash string) string {
	return fmt.Sprintf("auth:refresh:%s", tokenHash)
}

func APIKeyKey(prefix string) string {
	return fmt.Sprintf("apikey:%s", prefix)
}

func SystemCategoriesKey(categoryType int16) string {
	return fmt.Sprintf("categories:system:type:%d", categoryType)
}

func MonthlyStatsKey(userID uuid.UUID, month string) string {
	return fmt.Sprintf("stats:monthly:%s:%s", userID.String(), month)
}

func DailyStatsKey(userID uuid.UUID, month string) string {
	return fmt.Sprintf("stats:daily:%s:%s", userID.String(), month)
}

func CategoryStatsKey(userID uuid.UUID, month string, categoryType int16) string {
	return fmt.Sprintf("stats:category:%s:%s:%d", userID.String(), month, categoryType)
}

func MonthKey(t time.Time, loc *time.Location) string {
	return t.In(loc).Format("2006-01")
}
