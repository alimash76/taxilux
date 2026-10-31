package migrate

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Run(db *sql.DB) error {
	dir := os.Getenv("MIGRATIONS_DIR")
	if dir == "" {
		dir = "/migrations"
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version    TEXT PRIMARY KEY,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migrations directory %q: %w", dir, err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, file := range files {
		version := strings.TrimSuffix(file, ".up.sql")

		var applied int
		if err := db.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = ?`, version).Scan(&applied); err != nil {
			return fmt.Errorf("check migration %q: %w", version, err)
		}
		if applied > 0 {
			continue
		}

		raw, err := os.ReadFile(filepath.Join(dir, file))
		if err != nil {
			return fmt.Errorf("read migration %q: %w", file, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin migration %q: %w", version, err)
		}
		if _, err := tx.Exec(string(raw)); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("apply migration %q failed: %w; rollback: %w", version, err, rbErr)
			}
			return fmt.Errorf("apply migration %q: %w", version, err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, version); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("record migration %q failed: %w; rollback: %w", version, err, rbErr)
			}
			return fmt.Errorf("record migration %q: %w", version, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %q: %w", version, err)
		}
	}

	return nil
}
