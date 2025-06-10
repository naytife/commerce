package resolver

import "github.com/jackc/pgx/v5/pgtype"

func safeStringDereference(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

// stringDereference returns the value of a string pointer or an empty string if nil
func stringDereference(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// numericToFloat64 converts a pgtype.Numeric to a float64
func numericToFloat64(n *pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}

	// Convert the numeric to a float64
	f, _ := n.Int.Float64()
	return f
}
