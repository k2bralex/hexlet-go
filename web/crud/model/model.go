package model

import "github.com/satori/go.uuid"

type (
	GetTaskResponse struct {
		UUID     uuid.UUID `json:"id"`
		Desc     string    `json:"description"`
		Deadline int64     `json:"deadline"`
	}

	CreateTaskRequest struct {
		Desc     string `json:"description" validate:"required"`
		Deadline int64  `json:"deadline" validate:""`
	}

	CreateTaskResponse struct {
		UUID int64 `json:"id" validate:"required, uuid4_rfc4122"`
	}

	UpdateTaskRequest struct {
		Desc     string `json:"description" validate:"required"`
		Deadline int64  `json:"deadline" validate:""`
	}

	Task struct {
		UUID     uuid.UUID
		Desc     string
		Deadline int64
	}
)
