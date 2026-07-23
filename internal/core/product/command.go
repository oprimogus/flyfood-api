package product

// import (
// 	"context"

// 	"github.com/oprimogus/flyfood-api/internal/core/ownership"
// 	"github.com/oprimogus/flyfood-api/internal/infra/services/adapter/storage"
// 	"github.com/oprimogus/flyfood-api/pkg/xerrors"
// )

// type Command struct {
// 	r                Repository
// 	ownershipService ownership.Service
// 	s                storage.Service
// }

// func NewCommand(r Repository, ownershipService ownership.Service, s storage.Service) Command {
// 	return Command{
// 		r:                r,
// 		ownershipService: ownershipService,
// 		s:                s,
// 	}
// }

// func (s Command) NewProduct(ctx context.Context, ownerId string, input CreateProductDTO) error {
// 	ok, err := s.ownershipService.IsOwner(ctx, ownerId, input.StoreID)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return xerrors.NewWithContext(ctx, ownership.ErrNotOwnerOfResource).WithStatusUnauthorized()
// 	}

// 	newProduct, err := NewProduct(
// 		input.StoreID,
// 		input.Name,
// 		input.Tag,
// 		input.Description,
// 		input.SKU,
// 		input.Price,
// 		input.Type)
// 	if err != nil {
// 		return err
// 	}

// 	return s.r.Save(ctx, newProduct)
// }

// func (s Command) UpdateProduct(ctx context.Context, ownerId string, input UpdateProductDTO) error {
//     ok, err := s.ownershipService.IsOwner(ctx, ownerId, input.StoreID)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return xerrors.NewWithContext(ctx, ownership.ErrNotOwnerOfResource).WithStatusUnauthorized()
// 	}

// 	p, err := s.r.FindByID(ctx, input.ID)
// 	if err != nil {
// 		return err
// 	}

// 	err = p.Update(input.Type, input.Name, input.Description, input.SKU)
// 	if err != nil {
// 		return err
// 	}

// 	return s.r.Save(ctx, p)
// }

// func (s Command) IncreaseProductStock(ctx context.Context, ownerId string, input ChangeStockProductDTO) error {
//     ok, err := s.ownershipService.IsOwner(ctx, ownerId, input.StoreID)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return xerrors.NewWithContext(ctx, ownership.ErrNotOwnerOfResource).WithStatusUnauthorized()
// 	}

// 	p, err := s.r.FindByID(ctx, input.ID)
// 	if err != nil {
// 		return err
// 	}

// 	err = p.IncreaseStock(input.Quantity)
// 	if err != nil {
// 		return err
// 	}

// 	return s.r.Save(ctx, p)
// }

// func (s Command) DecreaseProductStock(ctx context.Context, ownerId string, input ChangeStockProductDTO) error {
//     ok, err := s.ownershipService.IsOwner(ctx, ownerId, input.StoreID)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return xerrors.NewWithContext(ctx, ownership.ErrNotOwnerOfResource).WithStatusUnauthorized()
// 	}

// 	p, err := s.r.FindByID(ctx, input.ID)
// 	if err != nil {
// 		return err
// 	}

// 	err = p.DecreaseStock(input.Quantity)
// 	if err != nil {
// 		return err
// 	}

// 	return s.r.Save(ctx, p)
// }

// func (s Command) ChangeProductImage(ctx context.Context, ownerID string, params UploadProductImageDTO) error {
//     ok, err := s.ownershipService.IsOwner(ctx, ownerId, input.StoreID)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return xerrors.NewWithContext(ctx, ownership.ErrNotOwnerOfResource).WithStatusUnauthorized()
// 	}

// 	p, err := s.r.FindByID(ctx, params.ProductID)
// 	if err != nil {
// 		return err
// 	}

// 	objectKey := p.ID + "-image" + params.Ext

// 	if p.Image != "" {
// 		err = s.s.RemoveFile(ctx, string(ProductBucket), objectKey)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	url, err := s.s.UploadFile(ctx, string(ProductBucket), objectKey, params.Image)
// 	if err != nil {
// 		return err
// 	}

// 	p.Image = url

// 	return s.r.Save(ctx, p)
// }
