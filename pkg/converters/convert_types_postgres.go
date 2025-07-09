package converters

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UuidToString(uuidVal pgtype.UUID) (*string, error) {
	if !uuidVal.Valid {
		return nil, nil
	}

	u, err := uuid.FromBytes(uuidVal.Bytes[:])
	if err != nil {
		return nil, err
	}

	str := u.String()
	return &str, nil
}

func StringToUUID(str string) (pgtype.UUID, error) {
	u, err := uuid.Parse(str)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("fail on convert string to uuid: %w", err)
	}

	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}, nil
}

func StringToText(str string) pgtype.Text {
	isValid := str != ""
	return pgtype.Text{
		String: str,
		Valid:  isValid,
	}
}

func Int4ToInt(value pgtype.Int4) int {
	if !value.Valid {
		return 0
	}
	return int(value.Int32)
}
