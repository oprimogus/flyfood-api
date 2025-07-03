package product

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/oprimogus/flyfood-api/internal/xvalidator"
)

const defaultScore = 500

func init() {
	err := xvalidator.AddValidations(validationsMap)
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
	ID               string         `json:"id" validate:"required,uuid"`
	StoreID          string         `json:"storeID" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	SKU              string         `json:"SKU" example:"XBOO168"`
	ActiveForSale    bool           `json:"active" validate:"boolean"`
	PromoActive      bool           `json:"promoActive" validate:"boolean"`
	Type             Type           `json:"type" validate:"required,productType" example:"FOOD"`
	Tag              string         `json:"tag" validate:"required" example:"Promotional 1"`
	Name             string         `json:"name" validate:"required,lte=25,gte=3" example:"Pizza Portuguesa"`
	Description      string         `json:"description" validate:"required,lte=255" example:"Pizza com queijo, azeitona, presunto"`
	StockQuantity    int            `json:"stockQuantity" validate:"number"`
	Score            int            `json:"score" validate:"number,required"`
	Image            string         `json:"image"`
	Details          map[string]any `json:"details"`
	Price            int            `json:"price" validate:"required,number,gt=0" example:"5990"`
	PromotionalPrice int            `json:"promotionalPrice" validate:"number,gte=0"`
}

func NewProduct(storeID, name, tag, description, sku string, price int, productType Type) (Product, error) {
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
		Tag:           tag,
		Price:         price,
		Score:         defaultScore,
	}
	if err := newProduct.Validate(); err != nil {
		return Product{}, err
	}
	return newProduct, nil
}

func (p *Product) Validate() error {
	return xvalidator.Validate(p)
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

func (p *Product) EnableForSale() error {
	if p.Price == 0 {
		return ErrPriceZeroWhenEnableProduct
	}
	if p.PromoActive && p.PromotionalPrice == 0 {
		return ErrPriceZeroWhenEnableProduct
	}
	p.ActiveForSale = true
	return nil
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

func (p *Product) ChangeTag(newTag string) error {
	p.Tag = newTag
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
	if p.PromotionalPrice == 0 {
		return
	}
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
