//go:build integration

package persistence

//
//import (
//	"github.com/oprimogus/cardapiogo/internal/core/address"
//	"github.com/oprimogus/cardapiogo/internal/core/customer"
//	"github.com/oprimogus/cardapiogo/internal/core/owner"
//	"github.com/oprimogus/cardapiogo/internal/core/store"
//	"github.com/oprimogus/cardapiogo/internal/core/store/product"
//	postgresDB "github.com/oprimogus/cardapiogo/internal/infrastructure/database/postgres"
//	"github.com/oprimogus/cardapiogo/test/integration"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/suite"
//	"golang.org/x/net/context"
//	"testing"
//)
//
//type ProductRepositoryTestSuite struct {
//	suite.Suite
//	mockPostgres       *integration.Container
//	productRepository  ProductRepository
//	storeRepository    StoreRepository
//	customerRepository CustomerRepository
//}
//
//func (s *ProductRepositoryTestSuite) SetupSuite() {
//	ctx := context.Background()
//	mockDB, err := integration.MakePostgres(ctx)
//	if err != nil {
//		assert.Error(s.T(), err)
//	}
//	s.mockPostgres = mockDB
//
//	db := postgresDB.GetTestInstance(mockDB.Port)
//	s.customerRepository = NewCustomerRepository(db)
//	s.productRepository = NewProductRepository(db)
//	s.storeRepository = NewStoreRepository(db)
//}
//
//func (s *ProductRepositoryTestSuite) TearDownSuite() {
//	ctx := context.Background()
//	s.mockPostgres.Kill(ctx)
//}
//
//func TestProductRepositorySuite(t *testing.T) {
//	suite.Run(t, new(ProductRepositoryTestSuite))
//}
//
//func (s *ProductRepositoryTestSuite) TestFindByIDAndSave() {
//	ctx := context.Background()
//	mockCustomer := customer.Customer{
//		ID:       "2575475247247542254",
//		Name:     "John",
//		LastName: "Marston",
//		CPF:      "fake cpf",
//		Phone:    "fake phone",
//		Email:    "fake email",
//		Addresses: []address.Address{
//			{
//				Name:         "test address",
//				AddressLine1: "test Street",
//				AddressLine2: "test Number",
//				Neighborhood: "test neighborhood",
//				City:         "test city",
//				State:        "test state",
//				PostalCode:   "test postal code",
//				Country:      "teste country",
//			},
//		},
//	}
//
//	// mock customer
//	err := s.customerRepository.Save(ctx, &mockCustomer)
//	assert.NoError(s.T(), err)
//
//	mockOwner := owner.NewOwner(mockCustomer.ID)
//
//	// become mock customer an owner
//	err = s.owne.SaveOwner(context.Background(), mockOwner)
//	assert.NoError(s.T(), err)
//
//	// store address
//	addr := address.Address{
//		AddressLine1: "rua 1",
//		AddressLine2: "879",
//		Neighborhood: "test",
//		City:         "test",
//		State:        "test",
//		PostalCode:   "11490-135",
//		Country:      "Brasil",
//	}
//
//	// owner create a store
//	mockStore, err := mockOwner.NewStore(
//		"63432495000148",
//		"Store test",
//		"store from test",
//		"+5513997590579",
//		addr,
//		store.Restaurant)
//	assert.NoError(s.T(), err)
//
//	//Save store
//	err = s.storeRepository.Save(ctx, &mockStore)
//	assert.NoError(s.T(), err)
//
//	// Create product
//	p, err := product.NewProduct(
//		mockStore.ID,
//		"product one", "product tag", "your desc", "P001", 2500, product.Food)
//	assert.NoError(s.T(), err)
//
//	// finally save product
//	err = s.productRepository.Save(ctx, p)
//	assert.NoError(s.T(), err)
//
//	//finally do test
//	persistentProduct, err := s.productRepository.FindByID(ctx, p.ID)
//	assert.NoError(s.T(), err)
//	assert.Equal(s.T(), p, persistentProduct)
//
//	st, err := s.storeRepository.FindStoreByID(ctx, mockStore.ID)
//	assert.NoError(s.T(), err)
//	assert.Equal(s.T(), []string{persistentProduct.ID}, st.Products)
//
//}
