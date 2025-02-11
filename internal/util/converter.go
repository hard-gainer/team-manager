package util

import "github.com/jackc/pgx/v5/pgtype"

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