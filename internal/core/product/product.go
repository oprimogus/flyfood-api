package product

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/oprimogus/cardapiogo/internal/xvalidator"
)

const defaultScore = 500

func init() {
	validator := xvalidator.GetPtInstance()
	err := validator.AddValidations(validationsMap)
	if err != nil {
		panic(err)
	}
}

type Type string

const (
	Food  Type = "FOOD"
	Water Type = "WATER"
)

func IsValidType(value string) bool {
	switch Type(value) {
	case Food, Water:
		return true
	default:
		return false
	}
}

type Product struct {
	ID               string                 `json:"id" validate:"required,uuid"`
	StoreID          string                 `json:"store_id" validate:"required,uuid"`
	SKU              string                 `json:"SKU"`
	ActiveForSale    bool                   `json:"active" validate:"boolean"`
	PromoActive      bool                   `json:"promo_active" validate:"boolean"`
	Type             Type                   `json:"type" validate:"required,productType"`
	Name             string                 `json:"name" validate:"required,lte=25,gte=3"`
	Description      string                 `json:"description" validate:"required,lte=255"`
	StockQuantity    int                    `json:"stock_quantity" validate:"number"`
	Score            int                    `json:"score" validate:"number,required"`
	Image            string                 `json:"image"`
	Details          map[string]interface{} `json:"details"`
	Price            int                    `json:"price" validate:"required,number,gt=0"`
	PromotionalPrice int                    `json:"promotional_price" validate:"number,gte=0"`
}

func NewProduct(storeID, name, description, sku string, price int, productType Type) (Product, error) {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		return Product{}, fmt.Errorf("fail on create  product id: %w", err)
	}
	newProduct := Product{
		ID:            uuidV7.String(),
		StoreID:       storeID,
		SKU:           sku,
		ActiveForSale: false,
		PromoActive:   false,
		Name:          name,
		Description:   description,
		Type:          productType,
		Price:         price,
		Score:         defaultScore,
	}
	if err := newProduct.Validate(); err != nil {
		return Product{}, err
	}
	return newProduct, nil
}

func (p *Product) Validate() error {
	return xvalidator.GetPtInstance().Validate(p)
}

func (p *Product) IncreaseStock(quantity int) error {
	if quantity < 0 {
		return ErrQuantityGreaterThanZero
	}
	sum := p.StockQuantity + quantity
	if sum < p.StockQuantity {
		return ErrInvalidStockQuantity
	}
	p.StockQuantity += quantity
	return p.Validate()
}

func (p *Product) IsAvailable(quantity int) bool {
	return p.StockQuantity >= quantity
}

func (p *Product) DecreaseStock(quantity int) error {
	if quantity < 0 {
		return ErrQuantityGreaterThanZero
	}
	sum := p.StockQuantity - quantity
	if sum < 0 {
		return ErrStockLessThanZero
	}
	p.StockQuantity += quantity
	return p.Validate()
}

func (p *Product) EnableForSale() {
	p.ActiveForSale = true
}

func (p *Product) DisableForSale() {
	p.ActiveForSale = false
}

func (p *Product) ChangePrice(price int) error {
	if price <= 0 {
		return ErrPriceZero
	}
	if price < p.PromotionalPrice {
		return ErrPriceLessThanPromoPrice
	}
	p.Price = price
	return p.Validate()
}

func (p *Product) ChangePromotionalPrice(price int) error {
	if price <= 0 {
		return ErrPriceZero
	}
	if price >= p.Price {
		return ErrPromoPriceGreaterThanPrice
	}
	p.PromotionalPrice = price
	return p.Validate()
}

func (p *Product) EnablePromotionalPrice() {
	p.PromoActive = true
}

func (p *Product) DisablePromotionalPrice() {
	p.PromoActive = false
}

func (p *Product) Update(productType Type, name, description, sku string) error {
	p.Name = name
	p.SKU = sku
	p.Description = description
	p.Type = productType
	return p.Validate()
}
