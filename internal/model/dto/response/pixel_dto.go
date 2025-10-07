package response

type PixelDto struct {
	ID     uint         `json:"id"`
	X      uint         `json:"x"`
	Y      uint         `json:"y"`
	Color  string       `json:"color"`
	Author UserShortDto `json:"author"`
}

func NewPixelDto(id, x, y uint, color string, author UserShortDto) *PixelDto {
	return &PixelDto{
		ID:     id,
		X:      x,
		Y:      y,
		Color:  color,
		Author: author,
	}
}
