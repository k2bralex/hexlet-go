package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"sort"
)

type (
	BinarySearchRequest struct {
		Numbers []int `json:"numbers"`
		Target  int   `json:"target"`
	}

	BinarySearchResponse struct {
		TargetIndex int    `json:"target_index"`
		Error       string `json:"error,omitempty"`
	}
)

const targetNotFound = -1

func main() {
	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	webApp.Post("/search", func(ctx *fiber.Ctx) error {
		var request BinarySearchRequest

		if err := ctx.BodyParser(&request); err != nil {
			ctx.Status(400)
			return ctx.JSON(BinarySearchResponse{
				TargetIndex: targetNotFound,
				Error:       "Invalid JSON",
			})
		}

		if !contains(request.Numbers, request.Target) {
			ctx.Status(404)
			return ctx.JSON(BinarySearchResponse{
				TargetIndex: targetNotFound,
				Error:       "Target was not found",
			})
		}

		return ctx.JSON(BinarySearchResponse{
			TargetIndex: sort.SearchInts(request.Numbers, request.Target),
		})
	})

	logrus.Fatal(webApp.Listen(":8080"))
}

func contains(s []int, str int) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
