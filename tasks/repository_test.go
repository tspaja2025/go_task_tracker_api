package tasks

import (
	"main/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateTask_Success(t *testing.T) {
	// Create database mock connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a database connection", err)
	}
	defer db.Close()

	// Set up input data parameters
	userID := 1
	now := time.Now()
	req := models.CreateTaskRequest{
		Title:       "Test task",
		Description: "Testing repository layer",
		Status:      "pending",
		Priority:    "medium",
		DueAt:       &now,
	}

	// Exact query pattern and the mock rows expected to return
	expectedQuery := `INSERT INTO tasks \(user_id, title, description, status, priority, due_at\)`

	mockRows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "status", "priority", "due_at", "created_at", "updated_at"}).AddRow(
		1, userID, req.Title, req.Description, req.Status, req.Priority, req.DueAt, now, now,
	)

	mock.ExpectQuery(expectedQuery).WithArgs(userID, req.Title, req.Description, req.Status, req.Priority, req.DueAt).WillReturnRows(mockRows)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
