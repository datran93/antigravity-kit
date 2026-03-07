package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const cacheExpirySeconds = 30 * 24 * 60 * 60 // 1 month

func initDB() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)
	dbPath := filepath.Join(exeDir, "research_cache.db")

	var errOpen error
	db, errOpen = sql.Open("sqlite3", dbPath)
	if errOpen != nil {
		return errOpen
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cache (
			key TEXT PRIMARY KEY,
			value TEXT,
			timestamp REAL
		)
	`)
	if err != nil {
		return err
	}
	return nil
}

func getCache(key string) *string {
	if db == nil {
		return nil
	}
	var value string
	var timestamp float64
	err := db.QueryRow("SELECT value, timestamp FROM cache WHERE key = ?", key).Scan(&value, &timestamp)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("DB query error: %v\n", err)
		}
		return nil
	}

	now := float64(time.Now().Unix())
	if now-timestamp < cacheExpirySeconds {
		return &value
	}
	return nil
}

func setCache(key string, value string) {
	if db == nil {
		return
	}
	now := float64(time.Now().Unix())
	_, err := db.Exec(`
			INSERT OR REPLACE INTO cache (key, value, timestamp)
			VALUES (?, ?, ?)
		`, key, value, now)
	if err != nil {
		log.Printf("DB insert error: %v\n", err)
	}
}

func generateCacheKey(prefix string, content string) string {
	hash := md5.Sum([]byte(content))
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(hash[:]))
}
