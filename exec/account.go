package exec

import (
	"fmt"
	"jiachen/model"
	"jiachen/store"
)

func ValidateUsername(username string) error {
	acc, err := store.Account.FindOne(model.Account{Username: username})
	if err == nil || acc != nil {
		return fmt.Errorf("username already exists")
	}
	return nil
}
