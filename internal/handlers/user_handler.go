package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"go-lambda-crud/internal/database"
	"go-lambda-crud/internal/models"
	"go-lambda-crud/pkg/logger"

	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserHandler struct {
	db     *database.Database
	logger *logger.Logger
}

func NewUserHandler(db *database.Database, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		db:     db,
		logger: logger,
	}
}

func (h *UserHandler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	h.logger.Info("Received request", zap.String("path", request.Path), zap.String("method", request.HTTPMethod))

	switch request.HTTPMethod {
	case "GET":
		if request.PathParameters["id"] != "" {
			return h.GetUser(ctx, request)
		}
		return h.GetUsers(ctx, request)
	case "POST":
		return h.CreateUser(ctx, request)
	case "PUT":
		return h.UpdateUser(ctx, request)
	case "DELETE":
		return h.DeleteUser(ctx, request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       `{"error": "Method not allowed"}`,
		}, nil
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req models.CreateUserRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid request body"}`,
		}, nil
	}

	user := models.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	if err := h.db.DB.Create(&user).Error; err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Failed to create user"}`,
		}, nil
	}

	response, _ := json.Marshal(user)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(response),
	}, nil
}

func (h *UserHandler) GetUsers(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var users []models.User
	if err := h.db.DB.Find(&users).Error; err != nil {
		h.logger.Error("Failed to get users", zap.Error(err))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Failed to get users"}`,
		}, nil
	}

	response, _ := json.Marshal(users)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, err := strconv.Atoi(request.PathParameters["id"])
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid user ID"}`,
		}, nil
	}

	var user models.User
	if err := h.db.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       `{"error": "User not found"}`,
			}, nil
		}
		h.logger.Error("Failed to get user", zap.Error(err))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Failed to get user"}`,
		}, nil
	}

	response, _ := json.Marshal(user)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, err := strconv.Atoi(request.PathParameters["id"])
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid user ID"}`,
		}, nil
	}

	var req models.UpdateUserRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid request body"}`,
		}, nil
	}

	var user models.User
	if err := h.db.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       `{"error": "User not found"}`,
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Failed to get user"}`,
		}, nil
	}

	if err := h.db.DB.Model(&user).Updates(models.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}).Error; err != nil {
		h.logger.Error("Failed to update user", zap.Error(err))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Failed to update user"}`,
		}, nil
	}

	response, _ := json.Marshal(user)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, err := strconv.Atoi(request.PathParameters["id"])
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid user ID"}`,
		}, nil
	}

	if err := h.db.DB.Delete(&models.User{}, id).Error; err != nil {
		h.logger.Error("Failed to delete user", zap.Error(err))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Failed to delete user"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message": "User deleted successfully"}`,
	}, nil
}
