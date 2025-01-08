package product

type CreateProductDTO struct {
	StoreID     string `json:"store_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string `json:"name" validate:"required,lte=25,gte=3" example:"Pizza Portuguesa"`
	Description string `json:"description" validate:"required,lte=255" example:"Pizza com queijo, azeitona, presunto"`
	Tag         string `json:"tag" validate:"required" example:"Promotional 1"`
	SKU         string `json:"SKU" example:"XBOO168"`
	Price       int    `json:"price" validate:"required,number,gt=0" example:"5990"`
	Type        Type   `json:"type" validate:"required,productType" example:"FOOD"`
}
