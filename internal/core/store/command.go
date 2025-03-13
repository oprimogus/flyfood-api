package store

import (
	"context"
	"github.com/oprimogus/cardapiogo/internal/core/owner"
	"github.com/oprimogus/cardapiogo/internal/core/store/product"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/adapter"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/adapter/storage"
)

type Command struct {
	r  Repository
	or owner.Repository
	s  storage.Service
}

func NewCommand(r Repository, or owner.Repository, f adapter.Factory) Command {
	return Command{r: r, or: or, s: f.NewStorageService()}
}

func (s Command) CreateNewStore(ctx context.Context, ownerID string, params CreateNewStoreDTO) error {

	ow, err := s.or.FindOwnerByID(ctx, ownerID)
	if err != nil {
		return owner.ErrNotOwner
	}

	st, err := NewStore(ow.ID, params.CNPJ, params.Name, params.Description, params.Phone, params.Address, params.Type)
	if err != nil {
		return err
	}

	return s.r.SaveStore(ctx, st)
}

func (s Command) Update(ctx context.Context, ownerID string, params UpdateStoreDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.ID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return owner.ErrNotOwner
	}

	err = st.UpdateStoreProfile(params.Name, params.Description, params.Phone, params.Address, params.Type, params.DeliveryTime)
	if err != nil {
		return err
	}

	return s.r.SaveStore(ctx, st)
}

func (s Command) AddStoreBusinessHour(ctx context.Context, ownerID string, params AddOrDeleteBusinessHourDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return owner.ErrNotOwner
	}

	err = st.AddNewBusinessHour(params.BusinessHours)
	if err != nil {
		return err
	}

	return s.r.SaveStore(ctx, st)
}

func (s Command) RemoveStoreBusinessHour(ctx context.Context, ownerID string, params AddOrDeleteBusinessHourDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return owner.ErrNotOwner
	}

	err = st.RemoveBusinessHour(params.BusinessHours)
	if err != nil {
		return err
	}

	return s.r.SaveStore(ctx, st)
}

func (s Command) AddStorePaymentMethod(ctx context.Context, ownerID string, params AddOrDeletePaymentMethodDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return owner.ErrNotOwner
	}

	err = st.AddPaymentMethod(params.PaymentMethods)
	if err != nil {
		return err
	}

	return s.r.SaveStore(ctx, st)
}

func (s Command) RemoveStorePaymentMethod(ctx context.Context, ownerID string, params AddOrDeletePaymentMethodDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return owner.ErrNotOwner
	}

	err = st.RemovePaymentMethod(params.PaymentMethods)
	if err != nil {
		return err
	}

	return s.r.SaveStore(ctx, st)
}

func (s Command) OpenOrCloseStore(ctx context.Context, ownerID string, params SetOpenStateDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return owner.ErrNotOwner
	}

	if params.IsOpen && !st.IsOpen {
		st.OpenStore()
		return s.r.SaveStore(ctx, st)
	}

	if !params.IsOpen && st.IsOpen {
		st.CloseStore()
		return s.r.SaveStore(ctx, st)
	}

	return nil
}

func (s Command) ActiveOrDeactivateStore(ctx context.Context, ownerID string, params SetActiveDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return owner.ErrNotOwner
	}

	if params.Active && !st.Active {
		st.Activate()
		return s.r.SaveStore(ctx, st)
	}

	if !params.Active && st.Active {
		st.Deactivate()
		return s.r.SaveStore(ctx, st)
	}

	return nil
}

func (s Command) ChangeStoreProfileImage(ctx context.Context, ownerID string, params UploadStoreImageDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return owner.ErrNotOwner
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

	return s.r.SaveStore(ctx, st)
}

func (s Command) ChangeStoreHeaderImage(ctx context.Context, ownerID string, params UploadStoreImageDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return owner.ErrNotOwner
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
	return s.r.SaveStore(ctx, st)
}

func (s Command) NewStoreProduct(ctx context.Context, ownerId string, input product.CreateProductDTO) error {
	isOwnerOf, err := s.or.IsOwnerOf(ctx, ownerId, input.StoreID)
	if err != nil {
		return err
	}

	if !isOwnerOf {
		return owner.ErrNotOwner
	}

	newProduct, err := product.NewProduct(
		input.StoreID,
		input.Name,
		input.Tag,
		input.Description,
		input.SKU,
		input.Price,
		input.Type)
	if err != nil {
		return err
	}

	return s.r.SaveProduct(ctx, newProduct)
}

func (s Command) UpdateProduct(ctx context.Context, ownerId string, input product.UpdateProductDTO) error {
	isOwnerOf, err := s.or.IsOwnerOf(ctx, ownerId, input.StoreID)
	if err != nil {
		return err
	}

	if !isOwnerOf {
		return owner.ErrNotOwner
	}

	p, err := s.r.FindStoreProductByID(ctx, input.ID)
	if err != nil {
		return err
	}

	err = p.Update(input.Type, input.Name, input.Description, input.SKU)
	if err != nil {
		return err
	}

	return s.r.SaveProduct(ctx, p)
}

func (s Command) IncreaseStock(ctx context.Context, ownerId string, input product.ChangeStockProductDTO) error {
	isOwnerOf, err := s.or.IsOwnerOf(ctx, ownerId, input.StoreID)
	if err != nil {
		return err
	}

	if !isOwnerOf {
		return owner.ErrNotOwner
	}

	p, err := s.r.FindStoreProductByID(ctx, input.ID)
	if err != nil {
		return err
	}

	err = p.IncreaseStock(input.Quantity)
	if err != nil {
		return err
	}

	return s.r.SaveProduct(ctx, p)
}

func (s Command) DecreaseStock(ctx context.Context, ownerId string, input product.ChangeStockProductDTO) error {
	isOwnerOf, err := s.or.IsOwnerOf(ctx, ownerId, input.StoreID)
	if err != nil {
		return err
	}

	if !isOwnerOf {
		return owner.ErrNotOwner
	}

	p, err := s.r.FindStoreProductByID(ctx, input.ID)
	if err != nil {
		return err
	}

	err = p.DecreaseStock(input.Quantity)
	if err != nil {
		return err
	}

	return s.r.SaveProduct(ctx, p)
}

func (s Command) ChangeProductImage(ctx context.Context, ownerID string, params product.UploadProductImageDTO) error {
	st, err := s.r.FindStoreByID(ctx, params.StoreID)
	if err != nil {
		return err
	}

	if st.OwnerID != ownerID {
		return owner.ErrNotOwner
	}
	p, err := s.r.FindStoreProductByID(ctx, params.ProductID)
	if err != nil {
		return err
	}

	objectKey := p.ID + "-image" + params.Ext

	if p.Image != "" {
		err = s.s.RemoveFile(ctx, string(ProductBucket), objectKey)
		if err != nil {
			return err
		}
	}

	url, err := s.s.UploadFile(ctx, string(ProductBucket), objectKey, params.Image)
	if err != nil {
		return err
	}

	p.Image = url

	return s.r.SaveProduct(ctx, p)
}
