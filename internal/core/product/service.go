package product

import "context"

type Service struct {
	r Repository
}

func NewService(r Repository) Service {
	return Service{r}
}

func (s Service) NewProduct(ctx context.Context, input CreateProductDTO) error {
	newProduct, err := NewProduct(
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

	return s.r.Save(ctx, newProduct)
}

//func (s Service) IncreaseStock(ctx context.Context, input CreateProductDTO) error {
//
//}
