package request

type UpdateUserDto struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=20"`
	Password string `json:"password,omitempty" validate:"omitempty,min=6"`
}
