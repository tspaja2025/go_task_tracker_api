package tasks

import (
	"encoding/json"
	"main/middleware"
	"main/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	repo     *Repository
	validate *validator.Validate
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{
		repo:     repo,
		validate: validator.New(),
	}
}

// Dispatch incoming route requests based on HTTP method and paths
func (h *Handler) Router(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")

	// Route is "/tasks"
	if idStr == "" || idStr == "/" {
		switch r.Method {
		case http.MethodPost:
			h.createTask(w, r, userID)
		case http.MethodGet:
			h.listTasks(w, r, userID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Endpoint contains explicit ID string
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task ID formatting"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTask(w, r, id, userID)
	case http.MethodPut:
		h.updateTask(w, r, id, userID)
	case http.MethodDelete:
		h.deleteTask(w, r, id, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request, userID int) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid body formatting"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	task, err := h.repo.Create(r.Context(), userID, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to store item"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request, id, userID int) {
	task, err := h.repo.GetByID(r.Context(), id, userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request, id, userID int) {
	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task, err := h.repo.Update(r.Context(), id, userID, req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request, id, userID int) {
	if err := h.repo.Delete(r.Context(), id, userID); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request, userID int) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	filters := models.TaskQueryFilters{
		Page:     page,
		Limit:    limit,
		Status:   q.Get("status"),
		Priority: q.Get("priority"),
		SortBy:   q.Get("sort_by"),
		Order:    strings.ToUpper(q.Get("order")),
	}

	tasks, err := h.repo.List(r.Context(), userID, filters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch tasks"})
		return
	}

	// Return an empty array instead of null if no tasks exits
	if tasks == nil {
		tasks = []models.Task{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
