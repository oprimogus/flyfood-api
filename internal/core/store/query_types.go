package store

import "github.com/oprimogus/cardapiogo/internal/core/address"

type QueryStore struct {
	ID             string          `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name           string          `json:"name" validate:"required,lte=25,gte=3" example:"Delicious Bakery"`
	Description    string          `json:"description" validate:"required,lte=255" example:"Best bakery in town with fresh pastries"`
	IsOpen         bool            `json:"isOpen" validate:"boolean" example:"false"`
	Phone          string          `json:"phone" validate:"required,phone" example:"+5511997590670"`
	Score          int             `json:"score" validate:"required,number" example:"500"`
	Address        address.Address `json:"address" validate:"required"`
	Type           Type            `json:"type" validate:"required,storeType" example:"RESTAURANT"`
	ProfileImage   string          `json:"profileImage" example:"https://example.com/profile.jpg"`
	HeaderImage    string          `json:"headerImage" example:"https://example.com/header.jpg"`
	BusinessHours  []BusinessHours `json:"businessHours" validate:"required,dive"`
	PaymentMethods []PaymentMethod `json:"paymentMethods" validate:"required,dive"`
}

type QueryStoreList struct {
	ID           string `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name         string `json:"name" validate:"required,lte=25,gte=3" example:"Delicious Bakery"`
	IsOpen       bool   `json:"isOpen" validate:"boolean" example:"false"`
	Score        int    `json:"score" validate:"required,number" example:"500"`
	Neighborhood string `json:"neighborhood" validate:"required,lte=25" example:"Downtown"`
	Type         Type   `json:"type" validate:"required,storeType" example:"RESTAURANT"`
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

func (st *Store) ToQueryStore() QueryStore {
	return QueryStore{
		ID:             st.ID,
		Name:           st.Name,
		Description:    st.Description,
		IsOpen:         st.IsOpen,
		Phone:          st.Phone,
		Score:          st.Score,
		Address:        st.Address,
		Type:           st.Type,
		ProfileImage:   st.ProfileImage,
		HeaderImage:    st.HeaderImage,
		BusinessHours:  st.BusinessHours,
		PaymentMethods: st.PaymentMethods,
	}
}
