package product

// import (
// 	"context"
// 	"github.com/oprimogus/flyfood-api/internal/infra/database"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )

// type ProductRepositoryTestSuite struct {
// 	suite.Suite
// 	mockPostgres *integration.Container
// 	productRepo  Repository
// 	db           *database.Postgres
// }

// func (s *ProductRepositoryTestSuite) SetupSuite() {
// 	ctx := context.Background()
// 	mockDB, err := integration.MakePostgres(ctx)
// 	if err != nil {
// 		s.T().Fatal("Failed to create test database:", err)
// 	}
// 	s.mockPostgres = mockDB

// 	db, err := database.GetPostgres(ctx)
// 	if err != nil {
// 		s.T().Fatal("Failed to get Postgres instance:", err)
// 	}
// 	s.db = db
// 	s.productRepo = NewRepository(s.db)
// }

// func (s *ProductRepositoryTestSuite) SetupTest() {
// 	// Preparar dados mínimos necessários para foreign keys
// 	s.seedMinimalTestData()
// }

// func (s *ProductRepositoryTestSuite) TearDownTest() {
// 	// Limpar dados de teste após cada teste
// 	s.cleanupTestData()
// }

// func (s *ProductRepositoryTestSuite) TearDownSuite() {
// 	ctx := context.Background()
// 	s.mockPostgres.Kill(ctx)
// }

// // seedMinimalTestData insere apenas os dados mínimos necessários para satisfazer foreign keys
// func (s *ProductRepositoryTestSuite) seedMinimalTestData() {
// 	ctx := context.Background()

// 	// Inserir customer de teste (necessário para owner)
// 	_, err := s.db.Exec(ctx, `
// 		INSERT INTO customers (id, name, last_name, email, phone)
// 		VALUES ('test-owner-001', 'Test', 'Owner', 'test@example.com', '+5511999999999')
// 		ON CONFLICT (id) DO NOTHING
// 	`)
// 	assert.NoError(s.T(), err, "Failed to seed test customer")

// 	// Inserir owner de teste (necessário para store)
// 	_, err = s.db.Exec(ctx, `
// 		INSERT INTO owners (id, signature_active)
// 		VALUES ('test-owner-001', false)
// 		ON CONFLICT (id) DO NOTHING
// 	`)
// 	assert.NoError(s.T(), err, "Failed to seed test owner")

// 	// Inserir store de teste (necessário para product)
// 	_, err = s.db.Exec(ctx, `
// 		INSERT INTO stores (
// 			id, owner_id, cnpj, name, description, active, is_open, phone, score, type,
// 			address_line_1, address_line_2, neighborhood, city, state, postal_code, country
// 		)
// 		VALUES (
// 			'test-store-001', 'test-owner-001', '12345678000190', 'Test Store', 'Test Description',
// 			false, false, '+5511999999999', 500, 'RESTAURANT',
// 			'Test Street 123', 'Suite 100', 'Test Neighborhood', 'Test City', 'TS', '12345-678', 'Brasil'
// 		)
// 		ON CONFLICT (id) DO NOTHING
// 	`)
// 	assert.NoError(s.T(), err, "Failed to seed test store")
// }

// // cleanupTestData remove apenas os dados criados durante os testes
// func (s *ProductRepositoryTestSuite) cleanupTestData() {
// 	ctx := context.Background()

// 	// Remover produtos de teste
// 	_, _ = s.db.Exec(ctx, "DELETE FROM products WHERE store_id = 'test-store-001'")
// }

// // createTestProduct cria um produto básico para testes
// func (s *ProductRepositoryTestSuite) createTestProduct(name, tag, sku string, price int) (Product, error) {
// 	return NewProduct(
// 		"test-store-001",
// 		name,
// 		tag,
// 		"Test product description",
// 		sku,
// 		price,
// 		Food,
// 	)
// }

// func (s *ProductRepositoryTestSuite) TestSave_BasicProduct() {
// 	ctx := context.Background()

// 	// Criar produto básico para teste
// 	testProduct, err := s.createTestProduct("Basic Test Product", "test-tag", "SKU001", 2500)
// 	assert.NoError(s.T(), err, "Should create test product")

// 	// TESTE: Salvar produto
// 	err = s.productRepo.Save(ctx, testProduct)
// 	assert.NoError(s.T(), err, "Should save basic product without error")

// 	// VERIFICAÇÃO: Recuperar e comparar
// 	savedProduct, err := s.productRepo.FindByID(ctx, testProduct.ID)
// 	assert.NoError(s.T(), err, "Should find saved product")

