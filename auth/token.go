package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// For production, load from environment variables
var jwtSecret = []byte("super-secret-development-key")

// Create both an access token and a refresh token
func GenerateTokens(userID int) (string, string, time.Time, error) {
	// Create access token (15 min expiration timer)
	accessExpiration := time.Now().Add(15 * time.Minute)
	accessTokenClaims := jwt.MapClaims{
		"sub": userID,
		"exp": accessExpiration.Unix(),
		"iat": time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessStr, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Create refresh token (7 days expiration timer)
	refreshExpiration := time.Now().Add(7 * 24 * time.Hour)
	refreshTokenClaims := jwt.MapClaims{
		"sub": userID,
		"exp": refreshExpiration.Unix(),
		"iat": time.Now().Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshStr, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessStr, refreshStr, refreshExpiration, nil
}

// Convert raw token string into secure hash for database
func HashToken(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Validate a token string and return if valid
func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		// Ensure the signing method is HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenRequiredClaimMissing
}
