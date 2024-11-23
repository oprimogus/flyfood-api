package product

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProduct(t *testing.T) {
	tests := []struct {
		name        string
		storeID     string
		nameProduct string
		description string
		sku         string
		price       int
		pType       Type
	}{
		{
			name:        "should create a new product",
			storeID:     "01933d0f-6d9e-75f1-8448-b35edb94f1e2",
			nameProduct: "product one",
			description: "desc product 1",
			sku:         "P001",
			price:       2500,
			pType:       Food,
		},
	}

	for _, tt := range tests {
		_, err := NewProduct(tt.storeID, tt.nameProduct, tt.description, tt.sku, tt.price, tt.pType)
		assert.NoError(t, err, tt.name)
		if err != nil {
			errJson, _ := json.Marshal(err)
			fmt.Println(string(errJson))
		}
	}
}

func TestProduct_IncreaseStock(t *testing.T) {
	tests := []struct {
		name        string
		newQuantity int
		expected    error
	}{
		{
			name:        "should increase stock",
			newQuantity: 100,
			expected:    nil,
		},
		{
			name:        "should return error with invalid quantity",
			newQuantity: -100,
			expected:    ErrQuantityGreaterThanZero,
		},
	}

	for _, test := range tests {
		p, err := NewProduct("01933d0f-6d9e-75f1-8448-b35edb94f1e2", "product one", "your desc", "P001", 2500, Food)
		assert.NoError(t, err, test.name)

		err = p.IncreaseStock(test.newQuantity)
		assert.Equal(t, test.expected, err, test.name)

	}
}

func TestProduct_DecreaseStock(t *testing.T) {
	tests := []struct {
		name     string
		remove   int
		expected error
	}{
		{
			name:     "should decrease stock",
			remove:   100,
			expected: nil,
		},
		{
			name:     "should return error with invalid quantity",
			remove:   -100,
			expected: ErrQuantityGreaterThanZero,
		},
	}

	for _, test := range tests {
		p, err := NewProduct("01933d0f-6d9e-75f1-8448-b35edb94f1e2", "product one", "your desc", "P001", 2500, Food)
		assert.NoError(t, err, test.name)

		err = p.IncreaseStock(100)
		assert.NoError(t, err, test.name)

		err = p.DecreaseStock(test.remove)
		assert.Equal(t, test.expected, err, test.name)

	}
}

func TestProduct_ChangePrice(t *testing.T) {
	tests := []struct {
		name          string
		value         int
		hasPromoPrice bool
		promoPrice    int
		expect        error
	}{
		{
			name:          "should change price",
			value:         5000,
			hasPromoPrice: false,
			promoPrice:    2000,
			expect:        nil,
		},
		{
			name:          "should change price when has promotional price",
			value:         5000,
			hasPromoPrice: true,
			promoPrice:    2000,
			expect:        nil,
		},
		{
			name:          "should return err when try define price < promotional price",
			value:         1500,
			hasPromoPrice: true,
			promoPrice:    2000,
			expect:        ErrPriceLessThanPromoPrice,
		},
		{
			name:          "should return err when try define price <= 0",
			value:         0,
			hasPromoPrice: true,
			promoPrice:    2000,
			expect:        ErrPriceZero,
		},
	}

	for _, test := range tests {
		p, err := NewProduct("01933d0f-6d9e-75f1-8448-b35edb94f1e2", "product one", "your desc", "P001", 2500, Food)
		assert.NoError(t, err, test.name)

		if test.hasPromoPrice {
			err = p.ChangePromotionalPrice(test.promoPrice)
			assert.NoError(t, err, test.name)
		}

		err = p.ChangePrice(test.value)
		assert.Equal(t, test.expect, err, test.name)
	}
}

func TestProduct_ChangePromotionalPrice(t *testing.T) {
	tests := []struct {
		name       string
		value      int
		promoPrice int
		expect     error
	}{
		{
			name:       "should change promotional price",
			value:      2400,
			promoPrice: 2000,
			expect:     nil,
		},
		{
			name:       "should return err when try define promotional price > normal price",
			value:      6000,
			promoPrice: 2000,
			expect:     ErrPromoPriceGreaterThanPrice,
		},
		{
			name:       "should return err when try define promotional price <= 0",
			value:      0,
			promoPrice: 2000,
			expect:     ErrPriceZero,
		},
	}

	for _, test := range tests {
		p, err := NewProduct("01933d0f-6d9e-75f1-8448-b35edb94f1e2", "product one", "your desc", "P001", 2500, Food)
		assert.NoError(t, err, test.name)

		err = p.ChangePromotionalPrice(test.promoPrice)
		assert.NoError(t, err, test.name)

		err = p.ChangePromotionalPrice(test.value)
		assert.Equal(t, test.expect, err, test.name)
	}

}
