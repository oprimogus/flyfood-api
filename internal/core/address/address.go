package address

import "github.com/google/uuid"

type Address struct {
	ID           uuid.UUID `json:"id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name         string    `json:"name" validate:"lte=25" example:"Home"`
	Default      bool      `json:"default" example:"true"`
	AddressLine1 string    `json:"addressLine1" validate:"required,lte=40" example:"123 Main Street"`
	AddressLine2 string    `json:"addressLine2" validate:"required,lte=20" example:"Apt 4B"`
	Neighborhood string    `json:"neighborhood" validate:"required,lte=25" example:"Downtown"`
	City         string    `json:"city" validate:"required,lte=25" example:"New York"`
	State        string    `json:"state" validate:"required,lte=15" example:"NY"`
	PostalCode   string    `json:"postalCode" validate:"required,lte=15" example:"10001"`
	Latitude     float64   `json:"latitude" example:"40.7128"`
	Longitude    float64   `json:"longitude" example:"-74.0060"`
	Country      string    `json:"country" validate:"required,lte=15" example:"United States"`
}
