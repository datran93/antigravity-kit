package main

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	enginesCache = make(map[string]*gorm.DB)
	redisCache   = make(map[string]*redis.Client)
	cacheMutex   sync.Mutex
)

func getGormDB(connStr string) (*gorm.DB, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if db, ok := enginesCache[connStr]; ok {
		// Ping to ensure health
		sqlDB, err := db.DB()
		if err == nil && sqlDB.Ping() == nil {
			return db, nil
		}
	}

	var dialector gorm.Dialector

	if strings.HasPrefix(connStr, "postgresql://") || strings.HasPrefix(connStr, "postgres://") {
		dialector = postgres.Open(connStr)
	} else if strings.HasPrefix(connStr, "mysql+pymysql://") || strings.HasPrefix(connStr, "mysql+aiomysql://") || strings.HasPrefix(connStr, "mysql://") {
		connStrClean := strings.ReplaceAll(connStr, "mysql+pymysql://", "mysql://")
		connStrClean = strings.ReplaceAll(connStrClean, "mysql+aiomysql://", "mysql://")
		parsed, err := url.Parse(connStrClean)
		if err != nil {
			return nil, fmt.Errorf("invalid mysql url: %v", err)
		}
		pass, _ := parsed.User.Password()
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", parsed.User.Username(), pass, parsed.Host, strings.TrimPrefix(parsed.Path, "/"))
		dialector = mysql.Open(dsn)
	} else if strings.HasPrefix(connStr, "sqlite://") {
		return nil, fmt.Errorf("unsupported or unrecognized database connection string format")
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		DisableAutomaticPing: false,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(5)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	enginesCache[connStr] = db
	return db, nil
}

func getRedisClient(connStr string) (*redis.Client, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if client, ok := redisCache[connStr]; ok {
		return client, nil
	}

	opts, err := redis.ParseURL(connStr)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	redisCache[connStr] = client
	return client, nil
}

func isRedis(connStr string) bool {
	return strings.HasPrefix(connStr, "redis://") || strings.HasPrefix(connStr, "rediss://")
}
