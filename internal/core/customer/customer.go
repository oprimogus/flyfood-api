package customer

import (
	"github.com/oprimogus/flyfood-api/internal/core/address"
	"github.com/oprimogus/flyfood-api/internal/xvalidator"
)

type Customer struct {
	ID        string            `json:"id" validate:"required" example:"295105940221919239"`
	Name      string            `json:"name" validate:"required,alpha,gte=3,lte=25" example:"John"`
	LastName  string            `json:"lastName" validate:"required,gte=3,lte=60" example:"Doe"`
	CPF       string            `json:"cpf,omitzero" validate:"cpf" example:"52024227090"`
	Email     string            `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Phone     string            `json:"phone" validate:"phone" example:"+5513997590579"`
	Addresses []address.Address `json:"addresses" validate:"dive"`
}

func (c *Customer) Validate() error {
	return xvalidator.GetPtInstance().Validate(c)
}

func NewCustomer(id, name, lastName, cpf, email, phone string) (*Customer, error) {
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

func (c *Customer) RemoveAddress(addrToRemove address.Address) error {
	if len(c.Addresses) == 0 {
		return ErrThereIsNoAddresses
	}
	var addrs []address.Address
	for _, v := range c.Addresses {
		if v != addrToRemove {
			addrs = append(addrs, v)
		}
	}

	c.Addresses = addrs
	return c.Validate()
}
