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


