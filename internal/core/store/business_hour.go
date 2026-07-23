package store

import (
	"fmt"
	"strconv"
	"strings"
)

const MaxMinutesInDay = 1439

// MinutesOfDay representa os minutos decorridos desde 00:00 (0 a 1439)
type MinutesOfDay int16

// NewMinutesOfDayFromHHMM converte "08:30" -> 510 minutos
func NewMinutesOfDayFromHHMM(hhmm string) (MinutesOfDay, error) {
	parts := strings.Split(hhmm, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("formato inválido, use HH:MM")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil || hours < 0 || hours > 23 {
		return 0, fmt.Errorf("hora inválida")
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil || minutes < 0 || minutes > 59 {
		return 0, fmt.Errorf("minutos inválidos")
	}

	return MinutesOfDay(hours*60 + minutes), nil
}

func NewMinutesOfDayFromInt(totalMinutes int) (MinutesOfDay, error) {
	if totalMinutes < 0 || totalMinutes > MaxMinutesInDay {
		return 0, fmt.Errorf("minutos inválidos: deve estar entre 0 e %d", MaxMinutesInDay)
	}

	return MinutesOfDay(totalMinutes), nil
}

func (m MinutesOfDay) String() string {
	hours := m / 60
	minutes := m % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}


type BusinessHours struct {
	WeekDay     int    `json:"weekDay" validate:"min=0,max=6" example:"1"`
	OpeningTime MinutesOfDay `json:"openingTime" validate:"businessHour,required" example:"09:00"`
	ClosingTime MinutesOfDay `json:"closingTime" validate:"businessHour,required" example:"15:00"`
}

func isValidHour(input string) bool {
	_, err := NewMinutesOfDayFromHHMM(input)
	return err == nil
}

func isValidTimeRange(start, end string) error {
	if start == end {
		return ErrOpeningHourEqualClosingHour
	}
	_, err := NewMinutesOfDayFromHHMM(start)
	if err != nil {
		return ErrInvalidHour
	}

	_, err = NewMinutesOfDayFromHHMM(end)
	if err != nil {
		return ErrInvalidHour
	}

	return nil
}

func IsValidBusinessHour(bh BusinessHours) error {
	if bh.WeekDay > 6 || bh.WeekDay < 0 {
		return ErrInvalidWeekDayBusinessHour
	}
	return nil
}
