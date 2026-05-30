package tasks

import (
	"context"
	"errors"
	"fmt"
	"main/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create stores brand new task linked to the calling user
func (r *Repository) Create(ctx context.Context, userID int, req models.CreateTaskRequest) (*models.Task, error) {
	// Set default values if omitted by request payload
	if req.Status == "" {
		req.Status = "pending"
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}

	query := `INSERT INTO tasks (user_id, title, description, status, priority, due_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, user_id, title, description, status, priority, due_at, created_at, updated_at;`

	var t models.Task
	err := r.db.QueryRow(ctx, query, userID, req.Title, req.Description, req.Status, req.Priority, req.DueAt).Scan(
		&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueAt, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	return &t, nil
}

// Find a single task, ensure that the task belongs to the requesting user
func (r *Repository) GetByID(ctx context.Context, id, userID int) (*models.Task, error) {
	query := `SELECT id, user_id, title, description, status, priority, due_at, created_at, updated_at FROM tasks WHERE id = $1 AND user_id = $2;`

	var t models.Task
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueAt, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &t, nil
}

// Rewrite changed parameters on an existing task asset
func (r *Repository) Update(ctx context.Context, id, userID int, req models.UpdateTaskRequest) (*models.Task, error) {
	// Fetch current version to ensure ownership exists
	current, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Overwrite field that were provided in patch payload
	if req.Title != nil {
		current.Title = *req.Title
	}
	if req.Description != nil {
		current.Description = *req.Description
	}
	if req.Status != nil {
		current.Status = *req.Status
	}
	if req.Priority != nil {
		current.Priority = *req.Priority
	}
	if req.DueAt != nil {
		current.DueAt = req.DueAt
	}

	query := `UPDATE TASKS SET title = $1, description = $2, status = $3, priority = $4, due_at = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6 AND user_id = $7 RETURNING id, user_id, title, description, status, priority, due_at, created_at, updated_at;`

	var t models.Task
	err = r.db.QueryRow(ctx, query, current.Title, current.Description, current.Status, current.DueAt, id, userID).Scan(
		&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueAt, &t.CreatedAt, &t.UpdatedAt,
	)
	return &t, err
}

// Remove a task from database matching both IDs
func (r *Repository) Delete(ctx context.Context, id, userID int) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("task not found or unauthorized")
	}
	return nil
}
