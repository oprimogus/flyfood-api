package persistence

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oprimogus/cardapiogo/internal/core/address"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	postgresDB "github.com/oprimogus/cardapiogo/internal/infrastructure/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/database/sqlc"
	"github.com/oprimogus/cardapiogo/pkg/converters"
	"log/slog"
	"strconv"
	"strings"
	"sync"
)

type StoreRepository struct {
	db *pgxpool.Pool
	q  sqlc.Querier
}

func NewStoreRepository(db *postgresDB.Database) StoreRepository {
	return StoreRepository{db: db.GetDB(), q: sqlc.New(db.GetDB())}
}

func (r StoreRepository) FindStoreByID(ctx context.Context, id string) (*store.Store, error) {
	idPg, err := converters.StringToUUID(id)
	if err != nil {
		return nil, err
	}

	st, err := r.q.FindStoreByID(ctx, idPg)
	if err != nil {
		return nil, err
	}

	bhs := make([]store.BusinessHours, len(st.BusinessHours))
	for i, bh := range st.BusinessHours {
		arrayText := strings.Split(strings.Trim(bh, "{}"), " ")
		if len(arrayText) == 0 {
			continue
		}
		weekDay, err := strconv.Atoi(arrayText[0])
		if err != nil {
			return nil, err
		}
		bhs[i] = store.BusinessHours{
			WeekDay:     weekDay,
			OpeningTime: arrayText[1],
			ClosingTime: arrayText[2],
		}
	}

	pms := make([]store.PaymentMethod, len(st.PaymentMethods))
	for i, pm := range st.PaymentMethods {
		arrayText := strings.Trim(pm, "{}")
		if arrayText == "" {
			continue
		}
		pms[i] = store.PaymentMethod(arrayText)
	}

	return &store.Store{
		ID:          id,
		OwnerID:     int(st.OwnerID),
		CNPJ:        st.Cnpj,
		Name:        st.Name,
		Description: st.Description,
		Active:      st.Active,
		IsOpen:      st.IsOpen,
		Phone:       st.Phone,
		Score:       int(st.Score),
		Type:        store.Type(st.Type),
		Address: address.Address{
			AddressLine1: st.AddressLine1,
			AddressLine2: st.AddressLine2,
			Neighborhood: st.Neighborhood,
			City:         st.City,
			State:        st.State,
			PostalCode:   st.PostalCode,
			Country:      st.Country,
			Latitude:     st.Latitude.String,
			Longitude:    st.Longitude.String,
		},
		BusinessHours:  bhs,
		PaymentMethods: pms,
		Products:       []string{},
	}, nil

}

func (r StoreRepository) FindOwnerByID(ctx context.Context, id int) (*store.Owner, error) {
	owner, err := r.q.FindOwnerByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	return &store.Owner{
		ID:              int(owner.ID),
		SignatureActive: owner.SignatureActive,
	}, nil
}

