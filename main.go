package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Email     string         `gorm:"unique;not" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	CreatedAt time.Time      `gorm:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func main() {
	// Initialize Database
	db := InitDB()

	r := gin.Default()

	// Public routes
	r.POST("/register", func(c *gin.Context) { /* Call register logic */ })
	r.POST("/login", func(c *gin.Context) { /* Call login logic */ })
	r.POST("/refresh", func(c *gin.Context) { /* Call refresh token logic */ })

	// Protected routes
	protected := r.Group("/")
	protected.Use(AuthMiddleware("jwt_secret_here"))
	{
		protected.POST("/tasks", func(c *gin.Context) { /* Create Task */ })
		protected.GET("/tasks", func(c *gin.Context) { /* List tasks with Scopes(paginate) */ })
		protected.PUT("/tasks", func(c *gin.Context) { /* Update Task + check ownership */ })
		protected.DELETE("/tasks", func(c *gin.Context) { /* Delete Task + chek ownership */ })
	}

	fmt.Println("Server running at: http://localhost:3000")
	r.Run(":3000")
}

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauhorized"})
			c.Abort()
			return
		}

		// Expecting "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		// Extract user ID and pass it down to handlers
		userID := uint(claims["user_id"].(float64))
		c.Set("userID", userID)
		c.Next()
	}
}

func InitDB() *gorm.DB {
	// Load from environment variables
	dsn := "host=localhost user=admin password=secretpassword dbname=task_tracker port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automigrate creates tables automatically based of go structs
	err = db.AutoMigrate(&User{}, &Task{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Database successfully connected and migrated!")
	return db
}

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page <= 0 {
			page = 1
		}

		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}
