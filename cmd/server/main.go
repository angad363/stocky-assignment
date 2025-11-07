package main

import (
	"github.com/angad363/stocky-assignment/internal/config"
	"github.com/angad363/stocky-assignment/internal/log"
	"github.com/angad363/stocky-assignment/internal/server"
)

func main() {
	cfg := config.Load()
	logger := log.NewLogger()

	srv := server.NewServer(logger)
	srv.Run(cfg.ServerPort)
}