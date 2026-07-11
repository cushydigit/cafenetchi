package mapper

import "github.com/jackc/pgx/v5/pgtype"

func StringToPgText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  s != "",
	}
}
