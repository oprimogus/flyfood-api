package store

// import (
// 	"context"

// 	"github.com/oprimogus/flyfood-api/internal/core/ownership"

// 	"github.com/oprimogus/flyfood-api/internal/infra/services/adapter"
// 	"github.com/oprimogus/flyfood-api/internal/infra/services/adapter/storage"
// 	"github.com/oprimogus/flyfood-api/internal/infra/services/nominatim"
// )

// type Command struct {
// 	r                Repository
// 	ownershipService ownership.Service
// 	s                storage.Service
// }

// func NewCommand(r Repository, f adapter.Factory) Command {
// 	return Command{r: r}
// }

// func (s Command) CreateNewStore(ctx context.Context, ownerID string, params CreateStoreDTO) error {
// 	err := s.ownershipService.EnsureOwnerExists(ctx, ownerID)
// 	if err != nil {
// 		return err
// 	}

// 	st, err := NewStore(ownerID, params.CNPJ, params.Name, params.Description, params.Phone, params.Address, params.Type)
// 	if err != nil {
// 		return err
// 	}

// 	geoData, err := nominatim.Search(ctx, nominatim.Query{
// 		Street:     params.Address.AddressLine1,
// 		City:       params.Address.City,
// 		State:      params.Address.State,
// 		Country:    params.Address.Country,
// 		PostalCode: params.Address.PostalCode,
// 	})
// 	if err == nil {
// 		st.Address.Latitude = geoData[0].Latitude
// 		st.Address.Longitude = geoData[0].Longitude
// 	}

// 	return s.r.Save(ctx, st)
// }

// func (s Command) UpdateStore(ctx context.Context, ownerID string, params UpdateStoreDTO) error {
// 	err := s.ownershipService.ValidateOwnership(ctx, ownerID, params.ID)
// 	if err != nil {
// 		return err
// 	}

// 	st, err := s.r.FindByID(ctx, params.ID)
// 	if err != nil {
// 		return err
// 	}

// 	err = st.UpdateStoreProfile(params.Name, params.Description, params.Phone, params.Address, params.Type, params.DeliveryTime)
// 	if err != nil {
// 		return err
// 	}

// 	return s.r.Save(ctx, st)
// }

// func (s Command) AddStoreBusinessHour(ctx context.Context, ownerID string, params AddOrDeleteBusinessHourDTO) error {
// 	err := s.ownershipService.ValidateOwnership(ctx, ownerID, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	st, err := s.r.FindByID(ctx, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	err = st.AddNewBusinessHour(params.BusinessHours)
// 	if err != nil {
// 		return err
// 	}

// 	return s.r.Save(ctx, st)
// }

// func (s Command) RemoveStoreBusinessHour(ctx context.Context, ownerID string, params AddOrDeleteBusinessHourDTO) error {
// 	err := s.ownershipService.ValidateOwnership(ctx, ownerID, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	st, err := s.r.FindByID(ctx, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	err = st.RemoveBusinessHour(params.BusinessHours)
// 	if err != nil {
// 		return err
// 	}

// 	return s.r.Save(ctx, st)
// }

// func (s Command) AddStorePaymentMethod(ctx context.Context, ownerID string, params AddOrDeletePaymentMethodDTO) error {
// 	err := s.ownershipService.ValidateOwnership(ctx, ownerID, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	st, err := s.r.FindByID(ctx, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	err = st.AddPaymentMethod(params.PaymentMethods)
// 	if err != nil {
// 		return err
// 	}

// 	return s.r.Save(ctx, st)
// }

// func (s Command) RemoveStorePaymentMethod(ctx context.Context, ownerID string, params AddOrDeletePaymentMethodDTO) error {
// 	err := s.ownershipService.ValidateOwnership(ctx, ownerID, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	st, err := s.r.FindByID(ctx, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	err = st.RemovePaymentMethod(params.PaymentMethods)
// 	if err != nil {
// 		return err
// 	}

// 	return s.r.Save(ctx, st)
// }

// func (s Command) OpenOrCloseStore(ctx context.Context, ownerID string, params SetOpenStateDTO) error {
// 	err := s.ownershipService.ValidateOwnership(ctx, ownerID, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	st, err := s.r.FindByID(ctx, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	if params.IsOpen && !st.IsOpen {
// 		st.OpenStore()
// 		return s.r.Save(ctx, st)
// 	}

// 	if !params.IsOpen && st.IsOpen {
// 		st.CloseStore()
// 		return s.r.Save(ctx, st)
// 	}

// 	return nil
// }

// func (s Command) ActiveOrDeactivateStore(ctx context.Context, ownerID string, params SetActiveDTO) error {
// 	err := s.ownershipService.ValidateOwnership(ctx, ownerID, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	st, err := s.r.FindByID(ctx, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	if params.Active && !st.Active {
// 		st.Activate()
// 		return s.r.Save(ctx, st)
// 	}

// 	if !params.Active && st.Active {
// 		st.Deactivate()
// 		return s.r.Save(ctx, st)
// 	}

// 	return nil
// }

// func (s Command) ChangeStoreProfileImage(ctx context.Context, ownerID string, params UploadStoreImageDTO) error {
// 	err := s.ownershipService.ValidateOwnership(ctx, ownerID, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	st, err := s.r.FindByID(ctx, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	objectKey := st.ID + "-profile-image" + params.Ext

// 	if st.ProfileImage != "" {
// 		err = s.s.RemoveFile(ctx, string(ProfileBucket), objectKey)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	url, err := s.s.UploadFile(ctx, string(ProfileBucket), objectKey, params.Image)
// 	if err != nil {
// 		return err
// 	}

// 	st.ProfileImage = url

// 	return s.r.Save(ctx, st)
// }

// func (s Command) ChangeStoreHeaderImage(ctx context.Context, ownerID string, params UploadStoreImageDTO) error {
// 	err := s.ownershipService.ValidateOwnership(ctx, ownerID, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	st, err := s.r.FindByID(ctx, params.StoreID)
// 	if err != nil {
// 		return err
// 	}

// 	objectKey := st.ID + "-header-image" + params.Ext

// 	if st.HeaderImage != "" {
// 		err = s.s.RemoveFile(ctx, string(HeaderBucket), objectKey)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	url, err := s.s.UploadFile(ctx, string(HeaderBucket), objectKey, params.Image)
// 	if err != nil {
// 		return err
// 	}

// 	st.HeaderImage = url
// 	return s.r.Save(ctx, st)
// }
