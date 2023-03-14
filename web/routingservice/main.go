package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

type (
	CreateLinkRequest struct {
		External string `json:"external"`
		Internal string `json:"internal"`
	}

	GetLinkResponse struct {
		Internal string `json:"internal"`
	}
)

var links = make(map[string]string)

func main() {
	webApp := fiber.New(fiber.Config{Immutable: true})
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	webApp.Get("/links/:external", GetLink)
	webApp.Post("/links", CreateLink)

	logrus.Fatal(webApp.Listen(":8080"))
}

func GetLink(ctx *fiber.Ctx) error {
	ext := ctx.Params("external")
	decoded, err := url.QueryUnescape(ext)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if _, ok := links[decoded]; !ok {
		ctx.Status(fiber.StatusNotFound)
		return ctx.SendString("link not found")
	}

	ctx.Status(http.StatusOK)
	return ctx.JSON(GetLinkResponse{
		Internal: links[decoded],
	})
}

func CreateLink(ctx *fiber.Ctx) error {
	var request CreateLinkRequest
	if err := ctx.BodyParser(&request); err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.SendString("invalid JSON")
	}
	links[request.External] = request.Internal
	return ctx.SendStatus(fiber.StatusOK)
}
