package request

type PlacePixelDto struct {
	X     uint   `json:"x" validate:"required"`
	Y     uint   `json:"y" validate:"required"`
	Color string `json:"color" validate:"required"`
}
