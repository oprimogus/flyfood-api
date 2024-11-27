//go:build integration

package persistence

import (
	"context"
	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/oprimogus/cardapiogo/internal/core/customer"
	postgresDB "github.com/oprimogus/cardapiogo/internal/infrastructure/database/postgres"
	"github.com/oprimogus/cardapiogo/test/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CustomerRepositoryTestSuite struct {
	suite.Suite
	mockPostgres *integration.Container
	customerRepo CustomerRepository
}

func (s *CustomerRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	mockDB, err := integration.MakePostgres(ctx)
	if err != nil {
		assert.Error(s.T(), err)
	}
	s.mockPostgres = mockDB

	db := postgresDB.GetTestInstance(mockDB.Port)

	s.customerRepo = NewCustomerRepository(db)
}

func (s *CustomerRepositoryTestSuite) TearDownSuite() {
	ctx := context.Background()
	s.mockPostgres.Kill(ctx)
}

func TestCustomerRepositorySuite(t *testing.T) {
	suite.Run(t, new(CustomerRepositoryTestSuite))
}

func (s *CustomerRepositoryTestSuite) TestFindByID() {
	testCases := []struct {
		name           string
		customer       *customer.Customer
		expectedError  error
		expectedResult *customer.Customer
	}{
		{
			name: "Should save with success customer without address",
			customer: &customer.Customer{
				ID:        4687487487864,
				Name:      "John",
				LastName:  "Marston",
				CPF:       "fake cpf",
				Phone:     "fake phone",
				Email:     "fake email",
				Addresses: []address.Address{},
			},
			expectedError: nil,
			expectedResult: &customer.Customer{
				ID:        4687487487864,
				Name:      "John",
				LastName:  "Marston",
				CPF:       "fake cpf",
				Phone:     "fake phone",
				Email:     "fake email",
				Addresses: []address.Address{},
			},
		},
		{
			name: "Should save with success customer with address",
			customer: &customer.Customer{
				ID:       4687487487864,
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
			expectedError: nil,
			expectedResult: &customer.Customer{
				ID:       4687487487864,
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
		},
	}

	for _, test := range testCases {
		ctx := context.Background()
		err := s.customerRepo.Save(ctx, test.customer)
		assert.NoError(s.T(), err, test.name)
		actual, err := s.customerRepo.FindByID(ctx, test.customer.ID)
		assert.NoError(s.T(), err, test.name)
		assert.Equal(s.T(), test.expectedResult, actual, test.name)
	}
}

func (s *CustomerRepositoryTestSuite) TestSave() {
	testCases := []struct {
		name             string
		savedCustomer    *customer.Customer
		inputCustomer    *customer.Customer
		expectedError    error
		expectedCustomer *customer.Customer
	}{
		{
			name: "Should save with success customer with new address and remove older",
			savedCustomer: &customer.Customer{
				ID:       4687487487864,
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
			inputCustomer: &customer.Customer{
				ID:       4687487487864,
				Name:     "John",
				LastName: "Marston",
				CPF:      "fake cpf",
				Phone:    "fake phone",
				Email:    "fake email",
				Addresses: []address.Address{
					{
						Name:         "test 2 address",
						AddressLine1: "test 2 Street",
						AddressLine2: "test 2 Number",
						Neighborhood: "test neighborhood",
						City:         "test city",
						State:        "test state",
						PostalCode:   "test postal code",
						Country:      "teste country",
					},
					{
						Name:         "test 3 address",
						AddressLine1: "test 3 Street",
						AddressLine2: "test 3 Number",
						Neighborhood: "test neighborhood",
						City:         "test city",
						State:        "test state",
						PostalCode:   "test postal code",
						Country:      "teste country",
					},
				},
			},
			expectedError: nil,
			expectedCustomer: &customer.Customer{
				ID:       4687487487864,
				Name:     "John",
				LastName: "Marston",
				CPF:      "fake cpf",
				Phone:    "fake phone",
				Email:    "fake email",
				Addresses: []address.Address{
					{
						Name:         "test 2 address",
						AddressLine1: "test 2 Street",
						AddressLine2: "test 2 Number",
						Neighborhood: "test neighborhood",
						City:         "test city",
						State:        "test state",
						PostalCode:   "test postal code",
						Country:      "teste country",
					},
					{
						Name:         "test 3 address",
						AddressLine1: "test 3 Street",
						AddressLine2: "test 3 Number",
						Neighborhood: "test neighborhood",
						City:         "test city",
						State:        "test state",
						PostalCode:   "test postal code",
						Country:      "teste country",
					},
				},
			},
		},
		{
			name: "Should save with success customer with new address and keep older",
			savedCustomer: &customer.Customer{
				ID:       4687487487864,
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
			inputCustomer: &customer.Customer{
				ID:       4687487487864,
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
					{
						Name:         "test 2 address",
						AddressLine1: "test 2 Street",
						AddressLine2: "test 2 Number",
						Neighborhood: "test neighborhood",
						City:         "test city",
						State:        "test state",
						PostalCode:   "test postal code",
						Country:      "teste country",
					},
					{
						Name:         "test 3 address",
						AddressLine1: "test 3 Street",
						AddressLine2: "test 3 Number",
						Neighborhood: "test neighborhood",
						City:         "test city",
						State:        "test state",
						PostalCode:   "test postal code",
						Country:      "teste country",
					},
				},
			},
			expectedError: nil,
			expectedCustomer: &customer.Customer{
				ID:       4687487487864,
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
					{
						Name:         "test 2 address",
						AddressLine1: "test 2 Street",
						AddressLine2: "test 2 Number",
						Neighborhood: "test neighborhood",
						City:         "test city",
						State:        "test state",
						PostalCode:   "test postal code",
						Country:      "teste country",
					},
					{
						Name:         "test 3 address",
						AddressLine1: "test 3 Street",
						AddressLine2: "test 3 Number",
						Neighborhood: "test neighborhood",
						City:         "test city",
						State:        "test state",
						PostalCode:   "test postal code",
						Country:      "teste country",
					},
				},
			},
		},
	}

	for _, test := range testCases {
		ctx := context.Background()
		err := s.customerRepo.Save(ctx, test.savedCustomer)
		assert.NoError(s.T(), err, test.name)

		err = s.customerRepo.Save(ctx, test.inputCustomer)
		assert.NoError(s.T(), err, test.name)

		actual, err := s.customerRepo.FindByID(ctx, test.savedCustomer.ID)
		assert.NoError(s.T(), err, test.name)
		assert.Equal(s.T(), test.expectedCustomer, actual, test.name)
	}
}
