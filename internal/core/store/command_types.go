package store

// import (
// 	"github.com/oprimogus/flyfood-api/internal/core/address"
// )

// type Type string

// const (
// 	Restaurant  Type = "RESTAURANT"
// 	Pharmacy    Type = "PHARMACY"
// 	Tobbaco     Type = "TOBBACO"
// 	Market      Type = "MARKET"
// 	Convenience Type = "CONVENIENCE"
// 	Pub         Type = "PUB"
// )





// type AddOrDeleteBusinessHourDTO struct {
// 	StoreID       string        `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
// 	BusinessHours BusinessHours `json:"businessHour"`
// }

// type AddOrDeletePaymentMethodDTO struct {
// 	StoreID        string        `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
// 	PaymentMethods PaymentMethod `json:"paymentMethod" validate:"required,paymentMethod"`
// }

// type SetOpenStateDTO struct {
// 	StoreID string `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
// 	IsOpen  bool   `json:"isOpen" validate:"boolean" example:"true"`
// }

// type SetActiveDTO struct {
// 	StoreID string `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
// 	Active  bool   `json:"active" validate:"boolean" example:"true"`
// }

// type UploadStoreImageDTO struct {
// 	StoreID string `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
// 	Image   []byte `json:"image"`
// 	Ext     string `json:"ext"`
// }
