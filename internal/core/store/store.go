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
	OwnerID        int             `json:"owner_id" validate:"required,number" example:"5692562784252"`
	CNPJ           string          `json:"cnpj" validate:"required,cnpj" example:"12345678000190"`
	Name           string          `json:"name" validate:"required,lte=25,gte=3" example:"Delicious Bakery"`
	Description    string          `json:"description" validate:"required,lte=255" example:"Best bakery in town with fresh pastries"`
	Active         bool            `json:"active" validate:"boolean" example:"true"`
	IsOpen         bool            `json:"is_open" validate:"boolean" example:"false"`
	Phone          string          `json:"phone" validate:"required,phone" example:"+5511997590670"`
	Score          int             `json:"score" validate:"required,number" example:"500"`
	Address        address.Address `json:"address" validate:"required"`
	Type           Type            `json:"type" validate:"required,storeType" example:"RESTAURANT"`
	ProfileImage   string          `json:"profile_image" example:"https://example.com/profile.jpg"`
	HeaderImage    string          `json:"header_image" example:"https://example.com/header.jpg"`
	BusinessHours  []BusinessHours `json:"business_hours" validate:"required,dive"`
	PaymentMethods []PaymentMethod `json:"payment_methods" validate:"required,dive"`
	Products       []string        `json:"products_id" validate:"required" example:"['prod1', 'prod2']"`
}

func (s *Store) Validate() error {
	return xvalidator.GetPtInstance().Validate(s)
}

func (s *Store) Activate() {
	s.Active = true
}

func (s *Store) Deactivate() {
	s.Active = false
}

func (s *Store) OpenStore() {
	s.IsOpen = true
}

func (s *Store) CloseStore() {
	s.IsOpen = false
}

func (s *Store) UpdateStoreProfile(name, description,
	phone string, address address.Address, types Type) error {
	s.Name = name
	s.Description = description
	s.Phone = phone
	s.Address = address
	s.Type = types
	return s.Validate()
}

func (s *Store) AddNewProduct(productID string) {
	s.Products = append(s.Products, productID)
}

func (s *Store) AddNewBusinessHour(hour BusinessHours) error {
	err := IsValidBusinessHour(hour)
	if err != nil {
		return err
	}

	for _, v := range s.BusinessHours {
		if hour == v {
			return ErrBusinessHourAlreadyExist
		}
	}
	s.BusinessHours = append(s.BusinessHours, hour)

	return nil
}

func (s *Store) RemoveBusinessHour(hour BusinessHours) error {
	err := IsValidBusinessHour(hour)
	if err != nil {
		return err
	}

	for i, v := range s.BusinessHours {
		if hour == v {
			s.BusinessHours = append(s.BusinessHours[:i], s.BusinessHours[i+1:]...)
			return nil
		}
	}

	return ErrBusinessHourNotExist
}

func (s *Store) AddPaymentMethod(paymentMethod PaymentMethod) error {
	for _, v := range s.PaymentMethods {
		if v == paymentMethod {
			return ErrPaymentMethodAlreadyDefined
		}
	}
	s.PaymentMethods = append(s.PaymentMethods, paymentMethod)
	return nil
}

func (s *Store) RemovePaymentMethod(paymentMethod PaymentMethod) error {
	for i, v := range s.PaymentMethods {
		if v == paymentMethod {
			s.PaymentMethods = append(s.PaymentMethods[:i], s.PaymentMethods[i+1:]...)
		}
		if i == len(s.PaymentMethods)-1 && v != paymentMethod {
			return ErrRemoveInvalidPaymentMethod
		}
	}
	return s.Validate()
}

func (s *Store) ChangeProfileImage(imageLink string) {
	s.ProfileImage = imageLink
}

func (s *Store) ChangeHeaderImage(imageLink string) {
	s.HeaderImage = imageLink
}
