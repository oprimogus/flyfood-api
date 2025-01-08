package store

import (
	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/oprimogus/cardapiogo/internal/xvalidator"
)

const (
	DefaultScore = 500
)

func init() {
	validator := xvalidator.GetPtInstance()
	err := validator.AddValidations(validationsMap)
	if err != nil {
		panic(err)
	}
}

type Store struct {
	ID             string          `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	OwnerID        int             `json:"ownerID" validate:"required,number" example:"5692562784252"`
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
	Products       []string        `json:"productsID" validate:"required" example:"['prod1', 'prod2']"`
}

func (st *Store) Validate() error {
	return xvalidator.GetPtInstance().Validate(st)
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

func (st *Store) UpdateStoreProfile(name, description,
	phone string, address address.Address, types Type) error {
	st.Name = name
	st.Description = description
	st.Phone = phone
	st.Address = address
	st.Type = types
	return st.Validate()
}

func (st *Store) AddNewProduct(productID string) {
	st.Products = append(st.Products, productID)
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
	for _, v := range st.PaymentMethods {
		if v == paymentMethod {
			return ErrPaymentMethodAlreadyDefined
		}
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
