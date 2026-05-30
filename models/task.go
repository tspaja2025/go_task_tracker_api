package models

import "time"

// Database layout for an individual task record
type Task struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueAt       *time.Time `json:"due_at"` // Allows field to be nil
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Validation expectations for new task
type CreateTaskRequest struct {
	Title       string     `json:"title" validate:"required,min=3,max=100"`
	Description string     `json:"description" validate:"max=1000"`
	Status      string     `json:"status" validate:"omitempty,oneof=pending in_progress completed"`
	Priority    string     `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueAt       *time.Time `json:"due_at"`
}

// Validation for altering a task
type UpdateTaskRequest struct {
	Title       *string    `json:"title" validate:"required,min=3,max=100"`
	Description *string    `json:"description" validate:"max=1000"`
	Status      *string    `json:"status" validate:"omitempty,oneof=pending in_progress completed"`
	Priority    *string    `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueAt       *time.Time `json:"due_at"`
}
