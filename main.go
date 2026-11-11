package main

import (
	"github.com/etzba/gosqlite/pkg/logger"
	"github.com/etzba/gosqlite/server"
)

func main() {
	logger := logger.New()
	server, err := server.New(logger, ":8080")
	if err != nil {
		panic(err)
	}

	if err := server.Run(); err != nil {
		panic(err)
	}
}
