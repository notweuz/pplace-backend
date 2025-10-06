package response

type AuthTokenDto struct {
	Token string `json:"token"`
}

func NewAuthTokenDto(token string) *AuthTokenDto {
	return &AuthTokenDto{Token: token}
}
