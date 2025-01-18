package store

import (
	"context"
	"github.com/oprimogus/cardapiogo/internal/core/customer"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/adapter"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/adapter/storage"
)

type CommandService struct {
	r  Repository
	cr customer.Repository
	s  storage.Service
}

func GetBuckets() []storage.Bucket {
	return []storage.Bucket{
		ProfileBucket,
		HeaderBucket,
	}
}

const (
	ProfileBucket storage.Bucket = "cardapiogo-store-profile-image"
	HeaderBucket  storage.Bucket = "cardapiogo-store-header-image"
)

func NewCommandService(r Repository, cr customer.Repository, f adapter.Factory) CommandService {
	return CommandService{r: r, cr: cr, s: f.NewStorageService()}
}

func (s CommandService) NewOwner(ctx context.Context, customerID string) error {
	_, err := s.cr.FindByID(ctx, customerID)
	if err != nil {
		return err
	}

	isOwner, err := s.r.IsOwner(ctx, customerID)
	if err != nil {
		return err
	}
	if !isOwner {
		ow := NewOwner(customerID)
		err = s.r.SaveOwner(ctx, ow)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s CommandService) CreateNewStore(ctx context.Context, ownerID string, params CreateNewStoreDTO) error {

	owner, err := s.r.FindOwnerByID(ctx, ownerID)
	if err != nil {
		return err
	}

	if owner.ID != ownerID {
		return ErrNotOwner
	}

	st, err := owner.NewStore(params.CNPJ, params.Name, params.Description, params.Phone, params.Address, params.Type)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, &st)
}

func (s CommandService) Update(ctx context.Context, ownerID string, params UpdateStoreDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.ID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return ErrNotOwner
	}

	err = st.UpdateStoreProfile(params.Name, params.Description, params.Phone, params.Address, params.Type)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s CommandService) AddBusinessHour(ctx context.Context, ownerID string, params AddOrDeleteBusinessHourDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return ErrNotOwner
	}

	err = st.AddNewBusinessHour(params.BusinessHours)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s CommandService) RemoveBusinessHour(ctx context.Context, ownerID string, params AddOrDeleteBusinessHourDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return ErrNotOwner
	}

	err = st.RemoveBusinessHour(params.BusinessHours)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s CommandService) AddPaymentMethod(ctx context.Context, ownerID string, params AddOrDeletePaymentMethodDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return ErrNotOwner
	}

	err = st.AddPaymentMethod(params.PaymentMethods)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s CommandService) RemovePaymentMethod(ctx context.Context, ownerID string, params AddOrDeletePaymentMethodDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return ErrNotOwner
	}

	err = st.RemovePaymentMethod(params.PaymentMethods)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, st)
}

func (s CommandService) OpenOrCloseStore(ctx context.Context, ownerID string, params SetOpenStateDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return ErrNotOwner
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

func (s CommandService) ActiveOrDeactivateStore(ctx context.Context, ownerID string, params SetActiveDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return ErrNotOwner
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

func (s CommandService) ChangeProfileImage(ctx context.Context, ownerID string, params UploadStoreImageDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return ErrNotOwner
	}

	objectKey := st.ID + "-profile-image" + params.Ext

	if st.ProfileImage != "" {
		err = s.s.RemoveFile(ctx, string(ProfileBucket), objectKey)
		if err != nil {
			return err
		}
	}

	url, err := s.s.UploadFile(ctx, string(ProfileBucket), objectKey, params.Image)
	if err != nil {
		return err
	}

	st.ProfileImage = url

	return s.r.Save(ctx, st)
}

func (s CommandService) ChangeHeaderImage(ctx context.Context, ownerID string, params UploadStoreImageDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return ErrNotOwner
	}

	objectKey := st.ID + "-header-image" + params.Ext

	if st.HeaderImage != "" {
		err = s.s.RemoveFile(ctx, string(HeaderBucket), objectKey)
		if err != nil {
			return err
		}
	}

	url, err := s.s.UploadFile(ctx, string(HeaderBucket), objectKey, params.Image)
	if err != nil {
		return err
	}

	st.HeaderImage = url
	return s.r.Save(ctx, st)
}
