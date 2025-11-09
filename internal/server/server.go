package server

import (
	"github.com/angad363/stocky-assignment/internal/price"
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

func NewServer(logger *logrus.Logger,  conn *sqlx.DB) *Server {
	r := gin.Default()

	// Initialize Redis for price service
	price.InitRedis()

	// Initialize price service
	priceService := price.NewPriceService(price.RedisConn)
	priceHandler := price.NewPriceHandler(priceService)

	price.StartPriceUpdater(priceService, conn)

	idemService := reward.NewIdempotencyService(price.RedisConn)
	rewardService := reward.NewRewardService(conn, priceService)
	rewardHandler := reward.NewRewardHandler(rewardService, idemService)

	userService := users.NewUserService(conn, rewardService)
	userHandler := users.NewUserHandler(userService)


	s := &Server{
		router: r,
		logger: logger,
	}

	s.registerRoutes(priceHandler, rewardHandler, userHandler)
	return s
}

func (s *Server) registerRoutes(priceHandler *price.PriceHandler,
								rewardHandler *reward.RewardHandler,
								userHandler *users.UserHandler,
							) {
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	s.router.GET("/price", priceHandler.GetPrice)
	s.router.POST("/reward", rewardHandler.CreateReward)
	s.router.POST("/register", userHandler.Register)
	s.router.GET("/today-stocks/:userId", rewardHandler.GetTodayRewards)
	s.router.GET("/historical-inr/:userId", rewardHandler.GetHistoricalINR)

}

func (s *Server) Start(port string) {
	s.logger.Infof("Starting server on port %s", port)
	s.router.Run(":" + port)
}