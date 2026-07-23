package ownership

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oprimogus/flyfood-api/internal/infra/database"
	"github.com/oprimogus/flyfood-api/internal/infra/database/sqlc"
)

type Repository interface {
	Save(ctx context.Context, owner *Owner) error
	FindByID(ctx context.Context, id uuid.UUID) (*Owner, error)
	FindByCustomerExternalID(ctx context.Context, externalID string) (*Owner, error)
	
}

type repository struct {
    db *pgxpool.Pool
	q  sqlc.Querier
}

func NewRepository(db *database.Postgres) Repository {
	return &repository{db: db.Pool, q: sqlc.New(db.Pool)}
}

func (r *repository) Save(ctx context.Context, owner *Owner) error {
    args := sqlc.SaveOwnerParams{
		ID:              database.ToPgTypeUUID(owner.ID),
		SignatureActive: owner.SignatureActive,
	}
    return r.q.SaveOwner(ctx, args)
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*Owner, error) {
	ow, err := r.q.FindOwnerByID(ctx, database.ToPgTypeUUID(id))
	if err != nil {
		return nil, err
	}
	return &Owner{
		ID:              database.ToUUID(ow.ID),
		SignatureActive: ow.SignatureActive,
		CreatedAt:       ow.CreatedAt.Time,
		UpdatedAt:       ow.UpdatedAt.Time,
		DeletedAt:       ow.DeletedAt.Time,
	}, nil
}

func (r *repository) FindByCustomerExternalID(ctx context.Context, externalID string) (*Owner, error) {
	ow, err := r.q.FindOwnerByCustomerExternalID(ctx, externalID)
	if err != nil {
	    if errors.Is(err, pgx.ErrNoRows) {
	        return nil, ErrNotAnOwner(ctx)
	    }
		return nil, err
	}
	return &Owner{
		ID:              database.ToUUID(ow.ID),
		SignatureActive: ow.SignatureActive,
		CreatedAt:       ow.CreatedAt.Time,
		UpdatedAt:       ow.UpdatedAt.Time,
		DeletedAt:       ow.DeletedAt.Time,
	}, nil
}
