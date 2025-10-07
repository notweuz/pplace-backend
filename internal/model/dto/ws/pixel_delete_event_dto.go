package ws

import "time"

type PixelDeleteEventDto struct {
	Action string    `json:"action"`
	ID     uint      `json:"id"`
	X      uint      `json:"x"`
	Y      uint      `json:"y"`
	Time   time.Time `json:"time"`
}

func NewPixelDeleteEventDto(action string, id uint, x, y uint, t time.Time) *PixelDeleteEventDto {
	return &PixelDeleteEventDto{
		Action: action,
		ID:     id,
		X:      x,
		Y:      y,
		Time:   t,
	}
}
