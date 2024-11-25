package store

import (
	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/oprimogus/cardapiogo/internal/core/store"
)

type CreateNewStoreDTO struct {
	CNPJ        string          `json:"cnpj" validate:"required,cnpj" example:"12345678000190"`
	Name        string          `json:"name" validate:"required,lte=25,gte=3" example:"Delicious Bakery"`
	Description string          `json:"description" validate:"required,lte=255" example:"Best bakery in town with fresh pastries"`
	Phone       string          `json:"phone" validate:"required,phone" example:"+5511997590670"`
	Address     address.Address `json:"address" validate:"required"`
	Type        store.Type      `json:"type" validate:"required,storeType" example:"RESTAURANT"`
}

type UpdateStoreDTO struct {
	ID          string          `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string          `json:"name" validate:"required,lte=25,gte=3" example:"Delicious Bakery"`
	Description string          `json:"description" validate:"required,lte=255" example:"Best bakery in town with fresh pastries"`
	Phone       string          `json:"phone" validate:"required,phone" example:"+5511997590670"`
	Address     address.Address `json:"address" validate:"required"`
	Type        store.Type      `json:"type" validate:"required,storeType" example:"RESTAURANT"`
}

type AddOrDeleteBusinessHourDTO struct {
	StoreID       string              `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	BusinessHours store.BusinessHours `json:"business_hour"`
}

type AddOrDeletePaymentMethodDTO struct {
	StoreID        string              `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	PaymentMethods store.PaymentMethod `json:"payment_method" validate:"required,paymentMethod"`
}

type SetOpenStateDTO struct {
	StoreID string `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	IsOpen  bool   `json:"is_open" validate:"boolean" example:"true"`
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
