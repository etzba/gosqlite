package dat

import (
	"context"
	"database/sql"

	"github.com/etzba/gosqlite/pkg/logger"

	_ "modernc.org/sqlite"
)

type Context struct {
	ctx    context.Context
	Logger logger.Logger
	Db     *sql.DB
}

func NewContext(logger logger.Logger) (*Context, error) {
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		logger.Error("failed to open db", err)
		return nil, err
	}

	return &Context{
		ctx:    context.Background(),
		Logger: logger,
		Db:     db,
	}, nil
}
