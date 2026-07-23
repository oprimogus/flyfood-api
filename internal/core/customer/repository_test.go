package customer

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/oprimogus/flyfood-api/internal/config"
	"github.com/oprimogus/flyfood-api/internal/core/address"
	"github.com/oprimogus/flyfood-api/internal/infra/database"
	"github.com/oprimogus/flyfood-api/pkg/testcontainers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func makeCustomer(id uuid.UUID, externalID, name, lastName, email, cpf string) Customer {
	return Customer{
		ID:         id,
		ExternalID: externalID,
		Name:       name,
		LastName:   lastName,
		CPF:        cpf,
		Phone:      "+5513997590579",
		Email:      email,
	}
}

func makeAddress(id uuid.UUID, name string, isDefault bool) address.Address {
	return address.Address{
		ID:           id,
		Name:         name,
		Default:      isDefault,
		AddressLine1: "Rua das Flores, 123",
		AddressLine2: "Apto 4",
		Neighborhood: "Centro",
		City:         "Santos",
		State:        "SP",
		PostalCode:   "11010000",
		Country:      "BR",
		Latitude:     -23.9608,
		Longitude:    -46.3336,
	}
}

// ── suite setup ───────────────────────────────────────────────────────────────

type CustomerRepositoryTestSuite struct {
	suite.Suite
	ctx          context.Context
	mockPostgres *testcontainers.PostgresContainer
	repo         Repository
}

func (s *CustomerRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()

	pgCfg := testcontainers.PostgresConfig{
		DatabaseName: "flyfood",
		Username:     "test",
		Password:     "test",
	}

	mockDB, err := testcontainers.MakePostgres(s.ctx, pgCfg)
	s.Require().NoError(err, "failed to start postgres container")
	s.mockPostgres = mockDB

	cfg := config.Get()
	cfg.Postgres.Host = "localhost"
	cfg.Postgres.Port = mockDB.Port
	cfg.Postgres.UserName = pgCfg.Username
	cfg.Postgres.Password = pgCfg.Password
	cfg.Postgres.DatabaseName = pgCfg.DatabaseName

	db, err := database.GetPostgres(s.ctx)
	s.Require().NoError(err, "failed to connect to postgres")

	err = db.Migrate(s.ctx)
	s.Require().NoError(err, "failed to migrate database")

	s.repo = NewRepository(db)
}

func (s *CustomerRepositoryTestSuite) TearDownSuite() {
	s.mockPostgres.Kill(s.ctx)
}

func TestCustomerRepositorySuite(t *testing.T) {
	suite.Run(t, new(CustomerRepositoryTestSuite))
}

// ── TestSave ──────────────────────────────────────────────────────────────────

func (s *CustomerRepositoryTestSuite) TestSave() {
	testCases := []struct {
		name        string
		input       Customer
		expectError bool
	}{
		{
			name:        "should save customer successfully",
			input:       makeCustomer(uuid.New(), "ext_save_001", "Arthur", "Morgan", "arthur@rdr.com", "52024227090"),
			expectError: false,
		},
		{
			name:        "should fail on duplicate external_id",
			input:       makeCustomer(uuid.New(), "ext_save_001", "John", "Marston", "john@rdr.com", "52024227091"),
			expectError: true,
		},
		{
			name:        "should fail on duplicate email",
			input:       makeCustomer(uuid.New(), "ext_save_002", "Dutch", "VanDerLinde", "arthur@rdr.com", "52024227092"),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			err := s.repo.Save(s.ctx, &tc.input)
			if tc.expectError {
				s.Error(err)
				return
			}
			s.NoError(err)

			found, err := s.repo.FindByID(s.ctx, tc.input.ID)
			s.NoError(err)
			s.Equal(tc.input.ID, found.ID)
			s.Equal(tc.input.ExternalID, found.ExternalID)
			s.Equal(tc.input.Email, found.Email)
		})
	}
}

// ── TestFindByID ──────────────────────────────────────────────────────────────

func (s *CustomerRepositoryTestSuite) TestFindByID() {
	existing := makeCustomer(uuid.New(), "ext_findid_001", "Lenny", "Summers", "lenny@rdr.com", "52024227097")
	s.Require().NoError(s.repo.Save(s.ctx, &existing))

	testCases := []struct {
		name        string
		id          uuid.UUID
		expectError bool
		expected    *Customer
	}{
		{
			name:        "should find customer by id",
			id:          existing.ID,
			expectError: false,
			expected:    &existing,
		},
		{
			name:        "should return error when customer does not exist",
			id:          uuid.New(),
			expectError: true,
			expected:    nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			result, err := s.repo.FindByID(s.ctx, tc.id)
			if tc.expectError {
				s.Error(err)
				s.Nil(result)
				return
			}
			s.NoError(err)
			s.Equal(tc.expected.ID, result.ID)
			s.Equal(tc.expected.ExternalID, result.ExternalID)
			s.Equal(tc.expected.Email, result.Email)
		})
	}
}

// ── TestFindByExternalID ──────────────────────────────────────────────────────

