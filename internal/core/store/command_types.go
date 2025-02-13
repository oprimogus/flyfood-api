package store

import (
	"github.com/oprimogus/cardapiogo/internal/core/address"
)

type Type string

const (
	Restaurant  Type = "RESTAURANT"
	Pharmacy    Type = "PHARMACY"
	Tobbaco     Type = "TOBBACO"
	Market      Type = "MARKET"
	Convenience Type = "CONVENIENCE"
	Pub         Type = "PUB"
)

func IsValidType(storeType string) bool {
	switch Type(storeType) {
	case Restaurant,
		Pharmacy,
		Tobbaco,
		Market,
		Convenience,
		Pub:
		return true
	default:
		return false
	}
}

type CreateNewStoreDTO struct {
	CNPJ        string          `json:"cnpj" validate:"required,cnpj" example:"12345678000190"`
	Name        string          `json:"name" validate:"required,lte=25,gte=3" example:"Delicious Bakery"`
	Description string          `json:"description" validate:"required,lte=255" example:"Best bakery in town with fresh pastries"`
	Phone       string          `json:"phone" validate:"required,phone" example:"+5511997590670"`
	Address     address.Address `json:"address" validate:"required"`
	Type        Type            `json:"type" validate:"required,storeType" example:"RESTAURANT"`
}

type UpdateStoreDTO struct {
	ID          string          `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string          `json:"name" validate:"required,lte=25,gte=3" example:"Delicious Bakery"`
	Description string          `json:"description" validate:"required,lte=255" example:"Best bakery in town with fresh pastries"`
	Phone       string          `json:"phone" validate:"required,phone" example:"+5511997590670"`
	Address     address.Address `json:"address" validate:"required"`
	Type        Type            `json:"type" validate:"required,storeType" example:"RESTAURANT"`
}

type AddOrDeleteBusinessHourDTO struct {
	StoreID       string        `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	BusinessHours BusinessHours `json:"businessHour"`
}

type AddOrDeletePaymentMethodDTO struct {
	StoreID        string        `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	PaymentMethods PaymentMethod `json:"paymentMethod" validate:"required,paymentMethod"`
}

type SetOpenStateDTO struct {
	StoreID string `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	IsOpen  bool   `json:"isOpen" validate:"boolean" example:"true"`
}

type SetActiveDTO struct {
	StoreID string `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Active  bool   `json:"active" validate:"boolean" example:"true"`
}

type UploadStoreImageDTO struct {
	StoreID string `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Image   []byte `json:"image"`
	Ext     string `json:"ext"`
}

type UploadProductImageDTO struct {
	StoreID   string `json:"store_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProductID string `json:"product_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Image     []byte `json:"image"`
	Ext       string `json:"ext"`
}
