package auth

import (
	"encoding/json"
	"main/models"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON format"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "Validation failed"})
		return
	}

	// Find user by email
	user, err := h.repo.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		// Do not tell user that the email was not found
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}

	// Generate auth tokens
	accessToken, refreshToken, refreshExpiration, err := GenerateTokens(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate tokens"})
		return
	}

	// Save hashed refresh token to database
	tokenhash := HashToken(refreshToken)
	err = h.repo.SaveRefreshToken(r.Context(), user.ID, tokenhash, refreshExpiration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save session"})
	}

	// Send tokens back to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
