package store

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/database/sqlc"
	"github.com/oprimogus/cardapiogo/pkg/converters"
)

type QueryService struct {
	q sqlc.Querier
	r Repository
}

func NewQueryService(r Repository, q sqlc.Querier) QueryService {
	return QueryService{
		q: q,
		r: r,
	}
}

func (q QueryService) GetStoreByID(ctx context.Context, id string) (*Store, error) {
	return q.r.FindStoreByID(ctx, id)
}

func (q QueryService) GetQueryStoreByID(ctx context.Context, id string) (QueryStore, error) {
	st, err := q.r.FindStoreByID(ctx, id)
	if err != nil {
		return QueryStore{}, err
	}
	return st.ToQueryStore(), nil
}

func (q QueryService) GetStoreByFilter(ctx context.Context, params QueryStoresInput) (*[]QueryStoreList, error) {
	offset := (params.Page - 1) * params.MaxItems
	args := sqlc.FindStoresByFilterParams{
		LimitItems:  int32(params.MaxItems),
		OffsetValue: int32(offset),
	}

	if params.Name != nil {
		args.Name = pgtype.Text{String: *params.Name, Valid: true}
	}

	if params.Score != nil {
		args.Name = pgtype.Text{String: *params.Name, Valid: true}
	}

	if params.IsOpen != nil {
		args.IsOpen = pgtype.Bool{Bool: *params.IsOpen, Valid: true}
	}

	if params.Type != nil {
		typeConv := string(*params.Type)
		args.Type = pgtype.Text{String: typeConv, Valid: true}
	}

	sts, err := q.q.FindStoresByFilter(ctx, args)
	if err != nil {
		return nil, err
	}
	stsOut := make([]QueryStoreList, len(sts))
	for i, st := range sts {
		idConv, err := converters.UuidToString(st.ID)
		if err != nil {
			return nil, err
		}
		stsOut[i] = QueryStoreList{
			ID:           *idConv,
			Name:         st.Name,
			IsOpen:       st.IsOpen,
			Score:        int(st.Score),
			Type:         Type(st.Type),
			ProfileImage: st.ProfileImage.String,
			Neighborhood: st.Neighborhood,
		}
	}

	return &stsOut, nil
}