// 	// Verificar campos principais
// 	assert.Equal(s.T(), testProduct.ID, savedProduct.ID)
// 	assert.Equal(s.T(), testProduct.StoreID, savedProduct.StoreID)
// 	assert.Equal(s.T(), testProduct.Name, savedProduct.Name)
// 	assert.Equal(s.T(), testProduct.Tag, savedProduct.Tag)
// 	assert.Equal(s.T(), testProduct.Description, savedProduct.Description)
// 	assert.Equal(s.T(), testProduct.SKU, savedProduct.SKU)
// 	assert.Equal(s.T(), testProduct.Price, savedProduct.Price)
// 	assert.Equal(s.T(), testProduct.Type, savedProduct.Type)
// 	assert.Equal(s.T(), testProduct.ActiveForSale, savedProduct.ActiveForSale)
// 	assert.Equal(s.T(), testProduct.PromoActive, savedProduct.PromoActive)
// 	assert.Equal(s.T(), testProduct.StockQuantity, savedProduct.StockQuantity)
// 	assert.Equal(s.T(), testProduct.Score, savedProduct.Score)
// 	assert.Equal(s.T(), testProduct.PromotionalPrice, savedProduct.PromotionalPrice)
// }

// func (s *ProductRepositoryTestSuite) TestSave_ProductWithStock() {
// 	ctx := context.Background()

// 	// Criar produto com estoque
// 	productWithStock, err := s.createTestProduct("Product With Stock", "stock-tag", "SKU002", 3500)
// 	assert.NoError(s.T(), err)

// 	// Adicionar estoque
// 	err = productWithStock.IncreaseStock(100)
// 	assert.NoError(s.T(), err)

// 	// TESTE: Salvar produto com estoque
// 	err = s.productRepo.Save(ctx, productWithStock)
// 	assert.NoError(s.T(), err, "Should save product with stock")

// 	// VERIFICAÇÃO: Verificar se o estoque foi salvo
// 	savedProduct, err := s.productRepo.FindByID(ctx, productWithStock.ID)
// 	assert.NoError(s.T(), err, "Should find saved product")
// 	assert.Equal(s.T(), 100, savedProduct.StockQuantity, "Should save stock quantity correctly")
// }

// func (s *ProductRepositoryTestSuite) TestSave_ProductWithPromotionalPrice() {
// 	ctx := context.Background()

// 	// Criar produto com preço promocional
// 	productWithPromo, err := s.createTestProduct("Product With Promo", "promo-tag", "SKU003", 5000)
// 	assert.NoError(s.T(), err)

// 	// Definir preço promocional
// 	err = productWithPromo.ChangePromotionalPrice(3500)
// 	assert.NoError(s.T(), err)

// 	productWithPromo.EnablePromotionalPrice()

// 	// TESTE: Salvar produto com preço promocional
// 	err = s.productRepo.Save(ctx, productWithPromo)
// 	assert.NoError(s.T(), err, "Should save product with promotional price")

// 	// VERIFICAÇÃO: Verificar preço promocional
// 	savedProduct, err := s.productRepo.FindByID(ctx, productWithPromo.ID)
// 	assert.NoError(s.T(), err, "Should find saved product")
// 	assert.Equal(s.T(), 3500, savedProduct.PromotionalPrice, "Should save promotional price correctly")
// 	assert.True(s.T(), savedProduct.PromoActive, "Should save promo active status")
// }

// func (s *ProductRepositoryTestSuite) TestSave_UpdateExistingProduct() {
// 	ctx := context.Background()

// 	// Criar e salvar produto inicial
// 	originalProduct, err := s.createTestProduct("Original Product", "original-tag", "SKU004", 1500)
// 	assert.NoError(s.T(), err)

// 	err = s.productRepo.Save(ctx, originalProduct)
// 	assert.NoError(s.T(), err)

// 	// Recuperar produto e modificar
// 	productToUpdate, err := s.productRepo.FindByID(ctx, originalProduct.ID)
// 	assert.NoError(s.T(), err)

// 	// Modificar dados
// 	err = productToUpdate.Update(Water, "Updated Product Name", "Updated description", "NEW-SKU")
// 	assert.NoError(s.T(), err)

// 	err = productToUpdate.ChangePrice(2500)
// 	assert.NoError(s.T(), err)

// 	// TESTE: Salvar produto atualizado
// 	err = s.productRepo.Save(ctx, productToUpdate)
// 	assert.NoError(s.T(), err, "Should update existing product")

// 	// VERIFICAÇÃO: Verificar se mudanças foram persistidas
// 	updatedProduct, err := s.productRepo.FindByID(ctx, originalProduct.ID)
// 	assert.NoError(s.T(), err, "Should find updated product")
// 	assert.Equal(s.T(), "Updated Product Name", updatedProduct.Name)
// 	assert.Equal(s.T(), "Updated description", updatedProduct.Description)
// 	assert.Equal(s.T(), "NEW-SKU", updatedProduct.SKU)
// 	assert.Equal(s.T(), Water, updatedProduct.Type)
// 	assert.Equal(s.T(), 2500, updatedProduct.Price)
// }

// func (s *ProductRepositoryTestSuite) TestSave_ProductWithComplexState() {
// 	ctx := context.Background()

// 	// Criar produto com estado complexo
// 	complexProduct, err := s.createTestProduct("Complex Product", "complex-tag", "COMPLEX-001", 8900)
// 	assert.NoError(s.T(), err)

// 	// Configurar estado complexo
// 	err = complexProduct.IncreaseStock(250)
// 	assert.NoError(s.T(), err)

