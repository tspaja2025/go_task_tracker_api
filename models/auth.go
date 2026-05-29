package models

// Define the payload for authentication
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Return JWT credentials to the client
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Define the payload for requesting new access token
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
