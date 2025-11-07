package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Router *gin.Engine
	Log    *logrus.Logger
}

func NewServer(log *logrus.Logger) *Server {
	r := gin.Default()
	s := &Server{
		Router: r,
		Log:    log,
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.Router.GET("/health", func(c *gin.Context) {
		s.Log.Info("Health check called")
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func (s *Server) Run(port string) {
	s.Log.Infof("Starting server on port %s", port)
	s.Router.Run(":" + port)
}