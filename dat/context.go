package dat

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/etzba/gosqlite/pkg/logger"

	_ "modernc.org/sqlite"
)

// Context acts as a bridge between the application and relational database
type Context struct {
	ctx    context.Context
	Logger logger.Logger
	Db     *sql.DB
}

// NewContext returns an instance of database context with db connection and logger
func NewContext(logger logger.Logger) (*Context, error) {
	logger.Info(fmt.Sprintf("Connecting to database %s and setting database context", dbName))
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
