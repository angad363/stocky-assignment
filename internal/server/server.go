package server

import (
	"time"

	"github.com/angad363/stocky-assignment/internal/price"
	referral "github.com/angad363/stocky-assignment/internal/referrals"
	"github.com/angad363/stocky-assignment/internal/reward"
	"github.com/angad363/stocky-assignment/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router *gin.Engine
	logger *logrus.Logger
}

func NewServer(logger *logrus.Logger, conn *sqlx.DB) *Server {
	r := gin.New()

	r.Use(gin.Recovery())

	// Global request logging middleware
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		logger.WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"latency":  duration.String(),
			"clientIP": c.ClientIP(),
		}).Info("HTTP request processed")
	})

	logger.Info("🔧 Initializing Redis and Price services...")

	price.InitRedis()
	priceService := price.NewPriceService(price.RedisConn)
	priceHandler := price.NewPriceHandler(priceService)

	price.StartPriceUpdater(priceService, conn)
	logger.Info("💹 Price updater started")

	idemService := reward.NewIdempotencyService(price.RedisConn)
	rewardService := reward.NewRewardService(conn, priceService)
	rewardHandler := reward.NewRewardHandler(rewardService, idemService)

	userService := users.NewUserService(conn, rewardService)
	userHandler := users.NewUserHandler(userService)

	referralService := referral.NewReferralService(conn, rewardService)
	referralHandler := referral.NewReferralHandler(referralService)

	s := &Server{
		router: r,
		logger: logger,
	}

	s.registerRoutes(priceHandler, rewardHandler, userHandler, referralHandler)

	logger.Info("✅ Routes registered successfully")

	return s
}

func (s *Server) registerRoutes(priceHandler *price.PriceHandler,
	rewardHandler *reward.RewardHandler,
	userHandler *users.UserHandler,
	referralHandler *referral.ReferralHandler,
) {
	s.logger.Info("🛣 Registering routes...")

	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
		s.logger.Debug("Health check endpoint called")
	})

	s.router.GET("/price", priceHandler.GetPrice)
	s.router.POST("/reward", rewardHandler.CreateReward)
	s.router.POST("/register", userHandler.Register)
	s.router.GET("/today-stocks/:userId", rewardHandler.GetTodayRewards)
	s.router.GET("/historical-inr/:userId", rewardHandler.GetHistoricalINR)
	s.router.GET("/stats/:userId", rewardHandler.GetUserStats)
	s.router.POST("/refer", referralHandler.CreateReferral)
	s.router.GET("/portfolio/:userId", rewardHandler.GetUserPortfolio)

	s.logger.Info("📡 All API routes registered")
}

func (s *Server) Start(port string) {
	s.logger.WithField("port", port).Info("Starting HTTP server")
	err := s.router.Run(":" + port)
	if err != nil {
		s.logger.WithError(err).Fatal("Failed to start server")
	}
}
