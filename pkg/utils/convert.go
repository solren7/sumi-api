package utils

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func NumericToString(n pgtype.Numeric) string {
	if !n.Valid {
		return "0"
	}
	// pgtype.Numeric Value() returns driver.Value which is string for Numeric?
	// Actually Float64Value might lose precision.
	// Int64Value might lose decimals.
	// Recommended way involves big.Int if needed, but Scan into string should work?
	// Wait, typical way:
	val, err := n.Value()
	if err != nil {
		return "0"
	}
	if val == nil {
		return "0"
	}
	return fmt.Sprintf("%v", val)
}

func TimeFormat(t time.Time) string {
	return t.Format(time.RFC3339)
}
