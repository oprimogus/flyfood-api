package customer

import (
	"github.com/google/uuid"
	"github.com/oprimogus/flyfood-api/pkg/validator"
)

type Customer struct {
	ID            uuid.UUID            `json:"id" validate:"required" example:"295105940221919239"`
	ExternalID    string            `json:"externalId" validate:"required" example:"295105940221919239"`
	Name          string            `json:"name" validate:"required,alpha,gte=3,lte=25" example:"John"`
	LastName      string            `json:"lastName" validate:"required,gte=3,lte=60" example:"Doe"`
	CPF           string            `json:"cpf,omitzero" validate:"cpf" example:"52024227090"`
	Email         string            `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Phone         string            `json:"phone" validate:"phone" example:"+5513997590579"`
}

func NewCustomer(externalID, name, lastName, email, cpf, phone string) (*Customer, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	newCustomer := &Customer{
		ID:          id,
		ExternalID:  externalID,
		Name:        name,
		LastName:    lastName,
		CPF:         cpf,
		Email:       email,
		Phone:       phone,
	}
	return newCustomer, nil
}

func (c *Customer) Validate() error {
	return validator.Validate(c)
}

func (c *Customer) UpdateProfile(name, lastName, phone string) error {
	c.Name = name
	c.LastName = lastName
	c.Phone = phone
	return c.Validate()
}

// func (c *Customer) AddNewAddress(addr address.CreateAddressDTO) error {
//     if len(c.Addresses) >= 5 {
//         return errMaxAddresses
//     }
//     id, err := uuid.NewV7()
//     if err != nil {
//         return err
//     }
// 	newAddr := address.Address{
// 		ID:           id,
// 		Name:         addr.Name,
// 		AddressLine1: addr.AddressLine1,
// 		AddressLine2: addr.AddressLine2,
// 		Neighborhood: addr.Neighborhood,
// 		City:         addr.City,
// 		State:        addr.State,
// 		PostalCode:   addr.PostalCode,
// 		Country:      addr.Country,
// 	}
// 	c.Addresses = append(c.Addresses, newAddr)
// 	return nil
// }