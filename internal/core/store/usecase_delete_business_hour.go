package store

import (
	"context"
	"fmt"
)

type useCaseDeleteBusinessHour struct {
	repository Repository
}

func newUseCaseDeleteBusinessHour(repository Repository) useCaseDeleteBusinessHour {
	return useCaseDeleteBusinessHour{repository: repository}
}

func (d useCaseDeleteBusinessHour) Execute(ctx context.Context, params StoreBusinessHoursParams) error {
	userID := ctx.Value("userID").(string)
	if userID == "" {
		return fmt.Errorf("invalid userID: '%s'", userID)
	}

	isOwner, errIsOwner := d.repository.IsOwner(ctx, params.ID, userID)
	if errIsOwner != nil {
		return fmt.Errorf("fail on get store data")
	}

	if !isOwner {
		return ErrNotOwner
	}

	businessHours := make([]BusinessHours, len(params.BusinessHours))
	for i, v := range params.BusinessHours {
		businessHoursConverted, err := v.Entity(params.TimeZone)
		if err != nil {
			return fmt.Errorf("fail on convert businessHour: %w", err)
		}
		businessHours[i] = businessHoursConverted
	}

	err := d.repository.DeleteBusinessHour(ctx, params.ID, businessHours)
	if err != nil {
		return err
	}
	return nil
}
