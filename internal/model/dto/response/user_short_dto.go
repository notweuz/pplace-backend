package response

type UserShortDto struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func NewUserShortDto(id uint, username string) *UserShortDto {
	return &UserShortDto{
		ID:       id,
		Username: username,
	}
}
