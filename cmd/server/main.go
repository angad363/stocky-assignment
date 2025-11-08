package main

import (
	"github.com/angad363/stocky-assignment/internal/config"
	"github.com/angad363/stocky-assignment/internal/db"
	"github.com/angad363/stocky-assignment/internal/log"
	"github.com/angad363/stocky-assignment/internal/server"
)

func main() {
	cfg := config.Load()
	logger := log.NewLogger()

	conn := db.Connect(cfg)
	defer conn.Close()

	srv := server.NewServer(logger, conn)
	srv.Start(cfg.ServerPort)
}
