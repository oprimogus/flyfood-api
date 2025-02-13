//go:build integration

package persistence

import (
	"context"
	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/oprimogus/cardapiogo/internal/core/customer"
	"github.com/oprimogus/cardapiogo/internal/core/owner"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	postgresDB "github.com/oprimogus/cardapiogo/internal/infrastructure/database/postgres"
	"github.com/oprimogus/cardapiogo/test/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type OwnerRepositoryTestSuite struct {
	suite.Suite
	mockPostgres *integration.Container
	repository   OwnerRepository
	customerRepo CustomerRepository
	storeRepo    StoreRepository
}

func (s *OwnerRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	mockDB, err := integration.MakePostgres(ctx)
	if err != nil {
		assert.Error(s.T(), err)
	}
	s.mockPostgres = mockDB

	db := postgresDB.GetTestInstance(mockDB.Port)
	s.repository = NewOwnerRepository(db)
	s.customerRepo = NewCustomerRepository(db)
	s.storeRepo = NewStoreRepository(db)
}

func (s *OwnerRepositoryTestSuite) TearDownSuite() {
	ctx := context.Background()
	s.mockPostgres.Kill(ctx)
}

func TestOwnerRepositorySuite(t *testing.T) {
	suite.Run(t, new(OwnerRepositoryTestSuite))
}

func (s *OwnerRepositoryTestSuite) TestFindOwnerByID() {
	ctx := context.Background()
	mockCustomer := customer.Customer{
		ID:       "356278453827428",
		Name:     "John",
		LastName: "Marston",
		CPF:      "fake cpf 345t643t",
		Phone:    "fake phone fefefge",
		Email:    "fake email 42376762",
		Addresses: []address.Address{
			{
				Name:         "test address",
				AddressLine1: "test Street",
				AddressLine2: "test Number",
				Neighborhood: "test neighborhood",
				City:         "test city",
				State:        "test state",
				PostalCode:   "test postal code",
				Country:      "teste country",
			},
		},
	}

	mockOwner := owner.NewOwner(mockCustomer.ID)

	err := s.customerRepo.Save(ctx, &mockCustomer)
	assert.NoError(s.T(), err)

	err = s.repository.SaveOwner(ctx, mockOwner)
	assert.NoError(s.T(), err)

	ow, err := s.repository.FindOwnerByID(ctx, mockCustomer.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockCustomer.ID, ow.ID)
	assert.Equal(s.T(), false, ow.SignatureActive)
}

func (s *OwnerRepositoryTestSuite) TestIsOwner() {

	testCases := []struct {
		name        string
		customer    customer.Customer
		isOwner     bool
		expectBool  bool
		expectError error
	}{
		{
			name: "should return true",
			customer: customer.Customer{
				ID:       "356278453827428",
				Name:     "John",
				LastName: "Marston",
				CPF:      "fake cpf",
				Phone:    "fake phone",
				Email:    "fake email",
				Addresses: []address.Address{
					{
						Name:         "test address",
						AddressLine1: "test Street",
						AddressLine2: "test Number",
						Neighborhood: "test neighborhood",
						City:         "test city",
						State:        "test state",
						PostalCode:   "test postal code",
						Country:      "teste country",
					},
				},
			},
			isOwner:     true,
			expectBool:  true,
			expectError: nil,
		},
		{
			name: "should return false",
			customer: customer.Customer{
				ID:       "356278453827429",
				Name:     "John",
				LastName: "Marstones",
				CPF:      "fake cpf a",
				Phone:    "fake phone a",
				Email:    "fake email a",
				Addresses: []address.Address{
					{
						Name:         "test address",
						AddressLine1: "test Street",
						AddressLine2: "test Number",
						Neighborhood: "test neighborhood",
						City:         "test city",
						State:        "test state",
						PostalCode:   "test postal code",
						Country:      "teste country",
					},
				},
			},
			isOwner:     false,
			expectBool:  false,
			expectError: nil,
		},
	}

	for _, tt := range testCases {
		ctx := context.Background()

		err := s.customerRepo.Save(ctx, &tt.customer)
		assert.NoError(s.T(), err)

		if tt.isOwner {
			mockOwner := owner.NewOwner(tt.customer.ID)
			err = s.repository.SaveOwner(ctx, mockOwner)
			assert.NoError(s.T(), err)
		}

		isOwner, err := s.repository.IsOwner(ctx, tt.customer.ID)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), tt.expectBool, isOwner, tt.name)
	}

}

