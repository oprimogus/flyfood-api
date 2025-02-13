package owner

import (
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
