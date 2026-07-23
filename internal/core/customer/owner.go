package customer

import "slices"

type OwnerProfile struct {
	SignatureActive bool     `json:"signature_active" validate:"bool"`
	StoresID        []string `json:"stores" validate:"required,uuid"`
}

func NewOwnerProfile() OwnerProfile {
	return OwnerProfile{
		SignatureActive: false,
		StoresID:        []string{},
	}
}

func (op *OwnerProfile) IsOwnerOf(storeID string) bool {
	return slices.Contains(op.StoresID, storeID)
}
