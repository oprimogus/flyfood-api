package product

import (
	"context"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	"log/slog"
)

type Service struct {
	r  Repository
	sr store.Repository
}

func NewService(r Repository, sr store.Repository) Service {
	return Service{r: r, sr: sr}
}

func (s Service) NewProduct(ctx context.Context, ownerId string, input CreateProductDTO) error {
	owner, err := s.sr.FindOwnerByID(ctx, ownerId)
	if err != nil {
		return err
	}

	slog.Info("Verifying owner result", "owner", owner)

	return nil

	//if owner.ID != ownerId {
	//	return store.ErrNotOwner
	//}
	//
	//newProduct, err := NewProduct(
	//	input.StoreID,
	//	input.Name,
	//	input.Tag,
	//	input.Description,
	//	input.SKU,
	//	input.Price,
	//	input.Type)
	//if err != nil {
	//	return err
	//}
	//
	//return s.r.Save(ctx, newProduct)
}

func (s Service) IncreaseStock(ctx context.Context, ownerId string, input ChangeStockProductDTO) error {

	p, err := s.r.FindByID(ctx, input.ID)
	if err != nil {
		return err
	}

	err = p.IncreaseStock(input.Quantity)
	if err != nil {
		return err
	}

	return s.r.Save(ctx, p)

}
