package main

import (
	"context"
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

	// Initialize handler
	userHandler = handlers.NewUserHandler(db, appLogger)

	appLogger.Info("Lambda function initialized successfully")
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return userHandler.HandleRequest(ctx, request)
}

func main() {
	lambda.Start(Handler)
}
