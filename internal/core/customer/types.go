package customer

type CreateProfileDTO struct {
	ID       int    `json:"id" validate:"required,number" example:"295105940221919239"`
	Name     string `json:"name" validate:"required,alpha,gte=3,lte=25" example:"John"`
	LastName string `json:"last_name" validate:"required,gte=3,lte=60" example:"Doe"`
	CPF      string `json:"cpf" validate:"cpf" example:"52024227090"`
	Email    string `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Phone    string `json:"phone" validate:"required,phone" example:"+5513997590579"`
}

type UpdateProfileDTO struct {
	Name     string `json:"name" validate:"required,alpha,gte=3,lte=25" example:"John"`
	LastName string `json:"last_name" validate:"required,gte=3,lte=60" example:"Doe"`
	Phone    string `json:"phone" validate:"required,phone" example:"+5513997590579"`
}
