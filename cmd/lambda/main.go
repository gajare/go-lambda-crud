package main

import (
	"context"
	"log"
	"os"

	"go-lambda-crud/configs"
	"go-lambda-crud/internal/database"
	"go-lambda-crud/internal/handlers"
	"go-lambda-crud/pkg/logger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var userHandler *handlers.UserHandler

func init() {
	log.Println("Initializing Lambda function...")

	// Load configuration
	cfg := configs.LoadConfig()

	// Initialize logger
	appLogger := logger.NewLogger()

	// Initialize database (handle errors gracefully)
	db, err := database.NewDatabase(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		appLogger.Error("Database connection failed", zap.Error(err))
		// Continue without database for health checks
	} else {
		userHandler = handlers.NewUserHandler(db, appLogger)
		appLogger.Info("Database connected successfully")
	}

	appLogger.Info("Lambda function initialized")
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Health check endpoint (works without database)
	if request.Path == "/health" && request.HTTPMethod == "GET" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"status": "healthy", "service": "go-crud-lambda"}`,
		}, nil
	}

	// If database is connected, use the full handler
	if userHandler != nil {
		return userHandler.HandleRequest(ctx, request)
	}

	// Database not available
	return events.APIGatewayProxyResponse{
		StatusCode: 503,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       `{"error": "Service unavailable", "message": "Database connection not available"}`,
	}, nil
}

func main() {
	// Check if running in Lambda environment
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		lambda.Start(Handler)
	} else {
		// Local development mode
		log.Println("Running in local mode (not Lambda)")

		// Start a simple HTTP server for testing
		log.Println("Local HTTP server would start here...")
		log.Println("Use 'go run cmd/local/main.go' for local development")
	}
}
