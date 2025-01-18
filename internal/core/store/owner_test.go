package store

import (
	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOwner_NewStore(t *testing.T) {

	mockOwner := NewOwner("45678948674")

	addr := address.Address{
		Name:         "Store name location",
		AddressLine1: "rua 1",
		AddressLine2: "879",
		Neighborhood: "test",
		City:         "test",
		State:        "test",
		PostalCode:   "11490-135",
		Country:      "Brasil",
	}

	store, err := mockOwner.NewStore(
		"63432495000148",
		"Store test",
		"store from test",
		"+5513997590579",
		addr,
		Restaurant)

	assert.NoError(t, err)
	assert.NotNil(t, store)
	assert.Equal(t, "45678948674", store.OwnerID)
	assert.Equal(t, DefaultScore, store.Score)
	assert.Equal(t, false, store.Active)
	assert.Equal(t, false, store.IsOpen)
	assert.Equal(t, store.Address, addr)
}
