package middleware

import (
	"context"
	"main/auth"
	"net/http"
	"strings"
)

// Private key type to prevent context collisions
type contextKey string

const UserIDKey contextKey = "userID"

// Inspect the JWT and inject the user ID into the request context
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Authorization header required}`, http.StatusUnauthorized)
			return
		}

		// "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, `{"error": "Authorization header must be Bearer token}`, http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		// Parse and validate the token
		claims, err := auth.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid or Expired access token}`, http.StatusUnauthorized)
			return
		}

		// Extract user ID
		sub, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, `{"error": "Invalid token payload"}`, http.StatusUnauthorized)
			return
		}
		userID := int(sub)

		// Add user ID to the request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// Pass the updated context to the handler
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Helper function for handlers to make it easier to get authenticated user's ID
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}
