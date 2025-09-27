package main

import (
	"log"
	"net/http"
	"os"

	"go-lambda-crud/configs"
	"go-lambda-crud/internal/database"
	"go-lambda-crud/internal/models"
	"go-lambda-crud/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := configs.LoadConfig()

	// Initialize logger
	appLogger := logger.NewLogger()

	// Initialize database
	db, err := database.NewDatabase(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		appLogger.Error("Failed to connect to database", zap.Error(err))
		panic(err)
	}

	appLogger.Info("Local server starting...")

	router := gin.Default()

	// User routes
	router.POST("/users", createUser(db, appLogger))
	router.GET("/users", getUsers(db, appLogger))
	router.GET("/users/:id", getUser(db, appLogger))
	router.PUT("/users/:id", updateUser(db, appLogger))
	router.DELETE("/users/:id", deleteUser(db, appLogger))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}

func createUser(db *database.Database, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := models.User{
			Name:  req.Name,
			Email: req.Email,
			Age:   req.Age,
		}

		if err := db.DB.Create(&user).Error; err != nil {
			logger.Error("Failed to create user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func getUsers(db *database.Database, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := db.DB.Find(&users).Error; err != nil {
			logger.Error("Failed to get users", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func getUser(db *database.Database, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user models.User

		if err := db.DB.First(&user, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			logger.Error("Failed to get user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func updateUser(db *database.Database, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req models.UpdateUserRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if err := db.DB.First(&user, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}

		if err := db.DB.Model(&user).Updates(models.User{
			Name:  req.Name,
			Email: req.Email,
			Age:   req.Age,
		}).Error; err != nil {
			logger.Error("Failed to update user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func deleteUser(db *database.Database, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := db.DB.Delete(&models.User{}, id).Error; err != nil {
			logger.Error("Failed to delete user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
