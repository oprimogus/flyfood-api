package store

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/oprimogus/flyfood-api/internal/core/address"
	"github.com/oprimogus/flyfood-api/pkg/validator"
)

const (
	DefaultScore = 500
)

type Store struct {
	ID             uuid.UUID       `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	OwnerID        uuid.UUID       `json:"ownerID" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	CNPJ           string          `json:"cnpj" validate:"required,cnpj" example:"12345678000190"`
	Name           string          `json:"name" validate:"required,lte=25,gte=3" example:"Delicious Bakery"`
	Description    string          `json:"description" validate:"required,lte=255" example:"Best bakery in town with fresh pastries"`
	Active         bool            `json:"active" validate:"boolean" example:"true"`
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

func NewStore(ownerID uuid.UUID, cnpj, name, description,
	phone string, address address.Address, types Type) (*Store, error) {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("fail on create store id: %w", err)
	}
	addressId, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("fail on create address id: %w", err)
	}

	newStore := Store{
		ID:             uuidV7,
		OwnerID:        ownerID,
		CNPJ:           cnpj,
		Name:           name,
		Description:    description,
		Active:         false,
		IsOpen:         false,
		Phone:          phone,
		Score:          DefaultScore,
		Address:        address,
		Type:           types,
		BusinessHours:  []BusinessHours{},
		PaymentMethods: []PaymentMethod{},
	}
	newStore.Address.ID = addressId
	if err := newStore.Validate(); err != nil {
		return nil, err
	}

	return &newStore, nil
}

func (st *Store) Validate() error {
	return validator.Validate(st)
}

func (st *Store) Activate() {
	st.Active = true
}

func (st *Store) Deactivate() {
	st.Active = false
}

func (st *Store) OpenStore() {
	st.IsOpen = true
}

func (st *Store) CloseStore() {
	st.IsOpen = false
}

func (st *Store) UpdateStore(name, description,
	phone string, address address.Address, types Type) error {
	st.Name = name
	st.Description = description
	st.Phone = phone
	st.Address = address
	st.Type = types
	return st.Validate()
}

func (st *Store) AddNewBusinessHour(hour BusinessHours) error {
	err := IsValidBusinessHour(hour)
	if err != nil {
		return err
	}

	for _, v := range st.BusinessHours {
		if hour == v {
			return ErrBusinessHourAlreadyExist
		}
	}
	st.BusinessHours = append(st.BusinessHours, hour)

	return nil
}

func (st *Store) RemoveBusinessHour(hour BusinessHours) error {
	err := IsValidBusinessHour(hour)
	if err != nil {
		return err
	}

	for i, v := range st.BusinessHours {
		if hour == v {
			st.BusinessHours = append(st.BusinessHours[:i], st.BusinessHours[i+1:]...)
			return nil
		}
	}

	return ErrBusinessHourNotExist
}

func (st *Store) AddPaymentMethod(paymentMethod PaymentMethod) error {
	if slices.Contains(st.PaymentMethods, paymentMethod) {
		return ErrPaymentMethodAlreadyDefined
	}
	st.PaymentMethods = append(st.PaymentMethods, paymentMethod)
	return nil
}

func (st *Store) RemovePaymentMethod(paymentMethod PaymentMethod) error {
	for i, v := range st.PaymentMethods {
		if v == paymentMethod {
			st.PaymentMethods = append(st.PaymentMethods[:i], st.PaymentMethods[i+1:]...)
		}
		if i == len(st.PaymentMethods)-1 && v != paymentMethod {
			return ErrRemoveInvalidPaymentMethod
		}
	}
	return st.Validate()
}

func (st *Store) ChangeProfileImage(imageLink string) {
	st.ProfileImage = imageLink
}

func (st *Store) ChangeHeaderImage(imageLink string) {
	st.HeaderImage = imageLink
}
