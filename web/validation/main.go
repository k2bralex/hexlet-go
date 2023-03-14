package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type User struct {
	ID      int64
	Email   string
	Age     int
	Country string
}

var users = map[int64]User{}

type (
	CreateUserRequest struct {
		ID      int64  `json:"id" validate:"required,number,min=1"`
		Email   string `json:"email" validate:"required,email,min=6,max=32"`
		Age     int    `json:"age" validate:"required,min=18,max=130"`
		Country string `json:"country" validate:"oneof=USA Germany France"`
	}
)

func main() {
	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	validate := validator.New()

	webApp.Post("/users", func(ctx *fiber.Ctx) error {
		var req CreateUserRequest
		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString("invalid JSON")
		}

		if err := validate.Struct(req); err != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}

		if _, ok := users[req.ID]; ok {
			return ctx.Status(fiber.StatusConflict).SendString("already exist")
		}

		users[req.ID] = User{
			ID:      req.ID,
			Email:   req.Email,
			Age:     req.Age,
			Country: req.Country,
		}

		return ctx.SendStatus(fiber.StatusOK)
	})

	logrus.Fatal(webApp.Listen(":8080"))
}
