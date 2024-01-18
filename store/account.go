package store

import (
	"jiachen/model"
)

type account struct {
	*collection[model.Account]
}

var Account = &account{
	&collection[model.Account]{
		name: "account",
	},
}
