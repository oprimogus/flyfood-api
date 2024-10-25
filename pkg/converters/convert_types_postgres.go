package converters

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ConvertUUIDToString converts pgtype.UUID to *string. Returns nil if UUID is not valid.
func ConvertUUIDToString(uuidVal pgtype.UUID) (*string, error) {
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

// ConvertStringToUUID converts string to pgtype.UUID.
func ConvertStringToUUID(str string) (pgtype.UUID, error) {
	u, err := uuid.Parse(str)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("fail on convert string to uuid: %w", err)
	}

	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}, nil
}

// ConvertInt4ToInt converts pgtype.Int4 to *int. Returns nil if Int4 is not valid.
func ConvertInt4ToInt(int4Val pgtype.Int4) (int, error) {
	if !int4Val.Valid {
		return 0, nil
	}

	result := int(int4Val.Int32)
	return result, nil
}

// ConvertIntToInt4 converts int to pgtype.Int4.
func ConvertIntToInt4(i int) pgtype.Int4 {
	return pgtype.Int4{
		Int32: int32(i),
		Valid: true,
	}
}

// ConvertInt32ToInt4 converts int32 to pgtype.int4.
func ConvertInt32ToInt4(int32Value int32) pgtype.Int4 {
	return pgtype.Int4{
		Valid: true,
		Int32: int32Value,
	}
}

// ConvertInt4ToInt32 converts pgtype.int4 to int32.
func ConvertInt4ToInt32(int4Value pgtype.Int4) (int32, error) {
	if !int4Value.Valid {
		return 0, nil
	}
	result := int32(int4Value.Int32)
	return result, nil
}

// ConvertTextToString converts pgtype.Text to *string. Returns nil if Text is not valid.
func ConvertTextToString(textVal pgtype.Text) (string, error) {
	if !textVal.Valid {
		return "", nil
	}

	return textVal.String, nil
}

// ConvertStringToText converts string to pgtype.Text.
func ConvertStringToText(str string) pgtype.Text {
	isValid := true
	if str == "" {
		isValid = false
	}
	return pgtype.Text{
		String: str,
		Valid:  isValid,
	}
}

// ConvertTimestamptzToTime converts pgtype.Timestamptz to *time.Time. Returns nil if Timestamptz is not valid.
func ConvertTimestamptzToTime(tzVal pgtype.Timestamptz) (time.Time, error) {
	if !tzVal.Valid {
		return time.Time{}, nil
	}

	return tzVal.Time, nil
}

// ConvertTimestampToTime converts pgtype.Timestamp to *time.Time. Returns nil if Timestamp is not valid.
func ConvertTimestamptToTime(tzVal pgtype.Timestamp) (time.Time, error) {
	if !tzVal.Valid {
		return time.Time{}, nil
	}

	return tzVal.Time, nil
}

// ConvertTimeToTimestamptz converts time.Time to pgtype.Timestamptz.
func ConvertTimeToTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}

// ConvertTimeToTimestamptz converts time.Time to pgtype.Timestamptz.
func ConvertTimeToTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: true,
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

func PgtypeTime(value time.Time) pgtype.Time {
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
