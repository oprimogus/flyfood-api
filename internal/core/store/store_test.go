package store

import (
	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStore_AddNewBusinessHour(t *testing.T) {

	tests := []struct {
		name               string
		actualBusinessHour []BusinessHours
		newBusinessHour    BusinessHours
		expected           error
	}{
		{
			name:               "should add new business hour",
			actualBusinessHour: []BusinessHours{},
			newBusinessHour: BusinessHours{
				WeekDay:     1,
				OpeningTime: "13:00",
				ClosingTime: "18:00",
			},
			expected: nil,
		},
		{
			name: "should return error when try add same business hour",
			actualBusinessHour: []BusinessHours{
				{
					WeekDay:     1,
					OpeningTime: "12:00",
					ClosingTime: "21:00",
				},
			},
			newBusinessHour: BusinessHours{
				WeekDay:     1,
				OpeningTime: "12:00",
				ClosingTime: "21:00",
			},
			expected: ErrBusinessHourAlreadyExist,
		},
		{
			name: "should return error when try add invalid business hour",
			actualBusinessHour: []BusinessHours{
				{
					WeekDay:     1,
					OpeningTime: "12:00",
					ClosingTime: "21:00",
				},
			},
			newBusinessHour: BusinessHours{
				WeekDay:     2,
				OpeningTime: "17:00",
				ClosingTime: "14:00",
			},
			expected: ErrClosingTimeBeforeOpeningTime,
		},
		{
			name: "should return error when try add invalid business hour where opening equal closing hour",
			actualBusinessHour: []BusinessHours{
				{
					WeekDay:     1,
					OpeningTime: "12:00",
					ClosingTime: "21:00",
				},
			},
			newBusinessHour: BusinessHours{
				WeekDay:     2,
				OpeningTime: "17:00",
				ClosingTime: "17:00",
			},
			expected: ErrOpeningHourEqualClosingHour,
		},
	}

	for _, test := range tests {

		mockOwner := NewOwner("8647877478")

		addr := address.Address{
			Name:         "Store name location",
			AddressLine1: "rua 1",
			AddressLine2: "879",
			Neighborhood: "test",
			City:         "test",
			State:        "test",
			PostalCode:   "11490-135",
			Country:      "Brasil",
		}

		storeCreated, err := mockOwner.NewStore(
			"63432495000148",
			"Store test",
			"store from test",
			"+5513997590579",
			addr,
			Restaurant)
		assert.NoError(t, err)

		storeCreated.BusinessHours = test.actualBusinessHour

		err = storeCreated.AddNewBusinessHour(test.newBusinessHour)
		assert.Equal(t, test.expected, err)
	}
}
