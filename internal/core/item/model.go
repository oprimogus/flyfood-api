package item

import "time"

type ItemType string

const (
	FoodItemType        ItemType = "FOOD"
	WaterGallonItemType ItemType = "WATERGALLON"
)

const defaultScore = 500

func IsValidItemType(itemType string) bool {
	switch ItemType(itemType) {
	case FoodItemType,
		WaterGallonItemType:
		return true
	default:
		return false
	}
}

type Item struct {
	ID             string                 `json:"id" validate:"required"`
	StoreID        string                 `json:"storeID" validate:"required,uuid"`
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
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	DeletedAt      time.Time              `json:"deleted_at"`
}
