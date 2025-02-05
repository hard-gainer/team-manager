package util

import "github.com/jackc/pgx/v5/pgtype"

func ToNullInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: int32(i),
		Valid: true,
	}
}
