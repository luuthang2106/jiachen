package api

import (
	"jiachen/cache"
	"jiachen/exec"
	"jiachen/model"
	"jiachen/store"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CreateAPI(c *fiber.Ctx) error {
	input := new(model.API)
	if err := c.BodyParser(input); err != nil {
		return err
	}
	path := strings.TrimSpace(input.Path)
	if path == "" {
		return c.Status(fiber.StatusBadRequest).SendString("path required")
	}

	if err := exec.ValidatePath(input.Path); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("path already exists")
	}

	ID, err := store.API.InsertOne(model.API{
		Path:     path,
		IsActive: input.IsActive,
	})
	if err != nil {
		return err
	}

	if input.IsActive != nil && *input.IsActive {
		cache.API.Set(input.Path, struct{}{})
	}

	return c.JSON(fiber.Map{
		"_id": ID.Hex(),
	})
}
