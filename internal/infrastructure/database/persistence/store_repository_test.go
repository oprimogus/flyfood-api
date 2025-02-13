//go:build integration

package persistence

import (
	"context"
	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/oprimogus/cardapiogo/internal/core/customer"
	"github.com/oprimogus/cardapiogo/internal/core/owner"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	"github.com/oprimogus/cardapiogo/internal/core/store/product"
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
	ownerRepo    OwnerRepository
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
	s.ownerRepo = NewOwnerRepository(db)
}

func (s *StoreRepositoryTestSuite) TearDownSuite() {
	ctx := context.Background()
	s.mockPostgres.Kill(ctx)
}

func TestStoreRepositorySuite(t *testing.T) {
	suite.Run(t, new(StoreRepositoryTestSuite))
}

func (s *StoreRepositoryTestSuite) TestSaveStore() {
	ctx := context.Background()
	mockCustomer := customer.Customer{
		ID:       "356278453827428",
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

	mockOwner := owner.NewOwner(mockCustomer.ID)

	err = s.ownerRepo.SaveOwner(ctx, mockOwner)
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
		"23754700000177",
		"Store test b",
		"store from test b",
		"+5513997590531",
		addr,
		store.Restaurant)
	assert.NoError(s.T(), err)

	mockStoreWithBusinessHour, err := store.NewStore(
		mockOwner.ID,
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

	err = mockStoreWithBusinessHour.AddNewBusinessHour(store.BusinessHours{
		WeekDay:     2,
		OpeningTime: "12:00",
		ClosingTime: "21:00",
	})
	assert.NoError(s.T(), err)

	mockStoreWithPaymentMethods, err := store.NewStore(
		mockOwner.ID,
		"43544583000124",
		"Store test 3",
		"store from test 3",
		"+5513997590572",
		addr,
		store.Restaurant)
	assert.NoError(s.T(), err)

	err = mockStoreWithPaymentMethods.AddPaymentMethod(store.Pix)
	assert.NoError(s.T(), err)

	err = mockStoreWithPaymentMethods.AddPaymentMethod(store.Bitcoin)
	assert.NoError(s.T(), err)

	tests := []struct {
		name   string
		store  store.Store
		expect error
	}{
		{
			name:   "Should save a new store",
			store:  mockStore,
			expect: nil,
		},
		{
			name:   "Should save a new store with businessHour",
			store:  mockStoreWithBusinessHour,
			expect: nil,
		},
		{
			name:   "Should save a new store with payment methods",
			store:  mockStoreWithPaymentMethods,
			expect: nil,
		},
	}

	for _, tt := range tests {
		err := s.storeRepo.SaveStore(ctx, tt.store)
		assert.Equal(s.T(), tt.expect, err, tt.name)

		savedStore, err := s.storeRepo.FindStoreByID(ctx, tt.store.ID)
		assert.NoError(s.T(), err, tt.name)
		assert.EqualValues(s.T(), tt.store, savedStore, tt.name)
	}
}

func (s *StoreRepositoryTestSuite) TestSaveProduct() {
	ctx := context.Background()
	mockCustomer := customer.Customer{
		ID:       "2575475247247542254",
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

	// mock customer
	err := s.customerRepo.Save(ctx, &mockCustomer)
	assert.NoError(s.T(), err)

	mockOwner := owner.NewOwner(mockCustomer.ID)

	// become mock customer an owner
	err = s.ownerRepo.SaveOwner(context.Background(), mockOwner)
	assert.NoError(s.T(), err)

	// store address
	addr := address.Address{
		AddressLine1: "rua 1",
		AddressLine2: "879",
		Neighborhood: "test",
		City:         "test",
		State:        "test",
		PostalCode:   "11490-135",
		Country:      "Brasil",
	}

	// owner create a store
	mockStore, err := store.NewStore(
		mockCustomer.ID,
		"63432495000148",
		"Store test",
		"store from test",
		"+5513997590579",
		addr,
		store.Restaurant)
	assert.NoError(s.T(), err)

	//Save store
	err = s.storeRepo.SaveStore(ctx, mockStore)
	assert.NoError(s.T(), err)

	// Create product
	p, err := product.NewProduct(
		mockStore.ID,
		"product one", "product tag", "your desc", "P001", 2500, product.Food)
	assert.NoError(s.T(), err)

	// finally save product
	err = s.storeRepo.SaveProduct(ctx, p)
	assert.NoError(s.T(), err)

	//finally do test
	persistentProduct, err := s.storeRepo.FindStoreProductByID(ctx, p.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), p, persistentProduct)

	st, err := s.storeRepo.FindStoreByID(ctx, mockStore.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockStore, st)

}
