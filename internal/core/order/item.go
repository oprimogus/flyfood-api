package order

type Item struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
	Amount    int    `json:"price" validate:"required,gt=0"`
	Details   string `json:"details"`
}
