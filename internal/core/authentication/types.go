package authentication

type SignInParams struct {
	Email    string `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Password string `json:"password" validate:"required" example:"*********"`
}

type RefreshParams struct {
	RefreshToken string `json:"refreshToken" validate:"required" example:"eyajbhhkd....."`
}
