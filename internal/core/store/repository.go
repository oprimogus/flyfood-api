package store

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oprimogus/flyfood-api/internal/core/address"
	"github.com/oprimogus/flyfood-api/internal/infra/database"
	"github.com/oprimogus/flyfood-api/internal/infra/database/sqlc"
	"github.com/oprimogus/flyfood-api/pkg/xerrors"
)

type Repository interface {
	FindByID(ctx context.Context, id string) (*Store, error)
	Save(ctx context.Context, s *Store) error
}

type repository struct {
	db *pgxpool.Pool
	q  sqlc.Querier
}

func NewRepository(db *database.Postgres) repository {
	return repository{db: db.Pool, q: sqlc.New(db.Pool)}
}

func (r repository) FindByID(ctx context.Context, id string) (*Store, error) {
	idPg, err := database.StringToUUID(id)
	if err != nil {
		return nil, xerrors.NewWithContext(ctx, fmt.Errorf("cannot generate uuid v7: %w", err))
	}

	storeSqlc, err := r.q.FindStoreByID(ctx, idPg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, xerrors.NewWithContext(ctx, ErrStoreNotFound).WithStatus(http.StatusNotFound)
		}
		return nil, xerrors.NewWithContext(ctx, err)
	}

	bhsSqlc, err := r.q.FindStoreBusinessHoursByStoreId(ctx, idPg)
	if err != nil {
		return nil, xerrors.NewWithContext(ctx, err)
	}

	pmSqlc, err := r.q.FindStorePaymentMethodsByStoreId(ctx, idPg)
	if err != nil {
		return nil, xerrors.NewWithContext(ctx, err)
	}

	st := &Store{
		ID:          storeSqlc.ID.Bytes,
		OwnerID:     storeSqlc.OwnerID.Bytes,
		CNPJ:        storeSqlc.Cnpj,
		Name:        storeSqlc.Name,
		Description: storeSqlc.Description.String,
		Active:      storeSqlc.Active,
		IsOpen:      storeSqlc.IsOpen,
		Phone:       storeSqlc.Phone,
		Score:       int(storeSqlc.Score),
		Address: address.Address{
			AddressLine1: storeSqlc.AddressLine1.String,
			AddressLine2: storeSqlc.AddressLine2.String,
			Neighborhood: storeSqlc.Neighborhood.String,
			City:         storeSqlc.City.String,
			State:        storeSqlc.State.String,
			PostalCode:   storeSqlc.PostalCode.String,
			Country:      storeSqlc.Country.String,
			Latitude:     storeSqlc.Latitude,
			Longitude:    storeSqlc.Longitude,
		},
		Type:           Type(storeSqlc.Type),
		ProfileImage:   storeSqlc.ProfileImage.String,
		HeaderImage:    storeSqlc.HeaderImage.String,
	}

	bhs := make([]BusinessHours, len(bhsSqlc))
	for i, bh := range bhsSqlc{
		op, err := NewMinutesOfDayFromInt(int(bh.OpenHour))
		if err != nil {
			return nil, xerrors.NewWithContext(ctx, err)
		}
		cl, err := NewMinutesOfDayFromInt(int(bh.ClosingHour))
		if err != nil {
			return nil, xerrors.NewWithContext(ctx, err)
		}
		bhs[i] = BusinessHours{
			WeekDay:     int(bh.Weekday),
			OpeningTime: op,
			ClosingTime: cl,
		}
	}
	st.BusinessHours = bhs

	pms := make([]PaymentMethod, len(pmSqlc))
	for i, pm := range pmSqlc {
		pms[i] = PaymentMethod(pm)
	}
	st.PaymentMethods = pms

	return st, nil
}

func (r repository) Save(ctx context.Context, st *Store) error {
    transaction, err := r.db.Begin(ctx)
    if err != nil {
        return xerrors.NewWithContext(ctx, err)
    }
    defer func() {
		_ = transaction.Rollback(ctx)
	}()
    
    tq := sqlc.New(transaction)

    argsSaveAddress := sqlc.SaveAddressParams{
    	ID:           database.ToPgTypeUUID(st.Address.ID),
    	Name:         st.Address.Name,
    	AddressLine1: st.Address.AddressLine1,
    	AddressLine2: st.Address.AddressLine2,
    	Neighborhood: st.Address.Neighborhood,
    	City:         st.Address.City,
    	State:        st.Address.State,
    	PostalCode:   st.Address.PostalCode,
    	Country:      st.Address.Country,
        Latitude:     st.Address.Latitude,
        Longitude:    st.Address.Longitude,
    }
    err = tq.SaveAddress(ctx, argsSaveAddress)
    if err != nil {
        return xerrors.NewWithContext(ctx, err)
    }

    argsSaveStore := sqlc.SaveStoreParams{
        ID:          database.ToPgTypeUUID(st.ID),
        OwnerID:     database.ToPgTypeUUID(st.OwnerID),
        AddressID:   database.ToPgTypeUUID(st.Address.ID),
        Cnpj:        st.CNPJ,
        Name:        st.Name,
        Description: database.StringToText(st.Description),
        Active:      st.Active,
        IsOpen:      st.IsOpen,
        Phone:       st.Phone,
        Score:       int32(st.Score),
        Type:        int32(st.Type),
    }
    err = tq.SaveStore(ctx, argsSaveStore)
    if err != nil {
        _ = transaction.Rollback(ctx)
        return xerrors.NewWithContext(ctx, err)
    }

    weekdays := make([]int16, len(st.BusinessHours))
    for i, bh := range st.BusinessHours {
        weekdays[i] = int16(bh.WeekDay)

        if err := tq.UpsertBusinessHour(ctx, sqlc.UpsertBusinessHourParams{
            StoreID:     database.ToPgTypeUUID(st.ID),
            Weekday:     int16(bh.WeekDay),
            OpenHour:    int16(bh.OpeningTime),
            ClosingHour: int16(bh.ClosingTime),
        }); err != nil {
            return xerrors.NewWithContext(ctx, err)
        }
    }

    if err = tq.DeleteBusinessHoursNotIn(ctx, sqlc.DeleteBusinessHoursNotInParams{
    	StoreID:       database.ToPgTypeUUID(st.ID),
    	BusinessHours: weekdays,
    }); err != nil {
        return xerrors.NewWithContext(ctx, err)
    }

    pms := make([]int32, len(st.PaymentMethods))
    for i, pm := range st.PaymentMethods {
        pms[i] = int32(pm)
        if err = tq.UpsertPaymentMethod(ctx, sqlc.UpsertPaymentMethodParams{
        	StoreID:       database.ToPgTypeUUID(st.ID),
        	PaymentMethod: pms[i],
        }); err != nil {
            return xerrors.NewWithContext(ctx, err)
        }
    }

    if err = tq.DeletePaymentMethodsNotIn(ctx, sqlc.DeletePaymentMethodsNotInParams{
    	StoreID:          database.ToPgTypeUUID(st.ID),
    	PaymentMethods:   pms,
    }); err != nil {
        return xerrors.NewWithContext(ctx, err)
    }
    
	return transaction.Commit(ctx)
}