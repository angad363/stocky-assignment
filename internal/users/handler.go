package users

import (
	"context"
	"net/http"

	"github.com/angad363/stocky-assignment/pkg/logger"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name string `json:"name" binding:"required"`
}

type RegisterResponse struct {
	User   User        `json:"user"`
	Reward interface{} `json:"reward"`
}

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warnf("Invalid registration request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, reward, err := h.service.CreateUser(context.Background(), req.Name)
	if err != nil {
		logger.Log.Errorf("Failed to create user '%s': %v", req.Name, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	logger.Log.WithField("user_name", req.Name).Info("New user registered successfully")
	c.JSON(http.StatusCreated, RegisterResponse{User: user, Reward: reward})
}
