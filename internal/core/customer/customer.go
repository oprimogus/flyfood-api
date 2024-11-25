package customer

import (
	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/oprimogus/cardapiogo/internal/xvalidator"
)

type Customer struct {
	ID        int               `json:"id" validate:"required,number" example:"295105940221919239"`
	Name      string            `json:"name" validate:"required,alpha,gte=3,lte=25" example:"John"`
	LastName  string            `json:"last_name" validate:"required,gte=3,lte=60" example:"Doe"`
	CPF       string            `json:"cpf" validate:"cpf" example:"52024227090"`
	Email     string            `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Phone     string            `json:"phone" validate:"required,phone" example:"+5513997590579"`
	Addresses []address.Address `json:"addresses" validate:"required,dive"`
	OrdersID  []string          `json:"orders_id"`
}

func (c *Customer) Validate() error {
	return xvalidator.GetPtInstance().Validate(c)
}

func NewCustomer(id int, name, lastName, cpf, email, phone string) (*Customer, error) {
	newCustomer := &Customer{
		ID:        id,
		Name:      name,
		LastName:  lastName,
		CPF:       cpf,
		Email:     email,
		Phone:     phone,
		Addresses: []address.Address{},
	}
	if err := newCustomer.Validate(); err != nil {
		return &Customer{}, err
	}
	return newCustomer, nil
}

func (c *Customer) UpdateProfile(name, lastName, phone string) error {
	c.Name = name
	c.LastName = lastName
	c.Phone = phone
	return c.Validate()
}

func (c *Customer) SaveNewAddress(address address.Address) error {
	if len(c.Addresses) >= 5 {
		return ErrMaxAddresses
	}
	c.Addresses = append(c.Addresses, address)

	return c.Validate()
}

func (c *Customer) RemoveAddress(address address.Address) error {
	if len(c.Addresses) == 0 {
		return ErrThereIsNoAddresses
	}
	for i, v := range c.Addresses {
		if v == address {
			c.Addresses = append(c.Addresses[:i], c.Addresses[i+1:]...)
		}
		if i == len(c.Addresses)-1 && v != address {
			return ErrTryRemoveInvalidAddress
		}
	}

	return c.Validate()
}
