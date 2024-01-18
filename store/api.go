package store

import (
	"jiachen/model"
)

type api struct {
	*collection[model.API]
}

var API = &api{
	&collection[model.API]{
		name: "api",
	},
}
