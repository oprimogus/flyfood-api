package product

type CreateProductDTO struct {
	StoreID     string `json:"storeID" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string `json:"name" validate:"required,lte=25,gte=3" example:"Pizza Portuguesa"`
	Description string `json:"description" validate:"required,lte=255" example:"Pizza com queijo, azeitona, presunto"`
	Tag         string `json:"tag" validate:"required" example:"Promotional 1"`
	SKU         string `json:"SKU" example:"XBOO168"`
	Price       int    `json:"price" validate:"required,number,gt=0" example:"5990"`
	Type        Type   `json:"type" validate:"required,productType" example:"FOOD"`
}

type UpdateProductDTO struct {
	ID          string `json:"id" validate:"required,uuid" example:"0194900a-8909-755a-bd61-ec7a18224200"`
	StoreID     string `json:"storeID" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string `json:"name" validate:"required,lte=25,gte=3" example:"Pizza Portuguesa"`
	Description string `json:"description" validate:"required,lte=255" example:"Pizza com queijo, azeitona, presunto"`
	Tag         string `json:"tag" validate:"required" example:"Promotional 1"`
	SKU         string `json:"SKU" example:"XBOO168"`
	Price       int    `json:"price" validate:"required,number,gt=0" example:"5990"`
	Type        Type   `json:"type" validate:"required,productType" example:"FOOD"`
}

type UploadProductImageDTO struct {
	StoreID   string `json:"storeID" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProductID string `json:"productID" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Image     []byte `json:"image"`
	Ext       string `json:"ext"`
}

type ChangeStockProductDTO struct {
	ID       string `json:"id" validate:"required,uuid" example:"0194900a-8909-755a-bd61-ec7a18224200"`
	StoreID  string `json:"storeID" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Quantity int    `json:"quantity" validate:"number"`
}

type RemoveProductDTO struct {
	ID      string `json:"id" validate:"required,uuid" example:"0194900a-8909-755a-bd61-ec7a18224200"`
	StoreID string `json:"storeID" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type ProductDTO struct {
	ID               string                 `json:"id" validate:"required,uuid"`
	SKU              string                 `json:"SKU" example:"XBOO168"`
	PromoActive      bool                   `json:"promoActive" validate:"boolean"`
	Type             Type                   `json:"type" validate:"required,productType" example:"FOOD"`
	Tag              string                 `json:"tag" validate:"required" example:"Promotional 1"`
	Name             string                 `json:"name" validate:"required,lte=25,gte=3" example:"Pizza Portuguesa"`
	Description      string                 `json:"description" validate:"required,lte=255" example:"Pizza com queijo, azeitona, presunto"`
	Score            int                    `json:"score" validate:"number,required"`
	Image            string                 `json:"image"`
	Details          map[string]interface{} `json:"details"`
	Price            int                    `json:"price" validate:"required,number,gt=0" example:"5990"`
	PromotionalPrice int                    `json:"promotionalPrice" validate:"number,gte=0"`
}
