package store

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/oprimogus/flyfood-api/internal/config"
	"github.com/oprimogus/flyfood-api/internal/core/address"
	"github.com/oprimogus/flyfood-api/internal/infra/database"
	"github.com/oprimogus/flyfood-api/pkg/testcontainers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StoreRepositoryTestSuite struct {
	suite.Suite
	mockPostgres *testcontainers.PostgresContainer
	storeRepo    Repository
	db           *database.Postgres
	ownerID      uuid.UUID
}

func (s *StoreRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	pgCfg := testcontainers.PostgresConfig {
		DatabaseName: "flyfood",
		Username:     "test",
		Password:     "test",
	}

	mockDB, err := testcontainers.MakePostgres(ctx, pgCfg)
	if err != nil {
		s.T().Fatal("Failed to start postgres container:", err)
	}
	s.mockPostgres = mockDB

	cfg := config.Get()
	cfg.Postgres.Host = "localhost"
	cfg.Postgres.Port = mockDB.Port
	cfg.Postgres.UserName = pgCfg.Username
	cfg.Postgres.Password = pgCfg.Password
	cfg.Postgres.DatabaseName = pgCfg.DatabaseName

	s.db, err = database.GetPostgres(ctx)
	if err != nil {
		s.T().Fatal("Failed to connect to test database:", err)
	}
	s.storeRepo = NewRepository(s.db)
	s.ownerID, err = uuid.NewV7()
	s.NoError(err)
}

func (s *StoreRepositoryTestSuite) SetupTest() {
	s.seedMinimalTestData()
}

func (s *StoreRepositoryTestSuite) TearDownTest() {
	s.cleanupTestData()
}

func (s *StoreRepositoryTestSuite) TearDownSuite() {
	s.mockPostgres.Kill(context.Background())
}

// ============================================================
// Helpers
// ============================================================

func (s *StoreRepositoryTestSuite) seedMinimalTestData() {
	ctx := context.Background()

	_, err := s.db.Exec(ctx, `
		INSERT INTO customer (id, external_id, name, last_name, email, phone)
		VALUES ($1, 'ext-test-owner', 'Test', 'Owner', 'owner@test.com', '+5511999999999')
		ON CONFLICT (id) DO NOTHING
	`, s.ownerID)
	assert.NoError(s.T(), err, "seed customer")

	_, err = s.db.Exec(ctx, `
		INSERT INTO owner (id, signature_active, created_at)
		VALUES ($1, true, NOW())
		ON CONFLICT (id) DO NOTHING
	`, s.ownerID)
	assert.NoError(s.T(), err, "seed owner")
}

func (s *StoreRepositoryTestSuite) cleanupTestData() {
	ctx := context.Background()
	_, _ = s.db.Exec(ctx, `
		DELETE FROM store_business_hour
		WHERE store_id IN (SELECT id FROM store WHERE owner_id = $1)
	`, s.ownerID)
	_, _ = s.db.Exec(ctx, `
		DELETE FROM store_payment_method
		WHERE store_id IN (SELECT id FROM store WHERE owner_id = $1)
	`, s.ownerID)
	_, _ = s.db.Exec(ctx, "DELETE FROM store WHERE owner_id = $1", s.ownerID)
	_, _ = s.db.Exec(ctx, "DELETE FROM address WHERE id NOT IN (SELECT address_id FROM store)")
}

func (s *StoreRepositoryTestSuite) defaultAddress() address.Address {
	return address.Address{
		ID:           uuid.New(), // Atribui um tipo estruturado uuid.UUID válido
		Name:         "Test Address",
		AddressLine1: "Rua das Flores, 123",
		AddressLine2: "Apto 10",
		Neighborhood: "Centro",
		City:         "São Paulo",
		State:        "SP",
		PostalCode:   "01310-100",
		Country:      "Brasil",
		Latitude:     -23.5505,
		Longitude:    -46.6333,
	}
}

func (s *StoreRepositoryTestSuite) mustSave(st *Store) {
	err := s.storeRepo.Save(context.Background(), st)
	s.Require().NoError(err, "precondition: save store")
}

func (s *StoreRepositoryTestSuite) mustFind(id string) *Store {
	st, err := s.storeRepo.FindByID(context.Background(), id)
	s.Require().NoError(err, "precondition: find store")
	return st
}

