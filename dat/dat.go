package dat

import (
	"fmt"

	"github.com/etzba/gosqlite/pkg/logger"
	_ "modernc.org/sqlite"
)

const dbName = "gosqlite.db"
const schema = `
CREATE TABLE IF NOT EXISTS users (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	username      TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	profile_id    INTEGER
);
CREATE TABLE IF NOT EXISTS profiles (
	id      INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL REFERENCES users(id),
	posts   TEXT NOT NULL DEFAULT '' CHECK (length(posts) <= 50)
);`

// InitDB prepare the database and run db migration
func (c *Context) InitDB(logger logger.Logger) error {
	logger.Info(fmt.Sprintf("Creating schema in database %s", dbName))
	if _, err := c.Db.Exec(schema); err != nil {
		logger.Error("failed to execute schema", err)
		return err
	}
	return nil
}
