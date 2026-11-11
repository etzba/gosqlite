package main

import (
	"github.com/etzba/gosqlite/dat"
	"github.com/etzba/gosqlite/pkg/logger"
	"github.com/etzba/gosqlite/server"
)

func main() {
	logger := logger.New()
	if err := dat.InitDB(logger); err != nil {
		logger.Error("Failed to initilize database", err)
	}

	server := server.New(logger, ":8080")
	if err := server.Run(); err != nil {
		panic(err)
	}
}
