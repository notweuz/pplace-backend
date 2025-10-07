package response

import "time"

type UserDto struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	LastPlaced time.Time `json:"last_placed"`
}

func NewUserDto(id uint, username string, lastPlaced time.Time) *UserDto {
	return &UserDto{
		ID:         id,
		Username:   username,
		LastPlaced: lastPlaced,
	}
}
