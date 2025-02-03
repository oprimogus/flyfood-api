package product

import "fmt"

var (
	ErrInvalidStockQuantity       = fmt.Errorf("A soma do estoque atual com a entrada de itens no estoque não pode ser menor que o estoque atual")
	ErrStockLessThanZero          = fmt.Errorf("o estoque não pode ser menor que 0")
	ErrQuantityGreaterThanZero    = fmt.Errorf("a quantidade adicionada ao estoque deve ser maior que 0")
	ErrPriceLessThanPromoPrice    = fmt.Errorf("o preço normal precisa ser maior que o preço promocional")
	ErrPromoPriceGreaterThanPrice = fmt.Errorf("o preço promocional precisa ser menor que o preço normal")
	ErrPriceZero                  = fmt.Errorf("o preço não pode ser menor ou igual a zero")
)
