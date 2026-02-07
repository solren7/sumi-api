package utils

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// StringToNumeric converts a string amount to pgtype.Numeric
func StringToNumeric(s string) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	err := n.Scan(s)
	return n, err
}

// NumericToString converts pgtype.Numeric to string
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

// FormatNumeric formats pgtype.Numeric to a standard 2-decimal string if possible,
// or just returns the string representation.
func FormatNumeric(n pgtype.Numeric) string {
	if !n.Valid {
		return "0.00"
	}

	// Convert to float for formatting (simple approach for display)
	// Or parse string and reformat.
	f, err := n.Float64Value()
	if err != nil {
		return "0.00" // If conversion to float fails, return default formatted string
	}
	return fmt.Sprintf("%.2f", f.Float64)
}
