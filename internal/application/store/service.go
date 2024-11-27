package store

import (
	"context"
	"github.com/google/uuid"
	"github.com/oprimogus/cardapiogo/internal/core/customer"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/adapter"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/adapter/storage"
	"strconv"
)

type Service struct {
	r  store.Repository
	cr customer.Repository
	s  storage.Service
}

type Bucket string

const (
	ProfileBucket Bucket = "cardapiogo-store-profile-image"
	HeaderBucket  Bucket = "cardapiogo-store-header-image"
)

func NewService(r store.Repository, cr customer.Repository, f adapter.Factory) Service {
	return Service{r: r, cr: cr, s: f.NewStorageService()}
}

func (s *Service) NewOwner(ctx context.Context, customerID int) error {
	_, err := s.cr.FindByID(ctx, customerID)
	if err != nil {
		return err
	}

	isOwner, err := s.r.IsOwner(ctx, customerID)
	if err != nil {
		return err
	}
	if !isOwner {
		ow := store.NewOwner(customerID)
		err = s.r.SaveOwner(ctx, ow)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) CreateNewStore(ctx context.Context, ownerID string, params CreateNewStoreDTO) error {
	owID, err := strconv.Atoi(ownerID)
	if err != nil {
		return err
	}

	owner, err := s.r.FindOwnerByID(ctx, owID)
	if err != nil {
		return store.ErrNotOwner
	}

	st, err := owner.NewStore(params.CNPJ, params.Name, params.Description, params.Phone, params.Address, params.Type)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, &st)
}

func (s *Service) Update(ctx context.Context, ownerID string, params UpdateStoreDTO) error {
	owID, err := strconv.Atoi(ownerID)
	if err != nil {
		return err
	}

	st, err := s.r.FindStoreByID(ctx, params.ID)
	if err != nil {
		return err
	}

	if st.OwnerID != owID {
		return store.ErrNotOwner
	}

	err = st.UpdateStoreProfile(params.Name, params.Description, params.Phone, params.Address, params.Type)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s *Service) AddBusinessHour(ctx context.Context, ownerID string, params AddOrDeleteBusinessHourDTO) error {
	owID, err := strconv.Atoi(ownerID)
	if err != nil {
		return err
	}

	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != owID {
		return store.ErrNotOwner
	}

	err = st.AddNewBusinessHour(params.BusinessHours)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s *Service) RemoveBusinessHour(ctx context.Context, ownerID string, params AddOrDeleteBusinessHourDTO) error {
	owID, err := strconv.Atoi(ownerID)
	if err != nil {
		return err
	}

	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != owID {
		return store.ErrNotOwner
	}

	err = st.RemoveBusinessHour(params.BusinessHours)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s *Service) AddPaymentMethod(ctx context.Context, ownerID string, params AddOrDeletePaymentMethodDTO) error {
	owID, err := strconv.Atoi(ownerID)
	if err != nil {
		return err
	}

	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != owID {
		return store.ErrNotOwner
	}

	err = st.AddPaymentMethod(params.PaymentMethods)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s *Service) RemovePaymentMethod(ctx context.Context, ownerID string, params AddOrDeletePaymentMethodDTO) error {
	owID, err := strconv.Atoi(ownerID)
	if err != nil {
		return err
	}

	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != owID {
		return store.ErrNotOwner
	}

	err = st.RemovePaymentMethod(params.PaymentMethods)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s *Service) OpenOrCloseStore(ctx context.Context, ownerID string, params SetOpenStateDTO) error {
	owID, err := strconv.Atoi(ownerID)
	if err != nil {
		return err
	}

	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != owID {
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

func (s *Service) ActiveOrDeactivateStore(ctx context.Context, ownerID string, params SetActiveDTO) error {
	owID, err := strconv.Atoi(ownerID)
	if err != nil {
		return err
	}

	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != owID {
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

func (s *Service) ChangeProfileImage(ctx context.Context, ownerID string, params UploadStoreImageDTO) error {
	owID, err := strconv.Atoi(ownerID)
	if err != nil {
		return err
	}

	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != owID {
		return store.ErrNotOwner
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	objectKey := id.String() + params.Ext
	url, err := s.s.UploadFile(ctx, string(ProfileBucket), objectKey, params.Image)
	if err != nil {
		return err
	}

	st.ProfileImage = url

	return s.r.Save(ctx, st)
}

func (s *Service) ChangeHeaderImage(ctx context.Context, ownerID string, params UploadStoreImageDTO) error {
	owID, err := strconv.Atoi(ownerID)
	if err != nil {
		return err
	}

	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != owID {
		return store.ErrNotOwner
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	objectKey := id.String() + params.Ext
	url, err := s.s.UploadFile(ctx, string(HeaderBucket), objectKey, params.Image)

	st.HeaderImage = url

	return s.r.Save(ctx, st)
}
