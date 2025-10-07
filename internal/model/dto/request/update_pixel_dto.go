package request

type UpdatePixelDto struct {
	Color string `json:"color" validate:"required"`
}
