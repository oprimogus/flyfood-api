package address

type Address struct {
	Name         string `json:"name" validate:"lte=25" example:"Home"`
	AddressLine1 string `db:"address_line_1" json:"addressLine1" validate:"required,lte=40" example:"123 Main Street"`
	AddressLine2 string `db:"address_line_2" json:"addressLine2" validate:"required,lte=20" example:"Apt 4B"`
	Neighborhood string `db:"neighborhood" json:"neighborhood" validate:"required,lte=25" example:"Downtown"`
	City         string `db:"city" json:"city" validate:"required,lte=25" example:"New York"`
	State        string `db:"state" json:"state" validate:"required,lte=15" example:"NY"`
	PostalCode   string `db:"postal_code" json:"postalCode" validate:"required,lte=15" example:"10001"`
	Latitude     string `db:"latitude" json:"latitude" example:"40.7128"`
	Longitude    string `db:"longitude" json:"longitude" example:"-74.0060"`
	Country      string `db:"country" json:"country" validate:"required,lte=15" example:"United States"`
}
