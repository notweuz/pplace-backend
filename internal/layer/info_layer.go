package layer

import (
	"pplace_backend/internal/config"
	"pplace_backend/internal/controller"
	"pplace_backend/internal/service"
)

type InfoLayer struct {
	Service    *service.InfoService
	Controller *controller.InfoController
}

func NewInfoLayer(config *config.PPlaceConfig) InfoLayer {
	s := service.NewInfoService(config)
	c := controller.NewInfoController(&s)

	return InfoLayer{
		Service:    &s,
		Controller: &c,
	}
}