func (r StoreRepository) SaveOwner(ctx context.Context, ow store.Owner) error {
	err := r.q.SaveOwner(ctx, sqlc.SaveOwnerParams{
		ID:              int64(ow.ID),
		SignatureActive: ow.SignatureActive,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r StoreRepository) IsOwner(ctx context.Context, customerID int) (bool, error) {
	isOwner, err := r.q.IsOwner(ctx, int64(customerID))
	if err != nil {
		return false, err
	}

	return isOwner, nil
}

func (r StoreRepository) Save(ctx context.Context, st *store.Store) error {

	storeID, err := converters.StringToUUID(st.ID)
	if err != nil {
		return err
	}

	var businessHourRepo []sqlc.FindBusinessHourByStoreIDRow
	var errBusinessHourRepo error

	var paymentMethodRepo []sqlc.PaymentMethod
	var errPaymentMethodRepo error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		businessHourRepo, errBusinessHourRepo = r.q.FindBusinessHourByStoreID(ctx, storeID)
		if businessHourRepo == nil {
			businessHourRepo = []sqlc.FindBusinessHourByStoreIDRow{}
		}
	}()

	go func() {
		defer wg.Done()
		paymentMethodRepo, errPaymentMethodRepo = r.q.FindPaymentMethodsByStoreID(ctx, storeID)
		if paymentMethodRepo == nil {
			paymentMethodRepo = []sqlc.PaymentMethod{}
		}
	}()

	wg.Wait()

	if errBusinessHourRepo != nil {
		return errBusinessHourRepo
	}

	if errPaymentMethodRepo != nil {
		return errPaymentMethodRepo
	}

	bhSlice := make([]store.BusinessHours, len(businessHourRepo))
	for i, v := range businessHourRepo {
		bhSlice[i] = store.BusinessHours{
			WeekDay:     int(v.Weekday),
			OpeningTime: v.OpenHour,
			ClosingTime: v.ClosingHour,
		}
	}

	pmSlice := make([]store.PaymentMethod, len(paymentMethodRepo))
	for i, v := range paymentMethodRepo {
		pmSlice[i] = store.PaymentMethod(v)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qtx := sqlc.New(tx)

	profileImage := converters.StringToText(st.ProfileImage)
	headerImage := converters.StringToText(st.HeaderImage)

	latitude := converters.StringToText(st.Address.Latitude)
	longitude := converters.StringToText(st.Address.Longitude)

	argsSaveStore := sqlc.SaveStoreParams{
		ID:           storeID,
		OwnerID:      int64(st.OwnerID),
		Cnpj:         st.CNPJ,
		Name:         st.Name,
		Description:  st.Description,
		Active:       st.Active,
		IsOpen:       st.IsOpen,
		Phone:        st.Phone,
		Score:        int32(st.Score),
		Type:         sqlc.StoreType(st.Type),
		ProfileImage: profileImage,
		HeaderImage:  headerImage,
		AddressLine1: st.Address.AddressLine1,
		AddressLine2: st.Address.AddressLine2,
		Neighborhood: st.Address.Neighborhood,
		City:         st.Address.City,
		State:        st.Address.State,
		Country:      st.Address.Country,
		PostalCode:   st.Address.PostalCode,
		Latitude:     latitude,
		Longitude:    longitude,
	}
	err = qtx.SaveStore(ctx, argsSaveStore)
	if err != nil {
		return err
	}

	if hasChange(bhSlice, st.BusinessHours) {
		syncBusinessHour(ctx, qtx, storeID, businessHourRepo, st)
	}

	if hasChange(pmSlice, st.PaymentMethods) {
		syncPaymentMethods(ctx, qtx, storeID, paymentMethodRepo, st)
	}

	return tx.Commit(ctx)
}

func syncPaymentMethods(ctx context.Context, qtx *sqlc.Queries, id pgtype.UUID,
	repo []sqlc.PaymentMethod, storeToSave *store.Store) {
	if repo == nil {
		repo = []sqlc.PaymentMethod{}
	}
	argsDeletePaymentMethods := make([]sqlc.DeleteStorePaymentMethodsParams, len(repo))
	for i, v := range repo {
		argsDeletePaymentMethods[i] = sqlc.DeleteStorePaymentMethodsParams{
			ID:            id,
			PaymentMethod: v,
		}
	}
	resultDeletePaymentMethods := qtx.DeleteStorePaymentMethods(ctx, argsDeletePaymentMethods)
	resultDeletePaymentMethods.Exec(func(i int, err error) {
		if err != nil {
			slog.ErrorContext(ctx, "fail to delete payment method",
				"index", i,
				"error", err.Error())
		}
	})
	if err := resultDeletePaymentMethods.Close(); err != nil {
		slog.ErrorContext(ctx, "Failed to close batch", "error", err.Error())
	}

	argsSavePaymentMethods := make([]sqlc.SaveStorePaymentMethodsParams, len(storeToSave.PaymentMethods))
	for i, v := range storeToSave.PaymentMethods {
		argsSavePaymentMethods[i] = sqlc.SaveStorePaymentMethodsParams{
			ID:            id,
			PaymentMethod: sqlc.PaymentMethod(v),
		}
	}
	resultSavePaymentMethods := qtx.SaveStorePaymentMethods(ctx, argsSavePaymentMethods)
	resultSavePaymentMethods.Exec(func(i int, err error) {
		if err != nil {
			slog.ErrorContext(ctx, "fail to delete payment method",
				"index", i,
				"error", err.Error())
		}
	})
	if err := resultSavePaymentMethods.Close(); err != nil {
		slog.ErrorContext(ctx, "Failed to close batch", "error", err.Error())
	}

}

func syncBusinessHour(ctx context.Context, qtx sqlc.Querier,
	storeID pgtype.UUID, sliceRepo []sqlc.FindBusinessHourByStoreIDRow, storeToSave *store.Store) {
	if sliceRepo == nil {
		sliceRepo = []sqlc.FindBusinessHourByStoreIDRow{}
	}
	argsDeleteBusinessHour := make([]sqlc.DeleteStoreBusinessHourParams, len(sliceRepo))
	for i, v := range sliceRepo {
		argsDeleteBusinessHour[i] = sqlc.DeleteStoreBusinessHourParams{
			ID:          storeID,
			Weekday:     v.Weekday,
			OpenHour:    v.OpenHour,
			ClosingHour: v.ClosingHour,
		}
	}

	resultDeleteBhs := qtx.DeleteStoreBusinessHour(ctx, argsDeleteBusinessHour)
	resultDeleteBhs.Exec(func(i int, err error) {
		if err != nil {
			slog.Error("failed to delete business hour",
				"index", i,
				"error", err.Error())
		}
	})
	if err := resultDeleteBhs.Close(); err != nil {
		slog.ErrorContext(ctx, "Failed to close batch", "error", err.Error())
	}

	argsAddBusinessHour := make([]sqlc.SaveStoreBusinessHourParams, len(storeToSave.BusinessHours))
	for i, v := range storeToSave.BusinessHours {
		argsAddBusinessHour[i] = sqlc.SaveStoreBusinessHourParams{
			ID:          storeID,
			Weekday:     int32(v.WeekDay),
			OpenHour:    v.OpeningTime,
			ClosingHour: v.ClosingTime,
		}
	}

	resultSaveBhs := qtx.SaveStoreBusinessHour(ctx, argsAddBusinessHour)
	resultSaveBhs.Exec(func(i int, err error) {
		if err != nil {
			slog.Error("failed to save business hour",
				"index", i,
				"error", err.Error())
		}
	})
	if err := resultSaveBhs.Close(); err != nil {
		slog.ErrorContext(ctx, "Failed to close batch", "error", err.Error())
	}

}

func hasChange[T comparable](persistence []T, aggregate []T) bool {
	if persistence == nil {
		persistence = []T{}
	}
	if aggregate == nil {
		aggregate = []T{}
	}

	if len(persistence) != len(aggregate) {
		return true
	}
	mapRepo := make(map[T]bool, len(persistence))
	for _, v := range persistence {
		mapRepo[v] = true
	}

	mapAggregate := make(map[T]bool, len(aggregate))
	for _, v := range aggregate {
		mapAggregate[v] = true
	}

	for bh := range mapRepo {
		if !mapAggregate[bh] {
			return true
		}
	}

	for bh := range mapAggregate {
		if !mapRepo[bh] {
			return true
		}
	}

	return false
}
