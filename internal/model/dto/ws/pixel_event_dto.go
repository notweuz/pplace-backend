package ws

import "time"

type PixelEventDto struct {
	Action   string    `json:"action"`
	ID       uint      `json:"id"`
	X        uint      `json:"x"`
	Y        uint      `json:"y"`
	Color    string    `json:"color"`
	UserID   uint      `json:"userId"`
	Username string    `json:"username"`
	Time     time.Time `json:"time"`
}

func NewPixelEventDto(action, color, username string, id, x, y, userID uint, time time.Time) *PixelEventDto {
	return &PixelEventDto{
		Action:   action,
		ID:       id,
		X:        x,
		Y:        y,
		Color:    color,
		UserID:   userID,
		Username: username,
		Time:     time,
	}
}
