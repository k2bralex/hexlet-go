package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
	. "hexlet/web/crud/model"
	"log"
	"os"
)

var (
	taskIDCounter int64 = 1
	tasks               = make(map[int64]Task)
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
		Format: "$[{time}]: ${status} ${method} ${path} ${locals:requestid}\n",
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
	err := ctx.BodyParser(&create)
	if err != nil {
		return err
	}

	tasks[taskIDCounter] = Task{
		ID:       taskIDCounter,
		Desc:     create.Desc,
		Deadline: create.Deadline,
	}
	taskIDCounter++

	ctx.Status(fiber.StatusOK)
	return ctx.JSON(CreateTaskResponse{ID: taskIDCounter - 1})
}

func Update(ctx *fiber.Ctx) error {
	var update UpdateTaskRequest
	err := ctx.BodyParser(&update)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	task, ok := tasks[int64(id)]
	if !ok {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	task.Desc = update.Desc
	task.Deadline = update.Deadline
	tasks[task.ID] = task

	return ctx.SendStatus(fiber.StatusOK)
}

func Read(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	task, ok := tasks[int64(id)]
	if !ok {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	ctx.Status(fiber.StatusOK)
	return ctx.JSON(GetTaskResponse{
		ID:       task.ID,
		Desc:     task.Desc,
		Deadline: task.Deadline,
	})
}

func Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if _, ok := tasks[int64(id)]; !ok {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	delete(tasks, int64(id))
	return ctx.SendStatus(fiber.StatusOK)
}
