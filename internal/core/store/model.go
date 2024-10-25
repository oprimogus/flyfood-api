package store

import (
	"errors"
	"time"

	"github.com/oprimogus/cardapiogo/internal/core/address"
)

var ErrClosingTimeBeforeOpeningTime = errors.New("closingTime cannot be before openingTime")
var ErrOpeningTimeAfterClosingTime = errors.New("openingTime cannot be after closingTime")
var ErrNotOwner = errors.New("only owner can do this action")

const (
	defaultStoreScore  = 500
	BusinessHourLayout = "15:04:05"
)

type ShopType string

const (
	StoreShopRestaurant  ShopType = "restaurant"
	StoreShopPharmacy    ShopType = "pharmacy"
	StoreShopTobbaco     ShopType = "tobbaco"
	StoreShopMarket      ShopType = "market"
	StoreShopConvenience ShopType = "convenience"
	StoreShopPub         ShopType = "pub"
)

func IsValidShopType(shopType string) bool {
	switch ShopType(shopType) {
	case StoreShopRestaurant,
		StoreShopPharmacy,
		StoreShopTobbaco,
		StoreShopMarket,
		StoreShopConvenience,
		StoreShopPub:
		return true
	default:
		return false
	}
}

type PaymentMethod string

const (
	Credit PaymentMethod = "credit"
	Debit  PaymentMethod = "debit"
	Pix    PaymentMethod = "pix"
	Cash   PaymentMethod = "cash"
)

func IsValidPaymentMethod(PaymentMethodEnum string) bool {
	switch PaymentMethod(PaymentMethodEnum) {
	case Credit, Debit, Pix, Cash:
		return true
	default:
		return false
	}
}

func IsBusinessHourString(value string) bool {
	_, err := time.Parse(BusinessHourLayout, value)
	return err == nil
}

func IsValidBusinessHour(bh BusinessHours) (bool, error) {
	if bh.ClosingTime.Before(bh.OpeningTime) {
		return false, ErrClosingTimeBeforeOpeningTime
	}

	if bh.OpeningTime.After(bh.ClosingTime) {
		return false, ErrOpeningTimeAfterClosingTime
	}

	return true, nil
}

type BusinessHours struct {
	WeekDay     int       `json:"weekDay" validate:"min=0,max=6"`
	OpeningTime time.Time `json:"openingTime" validate:"required"`
	ClosingTime time.Time `json:"closingTime" validate:"required"`
	TimeZone    string    `json:"timeZone" validate:"required"`
}

func IsValidBusinessHourSlice(slice []BusinessHours) bool {
	businessHours := slice
	if len(businessHours) > 0 {
		mapFrequency := map[int]int{}
		for _, v := range businessHours {
			mapFrequency[v.WeekDay] += 1
		}
		for i := range businessHours {
			if mapFrequency[i] > 1 {
				return false
			}
		}
	}
	return true
}

type StoreFilter struct {
	Range     int      `json:"range"`
	Score     int      `json:"score"`
	Name      string   `json:"name"`
	City      string   `json:"city"`
	Latitude  string   `json:"latitude"`
	Longitude string   `json:"longitude"`
	Type      ShopType `json:"type"`
}

type Store struct {
	ID                 string          `json:"id" validate:"required,uuid"`
	CpfCnpj            string          `json:"cpfCnpj" validate:"required,cpfCnpj"`
	OwnerID            string          `json:"owner_id" validate:"required,uuid"`
	Name               string          `json:"name" validate:"required,lte=25"`
	ProfileImage       string          `json:"profileImage"`
	HeaderImage        string          `json:"headerImage"`
	Active             bool            `json:"active" validate:"required,boolean"`
	Phone              string          `json:"phone" validate:"required,phone"`
	Score              int             `json:"score" validate:"required,number"`
	Address            address.Address `json:"address" validate:"required"`
	Type               ShopType        `json:"type" validate:"required,shopType"`
	BusinessHours      []BusinessHours `json:"businessHours" validate:"dive"`
	PaymentMethodEnums []PaymentMethod `json:"paymentMethod" validate:"dive"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	DeletedAt          time.Time       `json:"deleted_at"`
}

func NewStore(ownerID, name, cpfCnpj, phone string, address address.Address, shopType ShopType) Store {
	store := Store{
		Name:               name,
		CpfCnpj:            cpfCnpj,
		OwnerID:            ownerID,
		Phone:              phone,
		Address:            address,
		Type:               shopType,
		Score:              defaultStoreScore,
		BusinessHours:      []BusinessHours{},
		PaymentMethodEnums: []PaymentMethod{},
	}
	return store
}
