package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var postLikes = map[string]int64{}

func main() {
	webApp := fiber.New(fiber.Config{Immutable: true})
	//webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Go to /likes/12345")
	})

	webApp.Get("/likes/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		likes, ok := postLikes[id]
		if !ok {
			return ctx.SendStatus(http.StatusNotFound)
		}
		ctx.Status(http.StatusOK)
		return ctx.SendString(strconv.FormatInt(likes, 10))
	})

	webApp.Post("/likes/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		if _, ok := postLikes[id]; !ok {
			postLikes[id] = 1
			return ctx.SendStatus(http.StatusCreated)
		}
		postLikes[id]++
		fmt.Println(postLikes)
		return ctx.SendStatus(http.StatusOK)
	})

	logrus.Fatal(webApp.Listen(":8080"))
}
