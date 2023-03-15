package main

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
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
		return c.SendStatus(fiber.StatusOK)
	})

	publicGroup := webApp.Group("")
	publicGroup.Post("/signup", SignUp)
	publicGroup.Post("/signin", SignIn)

	authorizedGroup := webApp.Group("")
	authorizedGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: secretKey,
		ContextKey: contextKeyUser,
	}))
	authorizedGroup.Get("/profile", ViewProfile)

	logrus.Fatal(webApp.Listen(webApiPort))
}

func SignUp(ctx *fiber.Ctx) error {
	var signUpReq SignUpRequest
	err := ctx.BodyParser(&signUpReq)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).SendString("invalid JSON")
	}

	if _, ok := users[signUpReq.Email]; ok {
		return ctx.SendStatus(fiber.StatusConflict)
	}

	users[signUpReq.Email] = User{
		Email:    signUpReq.Email,
		password: signUpReq.Email,
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func SignIn(ctx *fiber.Ctx) error {
	var signInReq SignInRequest
	err := ctx.BodyParser(&signInReq)
	if err != nil {
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}

	user, ok := users[signInReq.Email]
	if !ok {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	if user.password != signInReq.Password {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": user.Email})
	jwtToken, err := tkn.SignedString(secretKey)
	if err != nil {
		logrus.WithError(err).Error("JWT token signing")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.Status(fiber.StatusOK).JSON(SignInResponse{JWTToken: jwtToken})
}

func ViewProfile(ctx *fiber.Ctx) error {
	payload, ok := jwtPayloadFromRequest(ctx)
	if !ok {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	user, ok := users[payload["sub"].(string)]
	if !ok {
		return errors.New("user not found")
	}

	return ctx.Status(fiber.StatusOK).JSON(ProfileResponse{Email: user.Email})
}

func jwtPayloadFromRequest(ctx *fiber.Ctx) (jwt.MapClaims, bool) {
	jwtToken, ok := ctx.Context().Value(contextKeyUser).(*jwt.Token)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"jwt_token_context_value": ctx.Context().Value(contextKeyUser),
		}).Error("wrong type of JWT token in context")
		return nil, false
	}

	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"jwt_token_claims": jwtToken.Claims,
		}).Error("wrong type of JWT token claims")
		return nil, false
	}

	return payload, true
}
