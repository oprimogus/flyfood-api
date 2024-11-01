package item

import "time"

type Type string

const (
	Food        Type = "FOOD"
	WaterGallon Type = "WATERGALLON"
)

const defaultScore = 500

func IsValidItemType(itemType string) bool {
	switch Type(itemType) {
	case Food,
		WaterGallon:
		return true
	default:
		return false
	}
}

type Item struct {
	ID             int                    `json:"id" validate:"required,number" example:"38675"`
	StoreID        string                 `json:"storeID" validate:"required,uuid"`
	Type           Type                   `json:"type" validate:"required,itemType"`
	Name           string                 `json:"name" validate:"required,lte=25"`
	Description    string                 `json:"description" validate:"required,lte=50"`
	Score          int                    `json:"score" validate:"required,number"`
	Image          string                 `json:"image" validate:"required"`
	Active         bool                   `json:"active" validate:"required,boolean"`
	DiscountActive bool                   `json:"discountActive" validate:"required,boolean"`
	Details        map[string]interface{} `json:"details"`
	Price          int                    `json:"price" validate:"required,gte=0"`
	DiscountPrice  int                    `json:"discountPrice" validate:"required,gte=0"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
	DeletedAt      *time.Time             `json:"deletedAt"`
}
