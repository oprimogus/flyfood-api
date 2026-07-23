package database

import (
	"fmt"
	"time"

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

func TimeToTimestampz(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  value,
		Valid: true,
	}
}

func ToTime(value pgtype.Timestamptz) time.Time {
	if !value.Valid {
		return time.Time{}
	}
	return value.Time
}

func ToPgTypeUUID(value uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: value,
		Valid: true,
	}
}

func ToUUID(value pgtype.UUID) uuid.UUID {
	return uuid.UUID(value.Bytes)
}

func PgTimeToString(t pgtype.Time) string {
    if !t.Valid {
        return ""
    }
    totalSeconds := t.Microseconds / 1_000_000
    hours := totalSeconds / 3600
    minutes := (totalSeconds % 3600) / 60
    return fmt.Sprintf("%02d:%02d", hours, minutes)
}