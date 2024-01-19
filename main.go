package main

import (
	"fmt"
	"jiachen/api"
	"jiachen/cache"
	"jiachen/model"
	"jiachen/pool"
	"jiachen/store"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var (
	port            = 80
	routinePoolSize = 100
	uri             = "mongodb+srv://luuthang2106:Th%40ng210697@jiachen.sbxpssx.mongodb.net/"
)

func main() {
	s := store.NewStore(uri)
	store.Account.Init(s)
	store.API.Init(s)

	cache.API.Warmup(cache.WarmupAPICache)

	// start := time.Now()
	// arr, err1 := store.Account.Find(model.Account{Username: "linh"})
	// a, err2 := store.Account.FindOne(model.Account{Username: "linh"})
	// store.Account.UpdateOne(model.Account{Username: "linh"}, model.Account{Username: "thang"})
	// b, err3 := store.Account.FindOne(model.Account{Username: "thang"})
	// store.Account.UpdateOne(model.Account{Username: "thang"}, model.Account{Username: "linh"})
	// fmt.Println(arr, err1)
	// fmt.Println(a, err2)
	// fmt.Println(b, err3)
	// fmt.Println(time.Since(start).Milliseconds())
	rootHashedPassword, _ := bcrypt.GenerateFromPassword([]byte("root"), bcrypt.MinCost)
	store.Account.Upsert(model.Account{Username: "root"}, model.Account{Username: "root", Password: string(rootHashedPassword)})

	pool.NewPool(routinePoolSize)
	app := fiber.New()

	// apis := app.Group("/api")
	// apis.Get("/aaaa/:dddd", func(c *fiber.Ctx) error {
	// 	// fmt.Println(c.Route().Path)
	// 	return c.SendString("Hello, World!")
	// })

	// apis
	app.Post("/login", api.Login)
	app.Post("/account", api.ProtectPath, api.CreateAccount)
	app.Put("/account", api.ProtectPath, api.CreateAccount)
	app.Get("/thang", api.ProtectPath, api.CreateAccount)

	app.Listen(fmt.Sprintf(":%d", port))
}
