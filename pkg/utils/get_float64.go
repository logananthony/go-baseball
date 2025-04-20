package utils

import(
  "database/sql"
)


// Float helper
func GetFloat(n sql.NullFloat64) float64 {
    if n.Valid {
        return n.Float64
    }
    return 0.0 // or you can use NaN or -1 as a sentinel if needed
}

// String helper
func GetString(n sql.NullString) string {
    if n.Valid {
        return n.String
    }
    return "" // or return "NA", etc.
}

func GetInt(n sql.NullInt64) int {
    if n.Valid {
        return int(n.Int64) // cast to int
    }
    return 0 // or -1 or any fallback value
}


func StrToNull(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func IntToNull(i int) sql.NullInt32 {
	if i == -1 {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(i), Valid: true}
}
