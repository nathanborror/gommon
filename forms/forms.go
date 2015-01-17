package forms

import (
	"database/sql"
	"net/http"
	"strconv"
)

// NullFloat64 returns null or a float64 given a request object and
// a FormValue key.
func NullFloat64(key string, r *http.Request) sql.NullFloat64 {
	value := r.FormValue(key)
	parsed, _ := strconv.ParseFloat(value, 64)
	return sql.NullFloat64{Float64: parsed, Valid: true}
}

// Float64 returns a string
func Float64(key string, r *http.Request) float64 {
	value := r.FormValue(key)
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return parsed
}

// NullInt64 returns null or an int64 given a request object and
// a FormValue key.
func NullInt64(key string, r *http.Request) sql.NullInt64 {
	parsed, _ := strconv.ParseInt(r.FormValue(key), 0, 64)
	return sql.NullInt64{Int64: parsed, Valid: true}
}

// Int returns an integer.
func Int(key string, r *http.Request) int {
	parsed, err := strconv.Atoi(r.FormValue(key))
	if err != nil {
		return 0
	}
	return parsed
}

// NullString returns null or a string given a request object and
// a FormValue key.
func NullString(key string, r *http.Request) sql.NullString {
	parsed := r.FormValue(key)
	return sql.NullString{String: parsed, Valid: true}
}

// String returns a string.
func String(key string, r *http.Request) string {
	return r.FormValue(key)
}
