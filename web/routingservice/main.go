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
		return c.SendStatus(200)
	})

	webApp.Get("/links/:external", func(ctx *fiber.Ctx) error {
		ext := ctx.Params("external")
		decoded, err := url.QueryUnescape(ext)
		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if _, ok := links[decoded]; !ok {
			ctx.Status(http.StatusNotFound)
			return ctx.SendString("Link not found")
		}

		ctx.Status(http.StatusOK)
		return ctx.JSON(GetLinkResponse{
			Internal: links[decoded],
		})
	})

	webApp.Post("/links", func(ctx *fiber.Ctx) error {
		var request CreateLinkRequest
		if err := ctx.BodyParser(&request); err != nil {
			ctx.Status(http.StatusBadRequest)
			return ctx.SendString("Invalid JSON")
		}
		links[request.External] = request.Internal
		return ctx.SendStatus(http.StatusOK)
	})

	logrus.Fatal(webApp.Listen(":8080"))
}
