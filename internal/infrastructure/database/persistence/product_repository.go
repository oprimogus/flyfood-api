package persistence

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oprimogus/cardapiogo/internal/core/product"
	postgresDB "github.com/oprimogus/cardapiogo/internal/infrastructure/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/database/sqlc"
	"github.com/oprimogus/cardapiogo/pkg/converters"
)

type ProductRepository struct {
	db *pgxpool.Pool
	q  sqlc.Querier
}

func NewProductRepository(db *postgresDB.Database) ProductRepository {
	return ProductRepository{db: db.GetDB(), q: sqlc.New(db.GetDB())}
}

func (r ProductRepository) FindByID(ctx context.Context, id string) (product.Product, error) {
	convUUID, err := converters.StringToUUID(id)
	if err != nil {
		return product.Product{}, err
	}

	p, err := r.q.FindProductByID(ctx, convUUID)
	if err != nil {
		return product.Product{}, err
	}

	convStoreID, err := converters.UuidToString(p.StoreID)
	if err != nil {
		return product.Product{}, err
	}

	var details map[string]interface{}
	err = json.Unmarshal(p.Details, &details)
	if err != nil {
		return product.Product{}, err
	}

	return product.Product{
		ID:               id,
		StoreID:          *convStoreID,
		SKU:              p.Sku.String,
		ActiveForSale:    p.ActiveForSale,
		PromoActive:      p.PromoActive,
		Type:             product.Type(p.Type),
		Tag:              p.Tag,
		Name:             p.Name,
		Description:      p.Description,
		StockQuantity:    int(p.StockQuantity),
		Score:            int(p.Score),
		Image:            p.ImageUrl.String,
		Details:          details,
		Price:            int(p.Price),
		PromotionalPrice: int(p.PromotionalPrice.Int32),
	}, nil
}

func (r ProductRepository) Save(ctx context.Context, p product.Product) error {

	id, err := converters.StringToUUID(p.ID)
	if err != nil {
		return err
	}

	storeID, err := converters.StringToUUID(p.StoreID)
	if err != nil {
		return err
	}

	details, err := json.Marshal(p.Details)
	if err != nil {
		return err
	}

	args := sqlc.SaveProductParams{
		ID:               id,
		StoreID:          storeID,
		Sku:              converters.StringToText(p.SKU),
		ActiveForSale:    p.ActiveForSale,
		PromoActive:      p.PromoActive,
		Type:             sqlc.ProductType(p.Type),
		Tag:              p.Tag,
		Name:             p.Name,
		Description:      p.Description,
		StockQuantity:    int32(p.StockQuantity),
		Score:            int32(p.Score),
		ImageUrl:         converters.StringToText(p.Image),
		Details:          details,
		Price:            int32(p.Price),
		PromotionalPrice: pgtype.Int4{Int32: int32(p.PromotionalPrice), Valid: true},
	}

	return r.q.SaveProduct(ctx, args)
}
