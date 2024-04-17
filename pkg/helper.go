package pkg

import (
	"database/sql"
	"math/rand"
	"strconv"
	"time"
)

func NullStringToString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}

	return ""
}

func NullIntToInt(n sql.NullInt16) int {
	if n.Valid {
		return int(n.Int16)
	}

	return 0
}

func NullFloatToFloat(f sql.NullFloat64) float64 {
	if f.Valid {
		return f.Float64
	}

	return 0.0
}


func EncodeToStringFor(value int) string {

    return strconv.Itoa(value)
}

func EncodeTimestampToString(timestamp time.Time) string {
    return timestamp.Format(time.RFC3339)
}


func GenerateOTP() int {

	return rand.Intn(900000) + 100000
}
