package store

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/oprimogus/cardapiogo/internal/xvalidator"
)

type Owner struct {
	ID              string   `json:"id" validate:"required" example:"245727247525742"`
	SignatureActive bool     `json:"signature_active" validate:"bool"`
	StoresID        []string `json:"stores" validate:"required,uuid"`
}

func NewOwner(customerID string) Owner {
	return Owner{
		ID:              customerID,
		SignatureActive: false,
		StoresID:        []string{},
	}
}

func (o *Owner) Validate() error {
	return xvalidator.GetPtInstance().Validate(o)
}

func (o *Owner) NewStore(cnpj, name, description,
	phone string, address address.Address, types Type) (Store, error) {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		return Store{}, fmt.Errorf("fail on create store id: %w", err)
	}
	newStore := Store{
		ID:             uuidV7.String(),
		OwnerID:        o.ID,
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
		Products:       []string{},
	}
	o.StoresID = append(o.StoresID, newStore.ID)
	if err := newStore.Validate(); err != nil {
		return Store{}, err
	}

	return newStore, nil
}
