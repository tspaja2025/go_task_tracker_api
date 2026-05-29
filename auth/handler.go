package auth

import (
	"encoding/json"
	"main/models"
	"net/http"
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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	// Decode raw body JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON format"})
		return
	}

	// Validate struct tags
	if err := h.validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "Validation failed: " + err.Error()})
		return
	}

	// User persist via repository
	user, err := h.repo.CreateUser(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// Duplicate error handling hook for uniqueness constraints
		if strings.Contains(err.Error(), "duplicate key") {
			json.NewEncoder(w).Encode(map[string]string{"error": "Username or Email already exists"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	// Return created user object
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
