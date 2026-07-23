package product

// import (
// 	"context"
// 	"encoding/json"
// 	"github.com/oprimogus/flyfood-api/internal/infra/database"

// 	"github.com/jackc/pgx/v5/pgtype"
// 	"github.com/oprimogus/flyfood-api/internal/infra/database/sqlc"
// )

// type Repository interface {
// 	FindByID(ctx context.Context, id string) (Product, error)
// 	Save(ctx context.Context, p Product) error
// }

// type repository struct {
// 	db *database.Postgres
// 	q  sqlc.Querier
// }

// func NewRepository(db *database.Postgres) Repository {
// 	return repository{
// 		db: db,
// 		q:  sqlc.New(db),
// 	}
// }

// func (r repository) FindByID(ctx context.Context, id string) (Product, error) {
// 	convUUID, err := database.StringToUUID(id)
// 	if err != nil {
// 		return Product{}, err
// 	}

// 	p, err := r.q.FindProductByID(ctx, convUUID)
// 	if err != nil {
// 		return Product{}, err
// 	}

// 	convStoreID, err := database.UuidToString(p.StoreID)
// 	if err != nil {
// 		return Product{}, err
// 	}

// 	var details map[string]any
// 	err = json.Unmarshal(p.Details, &details)
// 	if err != nil {
// 		return Product{}, err
// 	}

// 	return Product{
// 		ID:               id,
// 		StoreID:          *convStoreID,
// 		SKU:              p.Sku.String,
// 		ActiveForSale:    p.ActiveForSale,
// 		PromoActive:      p.PromoActive,
// 		Type:             Type(p.Type),
// 		Tag:              p.Tag,
// 		Name:             p.Name,
// 		Description:      p.Description,
// 		StockQuantity:    int(p.StockQuantity),
// 		Score:            int(p.Score),
// 		Image:            p.ImageUrl.String,
// 		Details:          details,
// 		Price:            int(p.Price),
// 		PromotionalPrice: int(p.PromotionalPrice.Int32),
// 	}, nil
// }

// func (r repository) Save(ctx context.Context, p Product) error {
// 	id, err := database.StringToUUID(p.ID)
// 	if err != nil {
// 		return err
// 	}

// 	storeID, err := database.StringToUUID(p.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	details, err := json.Marshal(p.Details)
// 	if err != nil {
// 		return err
// 	}

// 	args := sqlc.SaveProductParams{
// 		ID:               id,
// 		StoreID:          storeID,
// 		Sku:              database.StringToText(p.SKU),
// 		ActiveForSale:    p.ActiveForSale,
// 		PromoActive:      p.PromoActive,
// 		Type:             sqlc.ProductType(p.Type),
// 		Tag:              p.Tag,
// 		Name:             p.Name,
// 		Description:      p.Description,
// 		StockQuantity:    int32(p.StockQuantity),
// 		Score:            int32(p.Score),
// 		ImageUrl:         database.StringToText(p.Image),
// 		Details:          details,
// 		Price:            int32(p.Price),
// 		PromotionalPrice: pgtype.Int4{Int32: int32(p.PromotionalPrice), Valid: p.PromotionalPrice != 0},
// 	}

// 	return r.q.SaveProduct(ctx, args)
// }
