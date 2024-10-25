package user

func Roles(rolesJSON []string) []Role {
	roles := make([]Role, len(rolesJSON))
	for i, v := range rolesJSON {
		roles[i] = Role(v)
	}
	return roles
}

type CreateProfileParams struct {
	Name     string `json:"name" validate:"required" example:"John"`
	LastName string `json:"lastName" validate:"required" example:"Doe"`
	Phone    string `json:"phone" validate:"required,phone" example:"13997590579"`
}

func (d CreateProfileParams) ToEntity() Profile {
	return Profile{
		Name:     d.Name,
		LastName: d.LastName,
		Phone:    d.Phone,
	}
}

type CreateParams struct {
	Email    string              `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Password string              `json:"password" validate:"required" example:"*********"`
	Profile  CreateProfileParams `json:"profile" validate:"required"`
}

func (d CreateParams) ToEntity() User {
	return User{
		Profile:  d.Profile.ToEntity(),
		Email:    d.Email,
		Password: d.Password,
	}
}

type UpdateProfileParams struct {
	Name     string `json:"name" validate:"required" example:"John"`
	LastName string `json:"lastName" validate:"required" example:"Doe"`
	Phone    string `json:"phone" validate:"required" example:"13997590579"`
}

type AddRolesParams struct {
	Roles []Role `json:"roles" validate:"required,role" enums:"consumer, delivery_man, owner, employee"`
}

func (d UpdateProfileParams) ToEntity() User {
	return User{
		Profile: Profile{
			Name:     d.Name,
			LastName: d.LastName,
			Phone:    d.Phone,
		},
	}
}

type LoginParams struct {
	Email    string `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Password string `json:"password" validate:"required" example:"*********"`
}

type UpdatePasswordParams struct {
	Password    string `json:"password" validate:"required" example:"********"`
	NewPassword string `json:"newPassword" validate:"required" example:"********"`
}