// 	err = complexProduct.ChangePromotionalPrice(6500)
// 	assert.NoError(s.T(), err)

// 	complexProduct.EnablePromotionalPrice()

// 	err = complexProduct.EnableForSale()
// 	assert.NoError(s.T(), err)

// 	err = complexProduct.ChangeTag("premium-product")
// 	assert.NoError(s.T(), err)

// 	// TESTE: Salvar produto complexo
// 	err = s.productRepo.Save(ctx, complexProduct)
// 	assert.NoError(s.T(), err, "Should save complex product")

// 	// VERIFICAÇÃO: Verificar todos os campos
// 	savedProduct, err := s.productRepo.FindByID(ctx, complexProduct.ID)
// 	assert.NoError(s.T(), err, "Should find complex product")

// 	assert.Equal(s.T(), "Complex Product", savedProduct.Name)
// 	assert.Equal(s.T(), "premium-product", savedProduct.Tag)
// 	assert.Equal(s.T(), 250, savedProduct.StockQuantity)
// 	assert.Equal(s.T(), 8900, savedProduct.Price)
// 	assert.Equal(s.T(), 6500, savedProduct.PromotionalPrice)
// 	assert.True(s.T(), savedProduct.ActiveForSale)
// 	assert.True(s.T(), savedProduct.PromoActive)
// 	assert.Equal(s.T(), Food, savedProduct.Type)
// }

// func (s *ProductRepositoryTestSuite) TestFindByID_NonExistentProduct() {
// 	ctx := context.Background()

// 	// TESTE: Tentar buscar produto que não existe
// 	_, err := s.productRepo.FindByID(ctx, "non-existent-product-id")

// 	// VERIFICAÇÃO: Deve retornar erro
// 	assert.Error(s.T(), err, "Should return error for non-existent product")
// }

// func (s *ProductRepositoryTestSuite) TestSave_DifferentProductTypes() {
// 	ctx := context.Background()

// 	testCases := []struct {
// 		name        string
// 		productType Type
// 		tag         string
// 		sku         string
// 		price       int
// 	}{
// 		{
// 			name:        "Food Product",
// 			productType: Food,
// 			tag:         "food-tag",
// 			sku:         "FOOD-001",
// 			price:       1500,
// 		},
// 		{
// 			name:        "Water Product",
// 			productType: Water,
// 			tag:         "water-tag",
// 			sku:         "WATER-001",
// 			price:       500,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.T().Run(tc.name, func(t *testing.T) {
// 			// Criar produto do tipo específico
// 			product, err := NewProduct(
// 				"test-store-001",
// 				tc.name,
// 				tc.tag,
// 				"Product description for "+tc.name,
// 				tc.sku,
// 				tc.price,
// 				tc.productType,
// 			)
// 			assert.NoError(t, err)

// 			// TESTE: Salvar produto
// 			err = s.productRepo.Save(ctx, product)
// 			assert.NoError(t, err, "Should save %s", tc.name)

// 			// VERIFICAÇÃO: Verificar tipo do produto
// 			savedProduct, err := s.productRepo.FindByID(ctx, product.ID)
// 			assert.NoError(t, err, "Should find saved %s", tc.name)
// 			assert.Equal(t, tc.productType, savedProduct.Type, "Should save correct product type")
// 			assert.Equal(t, tc.price, savedProduct.Price, "Should save correct price")
// 			assert.Equal(t, tc.sku, savedProduct.SKU, "Should save correct SKU")
// 		})
// 	}
// }

// func (s *ProductRepositoryTestSuite) TestSave_DefaultValues() {
// 	ctx := context.Background()

// 	// Criar produto básico (que deve ter valores padrão)
// 	defaultProduct, err := s.createTestProduct("Default Values Product", "default-tag", "DEFAULT-001", 1000)
// 	assert.NoError(s.T(), err)

// 	// TESTE: Salvar produto com valores padrão
// 	err = s.productRepo.Save(ctx, defaultProduct)
// 	assert.NoError(s.T(), err, "Should save product with default values")

// 	// VERIFICAÇÃO: Verificar valores padrão
// 	savedProduct, err := s.productRepo.FindByID(ctx, defaultProduct.ID)
// 	assert.NoError(s.T(), err, "Should find saved product")

// 	// Verificar valores padrão esperados
// 	assert.False(s.T(), savedProduct.ActiveForSale, "Should default to not active for sale")
// 	assert.False(s.T(), savedProduct.PromoActive, "Should default to no active promo")
// 	assert.Equal(s.T(), 0, savedProduct.StockQuantity, "Should default to zero stock")
// 	assert.Equal(s.T(), 500, savedProduct.Score, "Should have default score") // assumindo defaultScore = 500
// 	assert.Equal(s.T(), 0, savedProduct.PromotionalPrice, "Should default to zero promotional price")
// 	assert.Empty(s.T(), savedProduct.Image, "Should default to empty image")
// }

// func TestProductRepositorySuite(t *testing.T) {
// 	suite.Run(t, new(ProductRepositoryTestSuite))
// }
