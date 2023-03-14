package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type (
	GetTaskResponse struct {
		ID       int64  `json:"id"`
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	CreateTaskRequest struct {
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	CreateTaskResponse struct {
		ID int64 `json:"id"`
	}

	UpdateTaskRequest struct {
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	Task struct {
		ID       int64
		Desc     string
		Deadline int64
	}
)

var (
	taskIDCounter int64 = 1
	tasks               = make(map[int64]Task)
)

func main() {
	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	webApp.Post("/tasks", func(ctx *fiber.Ctx) error {
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
	})

	webApp.Patch("/tasks/:id", func(ctx *fiber.Ctx) error {
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
	})

	webApp.Get("/tasks/:id", func(ctx *fiber.Ctx) error {
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
	})

	webApp.Delete("/tasks/:id", func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		if _, ok := tasks[int64(id)]; !ok {
			return ctx.SendStatus(fiber.StatusNotFound)
		}

		delete(tasks, int64(id))
		return ctx.SendStatus(fiber.StatusOK)
	})

	logrus.Fatal(webApp.Listen(":8080"))
}
