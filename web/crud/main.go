package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	. "hexlet/web/crud/model"
	"log"
	"os"
)

var (
	tasks    = make(map[uuid.UUID]Task)
	validate = validator.New()
)

func main() {
	file, err := os.OpenFile(".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	webApp := fiber.New()
	webApp.Use(requestid.New())
	webApp.Use(logger.New(logger.Config{
		Format: "[${time}]: ${status} ${method} ${path} ${locals:requestid}\n",
		Output: file,
	}))

	webApp.Get("/", ServerAlive)
	webApp.Post("/tasks", Create)
	webApp.Patch("/tasks/:id", Update)
	webApp.Get("/tasks/:id", Read)
	webApp.Delete("/tasks/:id", Delete)

	logrus.Fatal(webApp.Listen(":8080"))
}

func ServerAlive(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func Create(ctx *fiber.Ctx) error {
	var create CreateTaskRequest
	if err := ctx.BodyParser(&create); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid JSON")
	}

	if err := validate.Struct(&create); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}

	newUUID := uuid.NewV4()
	tasks[newUUID] = Task{
		UUID:     newUUID,
		Desc:     create.Desc,
		Deadline: create.Deadline,
	}

	ctx.Status(fiber.StatusOK)
	return ctx.JSON(CreateTaskResponse{UUID: newUUID})
}

func Update(ctx *fiber.Ctx) error {
	var update UpdateTaskRequest
	if err := ctx.BodyParser(&update); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid JSON")
	}

	if err := validate.Struct(&update); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}

	id, err := uuid.FromString(ctx.Params("id"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	task, ok := tasks[id]
	if !ok {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	task.Desc = update.Desc
	task.Deadline = update.Deadline
	tasks[task.UUID] = task

	return ctx.SendStatus(fiber.StatusOK)
}

func Read(ctx *fiber.Ctx) error {
	id, err := uuid.FromString(ctx.Params("id"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	task, ok := tasks[id]
	if !ok {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	ctx.Status(fiber.StatusOK)
	return ctx.JSON(GetTaskResponse{
		UUID:     task.UUID,
		Desc:     task.Desc,
		Deadline: task.Deadline,
	})
}

func Delete(ctx *fiber.Ctx) error {
	id, err := uuid.FromString(ctx.Params("id"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if _, ok := tasks[id]; !ok {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	delete(tasks, id)
	return ctx.SendStatus(fiber.StatusOK)
}
