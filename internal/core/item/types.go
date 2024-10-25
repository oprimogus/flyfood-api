package item

type GetItemFilter struct {
	ID             string                 `json:"id" validate:"required"`
	StoreID        string                 `json:"storeID" validate:"required,uuid"`
	Type           ItemType               `json:"type" validate:"required,itemType"`
	Name           string                 `json:"name" validate:"required,lte=25"`
	Score          int                    `json:"score" validate:"required,number"`
	DiscountActive bool                   `json:"discountActive" validate:"required,boolean"`
	Details        map[string]interface{} `json:"details"`
	Price          int                    `json:"price" validate:"required,gte=0"`
	DiscountPrice  int                    `json:"discountPrice" validate:"required,gte=0"`
}

type GetItemByIDOutput struct {
	Type           ItemType               `json:"type" validate:"required,itemType"`
	Name           string                 `json:"name" validate:"required,lte=25"`
	Description    string                 `json:"description" validate:"required,lte=50"`
	Score          int                    `json:"score" validate:"required,number"`
	Image          string                 `json:"image" validate:"required"`
	Active         bool                   `json:"active" validate:"required,boolean"`
	DiscountActive bool                   `json:"discountActive" validate:"required,boolean"`
	Details        map[string]interface{} `json:"details"`
	Price          int                    `json:"price" validate:"required,gte=0"`
	DiscountPrice  int                    `json:"discountPrice" validate:"required,gte=0"`
}

type CreateItemInput struct {
	StoreID     string                 `json:"storeID" validate:"required,uuid"`
	Type        ItemType               `json:"type" validate:"required,itemType"`
	Name        string                 `json:"name" validate:"required,lte=25"`
	Description string                 `json:"description" validate:"required,lte=50"`
	Details     map[string]interface{} `json:"details"`
	Price       int                    `json:"price" validate:"required,gte=0"`
}

type UpdateItemInput struct {
	StoreID        string                 `json:"storeID" validate:"required,uuid"`
	Type           ItemType               `json:"type" validate:"required,itemType"`
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
