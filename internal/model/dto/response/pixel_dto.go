package response

type PixelDto struct {
	ID     uint    `json:"id"`
	X      uint    `json:"x"`
	Y      uint    `json:"y"`
	Color  string  `json:"color"`
	Author UserDto `json:"author"`
}
