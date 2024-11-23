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

	var foundedStore sqlc.FindStoreByIDRow
	var errFoundedStore error

	var bhs []sqlc.FindBusinessHourByStoreIDRow
	var errBhs error

	var pms []sqlc.PaymentMethod
	var errPms error

	var productsID []pgtype.UUID
	var errProductsID error

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		foundedStore, errFoundedStore = r.q.FindStoreByID(ctx, idPg)
	}()

	go func() {
		defer wg.Done()
		bhs, errBhs = r.q.FindBusinessHourByStoreID(ctx, idPg)
	}()

	go func() {
		defer wg.Done()
		pms, errPms = r.q.FindPaymentMethodsByStoreID(ctx, idPg)
	}()

	go func() {
		defer wg.Done()
		productsID, errProductsID = r.q.FindProductsIDByStoreID(ctx, idPg)
	}()

	wg.Wait()

	if errFoundedStore != nil {
		return nil, errFoundedStore
	}

	if errBhs != nil {
		return nil, errBhs
	}

	if errPms != nil {
		return nil, errPms
	}

	if errProductsID != nil {
		return nil, errProductsID
	}

	if bhs == nil {
		bhs = []sqlc.FindBusinessHourByStoreIDRow{}
	}
	if pms == nil {
		pms = []sqlc.PaymentMethod{}
	}
	if productsID == nil {
		productsID = []pgtype.UUID{}
	}

	bhsConverted := make([]store.BusinessHours, len(bhs))
	for i, v := range bhs {
		bhsConverted[i] = store.BusinessHours{
			WeekDay:     int(v.Weekday),
			OpeningTime: v.OpenHour,
			ClosingTime: v.ClosingHour,
		}
	}

	pmsConverted := make([]store.PaymentMethod, len(pms))
	for i, v := range pms {
		pmsConverted[i] = store.PaymentMethod(v)
	}

	productsIDConverted := make([]string, len(productsID))
	for i, v := range productsID {
		vConverted, err := converters.UuidToString(v)
		if err != nil {
			return nil, err
		}
		productsIDConverted[i] = *vConverted
	}

	ownerIDConverted, err := converters.UuidToString(foundedStore.OwnerID)
	if err != nil {
		return nil, err
	}

	return &store.Store{
		ID:          id,
		OwnerID:     *ownerIDConverted,
		CNPJ:        foundedStore.Cnpj,
		Name:        foundedStore.Name,
		Description: foundedStore.Description,
		Active:      foundedStore.Active,
		IsOpen:      foundedStore.IsOpen,
		Phone:       foundedStore.Phone,
		Score:       int(foundedStore.Score),
		Type:        store.Type(foundedStore.Type),
		Address: address.Address{
			AddressLine1: foundedStore.AddressLine1,
			AddressLine2: foundedStore.AddressLine2,
			Neighborhood: foundedStore.Neighborhood,
			City:         foundedStore.City,
			State:        foundedStore.State,
			PostalCode:   foundedStore.PostalCode,
			Country:      foundedStore.Country,
			Latitude:     foundedStore.Latitude.String,
			Longitude:    foundedStore.Longitude.String,
		},
		BusinessHours:  bhsConverted,
		PaymentMethods: pmsConverted,
		Products:       productsIDConverted,
	}, nil
}

func (r StoreRepository) FindOwnerByID(ctx context.Context, id string) (*store.Owner, error) {
	idPg, err := converters.StringToUUID(id)
	if err != nil {
		return nil, err
	}

	owner, err := r.q.FindOwnerByID(ctx, idPg)
	if err != nil {
		return nil, err
	}

	ownerIDConverted, err := converters.UuidToString(owner.ID)
	if err != nil {
		return nil, err
	}

	return &store.Owner{
		ID:              *ownerIDConverted,
		SignatureActive: owner.SignatureActive,
	}, nil
}

func (r StoreRepository) BecomeOwner(ctx context.Context, customerID string) error {
	idPg, err := converters.StringToUUID(customerID)
	if err != nil {
		return err
	}

	err = r.q.SetOwner(ctx, sqlc.SetOwnerParams{
		ID:              idPg,
		SignatureActive: false,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r StoreRepository) IsOwner(ctx context.Context, customerID string) (bool, error) {

	ownerID, err := converters.StringToUUID(customerID)
	if err != nil {
		return false, err
	}

	isOwner, err := r.q.IsOwner(ctx, ownerID)
	if err != nil {
		return false, err
	}

	return isOwner, nil
}

func (r StoreRepository) Save(ctx context.Context, storeToSave *store.Store) error {

	storeID, err := converters.StringToUUID(storeToSave.ID)
	if err != nil {
		return err
	}

	ownerID, err := converters.StringToUUID(storeToSave.OwnerID)
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

	profileImage := converters.StringToText(storeToSave.ProfileImage)
	headerImage := converters.StringToText(storeToSave.HeaderImage)

	latitude := converters.StringToText(storeToSave.Address.Latitude)
	longitude := converters.StringToText(storeToSave.Address.Longitude)

	argsSaveStore := sqlc.SaveStoreParams{
		ID:           storeID,
		OwnerID:      ownerID,
		Cnpj:         storeToSave.CNPJ,
		Name:         storeToSave.Name,
		Description:  storeToSave.Description,
		Active:       storeToSave.Active,
		IsOpen:       storeToSave.IsOpen,
		Phone:        storeToSave.Phone,
		Score:        int32(storeToSave.Score),
		Type:         sqlc.StoreType(storeToSave.Type),
		ProfileImage: profileImage,
		HeaderImage:  headerImage,
		AddressLine1: storeToSave.Address.AddressLine1,
		AddressLine2: storeToSave.Address.AddressLine2,
		Neighborhood: storeToSave.Address.Neighborhood,
		City:         storeToSave.Address.City,
		State:        storeToSave.Address.State,
		Country:      storeToSave.Address.Country,
		PostalCode:   storeToSave.Address.PostalCode,
		Latitude:     latitude,
		Longitude:    longitude,
	}
	err = qtx.SaveStore(ctx, argsSaveStore)
	if err != nil {
		return err
	}

	if hasChange(bhSlice, storeToSave.BusinessHours) {
		syncBusinessHour(ctx, qtx, storeID, businessHourRepo, storeToSave)
	}

	if hasChange(pmSlice, storeToSave.PaymentMethods) {
		syncPaymentMethods(ctx, qtx, storeID, paymentMethodRepo, storeToSave)
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