func (s *StoreRepositoryTestSuite) TestSave_MinimalStore() {
	ctx := context.Background()

	st, err := NewStore(
		s.ownerID,
		"24611859000103",
		"Minimal Store",
		"Descrição válida", // Modificado para passar nos validadores de tamanho da entidade, se houver
		"+5511997590670",
		s.defaultAddress(),
		Restaurant,
	)
	s.Require().NoError(err)

	err = s.storeRepo.Save(ctx, st)
	assert.NoError(s.T(), err)

	saved := s.mustFind(st.ID.String())
	assert.Equal(s.T(), st.ID, saved.ID)
	assert.Equal(s.T(), st.CNPJ, saved.CNPJ)
	assert.Equal(s.T(), st.Name, saved.Name)
	assert.Equal(s.T(), st.Phone, saved.Phone)
	assert.Equal(s.T(), st.Type, saved.Type)
	assert.False(s.T(), saved.Active)
	assert.False(s.T(), saved.IsOpen)
	assert.Equal(s.T(), DefaultScore, saved.Score) // Corrigido de 0 para DefaultScore (500)
	assert.Empty(s.T(), saved.BusinessHours)
	assert.Empty(s.T(), saved.PaymentMethods)
}

func (s *StoreRepositoryTestSuite) TestSave_WithBusinessHours() {
	ctx := context.Background()

	st, err := NewStore(
		s.ownerID, "24611859000103", "Store BH", "Descrição válida",
		"+5511997590670", s.defaultAddress(), Restaurant,
	)
	s.Require().NoError(err)

	op1, err := NewMinutesOfDayFromHHMM("08:00")
	s.Require().NoError(err)

	cl1, err := NewMinutesOfDayFromHHMM("18:00")
	s.Require().NoError(err)

	op2, err := NewMinutesOfDayFromHHMM("09:00")
	s.Require().NoError(err)

	cl2, err := NewMinutesOfDayFromHHMM("17:00")
	s.Require().NoError(err)

	op3, err := NewMinutesOfDayFromHHMM("10:00")
	s.Require().NoError(err)

	cl3, err := NewMinutesOfDayFromHHMM("22:00")
	s.Require().NoError(err)

	hours := []BusinessHours{
		{WeekDay: 1, OpeningTime: op1, ClosingTime: cl1},
		{WeekDay: 3, OpeningTime: op2, ClosingTime: cl2},
		{WeekDay: 5, OpeningTime: op3, ClosingTime: cl3},
	}
	for _, bh := range hours {
		s.Require().NoError(st.AddNewBusinessHour(bh))
	}

	err = s.storeRepo.Save(ctx, st)
	assert.NoError(s.T(), err)

	saved := s.mustFind(st.ID.String())
	assert.Len(s.T(), saved.BusinessHours, 3)

	for _, expected := range hours {
		found := false
		for _, actual := range saved.BusinessHours {
			if actual.WeekDay == expected.WeekDay &&
				actual.OpeningTime == expected.OpeningTime &&
				actual.ClosingTime == expected.ClosingTime {
				found = true
				break
			}
		}
		assert.True(s.T(), found, "missing business hour: %+v", expected)
	}
}

func (s *StoreRepositoryTestSuite) TestSave_WithPaymentMethods() {
	ctx := context.Background()

	st, err := NewStore(
		s.ownerID, "24611859000103", "Store PM", "Descrição válida",
		"+5511997590670", s.defaultAddress(), Pharmacy,
	)
	s.Require().NoError(err)

	methods := []PaymentMethod{Credit, Debit, Pix}
	for _, pm := range methods {
		s.Require().NoError(st.AddPaymentMethod(pm))
	}

	err = s.storeRepo.Save(ctx, st)
	assert.NoError(s.T(), err)

	saved := s.mustFind(st.ID.String())
	assert.Len(s.T(), saved.PaymentMethods, 3)

	for _, expected := range methods {
		found := false
		for _, actual := range saved.PaymentMethods {
			if actual == expected {
				found = true
				break
			}
		}
		assert.True(s.T(), found, "missing payment method: %s", expected)
	}
}

