package store

import (
	"strconv"

	"github.com/oprimogus/flyfood-api/internal/core/address"
	"github.com/oprimogus/flyfood-api/pkg/validator"
)

type Type int

const (
	Restaurant Type = iota + 1
	Pharmacy
	Tobacco
	Market
	Convenience
	Pub
)

func (t Type) String() string {
	switch t {
	case Restaurant:
		return "RESTAURANT"
	case Pharmacy:
		return "PHARMACY"
	case Tobacco:
		return "TOBACCO"
	case Market:
		return "MARKET"
	case Convenience:
		return "CONVENIENCE"
	case Pub:
		return "PUB"
	default:
		return strconv.Itoa(int(t))
	}
}

type CreateStoreDTO struct {
	CNPJ        string          `json:"cnpj" validate:"required,cnpj" example:"12345678000190"`
	Name        string          `json:"name" validate:"required,lte=25,gte=3" example:"Delicious Bakery"`
	Description string          `json:"description" validate:"required,lte=255" example:"Best bakery in town with fresh pastries"`
	Phone       string          `json:"phone" validate:"required,phone" example:"+5511997590670"`
	Address     address.Address `json:"address" validate:"required"`
	Type        Type            `json:"type" validate:"required,storeType" example:"RESTAURANT"`
}

func (s *CreateStoreDTO) Validate() error {
	return validator.Validate(s)
}

type UpdateStoreDTO struct {
	ID           string          `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name         string          `json:"name" validate:"required,lte=25,gte=3" example:"Delicious Bakery"`
	Description  string          `json:"description" validate:"required,lte=255" example:"Best bakery in town with fresh pastries"`
	Phone        string          `json:"phone" validate:"required,phone" example:"+5511997590670"`
	Address      address.Address `json:"address" validate:"required"`
	Type         Type            `json:"type" validate:"required,storeType" example:"RESTAURANT"`
	DeliveryTime int             `json:"deliveryTime" validate:"number" example:"40"`
}

func (s *UpdateStoreDTO) Validate() error {
	return validator.Validate(s)
}

func IsValidType(storeType int) bool {
	switch Type(storeType) {
	case Restaurant,
		Pharmacy,
		Tobacco,
		Market,
		Convenience,
		Pub:
		return true
	default:
		return false
	}
}
