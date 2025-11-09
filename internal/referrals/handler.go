package referral

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReferralHandler struct {
	service *ReferralService
}

func NewReferralHandler(service *ReferralService) *ReferralHandler {
	return &ReferralHandler{service: service}
}

type ReferralRequest struct {
	UserID     int    `json:"user_id" binding:"required"`
	FriendName string `json:"friend_name" binding:"required"`
}

func (h *ReferralHandler) CreateReferral(c *gin.Context) {
	var req ReferralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ref, reward, err := h.service.CreateReferral(context.Background(), req.UserID, req.FriendName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create referral"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Referral successful! Reward granted.",
		"referral": ref,
		"reward":   reward,
	})
}