package persistence

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/oprimogus/cardapiogo/internal/core/item"
	"github.com/oprimogus/cardapiogo/internal/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/database/sqlc"
	"github.com/oprimogus/cardapiogo/pkg/converters"
)

type ItemRepository struct {
	db      *postgres.Database
	querier *sqlc.Queries
}

func NewItemRepository(db *postgres.Database, querier *sqlc.Queries) *ItemRepository {
	return &ItemRepository{db: db, querier: querier}
}

func (i *ItemRepository) CreateItem(ctx context.Context, params item.CreateItemInput, score int) (id int, err error) {
	storeIDConverted, err := converters.StringToUUID(params.StoreID)
	if err != nil {
		return 0, err
	}

	args := sqlc.CreateItemParams{
		StoreID:        storeIDConverted,
		Type:           sqlc.ItemType(params.Type),
		Name:           params.Name,
		Score:          int32(score),
		Description:    params.Description,
		Active:         false,
		DiscountActive: false,
		Price:          int32(params.Price),
		DiscountPrice:  int32(params.Price),
	}

	value, errCreateItem := i.querier.CreateItem(ctx, args)
	if errCreateItem != nil {
		return 0, errCreateItem
	}

	return int(value), nil
}

func (i *ItemRepository) GetItemByID(ctx context.Context, id int) (item.Item, error) {

	itemRow, err := i.querier.GetItemByID(ctx, int64(id))
	if err != nil {
		return item.Item{}, err
	}

	var details map[string]interface{}
	if itemRow.Detail != nil {
		errUnmarshal := json.Unmarshal(itemRow.Detail, &details)
		if errUnmarshal != nil {
			return item.Item{}, fmt.Errorf("fail on convert details column to map: %w", errUnmarshal)
		}
	} else {
		details = nil
	}

	storeIDConverted, err := converters.UuidToString(itemRow.StoreID)
	if err != nil {
		return item.Item{}, err
	}

	createdAt, err := converters.TimestampToTime(itemRow.CreatedAt)
	if err != nil {
		return item.Item{}, err
	}

	updatedAt, err := converters.TimestampToTime(itemRow.UpdatedAt)
	if err != nil {
		return item.Item{}, err
	}

	deletedAt, err := converters.TimestampToTime(itemRow.DeletedAt)
	if err != nil {
		return item.Item{}, err
	}

	return item.Item{
		ID:             int(itemRow.ID),
		StoreID:        *storeIDConverted,
		Type:           item.Type(itemRow.Type),
		Name:           itemRow.Name,
		Description:    itemRow.Description,
		Score:          int(itemRow.Score),
		Image:          itemRow.Image.String,
		Active:         itemRow.Active,
		DiscountActive: itemRow.DiscountActive,
		Price:          int(itemRow.Price),
		DiscountPrice:  int(itemRow.DiscountPrice),
		Details:        details,
		CreatedAt:      *createdAt,
		UpdatedAt:      *updatedAt,
		DeletedAt:      deletedAt,
	}, nil
}

func (i *ItemRepository) GetItemsByFilter(ctx context.Context, filter item.GetItemFilterInput) (*[]item.GetItemFilterOutput, error) {
	args := sqlc.GetItemByFilterParams{
		City:     filter.City,
		Name:     filter.Name,
		Score:    int32(filter.Score),
		Type:     sqlc.ItemType(filter.Type),
		MaxPrice: int32(filter.MaxPrice),
	}

	items, err := i.querier.GetItemByFilter(ctx, args)
	if err != nil {
		return nil, err
	}

	itemsConverted := make([]item.GetItemFilterOutput, len(items))
	for i, v := range items {
		convertedStoreID, err := converters.UuidToString(v.StoreID)
		if err != nil {
			return nil, err
		}
		itemsConverted[i] = item.GetItemFilterOutput{
			ID:             int(v.ID),
			StoreID:        *convertedStoreID,
			Type:           item.Type(v.Type),
			Name:           v.Name,
			Score:          int(v.Score),
			DiscountActive: v.DiscountActive,
			Price:          int(v.Price),
			DiscountPrice:  int(v.DiscountPrice),
		}
	}

	return &itemsConverted, nil
}

func (i *ItemRepository) UpdateItem(ctx context.Context, params item.UpdateItemInput) error {
	args := sqlc.UpdateItemParams{
		ID:             int64(params.ID),
		Type:           sqlc.ItemType(params.Type),
		Name:           params.Name,
		Description:    params.Description,
		Active:         params.Active,
		DiscountActive: params.DiscountActive,
		Price:          int32(params.Price),
		DiscountPrice:  int32(params.DiscountPrice),
	}
	err := i.querier.UpdateItem(ctx, args)
	return err
}

func (i *ItemRepository) DeleteItem(ctx context.Context, id int) error {
	return i.querier.DeleteItem(ctx, int64(id))
}
