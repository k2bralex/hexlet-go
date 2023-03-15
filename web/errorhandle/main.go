/*package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
	"time"
)

type (
	// Структура HTTP-запроса на расчет диапазона дат
	DateRangeRequest struct {
		From Date `json:"from"`
		To   Date `json:"to"`
	}

	// Структура даты, которая хранит формат и значение
	Date struct {
		Value  string `json:"value"`
		Format string `json:"format"`
	}

	// Структура HTTP-ответа на расчет диапазона дат
	// Хранит значение в секундах
	DateRangeResponse struct {
		SecondsRange int64 `json:"seconds_range"`
	}
)

func main() {
	webApp := fiber.New()
	// Устанавливаем посредника, который будет
	// восстанавливать веб-приложение после паники
	webApp.Use(recover.New())

	// Настраиваем обработчик для веб-страницы со списком фильмов
	webApp.Post("/daterange", func(c *fiber.Ctx) error {
		logrus.Println("post")
		var req *DateRangeRequest
		if err := c.BodyParser(&req); err != nil {
			logrus.WithError(err).Info("body parser")
			return c.Status(fiber.StatusBadRequest).SendString("bad JSON")
		}

		logrus.Println(req)

		from, err := time.Parse(req.From.Format, req.From.Value)
		if err != nil {
			logrus.WithError(err).Info("parse 'from' date")
			return c.Status(fiber.StatusUnprocessableEntity).SendString("bad 'from' date")
		}
		to, err := time.Parse(req.To.Format, req.To.Value)
		if err != nil {
			logrus.WithError(err).Info("parse 'to' date")
			return c.Status(fiber.StatusUnprocessableEntity).SendString("bad 'to' date")
		}

		return c.JSON(DateRangeResponse{
			SecondsRange: int64(to.Sub(from).Seconds()),
		})
	})

	lErr := webApp.Listen(":8080")
	if lErr != nil {
		logrus.WithError(lErr).Fatal("listen port")
	}
}
*/

package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
	"time"
)

type (
	SendPushNotificationRequest struct {
		Message string `json:"message"`
		UserID  int64  `json:"user_id"`
	}

	PushNotification struct {
		Message string `json:"message"`
		UserID  int64  `json:"user_id"`
	}
)

var pushNotificationsQueue []PushNotification

func main() {
	webApp := fiber.New(fiber.Config{
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
	webApp.Use(recover.New())

	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	webApp.Post("/push/send", func(ctx *fiber.Ctx) error {
		var req SendPushNotificationRequest
		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
		}

		pushNotificationsQueue = append(pushNotificationsQueue, PushNotification{
			Message: req.Message,
			UserID:  req.UserID,
		})
		if len(pushNotificationsQueue) > 3 {
			panic("Queue is full")
		}

		return ctx.SendStatus(fiber.StatusOK)
	})

	logrus.Fatal(webApp.Listen(":8080"))
}
