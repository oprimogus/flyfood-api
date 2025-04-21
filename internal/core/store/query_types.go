package store

import (
	"github.com/oprimogus/flyfood-api/internal/core/address"
	"github.com/oprimogus/flyfood-api/internal/core/store/product"
)

type QueryStore struct {
	ID          string          `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string          `json:"name" validate:"required,lte=25,gte=3" example:"Delicious Bakery"`
	Description string          `json:"description" validate:"required,lte=255" example:"Best bakery in town with fresh pastries"`
	IsOpen      bool            `json:"isOpen" validate:"boolean" example:"false"`
	Phone       string          `json:"phone" validate:"required,phone" example:"+5511997590670"`
	Score       int             `json:"score" validate:"required,number" example:"500"`
	Address     address.Address `json:"address" validate:"required"`
	Type        Type            `json:"type" validate:"required,storeType" example:"RESTAURANT"`
	// DeliveryTime is defined in minutes
	DeliveryTime   int                  `json:"deliveryTime" validate:"number" example:"40"`
	ProfileImage   string               `json:"profileImage" example:"https://example.com/profile.jpg"`
	HeaderImage    string               `json:"headerImage" example:"https://example.com/header.jpg"`
	BusinessHours  []BusinessHours      `json:"businessHours" validate:"required,dive"`
	PaymentMethods []PaymentMethod      `json:"paymentMethods" validate:"required,dive"`
	Products       []product.ProductDTO `json:"products" validate:"required,dive"`
}

type QueryOwnerStore struct {
	ID          string          `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string          `json:"name" validate:"required,lte=25,gte=3" example:"Delicious Bakery"`
	Description string          `json:"description" validate:"required,lte=255" example:"Best bakery in town with fresh pastries"`
	IsOpen      bool            `json:"isOpen" validate:"boolean" example:"false"`
	Phone       string          `json:"phone" validate:"required,phone" example:"+5511997590670"`
	Score       int             `json:"score" validate:"required,number" example:"500"`
	Address     address.Address `json:"address" validate:"required"`
	Type        Type            `json:"type" validate:"required,storeType" example:"RESTAURANT"`
	// DeliveryTime is defined in minutes
	DeliveryTime   int               `json:"deliveryTime" validate:"number" example:"40"`
	ProfileImage   string            `json:"profileImage" example:"https://example.com/profile.jpg"`
	HeaderImage    string            `json:"headerImage" example:"https://example.com/header.jpg"`
	BusinessHours  []BusinessHours   `json:"businessHours" validate:"required,dive"`
	PaymentMethods []PaymentMethod   `json:"paymentMethods" validate:"required,dive"`
	Products       []product.Product `json:"products" validate:"required,dive"`
}

type QueryStoreList struct {
	ID           string `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name         string `json:"name" validate:"required,lte=25,gte=3" example:"Delicious Bakery"`
	IsOpen       bool   `json:"isOpen" validate:"boolean" example:"false"`
	Score        int    `json:"score" validate:"required,number" example:"500"`
	Neighborhood string `json:"neighborhood" validate:"required,lte=25" example:"Downtown"`
	Latitude     string `json:"latitude" example:"40.7128"`
	Longitude    string `json:"longitude" example:"-74.0060"`
	Type         Type   `json:"type" validate:"required,storeType" example:"RESTAURANT"`
	// DeliveryTime is defined in minutes
	DeliveryTime int    `json:"deliveryTime" validate:"number" example:"40"`
	ProfileImage string `json:"profileImage" example:"https://example.com/profile.jpg"`
}

type QueryStoresInput struct {
	Name           *string          `json:"name" example:"Delicious Bakery"`
	IsOpen         *bool            `json:"isOpen" example:"false"`
	Score          *int             `json:"score" example:"500"`
	Type           *Type            `json:"type" example:"RESTAURANT"`
	PaymentMethods *[]PaymentMethod `json:"paymentMethods"`
	City           *string          `json:"city" example:"New York"`
	Page           int              `json:"page" example:"1" example:"1"`
	MaxItems       int              `json:"maxItems" example:"10" example:"10"`
}

func (st *Store) ToQueryStore(products []product.ProductDTO) QueryStore {
	return QueryStore{
		ID:             st.ID,
		Name:           st.Name,
		Description:    st.Description,
		IsOpen:         st.IsOpen,
		Phone:          st.Phone,
		Score:          st.Score,
		Address:        st.Address,
		Type:           st.Type,
		DeliveryTime:   st.DeliveryTime,
		ProfileImage:   st.ProfileImage,
		HeaderImage:    st.HeaderImage,
		BusinessHours:  st.BusinessHours,
		PaymentMethods: st.PaymentMethods,
		Products:       products,
	}
}

func (st *Store) ToQueryOwnerStore(products []product.Product) QueryOwnerStore {
	return QueryOwnerStore{
		ID:             st.ID,
		Name:           st.Name,
		Description:    st.Description,
		IsOpen:         st.IsOpen,
		Phone:          st.Phone,
		Score:          st.Score,
		Address:        st.Address,
		Type:           st.Type,
		DeliveryTime:   st.DeliveryTime,
		ProfileImage:   st.ProfileImage,
		HeaderImage:    st.HeaderImage,
		BusinessHours:  st.BusinessHours,
		PaymentMethods: st.PaymentMethods,
		Products:       products,
	}
}

type QueryOwnerStoreList struct {
	ID     string `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name   string `json:"name" validate:"required,lte=25,gte=3" example:"Delicious Bakery"`
	Active bool   `json:"active" validate:"boolean" example:"true"`
	IsOpen bool   `json:"isOpen" validate:"boolean" example:"false"`
	Score  int    `json:"score" validate:"required,number" example:"500"`
	Type   Type   `json:"type" validate:"required,storeType" example:"RESTAURANT"`
	// DeliveryTime is defined in minutes
	DeliveryTime int    `json:"deliveryTime" validate:"number" example:"40"`
	ProfileImage string `json:"profileImage" example:"https://example.com/profile.jpg"`
	City         string `json:"city" validate:"required,lte=25" example:"New York"`
	State        string `json:"state" validate:"required,lte=15" example:"NY"`
	Country      string `json:"country" validate:"required,lte=15" example:"United States"`
}
