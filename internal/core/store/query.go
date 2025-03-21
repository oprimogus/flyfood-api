package store

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oprimogus/flyfood-api/internal/core"
	"github.com/oprimogus/flyfood-api/internal/core/store/product"
	"github.com/oprimogus/flyfood-api/internal/infrastructure/database/sqlc"
	"github.com/oprimogus/flyfood-api/pkg/converters"
)

type Query struct {
	q sqlc.Querier
	r Repository
}

func NewQueryService(r Repository, q sqlc.Querier) Query {
	return Query{
		q: q,
		r: r,
	}
}

func (q Query) GetStoreByID(ctx context.Context, id string) (Store, error) {
	return q.r.FindStoreByID(ctx, id)
}

func (q Query) GetQueryStoreByID(ctx context.Context, id string) (QueryStore, error) {
	st, err := q.r.FindStoreByID(ctx, id)
	if err != nil {
		return QueryStore{}, err
	}

	stIDConv, err := converters.StringToUUID(id)
	if err != nil {
		return QueryStore{}, err
	}

	psSqlc, err := q.q.FindProductsByStoreID(ctx, stIDConv)
	if err != nil {
		return QueryStore{}, err
	}

	ps := make([]product.ProductDTO, len(psSqlc))
	for i, p := range psSqlc {
		pID, err := converters.UuidToString(p.ID)
		if err != nil {
			return QueryStore{}, err
		}

		ps[i] = product.ProductDTO{
			ID:               *pID,
			SKU:              p.Sku.String,
			PromoActive:      p.PromoActive,
			Type:             product.Type(p.Type),
			Tag:              p.Tag,
			Name:             p.Name,
			Description:      p.Description,
			Score:            int(p.Score),
			Image:            p.ImageUrl.String,
			Details:          nil,
			Price:            int(p.Price),
			PromotionalPrice: converters.Int4ToInt(p.PromotionalPrice),
		}
	}

	return st.ToQueryStore(ps), nil
}

func (q Query) GetQueryOwnerStoreByID(ctx context.Context, id string) (QueryOwnerStore, error) {
	st, err := q.r.FindStoreByID(ctx, id)
	if err != nil {
		return QueryOwnerStore{}, err
	}

	stIDConv, err := converters.StringToUUID(id)
	if err != nil {
		return QueryOwnerStore{}, err
	}

	psSqlc, err := q.q.FindProductsByStoreID(ctx, stIDConv)
	if err != nil {
		return QueryOwnerStore{}, err
	}

	ps := make([]product.Product, len(psSqlc))
	for i, p := range psSqlc {
		pID, err := converters.UuidToString(p.ID)
		if err != nil {
			return QueryOwnerStore{}, err
		}

		ps[i] = product.Product{
			ID:               *pID,
			StoreID:          st.ID,
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
			Details:          nil,
			Price:            int(p.Price),
			PromotionalPrice: converters.Int4ToInt(p.PromotionalPrice),
		}
	}

	return st.ToQueryOwnerStore(ps), nil
}

func (q Query) GetStoreByFilter(ctx context.Context, params QueryStoresInput) (core.Pagination[QueryStoreList], error) {
	offset := (params.Page - 1) * params.MaxItems
	args := sqlc.FindStoresByFilterParams{
		LimitItems:  int32(params.MaxItems),
		OffsetValue: int32(offset),
	}

	if params.City != nil {
		args.City = *params.City
	}

	if params.Name != nil {
		args.Name = *params.Name
	}

	if params.Score != nil {
		args.Score = int32(*params.Score)
	}

	if params.IsOpen != nil {
		args.IsOpen = pgtype.Bool{Bool: *params.IsOpen, Valid: true}
	}

	if params.Type != nil {
		args.Type = string(*params.Type)
	}

	sts, err := q.q.FindStoresByFilter(ctx, args)
	if err != nil {
		return core.Pagination[QueryStoreList]{}, err
	}
	stsOut := make([]QueryStoreList, len(sts))
	for i, st := range sts {
		idConv, err := converters.UuidToString(st.ID)
		if err != nil {
			return core.Pagination[QueryStoreList]{}, err
		}
		stsOut[i] = QueryStoreList{
			ID:           *idConv,
			Name:         st.Name,
			IsOpen:       st.IsOpen,
			Score:        int(st.Score),
			Type:         Type(st.Type),
			DeliveryTime: converters.Int4ToInt(st.DeliveryTime),
			ProfileImage: st.ProfileImage.String,
			Neighborhood: st.Neighborhood,
			Latitude:     st.Latitude.String,
			Longitude:    st.Longitude.String,
		}
	}
	var totalItems int
	if len(sts) > 0 {
		totalItems = int(sts[0].TotalItems)
	} else {
		totalItems = 0
	}

	return core.Paginate(stsOut, params.Page, params.MaxItems, totalItems), nil
}

func (q Query) GetOwnerStores(ctx context.Context, ownerID string) ([]QueryOwnerStoreList, error) {
	sts, err := q.q.FindOwnerStores(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	stsOut := make([]QueryOwnerStoreList, len(sts))
	for i, st := range sts {
		idConv, err := converters.UuidToString(st.ID)
		if err != nil {
			return nil, err
		}
		stsOut[i] = QueryOwnerStoreList{
			ID:           *idConv,
			Name:         st.Name,
			Active:       st.Active,
			IsOpen:       st.IsOpen,
			Score:        int(st.Score),
			Type:         Type(st.Type),
			DeliveryTime: converters.Int4ToInt(st.DeliveryTime),
			ProfileImage: st.ProfileImage.String,
			City:         st.City,
			State:        st.State,
			Country:      st.Country,
		}
	}

	return stsOut, nil
}
