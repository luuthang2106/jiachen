package api

import (
	"jiachen/model"
	"jiachen/store"
	"jiachen/util"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	tokenSecretKey = "L2bZwMsKffkTBA9z+1Ajh0ggDII/xf71ab0VLiLsCBg="
)

func ProtectPath(c *fiber.Ctx) error {
	path := c.Route().Path

	tokenStr := c.Cookies("token")
	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("token invalid")
	}
	// Giải mã token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecretKey), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}
	// Lấy claims từ token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("token invalid")
	}

	username := claims["username"].(string)

	if username == "root" {
		return c.Next()
	}

	acc, err := store.Account.FindOne(model.Account{Username: username})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if acc.APIs == nil || len(*acc.APIs) == 0 {
		return c.Status(fiber.StatusForbidden).SendString("account don't have any permissions")
	}

	if util.Contains(*acc.APIs, path) {
		return c.Next()
	}

	return c.Status(fiber.StatusForbidden).SendString("account don't have any permissions")
}

func Login(c *fiber.Ctx) error {
	input := new(model.Account)
	if err := c.BodyParser(&input); err != nil {
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
	acc, err := store.Account.FindOne(model.Account{Username: input.Username})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(input.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("password wrong")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = acc.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Thời gian hết hạn của token: 1 ngày

	// Ký token với secret key và tạo chuỗi token
	tokenString, err := token.SignedString([]byte(tokenSecretKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24), // Thời gian hết hạn của cookie: 1 ngày
		HTTPOnly: true,
	})

	return c.SendString("login successfully")
}
