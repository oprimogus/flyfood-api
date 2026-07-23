package ownership

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oprimogus/flyfood-api/internal/config"
	"github.com/oprimogus/flyfood-api/internal/core/customer"
	"github.com/oprimogus/flyfood-api/internal/infra/database"
	"github.com/oprimogus/flyfood-api/pkg/testcontainers"
	"github.com/stretchr/testify/suite"
)

type OwnershipRepositoryTestSuite struct {
	suite.Suite
	ctx          context.Context
	mockDB       *testcontainers.PostgresContainer
	repo         Repository
	customerRepo customer.Repository
}

func (s *OwnershipRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	pgCfg := testcontainers.PostgresConfig{
		DatabaseName: "flyfood",
		Username:     "test",
		Password:     "test",
	}

	mockDB, err := testcontainers.MakePostgres(s.ctx, pgCfg)
	s.NoError(err, "fail on create mock postgres")
	s.mockDB = mockDB

	cfg := config.Get()
	cfg.Postgres.Host = "localhost"
	cfg.Postgres.Port = mockDB.Port
	cfg.Postgres.UserName = pgCfg.Username
	cfg.Postgres.Password = pgCfg.Password
	cfg.Postgres.DatabaseName = pgCfg.DatabaseName

	db, err := database.GetPostgres(s.ctx)
	s.NoError(err, "fail on get connection with database")

	err = db.Migrate(s.ctx)
	s.NoError(err, "failed to migrate database")

	s.repo = NewRepository(db)
	s.customerRepo = customer.NewRepository(db)
}

func (s *OwnershipRepositoryTestSuite) TearDownSuite() {
	s.mockDB.Kill(s.ctx)
}

func (s *OwnershipRepositoryTestSuite) createTestCustomer(id uuid.UUID, cpf string) {
	c := &customer.Customer{
		ID:         id,
		ExternalID: "zitadel_" + id.String(),
		Name:       "Dono",
		LastName:   "Da Silva",
		CPF:        cpf,
		Email:      id.String() + "@teste.com",
		Phone:      "11999999999",
	}
	err := s.customerRepo.Save(s.ctx, c)
	s.NoError(err, "fail on create setup customer dependency")
}

func TestOwnershipRepository(t *testing.T) {
	suite.Run(t, new(OwnershipRepositoryTestSuite))
}

func (s *OwnershipRepositoryTestSuite) TestSave() {
	s.Run("Should save an owner with success", func() {
		ownerID, err := uuid.NewV7()
		s.NoError(err, "fail on generate UUID V7")

		s.createTestCustomer(ownerID, "41603346015")
		newOwner := &Owner{
			ID:              ownerID,
			SignatureActive: true,
		}

		err = s.repo.Save(s.ctx, newOwner)
		s.NoError(err)

		saved, err := s.repo.FindByID(s.ctx, ownerID)
		s.NoError(err)
		s.NotNil(saved)
		s.Equal(newOwner.ID, saved.ID)
		s.Equal(newOwner.SignatureActive, saved.SignatureActive)
	})
}

func (s *OwnershipRepositoryTestSuite) TestFindByID() {
	s.Run("Should find an existing owner by ID", func() {
		ownerID, err := uuid.NewV7()
		s.NoError(err, "fail on generate UUID V7")
		s.createTestCustomer(ownerID, "07064297027")

		expectedOwner := &Owner{
			ID:              ownerID,
			SignatureActive: false,
		}

		err = s.repo.Save(s.ctx, expectedOwner)
		s.NoError(err)

		actualOwner, err := s.repo.FindByID(s.ctx, ownerID)
		s.NoError(err)
		s.NotNil(actualOwner)
		s.Equal(expectedOwner.ID, actualOwner.ID)
		s.Equal(expectedOwner.SignatureActive, actualOwner.SignatureActive)
		s.WithinDuration(time.Now(), actualOwner.CreatedAt, 2*time.Second) // Garante que a data de criação bate com o "agora"
	})

	s.Run("Should return error when owner does not exist", func() {
		randomID, err := uuid.NewV7()
		s.NoError(err, "fail on generate UUID V7")

		actualOwner, err := s.repo.FindByID(s.ctx, randomID)
		s.Error(err)
		s.Nil(actualOwner)
	})
}
