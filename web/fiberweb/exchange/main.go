package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
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
	webApp := fiber.New()
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	file, err := os.OpenFile("log.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()
	logger.SetOutput(file)

	webApp.Get("/convert", func(ctx *fiber.Ctx) error {
		from := ctx.Query("from")
		to := ctx.Query("to")

		if _, ok := exchangeRate[from+"/"+to]; !ok {
			ctx.Status(404)
			return ctx.SendString("404 Not Found")
		}

		ctx.Status(200)
		return ctx.SendString(fmt.Sprintf("%.2f", exchangeRate[from+"/"+to]))
	})

	logrus.Fatal(webApp.Listen(":8000"))
}
