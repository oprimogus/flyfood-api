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

type itemRepository struct {
	db      *postgres.PostgresDatabase
	querier *sqlc.Queries
}

func NewItemRepository(db *postgres.PostgresDatabase, querier *sqlc.Queries) *itemRepository {
	return &itemRepository{db: db, querier: querier}
}

func (i *itemRepository) CreateItem(ctx context.Context, params item.CreateItemInput, score int) (id int, err error) {
	storeIDConverted, err := converters.ConvertStringToUUID(params.StoreID)
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

func (i *itemRepository) GetItemByID(ctx context.Context, id int) (item.GetItemByIDOutput, error) {

	itemRow, err := i.querier.GetItemByID(ctx, int64(id))
	if err != nil {
		return item.GetItemByIDOutput{}, err
	}

	var details map[string]interface{}
	if itemRow.Detail != nil {
		errUnmarshal := json.Unmarshal(itemRow.Detail, &details)
		if errUnmarshal != nil {
			return item.GetItemByIDOutput{}, fmt.Errorf("fail on convert details column to map: %w", errUnmarshal)
		}
	} else {
		details = nil
	}

	return item.GetItemByIDOutput{
		Type:           item.ItemType(itemRow.Type),
		Name:           itemRow.Name,
		Description:    itemRow.Description,
		Score:          int(itemRow.Score),
		Image:          itemRow.Image.String,
		Active:         itemRow.Active,
		DiscountActive: itemRow.DiscountActive,
		Price:          int(itemRow.Price),
		DiscountPrice:  int(itemRow.DiscountPrice),
		Details:        details,
	}, nil
}

func (i *itemRepository) GetItemsByFilter(ctx context.Context, filter item.GetItemFilter) (*[]item.GetItemFilter, error) {
	return nil, nil
}

func (i *itemRepository) UpdateItem(ctx context.Context, params item.UpdateItemInput) error {
	return nil
}

func (i *itemRepository) DeleteItem(ctx context.Context, id int) error {
	return nil
}
