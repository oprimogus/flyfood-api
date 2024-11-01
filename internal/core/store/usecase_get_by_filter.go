package store

import (
	"context"
)

type UseCaseGetByFilter struct {
	repository Repository
}

func NewUseCaseGetByFilter(repository Repository) UseCaseGetByFilter {
	return UseCaseGetByFilter{
		repository: repository,
	}
}

func (g UseCaseGetByFilter) Execute(ctx context.Context, params GetStoresFilterInput) (*[]GetStoreByFilterOutput, error) {
	stores, err := g.repository.FindByFilter(ctx, params)
	if err != nil {
		return nil, err
	}

	storesFiltered := make([]GetStoreByFilterOutput, len(*stores))
	for i, v := range *stores {
		businesHours := make([]BusinessHoursParams, len(v.BusinessHours))
		for i, v := range v.BusinessHours {
			businesHours[i] = BusinessHoursParams{
				WeekDay:     v.WeekDay,
				OpeningTime: v.OpeningTime.Format(BusinessHourLayout),
				ClosingTime: v.ClosingTime.Format(BusinessHourLayout),
			}
		}
		storesFiltered[i] = GetStoreByFilterOutput{
			ID:            v.ID,
			Name:          v.Name,
			Score:         v.Score,
			Type:          v.Type,
			ProfileImage:  v.ProfileImage,
			Neighborhood:  v.Address.Neighborhood,
			BusinessHours: businesHours,
		}
	}
	return &storesFiltered, nil
}
