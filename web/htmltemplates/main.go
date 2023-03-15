package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/sirupsen/logrus"
)

type (
	CreateItemRequest struct {
		Name  string `json:"name"`
		Price uint   `json:"price"`
	}

	Item struct {
		Name  string `json:"name"`
		Price uint   `json:"price"`
	}
)

var (
	items []Item
)

func main() {
	viewsEngine := html.New("./templates", ".html")
	webApp := fiber.New(fiber.Config{
		Views: viewsEngine,
	})
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	webApp.Post("/items", Create)
	webApp.Post("/items/view", ReadAll)

	logrus.Fatal(webApp.Listen(":8080"))
}

func Create(ctx *fiber.Ctx) error {
	var req CreateItemRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid JSON")
	}

	items = append(items, Item{
		Name:  req.Name,
		Price: req.Price,
	})

	return ctx.SendStatus(fiber.StatusOK)
}

func ReadAll(ctx *fiber.Ctx) error {
	return ctx.Render("prduct-list", items)
}
