package store

import (
	"context"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/adapter"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/adapter/storage"
)

type Service struct {
	r store.Repository
	s storage.Service
}

type Bucket string

const (
	ProfileBucket Bucket = "cardapiogo-store-profile-image"
	HeaderBucket  Bucket = "cardapiogo-store-header-image"
)

var buckets = map[string]string{
	"profile": "cardapiogo-store-profile-image",
	"header":  "cardapiogo-store-header-image",
}

func NewService(r store.Repository, f adapter.Factory) Service {
	return Service{r: r, s: f.NewStorageService()}
}

func (s *Service) CreateNewStore(ctx context.Context, params CreateNewStoreDTO) error {
	owner, err := s.r.FindOwnerByID(ctx, params.OwnerID)
	if err != nil {
		return store.ErrNotOwner
	}

	st, err := owner.NewStore(params.CNPJ, params.Name, params.Description, params.Phone, params.Address, params.Type)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, &st)
}

func (s *Service) Update(ctx context.Context, params UpdateStoreDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.ID)
	if err != nil {
		return err
	}

	if st.OwnerID != params.OwnerID {
		return store.ErrNotOwner
	}

	err = st.UpdateStoreProfile(params.Name, params.Description, params.Phone, params.Address, params.Type)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s *Service) AddBusinessHour(ctx context.Context, params AddOrDeleteBusinessHourDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != params.OwnerID {
		return store.ErrNotOwner
	}

	err = st.AddNewBusinessHour(params.BusinessHours)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s *Service) RemoveBusinessHour(ctx context.Context, params AddOrDeleteBusinessHourDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != params.OwnerID {
		return store.ErrNotOwner
	}

	err = st.RemoveBusinessHour(params.BusinessHours)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s *Service) AddPaymentMethod(ctx context.Context, params AddOrDeletePaymentMethodDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != params.OwnerID {
		return store.ErrNotOwner
	}

	err = st.AddPaymentMethod(params.PaymentMethods)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s *Service) RemovePaymentMethod(ctx context.Context, params AddOrDeletePaymentMethodDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != params.OwnerID {
		return store.ErrNotOwner
	}

	err = st.RemovePaymentMethod(params.PaymentMethods)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s *Service) OpenOrCloseStore(ctx context.Context, params SetOpenStateDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != params.OwnerID {
		return store.ErrNotOwner
	}

	if params.IsOpen && !st.IsOpen {
		st.OpenStore()
		return s.r.Save(ctx, st)
	}

	if !params.IsOpen && st.IsOpen {
		st.CloseStore()
		return s.r.Save(ctx, st)
	}

	return nil
}

func (s *Service) ActiveOrDeactivateStore(ctx context.Context, params SetActiveDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != params.OwnerID {
		return store.ErrNotOwner
	}

	if params.Active && !st.Active {
		st.Activate()
		return s.r.Save(ctx, st)
	}

	if !params.Active && st.Active {
		st.Deactivate()
		return s.r.Save(ctx, st)
	}

	return nil
}

func (s *Service) ChangeProfileImage(ctx context.Context, params ChangeProfileImageDTO) error {

	url, err := s.s.UploadFile(ctx)

	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != params.OwnerID {
		return store.ErrNotOwner
	}

	return nil
}