func (s *StoreRepositoryTestSuite) TestSave_WithAllData() {
	ctx := context.Background()

	st, err := NewStore(
		s.ownerID, "24611859000103", "Full Store", "Descrição completa",
		"+5511997590670", s.defaultAddress(), Pub,
	)
	s.Require().NoError(err)

	op1, err := NewMinutesOfDayFromHHMM("08:00")
	s.Require().NoError(err)

	cl1, err := NewMinutesOfDayFromHHMM("18:00")
	s.Require().NoError(err)

	op2, err := NewMinutesOfDayFromHHMM("10:00")
	s.Require().NoError(err)

	cl2, err := NewMinutesOfDayFromHHMM("23:00")
	s.Require().NoError(err)

	hours := []BusinessHours{
		{WeekDay: 1, OpeningTime: op1, ClosingTime: cl1},
		{WeekDay: 2, OpeningTime: op1, ClosingTime: cl1},
		{WeekDay: 5, OpeningTime: op2, ClosingTime: cl2},
		{WeekDay: 6, OpeningTime: op2, ClosingTime: cl2},
	}
	for _, bh := range hours {
		s.Require().NoError(st.AddNewBusinessHour(bh))
	}

	methods := []PaymentMethod{Credit, Debit, Pix, Cash}
	for _, pm := range methods {
		s.Require().NoError(st.AddPaymentMethod(pm))
	}

	err = s.storeRepo.Save(ctx, st)
	assert.NoError(s.T(), err)

	saved := s.mustFind(st.ID.String())
	assert.Len(s.T(), saved.BusinessHours, len(hours))
	assert.Len(s.T(), saved.PaymentMethods, len(methods))
	assert.Equal(s.T(), st.Description, saved.Description)
	assert.Equal(s.T(), st.Address.City, saved.Address.City)
	assert.InDelta(s.T(), st.Address.Latitude, saved.Address.Latitude, 0.0001)
	assert.InDelta(s.T(), st.Address.Longitude, saved.Address.Longitude, 0.0001)
}

// ============================================================
// Save — atualização (upsert)
// ============================================================

func (s *StoreRepositoryTestSuite) TestSave_UpdateProfile() {
	ctx := context.Background()

	st, err := NewStore(
		s.ownerID, "24611859000103", "Original Name", "Descrição original",
		"+5511997590670", s.defaultAddress(), Market,
	)
	s.Require().NoError(err)
	s.mustSave(st)

	updated := s.mustFind(st.ID.String())
	newAddr := s.defaultAddress()
	newAddr.City = "Campinas"

	// Ajustado de UpdateStoreProfile para o método correto: UpdateStore
	err = updated.UpdateStore("Updated Name", "Nova descrição", "+5511997590670", newAddr, Restaurant)
	s.Require().NoError(err)

	err = s.storeRepo.Save(ctx, updated)
	assert.NoError(s.T(), err)

	final := s.mustFind(st.ID.String())
	assert.Equal(s.T(), "Updated Name", final.Name)
	assert.Equal(s.T(), "Nova descrição", final.Description)
	assert.Equal(s.T(), "+5511997590670", final.Phone)
	assert.Equal(s.T(), "Campinas", final.Address.City)
	assert.Equal(s.T(), Restaurant, final.Type)
}

func (s *StoreRepositoryTestSuite) TestSave_UpdateBusinessHours_AddAndRemove() {
	ctx := context.Background()

	st, err := NewStore(
		s.ownerID, "24611859000103", "Store BH Update", "Descrição válida",
		"+5511997590670", s.defaultAddress(), Restaurant,
	)
	s.Require().NoError(err)

	op, err := NewMinutesOfDayFromHHMM("08:00")
	s.Require().NoError(err)

	cl, err := NewMinutesOfDayFromHHMM("18:00")
	s.Require().NoError(err)

	// Salva com SEG e TER
	s.Require().NoError(st.AddNewBusinessHour(BusinessHours{WeekDay: 1, OpeningTime: op, ClosingTime: cl}))
	s.Require().NoError(st.AddNewBusinessHour(BusinessHours{WeekDay: 2, OpeningTime: op, ClosingTime: cl}))
	s.mustSave(st)

	// Atualiza: remove TER, adiciona QUA
	updated := s.mustFind(st.ID.String())
	updated.BusinessHours = []BusinessHours{
		{WeekDay: 1, OpeningTime: op, ClosingTime: cl},
		{WeekDay: 3, OpeningTime: op, ClosingTime: cl},
	}

	err = s.storeRepo.Save(ctx, updated)
	assert.NoError(s.T(), err)

	final := s.mustFind(st.ID.String())
	assert.Len(s.T(), final.BusinessHours, 2)

	weekdays := make([]int, len(final.BusinessHours))
	for i, bh := range final.BusinessHours {
		weekdays[i] = bh.WeekDay
	}
	assert.Contains(s.T(), weekdays, 1, "SEG deve existir")
	assert.Contains(s.T(), weekdays, 3, "QUA deve existir")
	assert.NotContains(s.T(), weekdays, 2, "TER deve ter sido removido")
}

