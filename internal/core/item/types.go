package item

type GetItemFilterInput struct {
	Type     Type   `json:"type" validate:"itemType" enums:"FOOD, WATERGALLON"`
	Name     string `json:"name" validate:"lte=25" example:"Burguer"`
	Score    int    `json:"score" validate:"number" example:"419"`
	MaxPrice int    `json:"maxPrice" validate:"number" example:"4990"`
	City     string `json:"city" validate:"required" example:"Guarujá"`
}

type GetItemFilterOutput struct {
	ID             int    `json:"id" validate:"required,number" example:"246643"`
	StoreID        string `json:"storeID" validate:"required,uuid" example:"65293gfk-fgv3fgvf67-f38f378"`
	Type           Type   `json:"type" validate:"required,itemType" enums:"FOOD, WATERGALLON"`
	Name           string `json:"name" validate:"required,lte=25" example:"Burguer"`
	Score          int    `json:"score" validate:"required,number" example:"419"`
	DiscountActive bool   `json:"discountActive" validate:"required,boolean" example:"true"`
	Price          int    `json:"price" validate:"required,gte=0" example:"3990"`
	DiscountPrice  int    `json:"discountPrice" validate:"required,gte=0" example:"2990"`
}

type GetItemByIDOutput struct {
	StoreID        string                 `json:"storeID" validate:"required,uuid" example:"3279gbf23gb-fb23-6f239"`
	Type           Type                   `json:"type" validate:"required,itemType" enums:"FOOD, WATERGALLON"`
	Name           string                 `json:"name" validate:"required,lte=25" example:"Hambúrguer Artesanal"`
	Description    string                 `json:"description" validate:"required,lte=50" example:"Com blend de 150g e molho especial."`
	Score          int                    `json:"score" validate:"required,number" example:"475"`
	Image          string                 `json:"image" validate:"required" example:"https://cardapiogo.com.br/fgbjkgr7erg793bj3lk"`
	Active         bool                   `json:"active" validate:"required,boolean" example:"true"`
	DiscountActive bool                   `json:"discountActive" validate:"required,boolean" example:"true"`
	Details        map[string]interface{} `json:"details"`
	Price          int                    `json:"price" validate:"required,gte=0" example:"2990"`
	DiscountPrice  int                    `json:"discountPrice" validate:"required,gte=0" example:"2499"`
}

type CreateItemInput struct {
	StoreID     string                 `json:"storeID" validate:"required,uuid" example:"3279gbf23gb-fb23-6f239"`
	Type        Type                   `json:"type" validate:"required,itemType" enums:"FOOD, WATERGALLON"`
	Name        string                 `json:"name" validate:"required,lte=25" example:"Hamburguer"`
	Description string                 `json:"description" validate:"required,lte=50" example:"The best burguer!"`
	Details     map[string]interface{} `json:"details"`
	Price       int                    `json:"price" validate:"required,gte=0" example:"2990"`
}

type UpdateItemInput struct {
	ID             int                    `json:"id" validate:"required,number" example:"246643"`
	StoreID        string                 `json:"storeID" validate:"required,uuid"`
	Type           Type                   `json:"type" validate:"required,itemType" enums:"FOOD, WATERGALLON"`
	Active         bool                   `json:"active" validate:"required,boolean"`
	DiscountActive bool                   `json:"discountActive" validate:"required,boolean"`
	Name           string                 `json:"name" validate:"required,lte=25"`
	Description    string                 `json:"description" validate:"required,lte=50"`
	Details        map[string]interface{} `json:"details"`
	Price          int                    `json:"price" validate:"required,gte=0"`
	DiscountPrice  int                    `json:"discountPrice" validate:"required,gte=0"`
}

type CreatedItem struct {
	ID int `json:"id"`
}

type DeleteInput struct {
	ID int `json:"id" validate:"required,number" example:"246643"`
}
