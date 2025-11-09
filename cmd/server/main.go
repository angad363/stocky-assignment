package main

import (
	"github.com/angad363/stocky-assignment/internal/config"
	"github.com/angad363/stocky-assignment/internal/db"
	"github.com/angad363/stocky-assignment/internal/server"
	"github.com/angad363/stocky-assignment/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	logger.Init()
	logger.Log.Info("Starting Stocky Server initialization...")

	cfg := config.Load()

	conn := db.Connect(cfg)
	defer conn.Close()

	srv := server.NewServer(logger.Log, conn)
	srv.Start(cfg.ServerPort)
}
