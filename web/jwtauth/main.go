package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type (
	SignUpRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	SignInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	SignInResponse struct {
		JWTToken string `json:"jwt_token"`
	}

	ProfileResponse struct {
		Email string `json:"email"`
	}

	User struct {
		Email    string
		password string
	}
)

var (
	webApiPort = ":8080"

	users = map[string]User{}

	secretKey = []byte("qwerty123456")

	contextKeyUser = "user"
)

func main() {
	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// BEGIN (write your solution here) (write your solution here)

	// END

	logrus.Fatal(webApp.Listen(webApiPort))
}
