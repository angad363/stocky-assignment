package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRequest is the expected payload for /register
type RegisterRequest struct {
	Name string `json:"name" binding:"required"`
}

type RegisterResponse struct {
	User   User          `json:"user"`
	Reward interface{}   `json:"reward"`
}

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Register creates a user and auto-rewards them
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, reward, err := h.service.CreateUser(context.Background(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		User:   user,
		Reward: reward,
	})
}
