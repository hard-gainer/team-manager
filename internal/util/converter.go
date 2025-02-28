package util

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: true,
	}
}

func ParseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Now()
	}
	return date
}

func ToNullInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: int32(i),
		Valid: true,
	}
}

func ToNullInt8(i int64) pgtype.Int8 {
	return pgtype.Int8{
		Int64: int64(i),
		Valid: true,
	}
}