func (s *OwnerRepositoryTestSuite) TestIsOwnerOf() {

	testCases := []struct {
		name        string
		customer    customer.Customer
		isOwner     bool
		expectBool  bool
		expectError error
	}{
		{
			name: "should return true",
			customer: customer.Customer{
				ID:       "356278453827428",
				Name:     "John",
				LastName: "Marston",
				CPF:      "fake cpf",
				Phone:    "fake phone",
				Email:    "fake email",
				Addresses: []address.Address{
					{
						Name:         "test address",
						AddressLine1: "test Street",
						AddressLine2: "test Number",
						Neighborhood: "test neighborhood",
						City:         "test city",
						State:        "test state",
						PostalCode:   "test postal code",
						Country:      "teste country",
					},
				},
			},
			isOwner:     true,
			expectBool:  true,
			expectError: nil,
		},
		{
			name: "should return false",
			customer: customer.Customer{
				ID:       "356278453827429",
				Name:     "John",
				LastName: "Marstones",
				CPF:      "fake cpf a",
				Phone:    "fake phone a",
				Email:    "fake email a",
				Addresses: []address.Address{
					{
						Name:         "test address",
						AddressLine1: "test Street",
						AddressLine2: "test Number",
						Neighborhood: "test neighborhood",
						City:         "test city",
						State:        "test state",
						PostalCode:   "test postal code",
						Country:      "teste country",
					},
				},
			},
			isOwner:     false,
			expectBool:  false,
			expectError: nil,
		},
	}

	ctxDefault := context.Background()

	defaultCustomer := customer.Customer{
		ID:       "356278453827428",
		Name:     "John",
		LastName: "Marston",
		CPF:      "fake cpf",
		Phone:    "fake phone",
		Email:    "fake email",
		Addresses: []address.Address{
			{
				Name:         "test address",
				AddressLine1: "test Street",
				AddressLine2: "test Number",
				Neighborhood: "test neighborhood",
				City:         "test city",
				State:        "test state",
				PostalCode:   "test postal code",
				Country:      "teste country",
			},
		},
	}
	err := s.customerRepo.Save(ctxDefault, &defaultCustomer)
	assert.NoError(s.T(), err)

	defaultOwner := owner.NewOwner(defaultCustomer.ID)
	err = s.repository.SaveOwner(ctxDefault, defaultOwner)
	assert.NoError(s.T(), err)

	addr := address.Address{
		AddressLine1: "rua 1",
		AddressLine2: "879",
		Neighborhood: "test",
		City:         "test",
		State:        "test",
		PostalCode:   "11490-135",
		Country:      "Brasil",
	}

	defaultStore, err := store.NewStore(
		defaultOwner.ID,
		"40990221000179",
		"Store test r",
		"store from test r",
		"+5513997590590",
		addr,
		store.Restaurant)
	assert.NoError(s.T(), err)

	err = s.storeRepo.SaveStore(ctxDefault, defaultStore)

	for _, tt := range testCases {
		ctx := context.Background()

		err := s.customerRepo.Save(ctx, &tt.customer)
		assert.NoError(s.T(), err)

		if tt.isOwner {
			mockOwner := owner.NewOwner(tt.customer.ID)
			err = s.repository.SaveOwner(ctx, mockOwner)
			assert.NoError(s.T(), err)

			addr := address.Address{
				AddressLine1: "rua 1",
				AddressLine2: "879",
				Neighborhood: "test",
				City:         "test",
				State:        "test",
				PostalCode:   "11490-135",
				Country:      "Brasil",
			}

			mockStore, err := store.NewStore(
				mockOwner.ID,
				"63432495000148",
				"Store test",
				"store from test",
				"+5513997590579",
				addr,
				store.Restaurant)
			assert.NoError(s.T(), err)

			err = s.storeRepo.SaveStore(ctx, mockStore)
			assert.NoError(s.T(), err)

			isOwner, err := s.repository.IsOwnerOf(ctx, tt.customer.ID, mockStore.ID)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), tt.expectBool, isOwner, tt.name)
		} else {
			isOwner, err := s.repository.IsOwnerOf(ctx, tt.customer.ID, defaultStore.ID)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), tt.expectBool, isOwner, tt.name)
		}
	}
}
