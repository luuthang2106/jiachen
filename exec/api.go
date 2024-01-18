package exec

import (
	"fmt"
	"jiachen/model"
	"jiachen/store"
)

func ValidatePath(path string) error {
	api, err := store.API.FindOne(model.API{Path: path})
	if err == nil || api != nil {
		return fmt.Errorf("API already exists")
	}
	return nil
}
