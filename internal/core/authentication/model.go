package authentication

type JWT struct {
	AccessToken      string `json:"accessToken" example:"eysjs9a..."`
	IDToken          string `json:"idToken" example:"eysjg54g3ba..."`
	ExpiresIn        int    `json:"expiresIn" example:"300"`
	RefreshExpiresIn int    `json:"refreshExpiresIn" example:"300"`
	RefreshToken     string `json:"refreshToken" example:"eynmiks.ewij..."`
	TokenType        string `json:"tokenType" example:"Bearer"`
	NotBeforePolicy  int    `json:"notBeforePolicy" example:"150"`
	SessionState     string `json:"sessionState" example:"634g43y3..."`
	Scope            string `json:"scope" example:"email client user ..."`
}
