package main

import (
	"github.com/etzba/gosqlite/pkg/logger"
	"github.com/etzba/gosqlite/server"
)

func main() {
	logger := logger.New()
	server := server.New(logger, ":8080")
	if err := server.Run(); err != nil {
		panic(err)
	}
}
