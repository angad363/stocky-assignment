package reward

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	service     *RewardService
	idemService *IdempotencyService
}

func NewRewardHandler(service *RewardService, idem *IdempotencyService) *RewardHandler {
	return &RewardHandler{service: service, idemService: idem}
}

func (h *RewardHandler) CreateReward(c *gin.Context) {
	var req RewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	idemKey := c.GetHeader("Idempotency-Key")
	if idemKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key header is required"})
		return
	}

	ctx := context.Background()

	// Step 1: Check idempotency
	exists, _ := h.idemService.CheckOrSet(ctx, idemKey, nil)
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Duplicate reward request detected"})
		return
	}


	// Step 2: Create reward
	reward, err := h.service.CreateReward(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reward"})
		return
	}

	// Step 3: Save response in cache for idempotency replay
	_, _ = h.idemService.CheckOrSet(ctx, idemKey, reward)

	c.JSON(http.StatusCreated, reward)
}

func (h *RewardHandler) GetTodayRewards(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	ctx := context.Background()
	rewards, err := h.service.GetTodayRewards(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch rewards"})
		return
	}

	// make sure we return an empty list not null
	if rewards == nil {
		rewards = []Reward{}
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"rewards_today": rewards,
	})
}

func (h *RewardHandler) GetHistoricalINR(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	data, err := h.service.GetHistoricalINR(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch historical INR"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":        userID,
		"historical_inr": data,
	})
}

func (h *RewardHandler) GetUserStats(c *gin.Context) {
	userIDParam := c.Param("userId")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	var userID int
	_, err := fmt.Sscan(userIDParam, &userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	ctx := context.Background()
	todaySummary, totalValue, err := h.service.GetUserStats(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stats"})
		return
	}

	if todaySummary == nil {
		todaySummary = make(map[string]float64)
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":             userID,
		"today_summary":       todaySummary,
		"portfolio_value_inr": totalValue,
	})
}
