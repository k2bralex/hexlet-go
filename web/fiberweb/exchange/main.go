package main

import (
	fmt "fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

var exchangeRate = map[string]float64{
	"USD/EUR": 0.8,
	"EUR/USD": 1.25,
	"USD/GBP": 0.7,
	"GBP/USD": 1.43,
	"USD/JPY": 110,
	"JPY/USD": 0.0091,
}

func main() {
	file, err := os.OpenFile(".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	webApp := fiber.New()
	webApp.Use(logger.New(logger.Config{
		Output: file,
	}))

	webApp.Get("/convert", GetCurrency)

	logrus.Fatal(webApp.Listen(":8080"))
}

func GetCurrency(ctx *fiber.Ctx) error {
	from := ctx.Query("from")
	to := ctx.Query("to")

	if _, ok := exchangeRate[from+"/"+to]; !ok {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.Status(fiber.StatusOK).SendString(
		fmt.Sprintf("%.2f", exchangeRate[from+"/"+to]))
}

