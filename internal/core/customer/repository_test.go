// go:build integration

package customer

import (
	"context"
	"github.com/oprimogus/flyfood-api/internal/core/address"
	postgresDB "github.com/oprimogus/flyfood-api/internal/infrastructure/database/postgres"
	"github.com/oprimogus/flyfood-api/test/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CustomerRepositoryTestSuite struct {
	suite.Suite
	mockPostgres *integration.Container
	customerRepo Repository
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
		customer       *Customer
		expectedError  error
		expectedResult *Customer
	}{
		{
			name: "Should save with success customer without address",
			customer: &Customer{
				ID:        "4687487487864",
				Name:      "John",
				LastName:  "Marston",
				CPF:       "fake cpf",
				Phone:     "fake phone",
				Email:     "fake email",
				Addresses: []address.Address{},
			},
			expectedError: nil,
			expectedResult: &Customer{
				ID:        "4687487487864",
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
			customer: &Customer{
				ID:       "4687487487864",
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
			expectedResult: &Customer{
				ID:       "4687487487864",
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
		savedCustomer    *Customer
		inputCustomer    *Customer
		expectedError    error
		expectedCustomer *Customer
	}{
		{
			name: "Should save with success customer with new address and remove older",
			savedCustomer: &Customer{
				ID:       "4687487487864",
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
			inputCustomer: &Customer{
				ID:       "4687487487864",
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
			expectedCustomer: &Customer{
				ID:       "4687487487864",
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
			savedCustomer: &Customer{
				ID:       "4687487487864",
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
			inputCustomer: &Customer{
				ID:       "4687487487864",
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
			expectedCustomer: &Customer{
				ID:       "4687487487864",
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