func (s *StoreRepositoryTestSuite) TestSave_UpdateBusinessHours_RemoveAll() {
	ctx := context.Background()

	st, err := NewStore(
		s.ownerID, "24611859000103", "Store BH RemoveAll", "Descrição válida",
		"+5511997590670", s.defaultAddress(), Restaurant,
	)
	s.Require().NoError(err)
	
	op, err := NewMinutesOfDayFromHHMM("08:00")
	s.Require().NoError(err)

	cl, err := NewMinutesOfDayFromHHMM("18:00")
	s.Require().NoError(err)
	
	s.Require().NoError(st.AddNewBusinessHour(BusinessHours{WeekDay: 1, OpeningTime: op, ClosingTime: cl}))
	s.Require().NoError(st.AddNewBusinessHour(BusinessHours{WeekDay: 2, OpeningTime: op, ClosingTime: cl}))
	s.mustSave(st)

	updated := s.mustFind(st.ID.String())
	updated.BusinessHours = []BusinessHours{}

	err = s.storeRepo.Save(ctx, updated)
	assert.NoError(s.T(), err)

	final := s.mustFind(st.ID.String())
	assert.Empty(s.T(), final.BusinessHours, "todos os horários devem ter sido removidos")
}

func (s *StoreRepositoryTestSuite) TestSave_UpdatePaymentMethods_AddAndRemove() {
	ctx := context.Background()

	st, err := NewStore(
		s.ownerID, "24611859000103", "Store PM Update", "Descrição válida",
		"+5511997590670", s.defaultAddress(), Market,
	)
	s.Require().NoError(err)
	s.Require().NoError(st.AddPaymentMethod(Credit))
	s.Require().NoError(st.AddPaymentMethod(Cash))
	s.mustSave(st)

	updated := s.mustFind(st.ID.String())
	updated.PaymentMethods = []PaymentMethod{Credit, Pix} // remove Cash, adiciona Pix

	err = s.storeRepo.Save(ctx, updated)
	assert.NoError(s.T(), err)

	final := s.mustFind(st.ID.String())
	slog.Info("final", "paymentMethods", final)
	assert.Len(s.T(), final.PaymentMethods, 2)
	assert.Contains(s.T(), final.PaymentMethods, Credit)
	assert.Contains(s.T(), final.PaymentMethods, Pix)
	assert.NotContains(s.T(), final.PaymentMethods, Cash, "Cash deve ter sido removido")
}

func (s *StoreRepositoryTestSuite) TestSave_Idempotent() {
	ctx := context.Background()

	st, err := NewStore(
		s.ownerID, "24611859000103", "Idempotent Store", "Descrição válida",
		"+5511997590670", s.defaultAddress(), Convenience,
	)
	s.Require().NoError(err)

	err = s.storeRepo.Save(ctx, st)
	assert.NoError(s.T(), err)

	err = s.storeRepo.Save(ctx, st)
	assert.NoError(s.T(), err)

	saved := s.mustFind(st.ID.String())
	assert.Equal(s.T(), st.Name, saved.Name)
}

// ============================================================
// FindByID
// ============================================================

func (s *StoreRepositoryTestSuite) TestFindByID_NotFound() {
	_, err := s.storeRepo.FindByID(context.Background(), uuid.New().String())
	slog.Error("TestFindByID_NotFound", slog.Any("error", err))
	assert.Error(s.T(), err)
}

func (s *StoreRepositoryTestSuite) TestFindByID_InvalidID() {
	_, err := s.storeRepo.FindByID(context.Background(), "not-a-uuid")
	assert.Error(s.T(), err)
}

func (s *StoreRepositoryTestSuite) TestFindByID_AddressFields() {
	addr := s.defaultAddress()
	addr.City = "Guarujá"
	addr.State = "SP"
	addr.Latitude = -23.9934
	addr.Longitude = -46.2564

	st, err := NewStore(
		s.ownerID, "24611859000103", "Address Fields Store", "Descrição válida",
		"+5511997590670", addr, Restaurant,
	)
	slog.Error("TestFindByID_AddressFields", "err", err)
	s.Require().NoError(err)
	s.mustSave(st)

	saved := s.mustFind(st.ID.String())
	assert.Equal(s.T(), "Guarujá", saved.Address.City)
	assert.Equal(s.T(), "SP", saved.Address.State)
	assert.InDelta(s.T(), -23.9934, saved.Address.Latitude, 0.0001)
	assert.InDelta(s.T(), -46.2564, saved.Address.Longitude, 0.0001)
}

// ============================================================
// Entry point
// ============================================================

func TestStoreRepositorySuite(t *testing.T) {
	suite.Run(t, new(StoreRepositoryTestSuite))
}