func (s *CustomerRepositoryTestSuite) TestFindByExternalID() {
	existing := makeCustomer(uuid.New(), "ext_findextid_001", "Hosea", "Matthews", "hosea@rdr.com", "52024227093")
	s.Require().NoError(s.repo.Save(s.ctx, &existing))

	testCases := []struct {
		name        string
		externalID  string
		expectError bool
	}{
		{
			name:        "should find customer by external id",
			externalID:  existing.ExternalID,
			expectError: false,
		},
		{
			name:        "should return error when external id does not exist",
			externalID:  "nonexistent_ext_id",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			result, err := s.repo.FindByExternalID(s.ctx, tc.externalID)
			if tc.expectError {
				s.Error(err)
				s.Nil(result)
				return
			}
			s.NoError(err)
			s.Equal(existing.ID, result.ID)
			s.Equal(existing.ExternalID, result.ExternalID)
		})
	}
}

// ── TestSaveAddress ───────────────────────────────────────────────────────────

func (s *CustomerRepositoryTestSuite) TestSaveAddress() {
	customer := makeCustomer(uuid.New(), "ext_saveaddr_001", "Bill", "Williamson", "bill@rdr.com", "52024227098")
	s.Require().NoError(s.repo.Save(s.ctx, &customer))

	testCases := []struct {
		name        string
		customerID  uuid.UUID
		addr        address.Address
		expectError bool
	}{
		{
			name:        "should save address successfully",
			customerID:  customer.ID,
			addr:        makeAddress(uuid.New(), "Casa", true),
			expectError: false,
		},
		{
			name:        "should save second address successfully",
			customerID:  customer.ID,
			addr:        makeAddress(uuid.New(), "Trabalho", false),
			expectError: false,
		},
		{
			name:        "should fail with nonexistent customer",
			customerID:  uuid.New(),
			addr:        makeAddress(uuid.New(), "Outro", false),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			err := s.repo.SaveAddress(s.ctx, tc.customerID, tc.addr)
			if tc.expectError {
				s.Error(err)
				return
			}
			s.NoError(err)

			addresses, err := s.repo.FindAddressesByExternalCustomerID(s.ctx, customer.ExternalID)
			s.NoError(err)
			s.NotEmpty(addresses)

			found := false
			for _, a := range addresses {
				if a.ID == tc.addr.ID {
					found = true
					s.Equal(tc.addr.Name, a.Name)
					s.Equal(tc.addr.City, a.City)
					s.Equal(tc.addr.Default, a.Default)
					break
				}
			}
			s.True(found, "saved address should be found in customer's addresses")
		})
	}
}

// ── TestDeleteAddress ─────────────────────────────────────────────────────────

func (s *CustomerRepositoryTestSuite) TestDeleteAddress() {
	customer := makeCustomer(uuid.New(), "ext_deladdr_001", "Micah", "Bell", "micah@rdr.com", "52024227094")
	s.Require().NoError(s.repo.Save(s.ctx, &customer))

	addr := makeAddress(uuid.New(), "Esconderijo", true)
	s.Require().NoError(s.repo.SaveAddress(s.ctx, customer.ID, addr))

	testCases := []struct {
		name        string
		customerID  uuid.UUID
		addressID   uuid.UUID
		expectError bool
	}{
		{
			name:        "should fail when address does not belong to customer",
			customerID:  uuid.New(),
			addressID:   addr.ID,
			expectError: true,
		},
		{
			name:        "should delete address successfully",
			customerID:  customer.ID,
			addressID:   addr.ID,
			expectError: false,
		},
		{
			name:        "should fail when address is already deleted",
			customerID:  customer.ID,
			addressID:   addr.ID,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			err := s.repo.DeleteAddress(s.ctx, tc.customerID, tc.addressID)
			if tc.expectError {
				s.Error(err)
				return
			}
			s.NoError(err)

			addresses, err := s.repo.FindAddressesByExternalCustomerID(s.ctx, customer.ExternalID)
			s.NoError(err)

			for _, a := range addresses {
				s.NotEqual(addr.ID, a.ID, "deleted address should not appear in results")
			}
		})
	}
}

// ── TestFindAddressesByExternalCustomerID ─────────────────────────────────────

func (s *CustomerRepositoryTestSuite) TestFindAddressesByExternalCustomerID() {
	customer := makeCustomer(uuid.New(), "ext_findaddr_001", "Charles", "Smith", "charles@rdr.com", "52024227095")
	s.Require().NoError(s.repo.Save(s.ctx, &customer))

	addr1 := makeAddress(uuid.New(), "Casa", true)
	addr2 := makeAddress(uuid.New(), "Trabalho", false)
	s.Require().NoError(s.repo.SaveAddress(s.ctx, customer.ID, addr1))
	s.Require().NoError(s.repo.SaveAddress(s.ctx, customer.ID, addr2))

	testCases := []struct {
		name            string
		externalID      string
		expectError     bool
		expectedCount   int
	}{
		{
			name:          "should return all addresses for customer",
			externalID:    customer.ExternalID,
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:          "should return empty list for customer with no addresses",
			externalID:    "ext_findaddr_no_addr",
			expectError:   false,
			expectedCount: 0,
		},
	}

	// Setup: customer sem endereços
	noAddrCustomer := makeCustomer(uuid.New(), "ext_findaddr_no_addr", "Javier", "Escuella", "javier@rdr.com", "52024227096")
	s.Require().NoError(s.repo.Save(s.ctx, &noAddrCustomer))

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			result, err := s.repo.FindAddressesByExternalCustomerID(s.ctx, tc.externalID)
			if tc.expectError {
				s.Error(err)
				return
			}
			s.NoError(err)
			assert.Len(s.T(), result, tc.expectedCount)
		})
	}
}