package reward

import (
	"context"
	"net/http"
	"strconv"

	"github.com/angad363/stocky-assignment/pkg/logger"
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
		logger.Log.Warnf("Invalid reward request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	idemKey := c.GetHeader("Idempotency-Key")
	if idemKey == "" {
		logger.Log.Warn("Missing Idempotency-Key header")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key header is required"})
		return
	}

	ctx := context.Background()
	exists, _ := h.idemService.CheckOrSet(ctx, idemKey, nil)
	if exists {
		logger.Log.Warnf("Duplicate reward request detected for key: %s", idemKey)
		c.JSON(http.StatusConflict, gin.H{"error": "Duplicate reward request detected"})
		return
	}

	reward, err := h.service.CreateReward(ctx, req)
	if err != nil {
		logger.Log.Errorf("Failed to create reward: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reward"})
		return
	}

	_, _ = h.idemService.CheckOrSet(ctx, idemKey, reward)
	logger.Log.WithFields(map[string]interface{}{
		"user_id": req.UserID,
	}).Info("Reward successfully created")

	c.JSON(http.StatusCreated, reward)
}

func (h *RewardHandler) GetTodayRewards(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		logger.Log.Warn("Invalid userId in /today-stocks request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	rewards, err := h.service.GetTodayRewards(context.Background(), userID)
	if err != nil {
		logger.Log.Errorf("Failed to fetch today's rewards for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch rewards"})
		return
	}

	if rewards == nil {
		rewards = []Reward{}
	}

	logger.Log.WithField("user_id", userID).Info("Fetched today's rewards")
	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"rewards_today": rewards,
	})
}

func (h *RewardHandler) GetHistoricalINR(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		logger.Log.Warn("Invalid userId in /historical-inr request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	data, err := h.service.GetHistoricalINR(context.Background(), userID)
	if err != nil {
		logger.Log.Errorf("Failed to fetch historical INR for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch historical INR"})
		return
	}

	logger.Log.WithField("user_id", userID).Info("Fetched historical INR data")
	c.JSON(http.StatusOK, gin.H{
		"user_id":        userID,
		"historical_inr": data,
	})
}

func (h *RewardHandler) GetUserStats(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		logger.Log.Warn("Invalid userId in /stats request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	ctx := context.Background()
	todaySummary, totalValue, err := h.service.GetUserStats(ctx, userID)
	if err != nil {
		logger.Log.Errorf("Failed to fetch stats for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stats"})
		return
	}

	logger.Log.WithFields(map[string]interface{}{
		"user_id": userID,
		"value":   totalValue,
	}).Info("Fetched user stats")

	c.JSON(http.StatusOK, gin.H{
		"user_id":             userID,
		"today_summary":       todaySummary,
		"portfolio_value_inr": totalValue,
	})
}

func (h *RewardHandler) GetUserPortfolio(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		logger.Log.Warn("Invalid userId in /portfolio request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	portfolio, err := h.service.GetUserPortfolio(context.Background(), userID)
	if err != nil {
		logger.Log.Errorf("Failed to fetch portfolio for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch portfolio"})
		return
	}

	logger.Log.WithField("user_id", userID).Info("Fetched user portfolio")
	c.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"portfolio": portfolio,
	})
}
