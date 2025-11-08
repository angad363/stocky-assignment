package server

import (
	"github.com/angad363/stocky-assignment/internal/price"
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

	s := &Server{
		router: r,
		logger: logger,
	}

	s.registerRoutes(priceHandler)
	return s
}

func (s *Server) registerRoutes(priceHandler *price.PriceHandler) {
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	s.router.GET("/price", priceHandler.GetPrice)
}

func (s *Server) Start(port string) {
	s.logger.Infof("Starting server on port %s", port)
	s.router.Run(":" + port)
}