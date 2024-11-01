package persistence_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	"github.com/oprimogus/cardapiogo/internal/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/database/sqlc"
	"github.com/oprimogus/cardapiogo/internal/persistence"
	"github.com/oprimogus/cardapiogo/internal/services/adapter"
	"github.com/oprimogus/cardapiogo/test/integration"
)

type StoreRepositorySuite struct {
	suite.Suite
	ctx               context.Context
	repository        *persistence.StoreRepository
	pgContainer       *integration.Container
	keycloakContainer *integration.Container
}

func (s *StoreRepositorySuite) SetupSuite() {
	ctx := context.Background()
	s.ctx = ctx
	mockPostgres, err := integration.MakePostgres(s.ctx)
	if err != nil {
		panic("fail on initialize postgres with testContainers")
	}
	s.pgContainer = mockPostgres

	keycloak, err := integration.MakeKeycloak(s.ctx)
	if err != nil {
		panic(err)
	}
	s.keycloakContainer = keycloak

	db := postgres.GetInstance()
	querier := sqlc.New(db.GetDB())
	serviceFactory := adapter.NewServiceFactory()
	s.repository = persistence.NewStoreRepository(db, querier, serviceFactory.NewStorageService())
}

func (s *StoreRepositorySuite) TearDownSuite() {
	s.keycloakContainer.Kill(s.ctx)
	s.pgContainer.Kill(s.ctx)
}

func TestIntegrationStoreRepositorySuite(t *testing.T) {
	suite.Run(t, new(StoreRepositorySuite))
}

func (s *StoreRepositorySuite) TestCreate() {
	ctx := context.Background()
	fakeUserID, err := uuid.NewV7()
	if err != nil {
		s.T().Error("fail on generate fake uuid for test")
	}
	input := store.CreateParams{
		Name:    "Store test 1",
		CpfCnpj: "66063122000135",
		Phone:   "13997590576",
		Address: address.Address{
			AddressLine1: "RUA 1",
			AddressLine2: "657",
			Neighborhood: "Bairro test",
			City:         "Test city",
			State:        "SP",
			PostalCode:   "1213818",
			Country:      "Brasil",
		},
		Type: store.Market,
	}

	id, err := s.repository.Create(ctx, input.Entity(fakeUserID.String()))
	assert.Equal(s.T(), nil, err)

	storeTest, err := s.repository.FindByID(ctx, id)
	if err != nil {
		s.T().Fatalf("fail on verify data in database")
	}
	assert.Equal(s.T(), store.DefaultScore, storeTest.Score)

}
