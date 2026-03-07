package main

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// GetDBConnection opens a connection to the SQLite database and initializes the schema.
func GetDBConnection(workspacePath string) (*sql.DB, error) {
	if workspacePath == "" {
		return nil, fmt.Errorf("workspace_path is required")
	}

	absPath, err := filepath.Abs(workspacePath)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace path: %v", err)
	}

	dbPath := filepath.Join(absPath, "context.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := initializeSchema(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func initializeSchema(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS checkpoints (
            task_id TEXT PRIMARY KEY,
            description TEXT,
            status TEXT,
            completed_steps TEXT,
            next_steps TEXT,
            active_files TEXT,
            notes TEXT,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS intents (
            task_id TEXT PRIMARY KEY,
            tactic TEXT,
            locked_files TEXT,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS drift_tracker (
            task_id TEXT PRIMARY KEY,
            failure_count INTEGER DEFAULT 0,
            last_failed_at TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS anchors (
            key TEXT PRIMARY KEY,
            value TEXT,
            rule TEXT,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS file_annotations (
            file_path TEXT PRIMARY KEY,
            gotchas TEXT,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS knowledge_fts USING fts5(
            tactic_name,
            ki_path UNINDEXED,
            summary,
            decisions
        )`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to initialize schema query:\n%s\nerror: %v", query, err)
		}
	}

	return nil
}
