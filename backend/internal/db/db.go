package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	for _, pragma := range []string{"journal_mode=WAL", "foreign_keys=ON", "busy_timeout=5000"} {
		if _, err := database.Exec("PRAGMA " + pragma); err != nil {
			database.Close()
			return nil, fmt.Errorf("PRAGMA %s: %w", pragma, err)
		}
	}

	return database, nil
}
