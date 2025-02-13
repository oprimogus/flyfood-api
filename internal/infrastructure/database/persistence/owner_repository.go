package persistence

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oprimogus/cardapiogo/internal/core/owner"
	postgresDB "github.com/oprimogus/cardapiogo/internal/infrastructure/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/database/sqlc"
	"github.com/oprimogus/cardapiogo/pkg/converters"
)

type OwnerRepository struct {
	db *pgxpool.Pool
	q  sqlc.Querier
}

func NewOwnerRepository(db *postgresDB.Database) OwnerRepository {
	return OwnerRepository{db: db.GetDB(), q: sqlc.New(db.GetDB())}
}

func (r OwnerRepository) FindOwnerByID(ctx context.Context, id string) (*owner.Owner, error) {
	ow, err := r.q.FindOwnerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	stIds := make([]string, len(ow.StoreIds))
	for i, s := range ow.StoreIds {
		vConv, err := converters.UuidToString(s)
		if err != nil {
			return nil, err
		}
		stIds[i] = *vConv
	}

	return &owner.Owner{
		ID:              ow.ID,
		SignatureActive: ow.SignatureActive,
		StoresID:        stIds,
	}, nil
}

func (r OwnerRepository) SaveOwner(ctx context.Context, ow owner.Owner) error {
	err := r.q.SaveOwner(ctx, sqlc.SaveOwnerParams{
		ID:              ow.ID,
		SignatureActive: ow.SignatureActive,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r OwnerRepository) IsOwner(ctx context.Context, customerID string) (bool, error) {
	isOwner, err := r.q.IsOwner(ctx, customerID)
	if err != nil {
		return false, err
	}

	return isOwner, nil
}

func (r OwnerRepository) IsOwnerOf(ctx context.Context, customerID, storeID string) (bool, error) {
	idConv, err := converters.StringToUUID(storeID)
	if err != nil {
		return false, err
	}

	isOwner, err := r.q.IsOwnerOf(ctx, sqlc.IsOwnerOfParams{OwnerID: customerID, ID: idConv})
	if err != nil {
		return false, err
	}

	return isOwner, nil
}
