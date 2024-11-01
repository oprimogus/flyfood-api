package converters

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
	isValid := true
	if str == "" {
		isValid = false
	}
	return pgtype.Text{
		String: str,
		Valid:  isValid,
	}
}

func Time(value pgtype.Time) (time.Time, error) {
	baseDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	if value.Microseconds < 0 || value.Microseconds >= 24*60*60*1e6 {
		return time.Time{}, fmt.Errorf("time out of range")
	}

	duration := time.Duration(value.Microseconds) * time.Microsecond

	result := baseDate.Add(duration)
	return result, nil
}

func TimestampToTime(value pgtype.Timestamp) (*time.Time, error) {
	if value.Valid {
		return &value.Time, nil
	}
	return nil, nil
}

func TimeToPgTime(value time.Time) pgtype.Time {
	dayStart := time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
	microsecondsSinceMidnight := (value.Sub(dayStart)).Microseconds()

	if value.Hour() < 0 || value.Hour() >= 24 ||
		value.Minute() < 0 || value.Minute() >= 60 ||
		value.Second() < 0 || value.Second() >= 60 {
		return pgtype.Time{Valid: false}
	}

	return pgtype.Time{
		Microseconds: microsecondsSinceMidnight,
		Valid:        true,
	}
}
