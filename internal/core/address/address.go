package address

type Address struct {
	Name         string `json:"name" validate:"lte=25" example:"Home"`
	AddressLine1 string `json:"addressLine1" validate:"required,lte=40" example:"123 Main Street"`
	AddressLine2 string `json:"addressLine2" validate:"required,lte=20" example:"Apt 4B"`
	Neighborhood string `json:"neighborhood" validate:"required,lte=25" example:"Downtown"`
	City         string `json:"city" validate:"required,lte=25" example:"New York"`
	State        string `json:"state" validate:"required,lte=15" example:"NY"`
	PostalCode   string `json:"postalCode" validate:"required,lte=15" example:"10001"`
	Latitude     string `json:"latitude" example:"40.7128"`
	Longitude    string `json:"longitude" example:"-74.0060"`
	Country      string `json:"country" validate:"required,lte=15" example:"United States"`
}
