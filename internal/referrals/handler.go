package referral

import (
	"context"
	"net/http"

	"github.com/angad363/stocky-assignment/pkg/logger"
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
		logger.Log.Warnf("Invalid referral request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ref, reward, err := h.service.CreateReferral(context.Background(), req.UserID, req.FriendName)
	if err != nil {
		logger.Log.Errorf("Failed to create referral for user %d: %v", req.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create referral"})
		return
	}

	logger.Log.WithFields(map[string]interface{}{
		"user_id":     req.UserID,
		"friend_name": req.FriendName,
	}).Info("Referral created successfully")

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Referral successful! Reward granted.",
		"referral": ref,
		"reward":   reward,
	})
}
