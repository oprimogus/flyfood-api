//go:build integration

package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/oprimogus/cardapiogo/internal/core/customer"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	postgresDB "github.com/oprimogus/cardapiogo/internal/infrastructure/database/postgres"
	"github.com/oprimogus/cardapiogo/test/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type StoreRepositoryTestSuite struct {
	suite.Suite
	mockPostgres *integration.Container
	customerRepo CustomerRepository
	storeRepo    StoreRepository
}

func (s *StoreRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	mockDB, err := integration.MakePostgres(ctx)
	if err != nil {
		assert.Error(s.T(), err)
	}
	s.mockPostgres = mockDB

	db := postgresDB.GetTestInstance(mockDB.Port)
	s.customerRepo = NewCustomerRepository(db)
	s.storeRepo = NewStoreRepository(db)
}

func (s *StoreRepositoryTestSuite) TearDownSuite() {
	ctx := context.Background()
	s.mockPostgres.Kill(ctx)
}

func TestStoreRepositorySuite(t *testing.T) {
	suite.Run(t, new(StoreRepositoryTestSuite))
}

func (s *StoreRepositoryTestSuite) TestFindOwnerByID() {
	ctx := context.Background()
	mockCustomer := customer.Customer{
		ID:       "0193480f-9ad2-7dbc-bae0-ca1c9fae1246",
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

	err := s.customerRepo.Save(ctx, &mockCustomer)
	assert.NoError(s.T(), err)

	err = s.storeRepo.BecomeOwner(ctx, mockCustomer.ID)
	assert.NoError(s.T(), err)

	ow, err := s.storeRepo.FindOwnerByID(ctx, mockCustomer.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockCustomer.ID, ow.ID)
	assert.Equal(s.T(), false, ow.SignatureActive)
}

func (s *StoreRepositoryTestSuite) TestIsOwner() {

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
				ID:       "019322cc-de2d-7064-9516-46e2c1db90be",
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
				ID:       "019341e4-76b2-7a70-856c-4f7a7188c99f",
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
			err = s.storeRepo.BecomeOwner(ctx, tt.customer.ID)
			assert.NoError(s.T(), err)
		}

		isOwner, err := s.storeRepo.IsOwner(ctx, tt.customer.ID)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), tt.expectBool, isOwner, tt.name)
	}

}

func (s *StoreRepositoryTestSuite) TestSave() {
	ctx := context.Background()
	mockCustomer := customer.Customer{
		ID:       "01934806-981e-7e9d-adab-2086e95ee82e",
		Name:     "John",
		LastName: "Marston",
		CPF:      "fake cpf 21313",
		Phone:    "fake phone 134134",
		Email:    "fake email 431431413",
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

	err := s.customerRepo.Save(ctx, &mockCustomer)
	assert.NoError(s.T(), err)

	err = s.storeRepo.BecomeOwner(ctx, mockCustomer.ID)
	assert.NoError(s.T(), err)

	mockOwner := store.NewOwner(mockCustomer.ID)

	addr := address.Address{
		AddressLine1: "rua 1",
		AddressLine2: "879",
		Neighborhood: "test",
		City:         "test",
		State:        "test",
		PostalCode:   "11490-135",
		Country:      "Brasil",
	}

	mockStore, err := mockOwner.NewStore(
		"63432495000148",
		"Store test",
		"store from test",
		"+5513997590579",
		addr,
		store.Restaurant)
	assert.NoError(s.T(), err)

	mockStoreWithBusinessHour, err := mockOwner.NewStore(
		"65815550000104",
		"Store test 2",
		"store from test 2",
		"+5513997590571",
		addr,
		store.Restaurant)
	assert.NoError(s.T(), err)

	err = mockStoreWithBusinessHour.AddNewBusinessHour(store.BusinessHours{
		WeekDay:     1,
		OpeningTime: "12:00",
		ClosingTime: "21:00",
	})
	assert.NoError(s.T(), err)

	mockStoreWithPaymentMethods, err := mockOwner.NewStore(
		"43544583000124",
		"Store test 3",
		"store from test 3",
		"+5513997590572",
		addr,
		store.Restaurant)
	assert.NoError(s.T(), err)

	mockStoreWithPaymentMethods.AddPaymentMethod(store.Bitcoin)

	tests := []struct {
		name   string
		store  *store.Store
		expect error
	}{
		{
			name:   "Should save a new store",
			store:  &mockStore,
			expect: nil,
		},
		{
			name:   "Should save a new store with businessHour",
			store:  &mockStoreWithBusinessHour,
			expect: nil,
		},
		{
			name:   "Should save a new store with payment methods",
			store:  &mockStoreWithPaymentMethods,
			expect: nil,
		},
	}

	for _, tt := range tests {
		err := s.storeRepo.Save(ctx, tt.store)
		assert.Equal(s.T(), tt.expect, err, tt.name)

		savedStore, err := s.storeRepo.FindStoreByID(ctx, tt.store.ID)
		assert.NoError(s.T(), err, tt.name)
		assert.Equal(s.T(), tt.store, savedStore, tt.name)
		itemJson, _ := json.Marshal(savedStore)
		fmt.Println(string(itemJson))
	}
}
