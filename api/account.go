package api

import (
	"jiachen/exec"
	"jiachen/model"
	"jiachen/store"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateAccount(c *fiber.Ctx) error {
	input := new(model.Account)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	username := strings.TrimSpace(input.Username)
	if username == "" {
		return c.Status(fiber.StatusBadRequest).SendString("username required")
	}
	password := strings.TrimSpace(input.Password)
	if password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("password required")
	}
	if err := exec.ValidateUsername(input.Username); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("username already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	ID, err := store.Account.InsertOne(model.Account{
		Username: username,
		Password: string(hashedPassword),
	})
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"_id": ID.Hex(),
	})
}
