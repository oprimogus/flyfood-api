package store

import (
	"time"
)

const BusinessHourLayout = "15:04"

type BusinessHours struct {
	WeekDay     int    `json:"week_day" validate:"min=0,max=6" example:"1"`
	OpeningTime string `json:"opening_time" validate:"businessHour,required" example:"09:00"`
	ClosingTime string `json:"closing_time" validate:"businessHour,required" example:"15:00"`
}

func isValidHour(input string) bool {
	_, err := time.Parse(BusinessHourLayout, input)
	return err == nil
}

func isValidTimeRange(start, end string) error {
	if start == end {
		return ErrOpeningHourEqualClosingHour
	}
	startTime, err := time.Parse(BusinessHourLayout, start)
	if err != nil {
		return ErrInvalidHour
	}

	endTime, err := time.Parse(BusinessHourLayout, end)
	if err != nil {
		return ErrInvalidHour
	}

	if !endTime.After(startTime) {
		return ErrClosingTimeBeforeOpeningTime
	}

	return nil
}

func IsValidBusinessHour(bh BusinessHours) error {
	if bh.WeekDay > 6 || bh.WeekDay < 0 {
		return ErrInvalidWeekDayBusinessHour
	}

	if err := isValidTimeRange(bh.OpeningTime, bh.ClosingTime); err != nil {
		return err
	}

	return nil
}
