package model

import "github.com/satori/go.uuid"

type (
	GetTaskResponse struct {
		UUID     uuid.UUID `json:"id"`
		Desc     string    `json:"description"`
		Deadline int64     `json:"deadline"`
	}

	CreateTaskRequest struct {
		Desc     string `json:"description" validate:"required, min=1, max=150"`
		Deadline int64  `json:"deadline" validate:"required, gt=1678857313"`
	}

	CreateTaskResponse struct {
		UUID uuid.UUID `json:"id"`
	}

	UpdateTaskRequest struct {
		Desc     string `json:"description" validate:"required, min=1, max=150"`
		Deadline int64  `json:"deadline" validate:"required, gt=1678857313"`
	}

	Task struct {
		UUID     uuid.UUID
		Desc     string
		Deadline int64
	}
)
