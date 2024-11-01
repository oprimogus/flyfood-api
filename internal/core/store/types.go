package store

import (
	"fmt"
	"time"

	"github.com/oprimogus/cardapiogo/internal/core/address"
)

type CreateParams struct {
	CpfCnpj string          `json:"cpfCnpj" validate:"required,cpfCnpj" example:"83193927805"`
	Name    string          `json:"name" validate:"required,lte=25" example:"John"`
	Phone   string          `json:"phone" validate:"required,phone" example:"13997590579"`
	Address address.Address `json:"address" validate:"required"`
	Type    ShopType        `json:"type" validate:"required,shopType" enums:"restaurant, pharmacy, tobbaco, market, convenience, pub"`
}

type CreatedStore struct {
	ID string `json:"id"`
}

type UpdateParams struct {
	ID                 string          `json:"id" validate:"required"`
	Name               string          `json:"name" validate:"required,lte=25"`
	Phone              string          `json:"phone" validate:"required,phone"`
	Address            address.Address `json:"address" validate:"required"`
	Type               ShopType        `json:"type" validate:"required,shopType" enums:"restaurant, pharmacy, tobbaco, market, convenience, pub"`
	PaymentMethodEnums []PaymentMethod `json:"paymentMethod" validate:"dive" enums:"credit, debit, pix, cash"`
}

type StoreBusinessHoursParams struct {
	ID            string                `json:"id" validate:"required"`
	TimeZone      string                `json:"timeZone" validate:"required,timezone"`
	BusinessHours []BusinessHoursParams `json:"businessHours" validate:"dive"`
}

type BusinessHoursParams struct {
	WeekDay     int    `json:"weekDay" validate:"min=0,max=6"`
	OpeningTime string `json:"openingTime" validate:"required,businessHour"`
	ClosingTime string `json:"closingTime" validate:"required,businessHour"`
}

func (b *BusinessHoursParams) Entity(timeZone string) (BusinessHours, error) {

	zone, errLoadTimeZone := time.LoadLocation(timeZone)
	if errLoadTimeZone != nil {
		return BusinessHours{}, fmt.Errorf("fail on load timezone: %w", errLoadTimeZone)
	}

	openingTimeParsed, errOpeningTime := time.Parse(BusinessHourLayout, b.OpeningTime)
	if errOpeningTime != nil {
		return BusinessHours{}, fmt.Errorf("fail in parse openingTime: %w", errOpeningTime)
	}

	openingTime := time.Date(1970, time.January, 1, openingTimeParsed.Hour(), openingTimeParsed.Minute(), openingTimeParsed.Second(), 0, time.UTC)

	closingTimeParsed, errClosingTime := time.Parse(BusinessHourLayout, b.ClosingTime)
	if errClosingTime != nil {
		return BusinessHours{}, fmt.Errorf("fail in parse closingTime: %w", errClosingTime)
	}
	closingTime := time.Date(1970, time.January, 1, closingTimeParsed.Hour(), closingTimeParsed.Minute(), closingTimeParsed.Second(), 0, time.UTC)

	return BusinessHours{
		WeekDay:     b.WeekDay,
		OpeningTime: openingTime,
		ClosingTime: closingTime,
		TimeZone:    zone.String(),
	}, nil
}

type GetStoreByIdOutput struct {
	ID                 string                `json:"id"`
	Name               string                `json:"name"`
	ProfileImage       string                `json:"profileImage"`
	HeaderImage        string                `json:"headerImage"`
	Phone              string                `json:"phone"`
	Score              int                   `json:"score"`
	Address            AddressOutput         `json:"address"`
	Type               ShopType              `json:"type" enums:"restaurant, pharmacy, tobbaco, market, convenience, pub"`
	BusinessHours      []BusinessHoursParams `json:"businessHours"`
	PaymentMethodEnums []PaymentMethod       `json:"paymentMethod" enums:"credit, debit, pix, cash"`
}

type GetStoreByFilterOutput struct {
	ID            string                `json:"id"`
	Name          string                `json:"name"`
	ProfileImage  string                `json:"profileImage"`
	Score         int                   `json:"score"`
	Type          ShopType              `json:"type"`
	Neighborhood  string                `json:"neighborhood"`
	BusinessHours []BusinessHoursParams `json:"businessHours"`
}

type AddressOutput struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
}

func (c CreateParams) Entity(userID string) Store {
	return NewStore(userID, c.Name, c.CpfCnpj, c.Phone, c.Address, c.Type)
}

type GetStoresFilterInput struct {
	Range     int      `json:"range"`
	Score     int      `json:"score"`
	Name      string   `json:"name"`
	City      string   `json:"city"`
	Latitude  string   `json:"latitude"`
	Longitude string   `json:"longitude"`
	Type      ShopType `json:"type"`
}
