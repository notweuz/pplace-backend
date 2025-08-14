package request

type PixelPlaceDto struct {
	X     uint   `json:"x" validate:"required"`
	Y     uint   `json:"y" validate:"required"`
	Color string `json:"color" validate:"required"`
}
