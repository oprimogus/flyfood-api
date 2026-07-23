package product

import (
	"errors"
)

var (
	ErrInvalidStockQuantity       = errors.New("a soma do estoque atual com a entrada de itens no estoque não pode ser menor que o estoque atual")
	ErrStockLessThanZero          = errors.New("o estoque não pode ser menor que 0")
	ErrQuantityGreaterThanZero    = errors.New("a quantidade adicionada ao estoque deve ser maior que 0")
	ErrPriceLessThanPromoPrice    = errors.New("o preço normal precisa ser maior que o preço promocional")
	ErrPromoPriceGreaterThanPrice = errors.New("o preço promocional precisa ser menor que o preço normal")
	ErrPriceZero                  = errors.New("o preço não pode ser menor ou igual a zero")
	ErrPriceZeroWhenEnableProduct = errors.New("para habilitar o produto, seu preço deve ser definido como > 0")
)
